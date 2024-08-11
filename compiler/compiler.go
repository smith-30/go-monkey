package compiler

import (
	"fmt"

	"github.com/smith-30/go-monkey/ast"
	"github.com/smith-30/go-monkey/code"
	"github.com/smith-30/go-monkey/object"
)

// コンパイラが最後に実行した2つの命令（そのオペコードと実行された位置を含む）を追跡するように変更する。
type Compiler struct {
	instructions code.Instructions
	constants    []object.Object

	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
}

type EmittedInstruction struct {
	Opcode   code.Opcode
	Position int
}

func New() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants:    []object.Object{},
	}
}
func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, item := range node.Statements {
			err := c.Compile(item)
			if err != nil {
				return err
			}
		}
	case *ast.PrefixExpression:
		err := c.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "!":
			c.emit(code.OpBang)
		case "-":
			c.emit(code.OpMinus)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}

	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}
		// 式の評価が終わった段階で pop を呼び出して掃除する。
		c.emit(code.OpPop)
	case *ast.InfixExpression:
		// 右辺と左辺を逆転させることで、GreaterThan のみで対応可能にしている
		if node.Operator == "<" {
			err := c.Compile(node.Right)
			if err != nil {
				return err
			}
			err = c.Compile(node.Left)
			if err != nil {
				return err
			}
			c.emit(code.OpGreaterThan)
			return nil
		}

		err := c.Compile(node.Left)
		if err != nil {
			return err
		}
		err = c.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "+":
			c.emit(code.OpAdd)
		case "-":
			c.emit(code.OpSub)
		case "*":
			c.emit(code.OpMul)
		case "/":
			c.emit(code.OpDiv)
		case "==":
			c.emit(code.OpEqual)
		case "!=":
			c.emit(code.OpNotEqual)
		case ">":
			c.emit(code.OpGreaterThan)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}

	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		// vm が保持している constants の index を渡す。
		// vm はその index を利用し値を取り出す
		c.emit(code.OpConstant, c.addConstant(integer))
	case *ast.Boolean:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}

	case *ast.IfExpression:
		err := c.Compile(node.Condition)
		if err != nil {
			return err
		}

		// Emit an `OpJumpNotTruthy` with a bogus value
		// true ではない場合、どこにジャンプさせるかは後から決める
		// true 内部の評価をしないと false の場合にどこから処理が開始されるか決定できないため
		jumpNotTruthyPos := c.emit(code.OpJumpNotTruthy, 9999)

		// 内部式の評価
		err = c.Compile(node.Consequence)
		if err != nil {
			return err
		}

		// Consequence 評価後に pop で終わっている場合は削除する
		// なぜなら、IfExpression 評価後にPopは付与されるため重複してしまうから
		// let result = if (5 > 3) { 5 } else { 3 }; のような式に対応できない
		// ちなみに、上記の式のように値が未使用の単体の式は*ast.ExpressionStatementでラップされるから
		// Consequence 評価後に OpPop が付与されてしまう。これは Monkey の条件式が式であることによる。
		if c.lastInstructionIsPop() {
			c.removeLastPop()
		}

		// node.Alternative がない場合のみ、c.instruction の現在位置であるここにジャンプできる
		if node.Alternative == nil {
			//  バックパッチ, シングルパスコンパイラー
			// 内部式の評価が終わったのでjumpNotTruthyPos を書き換える
			afterConsequencePos := len(c.instructions)
			c.changeOperand(jumpNotTruthyPos, afterConsequencePos)
		} else {
			// Emit and `OpJump` with a bogus balue
			// else 式の評価後にどこに飛ぶか覚えていないといけない
			jumpPos := c.emit(code.OpJump, 9999)
			// Alternative がないときと同じように、式の評価は完了しているので
			// jumpNotTruthyPos を書き換える
			afterConsequencePos := len(c.instructions)
			c.changeOperand(jumpNotTruthyPos, afterConsequencePos)

			err := c.Compile(node.Alternative)
			if err != nil {
				return err
			}

			if c.lastInstructionIsPop() {
				c.removeLastPop()
			}

			// 式の評価後はどこに飛ばすかわかるので書き換える
			afterAlternativePos := len(c.instructions)
			c.changeOperand(jumpPos, afterAlternativePos)
		}

	// if { XXXXX; YYY; } など式内部のコンパイル
	case *ast.BlockStatement:
		for _, item := range node.Statements {
			err := c.Compile(item)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Compiler) lastInstructionIsPop() bool {
	return c.lastInstruction.Opcode == code.OpPop
}

func (c *Compiler) replaceInstruction(pos int, newInstruction []byte) {
	for i := 0; i < len(newInstruction); i++ {
		c.instructions[pos+i] = newInstruction[i]
	}
}

// 同じ型、同じ非変数の長さの命令だけを置き換える
func (c *Compiler) changeOperand(opPos int, operand int) {
	op := code.Opcode(c.instructions[opPos])
	newInstruction := code.Make(op, operand)
	c.replaceInstruction(opPos, newInstruction)
}

func (c *Compiler) removeLastPop() {
	c.instructions = c.instructions[:c.lastInstruction.Position]
	c.lastInstruction = c.previousInstruction
}

func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

// emit はコンパイラ用語で、"generate"（生成する）と "output"（出力する）を意味する
// 命令を生成し、それをprintしたり、ファイルに書き込んだり、
// メモリ上のコレクションに追加したりして、結果に追加する
func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := c.addInstruction(ins)

	c.setLastInstruction(op, pos)

	return pos
}

func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
	previous := c.lastInstruction
	last := EmittedInstruction{Opcode: op, Position: pos}
	c.previousInstruction = previous
	c.lastInstruction = last
}

func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.instructions)
	c.instructions = append(c.instructions, ins...)
	return posNewInstruction
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

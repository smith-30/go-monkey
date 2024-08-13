package compiler

// シンボルテーブルは、コンパイラやインタプリタにおいて、識別子とそれに関連する情報を結びつけるデータ構造です。
// レキシカル分析からコード生成までのすべての段階で使用され、特定の識別子（シンボルと呼ばれることもあります）に関する情報を保存および取得するために使用されます。
// その情報には、位置、スコープ、以前に宣言されたかどうか、関連付けられた値の型など
// 解釈やコンパイル中に有用と思われるものが含まれます。

// 私たちのシンボルテーブルは、識別子とスコープ、一意の番号を関連付けるために使用します。

// 現在の機能:

// 1. グローバルスコープ内の識別子を一意の番号と関連付ける
// 2. 指定された識別子に対して以前に関連付けられた番号を取得する

// これらの2つのメソッドは、一般的に「定義」と「解決」と呼ばれます。
// 特定のスコープ内で識別子を「定義」して、それに関連する情報を関連付けます。
// 後で、この情報に識別子を「解決」します。
// 情報自体を「シンボル」と呼びます。識別子はシンボルに関連付けられ、シンボル自体に情報が含まれます。

// 異なるスコープを区別するための識別子として使用される文字列型であること。
type SymbolScope string

const (
	GlobalScope SymbolScope = "GLOBAL"
)

// 識別子に関する情報を保持する構造体で、名前、スコープ、インデックスが含まれること。
type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

// 識別子とシンボルの対応関係を保持するデータ構造で、定義されたシンボルの数を追跡すること。
type SymbolTable struct {
	store          map[string]Symbol
	numDefinitions int
}

func NewSymbolTable() *SymbolTable {
	s := make(map[string]Symbol)
	return &SymbolTable{store: s}
}

func (s *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{Name: name, Index: s.numDefinitions, Scope: GlobalScope}
	s.store[name] = symbol
	s.numDefinitions++
	return symbol
}

func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	obj, ok := s.store[name]
	return obj, ok
}

package protocol

// Parser interface
type Parser interface {
	GetToken(data []byte) (interface{}, error)
}

type FuncParserHandler func(data []byte) (interface{}, error)
type funcParser struct{
	f FuncParserHandler
}
func NewFuncParser(f FuncParserHandler) Parser {
	return &funcParser{f: f}
}

func (it *funcParser) GetToken(data []byte) (interface{}, error) {
	return it.f(data)
}

package protocol

// Parser interface
type Parser interface {
	GetToken(data []byte) (interface{}, error)
}

type MultiplexParser interface {
	GetChannel(data []byte) (byte, error)
}

type TokenParser interface {
	GetToken(data []byte) (uint32, error)
}

type FuncMultiplexParser func(data []byte) (byte, error)
type funcMultiplexParser struct{
	f FuncMultiplexParser
}
func NewFuncMultiplexParser(f FuncMultiplexParser) MultiplexParser {{
	return &funcMultiplexParser{f: f}
}}
func (it *funcMultiplexParser) GetChannel(data []byte) (byte, error) {
	return it.f(data)
}


type FuncTokenParser func(data []byte) (uint32, error)
type funcTokenParser struct {
	f FuncTokenParser
}
func NewFuncTokenParser(f FuncTokenParser) TokenParser {
	return &funcTokenParser{f: f}
}
func (it *funcTokenParser) GetToken(data []byte) (uint32, error) {
	return it.f(data)
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

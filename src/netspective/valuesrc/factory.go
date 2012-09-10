package valuesrc

import (
	"strings"
	"fmt"
)

const vsDelim = ':'

type Value interface {
	// use type assertions to figure out the real type (http://golang.org/ref/spec#Type_assertions)
	GetValue() interface{}
	HasValue() bool
}

type ValueContext interface {
}

type ValueSource interface {
	EvaluateValue(vc ValueContext) Value
}

type Handler struct {
	name string
	construct func(string) (ValueSource, error)
}

type Handlers struct {
	handlers []Handler
	handlersByName map[string] Handler
}

func (h Handlers) RegisterHandler(name string, construct func(string) (ValueSource, error)) {
	handler := Handler{name, construct}
	h.handlers = append(h.handlers, handler)
	h.handlersByName[name] = handler
}

func (h Handlers) CreateValueSource(spec string, construct func(string) (ValueSource, error)) (ValueSource, error) {
	delimPos := strings.IndexRune(spec, vsDelim)
	if(delimPos > 0 && spec[delimPos] != '\\') {
		name := spec[:delimPos-1]
		handler, ok := h.handlersByName[name]
		if(! ok) {
			return nil, fmt.Errorf("ValueSource handler not found for '%s' in '%s'", name, spec)
		}

		params := spec[delimPos+1:]
		return handler.construct(params);
	}

	if(construct != nil) {
		return construct(spec)
	}
	return nil, nil
}

var Factory Handlers

func init() {
	Factory.handlersByName = make(map[string] Handler)
	Factory.RegisterHandler("text", NewTextLiteral)
	Factory.RegisterHandler("literal", NewLiteral)
}

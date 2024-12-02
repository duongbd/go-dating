package decode

import (
	"context"

	"github.com/labstack/echo/v4"
)

type Decoder interface {
	Decode(v interface{}) error
}

type Validator interface {
	Validate(v interface{}) error
}

// DecodeAndValidateRequest decodes and validates a request
func DecodeAndValidateRequest[T any](c context.Context, rq *T, decoder Decoder, validator Validator) error {
	if err := decoder.Decode(rq); err != nil {
		return err
	}

	if err := validator.Validate(rq); err != nil {
		return err

	}

	return nil
}

// EchoDecoder is a decoder for echo.Context as we use echo for this challenge
type EchoDecoder struct {
	C echo.Context
}

func (e *EchoDecoder) Decode(v interface{}) error {
	return e.C.Bind(v)
}

type EchoValidator struct {
	C echo.Context
}

func (e *EchoValidator) Validate(v interface{}) error {
	return e.C.Validate(v)
}

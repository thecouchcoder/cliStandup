package models

import "context"

type LLM interface {
	Generate(context.Context) (string, error)
}

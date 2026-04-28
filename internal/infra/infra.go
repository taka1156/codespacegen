package infra

import (
	"codespacegen/internal/infra/prompt"
	"io"
)

type CodespacePrompter = prompt.CodespacegenPrompter

func NewCodespacePrompter(reader io.Reader) *CodespacePrompter {
	return prompt.NewCodespacegenPrompter(reader)
}

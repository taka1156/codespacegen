package infra

import (
	"io"

	"github.com/taka1156/codespacegen/internal/infra/prompt"
)

type CodespacePrompter = prompt.CodespacegenPrompter

func NewCodespacePrompter(reader io.Reader) *CodespacePrompter {
	return prompt.NewCodespacegenPrompter(reader)
}

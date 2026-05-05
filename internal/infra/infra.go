package infra

import (
	"io"

	"github.com/taka1156/codespacegen/internal/infra/prompt"
	"github.com/taka1156/codespacegen/internal/infra/updater"
)

type CodespacegenPrompter = prompt.CodespacegenPrompter

func NewCodespacegenPrompter(reader io.Reader) *CodespacegenPrompter {
	return prompt.NewCodespacegenPrompter(reader)
}

type CodespacegenUpdater = updater.CodespacegenUpdater

func NewCodespacegenUpdater() *CodespacegenUpdater {
	return updater.NewCodespacegenUpdater()
}

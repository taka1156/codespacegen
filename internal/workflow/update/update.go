package update

type CodespacegenUpdater interface {
	Update(currentVersion string) error
}

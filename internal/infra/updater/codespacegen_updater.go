package updater

import (
	"fmt"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/taka1156/codespacegen/internal/domain/entity"
	"github.com/taka1156/codespacegen/internal/i18n"
)

type CodespacegenUpdater struct{}

func NewCodespacegenUpdater() *CodespacegenUpdater {
	return &CodespacegenUpdater{}
}

func (cu *CodespacegenUpdater) Update(currentVersion string) error {
	parsedCurrentVersion, err := semver.ParseTolerant(currentVersion)
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T("error_failed_to_parse_current_version"), err)
	}

	latest, err := selfupdate.UpdateSelf(parsedCurrentVersion, entity.DefaultRepositoryName)
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T("error_failed_to_check_latest_version"), err)
	}

	if latest.Version.Equals(parsedCurrentVersion) {
		fmt.Println(i18n.T("already_latest_version"))
		return nil
	}

	fmt.Print(
		i18n.T("success_update", map[string]interface{}{"Version": latest.Version, "ReleaseNotes": latest.ReleaseNotes}),
	)

	return nil
}

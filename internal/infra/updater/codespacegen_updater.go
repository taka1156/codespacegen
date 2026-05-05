package updater

import (
	"errors"
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
	parsedCurrentVersion, err := semver.Parse(currentVersion)
	if err != nil {
		return errors.New(i18n.T("error_failed_to_parse_current_version", map[string]interface{}{"Error": err.Error()}))
	}

	latest, err := selfupdate.UpdateSelf(parsedCurrentVersion, entity.DefaultRepositoryName)
	if err != nil {
		return errors.New(i18n.T("error_failed_to_check_latest_version", map[string]interface{}{"Error": err.Error()}))
	}

	fmt.Print(i18n.T("success_update", map[string]interface{}{"Version": latest.Version, "ReleaseNotes": latest.ReleaseNotes}))

	return nil
}

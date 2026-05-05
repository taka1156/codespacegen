package update

type UpdateCommandline struct {
	updateCodespacegenCommandline CodespacegenUpdater
}

func NewUpdateCommandline(
	updateCodespacegenCommandline CodespacegenUpdater,
) *UpdateCommandline {
	return &UpdateCommandline{
		updateCodespacegenCommandline: updateCodespacegenCommandline,
	}
}

func (uc *UpdateCommandline) Update(currentVersion string) error {
	err := uc.updateCodespacegenCommandline.Update(currentVersion)
	if err != nil {
		return err
	}

	return nil
}

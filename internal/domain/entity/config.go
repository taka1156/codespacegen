package entity

// CodespaceConfig holds values used to generate devcontainer artifacts.
type CodespaceConfig struct {
	ContainerName   string
	ServiceName     string
	WorkspaceFolder string
	BaseImage       string
	ComposeFileName string
	PortMapping     string
	InstallCommand  string
}

func (c CodespaceConfig) Validate() error {
	if c.ContainerName == "" {
		return ErrInvalidConfig("container name is required")
	}
	if c.ServiceName == "" {
		return ErrInvalidConfig("service name is required")
	}
	if c.WorkspaceFolder == "" {
		return ErrInvalidConfig("workspace folder is required")
	}
	if c.BaseImage == "" {
		return ErrInvalidConfig("base image is required")
	}
	if c.ComposeFileName == "" {
		return ErrInvalidConfig("compose file name is required")
	}
	return nil
}

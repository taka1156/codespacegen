package workdirprovider

import "os"

type WorkdirProvider struct{}

func NewWorkdirProvider() *WorkdirProvider {
	return &WorkdirProvider{}
}

func (w *WorkdirProvider) GetConfigOutputPath() (string, error) {
	return os.Getwd()
}

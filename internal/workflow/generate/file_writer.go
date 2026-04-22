package generate

type FileWriter interface {
	Write(path string, content string, overwrite bool) error
}

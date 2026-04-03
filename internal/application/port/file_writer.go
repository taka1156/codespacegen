package port

type FileWriter interface {
	Write(relativePath string, content string, overwrite bool) error
}

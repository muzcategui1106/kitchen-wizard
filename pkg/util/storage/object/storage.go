package object

type Storage interface {
	Save(string, []byte) error
	Delete(string) error
}

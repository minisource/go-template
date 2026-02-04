package dto

type IdName struct {
	Id   int
	Name string
}
type Name struct {
	Name string
}

type CreateFile struct {
	Name        string
	Directory   string
	Description string
	MimeType    string
}

type UpdateFile struct {
	Description string
}

type File struct {
	IdName
	Directory   string
	Description string
	MimeType    string
}

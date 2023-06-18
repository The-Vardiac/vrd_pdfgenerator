package repository

type Generate struct {
	Filename string
	Text     string
}

type GenerateFileEmailTemplate struct {
	UserName string
	Filename string
	FileUrl  string
}
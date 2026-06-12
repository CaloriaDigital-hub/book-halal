package commands



type UploadBookCommand struct {
	Title       string
	Author      string
	Description string
	Price       int
	TmpPDFPath  string
}


type UploadBookResult struct {
	BookID string
}
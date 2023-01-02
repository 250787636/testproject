package hecmcc

type FtpUploader struct {
	host     string
	username string
	password string
	dir      string
	isUpload bool
}

func NewFtpUploader(host, username, password, dir string, isUpload bool) *FtpUploader {
	uploader := new(FtpUploader)
	uploader.host = host
	uploader.username = username
	uploader.password = password
	uploader.dir = dir
	return uploader
}

func (f *FtpUploader) Upload(filePaths []string) error {
	if !f.isUpload {
		return nil
	}

	return nil
}

package models

import (
	"github.com/astaxie/beego"
	"io"
	"os"
)

func CopyFile(src, des string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		beego.Error(err)
	}
	defer srcFile.Close()

	desFile, err := os.Create(des)
	if err != nil {
		beego.Error(err)
	}
	defer desFile.Close()

	return io.Copy(desFile, srcFile)
}

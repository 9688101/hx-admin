package utils

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/gin-contrib/static"
)

// Credit: https://github.com/gin-contrib/static/issues/19

type embedFileSystem struct {
	http.FileSystem
}

func (e embedFileSystem) Exists(prefix string, path string) bool {
	_, err := e.Open(path)
	return err == nil
}

func EmbedFolder(fsEmbed embed.FS, targetPath string) static.ServeFileSystem {
	efs, err := fs.Sub(fsEmbed, targetPath)
	if err != nil {
		panic(err)
	}
	_, err = fsEmbed.Open(targetPath)
	if err != nil {
		fmt.Println("Failed to open targetPath:", err)
	} else {
		fmt.Println("Success!")
	}
	files, _ := fsEmbed.ReadDir(".")
	for _, file := range files {
		fmt.Println(1, file.Name())
	}
	return embedFileSystem{
		FileSystem: http.FS(efs),
	}
}

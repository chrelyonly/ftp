package main

import (
	"log"

	"goftp.io/server/core"
	"goftp.io/server/driver/file"
)

func main() {
	Name := "ftp"
	rootPath := "./test" //FTP根目录
	Port := 2121         //FTP 端口
	var perm = core.NewSimplePerm("cjh", "test")

	// Server options without hostname or port
	opt := &core.ServerOpts{
		Name: Name,
		Factory: &file.DriverFactory{
			RootPath: rootPath,
			Perm:     perm,
		},
		Auth: &core.SimpleAuth{
			Name:     "cjh",    // FTP 账号
			Password: "123456", // FTP 密码
		},
		Port:     Port,
		Hostname: "10.80.31.136",
	}
	// start ftp server
	s := core.NewServer(opt)
	err := s.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}

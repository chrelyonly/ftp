package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/jlaffaye/ftp"
)

func main() {
	// Do something with the FTP conn

	//for {
	//	newFunction()
	//	time.Sleep(3 * time.Second)
	//}

	//err := DownloadFtpFile("10.80.31.136", "cjh", "123456", "/", "info.json")
	//
	//if err != nil {
	//	fmt.Println("err=", err)
	//} else {
	//	fmt.Println("下载成功")
	//}

	DownloadDir("10.80.31.136", "cjh", "123456", "test_copy", "/")

	fmt.Println("下载成功")

}

// newFunction 上传
func newFunction() {
	file, err := os.Open("./info.json")
	if err != nil {
		log.Println("读取文件:")
		log.Println(err)
		return
	}

	c, err := ftp.Dial("10.80.31.136:2121", ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		log.Println("登录0:")
		log.Println(err)
		return
	}

	err = c.Login("cjh", "123456")
	if err != nil {
		log.Println("登录1:")
		log.Println(err)
		return
	}

	if err := c.Stor(`./info.json`, file); err != nil {
		log.Println("上传:")
		log.Println(err)
		return
	} else {
		log.Println("上传成功")
		os.Exit(0)
	}

}

// DownloadFtpFile 下载
func DownloadFtpFile(host, username, password string, path, fileName string) error {
	// 建立连接，默认用21端口
	c, err := ftp.Dial(host+":2121", ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return err
	}
	defer c.Quit()

	// 登录
	err = c.Login(username, password)
	if err != nil {
		return err
	}

	// 切换到path
	err = c.ChangeDir(path)
	if err != nil {
		return err
	}

	// 读取文件
	body, err := c.Retr(fileName)
	if err != nil {
		return err
	}
	defer body.Close()

	// 创建本地文件
	localFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer localFile.Close()

	// 下载到本地
	_, err = io.Copy(localFile, body)
	return nil
}

// DownloadDir 目录下载
func DownloadDir(host, username, password string, local string, remote string) {

	// 建立连接，默认用21端口
	conn, err := ftp.Dial(host+":2121", ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return
	}
	defer conn.Quit()

	// 登录
	err = conn.Login(username, password)
	if err != nil {
		return
	}
	_ = os.Mkdir(local, 0664)

	entries, err := conn.List(remote)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {

		if entry.Type == ftp.EntryTypeFile {
			downloadFile(conn, local+"//"+entry.Name, remote+"//"+entry.Name)
		}

		if entry.Type == ftp.EntryTypeFolder {
			if entry.Name != "." {
				DownloadDir(host, username, password, local+"//"+entry.Name, remote+"//"+entry.Name)
			}
		}
	}
}

func downloadFile(conn *ftp.ServerConn, local string, remote string) {

	res, err := conn.Retr(remote)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Close()

	outFile, err := os.Create(local)
	if err != nil {
		log.Fatal(err)
	}

	defer outFile.Close()

	_, err = io.Copy(outFile, res)
	if err != nil {
		log.Fatal(err)
	}
}

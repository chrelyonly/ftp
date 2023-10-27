package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"goftp.io/server/core"
	"goftp.io/server/driver/file"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Person struct {
	UserName  string `json:"userName"`
	Password  string `json:"password"`
	ApiProt   string `json:"apiProt"`
	Host      string `json:"host"`
	UploadDir string `json:"uploadDir"`
}

func main() {
	var arr = os.Args
	if len(arr) > 1 && arr[1] == "--version" {
		fmt.Println(color.GreenString("当前版本1.0"))
		return
	}
	fmt.Println(color.GreenString("简易ftp工具,支持全平台单文件,by_chrelyonly"))
	//打印当前时间
	currentTime := time.Now()
	fmt.Println(color.GreenString("当前时间: " + currentTime.Format("2006-01-02 15:04:05")))
	configPath := filepath.Base("config.json")
	//配置宝塔插件
	//configPath := "/www/server/panel/plugin/ftp_chrelyonly/config.json"
	//判断文件是否存在
	_, err := os.Stat(configPath)
	//配置信息
	var userName string
	var password string
	var apiProt int
	var host string
	var uploadDir string
	//文件对象
	var config *os.File
	if err != nil {
		fmt.Println(color.RedString("当前配置文件不存在,将进行初始化配置"))
		//创建文件
		config, err = os.Create(configPath)
		if err != nil {
			fmt.Println(color.RedString("配置文件创建失败,请检查是否有权限"))
			os.Exit(1)
		}
		fmt.Println(color.GreenString("请输入管理员用户名(默认: ftpadmin))"))
		_, err := fmt.Scanln(&userName)
		if err != nil {
			userName = "ftpadmin"
			fmt.Println(color.BlueString("当前用户名: ftpadmin"))
		}
		fmt.Println(color.GreenString("请输入管理员密码(默认: ftpadmin)"))
		_, err = fmt.Scanln(&password)
		if err != nil {
			password = "ftpadmin"
			fmt.Println(color.BlueString("当前密码: ftpadmin"))
		}
		fmt.Println(color.GreenString("请输入绑定ip: 默认127.0.0.1"))
		_, err = fmt.Scanln(&apiProt)
		if err != nil {
			host = "127.0.0.1"
			fmt.Println(color.BlueString("绑定ip: 127.0.0.1"))
		}
		fmt.Println(color.GreenString("请输入API接口端口(默认: 30000)"))
		_, err = fmt.Scanln(&apiProt)
		if err != nil {
			apiProt = 30000
			fmt.Println(color.BlueString("当前API端口: 30000"))
		}
		fmt.Println(color.GreenString("请输入映射位置(   /www/  || ./  || ....   ),默认./"))
		_, err = fmt.Scanln(&uploadDir)
		if err != nil {
			uploadDir = "./"
			fmt.Println(color.BlueString("当前映射位置: ./"))
		}
		//将配置信息保存
		person := Person{
			UserName:  userName,
			Password:  password,
			ApiProt:   strconv.Itoa(apiProt),
			UploadDir: uploadDir,
			Host:      host,
		}
		//将配置信息写入文件
		encoder := json.NewEncoder(config)
		err = encoder.Encode(person)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		//	如果存在文件则,从文件中读取配置
		fmt.Println(color.GreenString("当前配置文件存在,将从配置文件中读取配置"))
		//读取配置文件,json格式
		config, err = os.OpenFile(configPath, os.O_RDWR, 0666)
		if err != nil {
			fmt.Println(color.RedString("配置文件读取失败,请检查配置文件是否存在"))
			os.Exit(1)
		}
		// 从文件中读取 JSON 数据
		decoder := json.NewDecoder(config)
		var person Person
		err = decoder.Decode(&person)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		userName = person.UserName
		password = person.Password
		apiProt, _ = strconv.Atoi(person.ApiProt)
		uploadDir = person.UploadDir
		host = person.Host
	}

	//关闭文件
	err = config.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	Name := "简易ftp工具by_chrelyonly"
	//FTP根目录
	rootPath := uploadDir
	//FTP 端口
	Port := apiProt
	var perm = core.NewSimplePerm("chrelyonly", "test")

	// Server options without hostname or port
	opt := &core.ServerOpts{
		Name: Name,
		Factory: &file.DriverFactory{
			RootPath: rootPath,
			Perm:     perm,
		},
		Auth: &core.SimpleAuth{
			Name:     userName, // FTP 账号
			Password: password, // FTP 密码
		},
		Port:     Port,
		Hostname: host,
	}
	// start ftp server
	s := core.NewServer(opt)
	err = s.ListenAndServe()
	if err != nil {
		log.Fatal("服务启动失败:", err)
	}
}

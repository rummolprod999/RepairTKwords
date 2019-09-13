package main

import (
	"fmt"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Filelog string

var DirLog = "log_repair"
var DirTemp = "temp_repair"
var SetFile = "settings.json"
var FileLog Filelog
var mutex sync.Mutex
var Dsn string
var Server string
var Port int64
var Prefix = ""
var Typefz = 0

func CreateLogFile() {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dirlog := filepath.FromSlash(fmt.Sprintf("%s/%s", dir, DirLog))
	if _, err := os.Stat(dirlog); os.IsNotExist(err) {
		err := os.MkdirAll(dirlog, 0711)

		if err != nil {
			fmt.Println("Не могу создать папку для лога")
			os.Exit(1)
		}
	}
	t := time.Now()
	ft := t.Format("2006-01-02")
	FileLog = Filelog(filepath.FromSlash(fmt.Sprintf("%s/log_repair_%v.log", dirlog, ft)))
}

func CreateTempDir() {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dirtemp := filepath.FromSlash(fmt.Sprintf("%s/%s", dir, DirTemp))
	if _, err := os.Stat(dirtemp); os.IsNotExist(err) {
		err := os.MkdirAll(dirtemp, 0711)

		if err != nil {
			fmt.Println("Не могу создать папку для временных файлов")
			os.Exit(1)
		}
	} else {
		err = os.RemoveAll(dirtemp)
		if err != nil {
			fmt.Println("Не могу удалить папку для временных файлов")
		}
		err := os.MkdirAll(dirtemp, 0711)
		if err != nil {
			fmt.Println("Не могу создать папку для временных файлов")
			os.Exit(1)
		}
	}
}

func Logging(args ...interface{}) {
	mutex.Lock()
	file, err := os.OpenFile(string(FileLog), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	defer file.Close()
	if err != nil {
		fmt.Println("Ошибка записи в файл лога", err)
		return
	}
	fmt.Fprintf(file, "%v  ", time.Now())
	for _, v := range args {

		fmt.Fprintf(file, " %v", v)
	}
	//fmt.Fprintf(file, " %s", UrlXml)
	fmt.Fprintln(file, "")
	mutex.Unlock()
}

func ReadSetting() {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	filetemp := filepath.FromSlash(fmt.Sprintf("%s/%s", dir, SetFile))
	file, err := os.Open(filetemp)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	UserDb, err := jsonparser.GetString(b, "userdb")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if UserDb == "" {
		fmt.Println("Check file with settings")
		os.Exit(1)
	}
	PassDb, err := jsonparser.GetString(b, "passdb")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if PassDb == "" {
		fmt.Println("Check file with settings")
		os.Exit(1)
	}

	DbName, err := jsonparser.GetString(b, "db")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if DbName == "" {
		fmt.Println("Check file with settings")
		os.Exit(1)
	}

	Server, err = jsonparser.GetString(b, "server")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if Server == "" {
		fmt.Println("Check file with settings")
		os.Exit(1)
	}

	Port, err = jsonparser.GetInt(b, "port")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if Port == 0 {
		fmt.Println("Check file with settings")
		os.Exit(1)
	}
	Dsn = fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=true&readTimeout=60m&maxAllowedPacket=0&timeout=60m&writeTimeout=60m&autocommit=true&loc=Local", UserDb, PassDb, DbName)
}

func CreateEnv() {
	ReadSetting()
	CreateLogFile()
	CreateTempDir()
}

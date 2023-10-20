package conf

import (
    "fmt"
    "github.com/go-ini/ini"
    "os"
    "strconv"
    "sync"
)

var CntMutex sync.Mutex

var BlockedCom string
var SelectedType string
var PassedCnt int

var MesUrl string

func LoadConf() {
    fmt.Println("load")
    cfg, err := ini.Load("config.ini")
    if err != nil {
        file, err := os.Create("config.ini")
        if err != nil {
            fmt.Println("create file error, ", err)
            return
        }
        defer file.Close()
        file.WriteString("[runinfo]\n")
        file.WriteString("BlockedCom=\n")
        file.WriteString("SelectedType=\n")
        file.WriteString("PassedCnt=\n")
        cfg, err = ini.Load("config.ini")
        if err != nil {
            fmt.Println("get config failed, ", err)
            return
        }
    }
    section := cfg.Section("runinfo")
    BlockedCom = section.Key("BlockedCom").String()
    SelectedType = section.Key("SelectedType").String()
    PassedCnt, _ = section.Key("PassedCnt").Int()

    section = cfg.Section("mes")
    MesUrl = section.Key("Url").String()
}

func SyncConf() {
    cfg, err := ini.Load("config.ini")
    if err != nil {
        fmt.Println("load config failed, ", err)
        return
    }
    section := cfg.Section("runinfo")
    section.Key("BlockedCom").SetValue(BlockedCom)
    section.Key("SelectedType").SetValue(SelectedType)
    section.Key("PassedCnt").SetValue(strconv.Itoa(PassedCnt))
    err = cfg.SaveTo("config.ini")
    if err != nil {
        fmt.Println("save config failed, ", err)
    }
}

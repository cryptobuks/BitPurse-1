package test

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/astaxie/beego"
	"time"
	"fmt"
)

func init() {
	_, file, _, _ := runtime.Caller(1)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)
}

func TestConfig(t *testing.T) {
	ds := beego.AppConfig.String("dataSource")
	if len(ds) == 0 {
		t.Error("no data source")
	}
}

func TestTicker(t *testing.T) {
	ticker := time.NewTicker(time.Second)

	go func() {
		for t := range ticker.C {
			fmt.Println("Tick at", t)
		}
	}()
	time.Sleep(time.Second * 1600)
	ticker.Stop()
	fmt.Println("Ticker stopped")
}

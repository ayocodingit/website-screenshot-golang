package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/robfig/cron/v3"
)

func RemoveFiles() {
	fmt.Println("remove all public")
	dir, _ := ioutil.ReadDir("public")
	for _, d := range dir {
		if d.Name() != "readme" {
			os.RemoveAll(path.Join([]string{"public", d.Name()}...))
		}
	}
}

func InitCron() {
	scheduler := cron.New()
	defer scheduler.Stop()
	scheduler.AddFunc("@daily", func() { RemoveFiles() })
	scheduler.Start()
}

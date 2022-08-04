package utils

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

func GetScreenshot(screenShotUrl string, quality int) (filename string, err error) {

	var buf []byte

	var ext string = "png"
	if quality < 100 {
		ext = "jpeg"
	}

	var options []chromedp.ExecAllocatorOption
	options = append(options, chromedp.WindowSize(1280, 1280))
	options = append(options, chromedp.DefaultExecAllocatorOptions[:]...)

	actx, acancel := chromedp.NewExecAllocator(context.Background(), options...)

	defer acancel()

	ctx, cancel := chromedp.NewContext(actx)
	defer cancel()

	tasks := chromedp.Tasks{
		chromedp.Navigate(screenShotUrl),
		chromedp.CaptureScreenshot(&buf),
	}

	if err = chromedp.Run(ctx, tasks); err != nil {
		return
	}

	path := "public"

	filename = fmt.Sprintf("%d.%s", time.Now().UTC().Unix(), ext)

	filename = path + "/" + filename

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
	if err = ioutil.WriteFile(filename, buf, 0644); err != nil {
		return
	}

	return
}

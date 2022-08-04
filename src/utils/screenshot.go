package utils

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/chromedp/chromedp"
)

func GetScreenshot(screenShotUrl string, quality int) (filename string, err error) {
	var buf []byte
	var options []chromedp.ExecAllocatorOption

	options = append(options, chromedp.WindowSize(1400, 900))
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

	filename = fmt.Sprintf("%d.%s", time.Now().UTC().Unix(), "jpeg")
	filename = "public" + "/" + filename

	if err = ioutil.WriteFile(filename, buf, 0644); err != nil {
		return
	}

	return
}

package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	cron "github.com/robfig/cron/v3"
)

type Screenshot struct {
	URL     string `json:"url"`
	Quality int    `json:"quality"`
}

func RemoveFiles() {
	fmt.Println("remove all public")
	dir, _ := ioutil.ReadDir("public")
	for _, d := range dir {
		if d.Name() != "readme" {
			os.RemoveAll(path.Join([]string{"public", d.Name()}...))
		}
	}
}

func initCron() {
	scheduler := cron.New()
	defer scheduler.Stop()
	scheduler.AddFunc("@daily", func() { RemoveFiles() })
	go scheduler.Start()
}

func screenshotHandler(c *gin.Context) {
	var screenshot Screenshot
	if err := c.ShouldBindJSON(&screenshot); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"msg": err.Error()})
		return
	}
	url := location.Get(c)

	filename, err := getChromedpScreenShot(screenshot.URL, screenshot.Quality)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	res := url.String() + "/" + filename
	c.JSON(http.StatusOK, res)
}

func main() {
	initCron()

	r := gin.Default()
	r.Static("/public", "./public")
	r.Use(location.Default())
	r.Use(cors.Default())
	r.POST("/", screenshotHandler)

	r.Run()
}

func getChromedpScreenShot(site string, quality int) (filename string, err error) {
	//forming url to be captured
	screenShotUrl := fmt.Sprintf("https://%s/", site)

	//byte slice to hold captured image in bytes
	var buf []byte

	//setting image file extension to png but
	var ext string = "png"
	//if image quality is less than 100 file extension is jpeg
	if quality < 100 {
		ext = "jpeg"
	}

	//setting options for headless chrome to execute with
	var options []chromedp.ExecAllocatorOption
	options = append(options, chromedp.WindowSize(1400, 900))
	options = append(options, chromedp.DefaultExecAllocatorOptions[:]...)

	//setup context with options
	actx, acancel := chromedp.NewExecAllocator(context.Background(), options...)

	defer acancel()

	// create context
	ctx, cancel := chromedp.NewContext(actx)
	defer cancel()

	//configuring a set of tasks to be run
	tasks := chromedp.Tasks{
		//loads page of the URL
		chromedp.Navigate(screenShotUrl),
		//waits for 5 secs
		chromedp.Sleep(5 * time.Second),
		//Captures Screenshot with current window size
		chromedp.CaptureScreenshot(&buf),
		//captures full-page screenshot (uncomment to take fullpage screenshot)
		//chromedp.FullScreenshot(&buf,quality),
	}

	// running the tasks configured earlier and logging any errors
	if err = chromedp.Run(ctx, tasks); err != nil {
		return
	}
	//naming file using provided URL without "/"s and current unix datetime
	filename = fmt.Sprintf("%s-%d-standard.%s", strings.Replace(site, "/", "-", -1), time.Now().UTC().Unix(), ext)

	filename = "public/" + filename

	path := "public"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
	//write byte slice data of standard screenshot to file
	if err = ioutil.WriteFile(filename, buf, 0644); err != nil {
		return
	}

	return
}

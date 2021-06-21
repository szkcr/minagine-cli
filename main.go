package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/sclevine/agouti"
)

func newChromeDriver() *agouti.WebDriver {
	options := agouti.ChromeOptions(
		"args", []string{
			"--headless",
			"--window-size=1280,720",
			"--blink-settings=imagesEnabled=false",
			"--disable-gpu",
			"no-sandbox",
			"disable-dev-shm-usage",
		})
	return agouti.ChromeDriver(options)
}

type opts struct {
	Domain     string  `required:"true" short:"d" long:"domain" description:"domain of your account (tenant)"`
	User       string  `required:"true" short:"u" long:"user" description:"user id of your account"`
	Password   string  `required:"true" short:"p" long:"password" description:"password of your account"`
	Action     string  `required:"true" short:"a" long:"action" choice:"checkin" choice:"checkout" description:"action to be performed"`
	Force      bool    `required:"false" short:"f" long:"force" description:"[option] skip pre-check of working day"`
	WebhookURL *string `required:"false" short:"w" long:"webhook" description:"[option] webhook url for reporting action result"`
}

func exit(code int, msg string, webhookURL *string) {
	log.Println(msg)
	os.Exit(code)
}

func main() {
	opts := opts{}
	if _, err := flags.Parse(&opts); err != nil {
		exit(1, "invalid argument", nil)
	}
	domain := opts.Domain
	user := opts.User
	pw := opts.Password
	checkin := opts.Action == "checkin"

	// ブラウザ初期化
	driver := newChromeDriver()
	defer func() { _ = driver.Stop() }()
	if err := driver.Start(); err != nil {
		exit(1, fmt.Sprintf("failed to start browser: %v", err), opts.WebhookURL)
	}
	page, err := driver.NewPage()
	if err != nil {
		exit(1, fmt.Sprintf("failed to create browser page: %v", err), opts.WebhookURL)
	}

	// 出勤予定日かチェックする
	if !opts.Force {
		today := time.Now()
		isWorkingDay, status, err := isWorkingDay(page, domain, user, pw, today)
		if err != nil {
			exit(1, fmt.Sprintf("failed to check working day: %v", err), opts.WebhookURL)
		}
		if !isWorkingDay {
			exit(0, fmt.Sprintf("action was skipped (not working day: %s = %s)", today.Format("2006-01-02"), status), opts.WebhookURL)
		}
	}

	// アクション
	latestAction, err := dakoku(page, domain, user, pw, checkin)
	if err != nil {
		exit(1, fmt.Sprintf("failed to dakoku: %v)", err), opts.WebhookURL)
	}

	exit(0, fmt.Sprintf("dakoku was successful => %s", latestAction), opts.WebhookURL)
}

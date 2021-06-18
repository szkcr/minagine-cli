package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/sclevine/agouti"
)

const (
	MINAGINE_LOGIN_URL  = "https://tm.minagine.net/index.html"
	MINAGINE_SHIFT_URL  = "https://tm.minagine.net/work/wrktimemngmntshtself/sht"
	MINAGINE_DAKOKU_URL = "https://tm.minagine.net/work/wrktimergst"
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
	driver := agouti.ChromeDriver(options)
	return driver
}

func openPage(page *agouti.Page, targetURL string, domain string, user string, pw string) error {
	page.Navigate(targetURL)
	url, err := page.URL()
	if err != nil {
		return err
	}

	// ログイン画面へリダイレクトされたらログインする
	if url == MINAGINE_LOGIN_URL {
		if err := page.FindByID("user_cntrctr_dmn").Fill(domain); err != nil {
			return err
		}
		if err := page.FindByID("user_login").Fill(user); err != nil {
			return err
		}
		if err := page.FindByID("user_password").Fill(pw); err != nil {
			return err
		}
		if err := page.FindByID("login_form").Submit(); err != nil {
			return err
		}

		page.Navigate(targetURL)
		url, err = page.URL()
		if err != nil {
			return err
		}
	}

	if url != targetURL {
		return fmt.Errorf("failed to navigate targetURL: `%s`", targetURL)
	}

	// ログイン済みのページが正しく開けているか確認する
	if err := page.FindByXPath(`//*[@id="show_popup_info"]`).MouseToElement(); err != nil {
		return err
	}
	text, err := page.FindByXPath(`//*[@id="login_information"]/ul/li[1]/div[1]`).Text()
	if err != nil {
		return err
	}
	renderedID := strings.TrimSpace(strings.SplitN(text, `|`, 2)[1])
	if renderedID != user {
		return fmt.Errorf("loginID not matched (expected:`%s`, actual:`%s`)", user, renderedID)
	}

	return nil
}

func dakoku(page *agouti.Page, domain string, user string, pw string, starting bool) (string, error) {
	err := openPage(page, MINAGINE_DAKOKU_URL, domain, user, pw)
	if err != nil {
		return "", err
	}

	var button *agouti.Selection
	if starting {
		button = page.FindByID("button0")
	} else {
		button = page.FindByID("button1")
	}
	if err := button.Click(); err != nil {
		return "", err
	}

	flashText, err := page.FindByXPath(`//*[@id="_flash"]`).Text()
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(flashText) != "打刻しました。" {
		return "", fmt.Errorf("failed to confirm dakoku was succeeded")
	}

	latestAction, err := page.FindByXPath(`//*[@id="new_cndtn"]/table[3]/tbody/tr[1]`).Text()
	if err != nil {
		return "", err
	}

	return strings.ReplaceAll(strings.TrimSpace(latestAction), "\t", " "), nil
}

func isWorkingDay(page *agouti.Page, domain string, user string, pw string, today time.Time) (bool, error) {
	err := openPage(page, MINAGINE_SHIFT_URL, domain, user, pw)
	if err != nil {
		return false, err
	}

	rows := page.AllByXPath(`//*[@id="table_wrktimesht"]/tbody/tr`)
	numOfRow, err := rows.Count()
	if err != nil {
		return false, err
	}

	// 最初の4行はヘッダのためスキップ
	for i := 4; i < numOfRow; i++ {
		// `日`列をチェックして対象日との一致をチェック
		day, err := rows.At(i).FindByXPath(`td[1]`).Text()
		if err != nil {
			return false, err
		}
		if strconv.Itoa(today.Day()) == strings.TrimSpace(day) {
			// 出勤予定(`○`)かどうかをチェック
			status, err := rows.At(i).FindByXPath(`td[3]`).Text()
			if err != nil {
				return false, err
			}
			return strings.TrimSpace(status) == "○", nil
		}
	}

	return false, fmt.Errorf("could not find date in shift table")
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 4 {
		log.Fatalln("invalid len(args)")
		return
	}
	domain, user, pw := args[0], args[1], args[2]
	startOrEnd := args[3]
	switch startOrEnd {
	case "start", "end":
	default:
		log.Fatalln("invalid args[3]")
		return
	}

	driver := newChromeDriver()
	defer func() { _ = driver.Stop() }()
	if err := driver.Start(); err != nil {
		log.Fatal(err)
		return
	}
	page, err := driver.NewPage()
	if err != nil {
		log.Fatal(err)
		return
	}

	// 出勤予定日かチェックする
	today := time.Now()
	isWorkingDay, err := isWorkingDay(page, domain, user, pw, today)
	if err != nil {
		panic(err)
	}
	if !isWorkingDay {
		log.Printf("not working day: %s\n", today.Format("2006-01-02"))
		return
	}

	// 出勤予定日ならアクション
	latestAction, err := dakoku(page, domain, user, pw, startOrEnd == "start")
	if err != nil {
		panic(err)
	}
	log.Println(latestAction)
}

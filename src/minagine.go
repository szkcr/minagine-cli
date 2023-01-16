package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sclevine/agouti"
)

const (
	MINAGINE_LOGIN_URL    = "https://tm.minagine.net/index.html"
	MINAGINE_USERINFO_URL = "https://tm.minagine.net/employ/emplyself"
	MINAGINE_SHIFT_URL    = "https://tm.minagine.net/work/wrktimemngmntshtself/sht"
	MINAGINE_DAKOKU_URL   = "https://tm.minagine.net/work/wrktimergst"
)

func openPage(page *agouti.Page, targetURL string, domain string, user string, pw string) error {
	if err := page.Navigate(targetURL); err != nil {
		return err
	}
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

		// ログインできたことを確認
		if err := page.Navigate(MINAGINE_USERINFO_URL); err != nil {
			return err
		}
		renderedID, err := page.FindByXPath(`//*[@id="input_area"]/form/div[1]/table[3]/tbody/tr[2]/td/span[1]`).Text()
		if err != nil {
			return err
		}
		if renderedID != user {
			return fmt.Errorf("loginID not matched (expected:`%s`, actual:`%s`)", user, renderedID)
		}

		// 目的のページへ再度移動
		if err := page.Navigate(targetURL); err != nil {
			return err
		}
		url, err = page.URL()
		if err != nil {
			return err
		}
	}

	if url != targetURL {
		return fmt.Errorf("failed to navigate targetURL: `%s`", targetURL)
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

func isWorkingDay(page *agouti.Page, domain string, user string, pw string, today time.Time) (bool, string, error) {
	err := openPage(page, MINAGINE_SHIFT_URL, domain, user, pw)
	if err != nil {
		return false, "", err
	}

	rows := page.AllByXPath(`//*[@id="table_wrktimesht"]/tbody/tr`)
	numOfRow, err := rows.Count()
	if err != nil {
		return false, "", err
	}

	// 最初の4行はヘッダのためスキップ
	for i := 4; i < numOfRow; i++ {
		// `日`列をチェックして対象日との一致をチェック
		day, err := rows.At(i).FindByXPath(`td[1]`).Text()
		if err != nil {
			return false, "", err
		}
		if strconv.Itoa(today.Day()) == strings.TrimSpace(day) {
			// 出勤予定(`○`)かどうかをチェック
			status, err := rows.At(i).FindByXPath(`td[3]`).Text()
			if err != nil {
				return false, "", err
			}
			trimmed := strings.TrimSpace(status)
			return trimmed == "○", trimmed, nil
		}
	}

	return false, "", fmt.Errorf("could not find date in shift table")
}

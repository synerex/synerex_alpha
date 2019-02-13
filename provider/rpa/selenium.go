package main

import (
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/sclevine/agouti"
)

func main() {
	// Chromeを指定する
	driver := agouti.ChromeDriver(agouti.Browser("chrome"))
	if err := driver.Start(); err != nil {
		log.Fatalf("Failed to start driver: %v", err)
	}
	defer driver.Stop()

	page, err := driver.NewPage()
	if err != nil {
		log.Fatalf("Failed to open new page: %v", err)
	}

	// Sample of サイボウズOffice
	if err := page.Navigate("https://onlinedemo.cybozu.info/scripts/office10/ag.cgi?"); err != nil {
		log.Fatalf("Failed to navigate: %v", err)
	}

	//	time.Sleep(3 * time.Second)

	//
	pageContent, errPage := page.HTML() // get whole page
	if errPage != nil {
		println("Error:", errPage.Error())
	}

	// by using goquery, to obtain user lists.
	readerOfPage := strings.NewReader(pageContent)
	pageDom, pErr := goquery.NewDocumentFromReader(readerOfPage)
	if pErr != nil {
		println("PrintErr:", pErr)
	}
	//	println("DomText:", pageDom.Text())

	selectDom := pageDom.Find("select[name='_ID']").Children()
	//	selText := selectDom.Text()
	//	println("SelectText:", selText)
	users := make([]string, selectDom.Length())
	selectDom.Each(func(i int, sel *goquery.Selection) {
		tx := sel.Text()
		users[i] = tx
		println(i, tx)
	})

	//	nameX := page.FindByXPath("//div[@id='content-wrapper']/div/div/table/tbody/tr/td/center/div/table/tbody/tr[2]/td[2]/div/div/table[2]/tbody/tr[2]/td/table/tbody/tr[2]/td/select")
	nameX := page.FindByName("_ID")
	_, errn2 := nameX.Count()
	if errn2 != nil {
		println("Can't find Path", errn2.Error())
	}

	err = nameX.Select(users[1]) // ("高橋 健太")
	if err != nil {
		println("Select Error!", err.Error())
	}
	//	time.Sleep(3 * time.Second)

	submitButton := page.FindByName("Submit")
	_, e2 := submitButton.Count()
	if e2 != nil {
		println("Login Error!", e2.Error())
	}
	submitButton.Click() // ログインクリック

	println("Done!")

	time.Sleep(3 * time.Second)

	//	if err := page.FindByXPath("//*[@id=\"content-wrapper\"]/div[1]/div/table/tbody/tr/td/center/div/table/tbody/tr[2]/td[2]/div/div/table[2]/tbody/tr[2]/td/table/tbody/tr[7]/td/input").Click(); err != nil {
	//		log.Fatalf("Failed to submit: %v", err)
	//	}

	/*
		// Sample of desknets
		if err := page.Navigate("https://www.desknets.com/neo/trial/online.html"); err != nil {
			log.Fatalf("Failed to navigate: %v", err)
		}
		// PC用オンラインデモをクリック
		page.FindByID("ebislink3").Click()

		time.Sleep(5 * time.Second)

		// ログイン必須項目を選択

		name := page.AllByName("uid")

		ct2, errn2 := name.Count()
		fmt.Printf("Count is %d,  : %v\n", ct2, errn2)

		err3 :=	name.Click()
		if err3 != nil{
			println(err3.Error())
		}

		err = name.Select("鈴木誠") // ここの選択がエラーの原因っぽい

		time.Sleep(1 * time.Second)

		log.Printf("name.Select('2') is %v", err)
		// Submit

		page.FindByID("login-btn").Click()
	*/
	//	処理完了後、3秒間ブラウザを表示したままにする
	time.Sleep(3 * time.Second)

}

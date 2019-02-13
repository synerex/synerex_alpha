package main

import (
	"log"
	"time"

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
	name := page.FindByName("_ID")
	name.Select("17") // Select()が正常に動いていない
	time.Sleep(3 * time.Second)
	if err := page.FindByXPath("//*[@id=\"content-wrapper\"]/div[1]/div/table/tbody/tr/td/center/div/table/tbody/tr[2]/td[2]/div/div/table[2]/tbody/tr[2]/td/table/tbody/tr[7]/td/input").Click(); err != nil {
		log.Fatalf("Failed to submit: %v", err)
	}

	// // Sample of desknets
	// if err := page.Navigate("https://www.desknets.com/neo/trial/online.html"); err != nil {
	// 	log.Fatalf("Failed to navigate: %v", err)
	// }
	// // PC用オンラインデモをクリック
	// page.FindByID("ebislink3").Click()
	// // ログイン必須項目を選択
	// name := page.FindByName("uid")
	// name.Select("2") // ここの選択がエラーの原因っぽい
	// log.Printf("name.Select('2') is %v", name.Select("2"))
	// // Submit
	// if err := page.FindByClass("jlogin-submit").Submit(); err != nil {
	// 	log.Fatalf("Failed to submit: %v", err)
	// }

	// 処理完了後、3秒間ブラウザを表示したままにする
	time.Sleep(3 * time.Second)

}

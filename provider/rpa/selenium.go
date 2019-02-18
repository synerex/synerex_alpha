package main

import (
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/sclevine/agouti"
)

func main() {
	// set to use Chrome
	driver := agouti.ChromeDriver(agouti.Browser("chrome"))
	if err := driver.Start(); err != nil {
		log.Fatalf("Failed to start driver: %v", err)
	}
	defer driver.Stop()

	page, err := driver.NewPage()
	if err != nil {
		log.Fatalf("Failed to open new page: %v", err)
	}

	// sample of Cybozu Office
	if err := page.Navigate("https://onlinedemo.cybozu.info/scripts/office10/ag.cgi?"); err != nil {
		log.Fatalf("Failed to navigate: %v", err)
	}

	// get whole page
	pageContent, errPage := page.HTML()
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
		// println(i, tx)
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

	submitButton := page.FindByName("Submit")
	_, e2 := submitButton.Count()
	if e2 != nil {
		println("Login Error!", e2.Error())
	}
	submitButton.Click() // ログインクリック
	println("Logged in!")

	/* レイアウトが頻繁に変わるため
	schedule := page.FindByXPath("//*[@id='appIconMenuFrame']/div[2]/span[5]/a")
	err = schedule.Click()
	if err != nil {
		println("Cant Click Schedule:", err.Error())
	}
	*/

	// get whole page again
	pageContent, errPage = page.HTML()
	if errPage != nil {
		println("Error:", errPage.Error())
	}

	// by using goquery, to obtain Schedule
	readerOfPage = strings.NewReader(pageContent)
	pageDom, pErr = goquery.NewDocumentFromReader(readerOfPage)
	if pErr != nil {
		println("PrintErr:", pErr)
	}

	groupDom := pageDom.Find("select[name='GID']").Children()
	groups := make([]string, groupDom.Length())
	groupDom.Each(func(i int, g *goquery.Selection) {
		tx := g.Text()
		groups[i] = tx
		// println(i, tx)
	})

	groupX := page.FindByName("GID")
	_, err = groupX.Count()
	if err != nil {
		println("Can't find Path", err.Error())
	}

	err = groupX.Select(groups[10]) // "会議室"
	if err != nil {
		println("Select Error!", err.Error())
	}

	// get whole page again
	pageContent, errPage = page.HTML()
	if errPage != nil {
		println("Error:", errPage.Error())
	}

	// by using goquery, to obtain Schedule
	readerOfPage = strings.NewReader(pageContent)
	pageDom, pErr = goquery.NewDocumentFromReader(readerOfPage)
	if pErr != nil {
		println("PrintErr:", pErr)
	}

	calendarDom := pageDom.Find("#redraw > table > tbody").Children()
	rooms := make(map[string][]string, calendarDom.Length()) // 会議室
	calendarDom.Each(func(i int, sel *goquery.Selection) {
		if i == 0 { //  dates.
			sel.Children().Each(func(j int, cc *goquery.Selection) {
				// for each td
				if j == 0 {
					rooms["dates"] = []string{}
				} else {
					st := strings.TrimSpace(cc.Text())
					rooms["dates"] = append(rooms["dates"], st)
					// println("Dates", j, "[", st, "]")
				}
			})
		} else { //rooms
			rname := "none"
			sel.Children().Each(func(j int, cc *goquery.Selection) {
				if j == 0 {
					rname = strings.Trim(cc.Children().First().First().Text(), " \n")
					rname = strings.TrimSpace(rname)
					// println("RoomName:", i, rname)
					rooms[rname] = []string{}
				} else {
					st := strings.Trim(cc.Text(), " \n")
					st = strings.TrimSpace(st)
					rooms[rname] = append(rooms[rname], st)
					// println("RoomState:", j, "[", st, "]")
				}
			})
		}
	})

	// for k, v := range rooms {
	// 	fmt.Printf("rooms[%v]: %v\n", k, v)
	// }

	for i := 0; i < len(rooms["dates"]); i++ {
		println("----------")
		println(rooms["dates"][i])
		println("第一会議室:", rooms["第一会議室"][i])
		println("第二会議室:", rooms["第二会議室"][i])
		println("打合せルーム:", rooms["打合せルーム"][i])
	}

	// subscribe date from user
	userDate := "18（月）"

	// 会議室のカレンダーと比較する
	// 空室なら予約 or 満室なら別日程を促す
	for i := 0; i < len(rooms["dates"]); i++ {
		if rooms["dates"][i] == userDate {
			println("----------")
			println("userDate:", userDate)
			println("第一会議室:", rooms["第一会議室"][i])
			println("第二会議室:", rooms["第二会議室"][i])
			println("打合せルーム:", rooms["打合せルーム"][i])
		}
	}

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

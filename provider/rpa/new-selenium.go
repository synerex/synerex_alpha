package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/sclevine/agouti"
)

func getPageDOM(page *agouti.Page) *goquery.Document {
	// get whole page
	wholePage, err := page.HTML()
	if err != nil {
		fmt.Println("failed to get whole page:", err)
	}
	// use goquery
	readerOfPage := strings.NewReader(wholePage)
	pageDom, err := goquery.NewDocumentFromReader(readerOfPage)
	if err != nil {
		fmt.Println("failed to get page dom:", err)
	}
	return pageDom
}

func main() {
	// set of Chrome
	driver := agouti.ChromeDriver(agouti.Browser("chrome"))
	if err := driver.Start(); err != nil {
		println("", err)
		fmt.Println("Failed to start driver:", err)
	}
	defer driver.Stop()

	page, err := driver.NewPage()
	if err != nil {
		fmt.Println("Failed to open new page:", err)
	}

	// sample: Cybozu
	if err := page.Navigate("https://onlinedemo.cybozu.info/scripts/office10/ag.cgi?"); err != nil {
		fmt.Println("Failed to navigate:", err)
	}

	// Login: get user list
	usersDom := getPageDOM(page).Find("select[name='_ID']").Children()
	users := make([]string, usersDom.Length())
	usersDom.Each(func(i int, sel *goquery.Selection) {
		users[i] = sel.Text()
		fmt.Println(i, users[i])
	})

	// Login: set login user
	name := page.FindByName("_ID")
	if _, err := name.Count(); err != nil {
		fmt.Println("Can't find path", err)
	}
	name.Select(users[1]) // "高橋 健太"

	// Login: click login button
	submitBtn := page.FindByName("Submit")
	if _, err := submitBtn.Count(); err != nil {
		fmt.Println("Failed to login:", err)
	}
	submitBtn.Click()
	fmt.Println("Logged in:", users[1])

	// Schedule: get group list
	groupsDom := getPageDOM(page).Find("select[name='GID']").Children()
	groups := make([]string, groupsDom.Length())
	groupsDom.Each(func(i int, sel *goquery.Selection) {
		groups[i] = sel.Text()
		fmt.Println(i, groups[i])
	})

	// Schedule: move to meeting room page
	group := page.FindByName("GID")
	if _, err := group.Count(); err != nil {
		fmt.Println("Can't find path", err)
	}
	group.Select(groups[11]) // "会議室"

	// Schedule: get schedules
	schedulesDom := getPageDOM(page).Find("#redraw > table > tbody").Children()
	rooms := make(map[string][]string, schedulesDom.Length())
	schedulesDom.Each(func(i int, sel *goquery.Selection) {
		if i == 0 {
			sel.Children().Each(func(j int, cc *goquery.Selection) {
				if j == 0 {
					rooms["dates"] = []string{}
				} else {
					st := strings.TrimSpace(cc.Text())
					rooms["dates"] = append(rooms["dates"], st)
				}
			})
		} else {
			roomName := "none"
			sel.Children().Each(func(j int, cc *goquery.Selection) {
				if j == 0 {
					roomName = strings.Trim(cc.Children().First().First().Text(), " \n")
					roomName = strings.TrimSpace(roomName)
					rooms[roomName] = []string{}
				} else {
					st := strings.Trim(cc.Text(), "\n")
					st = strings.TrimSpace(st)
					rooms[roomName] = append(rooms[roomName], st)
				}
			})
		}
	})

	for k, v := range rooms {
		fmt.Printf("rooms[%v]: %v\n", k, v)
	}

	time.Sleep(5 * time.Second)
}

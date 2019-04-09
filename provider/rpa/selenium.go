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
		fmt.Println("Failed to get whole page:", err)
	}
	// use goquery
	readerOfPage := strings.NewReader(wholePage)
	pageDom, err := goquery.NewDocumentFromReader(readerOfPage)
	if err != nil {
		fmt.Println("Failed to get page dom:", err)
	}
	return pageDom
}

func login(page *agouti.Page) {
	// get user list
	usersDom := getPageDOM(page).Find("select[name='_ID']").Children()
	users := make([]string, usersDom.Length())
	usersDom.Each(func(i int, sel *goquery.Selection) {
		users[i] = sel.Text()
		fmt.Println(i, users[i])
	})
	// set login user
	name := page.FindByName("_ID")
	if _, err := name.Count(); err != nil {
		fmt.Println("Cannot find path:", err)
	}
	name.Select(users[1])
	// click login button
	submitBtn := page.FindByName("Submit")
	if _, err := submitBtn.Count(); err != nil {
		fmt.Println("Failed to login:", err)
	}
	submitBtn.Click()
}

func booking(page *agouti.Page) {
	reserveButton := page.FindByXPath("//*[@id=\"content-wrapper\"]/div[4]/div/div[1]/table/tbody/tr/td[1]/table/tbody/tr/td[1]/span/span/a")
	_, err := reserveButton.Count()
	if err != nil {
		fmt.Println("Cannot find path:", err)
	}
	reserveButton.Click()

	// set the date
	yearDom := getPageDOM(page).Find("select[name='SetDate.Year']").Children()
	monthDom := getPageDOM(page).Find("select[name='SetDate.Month']").Children()
	dayDom := getPageDOM(page).Find("select[name='SetDate.Day']").Children()
	startHourDom := getPageDOM(page).Find("select[name='SetTime.Hour']").Children()
	startMinuteDom := getPageDOM(page).Find("select[name='SetTime.Minute']").Children()
	endHourDom := getPageDOM(page).Find("select[name='EndTime.Hour']").Children()
	endMinuteDom := getPageDOM(page).Find("select[name='EndTime.Minute']").Children()

	years := make([]string, yearDom.Length())
	months := make([]string, monthDom.Length())
	days := make([]string, dayDom.Length())
	startHours := make([]string, startHourDom.Length())
	startMinutes := make([]string, startMinuteDom.Length())
	endHours := make([]string, endHourDom.Length())
	endMinutes := make([]string, endMinuteDom.Length())

	yearDom.Each(func(i int, g *goquery.Selection) {
		tx := g.Text()
		years[i] = tx
	})
	monthDom.Each(func(i int, g *goquery.Selection) {
		tx := g.Text()
		months[i] = tx
	})
	dayDom.Each(func(i int, g *goquery.Selection) {
		tx := g.Text()
		days[i] = tx
	})
	startHourDom.Each(func(i int, g *goquery.Selection) {
		tx := g.Text()
		startHours[i] = tx
	})
	startMinuteDom.Each(func(i int, g *goquery.Selection) {
		tx := g.Text()
		startMinutes[i] = tx
	})
	endHourDom.Each(func(i int, g *goquery.Selection) {
		tx := g.Text()
		endHours[i] = tx
	})
	endMinuteDom.Each(func(i int, g *goquery.Selection) {
		tx := g.Text()
		endMinutes[i] = tx
	})

	yearX := page.FindByName("SetDate.Year")
	_, err = yearX.Count()
	if err != nil {
		fmt.Println("Cannot find path:", err)
	}
	monthX := page.FindByName("SetDate.Month")
	_, err = monthX.Count()
	if err != nil {
		fmt.Println("Cannot find path:", err)
	}
	dayX := page.FindByName("SetDate.Day")
	_, err = dayX.Count()
	if err != nil {
		fmt.Println("Cannot find path:", err)
	}
	startHourX := page.FindByName("SetTime.Hour")
	_, err = startHourX.Count()
	if err != nil {
		fmt.Println("Cannot find path:", err)
	}
	startMinuteX := page.FindByName("SetTime.Minute")
	_, err = startMinuteX.Count()
	if err != nil {
		fmt.Println("Cannot find path:", err)
	}
	endHourX := page.FindByName("EndTime.Hour")
	_, err = endHourX.Count()
	if err != nil {
		fmt.Println("Cannot find path:", err)
	}
	endMinuteX := page.FindByName("EndTime.Minute")
	_, err = endMinuteX.Count()
	if err != nil {
		fmt.Println("Cannot find path:", err)
	}

	err = yearX.Select(years[22])
	if err != nil {
		fmt.Println("Select Error:", err)
	}
	err = monthX.Select(months[3])
	if err != nil {
		fmt.Println("Select Error:", err)
	}
	err = dayX.Select(days[20])
	if err != nil {
		fmt.Println("Select Error:", err)
	}
	err = startHourX.Select(startHours[11])
	if err != nil {
		fmt.Println("Select Error:", err)
	}
	err = startMinuteX.Select(startMinutes[2])
	if err != nil {
		fmt.Println("Select Error:", err)
	}
	err = endHourX.Select(endHours[12])
	if err != nil {
		fmt.Println("Select Error:", err)
	}
	err = endMinuteX.Select(endMinutes[2])
	if err != nil {
		fmt.Println("Select Error:", err)
	}

	// set the title
	title := page.FindByName("Detail")
	title.Fill("Test Booking")

	theRoomY := page.FindByXPath("//*[@id=\"content-wrapper\"]/div[4]/div/form/div[2]/table/tbody/tr/td/table/tbody/tr[2]/td/div/div[1]/div/table/tbody/tr[7]/td/table/tbody/tr[1]/td[3]/select/option[2]")
	theRoomY.Click()

	time.Sleep(2 * time.Second)

	// submit to make a reservation
	entryButton := page.FindByName("Entry")
	_, err = entryButton.Count()
	if err != nil {
		println("Login Error:", err)
	}
	entryButton.Click()
	println("Made a reservation.")
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

	// sample Cybozu
	if err := page.Navigate("https://onlinedemo.cybozu.info/scripts/office10/ag.cgi?"); err != nil {
		fmt.Println("Failed to navigate:", err)
	}

	// login
	login(page)

	// get group list
	groupsDom := getPageDOM(page).Find("select[name='GID']").Children()
	groups := make([]string, groupsDom.Length())
	groupsDom.Each(func(i int, sel *goquery.Selection) {
		groups[i] = sel.Text()
		fmt.Println(i, groups[i])
	})

	// move to meeting room page
	group := page.FindByName("GID")
	if _, err := group.Count(); err != nil {
		fmt.Println("Cannot find path:", err)
	}
	group.Select(groups[10]) // "会議室"

	// get schedules
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

	// make a reservation
	booking(page)

	time.Sleep(3 * time.Second)
}

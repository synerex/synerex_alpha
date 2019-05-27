package main

import (
	"fmt"
	"time"

	"github.com/sclevine/agouti"
)

func main() {
	driver := agouti.ChromeDriver(agouti.Browser("chrome"))
	if err := driver.Start(); err != nil {
		fmt.Println("Failed to start driver:", err)
	}
	defer driver.Stop()

	page, err := driver.NewPage()
	if err != nil {
		fmt.Println("Failed to open new page:", err)
	}

	if err := page.Navigate("https://golang.org/"); err != nil {
		fmt.Println("Failed to navigate:", err)
	}

	time.Sleep(3 * time.Second)
}

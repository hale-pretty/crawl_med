package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

type Company struct {
	Name    string
	Origin  string
	Address string
}

type Medicine struct {
	GPLH        string
	ExpiredDate string
	Name        string
	API         string
	Content     string
	Decision    string
	IssuedDate  string
	IssuedNo    string
	DosageForm  string
	Packaging   string
	Standard    string
	ShelfLife   string
	RegCom      *Company
	MfgCom      *Company
}

func processStr(value string) string {
	return strings.TrimSpace(value)
}

//creat a crawl func

func main() {
	// Get keyword
	reader := bufio.NewReader(os.Stdin)
	medicines := make([]Medicine, 0)
	fmt.Print("Enter search keyword: ")
	keyword, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("could not read input: %v", err)
	}
	keyword = strings.TrimSpace(keyword)

	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false), // Disable headless mode to show the browser
	})
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}
	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	if _, err = page.Goto("https://dichvucong.dav.gov.vn/congbothuoc/index"); err != nil {
		log.Fatalf("could not goto: %v", err)
	}

	// find search box
	searchBox := page.Locator("input[placeholder='Nhập từ khóa tìm kiếm theo Số GPLH và Tên thuốc...']")

	if err := searchBox.Fill(keyword); err != nil {
		log.Fatalf("could not fill keyword: %v", err)
	}

	if err := searchBox.Press("Enter"); err != nil {
		log.Fatalf("could not press enter: %v", err)
	}

	// Sleep for 1 second to wait for search results to load
	time.Sleep(1 * time.Second)

	nextPageNo := 2

	for {
		log.Printf("Crawl page %v", rune(nextPageNo)-1)
		trs, err := page.Locator("tr[class='ng-scope']").All()
		if err != nil {
			log.Fatalf("could not get table rows: %v", err)
		}
		for _, tr := range trs {
			tds, err := tr.Locator("td").All()
			if err != nil {
				log.Fatalf("could not get text content: %v", err)
			}
			values := make([]string, 0)
			for _, td := range tds {
				content, _ := td.TextContent()
				values = append(values, processStr(content))
			}
			medicine := Medicine{
				GPLH:        values[3],
				ExpiredDate: values[4],
				Name:        values[5],
				API:         values[6],
				Content:     values[7],
				Decision:    values[8],
				IssuedDate:  values[9],
				IssuedNo:    values[10],
				DosageForm:  values[11],
				Packaging:   values[12],
				Standard:    values[13],
				ShelfLife:   values[14],
				RegCom: &Company{
					Name:    values[15],
					Origin:  values[16],
					Address: values[17],
				},
				MfgCom: &Company{
					Name:    values[18],
					Origin:  values[19],
					Address: values[20],
				},
			}
			medicines = append(medicines, medicine)
		}

		log.Println(medicines)

		// Turn to next page
		nextPageNoStr := strconv.Itoa(nextPageNo)
		selector := fmt.Sprintf("a[data-page='%s']", nextPageNoStr)
		log.Println(selector + " next page")
		li := page.Locator(selector)

		exists, err := li.Count()

		if exists == 0 {
			break
		}
		li.Click()
		time.Sleep(1 * time.Second)
		nextPageNo += 1
	}

	if err = browser.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}
	if err = pw.Stop(); err != nil {
		log.Fatalf("could not stop Playwright: %v", err)
	}
}

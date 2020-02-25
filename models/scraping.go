package models

import (
	"../settings"
	"fmt"
	"github.com/gocolly/colly"
	"strconv"
	"strings"
)

type ScrapingItems struct {
	Url         string `json:"url"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Description string `json:"description"`
	Title       string `json:"title"`
	Country     string `json:"country"`
	Region      string `json:"region"`
	City        string `json:"city"`
	Phone       string `json:"phone"`
	Category    string `json:"category"`
	Subcategory string `json:"subcategory"`
}

type Scraping struct {
	items []ScrapingItems
}

func (then *Scraping) scrapingDetail(conf string) {
	scrap := ScrapingItems{
		Url: conf,
	}
	c := colly.NewCollector()
	c.IgnoreRobotsTxt = false
	c.OnHTML("#contact .name", func(e *colly.HTMLElement) {
		scrap.Name = strings.Replace(e.Text, "Nombre: ", "", -1)
	})
	//c.OnHTML("#contact .email", func(e *colly.HTMLElement) {
	//	if e.Text != "[email protected]" {
	//		scrap.Email += strings.Replace(e.Text, "E-mail: ", "", -1)
	//	}
	//})
	c.OnHTML("#contact .phone", func(e *colly.HTMLElement) {
		scrap.Phone = strings.Replace(e.Text, "Teléfono: ", "", -1)
	})
	c.OnHTML("#item-content div#description p", func(e *colly.HTMLElement) {
		scrap.Description += strings.Replace(e.Text, "<br>", "\n", -1)
	})
	c.OnHTML("h1 strong", func(e *colly.HTMLElement) {
		scrap.Title = e.Text
	})
	scrap.Country = "ES"
	scrap.Region = "07"
	scrap.City = "115345"
	url := strings.Split(conf, "/")
	scrap.Category = url[3]
	scrap.Subcategory = strings.Split(url[4], "_")[0]
	_ = c.Visit(conf)
	c.Wait()
	then.items = append(then.items, scrap)
}
func (then *Scraping) scrapingSearch(conf string) {
	var links = make(map[string]string)
	c := colly.NewCollector()
	c.IgnoreRobotsTxt = false
	c.OnHTML(".listing-basicinfo a.title[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		links[link] = link
	})
	_ = c.Visit(conf)
	c.Wait()
	if len(links) > 0 {
		for _, pag := range links {
			then.scrapingDetail(pag)
		}
	}
	//if len(paginate) > 0 {
	//	for _, pag := range paginate {
	//		then.scrapingSearch(pag)
	//	}
	//}
}
func (then *Scraping) ScrapingStart(config settings.ScrapingConfig) {
	for _, conf := range config.Urls {
		num := 1
		c := colly.NewCollector()
		c.IgnoreRobotsTxt = false
		c.OnHTML(".searchPaginationLast.list-last", func(e *colly.HTMLElement) {
			link := e.Attr("href")
			pag := link
			pag = strings.Replace(link, fmt.Sprintf("%s/iPage,", conf), "", -1)
			pag = strings.TrimSuffix(pag, "/")
			num, _ = strconv.Atoi(pag)
		})
		_ = c.Visit(conf)
		c.Wait()
		for i := 1; i <= num; i++ {
			pag := fmt.Sprintf("%s/iPage,%d/", conf, i)
			println(pag)
			then.scrapingSearch(pag)
		}
	}
}

func (then *Scraping) Results() []ScrapingItems {
	return then.items
}

func ScrapingSave(items []ScrapingItems) (err error) {
	db := settings.InstanceDb
	q := `INSERT INTO 
    items(title, description, price, contact_name, contact_phone, fk_user_id, fk_id_category, fk_id_subcategory, fk_city, scraping_url) 
    VALUES 
    `
	vals := []interface{}{}
	var count = 0
	for _, row := range items {
		//q += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, %d,%d, $%d, $%d),", count+1, count+2, count+3, count+4, count+5, count+6, count+7, count+8, count+9, count+10)
		q += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, (SELECT id FROM categories WHERE slug=$%d LIMIT 1), (SELECT id FROM subcategories WHERE slug=$%d LIMIT 1), $%d, $%d),", count+1, count+2, count+3, count+4, count+5, count+6, count+7, count+8, count+9, count+10)
		count += 10
		vals = append(vals, row.Title, row.Description, 0, row.Name, row.Phone, 1, row.Category, row.Subcategory, row.City, row.Url)
	}
	//trim the last ,
	q = q[0 : len(q)-1]
	//prepare the statement
	_, err = db.Exec(q, vals...)
	return
}

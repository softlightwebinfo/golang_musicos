package libs

import (
	"../settings"
	"fmt"
	"io/ioutil"
	"time"
)

var sl string = "sitemap.xml"

type urlSet []url
type url struct {
	Loc        string `xml:"loc"`
	LastMod    string `xml:"lastmod"`
	ChangeFreq string `xml:"changefreq"`
	Priority   string `xml:"priority"`
}
type Sitemap struct {
	templateStart  string
	templateTag    string
	templateEndTag string
	template       string
	urls           urlSet
}

func (e *Sitemap) New() {
	e.templateStart = `<?xml version="1.0" encoding="UTF-8"?>`
	e.templateTag = `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd">`
	e.templateEndTag = `</urlset>`
}
func (e *Sitemap) pages() {
	var pages = []string{
		"/",
		"/musicians",
		"/search",
		"/blog",
		"/publish",
		"/forum",
		"/stores",
		"/events",
		"/help",
		"/recovery",
		"/login",
		"/register",
		"/new-forum",
		"/categories-forum",
		"/search/web",
		"/search/images",
	}
	config := settings.Config()
	for i, page := range pages {
		t := time.Now()
		priority := "1.00"
		if i != 0 {
			priority = "0.85"
		}
		e.urls = append(e.urls, url{
			ChangeFreq: "daily",
			LastMod: fmt.Sprintf("%d-%02d-%02d",
				t.Year(), t.Month(), t.Day()),
			Loc:      fmt.Sprintf("%s%s", config.Web, page),
			Priority: priority,
		})
	}
}
func (e *Sitemap) Generate() {
	e.pages()
	e.profile()
	e.categories()
	e.items()
	e.blocs()
	e.forum()
	e.events()
}

func (e *Sitemap) categories() {
	db := settings.InstanceDb
	profile := "musisians"
	rows, _ := db.Query("SELECT slug from categories_subcategories")
	for rows.Next() {
		var last time.Time = time.Now()
		var slug = ""
		c := url{
			Loc:        slug,
			Priority:   "0.85",
			ChangeFreq: "daily",
			LastMod: fmt.Sprintf("%d-%02d-%02d",
				last.Year(), last.Month(), last.Day()),
		}
		err := rows.Scan(&slug)
		if err != nil {
			continue
		}
		c.Loc = fmt.Sprintf("%s/%s-%s", settings.Config().Web, profile, slug)
		c.LastMod = fmt.Sprintf("%d-%02d-%02d",
			last.Year(), last.Month(), last.Day())
		e.urls = append(e.urls, c)
	}
}
func (e *Sitemap) profile() {
	db := settings.InstanceDb
	profile := "profile"
	pages := []string{
		"wall",
		"about",
		"gallery",
		"gallery-videos",
		"events",
	}
	rows, _ := db.Query("SELECT up.slug, updated_at FROM users_profile up inner join users u on up.fk_id_user = u.id  WHERE up.is_public=true and active=true")
	for rows.Next() {
		slug := ""
		var last time.Time
		c := url{
			Loc:        slug,
			Priority:   "0.85",
			ChangeFreq: "daily",
			LastMod:    "",
		}
		err := rows.Scan(&slug, &last)
		if err != nil {
			continue
		}
		c.Loc = fmt.Sprintf("%s/%s/%s", settings.Config().Web, profile, slug)
		c.LastMod = fmt.Sprintf("%d-%02d-%02d",
			last.Year(), last.Month(), last.Day())
		e.urls = append(e.urls, c)
		for _, page := range pages {
			c.Loc = fmt.Sprintf("%s/%s/%s/%s", settings.Config().Web, profile, slug, page)
			e.urls = append(e.urls, c)
		}
	}
}
func (e *Sitemap) items() {
	db := settings.InstanceDb
	profile := "detail"
	rows, _ := db.Query("SELECT slugify(title), id, updated_at FROM items")
	for rows.Next() {
		slug := ""
		id := ""
		var last time.Time = time.Now()
		c := url{
			Loc:        slug,
			Priority:   "0.85",
			ChangeFreq: "daily",
			LastMod:    "",
		}
		err := rows.Scan(&slug, &id, &last)
		if err != nil {
			continue
		}
		c.Loc = fmt.Sprintf("%s/%s/%s--%s", settings.Config().Web, profile, slug, id)
		c.LastMod = fmt.Sprintf("%d-%02d-%02d",
			last.Year(), last.Month(), last.Day())
		e.urls = append(e.urls, c)
	}
}
func (e *Sitemap) blocs() {
	db := settings.InstanceDb
	profile := "blog"
	rows, _ := db.Query("SELECT slug, id, updated_at FROM blogs")
	for rows.Next() {
		slug := ""
		id := ""
		var last time.Time = time.Now()
		c := url{
			Loc:        slug,
			Priority:   "0.85",
			ChangeFreq: "daily",
			LastMod:    "",
		}
		err := rows.Scan(&slug, &id, &last)
		if err != nil {
			continue
		}
		c.Loc = fmt.Sprintf("%s/%s/%s--%s", settings.Config().Web, profile, slug, id)
		c.LastMod = fmt.Sprintf("%d-%02d-%02d",
			last.Year(), last.Month(), last.Day())
		e.urls = append(e.urls, c)
	}
}
func (e *Sitemap) forum() {
	db := settings.InstanceDb
	profile := "forum-detail"
	rows, _ := db.Query("SELECT slug, id, updated_at FROM forum")
	for rows.Next() {
		slug := ""
		id := ""
		var last time.Time = time.Now()
		c := url{
			Loc:        slug,
			Priority:   "0.85",
			ChangeFreq: "daily",
			LastMod:    "",
		}
		err := rows.Scan(&slug, &id, &last)
		if err != nil {
			continue
		}
		c.Loc = fmt.Sprintf("%s/%s/%s--%s", settings.Config().Web, profile, slug, id)
		c.LastMod = fmt.Sprintf("%d-%02d-%02d",
			last.Year(), last.Month(), last.Day())
		e.urls = append(e.urls, c)
	}
}
func (e *Sitemap) events() {
	db := settings.InstanceDb
	profile := "forum-detail"
	rows, _ := db.Query("SELECT slug, id, updated_at FROM events")
	for rows.Next() {
		slug := ""
		id := ""
		var last time.Time = time.Now()
		c := url{
			Loc:        slug,
			Priority:   "0.85",
			ChangeFreq: "daily",
			LastMod:    "",
		}
		err := rows.Scan(&slug, &id, &last)
		if err != nil {
			continue
		}
		c.Loc = fmt.Sprintf("%s/%s/%s?id=%s", settings.Config().Web, profile, slug, id)
		c.LastMod = fmt.Sprintf("%d-%02d-%02d",
			last.Year(), last.Month(), last.Day())
		e.urls = append(e.urls, c)
	}
}

func (e *Sitemap) Save() {
	e.template += e.templateStart
	e.template += e.templateTag
	template := `<url><loc>%s</loc><changefreq>%s</changefreq><priority>%s</priority><lastmod>%s</lastmod></url>`
	for _, url := range e.urls {
		e.template += fmt.Sprintf(
			template,
			url.Loc,
			url.ChangeFreq,
			url.Priority,
			url.LastMod,
		)
	}
	e.template += e.templateEndTag
	err := ioutil.WriteFile(fmt.Sprintf("../front_end/static/%s", sl), []byte(e.template), 0644)
	if err != nil {
		println("Error save file", err.Error())
	} else {
		println("Save file cron")
	}
}

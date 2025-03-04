package certfr_scrapping

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

type Item struct {
	Title       string
	Link        string
	Ref         string
	Date        string
	Status      string
	Description string
}

func PrettyPrint(items []Item) {
	for _, item := range items {
		fmt.Printf("-%s\n", item.Title)
		fmt.Printf("  %s\n", item.Link)
		fmt.Printf("  %s\n", item.Ref)
		fmt.Printf("  %s\n", item.Date)
		fmt.Printf("  %s\n", item.Status)
		fmt.Printf("  %s\n", item.Description)
		fmt.Println()
	}
}

func CollectAlert() (alerts []Item) {
	c := colly.NewCollector()

	c.OnHTML(".cert-alert", func(e *colly.HTMLElement) {
		alert := Item{
			Title:       strings.TrimSpace(e.ChildText("h3")),
			Link:        strings.TrimSpace(e.ChildAttr("a", "href")),
			Ref:         strings.TrimSpace(e.ChildText(".item-ref")),
			Date:        strings.TrimSpace(e.ChildText(".item-date")),
			Status:      strings.TrimSpace(e.ChildText(".item-status")),
			Description: strings.TrimSpace(e.ChildText(".item-excerpt")),
		}
		if alert.Title != "" {
			alert.Link = "https://www.cert.ssi.gouv.fr" + alert.Link
			alerts = append(alerts, alert)
		}
	})

	_ = c.Visit("https://www.cert.ssi.gouv.fr/alerte/")
	return
}

func CollectCti() (ctis []Item) {
	c := colly.NewCollector()

	c.OnHTML(".cert-cti", func(e *colly.HTMLElement) {
		cti := Item{
			Title:       strings.TrimSpace(e.ChildText("h3")),
			Link:        strings.TrimSpace(e.ChildAttr("a", "href")),
			Ref:         strings.TrimSpace(e.ChildText(".item-ref")),
			Date:        strings.TrimSpace(e.ChildText(".item-date")),
			Status:      strings.TrimSpace(e.ChildText(".item-status")),
			Description: strings.TrimSpace(e.ChildText(".item-excerpt")),
		}
		if cti.Title != "" {
			cti.Link = "https://www.cert.ssi.gouv.fr" + cti.Link
			ctis = append(ctis, cti)
		}
	})

	_ = c.Visit("https://www.cert.ssi.gouv.fr/cti/")
	return
}

func CollectAvis() (avis []Item) {
	c := colly.NewCollector()

	c.OnHTML(".cert-avis", func(e *colly.HTMLElement) {
		avi := Item{
			Title:       strings.TrimSpace(e.ChildText("h3")),
			Link:        strings.TrimSpace(e.ChildAttr("a", "href")),
			Ref:         strings.TrimSpace(e.ChildText(".item-ref")),
			Date:        strings.TrimSpace(e.ChildText(".item-date")),
			Status:      strings.TrimSpace(e.ChildText(".item-status")),
			Description: strings.TrimSpace(e.ChildText(".item-excerpt")),
		}
		if avi.Title != "" {
			avi.Link = "https://www.cert.ssi.gouv.fr" + avi.Link
			avis = append(avis, avi)
		}
	})

	_ = c.Visit("https://www.cert.ssi.gouv.fr/avis/")
	return
}

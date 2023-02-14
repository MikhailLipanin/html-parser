package parsing

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/spf13/viper"
)

type ErrorType struct {
	Id      string
	Message string
}

func Parse() []ErrorType {
	c := colly.NewCollector()

	var (
		data [][]string
		row  []string
	)
	c.OnHTML("div.table-wrap", func(e *colly.HTMLElement) {
		e.ForEach("td[class]", func(id int, e *colly.HTMLElement) {
			row = append(row, e.Text)
			if id&1 != 0 {
				data = append(data, row)
				row = make([]string, 0)
			}
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	c.Visit(viper.GetString("site"))

	var ret []ErrorType

	for _, row := range data {
		msg := ""
		for j := 1; j < len(row); j++ {
			msg += row[j]
		}
		ret = append(ret, ErrorType{
			Id:      row[0],
			Message: msg,
		})
	}

	return ret
}

package main

import (
	"fmt"
	"strconv"

	"github.com/gocolly/colly"
)

// Post is a representation of each reddit block
type Post struct {
	Title   string
	Link    string
	Upvotes int
	Flair   string
}

var err error

func main() {
	// Instantiate default collector
	c := colly.NewCollector()
	// a.SQnoC3ObvgnGjWt90zD9Z
	// div._1oQyIsiPHYt6nx7VOmd1sz.bE7JgM2ex7W3aF3zci5bm.D3IyhBGwXo9jPwz-Ka0Ve
	c.OnHTML("div._1oQyIsiPHYt6nx7VOmd1sz.bE7JgM2ex7W3aF3zci5bm.D3IyhBGwXo9jPwz-Ka0Ve", func(e *colly.HTMLElement) {
		findTitle(e)
		//findVotes(e)
		findFlair(e)
	})
	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on
	c.Visit("https://new.reddit.com/r/wallstreetbets/search/?q=-flair%3AMeme%20-flair%3ASatire%20-flair%3AShitpost&restrict_sr=1&t=day&sort=hot")
}

func findTitle(e *colly.HTMLElement) {
	container := e
	container.ForEach("a.SQnoC3ObvgnGjWt90zD9Z", func(_ int, elem *colly.HTMLElement) {
		fmt.Print(elem.Text)
		fmt.Println(elem.Attr("href"))

	})

}

func findVotes(e *colly.HTMLElement) []int64 {
	var (
		vote  int64
		votes []int64
	)
	e.ForEach("div._1rZYMD_4xY3gRcSS3p8ODO._3a2ZHWaih05DgAOtvu6cIo", func(_ int, elem *colly.HTMLElement) {
		voteString := elem.Text

		if voteString == "Vote" {
			vote = 0
		} else {

			if string(voteString[len(elem.Text)-1]) == "k" {
				votenum, err := strconv.ParseFloat(string(voteString[:len(elem.Text)-1]), 64)
				check(err)
				vote = int64(votenum * float64(1000))
			} else {
				vote, err = strconv.ParseInt(voteString, 10, 0)
				check(err)

			}
		}
		votes = append(votes, vote)
	})
	return votes
}

func findFlair(e *colly.HTMLElement) {
	e.ForEach("div._2X6EB3ZhEeXCh1eIVA64XM span", func(_ int, elem *colly.HTMLElement) {
		fmt.Print(elem.Text)
	})
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

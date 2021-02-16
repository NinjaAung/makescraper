package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/gocolly/colly"
)

// Post is a representation of each reddit block
type Post struct {
	Title   string
	Flair   string
	Link    string
	Upvotes int64
}

var err error

func main() {
	// Flags
	sortBy := flag.String("sort-by", "", "Sort by: Hot, Top, Relevance, New or Comments")
	postFrom := flag.String("post-from", "", "Post-from: Hour, 24, Week Month or Year")
	flag.Parse()

	// Vars
	var posts []Post
	baseURL := "https://new.reddit.com/r/wallstreetbets/search/?q=-flair%3AMeme%20-flair%3ASatire%20-flair%3AShitpost&restrict_sr=1"

	switch sort := *sortBy; sort {
	case "Relevance":
		baseURL += "&sort=relevance"
	case "Hot":
		baseURL += "&sort=hot"
	case "Top":
		baseURL += "&sort=top"
	case "New":
		baseURL += "&sort=new"
	case "Comments":
		baseURL += "&sort=comments"
	default:
		baseURL += "&sort=hot"
	}

	switch post := *postFrom; post {
	case "Hour":
		baseURL += "&t=hour"
	case "24":
		baseURL += "&t=day"
	case "Week":
		baseURL += "&t=weel"
	case "Month":
		baseURL += "&t=month"
	case "Year":
		baseURL += "&t=year"
	default:
		baseURL += "&t=day"
	}

	// Instantiate default collector
	c := colly.NewCollector()

	c.OnHTML("div._1oQyIsiPHYt6nx7VOmd1sz.bE7JgM2ex7W3aF3zci5bm.D3IyhBGwXo9jPwz-Ka0Ve", func(e *colly.HTMLElement) {
		Title, Link := findTitleAndLink(e)
		Upvotes := findVotes(e)
		Flair := findFlair(e)
		posts = append(posts, Post{Title, Flair, Link, Upvotes})
	})
	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on
	c.Visit(baseURL)

	// for _, v := range posts {
	// 	fmt.Printf("\nTitle: %s\nFlair: %s\nLink: %s\nUpvotes: %d\n", v.Title, v.Flair, v.Link, v.Upvotes)
	// }

	outputToJSON(posts, "output.json")

}
func outputToJSON(posts []Post, fileName string) {
	postJSON, err := json.MarshalIndent(posts, "", "    ")
	check(err)
	f, err := os.Create(fileName)
	check(err)
	f.Write(postJSON)
	f.Close()

}

func findTitleAndLink(e *colly.HTMLElement) (title, link string) {
	container := e
	container.ForEach("a.SQnoC3ObvgnGjWt90zD9Z", func(_ int, elem *colly.HTMLElement) {
		title = elem.Text
		link = elem.Attr("href")

	})
	return title, "https://new.reddit.com" + link
}

func findVotes(e *colly.HTMLElement) (vote int64) {

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
	})
	return vote
}

func findFlair(e *colly.HTMLElement) (flair string) {
	e.ForEach("div._2X6EB3ZhEeXCh1eIVA64XM span", func(_ int, elem *colly.HTMLElement) {
		flair = elem.Text
	})
	return flair
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

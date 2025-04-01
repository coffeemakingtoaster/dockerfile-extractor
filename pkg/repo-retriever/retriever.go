package reporetriever

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	"github.com/coffeemakingtoaster/dockerfile-extractor/pkg/db"
	"github.com/gocolly/colly/v2"
	"github.com/rs/zerolog/log"
)

type RankingRetriever struct {
	targetUrl string
	page      int
}

func NewRankingRetriever() RankingRetriever {
	return RankingRetriever{
		targetUrl: "https://gitstar-ranking.com/repositories",
		page:      1,
	}
}

func (r *RankingRetriever) Scrape() {
	currentCount, _ := db.GetPresentHashCount()
	log.Info().Msgf("Already have %d in db!", currentCount)

	if currentCount > 0 {
		return
	}

	links := []string{}
	c := colly.NewCollector()
	c.OnHTML(".hidden-xs.hidden-sm", func(e *colly.HTMLElement) {
		content := e.Text
		content = strings.ReplaceAll(content, "\n", "")
		if !strings.Contains(content, "/") {
			return
		}
		links = append(links, fmt.Sprintf("github.com/%s", content))
	})
	c.Limit(&colly.LimitRule{
		Parallelism: 1,
		RandomDelay: 5 * time.Second,
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Warn().Msgf("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	log.Info().Msg("Scraping...")
	for r.page <= 100 {
		c.Visit(fmt.Sprintf("https://gitstar-ranking.com/repositories?page=%d", r.page))
		r.page++
	}
	log.Info().Msgf("Scraped %d pages", r.page-1)

	count := 0

	conn := db.RetrieveDbConn()
	defer conn.Close()

	for _, link := range links {
		h := sha256.New()
		h.Write([]byte(link))
		hash := h.Sum(nil)
		err := db.AddToDB(conn, string(hash), link)
		if err != nil {
			log.Error().Msgf("Could not write to db: %s", err.Error())
			continue
		}
		count++
	}
	log.Info().Msgf("Wrote %d items to db!", count)
}

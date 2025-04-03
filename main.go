package main

import (
	reporetriever "github.com/coffeemakingtoaster/dockerfile-extractor/pkg/repo-retriever"
	"github.com/coffeemakingtoaster/dockerfile-extractor/pkg/scraper"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("Starting...")
	retriever := reporetriever.NewRankingRetriever()
	retriever.Scrape()

	scraper.ScrapeAllDBContents()
}

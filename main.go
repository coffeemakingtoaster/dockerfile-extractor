package main

import (
	reporetriever "github.com/coffeemakingtoaster/dockerfile-extractor/pkg/repo-retriever"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("Starting...")
	scraper := reporetriever.NewRankingRetriever()
	scraper.Scrape()
}

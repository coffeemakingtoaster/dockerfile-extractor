package main

import (
	"os"

	ghapi "github.com/coffeemakingtoaster/dockerfile-extractor/pkg/gh-api"
	reporetriever "github.com/coffeemakingtoaster/dockerfile-extractor/pkg/repo-retriever"
	"github.com/coffeemakingtoaster/dockerfile-extractor/pkg/scraper"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Make logs human friendly
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("Starting...")
	if ghapi.CanUseAuth() {
		log.Info().Msg("Will use detected pat for gh api calls")
	}
	retriever := reporetriever.NewRankingRetriever()
	retriever.Scrape()
	scraper.ScrapeAllDBContents()
}

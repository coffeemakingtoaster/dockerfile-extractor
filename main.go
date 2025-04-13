package main

import (
	"os"
	"strconv"
	"strings"

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
	switch os.Args[1] {
	case "scrape":
		log.Info().Msg("Scraping dockerfiles from db")
		if ghapi.CanUseAuth() {
			log.Info().Msg("Will use detected pat for gh api calls")
		}
		retriever := reporetriever.NewRankingRetriever()
		retriever.Scrape()
		scraper.ScrapeAllDBContents()
	case "gather":
		entryPoint := os.Args[2]
		if !strings.Contains(entryPoint, "/") || strings.Contains(entryPoint, "github.com") {
			panic("Invalid repo form. Just use the USERNAME/REPO")
		}
		scrapeLimit, err := strconv.Atoi(os.Args[3])
		if err != nil {
			log.Warn().Msgf("Could not parse %s defaulting to 50", os.Args[3])
			scrapeLimit = 50
		}
		crawlingRetriever := reporetriever.NewCrawlingRetriever(entryPoint)
		crawlingRetriever.Retrieve(scrapeLimit)
	default:
		log.Error().Msg("Unknown command")
	}
}

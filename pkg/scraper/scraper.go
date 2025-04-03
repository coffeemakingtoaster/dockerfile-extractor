package scraper

import (
	"github.com/coffeemakingtoaster/dockerfile-extractor/pkg/db"
	ghapi "github.com/coffeemakingtoaster/dockerfile-extractor/pkg/gh-api"
	"github.com/rs/zerolog/log"
)

func ScrapeAllDBContents() {
	rowIterator := db.GetRepoIterator()
	defer rowIterator.Close()
	for rowIterator.Next() {
		var repo string
		var hash string
		err := rowIterator.Scan(&hash, &repo)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		log.Info().Msgf("%s %s", hash, repo)
		ghapi.GetDockerfilesFrom(repo)
		//TODO: add sleep
	}
}

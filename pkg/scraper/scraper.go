package scraper

import (
	"time"

	"github.com/coffeemakingtoaster/dockerfile-extractor/pkg/db"
	ghapi "github.com/coffeemakingtoaster/dockerfile-extractor/pkg/gh-api"
	"github.com/rs/zerolog/log"
)

const targetPath = "./contents"

func ScrapeAllDBContents() {
	rowIterator := db.GetRepoIterator()
	defer rowIterator.Close()
	scraped_count := 0
	rowcount := 0
	startTime := time.Now()
	for rowIterator.Next() {
		rowcount++
		var repo string
		var hash string
		err := rowIterator.Scan(&hash, &repo)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		dockerfiles := ghapi.GetDockerfilesFrom(repo)
		saveAllDockerfilesToDisk(dockerfiles)
		scraped_count += len(dockerfiles)
		log.Info().Msgf("Scraped %d files (%d rows)", scraped_count, rowcount)
	}
	log.Info().Msgf("Scraping took %v.", time.Now().Sub(startTime))
}

func saveAllDockerfilesToDisk(dockerfiles []ghapi.DockerFileInformation) {
	for _, dockerfile := range dockerfiles {
		err := dockerfile.PopulateContent()
		if err != nil {
			log.Warn().Msg(err.Error())
		}
		err = dockerfile.SaveToFile(targetPath)
		if err != nil {
			log.Warn().Msg(err.Error())
		}
	}
}

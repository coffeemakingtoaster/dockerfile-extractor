package scraper

import (
	"time"

	"github.com/coffeemakingtoaster/dockerfile-extractor/pkg/db"
	ghapi "github.com/coffeemakingtoaster/dockerfile-extractor/pkg/gh-api"
	"github.com/rs/zerolog/log"
)

const targetPath = "./contents"

func ScrapeAllDBContents() {
	currentCount, doneCount, _ := db.GetPresentHashCount()
	log.Info().Msgf("%d items in DB (%d done, %d todo)", currentCount, doneCount, currentCount-doneCount)

	blockSize := 10
	totalScrapedDockerfiles := 0
	totalScrapedRepositories := 0
	startTime := time.Now()

	lastIterRows := 100
	for lastIterRows >= blockSize {
		scraped_count, lastIterRows := scrapeBlock(blockSize)
		totalScrapedDockerfiles += scraped_count
		totalScrapedRepositories += lastIterRows
		log.Info().Msgf("Found %d dockerfiles (%d repos)", totalScrapedDockerfiles, totalScrapedRepositories)
	}
	log.Info().Msgf("Scraping took %v.", time.Now().Sub(startTime))
}

func scrapeBlock(blockSize int) (int, int) {
	conn := db.RetrieveDbConn()

	tx, err := conn.Begin()

	if err != nil {
		panic(err)
	}

	rowIterator := db.GetRepoIterator(tx)
	defer conn.Close()
	defer rowIterator.Close()
	scraped_count := 0
	rowCount := 0
	for rowIterator.Next() && rowCount <= blockSize {
		rowCount++
		var repo string
		var hash string
		err := rowIterator.Scan(&hash, &repo)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		dockerfiles := ghapi.GetDockerfilesFrom(repo)
		saveAllDockerfilesToDisk(dockerfiles)
		scraped_count += len(dockerfiles)
		err = db.SetScrapedByHash(tx, hash)
		if err != nil {
			log.Warn().Msgf("Could not update status due to an error %s", err.Error())
		}
	}
	tx.Commit()
	return scraped_count, rowCount
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

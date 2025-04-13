package reporetriever

import (
	"math/rand"

	"github.com/coffeemakingtoaster/dockerfile-extractor/pkg/db"
	ghapi "github.com/coffeemakingtoaster/dockerfile-extractor/pkg/gh-api"
	"github.com/coffeemakingtoaster/dockerfile-extractor/pkg/util"
	"github.com/rs/zerolog/log"
)

type CrawlingRetriever struct {
	QueuedRepositories []string
}

func NewCrawlingRetriever(entryPoint string) CrawlingRetriever {
	return CrawlingRetriever{
		QueuedRepositories: []string{entryPoint},
	}
}

func (cr *CrawlingRetriever) Retrieve(limit int) {
	currentCount, _, _ := db.GetPresentHashCount()
	log.Info().Msgf("Already have %d in db!", currentCount)

	addedRepos := 0

	for len(cr.QueuedRepositories) > 0 && addedRepos < limit {
		log.Info().Msgf("Currently at %d out of %d", addedRepos, limit)
		currentRepo := cr.QueuedRepositories[0]
		cr.QueuedRepositories = cr.QueuedRepositories[1:]
		exists, _ := db.CheckIfAlreadyPresent(util.HashString(currentRepo))
		if !exists {
			conn := db.RetrieveDbConn()
			db.AddToDB(conn, util.HashString(currentRepo), currentRepo)
			conn.Close()
			addedRepos += 1
		} else {
			log.Debug().Msg("Entry already in db")
		}
		contributers := ghapi.GetRepositoryContributers(currentRepo)
		for i := range contributers {
			j := rand.Intn(i + 1)
			contributers[i], contributers[j] = contributers[j], contributers[i]
		}
		for _, contributor := range contributers {
			userRepos := ghapi.GetUserRepositories(contributor.Name)
			if len(userRepos) > 0 {
				cr.QueuedRepositories = append(cr.QueuedRepositories, userRepos...)
			}
			// No need to add more repositories for now
			if len(cr.QueuedRepositories) > 1000 {
				break
			}
		}
	}
	log.Info().Msgf("Added %d to db (%d were left in repo queue)", addedRepos, len(cr.QueuedRepositories))
}

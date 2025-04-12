package reporetriever

import (
	"github.com/coffeemakingtoaster/dockerfile-extractor/pkg/db"
	ghapi "github.com/coffeemakingtoaster/dockerfile-extractor/pkg/gh-api"
	"github.com/coffeemakingtoaster/dockerfile-extractor/pkg/util"
	"github.com/rs/zerolog/log"
)

type CrawlingRetriever struct {
	QueuedRepositories []string
	QueuedUsers        []string
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

	for len(cr.QueuedRepositories) > 0 && addedRepos <= limit {
		currentRepo := cr.QueuedRepositories[0]
		cr.QueuedRepositories = cr.QueuedRepositories[1:]
		exists, _ := db.CheckIfAlreadyPresent(util.HashString(currentRepo))
		if exists {
			continue
		}
		contributers := ghapi.GetRepositoryContributers(currentRepo)
		for _, contributor := range contributers {
			userRepos := ghapi.GetUserRepositories(contributor.Name)
			if len(userRepos) > 0 {
				cr.QueuedRepositories = append(cr.QueuedRepositories, userRepos...)
			}
		}
	}
}

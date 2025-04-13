# Scraper 

Crawl repositories and retriever their Dockerfiles.

## Useage

Set the `GH_TOKEN` env variable to a PAT to use auth for every api request. Unauthorized requests have a way smaller rate limiting threshold.

For gathering repositories from github into the sqlite db use the `gather` command. This also expects an entrypoint and a limit

```sh
./main gather torvalds/linux 100
```

For scraping the dockerfiles from repositories use the `scrape` command. This will scrape the dockerfiles from every repository in the db.
If there are no repositories in the db, the top 1000 repositories will be scraped from [this ranking site](https://gitstar-ranking.com/repositories).

```sh
./main scrape
```

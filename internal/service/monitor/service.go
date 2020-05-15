package monitor

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"telebot/telebot/CA/internal/domain/repository"
)

type Service struct {
	CrawlerRepository repository.CrawlerRepository
}

func (svc *Service) CrawlerService(map[string]int, error) (map[string]int, error) {
	var SiteList = make(map[string]int)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	var httpClient = &http.Client{
		Timeout:   time.Second * 10,
		Transport: tr,
	}

	for {
		for site, _ := range SiteList {
			response, err := httpClient.Get(site)
			if err != nil {
				SiteList[site] = 1
				log.Printf("Status of %s: %s", site, "1 - Connection refused")

			} else {
				log.Printf("Status of %s: %s", site, response.Status)
				SiteList[site] = response.StatusCode
			}
		}
		time.Sleep(time.Minute * 5)
		return SiteList, nil
	}
}

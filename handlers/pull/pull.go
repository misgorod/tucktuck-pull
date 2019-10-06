package pull

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

type pullResponse struct {
	Count    int
	Next     string
	Previous string
	Results  result
}

type result struct {
	Id              int
	Title           string
	Slug            string
	PublicationDate int64 `json:"publication_date"`
	Place           string
	Description     string
	Dates           interface{}
	BodyText        string `json:"body_text"`
	Location        interface{}
	Categories      []string
	TagLine         string `json:"tagline"`
	AgeRestriction  int    `json:"age_restriction"`
	Price           string
	SsFree          bool `json:"is_free"`
	Images          interface{}
	FavouritesCount int    `json:"favourites_count"`
	CommentsCount   int    `json:"comments_count"`
	SiteUrl         string `json:"site_url"`
	ShortTitle      string `json:"short_title"`
	Tags            []string
	Participants    []string
}

type PullHandler struct{}

func New(ctx context.Context) PullHandler {
	go func() {
		select {
		// Listening for a cancellation event
		case <-ctx.Done():
			log.WithError(ctx.Err()).Error(" received cancel from external context")
			return
		default:
		}
		logger := log.WithFields(log.Fields{
			"handler": "pullHandler",
			"method":  "new",
		})
		actualSince := time.Now().Unix()
		url := fmt.Sprintf("https://kudago.com/public-api/v1.4/events/?fields=id,publication_date,dates,title,short_title,slug,place,description,body_text,location,categories,tagline,age_restriction,price,is_free,images,favorites_count,comments_count,site_url,tags,participants&expand=images,place,location,dates,participants&text_format=text&location=msk&actual_since=%v", actualSince)
		response, err := http.Get(url)
		if err != nil {
			logger.WithError(err).Error("couldn't get response from kudago")
		}
		_, err = ioutil.ReadAll(response.Body)
		if err != nil {
			logger.WithError(err).Error("couldn't read response from kudago")
		}
		//json.Unmarshal(body, )

	}()
	return PullHandler{}
}

func (p *PullHandler) Get(w http.ResponseWriter, r *http.Request) {

}

package pull

import (
	"context"
	"encoding/json"
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
	Id              int `json:"id"`
	Title           string `json:"title"`
	Slug            string `json:"slug"`
	PublicationDate int64 `json:"publication_date"`
	Place           place `json:"place"`
	Description     string `json:"description"`
	Dates           []date `json:"dates"`
	BodyText        string `json:"body_text"`
	Location        location `json:"location"`
	Categories      []string `json:"categories"`
	TagLine         string `json:"tagline"`
	AgeRestriction  int    `json:"age_restriction"`
	Price           string `json:"price"`
	SsFree          bool `json:"is_free"`
	Images          []image `json:"images"`
	FavouritesCount int    `json:"favourites_count"`
	CommentsCount   int    `json:"comments_count"`
	SiteUrl         string `json:"site_url"`
	ShortTitle      string `json:"short_title"`
	Tags            []string `json:"tags"`
	Participants    []string `json:"participants"`
}

type image struct {
	Image string `json:"image"`
	Source source `json:"source"`
}

type source struct {
	Link string `json:"link"`
	Source string `json:"source"`
}

type place struct {
	Id int `json:"id"`
	Title string `json:"title"`
	Slug string `json:"slug"`
	Address string `json:"address"`
	Phone string `json:"phone"`
	IsStub bool `json:"is_stub"`
	SiteUrl string `json:"site_url"`
	Coordinates coordinates `json:"coords"`
	Subway string `json:"subway"`
	IsClosed bool `json:"is_closed"`
	Location string `json:"location"`
}

type location struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
	Timezone string `json:"timezone"`
	Coordinates coordinates `json:"coords"`
	Language string `json:"language"`
	Currency string `json:"currency"`
}

type coordinates struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

type date struct {
	StartDate string `json:"start_date"`
	StartTime string `json:"start_time"`
	Start int64 `json:"start"`
	EndDate string `json:"end_date"`
	EndTime string `json:"end_time"`
	End int64 `json:"end"`
	IsContinuous bool `json:"is_continuous"`
	IsEndless bool `json:"is_endless"`
	IsStartless bool `json:"is_startless"`
	Schedules []interface{}
	UsePlaceSchedule bool `json:"use_place_schedule"`
}

type Handler struct{}

func New(ctx context.Context) Handler {
	go func() {
		select {
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
		body, err = ioutil.ReadAll(response.Body)
		if err != nil {
			logger.WithError(err).Error("couldn't read response from kudago")
		}
		json.Unmarshal(body, )

	}()
	return Handler{}
}

func (p *Handler) Get(w http.ResponseWriter, r *http.Request) {

}

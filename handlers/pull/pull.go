package pull

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/misgorod/tucktuck-pull/common"
	"github.com/misgorod/tucktuck-pull/models"
	"github.com/misgorod/tucktuck-pull/repository"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

type pullResponse struct {
	Count    int             `json:"count"`
	Next     string          `json:"next"`
	Previous string          `json:"previous"`
	Results  []models.Result `json:"results"`
}

type Handler struct {
	client       *repository.Client  `json:"-"`
	upsertResult models.UpsertResult `json:"insertResult"`
}

type InsertResult struct {
	MatchedCount  int64
	ModifiedCount int64
	UpsertedCount int64
	UpsertedID    interface{}
	LastTime      time.Time
}

func (h *Handler) makeRequest(ctx context.Context, logger *log.Logger) error {
	actualSince := time.Now().Unix()
	url := fmt.Sprintf("https://kudago.com/public-api/v1.4/events/?fields=id,publication_date,dates,title,short_title,slug,place,description,body_text,location,categories,tagline,age_restriction,price,is_free,images,favorites_count,comments_count,site_url,tags,participants&expand=images,place,location,dates,participants&text_format=text&location=msk&actual_since=%v", actualSince)
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	var pullResponse pullResponse
	err = json.Unmarshal(body, &pullResponse)
	if err != nil {
		return err
	}
	data := make([]interface{}, len(pullResponse.Results))
	for _, result := range pullResponse.Results {
		data = append(data, result)
	}
	upsertResult, err := h.client.UpsertMany(ctx, data)
	if err != nil {
		return err
	}
	h.upsertResult = upsertResult
	logger.WithField("upsert result", upsertResult).Info("got events")
	return nil
}

func New(ctx context.Context, client *repository.Client) *Handler {
	handler := &Handler{client: client}
	logger := log.WithFields(log.Fields{
		"handler": "pullHandler",
		"method":  "new",
	})
	ticker := time.NewTicker(time.Hour)
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.WithError(ctx.Err()).Error(" received cancel from external context")
				return
			case <-ticker.C:
				err := handler.makeRequest(ctx, logger.Logger)
				if err != nil {
					logger.WithError(err).Error()
				}
			default:
			}
		}
	}()
	return handler
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	err := h.makeRequest(r.Context(), log.StandardLogger())
	if err != nil {
		common.RespondError(r.Context(), w, http.StatusInternalServerError, err)
		return
	}
	common.RespondJSON(r.Context(), w, http.StatusOK, h.upsertResult)
}

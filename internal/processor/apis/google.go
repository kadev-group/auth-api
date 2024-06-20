package apis

import (
	"auth-api/internal/models"
	"auth-api/internal/models/consts"
	"context"
	"encoding/json"
	"fmt"
	"github.com/doxanocap/pkg/errs"
	"github.com/doxanocap/pkg/gohttp"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"net/http"
	"time"
)

type GoogleAPI struct {
	api *oauth2.Config

	awaitTime time.Duration
	config    *models.Config
	log       *zap.Logger
}

func InitGoogleAPI(config *models.Config, log *zap.Logger) *GoogleAPI {
	return &GoogleAPI{
		log:       log,
		config:    config,
		awaitTime: consts.GoogleAwaitTime,
		api: &oauth2.Config{
			ClientID:     config.OAuth.GoogleAPI.ClientID,
			ClientSecret: config.OAuth.GoogleAPI.ClientSecret,
			Scopes: []string{
				consts.GoogleScopeUserProfile,
				consts.GoogleScopeEmail},
			Endpoint:    google.Endpoint,
			RedirectURL: config.OAuth.ServerCallBackURI,
		},
	}
}

func (g *GoogleAPI) GetUserProfileByToken(ctx context.Context, accessToken string) (*models.UserProfile, error) {
	url := fmt.Sprintf("https://www.googleapis.com/oauth2/v3/userinfo?access_token=%s", accessToken)

	response, err := gohttp.NewRequest().
		SetURL(url).
		SetMethod(http.MethodGet).
		SetRequestFormat(gohttp.FormatJSON).
		Execute(ctx)
	if err != nil {
		return nil, err
	}
	// TODO handle statuses

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	up := &models.UserProfile{}
	if err := json.Unmarshal(body, up); err != nil {
		return nil, err
	}
	return up, nil
}

func (g *GoogleAPI) NewRedirectURL(state string) (res *models.GoogleRedirectRes, err error) {
	oAuthURLParams := []oauth2.AuthCodeOption{
		oauth2.SetAuthURLParam("response_type", "code"),
	}
	return &models.GoogleRedirectRes{
		RedirectURL: g.api.AuthCodeURL(state, oAuthURLParams...),
	}, nil
}

func (g *GoogleAPI) Exchange(ctx context.Context, exchangeCode string) (token *oauth2.Token, err error) {
	token, err = g.api.Exchange(ctx, exchangeCode)
	if err != nil {
		return nil, errs.Wrap("g.exchange", err)
	}
	return
}

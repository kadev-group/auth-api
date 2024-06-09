package oauth

import (
	"auth-api/internal/interfaces"
	"auth-api/internal/models"
	"auth-api/internal/models/consts"
	"auth-api/internal/pkg/tools"
	"context"
	"encoding/json"
	"fmt"
	"github.com/doxanocap/pkg/errs"
	"github.com/doxanocap/pkg/gohttp"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"net/http"
	"time"
)

type GoogleAPI struct {
	appConfig *models.Config
	manager   interfaces.IManager
	awaitTime time.Duration

	api *oauth2.Config
}

func InitGoogleAPI(manager interfaces.IManager, config *models.Config) *GoogleAPI {
	return &GoogleAPI{
		appConfig: config,
		manager:   manager,
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

func (g *GoogleAPI) GetRedirectURL(ctx context.Context, state string) (res *models.GoogleRedirectRes, err error) {
	if !tools.IsUUID(state) {
		return nil, models.ErrInvalidOAuthState
	}

	oAuthURLParams := []oauth2.AuthCodeOption{
		oauth2.SetAuthURLParam("response_type", "code"),
	}
	value, err := g.manager.Repository().GoogleAPICodes().Get(ctx, state)
	if err != nil {
		return
	}
	if value != 0 {
		return nil, models.HttpConflict
	}
	if err = g.manager.Repository().GoogleAPICodes().Set(ctx, state, time.Now().Unix()); err != nil {
		return
	}
	return &models.GoogleRedirectRes{
		RedirectURL: g.api.AuthCodeURL(state, oAuthURLParams...),
	}, nil
}

func (g *GoogleAPI) HandleCallBack(ctx context.Context, state, exchangeCode string) (string, error) {
	insertedAt, err := g.manager.Repository().GoogleAPICodes().Get(ctx, state)
	if err != nil {
		return "", err
	}

	now := time.Now()
	if insertedAt > now.Unix() || now.Add(g.awaitTime).Unix() < insertedAt {
		return "", models.ErrStateNotFound
	}

	token, err := g.api.Exchange(ctx, exchangeCode)
	if err != nil {
		return "", errs.Wrap("g.exchange", err)
	}

	up, err := g.getUserInfo(ctx, token.AccessToken)
	if err != nil {
		return "", errs.Wrap("g.GetUserInfo", err)
	}

	if err = g.manager.Repository().GoogleAPICodes().Delete(ctx, state); err != nil {
		return "", err
	}

	user, err := g.manager.Repository().Users().FindByEmail(ctx, up.Email)
	if err != nil {
		return "", err
	}
	if user == nil {
		tokens, err := g.manager.Service().User().Create(ctx, &models.UserDTO{
			Email:        up.Email,
			Activated:    true,
			AuthProvider: models.GoogleAuth,
		})
		if err != nil {
			return "", err
		}

		return tokens.Tokens.AccessToken, nil
	}

	tokens, err := g.manager.Service().Auth().UpdateSession(ctx, user)
	if err != nil {
		return "", err
	}

	return tokens.AccessToken, nil
}

type userProfile struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Locale        string `json:"locale"`
}

func (g *GoogleAPI) getUserInfo(ctx context.Context, accessToken string) (*userProfile, error) {
	url := fmt.Sprintf("https://www.googleapis.com/oauth2/v3/userinfo?access_token=%s", accessToken)

	response, err := gohttp.NewRequest().
		SetURL(url).
		SetMethod(http.MethodGet).
		SetRequestFormat(gohttp.FormatJSON).
		Execute(ctx)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	up := &userProfile{}
	if err := json.Unmarshal(body, up); err != nil {
		return nil, err
	}
	return up, nil
}

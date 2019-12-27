package providers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/bitly/oauth2_proxy/api"
)

type COCOProvider struct {
	*ProviderData
}

func NewCOCOProvider(p *ProviderData) *COCOProvider {
	p.ProviderName = "COCO"
	if p.LoginURL.String() == "" {
		p.LoginURL = &url.URL{Scheme: "https",
			Host: "www.api.getcoco.buzz",
			Path: "/oauth/authorize",
			// ?granted_scopes=true
		}
	}
	if p.RedeemURL.String() == "" {
		p.RedeemURL = &url.URL{Scheme: "https",
			Host: "www.api.getcoco.buzz",
			Path: "/oauth/token",
		}
	}
	if p.ProfileURL.String() == "" {
		p.ProfileURL = &url.URL{Scheme: "https",
			Host: "www.api.getcoco.buzz",
			Path: "/user-manager/users/me",
		}
	}
	if p.ValidateURL.String() == "" {
		p.ValidateURL = p.ProfileURL
	}
	if p.Scope == "" {
		p.Scope = "profile"
	}
	return &COCOProvider{ProviderData: p}
}

func getCOCOHeader(access_token string) http.Header {
	header := make(http.Header)
	header.Set("Accept", "application/json")
	header.Set("x-li-format", "json")
	header.Set("Authorization", fmt.Sprintf("Bearer %s", access_token))
	return header
}

func (p *COCOProvider) GetEmailAddress(s *SessionState) (string, error) {
	if s.AccessToken == "" {
		return "", errors.New("missing access token")
	}
	req, err := http.NewRequest("GET", p.ProfileURL.String(), nil)
	if err != nil {
		return "", err
	}
	req.Header = getCOCOHeader(s.AccessToken)

	type result struct {
		Email string
	}
	var r result
	err = api.RequestJson(req, &r)
	if err != nil {
		return "", err
	}
	if r.username == "" {
		return "", errors.New("no email")
	}
	return r.username, nil
}

func (p *COCOProvider) ValidateSessionState(s *SessionState) bool {
	return validateToken(p, s.AccessToken, getCOCOHeader(s.AccessToken))
}

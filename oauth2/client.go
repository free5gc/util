package oauth2

import (
	"context"
	"fmt"
	"net/http"
	"time"

	nrf_management "github.com/ShouheiNishi/openapi5g/nrf/management"
	nrf_token "github.com/ShouheiNishi/openapi5g/nrf/token"
	utils_error "github.com/ShouheiNishi/openapi5g/utils/error"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"golang.org/x/oauth2"

	"github.com/free5gc/util/generics"
	"github.com/free5gc/util/httpclient"
)

var tokenCache generics.SyncMap[string, *oauth2.Token]

func GetOauth2RequestEditor(ctx context.Context, nfType nrf_management.NFType, nfId uuid.UUID, nrfUri string, scope string, targetNF nrf_management.NFType,
) (editor func(ctx context.Context, req *http.Request) error, err error) {
	token, err := getOauth2Token(ctx, nfType, nfId, nrfUri, scope, targetNF)
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context, req *http.Request) error {
		token.SetAuthHeader(req)
		return nil
	}, nil
}

func getOauth2Token(ctx context.Context, nfType nrf_management.NFType, nfId uuid.UUID, nrfUri string, scope string, targetNF nrf_management.NFType) (token *oauth2.Token, err error) {
	if cacheEntry, exist := tokenCache.Load(scope); exist {
		if cacheEntry.Expiry.IsZero() || cacheEntry.Expiry.After(time.Now()) {
			return cacheEntry, nil
		}
	}

	client, err := nrf_token.NewClientWithResponses(nrfUri, nrf_token.WithHTTPClient(httpclient.GetHttpClient(nrfUri)))
	if err != nil {
		return nil, fmt.Errorf("nrf_token.NewClientWithResponses: %w", err)
	}

	res, err := client.AccessTokenRequestWithFormdataBodyWithResponse(ctx, &nrf_token.AccessTokenRequestParams{
		AcceptEncoding: lo.ToPtr("application/json"),
	}, nrf_token.AccessTokenReq{
		GrantType:    nrf_token.ClientCredentials,
		NfInstanceId: nfId,
		Scope:        scope,
		NfType:       &nfType,
		TargetNfType: &targetNF,
	})

	if tokenRes := res.JSON200; tokenRes != nil {
		newToken := oauth2.Token{
			AccessToken: tokenRes.AccessToken,
			TokenType:   string(tokenRes.TokenType),
		}
		if tokenRes.ExpiresIn != nil {
			newToken.Expiry = time.Now().Add(time.Duration(*tokenRes.ExpiresIn) * time.Second)
		}
		tokenCache.Store(scope, &newToken)

		return &newToken, nil
	}

	return nil, utils_error.ExtractAndWrapOpenAPIError("nrf_token.AccessTokenRequestWithFormdataBodyWithResponse", res, err)
}

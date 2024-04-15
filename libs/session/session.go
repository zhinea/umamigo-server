package session

import (
	"errors"
	"github.com/zhinea/umamigo-server/entity"
	detectgeo "github.com/zhinea/umamigo-server/libs"
	myjwt "github.com/zhinea/umamigo-server/libs/jwt"
	"github.com/zhinea/umamigo-server/libs/queries"
	"log"
	"regexp"
)

type PayloadData struct {
	Headers map[string][]string
	body    entity.RequestPayload
	IP      string
	IsLocal bool
}

var HostnameRegex, _ = regexp.Compile("^[\\w-.]+$")

func UseSession(payload PayloadData) (entity.SessionClaims, error) {

	cacheToken := payload.Headers["x-umami-cache"]

	if len(cacheToken) > 0 {
		log.Println("Using cache token")
		result, err := myjwt.ParseToken(cacheToken[0])

		if err != nil {
			log.Println(err)
		} else {
			return *result, nil
		}
	}

	body := payload.body

	if !HostnameRegex.MatchString(body.Hostname) {
		return entity.SessionClaims{}, errors.New("invalid hostname")
	}

	website := queries.FindWebsite(body.ID)

	if website.ID == "" {
		return entity.SessionClaims{}, errors.New("website not found")
	}

	client := detectgeo.GetClientInfo(payload)

	return entity.SessionClaims{}, nil
}

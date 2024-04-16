package session

import (
	"errors"
	"github.com/zhinea/umamigo-server/entity"
	detectgeo "github.com/zhinea/umamigo-server/libs"
	myjwt "github.com/zhinea/umamigo-server/libs/jwt"
	"github.com/zhinea/umamigo-server/libs/queries"
	"github.com/zhinea/umamigo-server/utils"
	"log"
	"regexp"
	"time"
)

var HostnameRegex, _ = regexp.Compile("^[\\w-.]+$")

func UseSession(payload *entity.UseSessionPayloadData) (entity.JWTSessionClaims, error) {

	cacheToken := payload.Headers["x-umami-cache"]

	if len(cacheToken) > 0 {
		parsingTimer := time.Now()

		result, err := myjwt.ParseToken(cacheToken[0])

		log.Println("main->session: Parsing took ", time.Now().Sub(parsingTimer).String())

		if err != nil {
			log.Println(err)
		} else {
			return *result, nil
		}
	}

	validationTimer := time.Now()

	body := payload.Body

	if !HostnameRegex.MatchString(body.Hostname) {
		return entity.JWTSessionClaims{}, errors.New("invalid hostname")
	}
	log.Println("main->session: [hostname_regex] verification", time.Now().Sub(validationTimer).String())

	website := queries.FindWebsite(body.ID)

	if website.WebsiteID == "" {
		return entity.JWTSessionClaims{}, errors.New("website not found")
	}

	log.Println("main->session: Validation took ", time.Now().Sub(validationTimer).String())

	clientGEOTimer := time.Now()

	client := detectgeo.GetClientInfo(payload)

	log.Println("main->session: Client geo took ", time.Now().Sub(clientGEOTimer).String())

	uuidGenerateTimer := time.Now()

	sessionID := utils.UUID(
		body.ID,
		body.Hostname,
		client.IP,
		client.UserAgent,
	)
	visitID := utils.UUID(
		sessionID,
		utils.VisitSalt(),
	)

	log.Println("main->session: UUID generation took ", time.Now().Sub(uuidGenerateTimer).String())

	getSessionTimer := time.Now()
	session := queries.FindSession(sessionID)

	if session.SessionID == "" {
		session = entity.Session{
			SessionID:    sessionID,
			WebsiteID:    body.ID,
			Hostname:     body.Hostname,
			Browser:      client.Browser,
			OS:           client.OS,
			Device:       client.Device,
			Screen:       body.Screen,
			Language:     body.Language,
			Country:      client.Country,
			Subdivision1: client.Subdivision1,
			Subdivision2: client.Subdivision2,
			City:         client.City,
		}
		queries.CreateSession(session)
	}

	log.Println("main->session: Get session took ", time.Now().Sub(getSessionTimer).String())

	return entity.JWTSessionClaims{
		SessionClaims: entity.SessionClaims{
			Session: session,
			OwnerID: website.UserID,
			VisitID: visitID,
		},
	}, nil
}

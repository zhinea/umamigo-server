package services

import (
	uuid "github.com/satori/go.uuid"
	"github.com/zhinea/umamigo-server/entity"
	"github.com/zhinea/umamigo-server/libs/queries"
)

func SaveEvent(data entity.EventCreationPayload) {
	runAsSQL(data)
}

func runAsSQL(data entity.EventCreationPayload) {
	website := data.WebsiteEvent

	websiteEventID := uuid.NewV4().String()

	website.EventID = websiteEventID
	website.WebsiteID = data.SessionClaims.WebsiteID
	website.SessionID = data.SessionClaims.SessionID
	website.VisitID = data.SessionClaims.VisitID

	queries.CreateEvent(website)
}

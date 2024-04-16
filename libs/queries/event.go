package queries

import (
	"github.com/zhinea/umamigo-server/entity"
	"github.com/zhinea/umamigo-server/libs/database"
)

func CreateEvent(website entity.WebsiteEvent) {
	database.DB.Table("website_event").Create(&website)
}

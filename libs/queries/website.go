package queries

import (
	"github.com/zhinea/umamigo-server/entity"
	"github.com/zhinea/umamigo-server/libs/database"
)

func FindWebsite(websiteId string) entity.Website {
	var website entity.Website

	database.DB.First(&website, "website_id = ?", websiteId)

	return website
}

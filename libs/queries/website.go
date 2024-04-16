package queries

import (
	"context"
	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
	"github.com/zhinea/umamigo-server/entity"
	"github.com/zhinea/umamigo-server/libs/database"
	"github.com/zhinea/umamigo-server/utils"
	"log"
	"time"
)

func MarshallWebsite(website entity.Website) string {
	val, err := json.Marshal(website)
	if err != nil {
		panic(err)
	}
	return string(val)
}

func FindWebsite(websiteId string) entity.Website {
	website := entity.Website{}

	dbQueryTimer := time.Now()
	if utils.Cfg.EnableCache {
		res, err := database.Redis.Get(database.Ctx, "website:"+websiteId).Result()

		if err == redis.Nil {
		} else if err != nil {
			panic(err)
		} else {
			err = json.Unmarshal([]byte(res), &website)
			if err != nil {
				panic(err)
			}

			log.Println("main->session->queries: Website query took [cached]", time.Now().Sub(dbQueryTimer).String())
			return website
		}
	}
	ctx := context.TODO()

	database.DB.Raw("SELECT * FROM website WHERE website_id = ?", websiteId).Scan(&website)

	database.Redis.Set(ctx, "website:"+websiteId, MarshallWebsite(website), 0)

	log.Println("main->session->queries: Website query took ", time.Now().Sub(dbQueryTimer).String())

	return website
}

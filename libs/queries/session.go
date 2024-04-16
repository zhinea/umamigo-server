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

func MarshallSession(session entity.Session) string {
	val, err := json.Marshal(session)
	if err != nil {
		panic(err)
	}
	return string(val)
}

func FindSession(sessionId string) entity.Session {
	session := entity.Session{}
	ctx := context.TODO()

	dbQueryTimer := time.Now()

	if utils.Cfg.EnableCache {
		res, err := database.Redis.Get(database.Ctx, "session:"+sessionId).Result()

		if err == redis.Nil {
		} else if err != nil {
			panic(err)
		} else {
			log.Println("using cache")
			err = json.Unmarshal([]byte(res), &session)
			if err != nil {
				panic(err)
			}

			log.Println("main->session->queries: Session query took [cached]", time.Now().Sub(dbQueryTimer).String())
			return session
		}

	}

	//database.DB.Table("session").Take(&session, "session_id = ?", sessionId)
	database.DB.Raw("SELECT * FROM session WHERE session_id = ?", sessionId).Scan(&session)

	err := database.Redis.Set(ctx, "session:"+sessionId, MarshallSession(session), 86400).Err()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("main->session->queries: Session query took ", time.Now().Sub(dbQueryTimer).String())

	return session
}

func CreateSession(session entity.Session) {
	database.DB.Table("session").Create(&session)
}

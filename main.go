package main

import (
	"flag"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/zhinea/umamigo-server/entity"
	"github.com/zhinea/umamigo-server/libs/database"
	libjwt "github.com/zhinea/umamigo-server/libs/jwt"
	"github.com/zhinea/umamigo-server/libs/services"
	"github.com/zhinea/umamigo-server/libs/session"
	"github.com/zhinea/umamigo-server/utils"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

func createBatchEvents(events []entity.EventRequest, session entity.JWTSessionClaims) {
	defer utils.Recover()

	for index := range events {
		ev := events[index]

		payload := ev.Payload

		if ev.Type == "event" {

			urlPathQuery := strings.Split(utils.DecodeURIComponent(payload.Url), "?")
			referrerPathQuery := strings.Split(utils.DecodeURIComponent(payload.Referrer), "?")

			urlPath := "/"
			urlQuery := ""
			if len(urlPathQuery) > 0 {
				urlPath = urlPathQuery[0]
				if len(urlPathQuery) > 1 {
					urlQuery = urlPathQuery[1]
				}
			}

			referrerPath := ""
			referrerQuery := ""
			referrerDomain := ""
			if len(referrerPathQuery) > 0 {
				referrerPath = referrerPathQuery[0]
				if len(referrerPathQuery) > 1 {
					referrerQuery = referrerPathQuery[1]
				}
			}

			if strings.HasPrefix(referrerPath, "http") {
				refUrl, _ := url.Parse(referrerPath)
				referrerPath = refUrl.Path
				referrerQuery = refUrl.RawQuery
				referrerDomain = strings.Replace(refUrl.Hostname(), "www.", "", 1)
			}

			if os.Getenv("REMOVE_TRAILING_SLASH") != "" {
				re := regexp.MustCompile("(.+)/$")
				urlPath = re.ReplaceAllString(urlPath, "$1")
			}

			eventType := 1

			if payload.Name != "" {
				eventType = 2
			}

			constant := utils.Cfg.Umami.Constants
			urlLength := constant.UrlLength

			t := time.Now()

			// prevent before 2023-12-21 16:00:33 +0700 WIB
			if payload.T > 1703149233975 {
				t = time.Unix(int64(payload.T/1000), 0)
			}

			services.SaveEvent(
				entity.EventCreationPayload{
					SessionClaims: session.SessionClaims,
					WebsiteEvent: entity.WebsiteEvent{
						UrlPath:        utils.Substr(urlPath, 0, urlLength),
						UrlQuery:       utils.Substr(urlQuery, 0, urlLength),
						ReferrerPath:   utils.Substr(referrerPath, 0, urlLength),
						ReferrerQuery:  utils.Substr(referrerQuery, 0, urlLength),
						ReferrerDomain: utils.Substr(referrerDomain, 0, urlLength),
						PageTitle:      utils.Substr(payload.Title, 0, constant.PageTitleLength),
						EventName:      utils.Substr(payload.Name, 0, constant.EventNameLength),
						EventType:      eventType,
						CreatedAt:      t,
					},
					T:         payload.T,
					EventData: payload.Data,
				},
			)
		}
	}
}

func main() {
	cfgFilename := flag.String("config", utils.GetEnvPath(), "Config file path.")

	flag.Parse()
	// load configuration file
	utils.LoadConfig(*cfgFilename)

	database.Connect()

	http.HandleFunc("/api/send", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		totalTime := time.Now()
		validationTimer := time.Now()

		if r.Method != "POST" {
			fmt.Fprintf(w, "Only support POST method")
			return
		}

		var payload = new(entity.RequestPayload)

		err := json.NewDecoder(r.Body).Decode(payload)

		if err != nil {
			log.Println(err)
			fmt.Fprintf(w, "Invalid body request")
			return
		}

		log.Println("main: validation timer took", time.Now().Sub(validationTimer).String())

		sessionTimer := time.Now()

		ip, err := utils.GetIP(r)

		if err != nil {
			log.Fatal("error parsing IP")
		}

		sess, err := session.UseSession(&entity.UseSessionPayloadData{
			Headers: r.Header,
			Body:    *payload,
			IP:      ip,
			IsLocal: utils.Cfg.Env == "development",
		})

		if err != nil {
			log.Println(err)
			w.Write([]byte("Error parsing session"))
			return
		}

		if sess.RegisteredClaims.IssuedAt != nil && time.Now().Unix()-sess.RegisteredClaims.IssuedAt.Unix() > 1800 {
			sess.VisitID = utils.UUID(sess.ID, utils.VisitSalt())
		}

		sess.RegisteredClaims.IssuedAt = &jwt.NumericDate{Time: time.Now()}

		log.Println("main: session timer took", time.Now().Sub(sessionTimer).String())

		go createBatchEvents(payload.Events, sess)

		jwtCreationTimer := time.Now()
		jwtToken := libjwt.CreateToken(sess)
		log.Println("main: jwt creation took", time.Now().Sub(jwtCreationTimer).String())
		log.Println("main: total time took", time.Now().Sub(totalTime).String())
		log.Println("============")

		w.Write([]byte(jwtToken))
	})

	fmt.Println("starting web server at http://localhost:3022/")
	err := http.ListenAndServe(":3022", nil)
	if err != nil {
		log.Fatal(err)
		return
	}
}

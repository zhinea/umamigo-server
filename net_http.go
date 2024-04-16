package main

import (
	"flag"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/zhinea/umamigo-server/entity"
	"github.com/zhinea/umamigo-server/libs/database"
	libjwt "github.com/zhinea/umamigo-server/libs/jwt"
	"github.com/zhinea/umamigo-server/libs/session"
	"github.com/zhinea/umamigo-server/utils"
	"log"
	"net/http"
	"time"
)

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
			IsLocal: false,
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

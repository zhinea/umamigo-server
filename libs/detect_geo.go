package detectgeo

import (
	"github.com/mileusna/useragent"
	"github.com/oschwald/maxminddb-golang"
	"github.com/zhinea/umamigo-server/entity"
	"github.com/zhinea/umamigo-server/utils"
	"log"
	"net"
	"path"
	"time"
)

var mindDBInialized bool = false
var mindDB *maxminddb.Reader

func getLocation(payload *entity.UseSessionPayloadData) entity.GeoLocation {
	// ignore local ips
	if payload.IsLocal {
		return entity.GeoLocation{}
	}

	headers := payload.Headers

	if headers["cf-ipcountry"] != nil {
		return entity.GeoLocation{
			Country:      utils.DecodeURIComponent(utils.SoftTouch(headers["cf-ipcountry"])),
			Subdivision1: utils.DecodeURIComponent(utils.SoftTouch(headers["cf-region-code"])),
			City:         utils.DecodeURIComponent(utils.SoftTouch(headers["cf-city"])),
		}
	}

	if headers["x-vercel-ip-country"] != nil {
		return entity.GeoLocation{
			Country:      utils.DecodeURIComponent(utils.SoftTouch(headers["x-vercel-ip-country"])),
			Subdivision1: utils.DecodeURIComponent(utils.SoftTouch(headers["x-vercel-ip-country-region"])),
			City:         utils.DecodeURIComponent(utils.SoftTouch(headers["x-vercel-ip-city"])),
		}
	}

	if !mindDBInialized {
		log.Println("Initializing maxminddb")
		dir := path.Join(".", "geo")

		var err error
		mindDB, err = maxminddb.Open(path.Join(dir, "GeoLite2-City.mmdb"))

		if err != nil {
			log.Println(err)
		}
		mindDBInialized = true
	}

	ip := net.ParseIP(payload.IP)

	result := entity.GeoLocation{}

	err := mindDB.Lookup(ip, &result)
	if err != nil {
		return entity.GeoLocation{}
	}

	return result
}

func getAgent(userAgent string) entity.GeoAgent {
	b := useragent.Parse(userAgent)

	return entity.GeoAgent{
		Browser: b.Name,
		OS:      b.OS,
		Device:  b.Device,
	}
}

func GetClientInfo(payload *entity.UseSessionPayloadData) entity.GeoClientInfo {

	userAgentParsingTime := time.Now()
	userAgent := utils.SoftTouch(payload.Headers["User-Agent"])
	log.Println("main->session->geo: User agent parsing took ", time.Now().Sub(userAgentParsingTime).String())

	ip := payload.IP

	locationTimer := time.Now()
	location := getLocation(payload)

	country := location.Country
	subdivision1 := location.Subdivision1
	subdivision2 := location.Subdivision2
	city := location.City
	log.Println("main->session->geo: Location took ", time.Now().Sub(locationTimer).String())

	agentTimer := time.Now()
	agent := getAgent(userAgent)
	log.Println("main->session->geo: Agent parsing took ", time.Now().Sub(agentTimer).String())

	return entity.GeoClientInfo{
		UserAgent:    userAgent,
		Browser:      agent.Browser,
		OS:           agent.OS,
		IP:           ip,
		Country:      country,
		Subdivision1: subdivision1,
		Subdivision2: subdivision2,
		City:         city,
		Device:       agent.Device,
	}
}

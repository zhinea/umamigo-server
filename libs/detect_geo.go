package detectgeo

import (
	"github.com/dineshgowda24/browser"
	"github.com/oschwald/maxminddb-golang"
	"github.com/zhinea/umamigo-server/libs/session"
	"github.com/zhinea/umamigo-server/utils"
	"log"
	"net"
	"path"
)

type Location struct {
	Country      string `json:"country,omitempty"`
	Subdivision1 string `json:"subdivision1,omitempty"`
	Subdivision2 string `json:"subdivision2,omitempty"`
	City         string `json:"city,omitempty"`
}

type Agent struct {
	Browser string `json:"browser,omitempty"`
	OS      string `json:"os,omitempty"`
	Device  string `json:"device,omitempty"`
}

type ClientInfo struct {
	UserAgent    string `json:"user_agent,omitempty"`
	Browser      string `json:"browser,omitempty"`
	OS           string `json:"os,omitempty"`
	IP           string `json:"ip,omitempty"`
	Country      string `json:"country,omitempty"`
	Subdivision1 string `json:"subdivision1,omitempty"`
	Subdivision2 string `json:"subdivision2,omitempty"`
	City         string `json:"city,omitempty"`
	Device       string `json:"device,omitempty"`
}

var mindDBInialized bool = false
var mindDB maxminddb.Reader

func getLocation(payload session.PayloadData) Location {
	// ignore local ips
	if payload.IsLocal {
		return Location{}
	}

	headers := payload.Headers

	if headers["cf-ipcountry"] != nil {
		return Location{
			Country:      utils.DecodeURIComponent(headers["cf-ipcountry"][0]),
			Subdivision1: utils.DecodeURIComponent(headers["cf-region-code"][0]),
			City:         utils.DecodeURIComponent(headers["cf-city"][0]),
		}
	}

	if headers["x-vercel-ip-country"] != nil {
		return Location{
			Country:      utils.DecodeURIComponent(headers["x-vercel-ip-country"][0]),
			Subdivision1: utils.DecodeURIComponent(headers["x-vercel-ip-country-region"][0]),
			City:         utils.DecodeURIComponent(headers["x-vercel-ip-city"][0]),
		}
	}

	if !mindDBInialized {
		dir := path.Join(".", "geo")

		var err error
		mindDB, err = maxminddb.Open(path.Join(dir, "GeoLite2-City.mmdb"))

		if err != nil {
			log.Println(err)
		}
		mindDBInialized = true
	}

	ip := net.ParseIP(payload.IP)

	result := Location{}

	err := mindDB.Lookup(ip, &result)
	if err != nil {
		return Location{}
	}

	return result
}

func getAgent(userAgent string) Agent {
	b, err := browser.NewBrowser(userAgent)

	if err != nil {
		log.Println(err)
		return Agent{}
	}

	return Agent{
		Browser: b.Name(),
		OS:      b.Platform().Name(),
		Device:  b.Device().Name(),
	}
}

func GetClientInfo(payload session.PayloadData) ClientInfo {
	userAgent := payload.Headers["user-agent"]
	ip := payload.IP
	location := getLocation(payload)

	country := location.Country
	subdivision1 := location.Subdivision1
	subdivision2 := location.Subdivision2
	city := location.City

	agent := getAgent(userAgent[0])

	return ClientInfo{
		UserAgent:    userAgent[0],
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

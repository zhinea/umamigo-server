package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

func main() {
	if os.Getenv("VERCEL") != "" {
		fmt.Println("Vercel environment detected. Skipping geo setup.")
		os.Exit(0)
	}

	log.Println("Downloading geo database...")

	db := "GeoLite2-City"
	url := fmt.Sprintf("https://raw.githubusercontent.com/GitSquared/node-geolite2-redist/master/redist/%s.tar.gz", db)
	if licenseKey := os.Getenv("MAXMIND_LICENSE_KEY"); licenseKey != "" {
		url = fmt.Sprintf("https://download.maxmind.com/app/geoip_download?edition_id=%s&license_key=%s&suffix=tar.gz", db, licenseKey)
	}

	dest := path.Join(".", "geo")
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		os.Mkdir(dest, 0755)
	}

	if _, err := os.Stat(path.Join(dest, db+".mmdb")); err == nil {
		fmt.Println("Geo database already exists.")
		os.Exit(0)
	}

	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Error downloading database:", err)
		os.Exit(1)
	}
	defer res.Body.Close()

	gzr, err := gzip.NewReader(res.Body)
	if err != nil {
		fmt.Println("Error decompressing gzip:", err)
		os.Exit(1)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading tar header:", err)
			os.Exit(1)
		}

		if header.Typeflag == tar.TypeReg && path.Ext(header.Name) == ".mmdb" {
			filename := path.Join(dest, path.Base(header.Name))
			file, err := os.Create(filename)
			if err != nil {
				fmt.Println("Error creating file:", err)
				os.Exit(1)
			}
			defer file.Close()

			if _, err := io.Copy(file, tr); err != nil {
				fmt.Println("Error writing file:", err)
				os.Exit(1)
			}

			fmt.Println("Saved geo database:", filename)
		}
	}
}

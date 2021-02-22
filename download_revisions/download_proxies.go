package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type proxy struct {
	Name string `json:"name"`
}

type proxyList struct {
	List []proxy `json:"proxies"`
}

var baseURL string = "https://apigee.googleapis.com/v1"

func main() {

	project := flag.String("project-id", "", "GCP Project ID")
	token := flag.String("token", "", "Get with gcloud auth print-access-token")
	destinationFolder := flag.String("dest", "./", "destination folder for zip files.")

	flag.Parse()

	reqURL := baseURL + "/organizations/" + *project + "/apis"

	req, _ := http.NewRequest(
		"GET",
		reqURL,
		nil,
	)

	req.Header.Add("Authorization", "Bearer "+*token)
	req.Header.Add("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("Error: ", err)
	}

	data, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	fmt.Printf("body: %s\n", data)

	list := proxyList{
		List: []proxy{},
	}

	if err = json.Unmarshal(data, &list); err != nil {
		log.Fatal("Error: ", err)
	}

	for _, p := range list.List {
		reqURL = baseURL + "/organizations/" + *project + "/apis/" + p.Name + "/revisions"

		req, _ = http.NewRequest(
			"GET",
			reqURL,
			nil,
		)
		req.Header.Add("Authorization", "Bearer "+*token)

		res, err = http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal("Error: ", err)
		}

		data, _ = ioutil.ReadAll(res.Body)
		res.Body.Close()

		latestVersion := getLatestVersion(string(data))
		fmt.Println(latestVersion)

		reqURL = baseURL + "/organizations/" + *project + "/apis/" + p.Name + "/revisions/" + strconv.Itoa(latestVersion) + "?format=bundle"

		req, _ = http.NewRequest(
			"GET",
			reqURL,
			nil,
		)
		req.Header.Add("Authorization", "Bearer "+*token)
		req.Header.Add("Accept", "application/zip")

		res, err = http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal("Error: ", err)
		}

		data, _ = ioutil.ReadAll(res.Body)
		res.Body.Close()

		ioutil.WriteFile(*destinationFolder+p.Name+"-v"+strconv.Itoa(latestVersion)+".zip", data, 0777)

	}

}

func getLatestVersion(v string) int {
	re := regexp.MustCompile(`"[0-9]+"`)
	versions := re.FindAllString(v, -1)
	latest := 1
	for _, version := range versions {
		versionInt, _ := strconv.Atoi(strings.Trim(version, "\""))
		if versionInt > latest {
			latest = versionInt
		}
	}
	return latest
}

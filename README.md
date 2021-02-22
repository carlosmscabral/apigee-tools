# apigee-tools

## download_revisions

Simple [golang](https://golang.org) utility for downloading the latest revisions of all the API proxies in your Apigee org. Tested for the https://apigee.googleapis.com/v1 endpoint (Apigee X and Apigee Hybrid).

Assumes you have [golang](https://golang.org) installed in your machine, an Org/GCP Project created and access to a GCP token with all the necessary permissions.

### Usage:

go run download_proxies.go -dest _string_ -project-id _string_ -token _string_

Parameter explanation:

  

        -dest: destination folder for zip files. (default "./")
  

        -project-id: GCP Project ID *MANDATORY*
  

        -token: Get with gcloud auth print-access-token *MANDATORY*
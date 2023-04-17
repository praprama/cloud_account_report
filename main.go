package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type DomainsResult struct {
	Azure   AzureDomains	`json:"azure"`
	GCP 	GCPDomains		`json:"gcp"`
}

func main() {
	domainPtr := flag.String("domain", "", "the domain name to query")
	retryPtr := flag.Int("retry", 3, "the number of times to retry on error")
	outPtr := flag.String("output", "text", "specify the output format. Acceptable values 'text' and 'json'")
	flag.Parse()

	domain := *domainPtr
	if domain == "" {
		fmt.Println("Please provide a domain name with the -domain flag")
		flag.PrintDefaults()
		os.Exit(1)
	}

	retry := *retryPtr
	if retry < 0 {
		fmt.Println("Retry count cannot be negative")
		flag.PrintDefaults()
		os.Exit(1)
	}

	output := *outPtr
	if output != "text" && output != "json" {
		fmt.Println("Invalid value for output. Can be either 'text' or 'json'")
		flag.PrintDefaults()
		os.Exit(1)
	}

	url := fmt.Sprintf("https://crt.sh/?q=%s&output=json", domain)

	var certificates []map[string]interface{}
	var resp *http.Response
	var err error

	for i := 0; i <= retry; i++ {
		resp, err = http.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			break
		}
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("Error: %s\n", resp.Status)
		}
		if i < retry {
			fmt.Printf("Retrying in 5 seconds (attempt %d/%d)\n", i+1, retry)
			time.Sleep(5 * time.Second)
		}
	}
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to retrieve certificate information")
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	err = json.Unmarshal(body, &certificates)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	var commonNames []string
	for _, cert := range certificates {
		commonNames = append(commonNames, strings.Split(cert["name_value"].(string),"\n")...)
    commonNames = append(commonNames, cert["common_name"].(string))
	}

	uniqueCommonNames := unique(commonNames)

	az_data := az_enum(uniqueCommonNames, output)
	gcp_data := gcp_enum(uniqueCommonNames, output)

	if output == "json" {
		domainsResult := DomainsResult{
			Azure:   az_data,
			GCP: gcp_data,
		}
	
		domainsResultJSON, err := json.Marshal(domainsResult)
		if err != nil {
			fmt.Printf("Error marshaling domainsResult to JSON: %v\n", err)
			return
		}
	
		// Print the JSON string
		fmt.Println(string(domainsResultJSON))
	}
}

func unique(items []string) []string {
	encountered := map[string]bool{}
	var result []string
	for _, item := range items {
    if strings.Contains(item, ".") {
      if strings.HasPrefix(item, "*") {
        item = item[2:]
      }
  		if !encountered[item] {
  			encountered[item] = true
  			result = append(result, item)
  		}
    }
	}
	return result
}

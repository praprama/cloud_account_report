package main

import (
	//"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

var (
	gcp_mu_a    sync.Mutex
	gcp_available []string
)

var (
	gcp_mu_b    sync.Mutex
	gcp_unavailable []string
)

type GCPDomains struct {
	GCP_Available   []string `json:"available"`
	GCP_Unavailable []string `json:"unavailable"`
}

func gcp_enum(domains []string, output string) (GCPDomains){
	var wg sync.WaitGroup
	wg.Add(len(domains))

	for _, domain := range domains {
		go func(d string) {
			defer wg.Done()
			url := fmt.Sprintf("https://accounts.google.com/samlredirect?domain=%s", d)
			resp, err := http.Get(url)
			if err != nil {
				fmt.Printf("Error checking domain %s: %v\n", d, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusBadRequest {
				gcp_mu_b.Lock()
				defer gcp_mu_b.Unlock()
				gcp_unavailable = append(gcp_unavailable, d)
			} else {
				gcp_mu_a.Lock()
				defer gcp_mu_a.Unlock()
				gcp_available = append(gcp_available, d)
			}
		}(domain)
	}

	wg.Wait()

	if output == "text" {
		fmt.Println("Following domains are found to be managed on GCP")
		fmt.Println("================================================")
		fmt.Println(strings.Join(gcp_available, ", "))
		fmt.Println()
		fmt.Println()
		fmt.Println("Following domains are found to not be managed on GCP")
		fmt.Println("====================================================")
		fmt.Println(strings.Join(gcp_unavailable, ", "))
		fmt.Println()
		fmt.Println()
		fmt.Println("Done!")
	}

	gcpDomains := GCPDomains{
		GCP_Available:   gcp_available,
		GCP_Unavailable: gcp_unavailable,
	}

	return gcpDomains

	// Marshal the struct to JSON
	/*gcpDomainsJSON, err := json.Marshal(gcpDomains)
	if err != nil {
		fmt.Printf("Error marshaling Azure domains to JSON: %v\n", err)
		return
	}

	// Print the JSON string
	fmt.Println(string(gcpDomainsJSON))*/
}

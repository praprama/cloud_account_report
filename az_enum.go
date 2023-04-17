package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

var (
	az_mu_a    sync.Mutex
	az_available []string
)

var (
	az_mu_b    sync.Mutex
	az_unavailable []string
)

type AzureDomains struct {
	AZ_Available   []string `json:"available"`
	AZ_Unavailable []string `json:"unavailable"`
}

func az_enum(domains []string, output string) (AzureDomains){
	//domains = []string{"example.com", "google.com", "yahoo.com", "getfi.com"}

	var wg sync.WaitGroup
	wg.Add(len(domains))

	for _, domain := range domains {
		go func(d string) {
			defer wg.Done()
			resp, err := http.Get(fmt.Sprintf("https://login.microsoftonline.com/getuserrealm.srf?login=%s&json=1", d))
			if err != nil {
				fmt.Printf("Error for domain %s: %v\n", d, err)
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				fmt.Printf("Error for domain %s: received status code %d\n", d, resp.StatusCode)
				return
			}
			var raw json.RawMessage
			if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
				fmt.Printf("Error decoding JSON for domain %s: %v\n", d, err)
				return
			}

			var domains map[string]interface{}

			err = json.Unmarshal(raw, &domains)

			if err != nil {
				fmt.Printf("Error extracting State for domain %s: %v\n", d, err)
			}

			//fmt.Println(domains["NameSpaceType"])
			if domains["NameSpaceType"] != "Unknown" {
				az_mu_a.Lock()
				defer az_mu_a.Unlock()
				az_available = append(az_available, domains["Login"].(string))
			} else {
				az_mu_b.Lock()
				defer az_mu_b.Unlock()
				az_unavailable = append(az_unavailable, domains["Login"].(string))
			}
			
		}(domain)
	}

	wg.Wait()

	if output == "text" {
		fmt.Println("Following domains are found to be managed on Azure")
		fmt.Println("==================================================")
		fmt.Println(strings.Join(az_available, ", "))
		fmt.Println()
		fmt.Println()
		fmt.Println("Following domains are found to not be managed on Azure")
		fmt.Println("======================================================")
		fmt.Println(strings.Join(az_unavailable, ", "))
		fmt.Println()
		fmt.Println()
		fmt.Println("Done!")
	}

	azureDomains := AzureDomains{
		AZ_Available:   az_available,
		AZ_Unavailable: az_unavailable,
	}

	return azureDomains

	// Marshal the struct to JSON
	/*azureDomainsJSON, err := json.Marshal(azureDomains)
	if err != nil {
		fmt.Printf("Error marshaling Azure domains to JSON: %v\n", err)
		return azureDomains
	}*/
}

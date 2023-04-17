# Cloud Account Report
For a given organization name/pattern, get all unique domain names from crt.sh and then check if those are managed domains on Azure and GCP.

## Download
Download the file over here - [cloud_account_report]https://github.com/praprama/cloud_account_report/raw/main/bin/cloud_account_report)

## Usage
Run the binary downloaded with the options below. `-domain` is the only mandatory parameter.
```
% ./cloud_account_report
Please provide a domain name with the -domain flag
  -domain string
    	the domain name to query
  -output string
    	specify the output format. Acceptable values 'text' and 'json' (default "text")
  -retry int
    	the number of times to retry on error (default 3)
```

Data returned as JSON has the following format:
```
{
  "azure": {
    "available": [
        #list of Azure managed domains found on crt.sh.
    ],
    "unavailable": [
        #list of domains found on crt.sh that are not Azure managed.
    ]
  },
  "gcp": {
    "available": [
        #list of GCP managed domains found on crt.sh.
    ],
    "unavailable": [
        #list of domains found on crt.sh that are not GCP managed.
    ]
  }
}
```

The arrays `available` and `unavailable` will contain `null` if no domains are found part of that list.

## Limitations
For GCP, the utility only checks if the domain is a federated one or not. If it is a locally managed domain that is not federated on GCP, that is not checked today.
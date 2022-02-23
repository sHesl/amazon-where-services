package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	e := botoEndpoints()

	globalServices(e)
	byRegion(e)
	byPartition(e)
	byService(e)
}

type endpoints struct {
	Partitions []struct {
		Name     string `json:"partition"`
		Services map[string]struct {
			Endpoints map[string]struct {
				Deprecated *bool `json:"deprecated,omitempty"`
			} `json:"endpoints"`
			IsRegionalized *bool `json:"isRegionalized,omitempty"`
		} `json:"services"`
	} `json:"partitions"`
}

func botoEndpoints() endpoints {
	url := "https://raw.githubusercontent.com/boto/botocore/develop/botocore/data/endpoints.json"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("unable to download endpoints from botocore repo. %s", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("unable to read endpoints data from botocore repo. %s", err)
	}

	var e endpoints
	if err := json.Unmarshal(b, &e); err != nil {
		log.Fatalf("unable to read endpoints data from botocore repo. %s", err)
	}

	return e
}

func globalServices(e endpoints) {
	combined := make(map[string]map[string]map[string][]string)
	for _, s := range []string{"global-services", "regional-services", "single-region-services"} {
		combined[s] = make(map[string]map[string][]string)
	}

	for _, p := range e.Partitions {
		for serviceName, s := range p.Services {
			dir := "regional-services"
			if s.IsRegionalized != nil && !*s.IsRegionalized {
				dir = "global-services"
			}

			if len(s.Endpoints) == 1 && dir != "global-services" {
				dir = "single-region-services"
			}

			for e, ep := range s.Endpoints {
				if ep.Deprecated != nil && *ep.Deprecated {
					continue
				}

				e = resolvePseudoGlobalRegion(e)
				if combined[dir][p.Name] == nil {
					combined[dir][p.Name] = make(map[string][]string, 0)
				}

				if combined[dir][p.Name][e] == nil {
					combined[dir][p.Name][e] = make([]string, 0)
				}
				combined[dir][p.Name][e] = append(combined[dir][p.Name][e], serviceName)
			}
		}

		for dir, partitions := range combined {
			for partition, regions := range partitions {
				for region, services := range regions {
					f := filepath.Join(".", dir, partition)
					os.MkdirAll(f, os.ModePerm)
					os.WriteFile(f+"/"+region, prepFileInput(services), os.ModePerm)
				}
			}
		}
	}
}

func byRegion(e endpoints) {
	f := filepath.Join("./regions")
	os.MkdirAll(f, os.ModePerm)

	combined := make(map[string][]string)

	for _, p := range e.Partitions {
		for serviceName, s := range p.Services {
			for e := range s.Endpoints {
				e = resolvePseudoGlobalRegion(e)
				if combined[e] == nil {
					combined[e] = make([]string, 0)
				}
				combined[e] = append(combined[e], serviceName)
			}
		}
	}

	for region, services := range combined {
		os.WriteFile(f+"/"+region, prepFileInput(services), os.ModePerm)
	}
}

func byPartition(e endpoints) {
	f := filepath.Join("./partitions")
	os.MkdirAll(f, os.ModePerm)

	combined := make(map[string][]string)

	for _, p := range e.Partitions {
		for serviceName := range p.Services {
			combined[p.Name] = append(combined[p.Name], serviceName)
		}
		os.WriteFile(f+"/"+p.Name, prepFileInput(combined[p.Name]), os.ModePerm)
	}
}

func byService(e endpoints) {
	combined := make(map[string][]string)

	for _, p := range e.Partitions {
		for serviceName, s := range p.Services {
			for e := range s.Endpoints {
				e = resolvePseudoGlobalRegion(e)
				combined[serviceName] = append(combined[serviceName], e)
			}
		}
	}

	for service, regions := range combined {
		f := filepath.Join("./services")
		os.MkdirAll(f, os.ModePerm)
		os.WriteFile(f+"/"+service, prepFileInput(regions), os.ModePerm)
	}
}

// https://github.com/aws/aws-sdk-js-v3/blob/main/clients/client-iam/src/endpoints.ts
var psuedoRegionMap = map[string]string{
	"aws-global":             "us-east-1",
	"fips-aws-global":        "fips-us-east-1",
	"aws-global-fips":        "fips-us-east-1",
	"aws-cn-global":          "cn-north-1",
	"aws-us-gov-global":      "us-gov-west-1",
	"fips-aws-us-gov-global": "fips-us-gov-west-1",
	"aws-iso-global":         "us-iso-east-1",
	"aws-iso-b-global":       "us-isob-east-1",
}

func resolvePseudoGlobalRegion(r string) string {
	if swap, exists := psuedoRegionMap[r]; exists {
		return swap
	}

	return r
}

func prepFileInput(input []string) []byte {
	deduped := make(map[string]struct{})
	for _, s := range input {
		deduped[s] = struct{}{}
	}

	output := make([]string, 0)
	for k := range deduped {
		output = append(output, k)
	}

	sort.Strings(output)

	return []byte(strings.Join(output, "\n"))
}

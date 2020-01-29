package cfapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type cfzone struct {
	ID     string
	Name   string
	Status string
}

type cfrecord struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Type    string `json:"type"`
}

type cftransport struct {
	email string
	token string
	base  http.RoundTripper
}

func (t *cftransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("X-Auth-Email", t.email)
	req.Header.Add("X-Auth-Key", t.token)
	req.Header.Add("Content-Type", "application/json")
	base := t.base
	if base == nil {
		base = http.DefaultTransport
	}
	return base.RoundTrip(req)
}

type cfclient struct {
	client *http.Client
}

func NewCFClient(email string, token string) *cfclient {
	return &cfclient{
		client: &http.Client{
			Transport: &cftransport{
				email: email,
				token: token,
			},
		},
	}
}

func (c *cfclient) FindZoneByName(name string) (*cfzone, error) {
	resp, err := c.client.Get("https://api.cloudflare.com/client/v4/zones?name=" + name)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var zoneResult struct {
		Result  []cfzone
		Success bool
		Errors  []interface{}
	}
	err = json.Unmarshal(body, &zoneResult)
	if err != nil {
		return nil, err
	}
	if !zoneResult.Success {
		return nil, fmt.Errorf("failed to find zone for domain %v - %v", name, zoneResult.Errors)
	}

	if len(zoneResult.Result) == 0 {
		return nil, fmt.Errorf("zone %v not found", name)
	}

	return &zoneResult.Result[0], nil
}

func (c *cfclient) FindRecordByName(zone string, name string) (*cfrecord, error) {
	resp, err := c.client.Get("https://api.cloudflare.com/client/v4/zones/" + zone + "/dns_records?name=" + name)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var recordResult struct {
		Result  []cfrecord
		Success bool
		Errors  []interface{}
	}
	err = json.Unmarshal(body, &recordResult)
	if err != nil {
		return nil, err
	}
	if !recordResult.Success {
		return nil, fmt.Errorf("failed to find record %v for zone %v - %v", name, zone, recordResult.Errors)
	}

	if len(recordResult.Result) == 0 {
		return nil, fmt.Errorf("record %v not found for zone %v", name, zone)
	}

	return &recordResult.Result[0], nil
}

func (c *cfclient) UpdateRecord(zone string, newRecord *cfrecord) error {
	updateBytes, err := json.Marshal(newRecord)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(
		http.MethodPut,
		"https://api.cloudflare.com/client/v4/zones/"+zone+"/dns_records/"+newRecord.ID,
		bytes.NewReader(updateBytes),
	)
	if err != nil {
		return err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var updatedRecord struct {
		Result  cfrecord
		Success bool
		Errors  []interface{}
	}

	err = json.Unmarshal(body, &updatedRecord)
	if err != nil {
		return err
	}
	if !updatedRecord.Success {
		return fmt.Errorf("error update %v - %v", newRecord, updatedRecord.Errors)
	}

	return nil
}

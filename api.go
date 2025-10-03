package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const LimaCityBaseURL = "https://www.lima-city.de/usercp"

type DNSRecordResponse struct {
	Records []Record `json:"records"`
}

type Record struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Subdomain any    `json:"subdomain"` // string if record is subdomain, otherwise false
	Type      string `json:"type"`
	Content   string `json:"content"`
	Priority  int    `json:"priority"`
	TTL       int    `json:"ttl"`
}

type DNSRecordRequest struct {
	NameserverRecord NameserverRecord `json:"nameserver_record"`
}

type NameserverRecord struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Content  string `json:"content"`
	Priority string `json:"priority"`
	TTL      int    `json:"ttl"`
}

type StatusReponse struct {
	Status string `json:"status"`
}

func SendAPIRequest[R any, S any](
	httpMethod string,
	url string,
	token string,
	body *R,
) (*S, error) {
	req, err := http.NewRequest(httpMethod, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Method = httpMethod
	req.SetBasicAuth("api", token)
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req.Body = io.NopCloser(bytes.NewBuffer(bodyJSON))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	// NopCloser will never return a non-nill error
	//nolint:errcheck
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP request failed with status code %d", res.StatusCode)
	}
	if res.Body == nil {
		// In this case, "nil"/an empty body is not an "invalid" return value
		//nolint:nilnil
		return nil, nil
	}

	response := new(S)
	err = json.NewDecoder(res.Body).Decode(response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func GetRecords(
	domainID int,
	token string,
) (*DNSRecordResponse, error) {
	res, err := SendAPIRequest[any, DNSRecordResponse](
		"GET",
		fmt.Sprintf("%s/domains/%d/records.json", LimaCityBaseURL, domainID),
		token,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func CreateRecord(
	name string,
	recordType string,
	content string,
	priority string,
	ttl int,
	domainID int,
	token string,
) (*StatusReponse, error) {
	ns := NameserverRecord{
		Name:     name,
		Type:     recordType,
		Content:  content,
		Priority: priority,
		TTL:      ttl,
	}
	req := DNSRecordRequest{NameserverRecord: ns}

	res, err := SendAPIRequest[DNSRecordRequest, StatusReponse](
		"POST",
		fmt.Sprintf("%s/domains/%d/records.json", LimaCityBaseURL, domainID),
		token,
		&req,
	)
	if err != nil {
		return nil, err
	}
	if res.Status != "ok" {
		return nil, fmt.Errorf("Record create request failed with status %s", res.Status)
	}
	return res, nil
}

func UpdateRecord(
	name string,
	recordType string,
	content string,
	priority string,
	ttl int,
	domainID int,
	recordID int,
	token string,
) (*StatusReponse, error) {
	ns := NameserverRecord{
		Name:     name,
		Type:     recordType,
		Content:  content,
		Priority: priority,
		TTL:      ttl,
	}
	req := DNSRecordRequest{NameserverRecord: ns}

	res, err := SendAPIRequest[DNSRecordRequest, StatusReponse](
		"PUT",
		fmt.Sprintf("%s/domains/%d/records/%d", LimaCityBaseURL, domainID, recordID),
		token,
		&req,
	)
	if err != nil {
		return nil, err
	}
	if res.Status != "ok" {
		return nil, fmt.Errorf("Record update request failed with status %s", res.Status)
	}
	return res, nil
}

func DeleteRecord(
	domainID string,
	recordID int,
	token string,
) (*StatusReponse, error) {
	res, err := SendAPIRequest[any, StatusReponse](
		"DELETE",
		fmt.Sprintf("%s/domains/%s/records/%d", LimaCityBaseURL, domainID, recordID),
		token,
		nil,
	)
	if err != nil {
		return nil, err
	}
	if res.Status != "ok" {
		return nil, fmt.Errorf("Record delete request failed with status %s", res.Status)
	}
	return res, nil
}

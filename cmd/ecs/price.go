package ecs

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const instancePriceURL = "https://g.alicdn.com/aliyun/ecs-price-info-intl/2.0.7/price" +
	"/download/instancePrice.json"

// InstancePrice contains price of an instance type
type InstancePrice struct {
	RegionID     string
	InstanceType string
	OSType       string
	PricePerHour float64
}

// Price contains detailed pricing info of a pricing document
type Price struct {
	Currency        string                 `json:"currency"`
	Version         string                 `json:"version"`
	PublicationDate string                 `json:"publicationDate"`
	PricingInfo     map[string]PricingInfo `json:"pricingInfo"`
	Disclaimer      string                 `json:"disclaimer"`
	Type            string                 `json:"type"`
	Site            string                 `json:"site"`
	Description     string                 `json:"description"`
}

// PricingInfo contains pricing details
type PricingInfo struct {
	Hours  []Hour `json:"hours"`
	Months []Hour `json:"months"`
	Years  []Hour `json:"years"`
}

// Hour contains price per hour
type Hour struct {
	Price  string `json:"price"`
	Period string `json:"period"`
}

// QueryPriceList returns price list for all instances type from all region
func (c *Client) QueryPriceList() ([]InstancePrice, error) {
	priceList := []InstancePrice{}

	resp, err := http.Get(instancePriceURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data Price
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	for k, v := range data.PricingInfo {
		// example key: ap-southeast-2::ecs.sn1.medium::vpc::windows::optimized
		s := strings.Split(k, "::")
		hourlyPrice, err := strconv.ParseFloat(v.Hours[0].Price, 64)
		if err != nil {
			continue
		}

		priceList = append(priceList, InstancePrice{
			RegionID:     s[0],
			InstanceType: s[1],
			OSType:       s[3],
			PricePerHour: hourlyPrice,
		})
	}

	return priceList, nil
}

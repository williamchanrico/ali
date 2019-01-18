package ecs

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryPriceList(t *testing.T) {
	tc := []InstancePrice{
		InstancePrice{
			RegionID:     "ap-southeast-1",
			InstanceType: "ecs.t5-c1m1.xlarge",
			OSType:       "linux",
			PricePerHour: 0.09,
		},
		InstancePrice{
			RegionID:     "ap-southeast-1",
			InstanceType: "ecs.r5.6xlarge",
			OSType:       "linux",
			PricePerHour: 1.92,
		},
	}

	tcJSON := []byte(`
		{
		  "currency": "USD",
		  "version": "2.0.7",
		  "publicationDate": "2018-10-08T08:15:52Z",
		  "pricingInfo": {
			"ap-southeast-1::ecs.t5-c1m1.xlarge::vpc::linux::optimized": {
			  "hours": [
				{
				  "price": "0.09",
				  "period": "1"
				}
			  ],
			  "months": [
				{
				  "price": "45.74",
				  "period": "1"
				}
			  ],
			  "years": [
				{
				  "price": "466.56",
				  "period": "1"
				},
				{
				  "price": "933.1",
				  "period": "2"
				},
				{
				  "price": "1399.64",
				  "period": "3"
				},
				{
				  "price": "1866.18",
				  "period": "4"
				},
				{
				  "price": "2332.74",
				  "period": "5"
				}
			  ]
			},
			"ap-southeast-1::ecs.r5.6xlarge::vpc::linux::optimized": {
			  "hours": [
				{
				  "price": "1.92",
				  "period": "1"
				}
			  ],
			  "months": [
				{
				  "price": "883.01",
				  "period": "1"
				}
			  ],
			  "years": [
				{
				  "price": "9006.71",
				  "period": "1"
				},
				{
				  "price": "18013.4",
				  "period": "2"
				},
				{
				  "price": "27020.11",
				  "period": "3"
				},
				{
				  "price": "36026.8",
				  "period": "4"
				},
				{
				  "price": "45033.51",
				  "period": "5"
				}
			  ]
			}
		  },
		  "disclaimer": "This pricing list is for informational purposes only.The actual price completely depends on ecs-buy.aliyun.com",
		  "type": "Instance",
		  "site": "Intl",
		  "description": "The pricingInfo key structure is 'RegionId::InstanceType::NetworkType::OSType::IoOptimized'"
		}
	`)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected 'GET' request, got ‘%s’", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(tcJSON)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	resp, err := QueryPriceList(ts.URL)
	if err != nil {
		t.Errorf("Error querying price list: %s\n", err.Error())
	}

	assert.ElementsMatch(t, tc, resp)
}

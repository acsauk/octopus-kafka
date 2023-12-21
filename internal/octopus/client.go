package octopus

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
)

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	apiKey, baseURL string
	doer            Doer
}

type MeterPoint struct {
	GSP          string `json:"gsp"`
	MPAN         string `json:"mpan"`
	ProfileClass int    `json:"profile_class"`
}

func New(apiKey, baseURL string, doer Doer) *Client {
	return &Client{
		apiKey:  apiKey + ":",
		baseURL: baseURL,
		doer:    doer,
	}
}

func (c *Client) ElectricityMeterPoints(mpan string) MeterPoint {
	req, _ := http.NewRequest(http.MethodGet,
		c.baseURL+"/v1/electricity-meter-points/"+mpan+"/",
		nil,
	)
	req.Header.Set("Authorization", "Basic: "+base64.StdEncoding.EncodeToString([]byte(c.apiKey)))

	resp, _ := c.doer.Do(req)

	var meterPoint MeterPoint
	_ = json.NewDecoder(resp.Body).Decode(&meterPoint)

	return meterPoint
}

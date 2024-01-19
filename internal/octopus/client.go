package octopus

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type (
	Doer interface {
		Do(req *http.Request) (*http.Response, error)
	}

	Client struct {
		apiKey, baseURL string
		doer            Doer
	}

	MeterPoint struct {
		GSP          string `json:"gsp"`
		MPAN         string `json:"mpan"`
		ProfileClass int    `json:"profile_class"`
	}

	Register struct {
		Identifier           string `json:"identifier"`
		Rate                 string `json:"rate"`
		IsSettlementRegister bool   `json:"is_settlement_register"`
	}

	Meter struct {
		SerialNumber string     `json:"serial_number"`
		Registers    []Register `json:"registers"`
	}

	Agreement struct {
		TariffCode string     `json:"tariff_code"`
		ValidFrom  CustomTime `json:"valid_from,omitempty"`
		ValidTo    CustomTime `json:"valid_to,omitempty"`
	}

	ElectricityMeterPoint struct {
		MPAN                string      `json:"mpan"`
		ProfileClass        int         `json:"profile_class"`
		ConsumptionStandard int         `json:"consumption_standard"`
		Meters              []Meter     `json:"meters"`
		Agreements          []Agreement `json:"agreements"`
	}

	Property struct {
		Id                     json.Number             `json:"id"`
		MovedInAt              CustomTime              `json:"moved_in_at,omitempty"`
		MovedOutAt             CustomTime              `json:"moved_out_at"`
		AddressLine1           string                  `json:"address_line_1"`
		AddressLine2           string                  `json:"address_line_2"`
		AddressLine3           string                  `json:"address_line_3"`
		Town                   string                  `json:"town"`
		County                 string                  `json:"county"`
		Postcode               string                  `json:"postcode"`
		ElectricityMeterPoints []ElectricityMeterPoint `json:"electricity_meter_points"`
		GasMeterPoints         []interface{}           `json:"gas_meter_points"`
	}

	Account struct {
		Number     string     `json:"number"`
		Properties []Property `json:"properties"`
	}

	CustomTime struct {
		Time time.Time
	}
)

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	dateOnly := strings.Split(strings.Trim(string(b), "\""), "T")[0]
	if dateOnly == "" || dateOnly == "null" {
		return nil
	}

	t, err := time.Parse(time.DateOnly, dateOnly)
	if err != nil {
		return err
	}

	*ct = CustomTime{Time: t}
	return nil
}

func New(apiKey, baseURL string, doer Doer) *Client {
	return &Client{
		apiKey:  apiKey + ":",
		baseURL: baseURL,
		doer:    doer,
	}
}

func (c *Client) ElectricityMeterPoints(mpan string) (MeterPoint, error) {
	req, err := http.NewRequest(http.MethodGet,
		c.baseURL+"/electricity-meter-points/"+mpan+"/",
		nil,
	)
	if err != nil {
		return MeterPoint{}, err
	}

	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(c.apiKey)))

	resp, err := c.doer.Do(req)
	if err != nil {
		return MeterPoint{}, err
	}

	var meterPoint MeterPoint
	if err = json.NewDecoder(resp.Body).Decode(&meterPoint); err != nil {
		return MeterPoint{}, err
	}

	return meterPoint, nil
}

func (c *Client) Account(accountNumber string) (Account, error) {
	req, err := http.NewRequest(http.MethodGet,
		c.baseURL+"/accounts/"+accountNumber+"/",
		nil,
	)
	if err != nil {
		return Account{}, err
	}

	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(c.apiKey)))

	resp, err := c.doer.Do(req)
	if err != nil {
		return Account{}, err
	}

	var account Account
	if err = json.NewDecoder(resp.Body).Decode(&account); err != nil {
		return Account{}, err
	}

	return account, nil
}

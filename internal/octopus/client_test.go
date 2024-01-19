package octopus

import (
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	httpClient := http.DefaultClient

	assert.Equal(t, &Client{
		baseURL: "/base-url",
		apiKey:  "api-key:",
		doer:    httpClient,
	}, New(
		"api-key",
		"/base-url",
		httpClient),
	)
}

func TestClient_ElectricityMeterPoint(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet,
		"/base-url/v1/electricity-meter-points/mpan-number/",
		nil,
	)

	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("api-key:")))

	doer := NewMockDoer(t)
	doer.
		On("Do", req).
		Return(&http.Response{Body: io.NopCloser(strings.NewReader(`
{
    "gsp": "gsp",
    "mpan": "mpan-number",
    "profile_class": 1
}
`))}, nil)

	client := New("api-key", "/base-url", doer)

	meterPoint, err := client.ElectricityMeterPoints("mpan-number")

	assert.Nil(t, err)
	assert.Equal(t, MeterPoint{
		GSP:          "gsp",
		MPAN:         "mpan-number",
		ProfileClass: 1,
	}, meterPoint)
}

func TestClient_ElectricityMeterPointWhenRequestError(t *testing.T) {
	client := New("api-key", string([]byte{0x7f}), nil)

	_, err := client.ElectricityMeterPoints("mpan-number")

	assert.NotNil(t, err)
}

func TestClient_ElectricityMeterPointWhenDoerError(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet,
		"/base-url/v1/electricity-meter-points/mpan-number/",
		nil,
	)

	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("api-key:")))

	doer := NewMockDoer(t)
	doer.
		On("Do", req).
		Return(&http.Response{Body: io.NopCloser(strings.NewReader("body"))}, errors.New("err"))

	client := New("api-key", "/base-url", doer)

	_, err := client.ElectricityMeterPoints("mpan-number")

	assert.Equal(t, errors.New("err"), err)
}

func TestClient_ElectricityMeterPointWhenJSONError(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet,
		"/base-url/v1/electricity-meter-points/mpan-number/",
		nil,
	)

	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("api-key:")))

	doer := NewMockDoer(t)
	doer.
		On("Do", req).
		Return(&http.Response{Body: io.NopCloser(strings.NewReader("not JSON"))}, nil)

	client := New("api-key", "/base-url", doer)

	_, err := client.ElectricityMeterPoints("mpan-number")

	assert.NotNil(t, err)
}

func TestClient_Account(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet,
		"/base-url/accounts/account-number/",
		nil,
	)

	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("api-key:")))

	doer := NewMockDoer(t)
	doer.
		On("Do", req).
		Return(&http.Response{Body: io.NopCloser(strings.NewReader(`
{
  "number": "account-number",
  "properties": [
    {
      "id": 123455,
      "moved_in_at": "2000-01-01T00:00:00+01:00",
      "moved_out_at": null,
      "address_line_1": "Line 1",
      "address_line_2": "Line 2",
      "address_line_3": "Line 3",
      "town": "Town",
      "county": "County",
      "postcode": "Postcode",
      "electricity_meter_points": [
        {
          "mpan": "1200000000000",
          "profile_class": 0,
          "consumption_standard": 1428,
          "meters": [
            {
              "serial_number": "1111111111",
              "registers": [
                {
                  "identifier": "01",
                  "rate": "STANDARD",
                  "is_settlement_register": true
                }
              ]
            }
          ],
          "agreements": [
            {
              "tariff_code": "TARIFF-CODE-1",
              "valid_from": "2000-01-02T00:00:00Z",
              "valid_to": "2000-01-02T00:00:00Z"
            },
            {
              "tariff_code": "TARIFF-CODE-2",
              "valid_from": "2000-01-02T00:00:00Z",
              "valid_to": "2000-01-02T00:00:00Z"
            }
          ]
        }
      ],
      "gas_meter_points": []
    }
  ]
}
`))}, nil)

	client := New("api-key", "/base-url", doer)

	expectedMoveInDate, _ := time.Parse(time.DateOnly, "2000-01-01")
	expectedAgreementDate, _ := time.Parse(time.DateOnly, "2000-01-02")

	account, err := client.Account("account-number")

	assert.Nil(t, err)
	assert.Equal(t, "account-number", account.Number)
	assert.True(t, account.Properties[0].MovedInAt.Time.Equal(expectedMoveInDate))
	assert.True(t, account.Properties[0].ElectricityMeterPoints[0].Agreements[0].ValidFrom.Time.Equal(expectedAgreementDate))
	assert.True(t, account.Properties[0].ElectricityMeterPoints[0].Agreements[0].ValidTo.Time.Equal(expectedAgreementDate))
	assert.Equal(t, "TARIFF-CODE-1", account.Properties[0].ElectricityMeterPoints[0].Agreements[0].TariffCode)
	assert.True(t, account.Properties[0].ElectricityMeterPoints[0].Agreements[1].ValidFrom.Time.Equal(expectedAgreementDate))
	assert.True(t, account.Properties[0].ElectricityMeterPoints[0].Agreements[1].ValidTo.Time.Equal(expectedAgreementDate))
	assert.Equal(t, "TARIFF-CODE-2", account.Properties[0].ElectricityMeterPoints[0].Agreements[1].TariffCode)
}

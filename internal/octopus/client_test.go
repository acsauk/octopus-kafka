package octopus

import (
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

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

	req.Header.Set("Authorization", "Basic: "+base64.StdEncoding.EncodeToString([]byte("api-key:")))

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

	req.Header.Set("Authorization", "Basic: "+base64.StdEncoding.EncodeToString([]byte("api-key:")))

	doer := NewMockDoer(t)
	doer.
		On("Do", req).
		Return(&http.Response{Body: io.NopCloser(strings.NewReader("body"))}, errors.New("err"))

	client := New("api-key", "/base-url", doer)

	_, err := client.ElectricityMeterPoints("mpan-number")

	assert.Equal(t, errors.New("err"), err)
}

func TestClient_ElectricityMeterPointWhenJSONerror(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet,
		"/base-url/v1/electricity-meter-points/mpan-number/",
		nil,
	)

	req.Header.Set("Authorization", "Basic: "+base64.StdEncoding.EncodeToString([]byte("api-key:")))

	doer := NewMockDoer(t)
	doer.
		On("Do", req).
		Return(&http.Response{Body: io.NopCloser(strings.NewReader("not JSON"))}, nil)

	client := New("api-key", "/base-url", doer)

	_, err := client.ElectricityMeterPoints("mpan-number")

	assert.NotNil(t, err)
}

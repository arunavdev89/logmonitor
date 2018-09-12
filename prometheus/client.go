package prometheus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type PrometheusClient struct {
	Server *url.URL
}

func NewPrometheusClient(address string) (*PrometheusClient, error) {
	u, err := url.Parse(address)
	if err != nil {
		return nil, err
	}
	return &PrometheusClient{
		Server: u,
	}, nil
}

//prometheus query response
type QueryResponse struct {
	Status string             `json:"status"`
	Data   *QueryResponseData `json:"data"`
}

//prometheus query response data
type QueryResponseData struct {
	Result []*QueryResponseResult `json:"result"`
}

//prometheus query response result entry
type QueryResponseResult struct {
	Metric map[string]string     `json:"metric"`
	Values []*QueryResponseValue `json:"value"`
}

type QueryResponseValue interface{}

func (r *QueryResponseResult) Time() time.Time {
	v := r.Values[0]
	t := (*v).(float64)
	return time.Unix(int64(t), 0)
}

func (r *QueryResponseResult) Value() string {
	v := r.Values[1]
	return (*v).(string)

}

func (c *PrometheusClient) Query(query string, instant time.Time) (*QueryResponse, error) {
	u, err := url.Parse(fmt.Sprintf("./api/v1/query?query=%s&time=%s",
		url.QueryEscape(query),
		url.QueryEscape(fmt.Sprintf("%d", instant.Unix())),
	),
	)
	if err != nil {
		return nil, err
	}

	u = c.Server.ResolveReference(u)
	r, err := http.Get(u.String())

	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if 400 <= r.StatusCode {
		return nil, fmt.Errorf("error response: %s", string(b))
	}
	resp := &QueryResponse{}
	err = json.Unmarshal(b, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

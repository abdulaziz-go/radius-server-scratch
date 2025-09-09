package metrics

import (
	"fmt"
	cryptoUtil "radius-server/src/utils/crypto"
	"sync"
)

type RadiusRequestTypes string

var (
	AccessRequest   RadiusRequestTypes = "Access-Request"
	AccountingStart RadiusRequestTypes = "Accounting-Start"
	AccountingStop  RadiusRequestTypes = "Accounting-Stop"
	InterimUpdate   RadiusRequestTypes = "Interim-Update"
	CoA             RadiusRequestTypes = "CoA"
	Disconnect      RadiusRequestTypes = "Disconnect"
)

type RadiusResponseStatus string

var (
	Success RadiusResponseStatus = "Success"
	Failure RadiusResponseStatus = "Failure"
)

type Metric struct {
	RequestType RadiusRequestTypes
	Status      RadiusResponseStatus
	Value       float64
}

type ResponseTimes struct {
	Total float64
	Count int64
}

var ResponseTime = ResponseTimes{}

type AvgResponseTime struct {
	TotalCount        int64
	TotalResponseTime float64
}

var AvgResponseTimes = map[string]AvgResponseTime{}

var (
	AvarageResponseTime  map[string]Metric = map[string]Metric{}
	MaxResponseTime      map[string]Metric = map[string]Metric{}
	MinResponseTime      map[string]Metric = map[string]Metric{}
	TotalCountOfRequests map[string]Metric = map[string]Metric{}
)
var mu sync.RWMutex

func CreateRequestMetric(requestType RadiusRequestTypes, status RadiusResponseStatus, responseTime float64) error {
	key := cryptoUtil.HashString(string(requestType), string(status))

	mu.Lock()
	defer mu.Unlock()

	// Create Avarage response time
	if v, ok := AvgResponseTimes[key]; ok {
		v.TotalCount++
		v.TotalResponseTime += responseTime
		AvgResponseTimes[key] = v
	} else {
		AvgResponseTimes[key] = AvgResponseTime{
			TotalCount:        1,
			TotalResponseTime: responseTime,
		}
	}

	latestAvgResponsetime := AvgResponseTimes[key].TotalResponseTime / float64(AvgResponseTimes[key].TotalCount)
	if v, ok := AvarageResponseTime[key]; ok {
		v.Value = latestAvgResponsetime
		AvarageResponseTime[key] = v
	} else {
		AvarageResponseTime[key] = Metric{
			RequestType: requestType,
			Status:      status,
			Value:       latestAvgResponsetime,
		}
	}

	// Create Max response time
	if v, ok := MaxResponseTime[key]; ok {
		if v.Value < responseTime {
			v.Value = responseTime
			MaxResponseTime[key] = v
		}
	} else {
		MaxResponseTime[key] = Metric{
			RequestType: requestType,
			Status:      status,
			Value:       responseTime,
		}
	}

	// Create Min response time
	if v, ok := MinResponseTime[key]; ok {
		if v.Value > responseTime {
			v.Value = responseTime
			MinResponseTime[key] = v
		}
	} else {
		MinResponseTime[key] = Metric{
			RequestType: requestType,
			Status:      status,
			Value:       responseTime,
		}
	}

	// Create total requests
	if v, ok := TotalCountOfRequests[key]; ok {
		v.Value += 1
	} else {
		TotalCountOfRequests[key] = Metric{
			RequestType: requestType,
			Status:      status,
			Value:       1,
		}
	}

	return nil
}

func GetMetricsPromtheusFormatted() []string {
	response := []string{}
	mu.RLock()
	defer mu.RUnlock()
	// avg
	response = append(response, "# HELP radius_response_time_seconds_avg is avarage response time of each type of requests\n")
	response = append(response, "# TYPE radius_response_time_seconds_avg gauge\n")
	for _, v := range AvarageResponseTime {
		txt := fmt.Sprintf("radius_response_time_seconds_avg{request-type=\"%s\",status=\"%s\"} %.3f\n", v.RequestType, v.Status, v.Value)
		response = append(response, txt)
	}

	// max
	response = append(response, "# HELP radius_response_time_seconds_max is table to save latest max responsetime\n")
	response = append(response, "# TYPE radius_response_time_seconds_max gauge\n")
	for _, v := range MaxResponseTime {
		txt := fmt.Sprintf("radius_response_time_seconds_max{request-type=\"%s\",status=\"%s\"} %.3f\n", v.RequestType, v.Status, v.Value)
		response = append(response, txt)
	}

	// min
	response = append(response, "# HELP radius_response_time_seconds_min is table to save latest min responsetime\n")
	response = append(response, "# TYPE radius_response_time_seconds_min gauge\n")
	for _, v := range MinResponseTime {
		txt := fmt.Sprintf("radius_response_time_seconds_min{request-type=\"%s\",status=\"%s\"} %.3f\n", v.RequestType, v.Status, v.Value)
		response = append(response, txt)
	}

	// request count
	response = append(response, "# HELP radius_requests_total is table to saver total number of requests\n")
	response = append(response, "# TYPE radius_requests_total counter\n")
	for _, v := range MinResponseTime {
		txt := fmt.Sprintf("radius_requests_total{request-type=\"%s\",status=\"%s\"} %d\n", v.RequestType, v.Status, int64(v.Value))
		response = append(response, txt)
	}

	return response
}

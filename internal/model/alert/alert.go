package alert

import (
	"am-kafka-project/pkg/common"
	"time"
)

//	{
//		"version": "4",
//		"groupKey": <string>,              // key identifying the group of alerts (e.g. to deduplicate)
//		"truncatedAlerts": <int>,          // how many alerts have been truncated due to "max_alerts"
//		"status": "<resolved|firing>",
//		"receiver": <string>,
//		"groupLabels": <object>,
//		"commonLabels": <object>,
//		"commonAnnotations": <object>,
//		"externalURL": <string>,           // backlink to the Alertmanager.
//		"alerts": [
//		  {
//			"status": "<resolved|firing>",
//			"labels": <object>,
//			"annotations": <object>,
//			"startsAt": "<rfc3339>",
//			"endsAt": "<rfc3339>",
//			"generatorURL": <string>,      // identifies the entity that caused the alert
//			"fingerprint": <string>        // fingerprint to identify the alert
//		  },
//		  ...
//		]
//	}
type Alert struct {
	Status       string            `json:"status"`
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
	StartsAt     time.Time         `json:"startsAt"`
	EndsAt       time.Time         `json:"endsAt"`
	GeneratorURL string            `json:"generatorURL"`
	Fingerprint  string            `json:"fingerprint"`
}

type Alerts struct {
	Version           string            `json:"version"`
	GroupKey          string            `json:"groupKey"`
	TruncatedAlerts   int               `json:"truncatedAlerts"`
	Status            string            `json:"status"`
	Receiver          string            `json:"receiver"`
	GroupLabels       map[string]string `json:"groupLabels"`
	CommonLabels      map[string]string `json:"commonLabels"`
	CommonAnnotations map[string]string `json:"commonAnnotations"`
	ExternalURL       string            `json:"externalURL"`
	Alerts            []Alert           `json:"alerts"`
}

type KafkaAlert struct {
	Name         string    `json:"alertname"`
	Status       string    `json:"status"`
	Labels       string    `json:"labels"`
	Annotations  string    `json:"annotations"`
	StartsAt     time.Time `json:"starts_at"`
	EndsAt       time.Time `json:"ends_at"`
	GeneratorURL string    `json:"generator_url"`
	Fingerprint  string    `json:"fingerprint"`
}

func (a *Alert) ToKafkaAlert() *KafkaAlert {
	return &KafkaAlert{
		Name:         a.Labels["alertname"],
		Status:       a.Status,
		Labels:       common.MapStringJson(a.Labels),
		Annotations:  common.MapStringJson(a.Annotations),
		StartsAt:     a.StartsAt.Round(time.Second),
		EndsAt:       a.EndsAt.Round(time.Second),
		GeneratorURL: a.GeneratorURL,
		Fingerprint:  a.Fingerprint,
	}
}

package model

import (
	"fmt"
	"time"

	pm "github.com/prometheus/common/model"
)

// Alert is a generic representation of an alert in the Prometheus eco-system.
type Alert struct {
	// Label value pairs for purpose of aggregation, matching, and disposition
	// dispatching. This must minimally include an "alertname" label.
	Labels pm.LabelSet `json:"labels"`

	// Extra key/value information which does not define alert identity.
	Annotations pm.LabelSet `json:"annotations,omitempty"`

	// The known time range for this alert. Both ends are optional.
	StartsAt     *time.Time `json:"startsAt,omitempty"`
	EndsAt       *time.Time `json:"endsAt,omitempty"`
	GeneratorURL string     `json:"generatorURL,omitempty"`
}

func (a Alert) Validate() error {
	if err := a.Labels.Validate(); err != nil {
		return err
	}

	if _, ok := a.Labels["alertname"]; !ok {
		return fmt.Errorf("'alertname' label must be set")
	}

	if err := a.Annotations.Validate(); err != nil {
		return err
	}

	return nil

	for ln, lv := range a.Labels {
		if !ln.IsValid() {
			return fmt.Errorf("invalid name %q", ln)
		}
		if !lv.IsValid() {
			return fmt.Errorf("invalid value %q", lv)
		}
	}
	return nil
}

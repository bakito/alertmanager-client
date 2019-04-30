package cmd

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/bakito/alertmanager-client/pkg/model"
	pm "github.com/prometheus/common/model"
	"github.com/spf13/cobra"
	"gopkg.in/resty.v1"
)

var (
	startsAt    string
	endsAt      string
	targetURL   string
	token       string
	labels      []string
	annotations []string
)

// alertCmd represents the alert command
var alertCmd = &cobra.Command{
	Use:   "alert",
	Short: "Send alerts to alertmanager",
	RunE: func(cmd *cobra.Command, args []string) error {
		a := new(model.Alert)

		if ls, err := toLabelSet(labels); err == nil {
			a.Labels = ls
		} else {
			return err
		}

		if ls, err := toLabelSet(annotations); err == nil {
			a.Annotations = ls
		} else {
			return err
		}

		if s, err := toTime(startsAt); err == nil {
			a.StartsAt = s
		} else {
			return err
		}

		if s, err := toTime(endsAt); err == nil {
			a.EndsAt = s
		} else {
			return err
		}

		if err := a.Validate(); err != nil {
			return err
		}

		if targetURL != "" {
			r := resty.R().
				SetHeader("Content-Type", "application/json").
				SetBody(a)
			if token != "" {
				r.SetAuthToken(token)
			}
			resp, err := r.Post(targetURL)
			if err != nil {
				return err
			}
			fmt.Print(resp.StatusCode())
		} else {
			b, err := json.Marshal(a)
			if err != nil {
				return err
			}
			fmt.Printf("alert called %s\n", b)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(alertCmd)

	alertCmd.PersistentFlags().StringSliceVarP(&labels, "label", "l", make([]string, 0), "labels; url query formatted values. 'alertname' one or the labels")
	alertCmd.PersistentFlags().StringSliceVarP(&annotations, "annotation", "a", make([]string, 0), "annotation; url query formatted values")
	alertCmd.PersistentFlags().StringVarP(&startsAt, "startsAt", "s", "", "starts at")
	alertCmd.PersistentFlags().StringVarP(&endsAt, "endsAt", "e", "", "ends at")

	alertCmd.PersistentFlags().StringVarP(&targetURL, "targetURL", "t", "", "the target url to sent the alert to")
	alertCmd.PersistentFlags().StringVar(&token, "token", "", "the target auth token")

	cobra.MarkFlagRequired(alertCmd.PersistentFlags(), "label")
}

func toLabelSet(values []string) (pm.LabelSet, error) {
	ls := pm.LabelSet{}
	for _, l := range values {
		values, err := url.ParseQuery(l)
		if err != nil {
			return ls, err
		}
		for k, v := range values {
			ls[pm.LabelName(k)] = pm.LabelValue(v[0])
		}
	}

	return ls, nil
}

func toTime(str string) (*time.Time, error) {
	if str == "" {
		return nil, nil
	}
	layout := "2006-01-02T15:04:05.000Z"
	t, err := time.Parse(layout, str)

	if err != nil {
		return nil, err
	}
	return &t, nil
}

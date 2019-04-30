// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	labels      []string
	annotations []string
)

// alertCmd represents the alert command
var alertCmd = &cobra.Command{
	Use:   "alert",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

		b, err := json.Marshal(a)
		if err != nil {
			return err
		}

		fmt.Printf("alert called %s\n", b)

		if targetURL != "" {
			resp, err := resty.R().
				SetHeader("Content-Type", "application/json").
				SetBody(a).
				Post(targetURL)
			if err != nil {
				return err
			}
			fmt.Print(resp.StatusCode())
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(alertCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	alertCmd.PersistentFlags().StringSliceVarP(&labels, "label", "l", make([]string, 0), "labels; url query formatted values. 'alertname' one or the labels")
	alertCmd.PersistentFlags().StringSliceVarP(&annotations, "annotation", "a", make([]string, 0), "annotation; url query formatted values")
	alertCmd.PersistentFlags().StringVarP(&startsAt, "startsAt", "s", "", "starts at")
	alertCmd.PersistentFlags().StringVarP(&endsAt, "endsAt", "e", "", "ends at")

	alertCmd.PersistentFlags().StringVarP(&targetURL, "targetURL", "t", "", "the target url to sent the alert to")

	cobra.MarkFlagRequired(alertCmd.PersistentFlags(), "label")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// alertCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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

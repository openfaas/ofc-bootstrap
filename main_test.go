package main

import (
	"errors"
	"strings"
	"testing"

	"github.com/openfaas-incubator/ofc-bootstrap/pkg/types"
)

func Test_filterDNSFeature(t *testing.T) {
	tests := []struct {
		title           string
		plan            types.Plan
		expectedFeature string
		expectedErr     error
	}{
		{
			title:           "DNS Service provider is Google",
			plan:            types.Plan{TLSConfig: types.TLSConfig{DNSService: types.CloudDNS}},
			expectedFeature: types.GCPDNS,
			expectedErr:     nil,
		},
		{
			title:           "DNS Service provider is Amazon",
			plan:            types.Plan{TLSConfig: types.TLSConfig{DNSService: types.Route53}},
			expectedFeature: types.Route53DNS,
			expectedErr:     nil,
		},
		{
			title:           "DNS Service provider is Digital Ocean",
			plan:            types.Plan{TLSConfig: types.TLSConfig{DNSService: types.DigitalOcean}},
			expectedFeature: types.DODNS,
			expectedErr:     nil,
		},
		{
			title:           "DNS Service provider is Digital Ocean",
			plan:            types.Plan{TLSConfig: types.TLSConfig{DNSService: "unsupporteddns"}},
			expectedFeature: "",
			expectedErr:     errors.New("Error unavailable DNS service provider"),
		},
	}
	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			var planError error
			test.plan, planError = filterDNSFeature(test.plan)
			if planError != nil {
				if !strings.Contains(planError.Error(), test.expectedErr.Error()) {
					t.Errorf("Unexpected error message: %s", planError.Error())
				}
			}
			for _, feature := range test.plan.Features {
				if feature != test.expectedFeature {
					t.Errorf("Unexpected feature: %s", feature)
				}
			}
		})
	}
}

func Test_filterFeatures(t *testing.T) {
	tests := []struct {
		title            string
		planConfig       types.Plan
		expectedFeatures []string
		expectedError    error
	}{
		{
			title:            "Plan is empty only default feature is present",
			planConfig:       types.Plan{},
			expectedFeatures: []string{types.DefaultFeature},
			expectedError:    nil,
		},
		{
			title: "Every field which defines populated feature is populated",
			planConfig: types.Plan{
				Github: types.Github{
					AppID:          "example ID",
					PrivateKeyFile: "~/home/private_key.pem",
				},
				TLS: true,
				TLSConfig: types.TLSConfig{
					DNSService: types.Route53,
				},
				EnableOAuth: true,
			},
			expectedFeatures: []string{types.DefaultFeature, types.GitHub, types.Auth, types.Route53DNS},
			expectedError:    nil,
		},
		{
			title: "Auth and TLS are enabled",
			planConfig: types.Plan{
				TLS: true,
				TLSConfig: types.TLSConfig{
					DNSService: types.Route53,
				},
				EnableOAuth: true,
			},
			expectedFeatures: []string{types.DefaultFeature, types.Auth, types.Route53DNS},
			expectedError:    nil,
		},
	}
	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			test.planConfig, _ = filterFeatures(test.planConfig)
			for _, feature := range test.planConfig.Features {
				t.Logf("Searching for feature: %s, ", feature)
				for allFeatures, expectedFeature := range test.expectedFeatures {
					if feature == expectedFeature {
						break
					}
					if allFeatures == len(test.expectedFeatures) {
						t.Errorf("Feature: %s not found in the expected features", feature)
					}
				}
			}
		})
	}
}

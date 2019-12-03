package cmd

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
			title:           "DNS Service provider is Cloudflare",
			plan:            types.Plan{TLSConfig: types.TLSConfig{DNSService: "cloudflare"}},
			expectedFeature: types.CloudflareDNS,
			expectedErr:     nil,
		},
		{
			title:           "DNS Service provider is not supported",
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
				wantErr := ""
				if test.expectedErr != nil {
					wantErr = test.expectedErr.Error()
				}

				if strings.Contains(planError.Error(), wantErr) == false || len(wantErr) == 0 {
					t.Errorf("Got plan error: %s", planError.Error())
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
			title: "Every field which defines populated feature is populated for github",
			planConfig: types.Plan{
				SCM: types.GitHubSCM,
				TLS: true,
				TLSConfig: types.TLSConfig{
					DNSService: types.Route53,
				},
				EnableOAuth: true,
			},
			expectedFeatures: []string{types.DefaultFeature, types.GitHubFeature, types.Auth, types.Route53DNS},
			expectedError:    nil,
		},
		{
			title: "AWS ECR is detected as a feature",
			planConfig: types.Plan{
				EnableECR: true,
			},
			expectedFeatures: []string{types.DefaultFeature, types.ECRFeature},
			expectedError:    nil,
		},
		{
			title: "AWS ECR is disabled when not set as a feature",
			planConfig: types.Plan{
				EnableECR: false,
			},
			expectedFeatures: []string{types.DefaultFeature},
			expectedError:    nil,
		},
		{
			title: "Every field which defines populated feature is populated for gitlab",
			planConfig: types.Plan{
				SCM: types.GitLabSCM,
				TLS: true,
				TLSConfig: types.TLSConfig{
					DNSService: types.Route53,
				},
				EnableOAuth: true,
			},
			expectedFeatures: []string{types.DefaultFeature, types.GitLabFeature, types.Auth, types.Route53DNS},
			expectedError:    nil,
		},
		{
			title: "Example in which the function throws error in this case the SCM field is empty",
			planConfig: types.Plan{
				TLS: true,
				TLSConfig: types.TLSConfig{
					DNSService: types.Route53,
				},
				EnableOAuth: true,
			},
			expectedFeatures: []string{types.DefaultFeature},
			expectedError:    errors.New("Error while filtering features"),
		},
		{
			title: "Auth and TLS are enabled along with GitLab",
			planConfig: types.Plan{
				TLS: true,
				TLSConfig: types.TLSConfig{
					DNSService: types.Route53,
				},
				EnableOAuth: true,
				SCM:         types.GitHubSCM,
			},
			expectedFeatures: []string{types.DefaultFeature, types.Auth, types.Route53DNS},
			expectedError:    nil,
		},
	}

	for _, test := range tests {

		t.Run(test.title, func(t *testing.T) {

			var filterError error
			test.planConfig, filterError = filterFeatures(test.planConfig)
			t.Logf("Features in the plan: %v", test.planConfig.Features)

			if filterError != nil && test.expectedError != nil {

				if !strings.Contains(filterError.Error(), test.expectedError.Error()) {
					t.Errorf("Expected error to contain: `%s` got: `%s`", test.expectedError.Error(), filterError.Error())
				}

			}

			for _, expectedFeature := range test.expectedFeatures {
				for allPlanFeatures, enabledFeature := range test.planConfig.Features {
					if len(test.planConfig.Features) == 0 {
						t.Errorf("Feature 'default' should always be present")
					}
					if expectedFeature == enabledFeature {
						break
					}
					if allPlanFeatures == len(test.planConfig.Features)-1 {
						t.Errorf("Feature: '%s' not found in: %v", expectedFeature, test.planConfig.Features)
					}
				}
			}
		})
	}
}

func Test_filterGitRepositoryManager(t *testing.T) {
	tests := []struct {
		title           string
		planConfig      types.Plan
		expectedFeature []string
		expectedError   error
	}{
		{
			title: "SCM field is populated for gitlab",
			planConfig: types.Plan{
				SCM: types.GitLabSCM,
			},
			expectedFeature: []string{types.GitLabFeature},
			expectedError:   nil,
		},
		{
			title: "SCM field is populated for github",
			planConfig: types.Plan{
				SCM: types.GitHubSCM,
			},
			expectedFeature: []string{types.GitHubFeature},
			expectedError:   nil,
		},
		{
			title: "SCM field is populated for with unsupported Git repository manager",
			planConfig: types.Plan{
				SCM: "bitbucket",
			},
			expectedFeature: []string{},
			expectedError:   errors.New("Error unsupported Git repository manager: bitbucket"),
		},
	}
	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			var configErr error
			test.planConfig, configErr = filterGitRepositoryManager(test.planConfig)
			if configErr != nil && test.expectedError != nil {
				if !strings.EqualFold(configErr.Error(), test.expectedError.Error()) {
					t.Errorf("Expected error: '%s' got: '%s'", test.expectedError, configErr)
				}
			}
		})
	}
}

// Copyright (c) OpenFaaS Author(s) 2019. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package types

import (
	"testing"
)

func Test_mergePlans_Empty(t *testing.T) {

	_, err := MergePlans([]Plan{})

	if err == nil {
		t.Errorf("Expected an error for no plans")
		t.Fail()
	}
	want := "at least one plan is required"
	if err.Error() != want {
		t.Errorf("no plans error want: %s, got: %s", want, err.Error())
		t.Fail()
	}
}

func Test_mergePlans_OnlyOneItem(t *testing.T) {

	plan1 := Plan{
		OpenFaaSCloudVersion: "master",
	}

	planOut, err := MergePlans([]Plan{plan1})

	if err != nil {
		t.Errorf("Got error for a single plan, expected no error: %s", err.Error())
		t.Fail()
	}

	if plan1.OpenFaaSCloudVersion != planOut.OpenFaaSCloudVersion {
		t.Errorf("OpenFaaSCloudVersion want: %s, but got: %s", plan1.OpenFaaSCloudVersion, planOut.OpenFaaSCloudVersion)
	}
}

func Test_mergePlans_MergeEmptyItemsFromBoth(t *testing.T) {

	plan1 := Plan{
		OpenFaaSCloudVersion: "master",
	}

	plan2 := Plan{
		CustomersURL: "https://127.0.0.1:8443/customers",
	}

	planOut, err := MergePlans([]Plan{plan1, plan2})

	if err != nil {
		t.Errorf("Got error, expected no error: %s", err.Error())
		t.Fail()
	}

	if planOut.OpenFaaSCloudVersion != plan1.OpenFaaSCloudVersion {
		t.Errorf("OpenFaaSCloudVersion want: %s, but got: %s", plan1.OpenFaaSCloudVersion, planOut.OpenFaaSCloudVersion)
	}

	if planOut.CustomersURL != plan2.CustomersURL {
		t.Errorf("CustomersURL want: %s, but got: %s", plan2.CustomersURL, planOut.CustomersURL)
	}
}

func Test_mergePlans_PlanValuesOverwriteAccordingToOrder(t *testing.T) {

	plan1 := Plan{
		OpenFaaSCloudVersion: "0.12.0",
	}

	plan2 := Plan{
		OpenFaaSCloudVersion: "0.11.0",
	}

	planOut, err := MergePlans([]Plan{plan1, plan2})

	if err != nil {
		t.Errorf("Got error, expected no error: %s", err.Error())
		t.Fail()
	}

	wantVer := plan2.OpenFaaSCloudVersion
	if planOut.OpenFaaSCloudVersion != wantVer {
		t.Errorf("OpenFaaSCloudVersion want: %s, but got: %s", wantVer, planOut.OpenFaaSCloudVersion)
	}

}

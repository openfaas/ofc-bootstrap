// Copyright (c) OpenFaaS Author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package types

import (
	"fmt"

	"github.com/imdario/mergo"
	"github.com/jinzhu/copier"
)

// MergePlans combines one or more plan with a manual merge for
// the list of secrets.
func MergePlans(plans []Plan) (*Plan, error) {
	var err error
	masterPlan := &Plan{}

	if len(plans) == 1 {
		return &plans[0], nil
	}

	if len(plans) == 0 {
		return masterPlan, fmt.Errorf("at least one plan is required")
	}

	for _, plan := range plans {
		mergeErr := mergo.Merge(masterPlan, plan, mergo.WithOverride)
		if mergeErr != nil {
			return masterPlan, mergeErr
		}
	}

	patchSecrets(masterPlan, plans)

	return masterPlan, err
}

func patchSecrets(masterPlan *Plan, plans []Plan) {
	masterList := []KeyValueNamespaceTuple{}

	// Read each plan
	for _, plan := range plans {

		// Process each secret
		for _, v := range plan.Secrets {

			// Apply to master list
			index := -1
			for i, mv := range masterList {
				if mv.Name == v.Name {
					index = i
					break
				}
			}

			if index == -1 {
				item := KeyValueNamespaceTuple{}
				copier.Copy(&item, &v)
				masterList = append(masterList, item)
			} else {
				copier.Copy(&masterList[index], &v)
			}
		}
	}
	masterPlan.Secrets = masterList
}

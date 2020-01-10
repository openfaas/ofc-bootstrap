// Copyright (c) OpenFaaS Author(s) 2019. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package cmd

import (
	"fmt"

	"github.com/imdario/mergo"
	"github.com/openfaas-incubator/ofc-bootstrap/pkg/types"
)

func mergePlans(plans []types.Plan) (*types.Plan, error) {
	var err error
	masterPlan := &types.Plan{}

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

	return masterPlan, err
}

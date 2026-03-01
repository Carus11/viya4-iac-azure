// Copyright © 2025, SAS Institute Inc., Cary, NC, USA. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package nondefaultplan

import (
	"test/helpers"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlanACRStandard(t *testing.T) {
	t.Parallel()

	variables := helpers.GetDefaultPlanVars(t)
	variables["prefix"] = "acr-standard"
	variables["create_container_registry"] = true
	variables["container_registry_admin_enabled"] = true
	variables["container_registry_sku"] = "Standard"

	tests := map[string]helpers.TestCase{
		"acrGeoRepsNotExistTest": {
			Expected:          "[]",
			ResourceMapName:   "azurerm_container_registry.acr[0]",
			AttributeJsonPath: "{$.georeplications}",
			Message:           "Geo-replications found when they should not be present",
		},
		"nameTest": {
			Expected:          "acr",
			ResourceMapName:   "azurerm_container_registry.acr[0]",
			AttributeJsonPath: "{$.name}",
			AssertFunction:    assert.Contains,
			Message:           "ACR name does not contain 'acr'",
		},
		"skuTest": {
			Expected:          "Standard",
			ResourceMapName:   "azurerm_container_registry.acr[0]",
			AttributeJsonPath: "{$.sku}",
			Message:           "Unexpected ACR SKU value",
		},
		"adminEnabledTest": {
			Expected:          "true",
			ResourceMapName:   "azurerm_container_registry.acr[0]",
			AttributeJsonPath: "{$.admin_enabled}",
			Message:           "Unexpected ACR admin_enabled value",
		},
	}

	plan := helpers.GetPlanFromCache(t, variables)
	helpers.RunTests(t, tests, plan)
}

// Verify ACR premium
func TestPlanACRPremium(t *testing.T) {
	t.Parallel()

	variables := helpers.GetDefaultPlanVars(t)
	variables["prefix"] = "acr-premium"
	variables["create_container_registry"] = true
	variables["container_registry_admin_enabled"] = true
	variables["container_registry_sku"] = "Premium"
	variables["container_registry_geo_replica_locs"] = []string{"southeastus5", "southeastus3"}

	tests := map[string]helpers.TestCase{
		"locationsTest": {
			Expected:          "southeastus3 southeastus5",
			ResourceMapName:   "azurerm_container_registry.acr[0]",
			AttributeJsonPath: "{$.georeplications[*].location}",
			Message:           "Geo-replications do not match expected values",
		},
		"nameTest": {
			Expected:          "acr",
			ResourceMapName:   "azurerm_container_registry.acr[0]",
			AttributeJsonPath: "{$.name}",
			AssertFunction:    assert.Contains,
			Message:           "ACR name does not contain 'acr'",
		},
		"skuTest": {
			Expected:          "Premium",
			ResourceMapName:   "azurerm_container_registry.acr[0]",
			AttributeJsonPath: "{$.sku}",
			Message:           "Unexpected ACR SKU value",
		},
		"adminEnabledTest": {
			Expected:          "true",
			ResourceMapName:   "azurerm_container_registry.acr[0]",
			AttributeJsonPath: "{$.admin_enabled}",
			Message:           "Unexpected ACR admin_enabled value",
		},
	}

	plan := helpers.GetPlan(t, variables)
	helpers.RunTests(t, tests, plan)
}

// Verify ACR private endpoint configuration
func TestPlanACRPrivate(t *testing.T) {
	t.Parallel()

	variables := helpers.GetDefaultPlanVars(t)
	variables["prefix"] = "acr-private"
	variables["create_container_registry"] = true
	variables["container_registry_private"] = true

	tests := map[string]helpers.TestCase{
		"skuTest": {
			Expected:          "Premium",
			ResourceMapName:   "azurerm_container_registry.acr[0]",
			AttributeJsonPath: "{$.sku}",
			Message:           "Unexpected ACR SKU value for private registry",
		},
		"publicAccessTest": {
			Expected:          "false",
			ResourceMapName:   "azurerm_container_registry.acr[0]",
			AttributeJsonPath: "{$.public_network_access_enabled}",
			Message:           "Unexpected ACR public_network_access_enabled value",
		},
		"privateEndpointTest": {
			Expected:          "acr-private-acr-pe",
			ResourceMapName:   "azurerm_private_endpoint.acr[0]",
			AttributeJsonPath: "{$.name}",
			Message:           "Private endpoint for ACR was not created or has incorrect name",
		},
		"dnsZoneTest": {
			Expected:          "privatelink.azurecr.io",
			ResourceMapName:   "azurerm_private_dns_zone.acr[0]",
			AttributeJsonPath: "{$.name}",
			Message:           "Private DNS zone for ACR was not created or has incorrect name",
		},
		"roleAssignmentTest": {
			Expected:          "AcrPull",
			ResourceMapName:   "azurerm_role_assignment.acr_pull[0]",
			AttributeJsonPath: "{$.role_definition_name}",
			Message:           "AcrPull role assignment was not created",
		},
	}

	plan := helpers.GetPlanFromCache(t, variables)
	helpers.RunTests(t, tests, plan)
}

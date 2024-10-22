// Copyright 2021 The Terraformer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package okta

import (
	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	"github.com/okta/okta-sdk-golang/v5/okta"
)

type NetworkZoneGenerator struct {
	OktaService
}

func (g NetworkZoneGenerator) createResources(networkZoneList []okta.ListNetworkZones200ResponseInner) []terraformutils.Resource {
	var resources []terraformutils.Resource
	for _, networkZone := range networkZoneList {
		var zone []terraformutils.Resource
		if networkZone.DynamicNetworkZone != nil {

			zone = DynamicNetworkZone(networkZone)
			resources = append(resources, zone...)

		}
		if networkZone.IPNetworkZone != nil {
			zone = IPNetworkZone(networkZone)
			resources = append(resources, zone...)
		}
		//resources = append(resources, zone...)
	}

	return resources

}

func (g *NetworkZoneGenerator) InitResources() error {
	ctx, client, err := g.ClientV5()
	if err != nil {
		return err
	}

	output, resp, err := client.NetworkZoneAPI.ListNetworkZones(ctx).Execute()
	if err != nil {
		return err
	}

	for resp.HasNextPage() {
		var networkZoneSet []okta.ListNetworkZones200ResponseInner
		resp, _ = resp.Next(&networkZoneSet)
		output = append(output, networkZoneSet...)
	}

	g.Resources = g.createResources(output)
	return nil
}

func IPNetworkZone(zone okta.ListNetworkZones200ResponseInner) []terraformutils.Resource {
	var resource []terraformutils.Resource
	resource = append(resource, terraformutils.NewResource(
		*zone.IPNetworkZone.Id,
		zone.IPNetworkZone.Name,
		"okta_network_zone",
		"okta",
		map[string]string{
			"name": zone.IPNetworkZone.Name,
			"type": zone.IPNetworkZone.Type,
		},
		[]string{},
		attributesIPNetworkZone(zone),
	))
	return resource
}

func attributesIPNetworkZone(zone okta.ListNetworkZones200ResponseInner) map[string]interface{} {
	attributes := map[string]interface{}{}
	attributes["usage"] = zone.IPNetworkZone.Usage
	if zone.IPNetworkZone.Proxies != nil {
		attributes["proxies"] = zone.IPNetworkZone.Proxies
	}
	if zone.IPNetworkZone.Gateways != nil {
		attributes["gateways"] = zone.IPNetworkZone.Gateways
	}

	return attributes
}

func DynamicNetworkZone(zone okta.ListNetworkZones200ResponseInner) []terraformutils.Resource {
	var resource []terraformutils.Resource
	resource = append(resource, terraformutils.NewResource(
		*zone.DynamicNetworkZone.Id,
		zone.DynamicNetworkZone.Name,
		"okta_network_zone",
		"okta",
		map[string]string{
			"name": zone.DynamicNetworkZone.Name,
			"type": zone.DynamicNetworkZone.Type,
		},
		[]string{},
		attributesDynamicNetworkZone(zone),
	))
	return resource

}
func attributesDynamicNetworkZone(zone okta.ListNetworkZones200ResponseInner) map[string]interface{} {
	attributes := map[string]interface{}{}
	attributes["usage"] = *zone.DynamicNetworkZone.Usage
	attributes["proxytype"] = *zone.DynamicNetworkZone.ProxyType
	attributes["status"] = *zone.DynamicNetworkZone.Status
	if zone.DynamicNetworkZone.Locations != nil {
		dynamicLocations := []string{}
		for _, location := range zone.DynamicNetworkZone.Locations {
			locationStr := ""
			if location.Region != nil {
				locationStr = *location.Region
			} else {
				locationStr = *location.Country
			}

			// Append the formatted location string to dynamicLocations slice
			dynamicLocations = append(dynamicLocations, locationStr)
		}
		attributes["dynamic_locations"] = dynamicLocations
	}
	attributes["asns"] = zone.DynamicNetworkZone.GetAsns()

	return attributes
}

/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package output

import "reflect"

// Nat : mapping of a nat component
type Nat struct {
	NatAWSID               string   `json:"nat_gateway_aws_id"`
	Name                   string   `json:"name"`
	NatType                string   `json:"nat_type"`
	PublicNetwork          string   `json:"public_network"`
	RoutedNetworks         []string `json:"routed_networks"`
	RoutedNetworkAWSIDs    []string `json:"routed_networks_aws_ids"`
	PublicNetworkAWSID     string   `json:"public_network_aws_id"`
	NatGatewayAWSID        string   `json:"nat_gateway_aws_id"`
	NatGatewayAllocationID string   `json:"nat_gateway_allocation_id"`
	NatGatewayAllocationIP string   `json:"nat_gateway_allocation_ip"`
	DatacenterType         string   `json:"datacenter_type"`
	DatacenterName         string   `json:"datacenter_name"`
	DatacenterRegion       string   `json:"datacenter_region"`
	DatacenterToken        string   `json:"datacenter_token"`
	DatacenterSecret       string   `json:"datacenter_secret"`
	VpcID                  string   `json:"vpc_id"`
	Service                string   `json:"service"`
	Status                 string   `json:"status"`
	Exists                 bool
}

// HasChanged diff's the two items and returns true if there have been any changes
func (n *Nat) HasChanged(on *Nat) bool {
	if !reflect.DeepEqual(n.RoutedNetworks, on.RoutedNetworks) {
		return true
	}

	return false
}

func hasNetwork(networks []string, name string) bool {
	for _, network := range networks {
		if network == name {
			return true
		}
	}

	return false
}

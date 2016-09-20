/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package output

// Instance : mapping of an instance component
type VPC struct {
	VpcID     string `json:"vpc_id"`
	VpcSubnet string `json:"vpc_subnet"`
	Type      string `json:"_type"`
	Status    string `json:"status"`
	Exists    bool
}

// HasChanged diff's the two items and returns true if there have been any changes
func (v *VPC) HasChanged(ov *VPC) bool {
	if v.VpcSubnet != ov.VpcSubnet {
		return true
	}

	return false
}

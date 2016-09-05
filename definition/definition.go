/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"unicode/utf8"
)

// Definition ...
type Definition struct {
	Name           string          `json:"name"`
	Datacenter     string          `json:"datacenter"`
	ErnestIP       []string        `json:"ernest_ip"`
	ServiceIP      string          `json:"service_ip"`
	Networks       []Network       `json:"networks"`
	Instances      []Instance      `json:"instances"`
	SecurityGroups []SecurityGroup `json:"security_groups"`
	NatGateways    []NatGateway    `json:"nat_gateways"`
}

// New returns a new Definition
func New() *Definition {
	return &Definition{
		ErnestIP:       make([]string, 0),
		Networks:       make([]Network, 0),
		Instances:      make([]Instance, 0),
		SecurityGroups: make([]SecurityGroup, 0),
	}
}

// FromJSON creates a definition from json
func FromJSON(data []byte) (*Definition, error) {
	var d Definition

	err := json.Unmarshal(data, d)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

// ValidateName checks if service is valid
func (d *Definition) validateName() error {
	// Check if service name is null
	if d.Name == "" {
		return errors.New("Service name should not be null")
	}

	// Check if service name is > 50 characters
	if utf8.RuneCountInString(d.Name) > 50 {
		return fmt.Errorf("Datacenter name can't be greater than %d characters", AWSMAXNAME)
	}
	return nil
}

func (d *Definition) validateDatacenter() error {
	if d.Datacenter == "" {
		return errors.New("Datacenter not specified")
	}
	return nil
}

func (d *Definition) validateServiceIP() error {
	if d.ServiceIP == "" {
		return nil
	}
	if net.ParseIP(d.ServiceIP) == nil {
		return errors.New("ServiceIP is not a valid IP")
	}
	return nil
}

// Validate the definition
func (d *Definition) Validate() error {
	// Validate Definition
	err := d.validateName()
	if err != nil {
		return err
	}

	err = d.validateServiceIP()
	if err != nil {
		return err
	}

	// Validate Datacenter
	err = d.validateDatacenter()
	if err != nil {
		return err
	}

	// Validate Networks
	for _, n := range d.Networks {
		err := n.Validate()
		if err != nil {
			return err
		}
	}

	// Validate Instances
	for _, i := range d.Instances {
		nw := d.FindNetwork(i.Network)

		err := i.Validate(nw)
		if err != nil {
			return err
		}
	}

	// Validate Security Groups
	for _, sg := range d.SecurityGroups {
		err := sg.Validate(d.Networks)
		if err != nil {
			return err
		}
	}

	// Validate Nat Gateways
	for _, ng := range d.NatGateways {
		err := ng.Validate(d.Networks)
		if err != nil {
			return err
		}
	}

	if hasDuplicateNetworks(d.Networks) {
		return errors.New("Duplicate network names found")
	}

	if hasDuplicateInstance(d.Instances) {
		return errors.New("Duplicate instance names found")
	}

	return nil
}

// GeneratedName returns the generated service name
func (d *Definition) GeneratedName() string {
	return d.Datacenter + "-" + d.Name + "-"
}

// FindNetwork returns a network matched by name
func (d *Definition) FindNetwork(name string) *Network {
	for _, network := range d.Networks {
		if network.Name == name {
			return &network
		}
	}
	return nil
}

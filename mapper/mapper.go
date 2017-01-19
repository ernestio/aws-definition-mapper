/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"github.com/ernestio/aws-definition-mapper/definition"
	"github.com/ernestio/aws-definition-mapper/output"
)

// ConvertPayload will build an FSMMessage based on an input definition
func ConvertPayload(p *definition.Payload) *output.FSMMessage {
	m := output.FSMMessage{
		ID:          p.ServiceID,
		Service:     p.ServiceID,
		ServiceName: p.Service.Name,
		ClientName:  p.Client.Name,
		Type:        p.Datacenter.Type,
	}

	// Map datacenters
	m.Datacenters.Items = MapDatacenters(p.Datacenter)

	// Map VPCs
	m.VPCs.Items = MapVPCs(p)

	// Map networks
	m.Networks.Items = MapNetworks(p.Service)

	// Map instances
	m.Instances.Items = MapInstances(p.Service)

	// Map firewalls
	m.Firewalls.Items = MapSecurityGroups(p.Service)

	// Map nats/port forwarding
	m.Nats.Items = MapNats(p.Service)

	// Map ELB's
	m.ELBs.Items = MapELBs(p.Service)

	// Map S3 buckets
	m.S3s.Items = MapS3Buckets(p.Service)

	// Map Route53 zones
	m.Route53s.Items = MapRoute53Zones(p.Service)

	// Map RDS clusters
	m.RDSClusters.Items = MapRDSClusters(p.Service)

	// Map RDS instances
	m.RDSInstances.Items = MapRDSInstances(p.Service)

	// Map EBS volumes
	m.EBSVolumes.Items = MapEBSVolumes(p.Service)

	return &m
}

// ConvertFSMMessage : Convery an output FSMMessage to an input definition
func ConvertFSMMessage(m *output.FSMMessage) *definition.Payload {
	p := definition.Payload{
		ServiceID: m.Service,
	}

	p.Service.Name = m.ServiceName
	p.Client.Name = m.ClientName

	if len(m.VPCs.Items) > 0 {
		p.Service.VpcID = m.VPCs.Items[0].VpcID
		p.Service.VpcSubnet = m.VPCs.Items[0].VpcSubnet
	}

	p.Service.Networks = MapDefinitionNetworks(m)

	p.Service.Instances = MapDefinitionInstances(m)

	p.Service.SecurityGroups = MapDefinitionSecurityGroups(m)

	p.Service.NatGateways = MapDefinitionNats(m)

	p.Service.ELBs = MapDefinitionELBs(m)

	p.Service.EBSVolumes = MapDefinitionEBSVolumes(m)

	p.Service.S3Buckets = MapDefinitionS3Buckets(m)

	p.Service.Route53Zones = MapDefinitionRoute53Zones(m)

	p.Service.RDSInstances = MapDefinitionRDSInstances(m)

	p.Service.RDSClusters = MapDefinitionRDSClusters(m)

	p.Datacenter = MapDefinitionDatacenter(m.Datacenters.Items)

	return &p
}

// MapProviderData will map any information generated by a provider that is not
// deductible from the input definition
func MapProviderData(m, om *output.FSMMessage) {
	// Map network ID's
	for i, network := range m.Networks.Items {
		nw := om.FindNetwork(network.Name)
		if nw != nil {
			m.Networks.Items[i].NetworkAWSID = nw.NetworkAWSID
			m.Networks.Items[i].AvailabilityZone = nw.AvailabilityZone
			m.Networks.Items[i].DatacenterType = "$(datacenters.items.0.type)"
			m.Networks.Items[i].DatacenterName = "$(datacenters.items.0.name)"
			m.Networks.Items[i].AccessKeyID = "$(datacenters.items.0.aws_access_key_id)"
			m.Networks.Items[i].SecretAccessKey = "$(datacenters.items.0.aws_secret_access_key)"
			m.Networks.Items[i].DatacenterRegion = "$(datacenters.items.0.region)"
			m.Networks.Items[i].VpcID = "$(vpcs.items.0.vpc_id)"
		}
	}

	// Map instance ID's
	for i, instance := range m.Instances.Items {
		in := om.FindInstance(instance.Name)
		if in != nil {
			m.Instances.Items[i].InstanceAWSID = in.InstanceAWSID
			m.Instances.Items[i].PublicIP = in.PublicIP
			m.Instances.Items[i].ElasticIP = in.ElasticIP
			m.Instances.Items[i].ElasticIPAWSID = in.ElasticIPAWSID
			m.Instances.Items[i].DatacenterType = "$(datacenters.items.0.type)"
			m.Instances.Items[i].DatacenterName = "$(datacenters.items.0.name)"
			m.Instances.Items[i].AccessKeyID = "$(datacenters.items.0.aws_access_key_id)"
			m.Instances.Items[i].SecretAccessKey = "$(datacenters.items.0.aws_secret_access_key)"
			m.Instances.Items[i].DatacenterRegion = "$(datacenters.items.0.region)"
			m.Instances.Items[i].VpcID = "$(vpcs.items.0.vpc_id)"
		}
	}

	// Map firewall ID's
	for i, firewall := range m.Firewalls.Items {
		fw := om.FindFirewall(firewall.Name)
		if fw != nil {
			m.Firewalls.Items[i].SecurityGroupAWSID = fw.SecurityGroupAWSID
			m.Firewalls.Items[i].ProviderType = "$(datacenters.items.0.type)"
			m.Firewalls.Items[i].DatacenterType = "$(datacenters.items.0.type)"
			m.Firewalls.Items[i].DatacenterName = "$(datacenters.items.0.name)"
			m.Firewalls.Items[i].AccessKeyID = "$(datacenters.items.0.aws_access_key_id)"
			m.Firewalls.Items[i].SecretAccessKey = "$(datacenters.items.0.aws_secret_access_key)"
			m.Firewalls.Items[i].DatacenterRegion = "$(datacenters.items.0.region)"
			m.Firewalls.Items[i].VpcID = "$(vpcs.items.0.vpc_id)"
		}
	}

	// Map nat ID's
	for i, nat := range m.Nats.Items {
		nt := om.FindNat(nat.Name)
		if nt != nil {
			m.Nats.Items[i].NatGatewayAWSID = nt.NatGatewayAWSID
			m.Nats.Items[i].NatGatewayAllocationID = nt.NatGatewayAllocationID
			m.Nats.Items[i].NatGatewayAllocationIP = nt.NatGatewayAllocationIP
			m.Nats.Items[i].ProviderType = "$(datacenters.items.0.type)"
			m.Nats.Items[i].DatacenterType = "$(datacenters.items.0.type)"
			m.Nats.Items[i].DatacenterName = "$(datacenters.items.0.name)"
			m.Nats.Items[i].AccessKeyID = "$(datacenters.items.0.aws_access_key_id)"
			m.Nats.Items[i].SecretAccessKey = "$(datacenters.items.0.aws_secret_access_key)"
			m.Nats.Items[i].DatacenterRegion = "$(datacenters.items.0.region)"
			m.Nats.Items[i].VpcID = "$(vpcs.items.0.vpc_id)"
		}
	}

	// Map elb data
	for i, elb := range m.ELBs.Items {
		lb := om.FindELB(elb.Name)
		if lb != nil {
			m.ELBs.Items[i].DNSName = lb.DNSName
			m.ELBs.Items[i].Type = "$(datacenters.items.0.type)"
			m.ELBs.Items[i].DatacenterType = "$(datacenters.items.0.type)"
			m.ELBs.Items[i].DatacenterName = "$(datacenters.items.0.name)"
			m.ELBs.Items[i].AccessKeyID = "$(datacenters.items.0.aws_access_key_id)"
			m.ELBs.Items[i].SecretAccessKey = "$(datacenters.items.0.aws_secret_access_key)"
			m.ELBs.Items[i].DatacenterRegion = "$(datacenters.items.0.region)"
			m.ELBs.Items[i].VpcID = "$(vpcs.items.0.vpc_id)"
		}
	}

	// Map ebs data
	for i, ebs := range m.EBSVolumes.Items {
		volume := om.FindEBSVolume(ebs.Name)
		if volume != nil {
			m.EBSVolumes.Items[i].VolumeAWSID = volume.VolumeAWSID
			m.EBSVolumes.Items[i].DatacenterType = "$(datacenters.items.0.type)"
			m.EBSVolumes.Items[i].DatacenterName = "$(datacenters.items.0.name)"
			m.EBSVolumes.Items[i].SecretAccessKey = "$(datacenters.items.0.aws_secret_access_key)"
			m.EBSVolumes.Items[i].AccessKeyID = "$(datacenters.items.0.aws_access_key_id)"
			m.EBSVolumes.Items[i].DatacenterRegion = "$(datacenters.items.0.region)"
		}
	}

	for i, zone := range m.Route53s.Items {
		z := om.FindRoute53(zone.Name)
		if z != nil {
			m.Route53s.Items[i].HostedZoneID = z.HostedZoneID
			m.Route53s.Items[i].DatacenterName = "$(datacenters.items.0.name)"
			m.Route53s.Items[i].SecretAccessKey = "$(datacenters.items.0.aws_secret_access_key)"
			m.Route53s.Items[i].AccessKeyID = "$(datacenters.items.0.aws_access_key_id)"
			m.Route53s.Items[i].DatacenterRegion = "$(datacenters.items.0.region)"
			m.Route53s.Items[i].VPCID = "$(vpcs.items.0.vpc_id)"
		}
	}

	for i, s3 := range m.S3s.Items {
		z := om.FindS3(s3.Name)
		if z != nil {
			m.S3s.Items[i].DatacenterName = "$(datacenters.items.0.name)"
			m.S3s.Items[i].SecretAccessKey = "$(datacenters.items.0.aws_secret_access_key)"
			m.S3s.Items[i].AccessKeyID = "$(datacenters.items.0.aws_access_key_id)"
			m.S3s.Items[i].DatacenterRegion = "$(datacenters.items.0.region)"
		}
	}

	for i, cluster := range m.RDSClusters.Items {
		c := om.FindRDSCluster(cluster.Name)
		if c != nil {
			m.RDSClusters.Items[i].ARN = c.ARN
			m.RDSClusters.Items[i].SecretAccessKey = "$(datacenters.items.0.aws_secret_access_key)"
			m.RDSClusters.Items[i].AccessKeyID = "$(datacenters.items.0.aws_access_key_id)"
			m.RDSClusters.Items[i].DatacenterRegion = "$(datacenters.items.0.region)"
			m.RDSClusters.Items[i].Endpoint = c.Endpoint
		}
	}

	for i, instance := range m.RDSInstances.Items {
		in := om.FindRDSInstance(instance.Name)
		if in != nil {
			m.RDSInstances.Items[i].ARN = in.ARN
			m.RDSInstances.Items[i].SecretAccessKey = "$(datacenters.items.0.aws_secret_access_key)"
			m.RDSInstances.Items[i].AccessKeyID = "$(datacenters.items.0.aws_access_key_id)"
			m.RDSInstances.Items[i].DatacenterRegion = "$(datacenters.items.0.region)"
			m.RDSInstances.Items[i].Endpoint = in.Endpoint
		}
	}
}

func mapTags(name, service string) map[string]string {
	tags := make(map[string]string)

	tags["Name"] = name
	tags["ernest.service"] = service

	return tags
}

func mapTagsServiceOnly(service string) map[string]string {
	tags := make(map[string]string)

	tags["ernest.service"] = service

	return tags
}

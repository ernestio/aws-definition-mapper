/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package integration

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"strconv"
	"testing"

	aes "github.com/ernestio/crypto/aes"
	"github.com/nats-io/nats"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAWSHappyPath(t *testing.T) {
	var service = "aws"
	var encryptedPwd string
	var encryptedUsr string

	crypto := aes.New()
	key := os.Getenv("ERNEST_CRYPTO_KEY")

	service = service + strconv.Itoa(rand.Intn(9999999))

	neSub := make(chan *nats.Msg, 1)
	inSub := make(chan *nats.Msg, 1)
	fiSub := make(chan *nats.Msg, 1)
	naSub := make(chan *nats.Msg, 1)
	lbSub := make(chan *nats.Msg, 1)
	s3Sub := make(chan *nats.Msg, 1)
	rdsSub := make(chan *nats.Msg, 1)

	basicSetup("aws")

	Convey("Given I have a non existing aws definition", t, func() {
		Convey("When I apply aws1.yml", func() {
			f := getDefinitionPathAWS("aws1.yml", service)

			subNeC, _ := n.ChanSubscribe("network.create.aws-fake", neSub)
			subInC, _ := n.ChanSubscribe("instance.create.aws-fake", inSub)
			subFiC, _ := n.ChanSubscribe("firewall.create.aws-fake", fiSub)

			_, err := ernest("service", "apply", f)

			Convey("Then I should create a valid service", func() {
				if err != nil {
					log.Println(err.Error())
				}

				event := awsNetworkEvent{}
				eventI := awsInstanceEvent{}
				eventF := awsFirewallEvent{}

				msg, err := waitMsg(neSub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &event)
				_ = subNeC.Unsubscribe()
				msg, err = waitMsg(inSub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &eventI)
				_ = subInC.Unsubscribe()
				msg, err = waitMsg(fiSub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &eventF)
				_ = subFiC.Unsubscribe()

				Info("And should call network creator connector with valid fields", " ", 6)
				So(event.Type, ShouldEqual, "aws-fake")
				So(event.DatacenterRegion, ShouldEqual, "fake")

				encryptedPwd, _ = crypto.Decrypt(event.DatacenterAccessKey, key)
				encryptedUsr, _ = crypto.Decrypt(event.DatacenterAccessToken, key)
				So("fake_up_to_16_characters", ShouldEqual, encryptedPwd)
				So("up_to_16_characters_secret", ShouldEqual, encryptedUsr)
				So(event.DatacenterVpcID, ShouldEqual, "fakeaws")
				So(event.NetworkSubnet, ShouldEqual, "10.1.0.0/24")

				Info("And should call firewall creator connector with valid fields", " ", 6)
				So(eventF.Type, ShouldEqual, "aws-fake")
				So(eventF.DatacenterRegion, ShouldEqual, "fake")
				encryptedPwd, _ = crypto.Decrypt(eventF.DatacenterAccessKey, key)
				encryptedUsr, _ = crypto.Decrypt(eventF.DatacenterAccessToken, key)
				So("fake_up_to_16_characters", ShouldEqual, encryptedPwd)
				So("up_to_16_characters_secret", ShouldEqual, encryptedUsr)
				So(eventF.DatacenterVPCID, ShouldEqual, "fakeaws")
				So(eventF.SecurityGroupName, ShouldEqual, "fakeaws-"+service+"-web-sg-1")
				So(len(eventF.SecurityGroupRules.Egress), ShouldEqual, 1)
				So(eventF.SecurityGroupRules.Egress[0].IP, ShouldEqual, "10.1.1.11/32")
				So(eventF.SecurityGroupRules.Egress[0].From, ShouldEqual, 80)
				So(eventF.SecurityGroupRules.Egress[0].To, ShouldEqual, 80)
				So(eventF.SecurityGroupRules.Egress[0].Protocol, ShouldEqual, "-1")
				So(len(eventF.SecurityGroupRules.Ingress), ShouldEqual, 1)
				So(eventF.SecurityGroupRules.Ingress[0].IP, ShouldEqual, "10.1.1.11/32")
				So(eventF.SecurityGroupRules.Ingress[0].From, ShouldEqual, 80)
				So(eventF.SecurityGroupRules.Ingress[0].To, ShouldEqual, 80)
				So(eventF.SecurityGroupRules.Ingress[0].Protocol, ShouldEqual, "-1")
				So(eventF.Status, ShouldEqual, "processing")

				Info("And should call instance creator connector with valid fields", " ", 6)
				So(eventI.Type, ShouldEqual, "aws-fake")
				So(eventI.DatacenterRegion, ShouldEqual, "fake")
				encryptedPwd, _ = crypto.Decrypt(eventI.DatacenterAccessKey, key)
				encryptedUsr, _ = crypto.Decrypt(eventI.DatacenterAccessToken, key)
				So("fake_up_to_16_characters", ShouldEqual, encryptedPwd)
				So("up_to_16_characters_secret", ShouldEqual, encryptedUsr)
				So(eventI.DatacenterVpcID, ShouldEqual, "fakeaws")
				So(eventI.NetworkAWSID, ShouldEqual, "foo")
				So(len(eventI.SecurityGroupAWSIDs), ShouldEqual, 1)
				So(eventI.SecurityGroupAWSIDs[0], ShouldEqual, "foo")
				So(eventI.InstanceName, ShouldEqual, "fakeaws-"+service+"-web-1")
				So(eventI.InstanceImage, ShouldEqual, "ami-6666f915")
				So(eventI.InstanceType, ShouldEqual, "e1.micro")
				So(eventI.Status, ShouldEqual, "processing")

			})

		})

		Convey("When I apply aws2.yml", func() {
			f := getDefinitionPathAWS("aws2.yml", service)
			subInC, _ := n.ChanSubscribe("instance.create.aws-fake", inSub)
			_, err := ernest("service", "apply", f)
			Convey("Then it should create a new xx-web-2 instance", func() {
				if err != nil {
					log.Println(err.Error())
				}

				eventI := awsInstanceEvent{}

				msg, err := waitMsg(inSub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &eventI)
				_ = subInC.Unsubscribe()

				Info("And should call instance creator connector with valid fields", " ", 6)
				So(eventI.Type, ShouldEqual, "aws-fake")
				So(eventI.DatacenterRegion, ShouldEqual, "fake")
				encryptedPwd, _ = crypto.Decrypt(eventI.DatacenterAccessKey, key)
				encryptedUsr, _ = crypto.Decrypt(eventI.DatacenterAccessToken, key)
				So("fake_up_to_16_characters", ShouldEqual, encryptedPwd)
				So("up_to_16_characters_secret", ShouldEqual, encryptedUsr)
				So(eventI.DatacenterVpcID, ShouldEqual, "fakeaws")
				So(eventI.NetworkAWSID, ShouldEqual, "foo")
				So(len(eventI.SecurityGroupAWSIDs), ShouldEqual, 1)
				So(eventI.SecurityGroupAWSIDs[0], ShouldEqual, "foo")
				So(eventI.InstanceName, ShouldEqual, "fakeaws-"+service+"-web-2")
				So(eventI.InstanceImage, ShouldEqual, "ami-6666f915")
				So(eventI.InstanceType, ShouldEqual, "e1.micro")
				So(eventI.Status, ShouldEqual, "processing")
			})
		})

		Convey("When I apply aws3.yml", func() {
			f := getDefinitionPathAWS("aws3.yml", service)
			subInD, _ := n.ChanSubscribe("instance.delete.aws-fake", inSub)
			_, err := ernest("service", "apply", f)
			Convey("Then it should delete xx-web-2 instance", func() {
				if err != nil {
					log.Println(err.Error())
				}

				eventI := awsInstanceEvent{}

				msg, err := waitMsg(inSub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &eventI)
				_ = subInD.Unsubscribe()

				Info("And should call instance creator connector with valid fields", " ", 6)
				So(eventI.Type, ShouldEqual, "aws-fake")
				So(eventI.DatacenterRegion, ShouldEqual, "fake")
				encryptedPwd, _ = crypto.Decrypt(eventI.DatacenterAccessKey, key)
				encryptedUsr, _ = crypto.Decrypt(eventI.DatacenterAccessToken, key)
				So("fake_up_to_16_characters", ShouldEqual, encryptedPwd)
				So("up_to_16_characters_secret", ShouldEqual, encryptedUsr)
				So(eventI.DatacenterVpcID, ShouldEqual, "fakeaws")
				So(eventI.NetworkAWSID, ShouldEqual, "foo")
				So(len(eventI.SecurityGroupAWSIDs), ShouldEqual, 1)
				So(eventI.SecurityGroupAWSIDs[0], ShouldEqual, "foo")
				So(eventI.InstanceName, ShouldEqual, "fakeaws-"+service+"-web-2")
				So(eventI.InstanceImage, ShouldEqual, "ami-6666f915")
				So(eventI.InstanceType, ShouldEqual, "e1.micro")
				So(eventI.Status, ShouldEqual, "processing")
			})
		})

		Convey("When I apply aws4.yml", func() {
			f := getDefinitionPathAWS("aws4.yml", service)
			subInC, _ := n.ChanSubscribe("instance.update.aws-fake", inSub)
			_, err := ernest("service", "apply", f)
			Convey("Then it should update xx-web-1 instance", func() {
				if err != nil {
					log.Println(err.Error())
				}

				eventI := awsInstanceEvent{}

				msg, err := waitMsg(inSub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &eventI)
				_ = subInC.Unsubscribe()

				Info("And should call instance creator connector with valid fields", " ", 6)
				So(eventI.Type, ShouldEqual, "aws-fake")
				So(eventI.DatacenterRegion, ShouldEqual, "fake")
				encryptedPwd, _ = crypto.Decrypt(eventI.DatacenterAccessKey, key)
				encryptedUsr, _ = crypto.Decrypt(eventI.DatacenterAccessToken, key)
				So("fake_up_to_16_characters", ShouldEqual, encryptedPwd)
				So("up_to_16_characters_secret", ShouldEqual, encryptedUsr)
				So(eventI.DatacenterVpcID, ShouldEqual, "fakeaws")
				So(eventI.NetworkAWSID, ShouldEqual, "foo")
				So(len(eventI.SecurityGroupAWSIDs), ShouldEqual, 0)
				So(eventI.InstanceName, ShouldEqual, "fakeaws-"+service+"-web-1")
				So(eventI.InstanceImage, ShouldEqual, "ami-6666f915")
				So(eventI.InstanceType, ShouldEqual, "e1.micro")
				So(eventI.Status, ShouldEqual, "processing")
			})
		})

		Convey("When I apply aws5.yml", func() {
			f := getDefinitionPathAWS("aws5.yml", service)
			subFiU, _ := n.ChanSubscribe("firewall.update.aws-fake", fiSub)
			_, err := ernest("service", "apply", f)
			Convey("Then it should add an Ingress rule to existing firewall", func() {
				if err != nil {
					log.Println(err.Error())
				}

				eventF := awsFirewallEvent{}

				msg, err := waitMsg(fiSub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &eventF)
				_ = subFiU.Unsubscribe()

				Info("And should call firewall updater connector with valid fields", " ", 6)
				So(eventF.Type, ShouldEqual, "aws-fake")
				So(eventF.DatacenterRegion, ShouldEqual, "fake")
				encryptedPwd, _ = crypto.Decrypt(eventF.DatacenterAccessKey, key)
				encryptedUsr, _ = crypto.Decrypt(eventF.DatacenterAccessToken, key)
				So("fake_up_to_16_characters", ShouldEqual, encryptedPwd)
				So("up_to_16_characters_secret", ShouldEqual, encryptedUsr)
				So(eventF.DatacenterVPCID, ShouldEqual, "fakeaws")
				So(eventF.SecurityGroupName, ShouldEqual, "fakeaws-"+service+"-web-sg-1")
				So(len(eventF.SecurityGroupRules.Egress), ShouldEqual, 1)
				So(eventF.SecurityGroupRules.Egress[0].IP, ShouldEqual, "10.1.1.11/32")
				So(eventF.SecurityGroupRules.Egress[0].From, ShouldEqual, 80)
				So(eventF.SecurityGroupRules.Egress[0].To, ShouldEqual, 80)
				So(eventF.SecurityGroupRules.Egress[0].Protocol, ShouldEqual, "-1")
				So(len(eventF.SecurityGroupRules.Ingress), ShouldEqual, 2)
				So(eventF.SecurityGroupRules.Ingress[0].IP, ShouldEqual, "10.1.1.11/32")
				So(eventF.SecurityGroupRules.Ingress[0].From, ShouldEqual, 80)
				So(eventF.SecurityGroupRules.Ingress[0].To, ShouldEqual, 80)
				So(eventF.SecurityGroupRules.Ingress[0].Protocol, ShouldEqual, "-1")
				So(eventF.SecurityGroupRules.Ingress[1].IP, ShouldEqual, "10.1.1.11/32")
				So(eventF.SecurityGroupRules.Ingress[1].From, ShouldEqual, 22)
				So(eventF.SecurityGroupRules.Ingress[1].To, ShouldEqual, 22)
				So(eventF.SecurityGroupRules.Ingress[1].Protocol, ShouldEqual, "-1")
				So(eventF.Status, ShouldEqual, "processing")
			})
		})

		Convey("When I apply aws6.yml", func() {
			f := getDefinitionPathAWS("aws6.yml", service)
			subFiU, _ := n.ChanSubscribe("firewall.update.aws-fake", fiSub)
			_, err := ernest("service", "apply", f)
			Convey("Then it should add an Egress rule to existing firewall", func() {
				if err != nil {
					log.Println(err.Error())
				}

				eventF := awsFirewallEvent{}

				msg, err := waitMsg(fiSub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &eventF)
				_ = subFiU.Unsubscribe()

				Info("And should call firewall updater connector with valid fields", " ", 6)
				So(eventF.Type, ShouldEqual, "aws-fake")
				So(eventF.DatacenterRegion, ShouldEqual, "fake")
				encryptedPwd, _ = crypto.Decrypt(eventF.DatacenterAccessKey, key)
				encryptedUsr, _ = crypto.Decrypt(eventF.DatacenterAccessToken, key)
				So("fake_up_to_16_characters", ShouldEqual, encryptedPwd)
				So("up_to_16_characters_secret", ShouldEqual, encryptedUsr)
				So(eventF.DatacenterVPCID, ShouldEqual, "fakeaws")
				So(eventF.SecurityGroupName, ShouldEqual, "fakeaws-"+service+"-web-sg-1")
				So(len(eventF.SecurityGroupRules.Egress), ShouldEqual, 2)
				So(eventF.SecurityGroupRules.Egress[0].IP, ShouldEqual, "10.1.1.11/32")
				So(eventF.SecurityGroupRules.Egress[0].From, ShouldEqual, 80)
				So(eventF.SecurityGroupRules.Egress[0].To, ShouldEqual, 80)
				So(eventF.SecurityGroupRules.Egress[0].Protocol, ShouldEqual, "-1")
				So(eventF.SecurityGroupRules.Egress[1].IP, ShouldEqual, "10.1.1.11/32")
				So(eventF.SecurityGroupRules.Egress[1].From, ShouldEqual, 22)
				So(eventF.SecurityGroupRules.Egress[1].To, ShouldEqual, 22)
				So(eventF.SecurityGroupRules.Egress[1].Protocol, ShouldEqual, "-1")
				So(len(eventF.SecurityGroupRules.Ingress), ShouldEqual, 2)
				So(eventF.SecurityGroupRules.Ingress[0].IP, ShouldEqual, "10.1.1.11/32")
				So(eventF.SecurityGroupRules.Ingress[0].From, ShouldEqual, 80)
				So(eventF.SecurityGroupRules.Ingress[0].To, ShouldEqual, 80)
				So(eventF.SecurityGroupRules.Ingress[0].Protocol, ShouldEqual, "-1")
				So(eventF.SecurityGroupRules.Ingress[1].IP, ShouldEqual, "10.1.1.11/32")
				So(eventF.SecurityGroupRules.Ingress[1].From, ShouldEqual, 22)
				So(eventF.SecurityGroupRules.Ingress[1].To, ShouldEqual, 22)
				So(eventF.SecurityGroupRules.Ingress[1].Protocol, ShouldEqual, "-1")
				So(eventF.Status, ShouldEqual, "processing")
			})
		})

		Convey("When I apply aws7.yml", func() {
			f := getDefinitionPathAWS("aws7.yml", service)
			subFiU, _ := n.ChanSubscribe("firewall.update.aws-fake", fiSub)
			_, err := ernest("service", "apply", f)
			Convey("Then it should delete previously added egress and ingress rules from  existing firewall", func() {
				if err != nil {
					log.Println(err.Error())
				}

				eventF := awsFirewallEvent{}

				msg, err := waitMsg(fiSub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &eventF)
				_ = subFiU.Unsubscribe()

				Info("And should call firewall updater connector with valid fields", " ", 6)
				So(eventF.Type, ShouldEqual, "aws-fake")
				So(eventF.DatacenterRegion, ShouldEqual, "fake")
				encryptedPwd, _ = crypto.Decrypt(eventF.DatacenterAccessKey, key)
				encryptedUsr, _ = crypto.Decrypt(eventF.DatacenterAccessToken, key)
				So("fake_up_to_16_characters", ShouldEqual, encryptedPwd)
				So("up_to_16_characters_secret", ShouldEqual, encryptedUsr)
				So(eventF.DatacenterVPCID, ShouldEqual, "fakeaws")
				So(eventF.SecurityGroupName, ShouldEqual, "fakeaws-"+service+"-web-sg-1")
				So(len(eventF.SecurityGroupRules.Egress), ShouldEqual, 1)
				So(eventF.SecurityGroupRules.Egress[0].IP, ShouldEqual, "10.1.1.11/32")
				So(eventF.SecurityGroupRules.Egress[0].From, ShouldEqual, 80)
				So(eventF.SecurityGroupRules.Egress[0].To, ShouldEqual, 80)
				So(eventF.SecurityGroupRules.Egress[0].Protocol, ShouldEqual, "-1")
				So(len(eventF.SecurityGroupRules.Ingress), ShouldEqual, 1)
				So(eventF.SecurityGroupRules.Ingress[0].IP, ShouldEqual, "10.1.1.11/32")
				So(eventF.SecurityGroupRules.Ingress[0].From, ShouldEqual, 80)
				So(eventF.SecurityGroupRules.Ingress[0].To, ShouldEqual, 80)
				So(eventF.SecurityGroupRules.Ingress[0].Protocol, ShouldEqual, "-1")
				So(eventF.Status, ShouldEqual, "processing")
			})
		})

		Convey("When I apply aws8.yml", func() {
			f := getDefinitionPathAWS("aws8.yml", service)
			subNeC, _ := n.ChanSubscribe("network.create.aws-fake", neSub)
			_, err := ernest("service", "apply", f)
			Convey("Then it should create the new 10.2.0.0/24 network", func() {
				if err != nil {
					log.Println(err.Error())
				}

				event := awsNetworkEvent{}

				msg, err := waitMsg(neSub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &event)
				_ = subNeC.Unsubscribe()

				Info("And should call network creator connector with valid fields", " ", 6)
				So(event.Type, ShouldEqual, "aws-fake")
				So(event.DatacenterRegion, ShouldEqual, "fake")
				encryptedPwd, _ = crypto.Decrypt(event.DatacenterAccessKey, key)
				encryptedUsr, _ = crypto.Decrypt(event.DatacenterAccessToken, key)
				So("fake_up_to_16_characters", ShouldEqual, encryptedPwd)
				So("up_to_16_characters_secret", ShouldEqual, encryptedUsr)
				So(event.DatacenterVpcID, ShouldEqual, "fakeaws")
				So(event.NetworkSubnet, ShouldEqual, "10.2.0.0/24")
			})
		})

		Convey("When I apply aws9.yml", func() {
			f := getDefinitionPathAWS("aws9.yml", service)
			subNeC, _ := n.ChanSubscribe("network.delete.aws-fake", neSub)
			_, err := ernest("service", "apply", f)
			Convey("Then it should delete network 10.2.0.0/24", func() {
				if err != nil {
					log.Println(err.Error())
				}

				event := awsNetworkEvent{}

				msg, err := waitMsg(neSub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &event)
				_ = subNeC.Unsubscribe()

				Info("And should call network deleter connector with valid fields", " ", 6)
				So(event.Type, ShouldEqual, "aws-fake")
				So(event.DatacenterRegion, ShouldEqual, "fake")
				encryptedPwd, _ = crypto.Decrypt(event.DatacenterAccessKey, key)
				encryptedUsr, _ = crypto.Decrypt(event.DatacenterAccessToken, key)
				So("fake_up_to_16_characters", ShouldEqual, encryptedPwd)
				So("up_to_16_characters_secret", ShouldEqual, encryptedUsr)
				So(event.DatacenterVpcID, ShouldEqual, "fakeaws")
				So(event.NetworkSubnet, ShouldEqual, "10.2.0.0/24")

			})
		})

		Convey("When I apply aws10.yml", func() {
			f := getDefinitionPathAWS("aws10.yml", service)
			subNeC, _ := n.ChanSubscribe("network.create.aws-fake", neSub)
			subInC, _ := n.ChanSubscribe("instance.create.aws-fake", inSub)
			_, err := ernest("service", "apply", f)
			Convey("Then it should create the new 10.2.0.0/24 network", func() {
				if err != nil {
					log.Println(err.Error())
				}

				event := awsNetworkEvent{}

				msg, err := waitMsg(neSub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &event)
				_ = subNeC.Unsubscribe()

				eventI := awsInstanceEvent{}

				msg, err = waitMsg(inSub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &eventI)
				_ = subInC.Unsubscribe()

				Info("And should call instance creator connector with valid fields", " ", 6)
				So(eventI.Type, ShouldEqual, "aws-fake")
				So(eventI.DatacenterRegion, ShouldEqual, "fake")
				encryptedPwd, _ = crypto.Decrypt(eventI.DatacenterAccessKey, key)
				encryptedUsr, _ = crypto.Decrypt(eventI.DatacenterAccessToken, key)
				So("fake_up_to_16_characters", ShouldEqual, encryptedPwd)
				So("up_to_16_characters_secret", ShouldEqual, encryptedUsr)
				So(eventI.DatacenterVpcID, ShouldEqual, "fakeaws")
				So(eventI.InstanceName, ShouldEqual, "fakeaws-"+service+"-bknd-1")
				So(eventI.InstanceImage, ShouldEqual, "ami-6666f915")
				So(eventI.InstanceType, ShouldEqual, "e1.micro")
				So(eventI.Status, ShouldEqual, "processing")

				Info("And should call network creator connector with valid fields", " ", 6)
				So(event.Type, ShouldEqual, "aws-fake")
				So(event.DatacenterRegion, ShouldEqual, "fake")
				encryptedPwd, _ = crypto.Decrypt(event.DatacenterAccessKey, key)
				encryptedUsr, _ = crypto.Decrypt(event.DatacenterAccessToken, key)
				So("fake_up_to_16_characters", ShouldEqual, encryptedPwd)
				So("up_to_16_characters_secret", ShouldEqual, encryptedUsr)
				So(event.DatacenterVpcID, ShouldEqual, "fakeaws")
				So(event.NetworkSubnet, ShouldEqual, "10.2.0.0/24")
			})
		})

		Convey("When I apply aws11.yml", func() {
			f := getDefinitionPathAWS("aws11.yml", service)
			subNeD, _ := n.ChanSubscribe("network.delete.aws-fake", neSub)
			subInD, _ := n.ChanSubscribe("instance.delete.aws-fake", inSub)
			_, err := ernest("service", "apply", f)
			Convey("Then it should delete the 10.2.0.0/24 network", func() {
				if err != nil {
					log.Println(err.Error())
				}

				eventI := awsInstanceEvent{}

				msg, err := waitMsg(inSub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &eventI)
				_ = subInD.Unsubscribe()

				event := awsNetworkEvent{}

				msg, err = waitMsg(neSub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &event)
				_ = subNeD.Unsubscribe()

				Info("And should call instance deleter connector with valid fields", " ", 6)
				So(eventI.Type, ShouldEqual, "aws-fake")
				So(eventI.DatacenterRegion, ShouldEqual, "fake")
				encryptedPwd, _ = crypto.Decrypt(eventI.DatacenterAccessKey, key)
				encryptedUsr, _ = crypto.Decrypt(eventI.DatacenterAccessToken, key)
				So("fake_up_to_16_characters", ShouldEqual, encryptedPwd)
				So("up_to_16_characters_secret", ShouldEqual, encryptedUsr)
				So(eventI.InstanceName, ShouldEqual, "fakeaws-"+service+"-bknd-1")
				So(eventI.InstanceImage, ShouldEqual, "ami-6666f915")
				So(eventI.InstanceType, ShouldEqual, "e1.micro")
				So(eventI.Status, ShouldEqual, "processing")

				Info("And should call network deleter connector with valid fields", " ", 6)
				So(event.Type, ShouldEqual, "aws-fake")
				So(event.DatacenterRegion, ShouldEqual, "fake")
				encryptedPwd, _ = crypto.Decrypt(event.DatacenterAccessKey, key)
				encryptedUsr, _ = crypto.Decrypt(event.DatacenterAccessToken, key)
				So("fake_up_to_16_characters", ShouldEqual, encryptedPwd)
				So("up_to_16_characters_secret", ShouldEqual, encryptedUsr)
				So(event.DatacenterVpcID, ShouldEqual, "fakeaws")
				So(event.NetworkSubnet, ShouldEqual, "10.2.0.0/24")
			})
		})

		Convey("When I apply aws12.yml", func() {
			f := getDefinitionPathAWS("aws12.yml", service)
			subNeC, _ := n.ChanSubscribe("network.create.aws-fake", neSub)
			subNaC, _ := n.ChanSubscribe("nat.create.aws-fake", naSub)
			_, err := ernest("service", "apply", f)
			Convey("Then it should create the new 10.2.0.0/24 network", func() {
				if err != nil {
					log.Println(err.Error())
				}

				event := awsNetworkEvent{}

				msg, err := waitMsg(neSub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &event)
				_ = subNeC.Unsubscribe()

				eventN := awsNatEvent{}

				msg, err = waitMsg(naSub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &eventN)
				_ = subNaC.Unsubscribe()

				Info("And should call nat creator connector with valid fields", " ", 6)
				So(eventN.Type, ShouldEqual, "aws-fake")
				So(eventN.DatacenterRegion, ShouldEqual, "fake")
				encryptedPwd, _ = crypto.Decrypt(eventN.DatacenterAccessKey, key)
				encryptedUsr, _ = crypto.Decrypt(eventN.DatacenterAccessToken, key)
				So("fake_up_to_16_characters", ShouldEqual, encryptedPwd)
				So("up_to_16_characters_secret", ShouldEqual, encryptedUsr)
				So(eventN.DatacenterVPCID, ShouldEqual, "fakeaws")
				So(eventN.PublicNetwork, ShouldEqual, "fakeaws-"+service+"-web")
				So(len(eventN.RoutedNetworks), ShouldEqual, 1)
				So(eventN.RoutedNetworks[0], ShouldEqual, "fakeaws-"+service+"-db")
				So(eventN.Status, ShouldEqual, "processing")

				Info("And should call network creator connector with valid fields", " ", 6)
				So(event.Type, ShouldEqual, "aws-fake")
				So(event.DatacenterRegion, ShouldEqual, "fake")
				encryptedPwd, _ = crypto.Decrypt(event.DatacenterAccessKey, key)
				encryptedUsr, _ = crypto.Decrypt(event.DatacenterAccessToken, key)
				So("fake_up_to_16_characters", ShouldEqual, encryptedPwd)
				So("up_to_16_characters_secret", ShouldEqual, encryptedUsr)
				So(event.DatacenterVpcID, ShouldEqual, "fakeaws")
				So(event.NetworkSubnet, ShouldEqual, "10.2.0.0/24")
				So(event.NetworkIsPublic, ShouldBeFalse)
			})
		})

		Convey("When I apply aws13.yml", func() {
			f := getDefinitionPathAWS("aws13.yml", service)
			subLBC, _ := n.ChanSubscribe("elb.create.aws-fake", lbSub)
			subS3, _ := n.ChanSubscribe("s3.create.aws-fake", s3Sub)
			_, err := ernest("service", "apply", f)
			Convey("Then it should create the new elb-1 elb", func() {
				if err != nil {
					log.Println(err.Error())
				}

				eventLB := awsELBEvent{}
				msg, err := waitMsg(lbSub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &eventLB)
				_ = subLBC.Unsubscribe()

				eventS3 := awsS3Event{}
				msg, err = waitMsg(s3Sub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &eventS3)
				_ = subS3.Unsubscribe()

				Info("And should call elb creator connector with valid fields", " ", 6)
				So(eventLB.Type, ShouldEqual, "aws-fake")
				So(eventLB.DatacenterRegion, ShouldEqual, "fake")
				So(eventLB.DatacenterToken, ShouldEqual, "R80hNNKU3hCl2MbGC30w+azPA6lTrrtUDYQo6VWLUIVBDI7Yh5MCX3aarKGCQVkfXvc=")
				So(eventLB.DatacenterSecret, ShouldEqual, "iqMQUK8jivJ3tm4QAJSBcZzRoM7dTzQx/bVx6uUc1eUf6epaap4Xol/P0XUVAFxf")
				So(eventLB.VpcID, ShouldEqual, "fakeaws")
				So(eventLB.Name, ShouldEqual, "fakeaws-"+service+"-elb-1")
				So(len(eventLB.InstanceNames), ShouldEqual, 1)
				So(len(eventLB.InstanceAWSIDs), ShouldEqual, 1)
				So(len(eventLB.SecurityGroupAWSIDs), ShouldEqual, 1)
				So(eventLB.InstanceNames[0], ShouldEqual, "fakeaws-"+service+"-web-1")
				So(eventLB.SecurityGroupAWSIDs[0], ShouldEqual, "foo")
				So(len(eventLB.Listeners), ShouldEqual, 1)
				So(eventLB.Listeners[0].ToPort, ShouldEqual, 80)
				So(eventLB.Listeners[0].FromPort, ShouldEqual, 80)
				So(eventLB.Listeners[0].Protocol, ShouldEqual, "HTTP")
				So(eventLB.Listeners[0].SSLCert, ShouldEqual, "")

				Info("And should call s3 creator connector with valid fields", " ", 6)
				So(eventS3.Name, ShouldEqual, "bucket-1")
				So(eventS3.ACL, ShouldEqual, "")
				So(eventS3.BucketLocation, ShouldEqual, "eu-west-1")
				So(len(eventS3.Grantees), ShouldEqual, 1)
				g := eventS3.Grantees[0]
				So(g.ID, ShouldEqual, "foo@r3labs.io")
				So(g.Type, ShouldEqual, "emailaddress")
				So(g.Permissions, ShouldEqual, "FULL_CONTROL")
			})
		})

		Convey("When I apply aws14.yml", func() {
			f := getDefinitionPathAWS("aws14.yml", service)
			subLBU, _ := n.ChanSubscribe("elb.update.aws-fake", lbSub)
			subS3, _ := n.ChanSubscribe("s3.update.aws-fake", s3Sub)
			_, err := ernest("service", "apply", f)
			Convey("Then it should update the elb-1 elb", func() {
				if err != nil {
					log.Println(err.Error())
				}

				eventLB := awsELBEvent{}
				msg, err := waitMsg(lbSub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &eventLB)
				_ = subLBU.Unsubscribe()

				eventS3 := awsS3Event{}
				msg, err = waitMsg(s3Sub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &eventS3)
				_ = subS3.Unsubscribe()

				Info("And should call elb updater connector with valid fields", " ", 6)
				So(eventLB.Type, ShouldEqual, "aws-fake")
				So(eventLB.DatacenterRegion, ShouldEqual, "fake")
				So(eventLB.DatacenterToken, ShouldEqual, "R80hNNKU3hCl2MbGC30w+azPA6lTrrtUDYQo6VWLUIVBDI7Yh5MCX3aarKGCQVkfXvc=")
				So(eventLB.DatacenterSecret, ShouldEqual, "iqMQUK8jivJ3tm4QAJSBcZzRoM7dTzQx/bVx6uUc1eUf6epaap4Xol/P0XUVAFxf")
				So(eventLB.VpcID, ShouldEqual, "fakeaws")
				So(eventLB.Name, ShouldEqual, "fakeaws-"+service+"-elb-1")
				So(len(eventLB.InstanceNames), ShouldEqual, 1)
				So(len(eventLB.InstanceAWSIDs), ShouldEqual, 1)
				So(len(eventLB.SecurityGroupAWSIDs), ShouldEqual, 1)
				So(eventLB.InstanceNames[0], ShouldEqual, "fakeaws-"+service+"-web-1")
				So(eventLB.SecurityGroupAWSIDs[0], ShouldEqual, "foo")
				So(len(eventLB.Listeners), ShouldEqual, 2)
				So(eventLB.Listeners[0].ToPort, ShouldEqual, 80)
				So(eventLB.Listeners[0].FromPort, ShouldEqual, 80)
				So(eventLB.Listeners[0].Protocol, ShouldEqual, "HTTP")
				So(eventLB.Listeners[0].SSLCert, ShouldEqual, "")
				So(eventLB.Listeners[1].ToPort, ShouldEqual, 443)
				So(eventLB.Listeners[1].FromPort, ShouldEqual, 443)
				So(eventLB.Listeners[1].Protocol, ShouldEqual, "HTTPS")
				So(eventLB.Listeners[1].SSLCert, ShouldEqual, "foo")

				Info("And should call s3 creator connector with valid fields", " ", 6)
				So(eventS3.Name, ShouldEqual, "bucket-1")
				So(eventS3.ACL, ShouldEqual, "")
				So(eventS3.BucketLocation, ShouldEqual, "eu-west-1")
				So(len(eventS3.Grantees), ShouldEqual, 2)
				g := eventS3.Grantees[0]
				So(g.ID, ShouldEqual, "foo@r3labs.io")
				So(g.Type, ShouldEqual, "emailaddress")
				So(g.Permissions, ShouldEqual, "FULL_CONTROL")
				g = eventS3.Grantees[1]
				So(g.ID, ShouldEqual, "bar@r3labs.io")
				So(g.Type, ShouldEqual, "emailaddress")
				So(g.Permissions, ShouldEqual, "WRITE")
			})
		})

		Convey("When I apply aws15.yml", func() {
			f := getDefinitionPathAWS("aws15.yml", service)
			subLBD, _ := n.ChanSubscribe("elb.delete.aws-fake", lbSub)
			subS3, _ := n.ChanSubscribe("s3.delete.aws-fake", s3Sub)
			_, err := ernest("service", "apply", f)
			Convey("Then it should delete the elb-1 elb", func() {
				if err != nil {
					log.Println(err.Error())
				}

				eventLB := awsELBEvent{}

				msg, err := waitMsg(lbSub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &eventLB)
				_ = subLBD.Unsubscribe()

				eventS3 := awsS3Event{}
				msg, err = waitMsg(s3Sub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &eventS3)
				_ = subS3.Unsubscribe()

				Info("And should call elb updater connector with valid fields", " ", 6)
				So(eventLB.Type, ShouldEqual, "aws-fake")
				So(eventLB.DatacenterRegion, ShouldEqual, "fake")
				So(eventLB.DatacenterToken, ShouldEqual, "R80hNNKU3hCl2MbGC30w+azPA6lTrrtUDYQo6VWLUIVBDI7Yh5MCX3aarKGCQVkfXvc=")
				So(eventLB.DatacenterSecret, ShouldEqual, "iqMQUK8jivJ3tm4QAJSBcZzRoM7dTzQx/bVx6uUc1eUf6epaap4Xol/P0XUVAFxf")
				So(eventLB.VpcID, ShouldEqual, "fakeaws")
				So(eventLB.Name, ShouldEqual, "fakeaws-"+service+"-elb-1")

				Info("And should call s3 creator connector with valid fields", " ", 6)
				So(eventS3.Name, ShouldEqual, "bucket-1")
				So(eventS3.ACL, ShouldEqual, "")
				So(eventS3.BucketLocation, ShouldEqual, "eu-west-1")
				So(len(eventS3.Grantees), ShouldEqual, 2)
				g := eventS3.Grantees[0]
				So(g.ID, ShouldEqual, "foo@r3labs.io")
				So(g.Type, ShouldEqual, "emailaddress")
				So(g.Permissions, ShouldEqual, "FULL_CONTROL")
				g = eventS3.Grantees[1]
				So(g.ID, ShouldEqual, "bar@r3labs.io")
				So(g.Type, ShouldEqual, "emailaddress")
				So(g.Permissions, ShouldEqual, "WRITE")
			})
		})

		Convey("When I apply aws16.yml", func() {
			f := getDefinitionPathAWS("aws16.yml", service)
			subRDS, _ := n.ChanSubscribe("rds_cluster.create.aws-fake", rdsSub)
			_, err := ernest("service", "apply", f)
			Convey("Then it should update the 'aurora' rds cluster", func() {
				if err != nil {
					log.Println(err.Error())
				}

				eventRDS := awsRDSClusterEvent{}
				msg, err := waitMsg(rdsSub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &eventRDS)
				_ = subRDS.Unsubscribe()

				Info("And should call rds cluster creator connector with valid fields", " ", 6)
				So(eventRDS.ProviderType, ShouldEqual, "aws-fake")
				So(eventRDS.DatacenterRegion, ShouldEqual, "fake")
				So(eventRDS.DatacenterToken, ShouldEqual, "R80hNNKU3hCl2MbGC30w+azPA6lTrrtUDYQo6VWLUIVBDI7Yh5MCX3aarKGCQVkfXvc=")
				So(eventRDS.DatacenterSecret, ShouldEqual, "iqMQUK8jivJ3tm4QAJSBcZzRoM7dTzQx/bVx6uUc1eUf6epaap4Xol/P0XUVAFxf")
				So(eventRDS.Name, ShouldEqual, "fakeaws-"+service+"-aurora")
				So(eventRDS.Engine, ShouldEqual, "aurora")
				So(eventRDS.Port, ShouldEqual, 3306)
				So(eventRDS.DatabaseName, ShouldEqual, "test")
				So(eventRDS.DatabaseUsername, ShouldEqual, "test")
				So(eventRDS.DatabasePassword, ShouldEqual, "testpass")
				So(eventRDS.BackupRetention, ShouldEqual, 1)
			})
		})

		Convey("When I apply aws17.yml", func() {
			f := getDefinitionPathAWS("aws17.yml", service)
			subRDS, _ := n.ChanSubscribe("rds_cluster.update.aws-fake", rdsSub)
			_, err := ernest("service", "apply", f)
			Convey("Then it should update the 'aurora' rds cluster", func() {
				if err != nil {
					log.Println(err.Error())
				}

				eventRDS := awsRDSClusterEvent{}
				msg, err := waitMsg(rdsSub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &eventRDS)
				_ = subRDS.Unsubscribe()

				Info("And should call rds cluster creator connector with valid fields", " ", 6)
				So(eventRDS.ProviderType, ShouldEqual, "aws-fake")
				So(eventRDS.DatacenterRegion, ShouldEqual, "fake")
				So(eventRDS.DatacenterToken, ShouldEqual, "R80hNNKU3hCl2MbGC30w+azPA6lTrrtUDYQo6VWLUIVBDI7Yh5MCX3aarKGCQVkfXvc=")
				So(eventRDS.DatacenterSecret, ShouldEqual, "iqMQUK8jivJ3tm4QAJSBcZzRoM7dTzQx/bVx6uUc1eUf6epaap4Xol/P0XUVAFxf")
				So(eventRDS.Name, ShouldEqual, "fakeaws-"+service+"-aurora")
				So(eventRDS.Engine, ShouldEqual, "aurora")
				So(eventRDS.Port, ShouldEqual, 3306)
				So(eventRDS.DatabaseName, ShouldEqual, "test")
				So(eventRDS.DatabaseUsername, ShouldEqual, "test")
				So(eventRDS.DatabasePassword, ShouldEqual, "testpass-2")
				So(eventRDS.BackupRetention, ShouldEqual, 1)
			})
		})

		Convey("When I apply aws18.yml", func() {
			f := getDefinitionPathAWS("aws18.yml", service)
			subRDS, _ := n.ChanSubscribe("rds_instance.create.aws-fake", rdsSub)
			_, err := ernest("service", "apply", f)
			Convey("Then it should create the 'test-1' rds instance", func() {
				if err != nil {
					log.Println(err.Error())
				}

				eventRDS := awsRDSInstanceEvent{}
				msg, err := waitMsg(rdsSub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &eventRDS)
				_ = subRDS.Unsubscribe()

				Info("And should call rds cluster creator connector with valid fields", " ", 6)
				So(eventRDS.ProviderType, ShouldEqual, "aws-fake")
				So(eventRDS.DatacenterRegion, ShouldEqual, "fake")
				So(eventRDS.DatacenterToken, ShouldEqual, "R80hNNKU3hCl2MbGC30w+azPA6lTrrtUDYQo6VWLUIVBDI7Yh5MCX3aarKGCQVkfXvc=")
				So(eventRDS.DatacenterSecret, ShouldEqual, "iqMQUK8jivJ3tm4QAJSBcZzRoM7dTzQx/bVx6uUc1eUf6epaap4Xol/P0XUVAFxf")
				So(eventRDS.Name, ShouldEqual, "fakeaws-"+service+"-test-1")
				So(eventRDS.Size, ShouldEqual, "db.r3.large")
				So(eventRDS.Cluster, ShouldEqual, "fakeaws-"+service+"-aurora")
			})
		})

		Convey("When I apply aws19.yml", func() {
			f := getDefinitionPathAWS("aws19.yml", service)
			subRDS, _ := n.ChanSubscribe("rds_instance.update.aws-fake", rdsSub)
			_, err := ernest("service", "apply", f)
			Convey("Then it should update the 'test-1' rds instance", func() {
				if err != nil {
					log.Println(err.Error())
				}

				eventRDS := awsRDSInstanceEvent{}
				msg, err := waitMsg(rdsSub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &eventRDS)
				_ = subRDS.Unsubscribe()

				Info("And should call rds cluster creator connector with valid fields", " ", 6)
				So(eventRDS.ProviderType, ShouldEqual, "aws-fake")
				So(eventRDS.DatacenterRegion, ShouldEqual, "fake")
				So(eventRDS.DatacenterToken, ShouldEqual, "R80hNNKU3hCl2MbGC30w+azPA6lTrrtUDYQo6VWLUIVBDI7Yh5MCX3aarKGCQVkfXvc=")
				So(eventRDS.DatacenterSecret, ShouldEqual, "iqMQUK8jivJ3tm4QAJSBcZzRoM7dTzQx/bVx6uUc1eUf6epaap4Xol/P0XUVAFxf")
				So(eventRDS.Name, ShouldEqual, "fakeaws-"+service+"-test-1")
				So(eventRDS.Size, ShouldEqual, "db.r3.xlarge")
				So(eventRDS.Cluster, ShouldEqual, "fakeaws-"+service+"-aurora")
			})
		})

		Convey("When I apply aws20.yml", func() {
			f := getDefinitionPathAWS("aws20.yml", service)
			subRDS, _ := n.ChanSubscribe("rds_instance.delete.aws-fake", rdsSub)
			_, err := ernest("service", "apply", f)
			Convey("Then it should delete the 'test-1' rds instance", func() {
				if err != nil {
					log.Println(err.Error())
				}

				eventRDS := awsRDSInstanceEvent{}
				msg, err := waitMsg(rdsSub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &eventRDS)
				_ = subRDS.Unsubscribe()

				Info("And should call rds cluster creator connector with valid fields", " ", 6)
				So(eventRDS.ProviderType, ShouldEqual, "aws-fake")
				So(eventRDS.DatacenterRegion, ShouldEqual, "fake")
				So(eventRDS.DatacenterToken, ShouldEqual, "R80hNNKU3hCl2MbGC30w+azPA6lTrrtUDYQo6VWLUIVBDI7Yh5MCX3aarKGCQVkfXvc=")
				So(eventRDS.DatacenterSecret, ShouldEqual, "iqMQUK8jivJ3tm4QAJSBcZzRoM7dTzQx/bVx6uUc1eUf6epaap4Xol/P0XUVAFxf")
				So(eventRDS.Name, ShouldEqual, "fakeaws-"+service+"-test-1")
			})
		})

		Convey("When I apply aws21.yml", func() {
			f := getDefinitionPathAWS("aws21.yml", service)
			subRDS, _ := n.ChanSubscribe("rds_cluster.delete.aws-fake", rdsSub)
			_, err := ernest("service", "apply", f)
			Convey("Then it should delete the 'aurora' rds cluster", func() {
				if err != nil {
					log.Println(err.Error())
				}

				eventRDS := awsRDSClusterEvent{}
				msg, err := waitMsg(rdsSub)
				So(err, ShouldBeNil)
				_ = json.Unmarshal(msg.Data, &eventRDS)
				_ = subRDS.Unsubscribe()

				Info("And should call rds cluster creator connector with valid fields", " ", 6)
				So(eventRDS.ProviderType, ShouldEqual, "aws-fake")
				So(eventRDS.DatacenterRegion, ShouldEqual, "fake")
				So(eventRDS.DatacenterToken, ShouldEqual, "R80hNNKU3hCl2MbGC30w+azPA6lTrrtUDYQo6VWLUIVBDI7Yh5MCX3aarKGCQVkfXvc=")
				So(eventRDS.DatacenterSecret, ShouldEqual, "iqMQUK8jivJ3tm4QAJSBcZzRoM7dTzQx/bVx6uUc1eUf6epaap4Xol/P0XUVAFxf")
				So(eventRDS.Name, ShouldEqual, "fakeaws-"+service+"-aurora")
			})
		})

	})
}

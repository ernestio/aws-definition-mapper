@aws @add_and_attach_network
Feature: Service apply

  Scenario: Add a network and attach a new instance
    Given I setup ernest with target "https://ernest.local"
    And I setup a new service name
    When I'm logged in as "usr" / "pwd"
    And I apply the definition "aws9.yml"
    And I start recording
    And I apply the definition "aws10.yml"
    And I stop recording
    Then an event "network.create.aws-fake" should be called exactly "1" times
    And all "network.create.aws-fake" messages should contain a field "_type" with "aws-fake"
    And all "network.create.aws-fake" messages should contain a field "datacenter_region" with "fake"
    And all "network.create.aws-fake" messages should contain a field "vpc_id" with "fakeaws"
    And all "network.create.aws-fake" messages should contain a field "range" with "10.2.0.0/24"
    And all "network.create.aws-fake" messages should contain an encrypted field "aws_access_key_id" with "up_to_16_characters_secret"
    And all "network.create.aws-fake" messages should contain an encrypted field "aws_secret_access_key" with "fake_up_to_16_characters"
    Then an event "instance.create.aws-fake" should be called exactly "1" times
    And all "instance.create.aws-fake" messages should contain a field "_type" with "aws-fake"
    And all "instance.create.aws-fake" messages should contain a field "datacenter_region" with "fake"
    And all "instance.create.aws-fake" messages should contain an encrypted field "aws_access_key_id" with "up_to_16_characters_secret"
    And all "instance.create.aws-fake" messages should contain an encrypted field "aws_secret_access_key" with "fake_up_to_16_characters"
    And all "instance.create.aws-fake" messages should contain a field "vpc_id" with "fakeaws"
    And message "instance.create.aws-fake" number "0" should contain "null" as json field "security_group_aws_ids.0"
    And all "instance.create.aws-fake" messages should contain a field "name" with "fakeaws-$(name)-bknd-1"
    And all "instance.create.aws-fake" messages should contain a field "image" with "ami-6666f915"
    And all "instance.create.aws-fake" messages should contain a field "instance_type" with "e1.micro"
    And all "instance.create.aws-fake" messages should contain a field "status" with "processing"


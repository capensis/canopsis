# file: create_alarm.feature
Feature: create an alarm for a resource
  I need to be able to create an alarm for a resource
  An alarm document should be created in the db
  Three entities should be created in the db: for a component,
  for a connector, for a resource

  Scenario: given an event for a resource should create alarm and entities
    Given I am admin
    When I send an event:
    """
      {
        "connector" : "test_post_connector_create_alarm",
        "connector_name" : "test_post_connector_name_create_alarm",
        "source_type" : "resource",
        "event_type" : "check",
        "component" : "test_post_component_create_alarm",
        "resource" : "test_post_resource_create_alarm",
        "state" : 1,
        "output" : "noveo alarm"
      }
    """
    Then the response code should be 200
    When I wait the end of event processing
    Then an entity test_post_component_create_alarm should be in the db
    And an entity test_post_connector_create_alarm/test_post_connector_name_create_alarm should be in the db
    And an entity test_post_resource_create_alarm/test_post_component_create_alarm should be in the db
    And an alarm test_post_resource_create_alarm/test_post_component_create_alarm should be in the db


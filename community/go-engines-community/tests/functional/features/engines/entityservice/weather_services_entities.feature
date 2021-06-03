Feature: get service entities
  I need to be able to get service entities

  Scenario: given service for one entity should get one entity
    Given I am admin
    When I do POST /api/v4/entityservices:
    """
    {
      "_id": "test-service-weather-entity-1",
      "name": "test-service-weather-entity-1",
      "output_template": "Test-service-weather-entity-1",
      "category": "test-category-service-weather-entities",
      "enabled": true,
      "impact_level": 1,
      "entity_patterns": [
        {"name": "test-resource-service-weather-entity-1"}
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I send an event:
    """
    {
      "connector" : "test-connector-service-weather-entity-1",
      "connector_name" : "test-connector_name-service-weather-entity-1",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-service-weather-entity-1",
      "resource" : "test-resource-service-weather-entity-1",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services/test-service-weather-entity-1
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-resource-service-weather-entity-1/test-component-service-weather-entity-1",
          "name": "test-resource-service-weather-entity-1",
          "connector": "test-connector-service-weather-entity-1",
          "connector_name": "test-connector_name-service-weather-entity-1",
          "component": "test-component-service-weather-entity-1",
          "resource": "test-resource-service-weather-entity-1",
          "state": {"val": 2},
          "status": {"val": 1},
          "impact_state": 2,
          "impact_level": 1,
          "is_grey": false,
          "icon": "major",
          "ack": null,
          "snooze": null,
          "ticket": null,
          "infos": {},
          "pbehaviors": []
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """

  Scenario: given service for multiple entities should get multiple entities
    Given I am admin
    When I do POST /api/v4/entityservices:
    """
    {
      "_id": "test-service-weather-entity-2",
      "name": "test-service-weather-entity-2",
      "output_template": "Test-service-weather-entity-2",
      "category": "test-category-service-weather-entities",
      "enabled": true,
      "impact_level": 1,
      "entity_patterns": [
        {"name": "test-resource-service-weather-entity-2-1"},
        {"name": "test-resource-service-weather-entity-2-2"},
        {"name": "test-resource-service-weather-entity-2-3"}
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I send an event:
    """
    {
      "connector" : "test-connector-service-weather-entity-2",
      "connector_name" : "test-connector_name-service-weather-entity-2",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-service-weather-entity-2",
      "resource" : "test-resource-service-weather-entity-2-1",
      "state" : 1,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """
    {
      "connector" : "test-connector-service-weather-entity-2",
      "connector_name" : "test-connector_name-service-weather-entity-2",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-service-weather-entity-2",
      "resource" : "test-resource-service-weather-entity-2-2",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """
    {
      "connector" : "test-connector-service-weather-entity-2",
      "connector_name" : "test-connector_name-service-weather-entity-2",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-service-weather-entity-2",
      "resource" : "test-resource-service-weather-entity-2-3",
      "state" : 3,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services/test-service-weather-entity-2
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-resource-service-weather-entity-2-1/test-component-service-weather-entity-2",
          "name": "test-resource-service-weather-entity-2-1",
          "connector": "test-connector-service-weather-entity-2",
          "connector_name": "test-connector_name-service-weather-entity-2",
          "component": "test-component-service-weather-entity-2",
          "resource": "test-resource-service-weather-entity-2-1",
          "state": {"val": 1},
          "status": {"val": 1},
          "impact_state": 1,
          "impact_level": 1,
          "is_grey": false,
          "icon": "minor",
          "ack": null,
          "snooze": null,
          "ticket": null,
          "infos": {},
          "pbehaviors": []
        },
        {
          "_id": "test-resource-service-weather-entity-2-2/test-component-service-weather-entity-2",
          "name": "test-resource-service-weather-entity-2-2",
          "connector": "test-connector-service-weather-entity-2",
          "connector_name": "test-connector_name-service-weather-entity-2",
          "component": "test-component-service-weather-entity-2",
          "resource": "test-resource-service-weather-entity-2-2",
          "state": {"val": 2},
          "status": {"val": 1},
          "impact_state": 2,
          "impact_level": 1,
          "is_grey": false,
          "icon": "major",
          "ack": null,
          "snooze": null,
          "ticket": null,
          "infos": {},
          "pbehaviors": []
        },
        {
          "_id": "test-resource-service-weather-entity-2-3/test-component-service-weather-entity-2",
          "name": "test-resource-service-weather-entity-2-3",
          "connector": "test-connector-service-weather-entity-2",
          "connector_name": "test-connector_name-service-weather-entity-2",
          "component": "test-component-service-weather-entity-2",
          "resource": "test-resource-service-weather-entity-2-3",
          "state": {"val": 3},
          "status": {"val": 1},
          "impact_state": 3,
          "impact_level": 1,
          "is_grey": false,
          "icon": "critical",
          "ack": null,
          "snooze": null,
          "ticket": null,
          "infos": {},
          "pbehaviors": []
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 3
      }
    }
    """

  Scenario: given service for one entity and no open alarms should get one entity with ok state
    Given I am admin
    When I do POST /api/v4/entityservices:
    """
    {
      "_id": "test-service-weather-entity-3",
      "name": "test-service-weather-entity-3",
      "output_template": "Test-service-weather-entity-3",
      "category": "test-category-service-weather-entities",
      "enabled": true,
      "impact_level": 1,
      "entity_patterns": [
        {"name": "test-resource-service-weather-entity-3"}
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I send an event:
    """
    {
      "connector" : "test-connector-service-weather-entity-3",
      "connector_name" : "test-connector_name-service-weather-entity-3",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-service-weather-entity-3",
      "resource" : "test-resource-service-weather-entity-3",
      "state" : 0,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/weather-services/test-service-weather-entity-3
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-resource-service-weather-entity-3/test-component-service-weather-entity-3",
          "name": "test-resource-service-weather-entity-3",
          "connector": "",
          "connector_name": "",
          "component": "",
          "resource": "",
          "state": {"val": 0},
          "status": {"val": 0},
          "impact_state": 0,
          "impact_level": 1,
          "is_grey": false,
          "icon": "ok",
          "ack": null,
          "snooze": null,
          "ticket": null,
          "infos": {},
          "pbehaviors": []
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
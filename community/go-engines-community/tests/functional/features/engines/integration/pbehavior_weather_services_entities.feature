Feature: get service entities
  I need to be able to get service entities

  Scenario: given service for one entity with maintenance pbehavior
    should get one entity
    Given I am admin
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-entity-1",
      "connector_name": "test-connector-name-pbehavior-weather-service-entity-1",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-entity-1",
      "resource" : "test-resource-pbehavior-weather-service-entity-1",
      "state" : 0,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-pbehavior-weather-service-entity-1",
      "name": "test-pbehavior-weather-service-entity-1",
      "output_template": "test-pbehavior-weather-service-entity-1",
      "category": "test-category-pbehavior-weather-service-entity",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [{"name": "test-resource-pbehavior-weather-service-entity-1"}]
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-entity-1",
      "tstart": {{ now.Unix }},
      "tstop": {{ (now.Add (parseDuration "10m")).Unix }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-pbehavior-weather-service-entity-1"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait 1s
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-entity-1",
      "connector_name": "test-connector-name-pbehavior-weather-service-entity-1",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-entity-1",
      "resource" : "test-resource-pbehavior-weather-service-entity-1",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services/test-pbehavior-weather-service-entity-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-resource-pbehavior-weather-service-entity-1",
          "state": {"val": 2},
          "status": {"val": 1},
          "is_grey": true,
          "icon": "maintenance",
          "pbehaviors": [
            {
              "name": "test-pbehavior-weather-service-entity-1"
            }
          ]
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

  Scenario: given service for one entity with active pbehavior
    should get one entity
    Given I am admin
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-entity-2",
      "connector_name": "test-connector-name-pbehavior-weather-service-entity-2",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-entity-2",
      "resource" : "test-resource-pbehavior-weather-service-entity-2",
      "state" : 0,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-pbehavior-weather-service-entity-2",
      "name": "test-pbehavior-weather-service-entity-2",
      "output_template": "test-pbehavior-weather-service-entity-2",
      "category": "test-category-pbehavior-weather-service-entity",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [{"name": "test-resource-pbehavior-weather-service-entity-2"}]
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-entity-2",
      "tstart": {{ now.Unix }},
      "tstop": {{ (now.Add (parseDuration "10m")).Unix }},
      "type": "test-active-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-pbehavior-weather-service-entity-2"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait 1s
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-entity-2",
      "connector_name": "test-connector-name-pbehavior-weather-service-entity-2",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-entity-2",
      "resource" : "test-resource-pbehavior-weather-service-entity-2",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services/test-pbehavior-weather-service-entity-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-resource-pbehavior-weather-service-entity-2",
          "state": {"val": 2},
          "status": {"val": 1},
          "is_grey": false,
          "icon": "major",
          "pbehaviors": [
            {
              "name": "test-pbehavior-weather-service-entity-2"
            }
          ]
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

  Scenario: given service for one entity with maintenance pbehavior
    and another entity without pbehavior should get multiple entities
    Given I am admin
    When I send an event:
    """json
    {
      "connector" :  "test-connector-pbehavior-weather-service-entity-3",
      "connector_name": "test-connector-name-pbehavior-weather-service-entity-3",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-entity-3",
      "resource" : "test-resource-pbehavior-weather-service-entity-3-1",
      "state" : 0,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-pbehavior-weather-service-entity-3",
      "name": "test-pbehavior-weather-service-entity-3",
      "output_template": "test-pbehavior-weather-service-entity-3",
      "category": "test-category-pbehavior-weather-service-entity",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [
         {"name": "test-resource-pbehavior-weather-service-entity-3-1"},
         {"name": "test-resource-pbehavior-weather-service-entity-3-2"}
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-entity-3-1",
      "tstart": {{ now.Unix }},
      "tstop": {{ (now.Add (parseDuration "10m")).Unix }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-pbehavior-weather-service-entity-3-1"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait 1s
    When I send an event:
    """json
    {
      "connector" :  "test-connector-pbehavior-weather-service-entity-3",
      "connector_name": "test-connector-name-pbehavior-weather-service-entity-3",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-entity-3",
      "resource" : "test-resource-pbehavior-weather-service-entity-3-1",
      "state" : 3,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-entity-3",
      "connector_name": "test-connector-name-pbehavior-weather-service-entity-3",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-entity-3",
      "resource" : "test-resource-pbehavior-weather-service-entity-3-2",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services/test-pbehavior-weather-service-entity-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-resource-pbehavior-weather-service-entity-3-1",
          "state": {"val": 3},
          "status": {"val": 1},
          "is_grey": true,
          "icon": "maintenance",
          "pbehaviors": [
            {
              "name": "test-pbehavior-weather-service-entity-3-1"
            }
          ]
        },
        {
          "name": "test-resource-pbehavior-weather-service-entity-3-2",
          "state": {"val": 2},
          "status": {"val": 1},
          "is_grey": false,
          "icon": "major",
          "pbehaviors": []
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 2
      }
    }
    """

  Scenario: given service's entities should be marked grey if service is in pbehavior
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-pbehavior-weather-service-entity-4",
      "name": "test-pbehavior-weather-service-entity-4",
      "output_template": "test-pbehavior-weather-service-entity-4",
      "category": "test-category-pbehavior-weather-service-entity",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [
         {"name": "test-resource-pbehavior-weather-service-entity-4-1"},
         {"name": "test-resource-pbehavior-weather-service-entity-4-2"}
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-entity-4-1",
      "tstart": {{ now.Unix }},
      "tstop": {{ (now.Add (parseDuration "10m")).Unix }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "_id": "test-pbehavior-weather-service-entity-4"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait 1s
    When I send an event:
    """json
    {
      "connector" :  "test-connector-pbehavior-weather-service-entity-4",
      "connector_name": "test-connector-name-pbehavior-weather-service-entity-4",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-entity-4",
      "resource" : "test-resource-pbehavior-weather-service-entity-4-1",
      "state" : 3,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-entity-4",
      "connector_name": "test-connector-name-pbehavior-weather-service-entity-4",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-entity-4",
      "resource" : "test-resource-pbehavior-weather-service-entity-4-2",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services/test-pbehavior-weather-service-entity-4
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-resource-pbehavior-weather-service-entity-4-1",
          "state": {"val": 3},
          "status": {"val": 1},
          "is_grey": true,
          "pbehaviors": []
        },
        {
          "name": "test-resource-pbehavior-weather-service-entity-4-2",
          "state": {"val": 2},
          "status": {"val": 1},
          "is_grey": true,
          "pbehaviors": []
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 2
      }
    }
    """

  Scenario: given service's entities should be marked grey if service is in OK state and in pbehavior
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-pbehavior-weather-service-entity-5",
      "name": "test-pbehavior-weather-service-entity-5",
      "output_template": "test-pbehavior-weather-service-entity-5",
      "category": "test-category-pbehavior-weather-service-entity",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [
         {"name": "test-resource-pbehavior-weather-service-entity-5-1"},
         {"name": "test-resource-pbehavior-weather-service-entity-5-2"},
         {"name": "test-resource-pbehavior-weather-service-entity-5-3"}
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-entity-5-1",
      "comments":[{"author":"root","message":"pause test-pbehavior-weather-service-entity-5"}],
      "tstart": {{ now.Unix }},
      "tstop": {{ (now.Add (parseDuration "10m")).Unix }},
      "type": "test-default-pause-type",
      "reason": "test-reason-to-engine",
      "filter": {
        "_id": {
          "$in": [
            "test-resource-pbehavior-weather-service-entity-5-1/test-component-pbehavior-weather-service-entity-5",
            "test-resource-pbehavior-weather-service-entity-5-2/test-component-pbehavior-weather-service-entity-5"
          ]
        }
      }
    }
    """
    Then the response code should be 201
    When I wait 1s
    When I send an event:
    """json
    {
      "connector" :  "test-connector-pbehavior-weather-service-entity-5",
      "connector_name": "test-connector-name-pbehavior-weather-service-entity-5",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-entity-5",
      "resource" : "test-resource-pbehavior-weather-service-entity-5-1",
      "state" : 0,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-entity-5",
      "connector_name": "test-connector-name-pbehavior-weather-service-entity-5",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-entity-5",
      "resource" : "test-resource-pbehavior-weather-service-entity-5-2",
      "state" : 0,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-entity-5",
      "connector_name": "test-connector-name-pbehavior-weather-service-entity-5",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-entity-5",
      "resource" : "test-resource-pbehavior-weather-service-entity-5-3",
      "state" : 0,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I wait the next periodical process
    When I do GET /api/v4/weather-services/test-pbehavior-weather-service-entity-5
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-resource-pbehavior-weather-service-entity-5-1",
          "state": {"val": 0},
          "status": {"val": 0},
          "is_grey": true,
          "pbehaviors": [{"name": "test-pbehavior-weather-service-entity-5-1"}]
        },
        {
          "name": "test-resource-pbehavior-weather-service-entity-5-2",
          "state": {"val": 0},
          "status": {"val": 0},
          "is_grey": true,
          "pbehaviors": [{"name": "test-pbehavior-weather-service-entity-5-1"}]
        },
        {
          "name": "test-resource-pbehavior-weather-service-entity-5-3",
          "state": {"val": 0},
          "status": {"val": 0},
          "is_grey": false,
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

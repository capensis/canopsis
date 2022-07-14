Feature: get service weather
  I need to be able to get service weather

  Scenario: given one dependency with maintenance pbehavior should get maintenance icon and grey flag
    Given I am admin
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-1",
      "connector_name": "test-connector-name-pbehavior-weather-service-1",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-1",
      "resource" : "test-resource-pbehavior-weather-service-1",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-1",
      "output_template": "Test service weather 1",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [{"name": "test-resource-pbehavior-weather-service-1"}],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-1",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-pbehavior-weather-service-1"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services?filter={"name":"test-pbehavior-weather-service-1"}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-1",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "maintenance",
          "secondary_icon": "",
          "is_grey": true,
          "alarm_counters": [
            {
              "count": 1,
              "type": {
                "_id": "test-maintenance-type-to-engine",
                "description": "Engine maintenance",
                "icon_name": "test-maintenance-to-engine-icon",
                "name": "Engine maintenance",
                "priority": 18,
                "type": "maintenance"
              }
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

  Scenario: given one dependency with active pbehavior should get state icon
    Given I am admin
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-2",
      "connector_name": "test-connector-name-pbehavior-weather-service-2",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-2",
      "resource" : "test-resource-pbehavior-weather-service-2",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-2",
      "output_template": "Test service weather 2",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [{"name": "test-resource-pbehavior-weather-service-2"}],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-2",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-active-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-pbehavior-weather-service-2"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait the end of event processing
    When I do GET /api/v4/weather-services?filter={"name":"test-pbehavior-weather-service-2"}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-2",
          "state": {"val": 2},
          "status": {"val": 1},
          "icon": "major",
          "secondary_icon": "",
          "is_grey": false,
          "alarm_counters": []
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

  Scenario: given dependency with maintenance pbehavior
    and dependency without pbehavior should get secondary maintenance icon
    Given I am admin
    When I send an event:
    """json
    {
      "connector" :  "test-connector-pbehavior-weather-service-3",
      "connector_name": "test-connector-name-pbehavior-weather-service-3",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-3",
      "resource" : "test-resource-pbehavior-weather-service-3-1",
      "state" : 3,
      "output" : "noveo alarm"
    }
    """
    When I send an event:
    """json
    {
      "connector" :  "test-connector-pbehavior-weather-service-3",
      "connector_name": "test-connector-name-pbehavior-weather-service-3",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-3",
      "resource" : "test-resource-pbehavior-weather-service-3-2",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-3",
      "output_template": "Test service weather 3",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [
        {"name": "test-resource-pbehavior-weather-service-3-1"},
        {"name": "test-resource-pbehavior-weather-service-3-2"}
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-3-1",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-pbehavior-weather-service-3-1"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services?filter={"name":"test-pbehavior-weather-service-3"}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-3",
          "state": {"val": 2},
          "status": {"val": 1},
          "icon": "major",
          "secondary_icon": "maintenance",
          "is_grey": false,
          "alarm_counters": [
            {
              "count": 1,
              "type": {
                "_id": "test-maintenance-type-to-engine",
                "description": "Engine maintenance",
                "icon_name": "test-maintenance-to-engine-icon",
                "name": "Engine maintenance",
                "priority": 18,
                "type": "maintenance"
              }
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

  Scenario: given service with maintenance pbehavior should get maintenance icon and grey flag
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-4",
      "output_template": "Test service weather 4",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [{"name": "test-resource-pbehavior-weather-service-4"}],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-4",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-pbehavior-weather-service-4"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector" :  "test-connector-pbehavior-weather-service-4",
      "connector_name": "test-connector-name-pbehavior-weather-service-4",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-4",
      "resource" : "test-resource-pbehavior-weather-service-4",
      "state" : 3,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services?filter={"name":"test-pbehavior-weather-service-4"}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-4",
          "state": {"val": 3},
          "status": {"val": 1},
          "icon": "maintenance",
          "secondary_icon": "",
          "is_grey": true,
          "alarm_counters": [],
          "pbehavior_info": {
            "canonical_type": "maintenance",
            "name": "test-pbehavior-weather-service-4"
          },
          "pbehaviors": [
            {
              "name": "test-pbehavior-weather-service-4",
              "author": {
                "_id": "root",
                "name": "root"
              },
              "reason": {
                "_id": "test-reason-to-engine",
                "name": "Test Engine",
                "description": "Test Engine"
              },
              "type": {
                "_id": "test-maintenance-type-to-engine",
                "icon_name": "test-maintenance-to-engine-icon",
                "name": "Engine maintenance",
                "type": "maintenance"
              }
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

  Scenario: given service with active pbehavior should get state icon
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-5",
      "output_template": "Test service weather 5",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [{"name": "test-resource-pbehavior-weather-service-5"}],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-5",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-active-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-pbehavior-weather-service-5"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector" :  "test-connector-pbehavior-weather-service-5",
      "connector_name": "test-connector-name-pbehavior-weather-service-5",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-5",
      "resource" : "test-resource-pbehavior-weather-service-5",
      "state" : 3,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services?filter={"name":"test-pbehavior-weather-service-5"}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-5",
          "state": {"val": 3},
          "status": {"val": 1},
          "icon": "critical",
          "secondary_icon": "",
          "is_grey": false,
          "alarm_counters": [],
          "pbehavior_info": {
            "canonical_type": "active",
            "name": "test-pbehavior-weather-service-5"
          },
          "pbehaviors": [
            {
              "name": "test-pbehavior-weather-service-5",
              "author": {
                "_id": "root",
                "name": "root"
              },
              "reason": {
                "_id": "test-reason-to-engine",
                "name": "Test Engine",
                "description": "Test Engine"
              },
              "type": {
                "_id": "test-active-type-to-engine",
                "icon_name": "test-active-to-engine-icon",
                "name": "Engine active",
                "type": "active"
              }
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

  Scenario: given service with maintenance pbehavior and one dependency with maintenance pbehavior
    should get maintenance icon and not get secondary icon
    Given I am admin
    When I send an event:
    """json
    {
      "connector" :  "test-connector-pbehavior-weather-service-6",
      "connector_name": "test-connector-name-pbehavior-weather-service-6",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-6",
      "resource" : "test-resource-pbehavior-weather-service-6",
      "state" : 3,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-6",
      "output_template": "Test service weather 6",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [{"name": "test-resource-pbehavior-weather-service-6"}],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-6-1",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "2h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-pbehavior-weather-service-6"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-6-2",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "2h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-pbehavior-weather-service-6"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services?filter={"name":"test-pbehavior-weather-service-6"}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-6",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "maintenance",
          "secondary_icon": "",
          "is_grey": true,
          "alarm_counters": [
            {
              "count": 1,
              "type": {
                "_id": "test-maintenance-type-to-engine",
                "description": "Engine maintenance",
                "icon_name": "test-maintenance-to-engine-icon",
                "name": "Engine maintenance",
                "priority": 18,
                "type": "maintenance"
              }
            }
          ],
          "pbehavior_info": {
            "canonical_type": "maintenance",
            "name": "test-pbehavior-weather-service-6-1"
          },
          "pbehaviors": [
            {
              "name": "test-pbehavior-weather-service-6-1",
              "author": {
                "_id": "root",
                "name": "root"
              },
              "reason": {
                "_id": "test-reason-to-engine",
                "name": "Test Engine",
                "description": "Test Engine"
              },
              "type": {
                "_id": "test-maintenance-type-to-engine",
                "icon_name": "test-maintenance-to-engine-icon",
                "name": "Engine maintenance",
                "type": "maintenance"
              }
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

  Scenario: given service with maintenance pbehavior and dependency with maintenance pbehavior
    and another dependency without pbehavior should get maintenance icon and maintenance secondary icon
    Given I am admin
    When I send an event:
    """json
    {
      "connector" :  "test-connector-pbehavior-weather-service-7",
      "connector_name": "test-connector-name-pbehavior-weather-service-7",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-7",
      "resource" : "test-resource-pbehavior-weather-service-7-1",
      "state" : 3,
      "output" : "noveo alarm"
    }
    """
    When I send an event:
    """json
    {
      "connector" :  "test-connector-pbehavior-weather-service-7",
      "connector_name": "test-connector-name-pbehavior-weather-service-7",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-7",
      "resource" : "test-resource-pbehavior-weather-service-7-2",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-7",
      "output_template": "Test service weather 7",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [
        {"name": "test-resource-pbehavior-weather-service-7-1"},
        {"name": "test-resource-pbehavior-weather-service-7-2"}
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-7-1",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "2h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-pbehavior-weather-service-7"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-7-2",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "2h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-pbehavior-weather-service-7-1"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services?filter={"name":"test-pbehavior-weather-service-7"}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-7",
          "state": {"val": 2},
          "status": {"val": 1},
          "icon": "maintenance",
          "secondary_icon": "maintenance",
          "is_grey": true,
          "alarm_counters": [
            {
              "count": 1,
              "type": {
                "_id": "test-maintenance-type-to-engine",
                "description": "Engine maintenance",
                "icon_name": "test-maintenance-to-engine-icon",
                "name": "Engine maintenance",
                "priority": 18,
                "type": "maintenance"
              }
            }
          ],
          "pbehavior_info": {
            "canonical_type": "maintenance",
            "name": "test-pbehavior-weather-service-7-1"
          },
          "pbehaviors": [
            {
              "name": "test-pbehavior-weather-service-7-1",
              "author": {
                "_id": "root",
                "name": "root"
              },
              "reason": {
                "_id": "test-reason-to-engine",
                "name": "Test Engine",
                "description": "Test Engine"
              },
              "type": {
                "_id": "test-maintenance-type-to-engine",
                "icon_name": "test-maintenance-to-engine-icon",
                "name": "Engine maintenance",
                "type": "maintenance"
              }
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

  Scenario: given service with maintenance pbehavior should get service by filter icon=maintenance
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-11",
      "output_template": "Test service weather 11",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [{"name": "test-resource-pbehavior-weather-service-11-1"}],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-11",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-pbehavior-weather-service-11"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait the end of event processing
    When I send an event:
    """json
    {
      "connector" :  "test-connector-pbehavior-weather-service-11",
      "connector_name": "test-connector-name-pbehavior-weather-service-11",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-11",
      "resource" : "test-resource-pbehavior-weather-service-11-1",
      "state" : 3,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services?filter={"name":"test-pbehavior-weather-service-11","icon":"maintenance"}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-11"
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

  Scenario: given one dependency with maintenance pbehavior should get service by filter icon=maintenance
    Given I am admin
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-12",
      "connector_name": "test-connector-name-pbehavior-weather-service-12",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-12",
      "resource" : "test-resource-pbehavior-weather-service-12",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-12",
      "output_template": "Test service weather 12",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [
        {"name": "test-resource-pbehavior-weather-service-12"}
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-12-1",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-pbehavior-weather-service-12"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services?filter={"name":"test-pbehavior-weather-service-12","icon":"maintenance"}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-12"
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

  Scenario: given service without pbehavior should not get service by filter icon=maintenance
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-13",
      "output_template": "Test service weather 13",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [{"name": "test-resource-pbehavior-weather-service-13-1"}],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector" :  "test-connector-pbehavior-weather-service-13",
      "connector_name": "test-connector-name-pbehavior-weather-service-13",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-13",
      "resource" : "test-resource-pbehavior-weather-service-13-1",
      "state" : 3,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services?filter={"name":"test-pbehavior-weather-service-13","icon":"maintenance"}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """

  Scenario: given dependency with maintenance pbehavior
    and another dependency without pbehavior should get service by filter secondary_icon=maintenance
    Given I am admin
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-14",
      "connector_name": "test-connector-name-pbehavior-weather-service-14",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-14",
      "resource" : "test-resource-pbehavior-weather-service-14-1",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-14",
      "connector_name": "test-connector-name-pbehavior-weather-service-14",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-14",
      "resource" : "test-resource-pbehavior-weather-service-14-2",
      "state" : 3,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-14",
      "output_template": "Test service weather 14",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [
        {"name": "test-resource-pbehavior-weather-service-14-1"},
        {"name": "test-resource-pbehavior-weather-service-14-2"}
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-14-1",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-pbehavior-weather-service-14-1"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services?filter={"name":"test-pbehavior-weather-service-14","secondary_icon":"maintenance"}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-14"
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

  Scenario: given dependency without pbehavior should not get service by filter secondary_icon=maintenance
    Given I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-pbehavior-weather-service-15",
      "output_template": "Test service weather 15",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [{"name": "test-resource-pbehavior-weather-service-15-1"}],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I send an event:
    """json
    {
      "connector" :  "test-connector-pbehavior-weather-service-15",
      "connector_name": "test-connector-name-pbehavior-weather-service-15",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-15",
      "resource" : "test-resource-pbehavior-weather-service-15-1",
      "state" : 3,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services?filter={"name":"test-pbehavior-weather-service-15","secondary_icon":"maintenance"}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """

  Scenario: given service with maintenance pbehavior without alarm should get maintenance icon
    Given I am admin
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-16",
      "connector_name": "test-connector-name-pbehavior-weather-service-16",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-16",
      "resource" : "test-resource-pbehavior-weather-service-16",
      "state" : 0,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-pbehavior-weather-service-16",
      "name": "test-pbehavior-weather-service-16",
      "output_template": "test-pbehavior-weather-service-16",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [
         {"name": "test-resource-pbehavior-weather-service-16"}
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-16",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "_id": "test-pbehavior-weather-service-16"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait the end of event processing
    When I do GET /api/v4/weather-services?filter={"name":"test-pbehavior-weather-service-16"}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-16",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "maintenance",
          "secondary_icon": "",
          "is_grey": true,
          "alarm_counters": [],
          "pbehavior_info": {
            "canonical_type": "maintenance",
            "name": "test-pbehavior-weather-service-16"
          },
          "pbehaviors": [
            {
              "name": "test-pbehavior-weather-service-16",
              "author": {
                "_id": "root",
                "name": "root"
              },
              "reason": {
                "_id": "test-reason-to-engine",
                "name": "Test Engine",
                "description": "Test Engine"
              },
              "type": {
                "_id": "test-maintenance-type-to-engine",
                "icon_name": "test-maintenance-to-engine-icon",
                "name": "Engine maintenance",
                "type": "maintenance"
              }
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

  Scenario: given dependency with maintenance pbehavior without alarm should get maintenance icon
    Given I am admin
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-17",
      "connector_name": "test-connector-name-pbehavior-weather-service-17",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-17",
      "resource" : "test-resource-pbehavior-weather-service-17",
      "state" : 0,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-pbehavior-weather-service-17",
      "name": "test-pbehavior-weather-service-17",
      "output_template": "test-pbehavior-weather-service-17",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [
         {"name": "test-resource-pbehavior-weather-service-17"}
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-17",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-pbehavior-weather-service-17"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services?filter={"name":"test-pbehavior-weather-service-17"}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-17",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "maintenance",
          "secondary_icon": "",
          "is_grey": true,
          "alarm_counters": [
            {
              "count": 1,
              "type": {
                "_id": "test-maintenance-type-to-engine",
                "description": "Engine maintenance",
                "icon_name": "test-maintenance-to-engine-icon",
                "name": "Engine maintenance",
                "priority": 18,
                "type": "maintenance"
              }
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

  Scenario: given dependencies with maintenance pbehavior should keep maintenance icon on alarm create and resolve
    Given I am admin
    When I do POST /api/v4/resolve-rules:
    """json
    {
      "_id": "test-resolve-rule-pbehavior-weather-service-18",
      "name": "test-resolve-rule-pbehavior-weather-service-18-name",
      "entity_patterns": [
         {"name": "test-resource-pbehavior-weather-service-18-1"},
         {"name": "test-resource-pbehavior-weather-service-18-2"}
      ],
      "duration": {
        "value": 1,
        "unit": "s"
      },
      "priority": 10
    }
    """
    Then the response code should be 201
    When I wait the next periodical process
    When I send an event:
    """json
    [
      {
        "connector" : "test-connector-pbehavior-weather-service-18",
        "connector_name": "test-connector-name-pbehavior-weather-service-18",
        "source_type": "resource",
        "event_type": "check",
        "component" :  "test-component-pbehavior-weather-service-18",
        "resource" : "test-resource-pbehavior-weather-service-18-1",
        "state" : 1,
        "output" : "noveo alarm"
      },
      {
        "connector" : "test-connector-pbehavior-weather-service-18",
        "connector_name": "test-connector-name-pbehavior-weather-service-18",
        "source_type": "resource",
        "event_type": "check",
        "component" :  "test-component-pbehavior-weather-service-18",
        "resource" : "test-resource-pbehavior-weather-service-18-2",
        "state" : 0,
        "output" : "noveo alarm"
      }
    ]
    """
    When I wait the end of 2 events processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-pbehavior-weather-service-18",
      "name": "test-pbehavior-weather-service-18",
      "output_template": "test-pbehavior-weather-service-18",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [
         {"name": "test-resource-pbehavior-weather-service-18-1"},
         {"name": "test-resource-pbehavior-weather-service-18-2"}
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-18",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$or":[
          {
            "name": "test-resource-pbehavior-weather-service-18-1"
          },
          {
            "name": "test-resource-pbehavior-weather-service-18-2"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait the end of 4 events processing
    When I do GET /api/v4/weather-services?filter={"name":"test-pbehavior-weather-service-18"}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-18",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "maintenance",
          "secondary_icon": "",
          "is_grey": true,
          "alarm_counters": [
            {
              "count": 2,
              "type": {
                "_id": "test-maintenance-type-to-engine",
                "description": "Engine maintenance",
                "icon_name": "test-maintenance-to-engine-icon",
                "name": "Engine maintenance",
                "priority": 18,
                "type": "maintenance"
              }
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
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-18",
      "connector_name": "test-connector-name-pbehavior-weather-service-18",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-18",
      "resource" : "test-resource-pbehavior-weather-service-18-2",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services?filter={"name":"test-pbehavior-weather-service-18"}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-18",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "maintenance",
          "secondary_icon": "",
          "is_grey": true,
          "alarm_counters": [
            {
              "count": 2,
              "type": {
                "_id": "test-maintenance-type-to-engine",
                "description": "Engine maintenance",
                "icon_name": "test-maintenance-to-engine-icon",
                "name": "Engine maintenance",
                "priority": 18,
                "type": "maintenance"
              }
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
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-18",
      "connector_name": "test-connector-name-pbehavior-weather-service-18",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-18",
      "resource" : "test-resource-pbehavior-weather-service-18-1",
      "state" : 0,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/weather-services?filter={"name":"test-pbehavior-weather-service-18"}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-18",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "maintenance",
          "secondary_icon": "",
          "is_grey": true,
          "alarm_counters": [
            {
              "count": 2,
              "type": {
                "_id": "test-maintenance-type-to-engine",
                "description": "Engine maintenance",
                "icon_name": "test-maintenance-to-engine-icon",
                "name": "Engine maintenance",
                "priority": 18,
                "type": "maintenance"
              }
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
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-18",
      "connector_name": "test-connector-name-pbehavior-weather-service-18",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-18",
      "resource" : "test-resource-pbehavior-weather-service-18-2",
      "state" : 0,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/weather-services?filter={"name":"test-pbehavior-weather-service-18"}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-18",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "maintenance",
          "secondary_icon": "",
          "is_grey": true,
          "alarm_counters": [
            {
              "count": 2,
              "type": {
                "_id": "test-maintenance-type-to-engine",
                "description": "Engine maintenance",
                "icon_name": "test-maintenance-to-engine-icon",
                "name": "Engine maintenance",
                "priority": 18,
                "type": "maintenance"
              }
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

  Scenario: given dependency with maintenance pbehavior should update service state correctly on alarm state change
    Given I am admin
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-19",
      "connector_name": "test-connector-name-pbehavior-weather-service-19",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-19",
      "resource" : "test-resource-pbehavior-weather-service-19",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-pbehavior-weather-service-19",
      "name": "test-pbehavior-weather-service-19",
      "output_template": "test-pbehavior-weather-service-19",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [
         {"name": "test-resource-pbehavior-weather-service-19"}
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-19",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter": {
        "name": "test-resource-pbehavior-weather-service-19"
      }
    }
    """
    Then the response code should be 201
    When I save response pbhID={{ .lastResponse._id }}
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services?filter={"name":"test-pbehavior-weather-service-19"}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-19",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "maintenance",
          "secondary_icon": "",
          "is_grey": true,
          "alarm_counters": [
            {
              "count": 1,
              "type": {
                "_id": "test-maintenance-type-to-engine",
                "description": "Engine maintenance",
                "icon_name": "test-maintenance-to-engine-icon",
                "name": "Engine maintenance",
                "priority": 18,
                "type": "maintenance"
              }
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
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-19",
      "connector_name": "test-connector-name-pbehavior-weather-service-19",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-19",
      "resource" : "test-resource-pbehavior-weather-service-19",
      "state" : 0,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/weather-services?filter={"name":"test-pbehavior-weather-service-19"}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-19",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "maintenance",
          "secondary_icon": "",
          "is_grey": true,
          "alarm_counters": [
            {
              "count": 1,
              "type": {
                "_id": "test-maintenance-type-to-engine",
                "description": "Engine maintenance",
                "icon_name": "test-maintenance-to-engine-icon",
                "name": "Engine maintenance",
                "priority": 18,
                "type": "maintenance"
              }
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
    When I do DELETE /api/v4/pbehaviors/{{ .pbhID }}
    Then the response code should be 204
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services?filter={"name":"test-pbehavior-weather-service-19"}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-19",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "ok",
          "secondary_icon": "",
          "is_grey": false,
          "alarm_counters": []
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
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-19",
      "connector_name": "test-connector-name-pbehavior-weather-service-19",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-19",
      "resource" : "test-resource-pbehavior-weather-service-19",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services?filter={"name":"test-pbehavior-weather-service-19"}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-19",
          "state": {"val": 2},
          "status": {"val": 1},
          "icon": "major",
          "secondary_icon": "",
          "is_grey": false,
          "alarm_counters": []
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

  Scenario: given entity service entity should decrease pbh filter if disabled and increase if enabled
    Given I am admin
    When I send an event:
    """json
    [
      {
        "connector" : "test-connector-pbehavior-weather-service-20",
        "connector_name": "test-connector-name-pbehavior-weather-service-20",
        "source_type": "resource",
        "event_type": "check",
        "component" :  "test-component-pbehavior-weather-service-20",
        "resource" : "test-resource-pbehavior-weather-service-20-1",
        "state" : 2,
        "output" : "noveo alarm"
      },
      {
        "connector" : "test-connector-pbehavior-weather-service-20",
        "connector_name": "test-connector-name-pbehavior-weather-service-20",
        "source_type": "resource",
        "event_type": "check",
        "component" :  "test-component-pbehavior-weather-service-20",
        "resource" : "test-resource-pbehavior-weather-service-20-2",
        "state" : 2,
        "output" : "noveo alarm"
      },
      {
        "connector" : "test-connector-pbehavior-weather-service-20",
        "connector_name": "test-connector-name-pbehavior-weather-service-20",
        "source_type": "resource",
        "event_type": "check",
        "component" :  "test-component-pbehavior-weather-service-20",
        "resource" : "test-resource-pbehavior-weather-service-20-3",
        "state" : 2,
        "output" : "noveo alarm"
      }
    ]
    """
    When I wait the end of 3 events processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-pbehavior-weather-service-20",
      "name": "test-pbehavior-weather-service-20",
      "output_template": "test-pbehavior-weather-service-20",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [
        {
          "name": {
            "regex_match": "test-resource-pbehavior-weather-service-20"
          }
        }
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-20",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$or":[
          {
            "name": "test-resource-pbehavior-weather-service-20-2"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services?filter={"name":"test-pbehavior-weather-service-20"}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-20",
          "state": {"val": 2},
          "status": {"val": 1},
          "icon": "major",
          "secondary_icon": "maintenance",
          "is_grey": false,
          "alarm_counters": [
            {
              "count": 1,
              "type": {
                "_id": "test-maintenance-type-to-engine",
                "description": "Engine maintenance",
                "icon_name": "test-maintenance-to-engine-icon",
                "name": "Engine maintenance",
                "priority": 18,
                "type": "maintenance"
              }
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
    When I do PUT /api/v4/entitybasics?_id=test-resource-pbehavior-weather-service-20-2/test-component-pbehavior-weather-service-20:
    """json
    {
      "enabled": false,
      "impact_level": 1,
      "sli_avail_state": 0,
      "infos": [],
      "impact": [
        "test-component-pbehavior-weather-service-20"
      ],
      "depends": [
        "test-connector-pbehavior-weather-service-20/test-connector-name-pbehavior-weather-service-20"
      ]
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services?filter={"name":"test-pbehavior-weather-service-20"}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-20",
          "state": {"val": 2},
          "status": {"val": 1},
          "icon": "major",
          "secondary_icon": "",
          "is_grey": false,
          "alarm_counters": []
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
    When I do PUT /api/v4/entitybasics?_id=test-resource-pbehavior-weather-service-20-2/test-component-pbehavior-weather-service-20:
    """json
    {
      "enabled": true,
      "impact_level": 1,
      "sli_avail_state": 0,
      "infos": [],
      "impact": [
        "test-component-pbehavior-weather-service-20"
      ],
      "depends": [
        "test-connector-pbehavior-weather-service-20/test-connector-name-pbehavior-weather-service-20"
      ]
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services?filter={"name":"test-pbehavior-weather-service-20"}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-20",
          "state": {"val": 2},
          "status": {"val": 1},
          "icon": "major",
          "secondary_icon": "maintenance",
          "is_grey": false,
          "alarm_counters": [
            {
              "count": 1,
              "type": {
                "_id": "test-maintenance-type-to-engine",
                "description": "Engine maintenance",
                "icon_name": "test-maintenance-to-engine-icon",
                "name": "Engine maintenance",
                "priority": 18,
                "type": "maintenance"
              }
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
    When I do PUT /api/v4/entitybasics?_id=test-resource-pbehavior-weather-service-20-2/test-component-pbehavior-weather-service-20:
    """json
    {
      "enabled": false,
      "impact_level": 1,
      "sli_avail_state": 0,
      "infos": [],
      "impact": [
        "test-component-pbehavior-weather-service-20"
      ],
      "depends": [
        "test-connector-pbehavior-weather-service-20/test-connector-name-pbehavior-weather-service-20"
      ]
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services?filter={"name":"test-pbehavior-weather-service-20"}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-20",
          "state": {"val": 2},
          "status": {"val": 1},
          "icon": "major",
          "secondary_icon": "",
          "is_grey": false,
          "alarm_counters": []
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
    When I do PUT /api/v4/entitybasics?_id=test-resource-pbehavior-weather-service-20-2/test-component-pbehavior-weather-service-20:
    """json
    {
      "enabled": true,
      "impact_level": 1,
      "sli_avail_state": 0,
      "infos": [],
      "impact": [
        "test-component-pbehavior-weather-service-20"
      ],
      "depends": [
        "test-connector-pbehavior-weather-service-20/test-connector-name-pbehavior-weather-service-20"
      ]
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services?filter={"name":"test-pbehavior-weather-service-20"}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-20",
          "state": {"val": 2},
          "status": {"val": 1},
          "icon": "major",
          "secondary_icon": "maintenance",
          "is_grey": false,
          "alarm_counters": [
            {
              "count": 1,
              "type": {
                "_id": "test-maintenance-type-to-engine",
                "description": "Engine maintenance",
                "icon_name": "test-maintenance-to-engine-icon",
                "name": "Engine maintenance",
                "priority": 18,
                "type": "maintenance"
              }
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

  Scenario: given entity service entity should be grey if active entity is disabled, while other entities are in pbh state,
    if disabled entity is enabled again, the entity service should be returned to the ok state.
    Given I am admin
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-21",
      "connector_name": "test-connector-name-pbehavior-weather-service-21",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-21",
      "resource" : "test-resource-pbehavior-weather-service-21-1",
      "state" : 0,
      "output" : "noveo alarm"
    }
    """
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-21",
      "connector_name": "test-connector-name-pbehavior-weather-service-21",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-21",
      "resource" : "test-resource-pbehavior-weather-service-21-2",
      "state" : 0,
      "output" : "noveo alarm"
    }
    """
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-21",
      "connector_name": "test-connector-name-pbehavior-weather-service-21",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-21",
      "resource" : "test-resource-pbehavior-weather-service-21-3",
      "state" : 0,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 3 events processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-pbehavior-weather-service-21",
      "name": "test-pbehavior-weather-service-21",
      "output_template": "test-pbehavior-weather-service-21",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [
        {
          "name": {
            "regex_match": "test-resource-pbehavior-weather-service-21"
          }
        }
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-21-1",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$or":[
          {
            "name": "test-resource-pbehavior-weather-service-21-1"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-21-3",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$or":[
          {
            "name": "test-resource-pbehavior-weather-service-21-3"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services?filter={"name":"test-pbehavior-weather-service-21"}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-21",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "ok",
          "secondary_icon": "maintenance",
          "is_grey": false,
          "alarm_counters": [
            {
              "count": 2,
              "type": {
                "_id": "test-maintenance-type-to-engine",
                "description": "Engine maintenance",
                "icon_name": "test-maintenance-to-engine-icon",
                "name": "Engine maintenance",
                "priority": 18,
                "type": "maintenance"
              }
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
    When I do PUT /api/v4/entitybasics?_id=test-resource-pbehavior-weather-service-21-2/test-component-pbehavior-weather-service-21:
    """json
    {
      "enabled": false,
      "impact_level": 1,
      "sli_avail_state": 0,
      "infos": [],
      "impact": [
        "test-component-pbehavior-weather-service-21"
      ],
      "depends": [
        "test-connector-pbehavior-weather-service-21/test-connector-name-pbehavior-weather-service-21"
      ]
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/weather-services?filter={"name":"test-pbehavior-weather-service-21"}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-21",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "maintenance",
          "secondary_icon": "",
          "is_grey": true,
          "alarm_counters": [
            {
              "count": 2,
              "type": {
                "_id": "test-maintenance-type-to-engine",
                "description": "Engine maintenance",
                "icon_name": "test-maintenance-to-engine-icon",
                "name": "Engine maintenance",
                "priority": 18,
                "type": "maintenance"
              }
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
    When I do PUT /api/v4/entitybasics?_id=test-resource-pbehavior-weather-service-21-2/test-component-pbehavior-weather-service-21:
    """json
    {
      "enabled": true,
      "impact_level": 1,
      "sli_avail_state": 0,
      "infos": [],
      "impact": [
        "test-component-pbehavior-weather-service-21"
      ],
      "depends": [
        "test-connector-pbehavior-weather-service-21/test-connector-name-pbehavior-weather-service-21"
      ]
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/weather-services?filter={"name":"test-pbehavior-weather-service-21"}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-21",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "ok",
          "secondary_icon": "maintenance",
          "is_grey": false,
          "alarm_counters": [
            {
              "count": 2,
              "type": {
                "_id": "test-maintenance-type-to-engine",
                "description": "Engine maintenance",
                "icon_name": "test-maintenance-to-engine-icon",
                "name": "Engine maintenance",
                "priority": 18,
                "type": "maintenance"
              }
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


  Scenario: given entity service entity should decrease pbh counter on mass disabled and increase on mass enabled actions
    Given I am admin
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-22",
      "connector_name": "test-connector-name-pbehavior-weather-service-22",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-22",
      "resource" : "test-resource-pbehavior-weather-service-22-1",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-22",
      "connector_name": "test-connector-name-pbehavior-weather-service-22",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-22",
      "resource" : "test-resource-pbehavior-weather-service-22-2",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-22",
      "connector_name": "test-connector-name-pbehavior-weather-service-22",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-22",
      "resource" : "test-resource-pbehavior-weather-service-22-3",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 3 events processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-pbehavior-weather-service-22",
      "name": "test-pbehavior-weather-service-22",
      "output_template": "test-pbehavior-weather-service-22",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [
        {
          "name": {
            "regex_match": "test-resource-pbehavior-weather-service-22"
          }
        }
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-22-1",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$or":[
          {
            "name": "test-resource-pbehavior-weather-service-22-1"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-22-2",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$or":[
          {
            "name": "test-resource-pbehavior-weather-service-22-2"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services?filter={"name":"test-pbehavior-weather-service-22"}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-22",
          "state": {"val": 2},
          "status": {"val": 1},
          "icon": "major",
          "secondary_icon": "maintenance",
          "is_grey": false,
          "alarm_counters": [
            {
              "count": 2,
              "type": {
                "_id": "test-maintenance-type-to-engine",
                "description": "Engine maintenance",
                "icon_name": "test-maintenance-to-engine-icon",
                "name": "Engine maintenance",
                "priority": 18,
                "type": "maintenance"
              }
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
    When I do PUT /api/v4/bulk/entities/disable:
    """json
    [
      {
        "_id": "test-resource-pbehavior-weather-service-22-1/test-component-pbehavior-weather-service-22"
      },
      {
        "_id": "test-resource-pbehavior-weather-service-22-2/test-component-pbehavior-weather-service-22"
      }
    ]
    """
    Then the response code should be 207
    When I wait the end of 4 events processing
    When I do GET /api/v4/weather-services?filter={"name":"test-pbehavior-weather-service-22"}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-22",
          "state": {"val": 2},
          "status": {"val": 1},
          "icon": "major",
          "secondary_icon": "",
          "is_grey": false,
          "alarm_counters": []
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
    When I do PUT /api/v4/bulk/entities/enable:
    """json
    [
      {
        "_id": "test-resource-pbehavior-weather-service-22-1/test-component-pbehavior-weather-service-22"
      },
      {
        "_id": "test-resource-pbehavior-weather-service-22-2/test-component-pbehavior-weather-service-22"
      }
    ]
    """
    Then the response code should be 207
    When I wait the end of 4 events processing
    When I do GET /api/v4/weather-services?filter={"name":"test-pbehavior-weather-service-22"}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-22",
          "state": {"val": 2},
          "status": {"val": 1},
          "icon": "major",
          "secondary_icon": "maintenance",
          "is_grey": false,
          "alarm_counters": [
            {
              "count": 2,
              "type": {
                "_id": "test-maintenance-type-to-engine",
                "description": "Engine maintenance",
                "icon_name": "test-maintenance-to-engine-icon",
                "name": "Engine maintenance",
                "priority": 18,
                "type": "maintenance"
              }
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

  Scenario: given entity service entity should be grey if active entities are disabled by bulk disable action, while other entities are in pbh state,
    if disabled entities are enabled again, the entity service should be returned to the ok state.
    Given I am admin
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-23",
      "connector_name": "test-connector-name-pbehavior-weather-service-23",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-23",
      "resource" : "test-resource-pbehavior-weather-service-23-1",
      "state" : 0,
      "output" : "noveo alarm"
    }
    """
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-23",
      "connector_name": "test-connector-name-pbehavior-weather-service-23",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-23",
      "resource" : "test-resource-pbehavior-weather-service-23-2",
      "state" : 0,
      "output" : "noveo alarm"
    }
    """
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-23",
      "connector_name": "test-connector-name-pbehavior-weather-service-23",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-23",
      "resource" : "test-resource-pbehavior-weather-service-23-3",
      "state" : 0,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 3 events processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-pbehavior-weather-service-23",
      "name": "test-pbehavior-weather-service-23",
      "output_template": "test-pbehavior-weather-service-23",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [
        {
          "name": {
            "regex_match": "test-resource-pbehavior-weather-service-23"
          }
        }
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-23-1",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$or":[
          {
            "name": "test-resource-pbehavior-weather-service-23-1"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services?filter={"name":"test-pbehavior-weather-service-23"}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-23",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "ok",
          "secondary_icon": "maintenance",
          "is_grey": false,
          "alarm_counters": [
            {
              "count": 1,
              "type": {
                "_id": "test-maintenance-type-to-engine",
                "description": "Engine maintenance",
                "icon_name": "test-maintenance-to-engine-icon",
                "name": "Engine maintenance",
                "priority": 18,
                "type": "maintenance"
              }
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
    When I do PUT /api/v4/bulk/entities/disable:
    """json
    [
      {
        "_id": "test-resource-pbehavior-weather-service-23-2/test-component-pbehavior-weather-service-23"
      },
      {
        "_id": "test-resource-pbehavior-weather-service-23-3/test-component-pbehavior-weather-service-23"
      }
    ]
    """
    Then the response code should be 207
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services?filter={"name":"test-pbehavior-weather-service-23"}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-23",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "maintenance",
          "secondary_icon": "",
          "is_grey": true,
          "alarm_counters": [
            {
              "count": 1,
              "type": {
                "_id": "test-maintenance-type-to-engine",
                "description": "Engine maintenance",
                "icon_name": "test-maintenance-to-engine-icon",
                "name": "Engine maintenance",
                "priority": 18,
                "type": "maintenance"
              }
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
    When I do PUT /api/v4/bulk/entities/enable:
    """json
    [
      {
        "_id": "test-resource-pbehavior-weather-service-23-2/test-component-pbehavior-weather-service-23"
      },
      {
        "_id": "test-resource-pbehavior-weather-service-23-3/test-component-pbehavior-weather-service-23"
      }
    ]
    """
    Then the response code should be 207
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services?filter={"name":"test-pbehavior-weather-service-23"}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-23",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "ok",
          "secondary_icon": "maintenance",
          "is_grey": false,
          "alarm_counters": [
            {
              "count": 1,
              "type": {
                "_id": "test-maintenance-type-to-engine",
                "description": "Engine maintenance",
                "icon_name": "test-maintenance-to-engine-icon",
                "name": "Engine maintenance",
                "priority": 18,
                "type": "maintenance"
              }
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

  Scenario: given entity service entity should recompute counters properly if entity service is enabled before the dependent entity in mass enable action
    Given I am admin
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-24",
      "connector_name": "test-connector-name-pbehavior-weather-service-24",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-24",
      "resource" : "test-resource-pbehavior-weather-service-24-1",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-24",
      "connector_name": "test-connector-name-pbehavior-weather-service-24",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-24",
      "resource" : "test-resource-pbehavior-weather-service-24-2",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-weather-service-24",
      "connector_name": "test-connector-name-pbehavior-weather-service-24",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-service-24",
      "resource" : "test-resource-pbehavior-weather-service-24-3",
      "state" : 3,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 3 events processing
    When I do POST /api/v4/entityservices:
    """json
    {
      "_id": "test-pbehavior-weather-service-24",
      "name": "test-pbehavior-weather-service-24",
      "output_template": "test-pbehavior-weather-service-24",
      "category": "test-category-pbehavior-weather-service",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [
        {
          "name": {
            "regex_match": "test-resource-pbehavior-weather-service-24"
          }
        }
      ],
      "sli_avail_state": 0
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-weather-service-24-1",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$or":[
          {
            "name": "test-resource-pbehavior-weather-service-24-1"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-services?filter={"name":"test-pbehavior-weather-service-24"}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-24",
          "state": {"val": 3},
          "status": {"val": 1},
          "icon": "critical",
          "secondary_icon": "maintenance",
          "is_grey": false,
          "alarm_counters": [
            {
              "count": 1,
              "type": {
                "_id": "test-maintenance-type-to-engine",
                "description": "Engine maintenance",
                "icon_name": "test-maintenance-to-engine-icon",
                "name": "Engine maintenance",
                "priority": 18,
                "type": "maintenance"
              }
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
    When I do PUT /api/v4/bulk/entities/disable:
    """json
    [
      {
        "_id": "test-pbehavior-weather-service-24"
      },
      {
        "_id": "test-resource-pbehavior-weather-service-24-1/test-component-pbehavior-weather-service-24"
      }
    ]
    """
    Then the response code should be 207
    When I wait the end of 2 events processing
    When I do PUT /api/v4/bulk/entities/enable:
    """json
    [
      {
        "_id": "test-pbehavior-weather-service-24"
      },
      {
        "_id": "test-resource-pbehavior-weather-service-24-1/test-component-pbehavior-weather-service-24"
      }
    ]
    """
    Then the response code should be 207
    When I wait the end of 3 events processing
    When I do GET /api/v4/weather-services?filter={"name":"test-pbehavior-weather-service-24"}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-pbehavior-weather-service-24",
          "state": {"val": 3},
          "status": {"val": 1},
          "icon": "critical",
          "secondary_icon": "maintenance",
          "is_grey": false,
          "alarm_counters": [
            {
              "count": 1,
              "type": {
                "_id": "test-maintenance-type-to-engine",
                "description": "Engine maintenance",
                "icon_name": "test-maintenance-to-engine-icon",
                "name": "Engine maintenance",
                "priority": 18,
                "type": "maintenance"
              }
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

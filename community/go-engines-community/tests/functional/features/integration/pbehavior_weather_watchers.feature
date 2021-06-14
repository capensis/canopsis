Feature: get watcher weather
  I need to be able to get watcher weather

  Scenario: given watcher for one entity with maintenance pbehavior
    should not get secondary icon
    Given I am admin
    When I do POST /api/v2/watcherng:
    """
    {
      "name": "test-pbehavior-weather-watcher-1",
      "type": "watcher",
      "enabled": true,
      "state": {"method":"worst"},
      "entities":[
        {"name": "test-resource-pbehavior-weather-watcher-1"}
      ],
      "output_template": "Test watcher weather 1"
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """
    {
      "enabled": true,
      "author": "root",
      "name": "test-pbehavior-weather-watcher-1",
      "tstart": {{ now.Unix }},
      "tstop": {{ (now.Add (parseDuration "10m")).Unix }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-pbehavior-weather-watcher-1"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait 1s
    When I send an event:
    """
    {
      "connector" : "test-connector-pbehavior-weather-watcher-1",
      "connector_name": "test-connector-name-pbehavior-weather-watcher-1",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-watcher-1",
      "resource" : "test-resource-pbehavior-weather-watcher-1",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-watchers?filter={"name":"test-pbehavior-weather-watcher-1"}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "name": "test-pbehavior-weather-watcher-1",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "maintenance",
          "secondary_icon": "",
          "color": "pause",
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

  Scenario: given watcher for one entity with active pbehavior
    should not get secondary icon
    Given I am admin
    When I do POST /api/v2/watcherng:
    """
    {
      "name": "test-pbehavior-weather-watcher-2",
      "type": "watcher",
      "enabled": true,
      "state": {"method":"worst"},
      "entities":[
        {"name": "test-resource-pbehavior-weather-watcher-2"}
      ],
      "output_template": "Test watcher weather 2"
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """
    {
      "enabled": true,
      "author": "root",
      "name": "test-pbehavior-weather-watcher-2",
      "tstart": {{ now.Unix }},
      "tstop": {{ (now.Add (parseDuration "10m")).Unix }},
      "type": "test-active-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-pbehavior-weather-watcher-2"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait 1s
    When I send an event:
    """
    {
      "connector" : "test-connector-pbehavior-weather-watcher-2",
      "connector_name": "test-connector-name-pbehavior-weather-watcher-2",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-watcher-2",
      "resource" : "test-resource-pbehavior-weather-watcher-2",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-watchers?filter={"name":"test-pbehavior-weather-watcher-2"}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "name": "test-pbehavior-weather-watcher-2",
          "state": {"val": 2},
          "status": {"val": 1},
          "icon": "major",
          "secondary_icon": "",
          "color": "major",
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

  Scenario: given watcher for one entity with maintenance pbehavior
    and another entity without pbehavior should get secondary icon
    Given I am admin
    When I do POST /api/v2/watcherng:
    """
    {
      "name": "test-pbehavior-weather-watcher-3",
      "type": "watcher",
      "enabled": true,
      "state": {"method":"worst"},
      "entities":[
        {"name": "test-resource-pbehavior-weather-watcher-3-1"},
        {"name": "test-resource-pbehavior-weather-watcher-3-2"}
      ],
      "output_template": "Test watcher weather 3"
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """
    {
      "enabled": true,
      "author": "root",
      "name": "test-pbehavior-weather-watcher-3-1",
      "tstart": {{ now.Unix }},
      "tstop": {{ (now.Add (parseDuration "10m")).Unix }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-pbehavior-weather-watcher-3-1"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait 1s
    When I send an event:
    """
    {
      "connector" :  "test-connector-pbehavior-weather-watcher-3",
      "connector_name": "test-connector-name-pbehavior-weather-watcher-3",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-watcher-3",
      "resource" : "test-resource-pbehavior-weather-watcher-3-1",
      "state" : 3,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """
    {
      "connector" : "test-connector-pbehavior-weather-watcher-3",
      "connector_name": "test-connector-name-pbehavior-weather-watcher-3",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-watcher-3",
      "resource" : "test-resource-pbehavior-weather-watcher-3-2",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-watchers?filter={"name":"test-pbehavior-weather-watcher-3"}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "name": "test-pbehavior-weather-watcher-3",
          "state": {"val": 2},
          "status": {"val": 1},
          "icon": "major",
          "secondary_icon": "maintenance",
          "color": "major",
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

  Scenario: given watcher with maintenance pbehavior for one entity
    should get maintenance icon and pause color
    Given I am admin
    When I do POST /api/v2/watcherng:
    """
    {
      "name": "test-pbehavior-weather-watcher-4",
      "type": "watcher",
      "enabled": true,
      "state": {"method":"worst"},
      "entities":[
        {"name": "test-resource-pbehavior-weather-watcher-4"}
      ],
      "output_template": "Test watcher weather 4"
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """
    {
      "enabled": true,
      "author": "root",
      "name": "test-pbehavior-weather-watcher-4",
      "tstart": {{ now.Unix }},
      "tstop": {{ (now.Add (parseDuration "10m")).Unix }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-pbehavior-weather-watcher-4"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait 1s
    When I send an event:
    """
    {
      "connector" :  "test-connector-pbehavior-weather-watcher-4",
      "connector_name": "test-connector-name-pbehavior-weather-watcher-4",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-watcher-4",
      "resource" : "test-resource-pbehavior-weather-watcher-4",
      "state" : 3,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-watchers?filter={"name":"test-pbehavior-weather-watcher-4"}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "name": "test-pbehavior-weather-watcher-4",
          "state": {"val": 3},
          "status": {"val": 1},
          "icon": "maintenance",
          "secondary_icon": "",
          "color": "pause",
          "alarm_counters": [],
          "pbehaviors": [
            {
              "name": "test-pbehavior-weather-watcher-4"
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

  Scenario: given watcher with active pbehavior for one entity
    should get state icon
    Given I am admin
    When I do POST /api/v2/watcherng:
    """
    {
      "name": "test-pbehavior-weather-watcher-5",
      "type": "watcher",
      "enabled": true,
      "state": {"method":"worst"},
      "entities":[
        {"name": "test-resource-pbehavior-weather-watcher-5"}
      ],
      "output_template": "Test watcher weather 5"
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """
    {
      "enabled": true,
      "author": "root",
      "name": "test-pbehavior-weather-watcher-5",
      "tstart": {{ now.Unix }},
      "tstop": {{ (now.Add (parseDuration "10m")).Unix }},
      "type": "test-active-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-pbehavior-weather-watcher-5"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait 1s
    When I send an event:
    """
    {
      "connector" :  "test-connector-pbehavior-weather-watcher-5",
      "connector_name": "test-connector-name-pbehavior-weather-watcher-5",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-watcher-5",
      "resource" : "test-resource-pbehavior-weather-watcher-5",
      "state" : 3,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-watchers?filter={"name":"test-pbehavior-weather-watcher-5"}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "name": "test-pbehavior-weather-watcher-5",
          "state": {"val": 3},
          "status": {"val": 1},
          "icon": "critical",
          "secondary_icon": "",
          "color": "critical",
          "alarm_counters": [],
          "pbehaviors": [
            {
              "name": "test-pbehavior-weather-watcher-5"
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

  Scenario: given watcher with maintenance pbehavior for one entity with maintenance pbehavior
    should get maintenance icon and not get secondary icon
    Given I am admin
    When I do POST /api/v2/watcherng:
    """
    {
      "name": "test-pbehavior-weather-watcher-6",
      "type": "watcher",
      "enabled": true,
      "state": {"method":"worst"},
      "entities":[
        {"name": "test-resource-pbehavior-weather-watcher-6"}
      ],
      "output_template": "Test watcher weather 6"
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """
    {
      "enabled": true,
      "author": "root",
      "name": "test-pbehavior-weather-watcher-6-1",
      "tstart": {{ now.Unix }},
      "tstop": {{ (now.Add (parseDuration "20m")).Unix }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-pbehavior-weather-watcher-6"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/pbehaviors:
    """
    {
      "enabled": true,
      "author": "root",
      "name": "test-pbehavior-weather-watcher-6-2",
      "tstart": {{ now.Unix }},
      "tstop": {{ (now.Add (parseDuration "20m")).Unix }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-pbehavior-weather-watcher-6"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait 1s
    When I send an event:
    """
    {
      "connector" :  "test-connector-pbehavior-weather-watcher-6",
      "connector_name": "test-connector-name-pbehavior-weather-watcher-6",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-watcher-6",
      "resource" : "test-resource-pbehavior-weather-watcher-6",
      "state" : 3,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-watchers?filter={"name":"test-pbehavior-weather-watcher-6"}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "name": "test-pbehavior-weather-watcher-6",
          "state": {"val": 0},
          "status": {"val": 0},
          "icon": "maintenance",
          "secondary_icon": "",
          "color": "pause",
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

  Scenario: given watcher with maintenance pbehavior for one entity with maintenance pbehavior
    and another entity without pbehavior should get maintenance icon and maintenance secondary icon
    Given I am admin
    When I do POST /api/v2/watcherng:
    """
    {
      "name": "test-pbehavior-weather-watcher-7",
      "type": "watcher",
      "enabled": true,
      "state": {"method":"worst"},
      "entities":[
        {"name": "test-resource-pbehavior-weather-watcher-7-1"},
        {"name": "test-resource-pbehavior-weather-watcher-7-2"}
      ],
      "output_template": "Test watcher weather 7"
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """
    {
      "enabled": true,
      "author": "root",
      "name": "test-pbehavior-weather-watcher-7-1",
      "tstart": {{ now.Unix }},
      "tstop": {{ (now.Add (parseDuration "20m")).Unix }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-pbehavior-weather-watcher-7"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/pbehaviors:
    """
    {
      "enabled": true,
      "author": "root",
      "name": "test-pbehavior-weather-watcher-7-2",
      "tstart": {{ now.Unix }},
      "tstop": {{ (now.Add (parseDuration "20m")).Unix }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-pbehavior-weather-watcher-7-1"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait 1s
    When I send an event:
    """
    {
      "connector" :  "test-connector-pbehavior-weather-watcher-7",
      "connector_name": "test-connector-name-pbehavior-weather-watcher-7",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-watcher-7",
      "resource" : "test-resource-pbehavior-weather-watcher-7-1",
      "state" : 3,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """
    {
      "connector" :  "test-connector-pbehavior-weather-watcher-7",
      "connector_name": "test-connector-name-pbehavior-weather-watcher-7",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-watcher-7",
      "resource" : "test-resource-pbehavior-weather-watcher-7-2",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-watchers?filter={"name":"test-pbehavior-weather-watcher-7"}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "name": "test-pbehavior-weather-watcher-7",
          "state": {"val": 2},
          "status": {"val": 1},
          "icon": "maintenance",
          "secondary_icon": "maintenance",
          "color": "pause",
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
          "pbehaviors": [
            {
              "name": "test-pbehavior-weather-watcher-7-1"
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

  Scenario: given watcher with maintenance pbehavior should get watcher by filter color=pause
    Given I am admin
    When I do POST /api/v2/watcherng:
    """
    {
      "name": "test-pbehavior-weather-watcher-8",
      "type": "watcher",
      "enabled": true,
      "state": {"method":"worst"},
      "entities":[
        {"name": "test-resource-pbehavior-weather-watcher-8-1"}
      ],
      "output_template": "Test watcher weather 8"
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """
    {
      "enabled": true,
      "author": "root",
      "name": "test-pbehavior-weather-watcher-8",
      "tstart": {{ now.Unix }},
      "tstop": {{ (now.Add (parseDuration "10m")).Unix }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-pbehavior-weather-watcher-8"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait 1s
    When I send an event:
    """
    {
      "connector" :  "test-connector-pbehavior-weather-watcher-8",
      "connector_name": "test-connector-name-pbehavior-weather-watcher-8",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-watcher-8",
      "resource" : "test-resource-pbehavior-weather-watcher-8-1",
      "state" : 3,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-watchers?filter={"name":"test-pbehavior-weather-watcher-8","color":"pause"}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "name": "test-pbehavior-weather-watcher-8",
          "icon": "maintenance",
          "secondary_icon": "",
          "color": "pause",
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

  Scenario: given watcher for two entities with maintenance pbehavior
  should get watcher by filter color=pause
    Given I am admin
    When I do POST /api/v2/watcherng:
    """
    {
      "name": "test-pbehavior-weather-watcher-9",
      "type": "watcher",
      "enabled": true,
      "state": {"method":"worst"},
      "entities":[
        {"name": "test-resource-pbehavior-weather-watcher-9-1"},
        {"name": "test-resource-pbehavior-weather-watcher-9-2"}
      ],
      "output_template": "Test watcher weather 9"
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """
    {
      "enabled": true,
      "author": "root",
      "name": "test-pbehavior-weather-watcher-9-1",
      "tstart": {{ now.Unix }},
      "tstop": {{ (now.Add (parseDuration "10m")).Unix }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-pbehavior-weather-watcher-9-1"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/pbehaviors:
    """
    {
      "enabled": true,
      "author": "root",
      "name": "test-pbehavior-weather-watcher-9-2",
      "tstart": {{ now.Unix }},
      "tstop": {{ (now.Add (parseDuration "10m")).Unix }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-pbehavior-weather-watcher-9-2"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait 1s
    When I send an event:
    """
    {
      "connector" :  "test-connector-pbehavior-weather-watcher-9",
      "connector_name": "test-connector-name-pbehavior-weather-watcher-9",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-watcher-9",
      "resource" : "test-resource-pbehavior-weather-watcher-9-1",
      "state" : 3,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """
    {
      "connector" : "test-connector-pbehavior-weather-watcher-9",
      "connector_name": "test-connector-name-pbehavior-weather-watcher-9",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-watcher-9",
      "resource" : "test-resource-pbehavior-weather-watcher-9-2",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-watchers?filter={"name":"test-pbehavior-weather-watcher-9","color":"pause"}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "name": "test-pbehavior-weather-watcher-9",
          "icon": "maintenance",
          "secondary_icon": "",
          "color": "pause",
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

  Scenario: given watcher without pbehavior should not get watcher by filter color=pause
    Given I am admin
    When I do POST /api/v2/watcherng:
    """
    {
      "name": "test-pbehavior-weather-watcher-10",
      "type": "watcher",
      "enabled": true,
      "state": {"method":"worst"},
      "entities":[
        {"name": "test-resource-pbehavior-weather-watcher-10-1"}
      ],
      "output_template": "Test watcher weather 10"
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I send an event:
    """
    {
      "connector" :  "test-connector-pbehavior-weather-watcher-10",
      "connector_name": "test-connector-name-pbehavior-weather-watcher-10",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-watcher-10",
      "resource" : "test-resource-pbehavior-weather-watcher-10-1",
      "state" : 3,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-watchers?filter={"name":"test-pbehavior-weather-watcher-10","color":"pause"}
    Then the response code should be 200
    Then the response body should contain:
    """
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

  Scenario: given watcher with maintenance pbehavior should get watcher by filter icon=maintenance
    Given I am admin
    When I do POST /api/v2/watcherng:
    """
    {
      "name": "test-pbehavior-weather-watcher-11",
      "type": "watcher",
      "enabled": true,
      "state": {"method":"worst"},
      "entities":[
        {"name": "test-resource-pbehavior-weather-watcher-11-1"}
      ],
      "output_template": "Test watcher weather 11"
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """
    {
      "enabled": true,
      "author": "root",
      "name": "test-pbehavior-weather-watcher-11",
      "tstart": {{ now.Unix }},
      "tstop": {{ (now.Add (parseDuration "10m")).Unix }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-pbehavior-weather-watcher-11"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait 1s
    When I send an event:
    """
    {
      "connector" :  "test-connector-pbehavior-weather-watcher-11",
      "connector_name": "test-connector-name-pbehavior-weather-watcher-11",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-watcher-11",
      "resource" : "test-resource-pbehavior-weather-watcher-11-1",
      "state" : 3,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-watchers?filter={"name":"test-pbehavior-weather-watcher-11","icon":"maintenance"}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "name": "test-pbehavior-weather-watcher-11",
          "icon": "maintenance",
          "secondary_icon": "",
          "color": "pause",
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

  Scenario: given watcher for two entities with maintenance pbehavior
  should get watcher by filter icon=maintenance
    Given I am admin
    When I do POST /api/v2/watcherng:
    """
    {
      "name": "test-pbehavior-weather-watcher-12",
      "type": "watcher",
      "enabled": true,
      "state": {"method":"worst"},
      "entities":[
        {"name": "test-resource-pbehavior-weather-watcher-12-1"},
        {"name": "test-resource-pbehavior-weather-watcher-12-2"}
      ],
      "output_template": "Test watcher weather 12"
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """
    {
      "enabled": true,
      "author": "root",
      "name": "test-pbehavior-weather-watcher-12-1",
      "tstart": {{ now.Unix }},
      "tstop": {{ (now.Add (parseDuration "10m")).Unix }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-pbehavior-weather-watcher-12-1"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait 1s
    When I send an event:
    """
    {
      "connector" :  "test-connector-pbehavior-weather-watcher-12",
      "connector_name": "test-connector-name-pbehavior-weather-watcher-12",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-watcher-12",
      "resource" : "test-resource-pbehavior-weather-watcher-12-1",
      "state" : 3,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """
    {
      "enabled": true,
      "author": "root",
      "name": "test-pbehavior-weather-watcher-12-2",
      "tstart": {{ now.Unix }},
      "tstop": {{ (now.Add (parseDuration "10m")).Unix }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-pbehavior-weather-watcher-12-2"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait 1s
    When I send an event:
    """
    {
      "connector" : "test-connector-pbehavior-weather-watcher-12",
      "connector_name": "test-connector-name-pbehavior-weather-watcher-12",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-watcher-12",
      "resource" : "test-resource-pbehavior-weather-watcher-12-2",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-watchers?filter={"name":"test-pbehavior-weather-watcher-12","icon":"maintenance"}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "name": "test-pbehavior-weather-watcher-12",
          "icon": "maintenance",
          "secondary_icon": "",
          "color": "pause",
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

  Scenario: given watcher without pbehavior should not get watcher by filter icon=maintenance
    Given I am admin
    When I do POST /api/v2/watcherng:
    """
    {
      "name": "test-pbehavior-weather-watcher-13",
      "type": "watcher",
      "enabled": true,
      "state": {"method":"worst"},
      "entities":[
        {"name": "test-resource-pbehavior-weather-watcher-13-1"}
      ],
      "output_template": "Test watcher weather 13"
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I send an event:
    """
    {
      "connector" :  "test-connector-pbehavior-weather-watcher-13",
      "connector_name": "test-connector-name-pbehavior-weather-watcher-13",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-watcher-13",
      "resource" : "test-resource-pbehavior-weather-watcher-13-1",
      "state" : 3,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-watchers?filter={"name":"test-pbehavior-weather-watcher-13","icon":"maintenance"}
    Then the response code should be 200
    Then the response body should contain:
    """
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

  Scenario: given watcher for one entity with maintenance pbehavior
    and another entity without pbehavior should get watcher by filter secondary_icon=maintenance
    Given I am admin
    When I do POST /api/v2/watcherng:
    """
    {
      "name": "test-pbehavior-weather-watcher-14",
      "type": "watcher",
      "enabled": true,
      "state": {"method":"worst"},
      "entities":[
        {"name": "test-resource-pbehavior-weather-watcher-14-1"},
        {"name": "test-resource-pbehavior-weather-watcher-14-2"}
      ],
      "output_template": "Test watcher weather 14"
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """
    {
      "enabled": true,
      "author": "root",
      "name": "test-pbehavior-weather-watcher-14-1",
      "tstart": {{ now.Unix }},
      "tstop": {{ (now.Add (parseDuration "10m")).Unix }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-pbehavior-weather-watcher-14-1"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait 1s
    When I send an event:
    """
    {
      "connector" :  "test-connector-pbehavior-weather-watcher-14",
      "connector_name": "test-connector-name-pbehavior-weather-watcher-14",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-watcher-14",
      "resource" : "test-resource-pbehavior-weather-watcher-14-1",
      "state" : 3,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """
    {
      "connector" : "test-connector-pbehavior-weather-watcher-14",
      "connector_name": "test-connector-name-pbehavior-weather-watcher-14",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-watcher-14",
      "resource" : "test-resource-pbehavior-weather-watcher-14-2",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-watchers?filter={"name":"test-pbehavior-weather-watcher-14","secondary_icon":"maintenance"}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "name": "test-pbehavior-weather-watcher-14",
          "icon": "major",
          "secondary_icon": "maintenance",
          "color": "major",
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

  Scenario: given watcher for one entity without pbehavior should not get watcher by filter secondary_icon=maintenance
    Given I am admin
    When I do POST /api/v2/watcherng:
    """
    {
      "name": "test-pbehavior-weather-watcher-15",
      "type": "watcher",
      "enabled": true,
      "state": {"method":"worst"},
      "entities":[
        {"name": "test-resource-pbehavior-weather-watcher-15-1"}
      ],
      "output_template": "Test watcher weather 15"
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I send an event:
    """
    {
      "connector" :  "test-connector-pbehavior-weather-watcher-15",
      "connector_name": "test-connector-name-pbehavior-weather-watcher-15",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-watcher-15",
      "resource" : "test-resource-pbehavior-weather-watcher-15-1",
      "state" : 3,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-watchers?filter={"name":"test-pbehavior-weather-watcher-15","secondary_icon":"maintenance"}
    Then the response code should be 200
    Then the response body should contain:
    """
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

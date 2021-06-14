Feature: get watcher entities
  I need to be able to get watcher entities

  Scenario: given watcher for one entity with maintenance pbehavior
    should get one entity
    Given I am admin
    When I do POST /api/v2/watcherng:
    """
    {
      "_id": "test-pbehavior-weather-watcher-entity-1",
      "name": "test-pbehavior-weather-watcher-entity-1",
      "type": "watcher",
      "enabled": true,
      "state": {"method":"worst"},
      "entities":[
        {"name": "test-resource-pbehavior-weather-watcher-entity-1"}
      ],
      "output_template": "Test watcher weather entity 1"
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """
    {
      "enabled": true,
      "author": "root",
      "name": "test-pbehavior-weather-watcher-entity-1",
      "tstart": {{ now.Unix }},
      "tstop": {{ (now.Add (parseDuration "10m")).Unix }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-pbehavior-weather-watcher-entity-1"
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
      "connector" : "test-connector-pbehavior-weather-watcher-entity-1",
      "connector_name": "test-connector-name-pbehavior-weather-watcher-entity-1",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-watcher-entity-1",
      "resource" : "test-resource-pbehavior-weather-watcher-entity-1",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-watchers/test-pbehavior-weather-watcher-entity-1
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "name": "test-resource-pbehavior-weather-watcher-entity-1",
          "state": {"val": 2},
          "status": {"val": 1},
          "color": "pause",
          "icon": "maintenance",
          "pbehaviors": [
            {
              "name": "test-pbehavior-weather-watcher-entity-1"
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
    should get one entity
    Given I am admin
    When I do POST /api/v2/watcherng:
    """
    {
      "_id": "test-pbehavior-weather-watcher-entity-2",
      "name": "test-pbehavior-weather-watcher-entity-2",
      "type": "watcher",
      "enabled": true,
      "state": {"method":"worst"},
      "entities":[
        {"name": "test-resource-pbehavior-weather-watcher-entity-2"}
      ],
      "output_template": "Test watcher weather entity 2"
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """
    {
      "enabled": true,
      "author": "root",
      "name": "test-pbehavior-weather-watcher-entity-2",
      "tstart": {{ now.Unix }},
      "tstop": {{ (now.Add (parseDuration "10m")).Unix }},
      "type": "test-active-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-pbehavior-weather-watcher-entity-2"
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
      "connector" : "test-connector-pbehavior-weather-watcher-entity-2",
      "connector_name": "test-connector-name-pbehavior-weather-watcher-entity-2",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-watcher-entity-2",
      "resource" : "test-resource-pbehavior-weather-watcher-entity-2",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-watchers/test-pbehavior-weather-watcher-entity-2
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "name": "test-resource-pbehavior-weather-watcher-entity-2",
          "state": {"val": 2},
          "status": {"val": 1},
          "color": "major",
          "icon": "major",
          "pbehaviors": [
            {
              "name": "test-pbehavior-weather-watcher-entity-2"
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

  Scenario: given watcher for one entity with maintenance pbehavior
    and another entity without pbehavior should get multiple entities
    Given I am admin
    When I do POST /api/v2/watcherng:
    """
    {
      "_id": "test-pbehavior-weather-watcher-entity-3",
      "name": "test-pbehavior-weather-watcher-entity-3",
      "type": "watcher",
      "enabled": true,
      "state": {"method":"worst"},
      "entities":[
        {"name": "test-resource-pbehavior-weather-watcher-entity-3-1"},
        {"name": "test-resource-pbehavior-weather-watcher-entity-3-2"}
      ],
      "output_template": "Test watcher weather entity 3"
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """
    {
      "enabled": true,
      "author": "root",
      "name": "test-pbehavior-weather-watcher-entity-3-1",
      "tstart": {{ now.Unix }},
      "tstop": {{ (now.Add (parseDuration "10m")).Unix }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-pbehavior-weather-watcher-entity-3-1"
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
      "connector" :  "test-connector-pbehavior-weather-watcher-entity-3",
      "connector_name": "test-connector-name-pbehavior-weather-watcher-entity-3",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-watcher-entity-3",
      "resource" : "test-resource-pbehavior-weather-watcher-entity-3-1",
      "state" : 3,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I send an event:
    """
    {
      "connector" : "test-connector-pbehavior-weather-watcher-entity-3",
      "connector_name": "test-connector-name-pbehavior-weather-watcher-entity-3",
      "source_type": "resource",
      "event_type": "check",
      "component" :  "test-component-pbehavior-weather-watcher-entity-3",
      "resource" : "test-resource-pbehavior-weather-watcher-entity-3-2",
      "state" : 2,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/weather-watchers/test-pbehavior-weather-watcher-entity-3
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "name": "test-resource-pbehavior-weather-watcher-entity-3-1",
          "state": {"val": 3},
          "status": {"val": 1},
          "color": "pause",
          "icon": "maintenance",
          "pbehaviors": [
            {
              "name": "test-pbehavior-weather-watcher-entity-3-1"
            }
          ]
        },
        {
          "name": "test-resource-pbehavior-weather-watcher-entity-3-2",
          "state": {"val": 2},
          "status": {"val": 1},
          "color": "major",
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
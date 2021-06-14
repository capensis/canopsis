Feature: no update watcher when entity is inactive
  I need to be able to not update watcher when pause or maintenance pbehavior is in action.

  Scenario: given watcher and maintenance pbehavior should not update watcher
    Given I am admin
    When I do POST /api/v2/watcherng:
    """
    {
      "name": "test-pbehavior-watcher-1",
      "type": "watcher",
      "enabled": true,
      "state": {"method":"worst"},
      "entities":[
        {"name": "test-resource-pbehavior-watcher-1"}
      ],
      "output_template": "Test watcher"
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """
    {
      "enabled": true,
      "author": "root",
      "name": "test-pbehavior-watcher-1",
      "tstart": {{ now.Unix }},
      "tstop": {{ (now.Add (parseDuration "10m")).Unix }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-pbehavior-watcher-1"
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
        "connector" : "test-connector-pbehavior-watcher-1",
        "connector_name" : "test-connector-name-pbehavior-watcher-1",
        "source_type" : "resource",
        "event_type" : "check",
        "component" : "test-component-pbehavior-watcher-1",
        "resource" : "test-resource-pbehavior-watcher-1",
        "state" : 1,
        "output" : "noveo alarm"
      }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.resource":"test-pbehavior-watcher-1"}]}
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

  Scenario: given watcher and active pbehavior should update watcher
    Given I am admin
    When I do POST /api/v2/watcherng:
    """
    {
      "_id": "test-pbehavior-watcher-2",
      "name": "test-pbehavior-watcher-2",
      "type": "watcher",
      "enabled": true,
      "state": {"method":"worst"},
      "entities":[
        {"name": "test-resource-pbehavior-watcher-2"}
      ],
      "output_template": "Test watcher"
    }
    """
    Then the response code should be 200
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """
    {
      "enabled": true,
      "author": "root",
      "name": "test-pbehavior-watcher-2",
      "tstart": {{ now.Unix }},
      "tstop": {{ (now.Add (parseDuration "10m")).Unix }},
      "type": "test-active-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-pbehavior-watcher-2"
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
        "connector" : "test-connector-pbehavior-watcher-2",
        "connector_name" : "test-connector-name-pbehavior-watcher-2",
        "source_type" : "resource",
        "event_type" : "check",
        "component" : "test-component-pbehavior-watcher-2",
        "resource" : "test-resource-pbehavior-watcher-2",
        "state" : 2,
        "output" : "noveo alarm"
      }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.component":"test-pbehavior-watcher-2"}]}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "v": {
            "state": {
              "val": 2
            },
            "status": {
              "val": 1
            },
            "connector" : "watcher",
            "connector_name" : "watcher",
            "component" :  "test-pbehavior-watcher-2"
          }
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

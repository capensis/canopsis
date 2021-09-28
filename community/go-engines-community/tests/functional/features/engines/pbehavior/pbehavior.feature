Feature: update alarm on pbehavior
  I need to be able to create pbehavior for alarm

  Scenario: given pbehavior should create alarm with pbehavior info
    Given I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-1",
      "tstart": {{ now.Unix }},
      "tstop": {{ (now.Add (parseDuration "10m")).Unix }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-pbehavior-1"
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
      "connector" : "test-connector-pbehavior-1",
      "connector_name" : "test-connector-name-pbehavior-1",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-pbehavior-1",
      "resource" : "test-resource-pbehavior-1",
      "state" : 1,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.resource":"test-resource-pbehavior-1"}]}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "pbehavior_info": {
              "name": "test-pbehavior-1"
            },
            "connector" : "test-connector-pbehavior-1",
            "connector_name" : "test-connector-name-pbehavior-1",
            "component" : "test-component-pbehavior-1",
            "resource" : "test-resource-pbehavior-1"
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

  Scenario: given pbehavior and alarm should update alarm pbehavior info
    Given I am admin
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-2",
      "connector_name" : "test-connector-name-pbehavior-2",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-pbehavior-2",
      "resource" : "test-resource-pbehavior-2",
      "state" : 1,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-2",
      "tstart": {{ now.Unix }},
      "tstop": {{ (now.Add (parseDuration "10m")).Unix }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-pbehavior-2"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait the end of event processing
    When I do GET /api/v4/alarms?filter={"$and":[{"v.resource":"test-resource-pbehavior-2"}]}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "pbehavior_info": {
              "name": "test-pbehavior-2"
            },
            "connector" : "test-connector-pbehavior-2",
            "connector_name" : "test-connector-name-pbehavior-2",
            "component" : "test-component-pbehavior-2",
            "resource" : "test-resource-pbehavior-2"
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

  Scenario: given pbehavior should update last alarm date of pbehavior
    Given I am admin
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-3",
      "connector_name" : "test-connector-name-pbehavior-3",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-pbehavior-3",
      "resource" : "test-resource-pbehavior-3",
      "state" : 1,
      "output" : "test-output-pbehavior-3"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-3",
      "tstart": {{ now.Unix }},
      "tstop": {{ (now.Add (parseDuration "10m")).Unix }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-pbehavior-3"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "last_alarm_date": null
    }
    """
    When I save response pbehaviorID={{ .lastResponse._id }}
    When I wait the end of event processing
    When I do GET /api/v4/pbehaviors/{{ .pbehaviorID }}
    Then the response code should be 200
    Then the response key "last_alarm_date" should not be "null"

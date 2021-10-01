Feature: get pbehavior
  I need to be able to get pbehavior

  Scenario: given pbehavior should return true status
    Given I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-api-1",
      "tstart": {{ now.UTC.Unix }},
      "tstop": {{ (now.UTC.Add (parseDuration "10m")).Unix }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-pbehavior-api-1"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait 1s
    When I do GET /api/v4/pbehaviors?search=test-pbehavior-api-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_active_status": true
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

  Scenario: given pbehavior should return true status
    Given I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-api-2",
      "tstart": {{ (now.Add (parseDuration "-24h")).Unix }},
      "tstop": {{ (now.Add (parseDuration "-23h50m50s")).Unix }},
      "rrule": "FREQ=DAILY",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-pbehavior-api-2"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait 1s
    When I do GET /api/v4/pbehaviors?search=test-pbehavior-api-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_active_status": true
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

  Scenario: given pbehavior should return false status
    Given I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": false,
      "name": "test-pbehavior-api-3",
      "tstart": {{ now.UTC.Unix }},
      "tstop": {{ (now.UTC.Add (parseDuration "10m")).Unix }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-pbehavior-api-3"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait 1s
    When I do GET /api/v4/pbehaviors?search=test-pbehavior-api-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "is_active_status": false
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

  Scenario: given pbehavior should return status
    Given I am admin
    When I send an event:
    """json
    {
      "connector" : "test-connector-pbehavior-api-4",
      "connector_name" : "test-connector-name-pbehavior-api-4",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-pbehavior-api-4",
      "resource" : "test-resource-pbehavior-api-4",
      "state" : 0,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-api-4-1",
      "tstart": {{ now.UTC.Unix }},
      "tstop": {{ (now.UTC.Add (parseDuration "10m")).Unix }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-pbehavior-api-4"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I wait 1s
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-api-4-2",
      "tstart": {{ (now.UTC.Add (parseDuration "10m")).Unix }},
      "tstop": {{ (now.UTC.Add (parseDuration "20m")).Unix }},
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "filter":{
        "$and":[
          {
            "name": "test-resource-pbehavior-api-4"
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
      "connector" : "test-connector-pbehavior-api-4",
      "connector_name" : "test-connector-name-pbehavior-api-4",
      "source_type" : "resource",
      "event_type" : "check",
      "component" : "test-component-pbehavior-api-4",
      "resource" : "test-resource-pbehavior-api-4",
      "state" : 1,
      "output" : "noveo alarm"
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/entities/pbehaviors?id=test-resource-pbehavior-api-4/test-component-pbehavior-api-4
    Then the response code should be 200
    Then the response body should contain:
    """
    [
      {
        "name": "test-pbehavior-api-4-1",
        "is_active_status": true
      },
      {
        "name": "test-pbehavior-api-4-2",
        "is_active_status": false
      }
    ]
    """

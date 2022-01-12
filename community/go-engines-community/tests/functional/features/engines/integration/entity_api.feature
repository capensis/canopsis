Feature: test dynamic entity api fields
  Scenario: should return entity pbh
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-entity-api-1",
      "connector_name": "test-connector-name-entity-api-1",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-entity-api-1",
      "resource": "test-resource-entity-api-1",
      "state": 0,
      "output": "test-output-entity-api-1"
    }
    """
    When I wait the end of event processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-entity-api-1",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "5m" }},
      "type": "test-default-maintenance-type",
      "reason": "test-reason-1",
      "filter":{
        "$and":[
          {
            "name": "test-resource-entity-api-1"
          }
        ]
      }
    }
    """
    When I wait the end of event processing
    When I do GET /api/v4/entities?search=test-resource-entity-api-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "pbehavior_info": {
            "name": "test-pbehavior-entity-api-1",
            "reason": "test-reason-1-name",
            "type": "test-default-maintenance-type"
          }
        }
      ],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 1
      }
    }
    """

  Scenario: shouldn return entity last event date
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-entity-api-2",
      "connector_name": "test-connector-name-entity-api-2",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-entity-api-2",
      "resource": "test-resource-entity-api-2",
      "state": 0,
      "output": "test-output-entity-api-2",
      "timestamp": {{ now }}
    }
    """
    When I save response checkEventTimestamp={{ (index .lastResponse.sent_events 0).timestamp }}
    When I wait the end of event processing
    When I do GET /api/v4/entities?search=test-resource-entity-api-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "last_event_date": {{ .checkEventTimestamp }}
        }
      ],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 1
      }
    }
    """    

  Scenario: should cound ko and ok events
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-entity-api-3",
      "connector_name": "test-connector-name-entity-api-3",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-entity-api-3",
      "resource": "test-resource-entity-api-3",
      "state": 0,
      "output": "test-output-entity-api-3"
    }
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-entity-api-3",
      "connector_name": "test-connector-name-entity-api-3",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-entity-api-3",
      "resource": "test-resource-entity-api-3",
      "state": 0,
      "output": "test-output-entity-api-3"
    }
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-entity-api-3",
      "connector_name": "test-connector-name-entity-api-3",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-entity-api-3",
      "resource": "test-resource-entity-api-3",
      "state": 0,
      "output": "test-output-entity-api-3"
    }
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-entity-api-3",
      "connector_name": "test-connector-name-entity-api-3",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-entity-api-3",
      "resource": "test-resource-entity-api-3",
      "state": 0,
      "output": "test-output-entity-api-3"
    }
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-entity-api-3",
      "connector_name": "test-connector-name-entity-api-3",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-entity-api-3",
      "resource": "test-resource-entity-api-3",
      "state": 1,
      "output": "test-output-entity-api-3"
    }
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-entity-api-3",
      "connector_name": "test-connector-name-entity-api-3",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-entity-api-3",
      "resource": "test-resource-entity-api-3",
      "state": 2,
      "output": "test-output-entity-api-3"
    }
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-entity-api-3",
      "connector_name": "test-connector-name-entity-api-3",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-entity-api-3",
      "resource": "test-resource-entity-api-3",
      "state": 3,
      "output": "test-output-entity-api-3"
    }
    """
    When I wait the end of 7 events processing
    When I do GET /api/v4/entities?search=test-resource-entity-api-3
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "ok_events": 4,
          "ko_events": 3
        }
      ],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 1
      }
    }
    """

  Scenario: shouldn't count statistic if an entity is in inactive pbh state
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-entity-api-4",
      "connector_name": "test-connector-name-entity-api-4",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-entity-api-4",
      "resource": "test-resource-entity-api-4",
      "state": 0,
      "output": "test-output-entity-api-4"
    }
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-entity-api-4",
      "connector_name": "test-connector-name-entity-api-4",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-entity-api-4",
      "resource": "test-resource-entity-api-4",
      "state": 0,
      "output": "test-output-entity-api-4"
    }
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-entity-api-4",
      "connector_name": "test-connector-name-entity-api-4",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-entity-api-4",
      "resource": "test-resource-entity-api-4",
      "state": 0,
      "output": "test-output-entity-api-4"
    }
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/entities?search=test-resource-entity-api-4
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "ok_events": 3,
          "ko_events": 0
        }
      ],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 1
      }
    }
    """
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-entity-api-4",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "5m" }},
      "type": "test-default-inactive-type",
      "reason": "test-reason-1",
      "filter":{
        "$and":[
          {
            "name": "test-resource-entity-api-4"
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
      "connector": "test-connector-entity-api-4",
      "connector_name": "test-connector-name-entity-api-4",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-entity-api-4",
      "resource": "test-resource-entity-api-4",
      "state": 0,
      "output": "test-output-entity-api-4"
    }
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-entity-api-4",
      "connector_name": "test-connector-name-entity-api-4",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-entity-api-4",
      "resource": "test-resource-entity-api-4",
      "state": 0,
      "output": "test-output-entity-api-4"
    }
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-entity-api-4",
      "connector_name": "test-connector-name-entity-api-4",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-entity-api-4",
      "resource": "test-resource-entity-api-4",
      "state": 1,
      "output": "test-output-entity-api-4"
    }
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/entities?search=test-resource-entity-api-4
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "ok_events": 3,
          "ko_events": 0
        }
      ],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 1
      }
    }
    """

  Scenario: shouldn count statistic if entity in not inactive pbh state
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-entity-api-5",
      "connector_name": "test-connector-name-entity-api-5",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-entity-api-5",
      "resource": "test-resource-entity-api-5",
      "state": 0,
      "output": "test-output-entity-api-5"
    }
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-entity-api-5",
      "connector_name": "test-connector-name-entity-api-5",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-entity-api-5",
      "resource": "test-resource-entity-api-5",
      "state": 0,
      "output": "test-output-entity-api-5"
    }
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-entity-api-5",
      "connector_name": "test-connector-name-entity-api-5",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-entity-api-5",
      "resource": "test-resource-entity-api-5",
      "state": 0,
      "output": "test-output-entity-api-5"
    }
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/entities?search=test-resource-entity-api-5
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "ok_events": 3,
          "ko_events": 0
        }
      ],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 1
      }
    }
    """
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-entity-api-5",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "5m" }},
      "type": "test-default-maintenance-type",
      "reason": "test-reason-1",
      "filter":{
        "$and":[
          {
            "name": "test-resource-entity-api-5"
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
      "connector": "test-connector-entity-api-5",
      "connector_name": "test-connector-name-entity-api-5",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-entity-api-5",
      "resource": "test-resource-entity-api-5",
      "state": 0,
      "output": "test-output-entity-api-5"
    }
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-entity-api-5",
      "connector_name": "test-connector-name-entity-api-5",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-entity-api-5",
      "resource": "test-resource-entity-api-5",
      "state": 0,
      "output": "test-output-entity-api-5"
    }
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-entity-api-5",
      "connector_name": "test-connector-name-entity-api-5",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-entity-api-5",
      "resource": "test-resource-entity-api-5",
      "state": 1,
      "output": "test-output-entity-api-5"
    }
    """
    When I wait the end of 3 events processing
    When I do GET /api/v4/entities?search=test-resource-entity-api-5
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "ok_events": 5,
          "ko_events": 1
        }
      ],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 1
      }
    }
    """

  Scenario: should return corresponding alarm's state
    When I am admin
    When I send an event:
    """json
    {
      "connector": "test-connector-entity-api-6",
      "connector_name": "test-connector-name-entity-api-6",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-entity-api-6",
      "resource": "test-resource-entity-api-6",
      "state": 2,
      "output": "test-output-entity-api-6"
    }
    """
    When I send an event:
    """json
    {
      "connector": "test-connector-entity-api-7",
      "connector_name": "test-connector-name-entity-api-7",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-component-entity-api-7",
      "resource": "test-resource-entity-api-7",
      "state": 0,
      "output": "test-output-entity-api-7"
    }
    """
    When I wait the end of 2 events processing
    When I do GET /api/v4/entities?search=test-resource-entity-api-6
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "state": 2
        }
      ],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 1
      }
    }
    """
    When I do GET /api/v4/entities?search=test-resource-entity-api-7
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "state": 0
        }
      ],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 1
      }
    }
    """

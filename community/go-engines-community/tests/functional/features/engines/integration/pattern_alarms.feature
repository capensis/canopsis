Feature: Get matched alarms by patterns
  I need to be able get matched alarms by patterns
  Only admin should be able get matched alarms by patterns

  @concurrent
  Scenario: given get unauth request should not allow access
    When I do POST /api/v4/patterns-alarms
    Then the response code should be 401

  @concurrent
  Scenario: given patterns alarms request with missing fields should return empty response
    When I am admin
    When I do POST /api/v4/patterns-alarms:
    """json
    {
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "alarms": []
    }
    """

  @concurrent
  Scenario: given patterns alarms request with invalid patterns format should return bad request error
    When I am admin
    When I do POST /api/v4/patterns-alarms:
    """json
    {
      "alarm_pattern": [
        []
      ],
      "entity_pattern": [
        []
      ],
      "pbehavior_pattern": [
        []
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "alarm_pattern": "AlarmPattern is invalid alarm pattern.",
        "entity_pattern": "EntityPattern is invalid entity pattern.",
        "pbehavior_pattern": "PbehaviorPattern is invalid pbehavior pattern."
      }
    }
    """

  @concurrent
  Scenario: given patterns alarms request should return matched alarms by patterns
    When I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-pattern-alarms-connector-1",
      "connector_name": "test-pattern-alarms-connector-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-pattern-alarms-component-1",
      "resource": "test-pattern-alarms-resource-1-1",
      "state": 2,
      "output": "test-pattern-alarms-output",
      "long_output": "test-pattern-alarms-long-output",
      "author": "test-pattern-alarms-author"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-pattern-alarms-connector-1",
      "connector_name": "test-pattern-alarms-connector-1-name",
      "source_type": "resource",
      "event_type": "check",
      "component":  "test-pattern-alarms-component-1",
      "resource": "test-pattern-alarms-resource-1-2",
      "state": 2,
      "output": "test-pattern-alarms-output",
      "long_output": "test-pattern-alarms-long-output",
      "author": "test-pattern-alarms-author"
    }
    """
    When I do GET /api/v4/alarms?search=test-pattern-alarms-resource-1-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-pattern-alarms-resource-1-1"
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
    When I save response alarmID1={{ (index .lastResponse.data 0)._id }}
    When I save response alarmName1={{ (index .lastResponse.data 0).v.display_name }}
    When I do GET /api/v4/alarms?search=test-pattern-alarms-resource-1-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "resource": "test-pattern-alarms-resource-1-2"
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
    When I save response alarmID2={{ (index .lastResponse.data 0)._id }}
    When I save response alarmName2={{ (index .lastResponse.data 0).v.display_name }}
    When I do POST /api/v4/patterns-alarms:
    """json
    {
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "value": "test-pattern-alarms-resource-1-1",
              "type": "eq"
            }
          }
        ],
        [
          {
            "field": "name",
            "cond": {
              "value": "test-pattern-alarms-resource-1-2",
              "type": "eq"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    Then the response array key "alarms" should contain:
    """json
    [
      {
        "_id": "{{ .alarmID1 }}",
        "name": "{{ .alarmName1 }}"
      },
      {
        "_id": "{{ .alarmID2 }}",
        "name": "{{ .alarmName2 }}"
      }
    ]
    """
    When I do POST /api/v4/patterns-alarms:
    """json
    {
      "search": "{{ .alarmName1 }}",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "value": "test-pattern-alarms-resource-1-1",
              "type": "eq"
            }
          }
        ],
        [
          {
            "field": "name",
            "cond": {
              "value": "test-pattern-alarms-resource-1-2",
              "type": "eq"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    Then the response array key "alarms" should contain:
    """json
    [
      {
        "_id": "{{ .alarmID1 }}",
        "name": "{{ .alarmName1 }}"
      }
    ]
    """
    When I do POST /api/v4/patterns-alarms:
    """json
    {
      "search": "{{ .alarmName2 }}",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "value": "test-pattern-alarms-resource-1-1",
              "type": "eq"
            }
          }
        ],
        [
          {
            "field": "name",
            "cond": {
              "value": "test-pattern-alarms-resource-1-2",
              "type": "eq"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    Then the response array key "alarms" should contain:
    """json
    [
      {
        "_id": "{{ .alarmID2 }}",
        "name": "{{ .alarmName2 }}"
      }
    ]
    """
    When I do POST /api/v4/patterns-alarms:
    """json
    {
      "search": "some",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "value": "test-pattern-alarms-resource-1-1",
              "type": "eq"
            }
          }
        ],
        [
          {
            "field": "name",
            "cond": {
              "value": "test-pattern-alarms-resource-1-2",
              "type": "eq"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    Then the response array key "alarms" should contain:
    """json
    []
    """
    When I do POST /api/v4/patterns-alarms:
    """json
    {
      "alarm_pattern": [
        [
          {
            "field": "v.resource",
            "cond": {
              "value": "test-pattern-alarms-resource-1-1",
              "type": "eq"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    Then the response array key "alarms" should contain:
    """json
    [
      {
        "_id": "{{ .alarmID1 }}",
        "name": "{{ .alarmName1 }}"
      }
    ]
    """
    When I do POST /api/v4/patterns-alarms:
    """json
    {
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "value": "test-pattern-alarms-resource-1-1",
              "type": "eq"
            }
          }
        ]
      ],
      "alarm_pattern": [
        [
          {
            "field": "v.resource",
            "cond": {
              "value": "test-pattern-alarms-resource-1-2",
              "type": "eq"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    Then the response array key "alarms" should contain:
    """json
    []
    """
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pattern-alarms-1",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-alarms-resource-1-2"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then I wait the end of events processing which contain:
    """json
    [
      {
        "event_type": "pbhenter",
        "resource": "test-pattern-alarms-resource-1-2"
      }
    ]
    """
    When I do POST /api/v4/patterns-alarms:
    """json
    {
      "pbehavior_pattern": [
        [
          {
            "field": "pbehavior_info.canonical_type",
            "cond": {
              "type": "eq",
              "value": "maintenance"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    Then the response array key "alarms" should contain:
    """json
    [
      {
        "_id": "{{ .alarmID2 }}",
        "name": "{{ .alarmName2 }}"
      }
    ]
    """

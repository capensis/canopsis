Feature: Count matches
  I need to be able to count matches by patterns
  Only admin should be able to count matches by patterns

  Scenario: given count request should return counts
    When I am admin
    When I send an event:
    """json
    [
      {
        "connector": "test-connector-pbehavior-pattern-count-1",
        "connector_name": "test-connector-name-pbehavior-pattern-count-1",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-pbehavior-pattern-count-1",
        "resource": "test-resource-pbehavior-pattern-count-1-1",
        "state": 1,
        "output": "noveo alarm"
      },
      {
        "connector": "test-connector-pbehavior-pattern-count-1",
        "connector_name": "test-connector-name-pbehavior-pattern-count-1",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-pbehavior-pattern-count-1",
        "resource": "test-resource-pbehavior-pattern-count-1-2",
        "state": 1,
        "output": "noveo alarm"
      },
      {
        "connector": "test-connector-pbehavior-pattern-count-1",
        "connector_name": "test-connector-name-pbehavior-pattern-count-1",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-pbehavior-pattern-count-1",
        "resource": "test-resource-pbehavior-pattern-count-1-3",
        "state": 1,
        "output": "noveo alarm"
      },
      {
        "connector": "test-connector-pbehavior-pattern-count-1",
        "connector_name": "test-connector-name-pbehavior-pattern-count-1",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-pbehavior-pattern-count-1",
        "resource": "test-resource-pbehavior-pattern-count-1-4",
        "state": 1,
        "output": "noveo alarm"
      }
    ]
    """
    When I wait the end of 4 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-pattern-count-1",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine",
      "entity_pattern": [
        [
          {
            "field": "type",
            "cond": {
              "type": "eq",
              "value": "resource"
            }
          },
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-component-pbehavior-pattern-count-1"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response pbehaviorID={{ .lastResponse._id }}
    When I wait the end of 4 events processing
    When I do PUT /api/v4/internal/user_interface:
    """json
    {
      "max_matched_items": 4
    }
    """
    Then the response code should be 200
    Then I wait 2s
    When I am noperms
    When I do POST /api/v4/patterns-count:
    """json
    {
      "pbehavior_pattern": [
        [
          {
            "field": "pbehavior_info.id",
            "cond": {
              "type": "eq",
              "value": "{{ .pbehaviorID }}"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "pbehavior_pattern": {
        "count": 4,
        "over_limit": false
      }
    }
    """
    When I am admin
    When I do PUT /api/v4/internal/user_interface:
    """json
    {
      "max_matched_items": 3
    }
    """
    Then the response code should be 200
    Then I wait 2s
    When I am noperms
    When I do POST /api/v4/patterns-count:
    """json
    {
      "pbehavior_pattern": [
        [
          {
            "field": "pbehavior_info.id",
            "cond": {
              "type": "eq",
              "value": "{{ .pbehaviorID }}"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "pbehavior_pattern": {
        "count": 4,
        "over_limit": true
      }
    }
    """

  Scenario: given count requests for pbh reasons
    When I am admin
    When I send an event:
    """json
    [
      {
        "connector": "test-connector-pbehavior-pattern-count-2",
        "connector_name": "test-connector-name-pbehavior-pattern-count-2",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-pbehavior-pattern-count-2",
        "resource": "test-resource-pbehavior-pattern-count-2-1",
        "state": 1,
        "output": "noveo alarm"
      },
      {
        "connector": "test-connector-pbehavior-pattern-count-2",
        "connector_name": "test-connector-name-pbehavior-pattern-count-2",
        "source_type": "resource",
        "event_type": "check",
        "component": "test-component-pbehavior-pattern-count-2",
        "resource": "test-resource-pbehavior-pattern-count-2-2",
        "state": 1,
        "output": "noveo alarm"
      }
    ]
    """
    When I wait the end of 2 events processing
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-pattern-count-2",
      "tstart": {{ now }},
      "tstop": {{ nowAdd "1h" }},
      "color": "#FFFFFF",
      "type": "test-maintenance-type-to-engine",
      "reason": "test-reason-to-engine-3",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-pattern-count-2-1"
            }
          }
        ],
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-pattern-count-2-2"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I wait the end of 2 events processing
    Then I wait 2s
    When I do POST /api/v4/patterns-count:
    """json
    {
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-pattern-count-2-1"
            }
          }
        ],
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-pbehavior-pattern-count-2-2"
            }
          }
        ]
      ],
      "pbehavior_pattern": [
        [
          {
            "field": "pbehavior_info.reason",
            "cond": {
              "type": "eq",
              "value": "test-reason-to-engine-3"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "pbehavior_pattern": {
        "count": 2,
        "over_limit": false
      }
    }
    """
    When I do POST /api/v4/patterns-count:
    """json
    {
      "pbehavior_pattern": [
        [
          {
            "field": "pbehavior_info.reason",
            "cond": {
              "type": "eq",
              "value": "test-reason-to-engine-not-exist"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "pbehavior_pattern": {
        "count": 0,
        "over_limit": false
      }
    }
    """

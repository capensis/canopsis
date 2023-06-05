Feature: Get a idle rule
  I need to be able to read a idle rule
  Only admin should be able to read a idle rule

  Scenario: given get all request should return idle rules
    When I am admin
    When I do GET /api/v4/idle-rules?search=test-idle-rule-to-get
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-idle-rule-to-get-1",
          "alarm_condition": "last_event",
          "old_alarm_patterns": null,
          "old_entity_patterns": null,
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-idle-rule-to-get-1-alarm"
                }
              }
            ]
          ],
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-idle-rule-to-get-1-resource"
                }
              }
            ]
          ],
          "author": {
            "_id": "root",
            "name": "root"
          },
          "description": "test-idle-rule-to-get-1-description",
          "disable_during_periods": [
            "pause"
          ],
          "duration": {
            "value": 3,
            "unit": "s"
          },
          "enabled": true,
          "name": "test-idle-rule-to-get-1-name",
          "operation": {
            "parameters": {
              "duration": {
                "value": 3,
                "unit": "s"
              },
              "output": "test-idle-rule-to-get-1-operation-output"
            },
            "type": "snooze"
          },
          "type": "alarm"
        },
        {
          "_id": "test-idle-rule-to-get-2",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "description": "test-idle-rule-to-get-2-description",
          "disable_during_periods": null,
          "duration": {
            "value": 3,
            "unit": "s"
          },
          "enabled": true,
          "old_alarm_patterns": null,
          "old_entity_patterns": null,
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-idle-rule-to-get-2-resource"
                }
              }
            ]
          ],
          "name": "test-idle-rule-to-get-2-name",
          "type": "entity"
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

  Scenario: given sort request should return sorted idle rules
    When I am admin
    When I do GET /api/v4/idle-rules?search=test-idle-rule-to-get&sort_by=name&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-idle-rule-to-get-2"
        },
        {
          "_id": "test-idle-rule-to-get-1"
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

  Scenario: given get all request and no auth user should not allow access
    When I do GET /api/v4/idle-rules
    Then the response code should be 401

  Scenario: given get all request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/idle-rules
    Then the response code should be 403

  Scenario: given get request should return idle rule
    When I am admin
    When I do GET /api/v4/idle-rules/test-idle-rule-to-get-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-idle-rule-to-get-1",
      "alarm_condition": "last_event",
      "old_alarm_patterns": null,
      "old_entity_patterns": null,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-idle-rule-to-get-1-alarm"
            }
          }
        ]
      ],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-idle-rule-to-get-1-resource"
            }
          }
        ]
      ],
      "author": {
        "_id": "root",
        "name": "root"
      },
      "description": "test-idle-rule-to-get-1-description",
      "disable_during_periods": [
        "pause"
      ],
      "duration": {
        "value": 3,
        "unit": "s"
      },
      "enabled": true,
      "name": "test-idle-rule-to-get-1-name",
      "operation": {
        "parameters": {
          "duration": {
            "value": 3,
            "unit": "s"
          },
          "output": "test-idle-rule-to-get-1-operation-output"
        },
        "type": "snooze"
      },
      "type": "alarm"
    }
    """

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/idle-rules/notexist
    Then the response code should be 401

  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/idle-rules/test-exception-to-get-1
    Then the response code should be 403

  Scenario: given invalid get request should return error
    When I am admin
    When I do GET /api/v4/idle-rules/notexist
    Then the response code should be 404

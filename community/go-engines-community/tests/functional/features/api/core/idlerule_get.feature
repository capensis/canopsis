Feature: Get a idle rule
  I need to be able to read a idle rule
  Only admin should be able to read a idle rule

  Scenario: given get all request should return idle rules
    When I am admin
    When I do GET /api/v4/idle-rules?search=test-idle-rule-to-get
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "test-idle-rule-to-get-1",
          "alarm_condition": "last_event",
          "alarm_patterns": [
            {
              "_id": "test-idle-rule-to-get-1-alarm"
            }
          ],
          "author": "test-idle-rule-to-get-1-author",
          "created": 1616567033,
          "description": "test-idle-rule-to-get-1-description",
          "disable_during_periods": [
            "pause"
          ],
          "duration": {
            "value": 3,
            "unit": "s"
          },
          "enabled": true,
          "entity_patterns": [
            {
              "name": "test-idle-rule-to-get-1-resource"
            }
          ],
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
          "priority": 1,
          "type": "alarm",
          "updated": 1616567033
        },
        {
          "_id": "test-idle-rule-to-get-2",
          "alarm_patterns": null,
          "author": "test-idle-rule-to-get-2-author",
          "created": 1616567033,
          "description": "test-idle-rule-to-get-2-description",
          "disable_during_periods": null,
          "duration": {
            "value": 5,
            "unit": "s"
          },
          "enabled": true,
          "entity_patterns": [
            {
              "name": "test-idle-rule-to-get-2-resource"
            }
          ],
          "name": "test-idle-rule-to-get-2-name",
          "priority": 2,
          "type": "entity",
          "updated": 1616567033
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
    When I do GET /api/v4/idle-rules?search=test-idle-rule-to-get&sort_by=duration&sort=desc
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
    Then the response body should be:
    """json
    {
      "_id": "test-idle-rule-to-get-1",
      "alarm_condition": "last_event",
      "alarm_patterns": [
        {
          "_id": "test-idle-rule-to-get-1-alarm"
        }
      ],
      "author": "test-idle-rule-to-get-1-author",
      "created": 1616567033,
      "description": "test-idle-rule-to-get-1-description",
      "disable_during_periods": [
        "pause"
      ],
      "duration": {
        "value": 3,
        "unit": "s"
      },
      "enabled": true,
      "entity_patterns": [
        {
          "name": "test-idle-rule-to-get-1-resource"
        }
      ],
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
      "priority": 1,
      "type": "alarm",
      "updated": 1616567033
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

Feature: Update a declare ticket rule
  I need to be able to update a declare ticket rule
  Only admin should be able to update a declare ticket rule

  Scenario: given update request should update rule
    When I am admin
    Then I do PUT /api/v4/cat/declare-ticket-rules/test-declare-ticket-rule-to-update-1:
    """json
    {
      "name": "test-declare-ticket-rule-to-update-1-name-updated",
      "system_name": "test-declare-ticket-rule-to-update-1-system-name-updated",
      "enabled": true,
      "emit_trigger": true,
      "webhooks": [
        {
          "request": {
            "url": "https://canopsis-test.com",
            "method": "GET",
            "auth": {
              "username": "test",
              "password": "test"
            },
            "skip_verify": true,
            "timeout": {
              "value": 30,
              "unit": "s"
            },
            "retry_count": 3,
            "retry_delay": {
              "value": 1,
              "unit": "s"
            }
          },
          "stop_on_fail": true
        },
        {
          "request": {
            "url": "https://canopsis-test.com",
            "method": "POST",
            "headers": {
              "Content-Type": "application/json"
            },
            "payload": "{\"name\": \"test-declare-ticket-rule-to-update-1-payload\"}",
            "skip_verify": true,
            "timeout": {
              "value": 30,
              "unit": "s"
            },
            "retry_count": 3,
            "retry_delay": {
              "value": 1,
              "unit": "s"
            }
          },
          "declare_ticket": {
            "is_regexp": false,
            "ticket_id": "_id",
            "ticket_url": "url",
            "ticket_custom": "custom"
          }
        }
      ],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-declare-ticket-rule-to-update-1-resource"
            }
          }
        ]
      ],
      "alarm_pattern": [
        [
          {
            "field": "v.resource",
            "cond": {
              "type": "eq",
              "value": "test-declare-ticket-rule-to-update-1-resource"
            }
          }
        ]
      ],
      "pbehavior_pattern": [
        [
          {
            "field": "pbehavior_info.canonical_type",
            "cond": {
              "type": "eq",
              "value": "active"
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
      "author": {
        "_id": "root",
        "name": "root"
      },
      "created": 1619083733,
      "name": "test-declare-ticket-rule-to-update-1-name-updated",
      "system_name": "test-declare-ticket-rule-to-update-1-system-name-updated",
      "enabled": true,
      "emit_trigger": true,
      "webhooks": [
        {
          "request": {
            "url": "https://canopsis-test.com",
            "method": "GET",
            "auth": {
              "username": "test",
              "password": "test"
            },
            "skip_verify": true,
            "timeout": {
              "value": 30,
              "unit": "s"
            },
            "retry_count": 3,
            "retry_delay": {
              "value": 1,
              "unit": "s"
            }
          },
          "stop_on_fail": true
        },
        {
          "request": {
            "url": "https://canopsis-test.com",
            "method": "POST",
            "headers": {
              "Content-Type": "application/json"
            },
            "payload": "{\"name\": \"test-declare-ticket-rule-to-update-1-payload\"}",
            "skip_verify": true,
            "timeout": {
              "value": 30,
              "unit": "s"
            },
            "retry_count": 3,
            "retry_delay": {
              "value": 1,
              "unit": "s"
            }
          },
          "declare_ticket": {
            "is_regexp": false,
            "ticket_id": "_id",
            "ticket_url": "url",
            "ticket_custom": "custom"
          }
        }
      ],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-declare-ticket-rule-to-update-1-resource"
            }
          }
        ]
      ],
      "alarm_pattern": [
        [
          {
            "field": "v.resource",
            "cond": {
              "type": "eq",
              "value": "test-declare-ticket-rule-to-update-1-resource"
            }
          }
        ]
      ],
      "pbehavior_pattern": [
        [
          {
            "field": "pbehavior_info.canonical_type",
            "cond": {
              "type": "eq",
              "value": "active"
            }
          }
        ]
      ]
    }
    """

  Scenario: given get request and no auth user should not allow access
    When I do PUT /api/v4/cat/declare-ticket-rules/test-declare-ticket-rule-not-found
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/cat/declare-ticket-rules/test-declare-ticket-rule-not-found
    Then the response code should be 403

  Scenario: given update request with not exist id should return not found error
    When I am admin
    When I do PUT /api/v4/cat/declare-ticket-rules/test-declare-ticket-rule-not-found:
    """json
    {
      "name": "test-declare-ticket-rule-not-found-name",
      "system_name": "test-declare-ticket-rule-not-found-system-name",
      "enabled": true,
      "emit_trigger": true,
      "webhooks": [
        {
          "request": {
            "url": "https://canopsis-test.com",
            "method": "POST"
          },
          "declare_ticket": {
            "ticket_id": "_id"
          }
        }
      ],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-declare-ticket-rule-not-found-resource"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """

  Scenario: given update request with missing fields should return bad request
    When I am admin
    Then I do PUT /api/v4/cat/declare-ticket-rules/test-declare-ticket-rule-to-update:
    """json
    {}
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "name": "Name is missing.",
        "enabled": "Enabled is missing.",
        "emit_trigger": "EmitTrigger is missing.",
        "webhooks": "Webhooks is missing.",
        "entity_pattern": "EntityPattern is missing.",
        "corporate_entity_pattern": "CorporateEntityPattern is missing.",
        "alarm_pattern": "AlarmPattern is missing.",
        "corporate_alarm_pattern": "CorporateAlarmPattern is missing.",
        "pbehavior_pattern": "PbehaviorPattern is missing.",
        "corporate_pbehavior_pattern": "CorporatePbehaviorPattern is missing."
      }
    }
    """

  Scenario: given update request with already exists name should return error
    When I am admin
    When I do PUT /api/v4/cat/declare-ticket-rules/test-declare-ticket-rule-to-update-1:
    """json
    {
      "name": "test-declare-ticket-rule-to-check-unique-name"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "name": "Name already exists."
      }
    }
    """

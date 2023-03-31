Feature: Create a declare ticket rule
  I need to be able to create a declare ticket rule
  Only admin should be able to create a declare ticket rule

  Scenario: given create request should return ok
    When I am admin
    When I do POST /api/v4/cat/declare-ticket-rules:
    """json
    {
      "name": "test-declare-ticket-rule-to-create-1-name",
      "system_name": "test-declare-ticket-rule-to-create-1-system-name",
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
            "payload": "{\"name\": \"test-declare-ticket-rule-to-create-1-payload\"}",
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
              "value": "test-declare-ticket-rule-to-create-1-resource"
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
              "value": "test-declare-ticket-rule-to-create-1-resource"
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
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "name": "test-declare-ticket-rule-to-create-1-name",
      "system_name": "test-declare-ticket-rule-to-create-1-system-name",
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
            "payload": "{\"name\": \"test-declare-ticket-rule-to-create-1-payload\"}",
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
              "value": "test-declare-ticket-rule-to-create-1-resource"
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
              "value": "test-declare-ticket-rule-to-create-1-resource"
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
    When I do GET /api/v4/cat/declare-ticket-rules/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "name": "test-declare-ticket-rule-to-create-1-name",
      "system_name": "test-declare-ticket-rule-to-create-1-system-name",
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
            "payload": "{\"name\": \"test-declare-ticket-rule-to-create-1-payload\"}",
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
              "value": "test-declare-ticket-rule-to-create-1-resource"
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
              "value": "test-declare-ticket-rule-to-create-1-resource"
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

  Scenario: given create request with missing fields should return bad request
    When I am admin
    When I do POST /api/v4/cat/declare-ticket-rules:
    """json
    {
    }
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
    When I do POST /api/v4/cat/declare-ticket-rules:
    """json
    {
      "webhooks": [
        {},
        {}
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "webhooks.0.request.url": "URL is missing.",
        "webhooks.0.request.method": "Method is missing.",
        "webhooks.0.stop_on_fail": "StopOnFail is missing.",
        "webhooks.1.request.url": "URL is missing.",
        "webhooks.1.request.method": "Method is missing."
      }
    }
    """
    When I do POST /api/v4/cat/declare-ticket-rules:
    """json
    {
      "webhooks": [
        {
          "declare_ticket": {
            "ticket_id": "_id"
          }
        },
        {
          "declare_ticket": {
            "ticket_id": "_id"
          }
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "webhooks.1.declare_ticket.ticket_id": "TicketID is not empty."
      }
    }
    """

  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/cat/declare-ticket-rules
    Then the response code should be 401

  Scenario: given create request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/cat/declare-ticket-rules
    Then the response code should be 403

  Scenario: given create request with already exists name should return error
    When I am admin
    When I do POST /api/v4/cat/declare-ticket-rules:
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

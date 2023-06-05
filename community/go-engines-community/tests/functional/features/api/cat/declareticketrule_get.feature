Feature: Get a declare ticket rule
  I need to be able to get a declare ticket rule
  Only admin should be able to get a declare ticket rule

  Scenario: given search request should return rules
    When I am admin
    When I do GET /api/v4/cat/declare-ticket-rules?search=test-declare-ticket-rule-to-get
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "test-declare-ticket-rule-to-get-1",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "created": 1619083733,
          "updated": 1619083733,
          "name": "test-declare-ticket-rule-to-get-1-name",
          "system_name": "test-declare-ticket-rule-to-get-1-system-name",
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
                "headers": null,
                "payload": "",
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
              "declare_ticket": null,
              "stop_on_fail": true
            },
            {
              "request": {
                "url": "https://canopsis-test.com",
                "method": "POST",
                "auth": null,
                "headers": {
                  "Content-Type": "application/json"
                },
                "payload": "{\"name\": \"test-declare-ticket-rule-to-get-1-payload\"}",
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
                "ticket_custom": "custom",
                "empty_response": false
              }
            }
          ],
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-declare-ticket-rule-to-get-1-resource"
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
                  "value": "test-declare-ticket-rule-to-get-1-resource"
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
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """

  Scenario: given get request should return rule
    When I am admin
    When I do GET /api/v4/cat/declare-ticket-rules/test-declare-ticket-rule-to-get-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-declare-ticket-rule-to-get-1",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "created": 1619083733,
      "updated": 1619083733,
      "name": "test-declare-ticket-rule-to-get-1-name",
      "system_name": "test-declare-ticket-rule-to-get-1-system-name",
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
            "headers": null,
            "payload": "",
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
          "declare_ticket": null,
          "stop_on_fail": true
        },
        {
          "request": {
            "url": "https://canopsis-test.com",
            "method": "POST",
            "auth": null,
            "headers": {
              "Content-Type": "application/json"
            },
            "payload": "{\"name\": \"test-declare-ticket-rule-to-get-1-payload\"}",
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
            "ticket_custom": "custom",
            "empty_response": false
          }
        }
      ],
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-declare-ticket-rule-to-get-1-resource"
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
              "value": "test-declare-ticket-rule-to-get-1-resource"
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

  Scenario: given get all request and no auth user should not allow access
    When I do GET /api/v4/cat/declare-ticket-rules
    Then the response code should be 401

  Scenario: given get all request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/cat/declare-ticket-rules
    Then the response code should be 403

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/cat/declare-ticket-rules/test-declare-ticket-rule-not-found
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/cat/declare-ticket-rules/test-declare-ticket-rule-not-found
    Then the response code should be 403

  Scenario: given get request with not exist id should return not found error
    When I am admin
    When I do GET /api/v4/cat/declare-ticket-rules/test-declare-ticket-rule-not-found
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """

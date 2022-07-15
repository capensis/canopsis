Feature: Get a resolve rule
  I need to be able to get a resolve rule
  Only admin should be able to get a resolve rule

  Scenario: given search request should return resolve rules
    When I am admin
    When I do GET /api/v4/resolve-rules?search=test-resolve-rule-to-get
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "test-resolve-rule-to-get-1",
          "alarm_pattern": [
            [
              {
                "field": "v.connector",
                "cond": {
                  "type": "eq",
                  "value": "test-resolve-rule-to-get-1-pattern"
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
                  "value": "test-resolve-rule-to-get-1-resource"
                }
              }
            ]
          ],
          "old_alarm_patterns": null,
          "old_entity_patterns": null,
          "author": {
            "_id": "root",
            "name": "root"
          },
          "created": 1619083733,
          "name": "test-resolve-rule-to-get-1-name",
          "description": "test-resolve-rule-to-get-1-description",
          "duration": {
            "value": 10,
            "unit": "s"
          },
          "priority": 1,
          "updated": 1619083733
        },
        {
          "_id": "test-resolve-rule-to-get-2",
          "alarm_pattern": [
            [
              {
                "field": "v.connector",
                "cond": {
                  "type": "eq",
                  "value": "test-resolve-rule-to-get-2-pattern"
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
                  "value": "test-resolve-rule-to-get-2-resource"
                }
              }
            ]
          ],
          "old_alarm_patterns": null,
          "old_entity_patterns": null,
          "author": {
            "_id": "root",
            "name": "root"
          },
          "created": 1619083733,
          "name": "test-resolve-rule-to-get-2-name",
          "description": "test-resolve-rule-to-get-2-description",
          "duration": {
            "value": 10,
            "unit": "s"
          },
          "priority": 2,
          "updated": 1619083733
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

  Scenario: given get request should return resolve rule
    When I am admin
    When I do GET /api/v4/resolve-rules/test-resolve-rule-to-get-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-resolve-rule-to-get-1",
      "alarm_pattern": [
        [
          {
            "field": "v.connector",
            "cond": {
              "type": "eq",
              "value": "test-resolve-rule-to-get-1-pattern"
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
              "value": "test-resolve-rule-to-get-1-resource"
            }
          }
        ]
      ],
      "old_alarm_patterns": null,
      "old_entity_patterns": null,
      "author": {
        "_id": "root",
        "name": "root"
      },
      "created": 1619083733,
      "name": "test-resolve-rule-to-get-1-name",
      "description": "test-resolve-rule-to-get-1-description",
      "duration": {
        "value": 10,
        "unit": "s"
      },
      "priority": 1,
      "updated": 1619083733
    }
    """

  Scenario: given get all request and no auth user should not allow access
    When I do GET /api/v4/resolve-rules
    Then the response code should be 401

  Scenario: given get all request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/resolve-rules
    Then the response code should be 403

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/resolve-rules/test-resolve-rule-to-get-1
    Then the response code should be 401

  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/resolve-rules/test-resolve-rule-to-get-1
    Then the response code should be 403

  Scenario: given get request with not exist id should return not found error
    When I am admin
    When I do GET /api/v4/resolve-rules/test-resolve-rule-not-found
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """

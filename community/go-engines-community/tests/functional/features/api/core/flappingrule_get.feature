Feature: Get a flapping rule
  I need to be able to get a flapping rule
  Only admin should be able to get a flapping rule

  Scenario: given search request should return flapping rules
    When I am admin
    When I do GET /api/v4/flapping-rules?search=test-flapping-rule-to-get
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "test-flapping-rule-to-get-2",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "created": 1619083533,
          "description": "flapping rule 2",
          "flapping_interval": {
            "seconds": 23,
            "unit": "s"
          },
          "flapping_freq_limit": 5,
          "alarm_patterns": null,
          "updated": 1619083533,
          "priority": 1
        },
        {
          "_id": "test-flapping-rule-to-get-1",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "created": 1619083733,
          "description": "flapping rule 1",
          "flapping_interval": {
            "seconds": 10,
            "unit": "s"
          },
          "flapping_freq_limit": 3,
          "alarm_patterns": [
            {
              "v": {
                "connector": "test-flapping-rule-to-get-1-pattern"
              }
            }
          ],
          "updated": 1619083733,
          "priority": 0
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

  Scenario: given get request should return flapping rule
    When I am admin
    When I do GET /api/v4/flapping-rules/test-flapping-rule-to-get-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-flapping-rule-to-get-1",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "created": 1619083733,
      "description": "flapping rule 1",
      "flapping_interval": {
        "seconds": 10,
        "unit": "s"
      },
      "flapping_freq_limit": 3,
      "alarm_patterns": [
        {
          "v": {
            "connector": "test-flapping-rule-to-get-1-pattern"
          }
        }
      ],
      "updated": 1619083733,
      "priority": 0
    }
    """

  Scenario: given get all request and no auth user should not allow access
    When I do GET /api/v4/flapping-rules
    Then the response code should be 401

  Scenario: given get all request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/flapping-rules
    Then the response code should be 403

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/flapping-rules/test-flapping-rule-to-get-1
    Then the response code should be 401

  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/flapping-rules/test-flapping-rule-to-get-1
    Then the response code should be 403

  Scenario: given get request with not exist id should return not found error
    When I am admin
    When I do GET /api/v4/flapping-rules/test-flapping-rule-not-found
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """

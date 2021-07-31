Feature: Get a baggot rule
  I need to be able to get a baggot rule
  Only admin should be able to get a baggot rule

  Scenario: given search request should return baggot rules
    When I am admin
    When I do GET /api/v4/baggot-rules?search=test-baggot-rule-to-get
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "test-baggot-rule-to-get-2",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "created": 1619083533,
          "description": "baggot rule 2",
          "duration": {
            "seconds": 10,
            "unit": "s"
          },
          "alarm_patterns": null,
          "entity_patterns": [{"name": "test-baggot-rule-to-get-2-pattern"}],
          "updated": 1619083533,
          "priority": 1
        },
        {
          "_id": "test-baggot-rule-to-get-1",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "created": 1619083733,
          "description": "baggot rule 1",
          "duration": {
            "seconds": 10,
            "unit": "s"
          },
          "alarm_patterns": [
            {
              "v": {
                "connector": "test-baggot-rule-to-get-1-pattern"
              }
            }
          ],
          "entity_patterns": null,
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

  Scenario: given get request should return baggot rule
    When I am admin
    When I do GET /api/v4/baggot-rules/test-baggot-rule-to-get-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-baggot-rule-to-get-1",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "description": "baggot rule 1",
      "duration": {
        "seconds": 10,
        "unit": "s"
      },
      "alarm_patterns": [
        {
          "v": {
            "connector": "test-baggot-rule-to-get-1-pattern"
          }
        }
      ],
      "entity_patterns": null,
      "created": 1619083733,
      "updated": 1619083733,
      "priority": 0
    }
    """

  Scenario: given get all request and no auth user should not allow access
    When I do GET /api/v4/baggot-rules
    Then the response code should be 401

  Scenario: given get all request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/baggot-rules
    Then the response code should be 403

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/baggot-rules/test-baggot-rule-to-get-1
    Then the response code should be 401

  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/baggot-rules/test-baggot-rule-to-get-1
    Then the response code should be 403

  Scenario: given get request with not exist id should return not found error
    When I am admin
    When I do GET /api/v4/baggot-rules/test-baggot-rule-not-found
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """

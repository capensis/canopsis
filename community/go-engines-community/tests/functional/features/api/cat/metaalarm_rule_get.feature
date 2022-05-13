Feature: Get a metaalarm-rule
  I need to be able to get a metaalarm-rule
  Only admin should be able to get a metaalarm-rule

  Scenario: given search request should return metaalarm-rule
    When I am admin
    When I do GET /api/v4/cat/metaalarmrules?search=test-metaalarm-to-get
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "data": [
        {
          "_id": "test-metaalarm-to-get-1",
          "auto_resolve": false,
          "config": {
            "time_interval": {
              "value": 10,
              "unit": "s"
            }
          },
          "name": "Test alarm get",
          "author": "test-metaalarm-to-get-1-author",
          "type": "complex",
          "output_template": "",
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-pattern-to-get-1-pattern"
                }
              }
            ]
          ]
        },
        {
          "_id": "test-metaalarm-to-get-2",
          "auto_resolve": false,
          "config": {
            "time_interval": {
              "value": 10,
              "unit": "s"
            }
          },
          "name": "Test alarm get",
          "author": "test-metaalarm-to-get-2-author",
          "type": "complex",
          "output_template": "",
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-pattern-to-get-2-pattern"
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
        "total_count": 2
      }
    }
    """

  Scenario: given get request should return metaalarm-rule
    When I am admin
    When I do GET /api/v4/cat/metaalarmrules/test-metaalarm-to-get-1
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "_id": "test-metaalarm-to-get-1",
      "auto_resolve": false,
      "config": {
        "time_interval": {
          "value": 10,
          "unit": "s"
        }
      },
      "name": "Test alarm get",
      "author": "test-metaalarm-to-get-1-author",
      "type": "complex",
      "output_template": "",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-get-1-pattern"
            }
          }
        ]
      ]
    }
    """

  Scenario: given get all request and no auth user should not allow access
    When I do GET /api/v4/cat/metaalarmrules
    Then the response code should be 401

  Scenario: given get all request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/cat/metaalarmrules
    Then the response code should be 403

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/cat/metaalarmrules/test-metaalarm-rule-to-get-1
    Then the response code should be 401

  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/cat/metaalarmrules/test-metaalarm-rule-to-get-1
    Then the response code should be 403

  Scenario: given get request with not exist id should return not found error
    When I am admin
    When I do GET /api/v4/cat/metaalarmrules/test-metaalarm-rule-not-found
    Then the response code should be 404
    Then the response body should be:
    """
    {
      "error": "Not found"
    }
    """

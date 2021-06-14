Feature: get test suite's historical-data
  I need to be able to get test suite's historical-data
  Only admin should be able to get test suite's historical-data

  Scenario: GET unauthorized
    When I do GET /api/v4/cat/junit/test-suites-history/7e5d2a21-be12-47d1-a286-0b46b6b2b99b
    Then the response code should be 401

  Scenario: GET without permissions
    When I am noperms
    When I do GET /api/v4/cat/junit/test-suites-history/7e5d2a21-be12-47d1-a286-0b46b6b2b99b
    Then the response code should be 403

  Scenario: GET historical-data success
    When I am admin
    When I do GET /api/v4/cat/junit/test-suites-history/7e5d2a21-be12-47d1-a286-0b46b6b2b99b
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "data": [
        {
          "_id": "5b4461d4-cfea-41dd-97cd-a5b15c34c1e3",
          "created": 1614782435,
          "last_update": 1614782435,
          "state": 3
        },
        {
          "_id": "5b4461d4-cfea-41dd-97cd-a5b15c34c1e2",
          "created": 1614782430,
          "last_update": 1614782430,
          "state": 1
        },
        {
          "_id": "5b4461d4-cfea-41dd-97cd-a5b15c34c1e1",
          "created": 1614782425,
          "last_update": 1614782425,
          "state": 2
        }
      ],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 3
      }
    }
    """

  Scenario: GET historical-data, test-suite not found
    When I am admin
    When I do GET /api/v4/cat/junit/test-suites-history/7e5d2a21-be12-47d1-a286-0b46b6b2b99b-not-found
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "data": [],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 0
      }
    }
    """
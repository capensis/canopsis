Feature: get a reason

  Scenario: given get all request should return reasons
    When I am admin
    When I do GET /api/v4/pbehavior-reasons?search=test-reason-to-get
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "data": [
        {
          "_id": "test-reason-to-get-1",
          "description": "test-reason-to-get-1-description",
          "name": "test-reason-to-get-1-name"
        },
        {
          "_id": "test-reason-to-get-2",
          "description": "test-reason-to-get-2-description",
          "name": "test-reason-to-get-2-name"
        },
        {
          "_id": "test-reason-to-get-3",
          "description": "test-reason-to-get-3-description",
          "name": "test-reason-to-get-3-name"
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 3
      }
    }
    """

  Scenario: given get all request should return reasons with flags
    When I am admin
    When I do GET /api/v4/pbehavior-reasons?search=test-reason-to-get&with_flags=true
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-reason-to-get-1",
          "deletable": true
        },
        {
          "_id": "test-reason-to-get-2",
          "deletable": false
        },
        {
          "_id": "test-reason-to-get-3",
          "deletable": false
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 3
      }
    }
    """

  Scenario: GET a PBehavior but unauthorized
    When I do GET /api/v4/pbehavior-reasons/test-reason-to-get-1
    Then the response code should be 401

  Scenario: GET a PBehavior but without permissions
    When I am noperms
    When I do GET /api/v4/pbehavior-reasons/test-reason-to-get-1
    Then the response code should be 403

  Scenario: Get a reason with success
    When I am admin
    When I do GET /api/v4/pbehavior-reasons/test-reason-to-get-1
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "_id": "test-reason-to-get-1",
      "description": "test-reason-to-get-1-description",
      "name": "test-reason-to-get-1-name"
    }
    """

  Scenario: Get a PBehavior with not found response
    When I am admin
    When I do GET /api/v4/pbehavior-reasons/test-reason-not-found
    Then the response code should be 404
    Then the response body should be:
    """
    {
      "error": "Not found"
    }
    """

Feature: Create an broadcast-message
  I need to be able to create a broadcast-message
  Only admin should be able to create a broadcast-message

  Scenario: given create request should return ok
    When I am admin
    When I do POST /api/v4/broadcast-message:
    """
    {
      "color": "rgb(159, 13, 0)",
      "message": "test-broadcast-create",
      "start": 1612139798,
      "end": 1612296900
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """
    {
      "color": "rgb(159, 13, 0)",
      "message": "test-broadcast-create",
      "start": 1612139798,
      "end": 1612296900
    }
    """
    When I do GET /api/v4/broadcast-message/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "color": "rgb(159, 13, 0)",
      "message": "test-broadcast-create",
      "start": 1612139798,
      "end": 1612296900
    }
    """

  Scenario: given create request with wrong type should return bad request
    When I am admin
    When I do POST /api/v4/broadcast-message:
    """
    {
      "color": "#23"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "color": "Color is not valid."
      }
    }
    """

  Scenario: given create request with missing fields should return bad request
    When I am admin
    When I do POST /api/v4/broadcast-message:
    """
    {
      "color": "rgb(159, 13, 0, 23)"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "color": "Color is not valid."
      }
    }
    """

  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/broadcast-message
    Then the response code should be 401

  Scenario: given create request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/broadcast-message
    Then the response code should be 403

  Scenario: given create request with already exists name should return error
    When I am admin
    When I do POST /api/v4/broadcast-message:
    """
    {
      "_id": "test-broadcast-to-check-id"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "_id": "ID already exists."
      }
    }
    """
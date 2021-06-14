Feature: Update an broadcast-message
  I need to be able to update an broadcast-message
  Only admin should be able to update an broadcast-message

  Scenario: given update request should update broadcast-message
    When I am admin
    Then I do PUT /api/v4/broadcast-message/test-broadcast-to-update:
    """
    {
      "color": "#32a852",
      "message": "test-broadcast-to-update-updated",
      "start": 1612139798,
      "end": 2612296900
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "_id": "test-broadcast-to-update",
      "color": "#32a852",
      "message": "test-broadcast-to-update-updated",
      "start": 1612139798,
      "end": 2612296900,
      "created": 1612139798
    }
    """

  Scenario: given get request and no auth user should not allow access
    When I do PUT /api/v4/broadcast-message/test-broadcast-to-update
    Then the response code should be 401

  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/broadcast-message/test-broadcast-to-update
    Then the response code should be 403

  Scenario: given update request with not exist id should return not found error
    When I am admin
    When I do PUT /api/v4/broadcast-message/test-broadcast-not-found:
    """
    {
      "color": "#32a852",
      "message": "test-broadcast-to-update-updated",
      "start": 1612139798,
      "end": 2612296900
    }
    """
    Then the response code should be 404
    Then the response body should be:
    """
    {
      "error": "Not found"
    }
    """
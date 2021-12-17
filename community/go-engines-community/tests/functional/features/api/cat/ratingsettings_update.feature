Feature: Update a rating settings
  I need to be able to update a rating settings
  Only admin should be able to update a rating settings

  Scenario: given update request should update rating settings
    When I am admin
    Then I do PUT /api/v4/cat/rating-settings/7:
    """json
    {
      "enabled": true
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/cat/rating-settings?search=test-rating-settings-to-update
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "id": 7,
          "label": "test-rating-settings-to-update",
          "enabled": true
        }
      ]
    }
    """
    Then I do PUT /api/v4/cat/rating-settings/7:
    """json
    {
      "enabled": false
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/cat/rating-settings?search=test-rating-settings-to-update
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "id": 7,
          "label": "test-rating-settings-to-update",
          "enabled": false
        }
      ]
    }
    """

  Scenario: given get request and no auth user should not allow access
    When I do PUT /api/v4/cat/rating-settings/10000
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/cat/rating-settings/10000
    Then the response code should be 403

  Scenario: given update request with not exist id should return not found error
    When I am admin
    When I do PUT /api/v4/cat/rating-settings/10000:
    """json
    {
      "enabled": false
    }
    """
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """

  Scenario: given update request with missing fields should return bad request
    When I am admin
    Then I do PUT /api/v4/cat/rating-settings/10000:
    """
    {}
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "enabled": "Enabled is missing."
      }
    }
    """

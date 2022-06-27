Feature: Bulk update a rating settings
  I need to be able to bulk update a rating settings
  Only admin should be able to bulk update a rating settings

  Scenario: given get request and no auth user should not allow access
    When I do PUT /api/v4/cat/rating-settings/bulk
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/cat/rating-settings/bulk
    Then the response code should be 403

  Scenario: given bulk update request should update rating settings
    When I am admin
    Then I do PUT /api/v4/cat/rating-settings/bulk:
    """json
    [
      {
          "id": 8,
          "enabled": true
      },
      {
          "id": 9,
          "enabled": false
      },
      {
          "id": 10,
          "enabled": true
      }
    ]
    """
    Then the response code should be 204
    When I do GET /api/v4/cat/rating-settings?search=test-rating-settings-to-bulk-update
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "id": 8,
          "label": "test-rating-settings-to-bulk-update-1",
          "enabled": true
        },
        {
          "id": 9,
          "label": "test-rating-settings-to-bulk-update-2",
          "enabled": false
        },
        {
          "id": 10,
          "label": "test-rating-settings-to-bulk-update-3",
          "enabled": true
        }
      ]
    }
    """

  Scenario: given bulk update request with missing fields should return bad request
    When I am admin
    Then I do PUT /api/v4/cat/rating-settings/bulk:
    """
    {}
    """
    Then the response code should be 400
    Then I do PUT /api/v4/cat/rating-settings/bulk:
    """
    [{}]
    """
    Then the response code should be 400
    Then I do PUT /api/v4/cat/rating-settings/bulk:
    """
    [
      {
          "id": 9
      },
      {
          "id": 10,
          "enabled": true
      }
    ]
    """
    Then the response code should be 400
    Then I do PUT /api/v4/cat/rating-settings/bulk:
    """
    [
      {
          "id": 9,
          "enabled": true
      },
      {
          "enabled": true
      }
    ]
    """
    Then the response code should be 400

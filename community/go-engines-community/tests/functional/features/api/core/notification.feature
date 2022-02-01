Feature: update notifications
  I need to be able to update notifications
  Only admin should be able to update notifications

  Scenario: GET unauthorized
    When I do GET /api/v4/notification
    Then the response code should be 401

  Scenario: GET without permissions
    When I am noperms
    When I do GET /api/v4/notification
    Then the response code should be 403

  Scenario: PUT unauthorized
    When I do PUT /api/v4/notification
    Then the response code should be 401

  Scenario: PUT without permissions
    When I am noperms
    When I do PUT /api/v4/notification
    Then the response code should be 403

  Scenario: GET and PUT ok
    When I am admin
    When I do PUT /api/v4/notification:
    """
    {
      "instruction": {
        "rate": false,
        "rate_frequency": {
          "value": 1,
          "unit": "s"
        }
      }
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "instruction": {
        "rate": false,
        "rate_frequency": {
          "value": 1,
          "unit": "s"
        }
      }
    }
    """
    When I do GET /api/v4/notification
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "instruction": {
        "rate": false,
        "rate_frequency": {
          "value": 1,
          "unit": "s"
        }
      }
    }
    """

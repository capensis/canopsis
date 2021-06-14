Feature: Update a heartbeat
  I need to be able to update a heartbeat
  Only admin should be able to update a heartbeat

  Scenario: given update request should update heartbeat
    When I am admin
    Then I do PUT /api/v4/heartbeats/test-heartbeat-to-update:
    """
    {
      "name": "test-heartbeat-to-update-name-updated",
      "description": "test-heartbeat-to-update-description-updated",
      "author": "root-updated",
      "pattern": {
        "name": "test-heartbeat-to-update-resource-updated"
      },
      "expected_interval": "1000h",
      "output": "test-heartbeat-to-update-1-output-updated"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "_id": "test-heartbeat-to-update",
      "name": "test-heartbeat-to-update-name-updated",
      "description": "test-heartbeat-to-update-description-updated",
      "author": "root-updated",
      "pattern": {
        "name": "test-heartbeat-to-update-resource-updated"
      },
      "expected_interval": "1000h",
      "output": "test-heartbeat-to-update-1-output-updated",
      "created": 1592215337
    }
    """

  Scenario: given get request and no auth user should not allow access
    When I do PUT /api/v4/heartbeats/test-heartbeat-to-update:
    """
    {
      "name": "test-heartbeat-to-update-name-1-updated",
      "description": "test-heartbeat-to-update-1-description-updated",
      "author": "root-updated",
      "pattern": {
        "name": "test-heartbeat-to-update-1-resource-updated"
      },
      "expected_interval": "1000h",
      "output": "test-heartbeat-to-update-1-output-updated"
    }
    """
    Then the response code should be 401

  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/heartbeats/test-heartbeat-to-update:
    """
    {
      "name": "test-heartbeat-to-update-1-name-updated",
      "description": "test-heartbeat-to-update-1-description-updated",
      "author": "root-updated",
      "pattern": {
        "name": "test-heartbeat-to-update-1-resource-updated"
      },
      "expected_interval": "1000h",
      "output": "test-heartbeat-to-update-1-output-updated"
    }
    """
    Then the response code should be 403

  Scenario: given invalid update request should return errors
    When I am admin
    When I do PUT /api/v4/heartbeats/test-heartbeat-to-update:
    """
    {}
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "name": "Name is missing.",
        "description": "Description is missing.",
        "author": "Author is missing.",
        "pattern": "Pattern is missing.",
        "expected_interval": "ExpectedInterval is missing."
      }
    }
    """

  Scenario: given update request with already exists name should return error
    When I am admin
    When I do PUT /api/v4/heartbeats/test-heartbeat-to-update:
    """
    {
      "name": "test-heartbeat-to-check-unique-name-name",
      "description": "test-heartbeat-to-update-1-description",
      "author": "root",
      "pattern": {
        "name": "test-heartbeat-to-update-1-resource"
      },
      "expected_interval": "1000h",
      "output": "test-heartbeat-to-update-1-output-updated"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
          "name": "Name already exists"
      }
    }
    """

  Scenario: given update request with not exist id should return not found error
    When I am admin
    When I do PUT /api/v4/heartbeats/test-heartbeat-not-found:
    """
    {
      "name": "test-heartbeat-not-found-name",
      "description": "test-heartbeat-not-found-description",
      "author": "root",
      "pattern": {
        "name": "test-heartbeat-not-found-resource"
      },
      "expected_interval": "1000h",
      "output": "test-heartbeat-not-found-output"
    }
    """
    Then the response code should be 404
    Then the response body should be:
    """
    {
      "error": "Not found"
    }
    """

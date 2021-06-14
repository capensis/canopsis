Feature: Create a heartbeat
  I need to be able to create a heartbeat
  Only admin should be able to create a heartbeat

  Scenario: given create request should return ok
    When I am admin
    When I do POST /api/v4/heartbeats:
    """
    {
      "name": "test-heartbeat-to-create-1-name",
      "description": "test-heartbeat-to-create-1-description",
      "author": "test-heartbeat-to-create-1-author",
      "pattern": {
        "name": "test-heartbeat-to-create-1-resource"
      },
      "expected_interval": "1000h",
      "output": "test-heartbeat-to-create-1-output"
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """
    {
      "name": "test-heartbeat-to-create-1-name",
      "description": "test-heartbeat-to-create-1-description",
      "author": "test-heartbeat-to-create-1-author",
      "pattern": {
        "name": "test-heartbeat-to-create-1-resource"
      },
      "expected_interval": "1000h",
      "output": "test-heartbeat-to-create-1-output"
    }
    """

  Scenario: given create request should return ok to get request
    When I am admin
    When I do POST /api/v4/heartbeats:
    """
    {
      "name": "test-heartbeat-to-create-2-name",
      "description": "test-heartbeat-to-create-2-description",
      "author": "test-heartbeat-to-create-2-author",
      "pattern": {
        "name": "test-heartbeat-to-create-2-resource"
      },
      "expected_interval": "1000h",
      "output": "test-heartbeat-to-create-2-output"
    }
    """
    When I do GET /api/v4/heartbeats/{{ .lastResponse._id}}
    Then the response code should be 200

  Scenario: given create request with custom _id should return ok to get request with that _id
    When I am admin
    When I do POST /api/v4/heartbeats:
    """
    {
      "_id": "custom-heartbeat-id",
      "name": "test-custom-heartbeat-id-name",
      "description": "test-custom-heartbeat-id-description",
      "author": "test-custom-heartbeat-id-author",
      "pattern": {
        "name": "test-custom-heartbeat-id-resource"
      },
      "expected_interval": "1000h",
      "output": "test-custom-heartbeat-id-output"
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/heartbeats/custom-heartbeat-id
    Then the response code should be 200

  Scenario: given create request with custom _id that already exists should return dup error
    When I am admin
    When I do POST /api/v4/heartbeats:
    """
    {
      "_id": "test-heartbeat-to-update",
      "name": "test-custom-heartbeat-id-2--name",
      "description": "test-custom-heartbeat-id-2-description",
      "author": "test-custom-heartbeat-id-2-author",
      "pattern": {
        "name": "test-custom-heartbeat-id-2-resource"
      },
      "expected_interval": "1000h",
      "output": "test-custom-heartbeat-id-2-output"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "_id": "ID already exists"
      }
    }
    """

  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/heartbeats:
    """
    {
      "name": "test-heartbeat-to-create-1-name",
      "description": "test-heartbeat-to-create-1-description",
      "author": "test-heartbeat-to-create-1-author",
      "pattern": {
        "name": "test-heartbeat-to-create-1-resource"
      },
      "expected_interval": "1000h",
      "output": "test-heartbeat-to-create-1-output"
    }
    """
    Then the response code should be 401

  Scenario: given create request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/heartbeats:
    """
    {
      "name": "test-heartbeat-to-create-1-name",
      "description": "test-heartbeat-to-create-1-description",
      "author": "test-heartbeat-to-create-1-author",
      "pattern": {
        "name": "test-heartbeat-to-create-1-resource"
      },
      "expected_interval": "1000h",
      "output": "test-heartbeat-to-create-1-output"
    }
    """
    Then the response code should be 403

  Scenario: given invalid create request should return errors
    When I am admin
    When I do POST /api/v4/heartbeats:
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

  Scenario: given create request with already exists name should return error
    When I am admin
    When I do POST /api/v4/heartbeats:
    """
    {
      "name": "test-heartbeat-to-check-unique-name-name",
      "description": "test-heartbeat-to-create-3-description",
      "author": "test-heartbeat-to-create-3-author",
      "pattern": {
        "name": "test-heartbeat-to-create-3-resource"
      },
      "expected_interval": "1000h",
      "output": "test-heartbeat-to-create-3-output"
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

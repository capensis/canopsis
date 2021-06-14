Feature: Bulk create heartbeats
  I need to be able to create multiple heartbeat
  Only admin should be able to create multiple heartbeat

  Scenario: given bulk create request should return ok
    When I am admin
    When I do POST /api/v4/bulk/heartbeats:
    """
    [
      {
        "name": "test-heartbeat-to-bulk-create-1-name",
        "description": "test-heartbeat-to-bulk-create-1-description",
        "pattern": {
          "name": "test-heartbeat-to-bulk-create-1-resource"
        },
        "expected_interval": "1000h",
        "output": "test-heartbeat-to-bulk-create-1-output"
      },
      {
        "name": "test-heartbeat-to-bulk-create-2-name",
        "description": "test-heartbeat-to-bulk-create-2-description",
        "pattern": {
          "name": "test-heartbeat-to-bulk-create-2-resource"
        },
        "expected_interval": "1000h",
        "output": "test-heartbeat-to-bulk-create-2-output"
      }
    ]
    """
    Then the response code should be 201
    Then the response body should contain:
    """
    [
      {
        "name": "test-heartbeat-to-bulk-create-1-name",
        "description": "test-heartbeat-to-bulk-create-1-description",
        "author": "root",
        "pattern": {
          "name": "test-heartbeat-to-bulk-create-1-resource"
        },
        "expected_interval": "1000h",
        "output": "test-heartbeat-to-bulk-create-1-output"
      },
      {
        "name": "test-heartbeat-to-bulk-create-2-name",
        "description": "test-heartbeat-to-bulk-create-2-description",
        "author": "root",
        "pattern": {
          "name": "test-heartbeat-to-bulk-create-2-resource"
        },
        "expected_interval": "1000h",
        "output": "test-heartbeat-to-bulk-create-2-output"
      }
    ]
    """

  Scenario: given bulk create request should return ok to get request
    When I am admin
    When I do POST /api/v4/bulk/heartbeats:
    """
    [
      {
        "name": "test-heartbeat-to-bulk-create-3-name",
        "description": "test-heartbeat-to-bulk-create-3-description",
        "pattern": {
          "name": "test-heartbeat-to-bulk-create-3-resource"
        },
        "expected_interval": "1000h",
        "output": "test-heartbeat-to-bulk-create-3-output"
      }
    ]
    """
    When I do GET /api/v4/heartbeats/{{ (index .lastResponse 0)._id}}
    Then the response code should be 200

  Scenario: given invalid bulk create request should return errors
    When I am admin
    When I do POST /api/v4/bulk/heartbeats:
    """
    [
      {}
    ]
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "0.name": "Name is missing.",
        "0.description": "Description is missing.",
        "0.pattern": "Pattern is missing.",
        "0.expected_interval": "ExpectedInterval is missing."
      }
    }
    """

  Scenario: given bulk create request with one invalid and one valid data
    should return errors
    When I am admin
    When I do POST /api/v4/bulk/heartbeats:
    """
    [
      {
        "name": "test-heartbeat-to-bulk-create-4-name",
        "description": "test-heartbeat-to-bulk-create-4-description",
        "pattern": {
          "name": "test-heartbeat-to-bulk-create-4-resource"
        },
        "expected_interval": "1000h",
        "output": "test-heartbeat-to-bulk-create-4-output"
      },
      {}
    ]
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "1.name": "Name is missing.",
        "1.description": "Description is missing.",
        "1.pattern": "Pattern is missing.",
        "1.expected_interval": "ExpectedInterval is missing."
      }
    }
    """

  Scenario: given bulk create request with already exists name should return error
    When I am admin
    When I do POST /api/v4/bulk/heartbeats:
    """
    [
      {
        "name": "test-heartbeat-to-check-unique-name-name",
        "description": "test-heartbeat-to-bulk-create-5-description",
        "pattern": {
          "name": "test-heartbeat-to-bulk-create-5-resource"
        },
        "expected_interval": "1000h",
        "output": "test-heartbeat-to-bulk-create-5-output"
      }
    ]
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
          "0.name": "Name already exists."
      }
    }
    """

  Scenario: given bulk create request with multiple items with the same name
    should return error
    When I am admin
    When I do POST /api/v4/bulk/heartbeats:
    """
    [
      {
        "name": "test-heartbeat-to-bulk-create-6-name",
        "description": "test-heartbeat-to-bulk-create-6-description",
        "pattern": {
          "name": "test-heartbeat-to-bulk-create-6-resource"
        },
        "expected_interval": "1000h",
        "output": "test-heartbeat-to-bulk-create-6-output"
      },
      {
        "name": "test-heartbeat-to-bulk-create-6-name",
        "description": "test-heartbeat-to-bulk-create-7-description",
        "pattern": {
          "name": "test-heartbeat-to-bulk-create-7-resource"
        },
        "expected_interval": "1000h",
        "output": "test-heartbeat-to-bulk-create-7-output"
      },
      {
        "name": "test-heartbeat-to-bulk-create-6-name",
        "description": "test-heartbeat-to-bulk-create-8-description",
        "pattern": {
          "name": "test-heartbeat-to-bulk-create-8-resource"
        },
        "expected_interval": "1000h",
        "output": "test-heartbeat-to-bulk-create-8-output"
      }
    ]
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "1.name": "Name already exists.",
        "2.name": "Name already exists."
      }
    }
    """

  Scenario: given bulk create request and no auth user should not allow access
    When I do POST /api/v4/bulk/heartbeats:
    """
    [
      {
        "name": "test-heartbeat-to-bulk-create-1-name",
        "description": "test-heartbeat-to-bulk-create-1-description",
        "pattern": {
          "name": "test-heartbeat-to-bulk-create-1-resource"
        },
        "expected_interval": "1000h",
        "output": "test-heartbeat-to-bulk-create-1-output"
      }
    ]
    """
    Then the response code should be 401

  Scenario: given bulk create request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/bulk/heartbeats:
    """
    [
      {
        "name": "test-heartbeat-to-bulk-create-1-name",
        "description": "test-heartbeat-to-bulk-create-1-description",
        "pattern": {
          "name": "test-heartbeat-to-bulk-create-1-resource"
        },
        "expected_interval": "1000h",
        "output": "test-heartbeat-to-bulk-create-1-output"
      }
    ]
    """
    Then the response code should be 403

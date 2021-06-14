Feature: Bulk update a heartbeats
  I need to be able to update multiple heartbeats
  Only admin should be able to update multiple heartbeat

  Scenario: given bulk update request should update heartbeat
    When I am admin
    Then I do PUT /api/v4/bulk/heartbeats:
    """
    [
      {
        "_id": "test-heartbeat-to-bulk-update-1",
        "name": "test-heartbeat-to-bulk-update-1-name-updated",
        "description": "test-heartbeat-to-bulk-update-1-description-updated",
        "author": "root-updated",
        "pattern": {
          "name": "test-heartbeat-to-bulk-update-1-updated"
        },
        "expected_interval": "1000h",
        "output": "test-heartbeat-to-bulk-update-1-output-updated"
      },
      {
        "_id": "test-heartbeat-to-bulk-update-2",
        "name": "test-heartbeat-to-bulk-update-2-name-updated",
        "description": "test-heartbeat-to-bulk-update-2-description-updated",
        "author": "root-updated",
        "pattern": {
          "name": "test-heartbeat-to-bulk-update-2-updated"
        },
        "expected_interval": "1000h",
        "output": "test-heartbeat-to-bulk-update-2-output-updated"
      }
    ]
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    [
      {
        "_id": "test-heartbeat-to-bulk-update-1",
        "name": "test-heartbeat-to-bulk-update-1-name-updated",
        "description": "test-heartbeat-to-bulk-update-1-description-updated",
        "author": "root-updated",
        "pattern": {
          "name": "test-heartbeat-to-bulk-update-1-updated"
        },
        "expected_interval": "1000h",
        "output": "test-heartbeat-to-bulk-update-1-output-updated",
        "created": 1592215337
      },
      {
        "_id": "test-heartbeat-to-bulk-update-2",
        "name": "test-heartbeat-to-bulk-update-2-name-updated",
        "description": "test-heartbeat-to-bulk-update-2-description-updated",
        "author": "root-updated",
        "pattern": {
          "name": "test-heartbeat-to-bulk-update-2-updated"
        },
        "expected_interval": "1000h",
        "output": "test-heartbeat-to-bulk-update-2-output-updated",
        "created": 1592215337
      }
    ]
    """

  Scenario: given bulk update request with not exist ids should return not found error
    When I am admin
    When I do PUT /api/v4/bulk/heartbeats:
    """
    [
      {
        "_id": "test-heartbeat-not-found-1",
        "name": "test-heartbeat-not-found-1-name",
        "description": "test-heartbeat-not-found-1-description",
        "author": "root",
        "pattern": {
          "name": "test-heartbeat-not-found-1-resource"
        },
        "expected_interval": "1000h",
        "output": "test-heartbeat-not-found-1-output"
      },
      {
        "_id": "test-heartbeat-not-found-2",
        "name": "test-heartbeat-not-found-2-name",
        "description": "test-heartbeat-not-found-2-description",
        "author": "root",
        "pattern": {
          "name": "test-heartbeat-not-found-2-resource"
        },
        "expected_interval": "1000h",
        "output": "test-heartbeat-not-found-2-output"
      }
    ]
    """
    Then the response code should be 404
    Then the response body should be:
    """
    {
      "error": "not found test-heartbeat-not-found-1,test-heartbeat-not-found-2"
    }
    """

  Scenario: given bulk update request with one exist id and one not exist id
    should return not found error
    When I am admin
    When I do PUT /api/v4/bulk/heartbeats:
    """
    [
      {
        "_id": "test-heartbeat-to-bulk-update-1",
        "name": "test-heartbeat-to-bulk-update-1-name-updated",
        "description": "test-heartbeat-to-bulk-update-1-description-updated",
        "author": "root-updated",
        "pattern": {
          "name": "test-heartbeat-to-bulk-update-1-updated"
        },
        "expected_interval": "1000h",
        "output": "test-heartbeat-to-bulk-update-1-output-updated"
      },
      {
        "_id": "test-heartbeat-not-found",
        "name": "test-heartbeat-not-found-name",
        "description": "test-heartbeat-not-found-description",
        "author": "root",
        "pattern": {
          "name": "test-heartbeat-not-found-resource"
        },
        "expected_interval": "1000h",
        "output": "test-heartbeat-not-found-output"
      }
    ]
    """
    Then the response code should be 404
    Then the response body should be:
    """
    {
      "error": "not found test-heartbeat-not-found"
    }
    """

  Scenario: given invalid bulk update request should return errors
    When I am admin
    When I do PUT /api/v4/bulk/heartbeats:
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
        "0._id": "ID is missing.",
        "0.name": "Name is missing.",
        "0.description": "Description is missing.",
        "0.author": "Author is missing.",
        "0.pattern": "Pattern is missing.",
        "0.expected_interval": "ExpectedInterval is missing."
      }
    }
    """

  Scenario: given bulk update request with one valid item and one invalid item
    should return errors
    When I am admin
    When I do PUT /api/v4/bulk/heartbeats:
    """
    [
      {
        "_id": "test-heartbeat-to-bulk-update-1",
        "name": "test-heartbeat-to-bulk-update-1-name-updated",
        "description": "test-heartbeat-to-bulk-update-1-description-updated",
        "author": "root-updated",
        "pattern": {
          "name": "test-heartbeat-to-bulk-update-1-updated"
        },
        "expected_interval": "1000h",
        "output": "test-heartbeat-to-bulk-update-1-output-updated"
      },
      {}
    ]
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "1._id": "ID is missing.",
        "1.name": "Name is missing.",
        "1.description": "Description is missing.",
        "1.author": "Author is missing.",
        "1.pattern": "Pattern is missing.",
        "1.expected_interval": "ExpectedInterval is missing."
      }
    }
    """

  Scenario: given bulk update request with already exists name should return error
    When I am admin
    When I do PUT /api/v4/bulk/heartbeats:
    """
    [
      {
        "_id": "test-heartbeat-to-bulk-update-1",
        "name": "test-heartbeat-to-check-unique-name-name",
        "description": "test-heartbeat-to-bulk-update-1-description",
        "author": "root",
        "pattern": {
          "name": "test-heartbeat-to-update-resource"
        },
        "expected_interval": "1000h",
        "output": "test-heartbeat-to-bulk-update-1-output"
      }
    ]
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
          "0.name": "Name already exists"
      }
    }
    """

  Scenario: given bulk update request with multiple items with the same name
    should return error
    When I am admin
    Then I do PUT /api/v4/bulk/heartbeats:
    """
    [
      {
        "_id": "test-heartbeat-to-bulk-update-1",
        "name": "test-heartbeat-to-bulk-update-1-name-same-updated",
        "description": "test-heartbeat-to-bulk-update-1-description-updated",
        "author": "root-updated",
        "pattern": {
          "name": "test-heartbeat-to-bulk-update-1-updated"
        },
        "expected_interval": "1000h",
        "output": "test-heartbeat-to-bulk-update-1-output-updated"
      },
      {
        "_id": "test-heartbeat-to-bulk-update-2",
        "name": "test-heartbeat-to-bulk-update-1-name-same-updated",
        "description": "test-heartbeat-to-bulk-update-2-description-updated",
        "author": "root-updated",
        "pattern": {
          "name": "test-heartbeat-to-bulk-update-2-updated"
        },
        "expected_interval": "1000h",
        "output": "test-heartbeat-to-bulk-update-2-output-updated"
      },
      {
        "_id": "test-heartbeat-to-bulk-update-3",
        "name": "test-heartbeat-to-bulk-update-1-name-same-updated",
        "description": "test-heartbeat-to-bulk-update-3-description-updated",
        "author": "root-updated",
        "pattern": {
          "name": "test-heartbeat-to-bulk-update-3-updated"
        },
        "expected_interval": "1000h",
        "output": "test-heartbeat-to-bulk-update-3-output-updated"
      }
    ]
    """
    Then the response code should be 400
    """
    {
      "errors": {
          "1.name": "Name already exists",
          "2.name": "Name already exists"
      }
    }
    """

  Scenario: given get request and no auth user should not allow access
    When I do PUT /api/v4/bulk/heartbeats:
    """
    [
      {
        "_id": "test-heartbeat-to-bulk-update-1",
        "name": "test-heartbeat-to-bulk-update-1-name-updated",
        "description": "test-heartbeat-to-bulk-update-1-description-updated",
        "author": "root-updated",
        "pattern": {
          "name": "test-heartbeat-to-bulk-update-1-updated"
        },
        "expected_interval": "1000h",
        "output": "test-heartbeat-to-bulk-update-1-output-updated"
      }
    ]
    """
    Then the response code should be 401

  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/bulk/heartbeats:
    """
    [
      {
        "_id": "test-heartbeat-to-bulk-update-1",
        "name": "test-heartbeat-to-bulk-update-1-name-updated",
        "description": "test-heartbeat-to-bulk-update-1-description-updated",
        "author": "root-updated",
        "pattern": {
          "name": "test-heartbeat-to-bulk-update-1-updated"
        },
        "expected_interval": "1000h",
        "output": "test-heartbeat-to-bulk-update-1-output-updated"
      }
    ]
    """
    Then the response code should be 403
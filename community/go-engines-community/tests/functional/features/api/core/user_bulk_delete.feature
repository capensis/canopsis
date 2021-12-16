Feature: Delete a user
  I need to be able to bulk delete a user
  Only admin should be able to bulk delete a user

  Scenario: given delete request and no auth user should not allow access
    When I do DELETE /api/v4/bulk/users:
    """json
    [
      {
        "_id": "test-user-to-bulk-delete-1"
      },
      {
        "_id": "test-user-to-bulk-delete-2"
      }
    ]
    """
    Then the response code should be 401

  Scenario: given delete request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do DELETE /api/v4/bulk/users:
    """json
    [
      {
        "_id": "test-user-to-bulk-delete-1"
      },
      {
        "_id": "test-user-to-bulk-delete-2"
      }
    ]
    """
    Then the response code should be 403

  Scenario: given delete request should delete user
    When I am admin
    When I do DELETE /api/v4/bulk/users:
    """json
    [
      {
        "_id": "test-user-to-bulk-delete-1"
      },
      {
        "_id": "test-user-to-bulk-delete-2"
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "id": "test-user-to-bulk-delete-1",
        "status": 200,
        "item": {
          "_id": "test-user-to-bulk-delete-1"
        }
      },
      {
        "id": "test-user-to-bulk-delete-2",
        "status": 200,
        "item": {
          "_id": "test-user-to-bulk-delete-2"
        }
      }
    ]
    """
    When I do GET /api/v4/users/test-user-to-bulk-delete-1
    Then the response code should be 404
    When I do GET /api/v4/users/test-user-to-bulk-delete-2
    Then the response code should be 404
    When I do DELETE /api/v4/bulk/users:
    """json
    [
      {
        "_id": "test-user-to-bulk-delete-1"
      },
      {
        "_id": "test-user-to-bulk-delete-2"
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "error": "Not found",
        "status": 404,
        "item": {
          "_id": "test-user-to-bulk-delete-1"
        }
      },
      {
        "error": "Not found",
        "status": 404,
        "item": {
          "_id": "test-user-to-bulk-delete-2"
        }
      }
    ]
    """

  Scenario: given delete request with empty ids should return error
    When I am admin
    When I do DELETE /api/v4/bulk/users:
    """json
    [
      {}
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "errors": {
          "_id": "ID is missing."
        },
        "item": {},
        "status": 400
      }
    ]
    """

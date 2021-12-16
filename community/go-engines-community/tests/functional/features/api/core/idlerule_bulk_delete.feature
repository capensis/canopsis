Feature: Delete a idlerule
  I need to be able to bulk delete a idlerule
  Only admin should be able to bulk delete a idlerule

  Scenario: given delete request and no auth idlerule should not allow access
    When I do DELETE /api/v4/bulk/idle-rules:
    """json
    [
      {
        "_id": "test-idle-rule-to-bulk-delete-1"
      },
      {
        "_id": "test-idle-rule-to-bulk-delete-2"
      }
    ]
    """
    Then the response code should be 401

  Scenario: given delete request and auth idlerule by api key without permissions should not allow access
    When I am noperms
    When I do DELETE /api/v4/bulk/idle-rules:
    """json
    [
      {
        "_id": "test-idle-rule-to-bulk-delete-1"
      },
      {
        "_id": "test-idle-rule-to-bulk-delete-2"
      }
    ]
    """
    Then the response code should be 403

  Scenario: given delete request should delete idlerule
    When I am admin
    When I do DELETE /api/v4/bulk/idle-rules:
    """json
    [
      {
        "_id": "test-idle-rule-to-bulk-delete-1"
      },
      {
        "_id": "test-idle-rule-to-bulk-delete-2"
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "id": "test-idle-rule-to-bulk-delete-1",
        "status": 200,
        "item": {
          "_id": "test-idle-rule-to-bulk-delete-1"
        }
      },
      {
        "id": "test-idle-rule-to-bulk-delete-2",
        "status": 200,
        "item": {
          "_id": "test-idle-rule-to-bulk-delete-2"
        }
      }
    ]
    """
    When I do GET /api/v4/idle-rules/test-idle-rule-to-bulk-delete-1
    Then the response code should be 404
    When I do GET /api/v4/idle-rules/test-idle-rule-to-bulk-delete-2
    Then the response code should be 404
    When I do DELETE /api/v4/bulk/idle-rules:
    """json
    [
      {
        "_id": "test-idle-rule-to-bulk-delete-1"
      },
      {
        "_id": "test-idle-rule-to-bulk-delete-2"
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
          "_id": "test-idle-rule-to-bulk-delete-1"
        }
      },
      {
        "error": "Not found",
        "status": 404,
        "item": {
          "_id": "test-idle-rule-to-bulk-delete-2"
        }
      }
    ]
    """

  Scenario: given delete request with empty ids should return error
    When I am admin
    When I do DELETE /api/v4/bulk/idle-rules:
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

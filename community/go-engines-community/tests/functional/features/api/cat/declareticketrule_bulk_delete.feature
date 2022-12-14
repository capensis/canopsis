Feature: Bulk delete declare ticket rules
  I need to be able to bulk delete declare ticket rules
  Only admin should be able to bulk delete declare ticket rules

  Scenario: given delete request should return multi status and should be handled independently
    When I am admin
    When I do DELETE /api/v4/cat/bulk/declare-ticket-rules:
    """json
    [
      {
        "_id": "test-declare-ticket-rule-to-bulk-delete-1"
      },
      {
        "_id": "test-declare-ticket-rule-to-bulk-delete-not-found"
      },
      {},
      [],
      {
        "_id": "test-declare-ticket-rule-to-bulk-delete-1"
      }
    ]
    """
    Then the response code should be 207
    Then the response body should be:
    """json
    [
      {
        "id": "test-declare-ticket-rule-to-bulk-delete-1",
        "status": 200,
        "item": {
          "_id": "test-declare-ticket-rule-to-bulk-delete-1"
        }
      },
      {
        "error": "Not found",
        "status": 404,
        "item": {
          "_id": "test-declare-ticket-rule-to-bulk-delete-not-found"
        }
      },
      {
        "errors": {
          "_id": "ID is missing."
        },
        "item": {},
        "status": 400
      },
      {
        "status": 400,
        "item": [],
        "error": "value doesn't contain object; it contains array"
      },
      {
        "error": "Not found",
        "status": 404,
        "item": {
          "_id": "test-declare-ticket-rule-to-bulk-delete-1"
        }
      }
    ]
    """
    When I do GET /api/v4/declare-ticket-rules/test-declare-ticket-rule-to-bulk-delete-1
    Then the response code should be 404

  Scenario: given bulk delete request and no auth user should not allow access
    When I do DELETE /api/v4/cat/bulk/declare-ticket-rules
    Then the response code should be 401

  Scenario: given bulk delete request and auth user without permissions should not allow access
    When I am noperms
    When I do DELETE /api/v4/cat/bulk/declare-ticket-rules
    Then the response code should be 403

Feature: Bulk delete patterns
  I need to be able to bulk delete patterns
  Only admin should be able to bulk delete patterns

  Scenario: given bulk bulk delete request should return multistatus and should be handled independently
    When I am noperms
    When I do DELETE /api/v4/bulk/patterns:
    """json
    [
      {
        "_id": "test-pattern-to-bulk-delete-1"
      },
      {
        "_id": "test-pattern-to-bulk-delete-2"
      },
      {
        "_id": "test-pattern-to-bulk-delete-3"
      },
      {
        "_id": "test-pattern-to-bulk-delete-not-found"
      },
      {},
      [],
      {
        "_id": "test-pattern-to-bulk-delete-2"
      },
      {
        "_id": "test-pattern-to-bulk-delete-corporate-1"
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "id": "test-pattern-to-bulk-delete-1",
        "status": 200,
        "item": {
          "_id": "test-pattern-to-bulk-delete-1"
        }
      },
      {
        "id": "test-pattern-to-bulk-delete-2",
        "status": 200,
        "item": {
          "_id": "test-pattern-to-bulk-delete-2"
        }
      },
      {
        "error": "Not found",
        "status": 404,
        "item": {
          "_id": "test-pattern-to-bulk-delete-3"
        }
      },
      {
        "error": "Not found",
        "status": 404,
        "item": {
          "_id": "test-pattern-to-bulk-delete-not-found"
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
          "_id": "test-pattern-to-bulk-delete-2"
        }
      },
      {
        "error": "Forbidden",
        "status": 403,
        "item": {
          "_id": "test-pattern-to-bulk-delete-corporate-1"
        }
      }
    ]
    """
    When I do GET /api/v4/patterns/test-pattern-to-bulk-delete-1
    Then the response code should be 404
    When I do GET /api/v4/patterns/test-pattern-to-bulk-delete-2
    Then the response code should be 404
    When I am admin
    When I do DELETE /api/v4/bulk/patterns:
    """json
    [
      {
        "_id": "test-pattern-to-bulk-delete-4"
      },
      {
        "_id": "test-pattern-to-bulk-delete-corporate-1"
      },
      {
        "_id": "test-pattern-to-bulk-delete-corporate-2"
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
          "_id": "test-pattern-to-bulk-delete-4"
        }
      },
      {
        "id": "test-pattern-to-bulk-delete-corporate-1",
        "status": 200,
        "item": {
          "_id": "test-pattern-to-bulk-delete-corporate-1"
        }
      },
      {
        "id": "test-pattern-to-bulk-delete-corporate-2",
        "status": 200,
        "item": {
          "_id": "test-pattern-to-bulk-delete-corporate-2"
        }
      }
    ]
    """
    When I do GET /api/v4/patterns/test-pattern-to-bulk-delete-corporate-1
    Then the response code should be 404
    When I do GET /api/v4/patterns/test-pattern-to-bulk-delete-corporate-2
    Then the response code should be 404

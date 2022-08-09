Feature: Bulk delete eventfilters
  I need to be able to bulk delete eventfilters
  Only admin should be able to bulk delete eventfilters

  Scenario: given bulk delete request and no auth eventfilter should not allow access
    When I do DELETE /api/v4/bulk/eventfilters
    Then the response code should be 401

  Scenario: given bulk delete request and auth by api key without permissions should not allow access
    When I am noperms
    When I do DELETE /api/v4/bulk/eventfilters
    Then the response code should be 403

  Scenario: given bulk delete request should return multi status and should be handled independently
    When I am admin
    When I do DELETE /api/v4/bulk/eventfilters:
    """json
    [
      {
        "_id": "test-eventfilter-to-bulk-delete-1"
      },
      {
        "_id": "test-eventfilter-to-bulk-delete-1"
      },
      {
        "_id": "test-eventfilter-to-bulk-delete-not-found"
      },
      {},
      [],
      {
        "_id": "test-eventfilter-to-bulk-delete-2"
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "id": "test-eventfilter-to-bulk-delete-1",
        "status": 200,
        "item": {
          "_id": "test-eventfilter-to-bulk-delete-1"
        }
      },
      {
        "error": "Not found",
        "status": 404,
        "item": {
          "_id": "test-eventfilter-to-bulk-delete-1"
        }
      },
      {
        "error": "Not found",
        "status": 404,
        "item": {
          "_id": "test-eventfilter-to-bulk-delete-not-found"
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
        "id": "test-eventfilter-to-bulk-delete-2",
        "status": 200,
        "item": {
          "_id": "test-eventfilter-to-bulk-delete-2"
        }
      }
    ]
    """
    When I do GET /api/v4/eventfilter/rules/test-eventfilter-to-bulk-delete-1
    Then the response code should be 404
    When I do GET /api/v4/eventfilter/rules/test-eventfilter-to-bulk-delete-2
    Then the response code should be 404

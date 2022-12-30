Feature: Bulk delete pbehaviors
  I need to be able to bulk delete pbehaviors
  Only admin should be able to bulk delete pbehaviors

  Scenario: given bulk delete request and no auth pbehavior should not allow access
    When I do DELETE /api/v4/bulk/entity-pbehaviors
    Then the response code should be 401

  Scenario: given bulk delete request and auth pbehavior by api key without permissions should not allow access
    When I am noperms
    When I do DELETE /api/v4/bulk/entity-pbehaviors
    Then the response code should be 403

  Scenario: given bulk delete request should return multi status and should be handled independently
    When I am admin
    When I do DELETE /api/v4/bulk/entity-pbehaviors:
    """json
    [
      {
        "entity": "test-entity-to-bulk-entity-pbehavior-delete-1/test-component-default",
        "origin": "test-pbehavior-to-bulk-entity-delete-1-origin"
      },
      {
        "entity": "test-entity-to-bulk-entity-pbehavior-delete-1/test-component-default",
        "origin": "test-pbehavior-to-bulk-entity-delete-1-origin"
      },
      {
        "entity": "test-entity-to-bulk-entity-pbehavior-delete-not-found/test-component-default",
        "origin": "test-pbehavior-to-bulk-entity-delete-1-origin"
      },
      {},
      [],
      {
        "entity": "test-entity-to-bulk-entity-pbehavior-delete-2/test-component-default",
        "origin": "test-pbehavior-to-bulk-entity-delete-2-origin"
      },
      {
        "entity": "test-entity-to-bulk-entity-pbehavior-delete-3/test-component-default",
        "origin": "test-pbehavior-to-bulk-entity-delete-3-origin"
      }
    ]
    """
    Then the response code should be 207
    Then the response body should be:
    """json
    [
      {
        "id": "test-pbehavior-to-bulk-entity-delete-1",
        "status": 200,
        "item": {
          "entity": "test-entity-to-bulk-entity-pbehavior-delete-1/test-component-default",
          "origin": "test-pbehavior-to-bulk-entity-delete-1-origin"
        }
      },
      {
        "error": "Not found",
        "status": 404,
        "item": {
          "entity": "test-entity-to-bulk-entity-pbehavior-delete-1/test-component-default",
          "origin": "test-pbehavior-to-bulk-entity-delete-1-origin"
        }
      },
      {
        "error": "Not found",
        "status": 404,
        "item": {
          "entity": "test-entity-to-bulk-entity-pbehavior-delete-not-found/test-component-default",
          "origin": "test-pbehavior-to-bulk-entity-delete-1-origin"
        }
      },
      {
        "errors": {
          "entity": "Entity is missing.",
          "origin": "Origin is missing."
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
        "id": "test-pbehavior-to-bulk-entity-delete-2",
        "status": 200,
        "item": {
          "entity": "test-entity-to-bulk-entity-pbehavior-delete-2/test-component-default",
          "origin": "test-pbehavior-to-bulk-entity-delete-2-origin"
        }
      },
      {
        "error": "Not found",
        "status": 404,
        "item": {
          "entity": "test-entity-to-bulk-entity-pbehavior-delete-3/test-component-default",
          "origin": "test-pbehavior-to-bulk-entity-delete-3-origin"
        }
      }
    ]
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-bulk-entity-delete-1
    Then the response code should be 404
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-bulk-entity-delete-2
    Then the response code should be 404

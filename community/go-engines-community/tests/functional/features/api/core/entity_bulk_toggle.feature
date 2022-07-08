Feature: Bulk toggle entityservices
  I need to be able to bulk toggle entityservices
  Only admin should be able to bulk toggle entityservices

  Scenario: given bulk update request and no auth should not allow access
    When I do PUT /api/v4/bulk/entities/enable
    Then the response code should be 401
    When I do PUT /api/v4/bulk/entities/disable
    Then the response code should be 401

  Scenario: given bulk update request and auth by api key without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/bulk/entities/enable
    Then the response code should be 403
    When I do PUT /api/v4/bulk/entities/disable
    Then the response code should be 403

  Scenario: given bulk update request should return multistatus and should be handled independently
    When I am admin
    When I do PUT /api/v4/bulk/entities/disable:
    """json
    [
      {
        "_id": "test-entity-to-bulk-toggle-resource-1/test-entity-to-bulk-toggle-component-1"
      },
      {
        "_id": "test-entity-to-bulk-toggle-resource-2/test-entity-to-bulk-toggle-component-1"
      },
      {},
      {
        "_id": "test-entity-to-bulk-toggle-resource-3/test-entity-to-bulk-toggle-component-1"
      },
      {
        "_id": "test-entity-to-bulk-toggle-service"
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "id": "test-entity-to-bulk-toggle-resource-1/test-entity-to-bulk-toggle-component-1",
        "status": 200,
        "item": {
          "_id": "test-entity-to-bulk-toggle-resource-1/test-entity-to-bulk-toggle-component-1"
        }
      },
      {
        "id": "test-entity-to-bulk-toggle-resource-2/test-entity-to-bulk-toggle-component-1",
        "status": 200,
        "item": {
          "_id": "test-entity-to-bulk-toggle-resource-2/test-entity-to-bulk-toggle-component-1"
        }
      },
      {
        "status": 400,
        "item": {},
        "errors": {
          "_id": "ID is missing."
        }
      },
      {
        "status": 404,
        "item": {
          "_id": "test-entity-to-bulk-toggle-resource-3/test-entity-to-bulk-toggle-component-1"
        },
        "error": "Not found"
      },
      {
        "id": "test-entity-to-bulk-toggle-service",
        "status": 200,
        "item": {
          "_id": "test-entity-to-bulk-toggle-service"
        }
      }
    ]
    """
    When I do GET /api/v4/entities?search=test-entity-to-bulk-toggle&sort_by=impact_level&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-entity-to-bulk-toggle-connector/test-entity-to-bulk-toggle-connector-name",
          "enabled": true
        },
        {
          "_id": "test-entity-to-bulk-toggle-component-1",
          "enabled": true
        },
        {
          "_id": "test-entity-to-bulk-toggle-resource-1/test-entity-to-bulk-toggle-component-1",
          "enabled": false
        },
        {
          "_id": "test-entity-to-bulk-toggle-resource-2/test-entity-to-bulk-toggle-component-1",
          "enabled": false
        },
        {
          "_id": "test-entity-to-bulk-toggle-service",
          "enabled": false
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 5
      }
    }
    """
    When I do PUT /api/v4/bulk/entities/enable:
    """json
    [
      {
        "_id": "test-entity-to-bulk-toggle-resource-1/test-entity-to-bulk-toggle-component-1"
      },
      {
        "_id": "test-entity-to-bulk-toggle-resource-2/test-entity-to-bulk-toggle-component-1"
      },
      {},
      {
        "_id": "test-entity-to-bulk-toggle-resource-3/test-entity-to-bulk-toggle-component-1"
      },
      {
        "_id": "test-entity-to-bulk-toggle-service"
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "id": "test-entity-to-bulk-toggle-resource-1/test-entity-to-bulk-toggle-component-1",
        "status": 200,
        "item": {
          "_id": "test-entity-to-bulk-toggle-resource-1/test-entity-to-bulk-toggle-component-1"
        }
      },
      {
        "id": "test-entity-to-bulk-toggle-resource-2/test-entity-to-bulk-toggle-component-1",
        "status": 200,
        "item": {
          "_id": "test-entity-to-bulk-toggle-resource-2/test-entity-to-bulk-toggle-component-1"
        }
      },
      {
        "status": 400,
        "item": {},
        "errors": {
          "_id": "ID is missing."
        }
      },
      {
        "status": 404,
        "item": {
          "_id": "test-entity-to-bulk-toggle-resource-3/test-entity-to-bulk-toggle-component-1"
        },
        "error": "Not found"
      },
      {
        "id": "test-entity-to-bulk-toggle-service",
        "status": 200,
        "item": {
          "_id": "test-entity-to-bulk-toggle-service"
        }
      }      
    ]
    """
    When I do GET /api/v4/entities?search=test-entity-to-bulk-toggle&sort_by=impact_level&sort=desc
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-entity-to-bulk-toggle-connector/test-entity-to-bulk-toggle-connector-name",
          "enabled": true
        },
        {
          "_id": "test-entity-to-bulk-toggle-component-1",
          "enabled": true
        },
        {
          "_id": "test-entity-to-bulk-toggle-resource-1/test-entity-to-bulk-toggle-component-1",
          "enabled": true
        },
        {
          "_id": "test-entity-to-bulk-toggle-resource-2/test-entity-to-bulk-toggle-component-1",
          "enabled": true
        },
        {
          "_id": "test-entity-to-bulk-toggle-service",
          "enabled": true
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 5
      }
    }
    """

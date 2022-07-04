Feature: Bulk update entities
  I need to be able to bulk update entities
  Only admin should be able to bulk update entities

  Scenario: given bulk update request and no auth should not allow access
    When I do PUT /api/v4/bulk/entitybasics
    Then the response code should be 401

  Scenario: given bulk update request and auth by api key without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/bulk/entitybasics
    Then the response code should be 403

  Scenario: given bulk update request should return multistatus and should be handled independently
    When I am admin
    When I do PUT /api/v4/bulk/entitybasics:
    """json
    [
      {
        "_id": "test-entitybasic-to-bulk-update-resource-1/test-entitybasic-to-bulk-update-component-1",
        "description": "test-entitybasic-to-bulk-update-resource-description-updated",
        "enabled": true,
        "category": "test-category-to-entitybasic-edit",
        "impact_level": 1,
        "sli_avail_state": 1,
        "impact": [
          "test-entitybasic-to-bulk-update-connector/test-entitybasic-to-bulk-update-connector-name"
        ],
        "depends": [
          "test-entitybasic-to-bulk-update-resource-1/test-entitybasic-to-bulk-update-component-1"
        ]
      },
      {
        "_id": "test-entitybasic-to-bulk-update-resource-1/test-entitybasic-to-bulk-update-component-1",
        "description": "test-entitybasic-to-bulk-update-resource-description-updated-twice",
        "enabled": true,
        "category": "test-category-to-entitybasic-edit",
        "impact_level": 1,
        "sli_avail_state": 1,
        "impact": [
          "test-entitybasic-to-bulk-update-connector/test-entitybasic-to-bulk-update-connector-name"
        ],
        "depends": [
          "test-entitybasic-to-bulk-update-resource-1/test-entitybasic-to-bulk-update-component-1"
        ]
      },
      {},
      {
        "_id": "test-entitybasic-to-bulk-update-resource-1/test-entitybasic-to-bulk-update-component-1",
        "description": "test-entitybasic-to-update-description",
        "enabled": true,
        "category": "test-category-not-exist",
        "impact_level": 11,
        "sli_avail_state": 4,
        "infos": [
          {}
        ],
        "impact": [
          "test-entity-not-exist"
        ],
        "depends": [
          "test-entity-not-exist"
        ]
      },
      {
        "_id": "test-entitybasic-to-bulk-update-resource-1/test-entitybasic-to-bulk-update-component-1",
        "impact": [
          "test-entitybasic-to-edit-impact-1",
          "test-entitybasic-to-edit-impact-1"
        ],
        "depends": [
          "test-entitybasic-to-edit-impact-2",
          "test-entitybasic-to-edit-impact-2"
        ]
      },
      {
        "_id": "test-entitybasic-to-bulk-update-resource-1/test-entitybasic-to-bulk-update-component-1",
        "impact": [
          "test-entitybasic-to-edit-impact-1"
        ],
        "depends": [
          "test-entitybasic-to-edit-impact-1"
        ]
      },
      {
        "_id": "test-entitybasic-to-bulk-update-resource-not-found/test-entitybasic-to-bulk-update-component-not-found",
        "description": "test-entitybasic-to-bulk-update-resource-description-updated",
        "enabled": true,
        "category": "test-category-to-entitybasic-edit",
        "impact_level": 1,
        "sli_avail_state": 1,
        "impact": [
          "test-entitybasic-to-bulk-update-connector/test-entitybasic-to-bulk-update-connector-name"
        ],
        "depends": [
          "test-entitybasic-to-bulk-update-resource-1/test-entitybasic-to-bulk-update-component-1"
        ]
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "id": "test-entitybasic-to-bulk-update-resource-1/test-entitybasic-to-bulk-update-component-1",
        "status": 200,
        "item": {
          "_id": "test-entitybasic-to-bulk-update-resource-1/test-entitybasic-to-bulk-update-component-1",
          "description": "test-entitybasic-to-bulk-update-resource-description-updated",
          "enabled": true,
          "category": "test-category-to-entitybasic-edit",
          "impact_level": 1,
          "sli_avail_state": 1,
          "impact": [
            "test-entitybasic-to-bulk-update-connector/test-entitybasic-to-bulk-update-connector-name"
          ],
          "depends": [
            "test-entitybasic-to-bulk-update-resource-1/test-entitybasic-to-bulk-update-component-1"
          ]
        }
      },
      {
        "id": "test-entitybasic-to-bulk-update-resource-1/test-entitybasic-to-bulk-update-component-1",
        "status": 200,
        "item": {
          "_id": "test-entitybasic-to-bulk-update-resource-1/test-entitybasic-to-bulk-update-component-1",
          "description": "test-entitybasic-to-bulk-update-resource-description-updated-twice",
          "enabled": true,
          "category": "test-category-to-entitybasic-edit",
          "impact_level": 1,
          "sli_avail_state": 1,
          "impact": [
            "test-entitybasic-to-bulk-update-connector/test-entitybasic-to-bulk-update-connector-name"
          ],
          "depends": [
            "test-entitybasic-to-bulk-update-resource-1/test-entitybasic-to-bulk-update-component-1"
          ]
        }
      },
      {
        "status": 400,
        "item": {},
        "errors": {
          "enabled": "Enabled is missing.",
          "impact_level": "ImpactLevel is missing.",
          "sli_avail_state": "SliAvailState is missing."
        }
      },
      {
        "status": 400,
        "item": {
          "_id": "test-entitybasic-to-bulk-update-resource-1/test-entitybasic-to-bulk-update-component-1",
          "description": "test-entitybasic-to-update-description",
          "enabled": true,
          "category": "test-category-not-exist",
          "impact_level": 11,
          "sli_avail_state": 4,
          "infos": [
            {}
          ],
          "impact": [
            "test-entity-not-exist"
          ],
          "depends": [
            "test-entity-not-exist"
          ]
        },
        "errors": {
          "impact_level": "ImpactLevel should be 10 or less.",
          "sli_avail_state": "SliAvailState should be 3 or less.",
          "category": "Category doesn't exist.",
          "impact": "Impacts doesn't exist.",
          "depends": "Depends doesn't exist.",
          "infos.0.name": "Name is missing."
        }
      },
      {
        "status": 400,
        "item": {
          "_id": "test-entitybasic-to-bulk-update-resource-1/test-entitybasic-to-bulk-update-component-1",
          "impact": [
            "test-entitybasic-to-edit-impact-1",
            "test-entitybasic-to-edit-impact-1"
          ],
          "depends": [
            "test-entitybasic-to-edit-impact-2",
            "test-entitybasic-to-edit-impact-2"
          ]
        },
        "errors": {
          "depends": "Depends contains duplicate values.",
          "impact": "Impacts contains duplicate values."
        }
      },
      {
        "status": 400,
        "item": {
          "impact": [
            "test-entitybasic-to-edit-impact-1"
          ],
          "depends": [
            "test-entitybasic-to-edit-impact-1"
          ]
        },
        "errors": {
          "depends": "Depends contains duplicate values with Impacts."
        }
      },
      {
        "status": 404,
        "item": {
          "_id": "test-entitybasic-to-bulk-update-resource-not-found/test-entitybasic-to-bulk-update-component-not-found",
          "description": "test-entitybasic-to-bulk-update-resource-description-updated",
          "enabled": true,
          "category": "test-category-to-entitybasic-edit",
          "impact_level": 1,
          "sli_avail_state": 1,
          "impact": [
            "test-entitybasic-to-bulk-update-connector/test-entitybasic-to-bulk-update-connector-name"
          ],
          "depends": [
            "test-entitybasic-to-bulk-update-resource-1/test-entitybasic-to-bulk-update-component-1"
          ]
        },
        "error": "Not found"
      }
    ]
    """
    When I do GET /api/v4/entitybasics?_id=test-entitybasic-to-bulk-update-resource-1/test-entitybasic-to-bulk-update-component-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-entitybasic-to-bulk-update-resource-1/test-entitybasic-to-bulk-update-component-1",
      "name": "test-entitybasic-to-bulk-update-resource-1",
      "enable_history": [],
      "measurements": null,
      "enabled": true,
      "infos": {},
      "type": "resource",
      "impact_level": 1,
      "category": {
        "_id": "test-category-to-entitybasic-edit",
        "name": "test-category-to-entitybasic-edit-name",
        "author": "test-category-to-entitybasic-edit-author",
        "created": 1592215337,
        "updated": 1592215337
      },
      "description": "test-entitybasic-to-bulk-update-resource-description-updated-twice",
      "sli_avail_state": 1
    }
    """

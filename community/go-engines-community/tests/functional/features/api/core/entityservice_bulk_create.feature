Feature: Bulk create entityservices
  I need to be able to bulk create entityservices
  Only admin should be able to bulk create entityservices

  Scenario: given bulk create request and no auth should not allow access
    When I do POST /api/v4/bulk/entityservices
    Then the response code should be 401

  Scenario: given bulk create request and auth by api key without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/bulk/entityservices
    Then the response code should be 403

  Scenario: given bulk create request should return multistatus and should be handled independently
    When I am admin
    When I do POST /api/v4/bulk/entityservices:
    """json
    [
      {
        "_id": "test-entityservice-to-bulk-create-1",
        "name": "test-entityservice-to-bulk-create-1-name",
        "output_template": "test-entityservice-to-bulk-create-1-output",
        "category": "test-category-to-entityservice-edit",
        "impact_level": 1,
        "enabled": true,
        "entity_pattern": [
          [
            {
              "field": "name",
              "cond": {
                "type": "eq",
                "value": "test-entityservice-to-bulk-create-1-pattern"
              }
            }
          ]
        ],
        "sli_avail_state": 1,
        "infos": [
          {
            "description": "test-entityservice-to-bulk-create-info-1-description",
            "name": "test-entityservice-to-bulk-create-info-1-name",
            "value": "test-entityservice-to-bulk-create-info-1-value"
          },
          {
            "description": "test-entityservice-to-bulk-create-info-2-description",
            "name": "test-entityservice-to-bulk-create-info-2-name",
            "value": false
          },
          {
            "description": "test-entityservice-to-bulk-create-info-3-description",
            "name": "test-entityservice-to-bulk-create-info-3-name",
            "value": 1022
          },
          {
            "description": "test-entityservice-to-bulk-create-info-4-description",
            "name": "test-entityservice-to-bulk-create-info-4-name",
            "value": 10.45
          },
          {
            "description": "test-entityservice-to-bulk-create-info-5-description",
            "name": "test-entityservice-to-bulk-create-info-5-name",
            "value": null
          },
          {
            "description": "test-entityservice-to-bulk-create-info-6-description",
            "name": "test-entityservice-to-bulk-create-info-6-name",
            "value": ["test-entityservice-to-bulk-create-info-6-value", false, 1022, 10.45, null]
          }
        ]
      },
      {
        "_id": "test-entityservice-to-bulk-create-1",
        "name": "test-entityservice-to-bulk-create-1-name",
        "output_template": "test-entityservice-to-bulk-create-1-output",
        "category": "test-category-to-entityservice-edit",
        "impact_level": 1,
        "enabled": true,
        "entity_pattern": [
          [
            {
              "field": "name",
              "cond": {
                "type": "eq",
                "value": "test-entityservice-to-bulk-create-1-pattern"
              }
            }
          ]
        ],
        "sli_avail_state": 1,
        "infos": [
          {
            "description": "test-entityservice-to-bulk-create-info-1-description",
            "name": "test-entityservice-to-bulk-create-info-1-name",
            "value": "test-entityservice-to-bulk-create-info-1-value"
          },
          {
            "description": "test-entityservice-to-bulk-create-info-2-description",
            "name": "test-entityservice-to-bulk-create-info-2-name",
            "value": false
          },
          {
            "description": "test-entityservice-to-bulk-create-info-3-description",
            "name": "test-entityservice-to-bulk-create-info-3-name",
            "value": 1022
          },
          {
            "description": "test-entityservice-to-bulk-create-info-4-description",
            "name": "test-entityservice-to-bulk-create-info-4-name",
            "value": 10.45
          },
          {
            "description": "test-entityservice-to-bulk-create-info-5-description",
            "name": "test-entityservice-to-bulk-create-info-5-name",
            "value": null
          },
          {
            "description": "test-entityservice-to-bulk-create-info-6-description",
            "name": "test-entityservice-to-bulk-create-info-6-name",
            "value": ["test-entityservice-to-bulk-create-info-6-value", false, 1022, 10.45, null]
          }
        ]
      },
      {},
      {
        "category": "test-category-not-exist",
        "infos": [
          {}
        ],
        "sli_avail_state": 4
      },
      {
        "name": "test-entityservice-to-check-unique-name-name"
      },
      [],
      {
        "_id": "test-entityservice-to-bulk-create-2",
        "name": "test-entityservice-to-bulk-create-2-name",
        "output_template": "test-entityservice-to-bulk-create-2-output",
        "category": "test-category-to-entityservice-edit",
        "impact_level": 1,
        "enabled": true,
        "entity_pattern": [
          [
            {
              "field": "name",
              "cond": {
                "type": "eq",
                "value": "test-entityservice-to-bulk-create-2-pattern"
              }
            }
          ]
        ],
        "sli_avail_state": 1,
        "infos": [
          {
            "description": "test-entityservice-to-bulk-create-info-1-description",
            "name": "test-entityservice-to-bulk-create-info-1-name",
            "value": "test-entityservice-to-bulk-create-info-1-value"
          },
          {
            "description": "test-entityservice-to-bulk-create-info-2-description",
            "name": "test-entityservice-to-bulk-create-info-2-name",
            "value": false
          },
          {
            "description": "test-entityservice-to-bulk-create-info-3-description",
            "name": "test-entityservice-to-bulk-create-info-3-name",
            "value": 1022
          },
          {
            "description": "test-entityservice-to-bulk-create-info-4-description",
            "name": "test-entityservice-to-bulk-create-info-4-name",
            "value": 10.45
          },
          {
            "description": "test-entityservice-to-bulk-create-info-5-description",
            "name": "test-entityservice-to-bulk-create-info-5-name",
            "value": null
          },
          {
            "description": "test-entityservice-to-bulk-create-info-6-description",
            "name": "test-entityservice-to-bulk-create-info-6-name",
            "value": ["test-entityservice-to-bulk-create-info-6-value", false, 1022, 10.45, null]
          }
        ]
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "id": "test-entityservice-to-bulk-create-1",
        "status": 200,
        "item": {
          "_id": "test-entityservice-to-bulk-create-1",
          "name": "test-entityservice-to-bulk-create-1-name",
          "output_template": "test-entityservice-to-bulk-create-1-output",
          "category": "test-category-to-entityservice-edit",
          "impact_level": 1,
          "enabled": true,
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-entityservice-to-bulk-create-1-pattern"
                }
              }
            ]
          ],
          "sli_avail_state": 1,
          "infos": [
            {
              "description": "test-entityservice-to-bulk-create-info-1-description",
              "name": "test-entityservice-to-bulk-create-info-1-name",
              "value": "test-entityservice-to-bulk-create-info-1-value"
            },
            {
              "description": "test-entityservice-to-bulk-create-info-2-description",
              "name": "test-entityservice-to-bulk-create-info-2-name",
              "value": false
            },
            {
              "description": "test-entityservice-to-bulk-create-info-3-description",
              "name": "test-entityservice-to-bulk-create-info-3-name",
              "value": 1022
            },
            {
              "description": "test-entityservice-to-bulk-create-info-4-description",
              "name": "test-entityservice-to-bulk-create-info-4-name",
              "value": 10.45
            },
            {
              "description": "test-entityservice-to-bulk-create-info-5-description",
              "name": "test-entityservice-to-bulk-create-info-5-name",
              "value": null
            },
            {
              "description": "test-entityservice-to-bulk-create-info-6-description",
              "name": "test-entityservice-to-bulk-create-info-6-name",
              "value": ["test-entityservice-to-bulk-create-info-6-value", false, 1022, 10.45, null]
            }
          ]
        }
      },
      {
        "status": 400,
        "item": {
          "_id": "test-entityservice-to-bulk-create-1",
          "name": "test-entityservice-to-bulk-create-1-name",
          "output_template": "test-entityservice-to-bulk-create-1-output",
          "category": "test-category-to-entityservice-edit",
          "impact_level": 1,
          "enabled": true,
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-entityservice-to-bulk-create-1-pattern"
                }
              }
            ]
          ],
          "sli_avail_state": 1,
          "infos": [
            {
              "description": "test-entityservice-to-bulk-create-info-1-description",
              "name": "test-entityservice-to-bulk-create-info-1-name",
              "value": "test-entityservice-to-bulk-create-info-1-value"
            },
            {
              "description": "test-entityservice-to-bulk-create-info-2-description",
              "name": "test-entityservice-to-bulk-create-info-2-name",
              "value": false
            },
            {
              "description": "test-entityservice-to-bulk-create-info-3-description",
              "name": "test-entityservice-to-bulk-create-info-3-name",
              "value": 1022
            },
            {
              "description": "test-entityservice-to-bulk-create-info-4-description",
              "name": "test-entityservice-to-bulk-create-info-4-name",
              "value": 10.45
            },
            {
              "description": "test-entityservice-to-bulk-create-info-5-description",
              "name": "test-entityservice-to-bulk-create-info-5-name",
              "value": null
            },
            {
              "description": "test-entityservice-to-bulk-create-info-6-description",
              "name": "test-entityservice-to-bulk-create-info-6-name",
              "value": ["test-entityservice-to-bulk-create-info-6-value", false, 1022, 10.45, null]
            }
          ]
        },
        "errors": {
          "_id": "ID already exists."
        }
      },
      {
        "status": 400,
        "item": {},
        "errors": {
          "enabled": "Enabled is missing.",
          "impact_level": "ImpactLevel is missing.",
          "name": "Name is missing.",
          "output_template": "OutputTemplate is missing.",
          "sli_avail_state": "SliAvailState is missing."
        }
      },
      {
        "status": 400,
        "item": {
          "category": "test-category-not-exist",
          "infos": [
            {}
          ],
          "sli_avail_state": 4
        },
        "errors": {
          "category": "Category doesn't exist.",
          "infos.0.name": "Name is missing.",
          "sli_avail_state": "SliAvailState should be 3 or less."
        }
      },
      {
        "status": 400,
        "item": {
          "name": "test-entityservice-to-check-unique-name-name"
        },
        "errors": {
          "name": "Name already exists."
        }
      },
      {
        "status": 400,
        "item": [],
        "error": "value doesn't contain object; it contains array"
      },
      {
        "id": "test-entityservice-to-bulk-create-2",
        "status": 200,
        "item": {
          "_id": "test-entityservice-to-bulk-create-2",
          "name": "test-entityservice-to-bulk-create-2-name",
          "output_template": "test-entityservice-to-bulk-create-2-output",
          "category": "test-category-to-entityservice-edit",
          "impact_level": 1,
          "enabled": true,
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-entityservice-to-bulk-create-2-pattern"
                }
              }
            ]
          ],
          "sli_avail_state": 1,
          "infos": [
            {
              "description": "test-entityservice-to-bulk-create-info-1-description",
              "name": "test-entityservice-to-bulk-create-info-1-name",
              "value": "test-entityservice-to-bulk-create-info-1-value"
            },
            {
              "description": "test-entityservice-to-bulk-create-info-2-description",
              "name": "test-entityservice-to-bulk-create-info-2-name",
              "value": false
            },
            {
              "description": "test-entityservice-to-bulk-create-info-3-description",
              "name": "test-entityservice-to-bulk-create-info-3-name",
              "value": 1022
            },
            {
              "description": "test-entityservice-to-bulk-create-info-4-description",
              "name": "test-entityservice-to-bulk-create-info-4-name",
              "value": 10.45
            },
            {
              "description": "test-entityservice-to-bulk-create-info-5-description",
              "name": "test-entityservice-to-bulk-create-info-5-name",
              "value": null
            },
            {
              "description": "test-entityservice-to-bulk-create-info-6-description",
              "name": "test-entityservice-to-bulk-create-info-6-name",
              "value": ["test-entityservice-to-bulk-create-info-6-value", false, 1022, 10.45, null]
            }
          ]
        }
      }
    ]
    """
    When I do GET /api/v4/entityservices/test-entityservice-to-bulk-create-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-entityservice-to-bulk-create-1",
      "category": {
        "_id": "test-category-to-entityservice-edit",
        "name": "test-category-to-entityservice-edit-name"
      },
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-entityservice-to-bulk-create-1-pattern"
            }
          }
        ]
      ],
      "impact_level": 1,
      "infos": {
        "test-entityservice-to-bulk-create-info-1-name": {
          "description": "test-entityservice-to-bulk-create-info-1-description",
          "name": "test-entityservice-to-bulk-create-info-1-name",
          "value": "test-entityservice-to-bulk-create-info-1-value"
        },
        "test-entityservice-to-bulk-create-info-2-name": {
          "description": "test-entityservice-to-bulk-create-info-2-description",
          "name": "test-entityservice-to-bulk-create-info-2-name",
          "value": false
        },
        "test-entityservice-to-bulk-create-info-3-name": {
          "description": "test-entityservice-to-bulk-create-info-3-description",
          "name": "test-entityservice-to-bulk-create-info-3-name",
          "value": 1022
        },
        "test-entityservice-to-bulk-create-info-4-name": {
          "description": "test-entityservice-to-bulk-create-info-4-description",
          "name": "test-entityservice-to-bulk-create-info-4-name",
          "value": 10.45
        },
        "test-entityservice-to-bulk-create-info-5-name": {
          "description": "test-entityservice-to-bulk-create-info-5-description",
          "name": "test-entityservice-to-bulk-create-info-5-name",
          "value": null
        },
        "test-entityservice-to-bulk-create-info-6-name": {
          "description": "test-entityservice-to-bulk-create-info-6-description",
          "name": "test-entityservice-to-bulk-create-info-6-name",
          "value": ["test-entityservice-to-bulk-create-info-6-value", false, 1022, 10.45, null]
        }
      },
      "name": "test-entityservice-to-bulk-create-1-name",
      "output_template": "test-entityservice-to-bulk-create-1-output",
      "sli_avail_state": 1,
      "type": "service"
    }
    """
    When I do GET /api/v4/entityservices/test-entityservice-to-bulk-create-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-entityservice-to-bulk-create-2",
      "category": {
        "_id": "test-category-to-entityservice-edit",
        "name": "test-category-to-entityservice-edit-name"
      },
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-entityservice-to-bulk-create-2-pattern"
            }
          }
        ]
      ],
      "impact_level": 1,
      "infos": {
        "test-entityservice-to-bulk-create-info-1-name": {
          "description": "test-entityservice-to-bulk-create-info-1-description",
          "name": "test-entityservice-to-bulk-create-info-1-name",
          "value": "test-entityservice-to-bulk-create-info-1-value"
        },
        "test-entityservice-to-bulk-create-info-2-name": {
          "description": "test-entityservice-to-bulk-create-info-2-description",
          "name": "test-entityservice-to-bulk-create-info-2-name",
          "value": false
        },
        "test-entityservice-to-bulk-create-info-3-name": {
          "description": "test-entityservice-to-bulk-create-info-3-description",
          "name": "test-entityservice-to-bulk-create-info-3-name",
          "value": 1022
        },
        "test-entityservice-to-bulk-create-info-4-name": {
          "description": "test-entityservice-to-bulk-create-info-4-description",
          "name": "test-entityservice-to-bulk-create-info-4-name",
          "value": 10.45
        },
        "test-entityservice-to-bulk-create-info-5-name": {
          "description": "test-entityservice-to-bulk-create-info-5-description",
          "name": "test-entityservice-to-bulk-create-info-5-name",
          "value": null
        },
        "test-entityservice-to-bulk-create-info-6-name": {
          "description": "test-entityservice-to-bulk-create-info-6-description",
          "name": "test-entityservice-to-bulk-create-info-6-name",
          "value": ["test-entityservice-to-bulk-create-info-6-value", false, 1022, 10.45, null]
        }
      },
      "name": "test-entityservice-to-bulk-create-2-name",
      "output_template": "test-entityservice-to-bulk-create-2-output",
      "sli_avail_state": 1,
      "type": "service"
    }
    """

Feature: Bulk update eventfilters
  I need to be able to bulk update eventfilters
  Only admin should be able to bulk update eventfilters

  Scenario: given bulk update request and no auth should not allow access
    When I do PUT /api/v4/bulk/eventfilters
    Then the response code should be 401

  Scenario: given bulk update request and auth by api key without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/bulk/eventfilters
    Then the response code should be 403

  Scenario: given bulk update request should return multi status and should be handled independently
    When I am admin
    When I do PUT /api/v4/bulk/eventfilters:
    """json
    [
      {
        "_id": "test-eventfilter-bulk-update-1",
        "description": "drop filter",
        "type": "drop",
        "event_pattern": [
          [
            {
              "field": "resource",
              "cond": {
                "type": "eq",
                "value": "test-eventfilter-to-update-pattern-updated"
              }
            }
          ]
        ],
        "priority": 1,
        "enabled": true
      },
      {
        "_id": "test-eventfilter-bulk-update-1",
        "description": "drop filter updated",
        "type": "drop",
        "event_pattern": [
          [
            {
              "field": "resource",
              "cond": {
                "type": "eq",
                "value": "test-eventfilter-to-update-pattern-updated"
              }
            }
          ]
        ],
        "priority": 1,
        "enabled": true
      },
      {
        "_id": "test-eventfilter-bulk-update-2",
        "description": "drop filter",
        "type": "drop",
        "event_pattern": [
          [
            {
              "field": "resource",
              "cond": {
                "type": "eq",
                "value": "test-eventfilter-to-update-pattern"
              }
            }
          ]
        ],
        "entity_pattern": [
          [
            {
              "field": "name",
              "cond": {
                "type": "eq",
                "value": "test-eventfilter-update-2-pattern-updated"
              }
            }
          ]
        ],
        "priority": 1,
        "enabled": true
      },
      {
        "_id": "test-eventfilter-bulk-update-3",
        "description": "drop filter",
        "type": "drop",
        "event_pattern": [
          [
            {
              "field": "resource",
              "cond": {
                "type": "eq",
                "value": "test-eventfilter-to-update-pattern-updated"
              }
            }
          ]
        ],
        "entity_pattern": [
          [
            {
              "field": "name",
              "cond": {
                "type": "eq",
                "value": "test-eventfilter-update-3-pattern-updated"
              }
            }
          ]
        ],
        "priority": 1,
        "enabled": true
      },
      {
        "_id": "test-eventfilter-bulk-update-4",
        "description": "drop filter",
        "type": "drop",
        "event_pattern": [
          [
            {
              "field": "resource",
              "cond": {
                "type": "eq",
                "value": "test-eventfilter-to-update-pattern"
              }
            }
          ]
        ],
        "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
        "priority": 1,
        "enabled": true
      },
      {
        "_id": "test-eventfilter-not-found",
        "description": "drop filter",
        "type": "drop",
        "event_pattern": [
          [
            {
              "field": "resource",
              "cond": {
                "type": "eq",
                "value": "test-eventfilter-to-update-pattern-updated"
              }
            }
          ]
        ],
        "priority": 1,
        "enabled": true
      },
      {
        "_id": "test-eventfilter-bulk-update-5",
        "description": "update change_entity",
        "type": "change_entity",
        "event_pattern": [
          [
            {
              "field": "resource",
              "cond": {
                "type": "eq",
                "value": "never be used change entity update test"
              }
            }
          ]
        ],
        "enabled": true
      },
      {
        "_id": "test-eventfilter-bulk-update-5",
        "type":"enrichment",
        "description":"Another entity copy",
        "event_pattern":[[
          {
            "field": "connector_bad",
            "cond": {
              "type": "eq",
              "value": "some"
            }
          }
        ]],
        "priority":0,
        "enabled":true,
        "config": {
          "actions":[{"value":"ExternalData.entity","name":"Entity","type":"copy"}],
          "on_success":"pass",
          "on_failure":"pass"
        },
        "external_data":{"entity":{"type":"entity"}}
      },
      {
        "_id": "test-eventfilter-bulk-update-5",
        "type":"enrichment",
        "description":"Another entity copy",
        "corporate_entity_pattern": "test-pattern-not-exist",
        "priority":0,
        "enabled":true,
        "config": {
          "actions":[{"value":"ExternalData.entity","name":"Entity","type":"copy"}],
          "on_success":"pass",
          "on_failure":"pass"
        },
        "external_data":{"entity":{"type":"entity"}}
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "item": {
          "_id": "test-eventfilter-bulk-update-1",
          "description": "drop filter",
          "type": "drop",
          "event_pattern": [
            [
              {
                "field": "resource",
                "cond": {
                  "type": "eq",
                  "value": "test-eventfilter-to-update-pattern-updated"
                }
              }
            ]
          ],
          "priority": 1,
          "enabled": true
        }
      },
      {
        "status": 200,
        "item": {
          "_id": "test-eventfilter-bulk-update-1",
          "description": "drop filter updated",
          "type": "drop",
          "event_pattern": [
            [
              {
                "field": "resource",
                "cond": {
                  "type": "eq",
                  "value": "test-eventfilter-to-update-pattern-updated"
                }
              }
            ]
          ],
          "priority": 1,
          "enabled": true
        }
      },
      {
        "status": 200,
        "item": {
          "_id": "test-eventfilter-bulk-update-2",
          "description": "drop filter",
          "type": "drop",
          "event_pattern": [
            [
              {
                "field": "resource",
                "cond": {
                  "type": "eq",
                  "value": "test-eventfilter-to-update-pattern"
                }
              }
            ]
          ],
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-eventfilter-update-2-pattern-updated"
                }
              }
            ]
          ],
          "priority": 1,
          "enabled": true
        }
      },
      {
        "status": 200,
        "item": {
          "_id": "test-eventfilter-bulk-update-3",
          "description": "drop filter",
          "type": "drop",
          "event_pattern": [
            [
              {
                "field": "resource",
                "cond": {
                  "type": "eq",
                  "value": "test-eventfilter-to-update-pattern-updated"
                }
              }
            ]
          ],
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-eventfilter-update-3-pattern-updated"
                }
              }
            ]
          ],
          "priority": 1,
          "enabled": true
        }
      },
      {
        "status": 200,
        "item": {
          "_id": "test-eventfilter-bulk-update-4",
          "description": "drop filter",
          "type": "drop",
          "event_pattern": [
            [
              {
                "field": "resource",
                "cond": {
                  "type": "eq",
                  "value": "test-eventfilter-to-update-pattern"
                }
              }
            ]
          ],
          "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
          "priority": 1,
          "enabled": true
        }
      },
      {
        "status": 404,
        "item": {
          "_id": "test-eventfilter-not-found",
          "description": "drop filter",
          "type": "drop",
          "event_pattern": [
            [
              {
                "field": "resource",
                "cond": {
                  "type": "eq",
                  "value": "test-eventfilter-to-update-pattern-updated"
                }
              }
            ]
          ],
          "priority": 1,
          "enabled": true
        },
        "error": "Not found"
      },
      {
        "status": 400,
        "item": {
          "_id": "test-eventfilter-bulk-update-5",
          "description": "update change_entity",
          "type": "change_entity",
          "event_pattern": [
            [
              {
                "field": "resource",
                "cond": {
                  "type": "eq",
                  "value": "never be used change entity update test"
                }
              }
            ]
          ],
          "enabled": true
        },
        "errors": {
          "config": "Config is missing."
        }
      },
      {
        "status": 400,
        "item": {
          "_id": "test-eventfilter-bulk-update-5",
          "type":"enrichment",
          "description":"Another entity copy",
          "event_pattern":[[
            {
              "field": "connector_bad",
              "cond": {
                "type": "eq",
                "value": "some"
              }
            }
          ]],
          "priority":0,
          "enabled":true,
          "config": {
            "actions":[{"value":"ExternalData.entity","name":"Entity","type":"copy"}],
            "on_success":"pass",
            "on_failure":"pass"
          },
          "external_data":{"entity":{"type":"entity"}}
        },
        "errors": {
          "event_pattern": "EventPattern is invalid event pattern."
        }
      },
      {
        "status": 400,
        "item": {
          "_id": "test-eventfilter-bulk-update-5",
          "type":"enrichment",
          "description":"Another entity copy",
          "corporate_entity_pattern": "test-pattern-not-exist",
          "priority":0,
          "enabled":true,
          "config": {
            "actions":[{"value":"ExternalData.entity","name":"Entity","type":"copy"}],
            "on_success":"pass",
            "on_failure":"pass"
          },
          "external_data":{"entity":{"type":"entity"}}
        },
        "errors": {
          "corporate_entity_pattern": "CorporateEntityPattern doesn't exist."
        }
      }
    ]
    """
    When I do GET /api/v4/eventfilter/rules?search=test-eventfilter-bulk-update
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-eventfilter-bulk-update-1",
          "author": "root",
          "description": "drop filter updated",
          "type": "drop",
          "event_pattern": [
            [
              {
                "field": "resource",
                "cond": {
                  "type": "eq",
                  "value": "test-eventfilter-to-update-pattern-updated"
                }
              }
            ]
          ],
          "priority": 1,
          "enabled": true,
          "created": 1608635535
        },
        {
          "_id": "test-eventfilter-bulk-update-2",
          "author": "root",
          "description": "drop filter",
          "type": "drop",
          "event_pattern": [
            [
              {
                "field": "resource",
                "cond": {
                  "type": "eq",
                  "value": "test-eventfilter-to-update-pattern"
                }
              }
            ]
          ],
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-eventfilter-update-2-pattern-updated"
                }
              }
            ]
          ],
          "priority": 1,
          "enabled": true,
          "created": 1608635535
        },
        {
          "_id": "test-eventfilter-bulk-update-3",
          "author": "root",
          "description": "drop filter",
          "type": "drop",
          "event_pattern": [
            [
              {
                "field": "resource",
                "cond": {
                  "type": "eq",
                  "value": "test-eventfilter-to-update-pattern-updated"
                }
              }
            ]
          ],
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-eventfilter-update-3-pattern-updated"
                }
              }
            ]
          ],
          "priority": 1,
          "enabled": true,
          "created": 1608635535
        },
        {
          "_id": "test-eventfilter-bulk-update-4",
          "author": "root",
          "description": "drop filter",
          "type": "drop",
          "event_pattern": [
            [
              {
                "field": "resource",
                "cond": {
                  "type": "eq",
                  "value": "test-eventfilter-to-update-pattern"
                }
              }
            ]
          ],
          "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
          "corporate_entity_pattern_title": "test-pattern-to-rule-edit-2-title",
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-pattern-to-rule-edit-2-pattern"
                }
              }
            ]
          ],
          "priority": 1,
          "enabled": true,
          "created": 1608635535
        },
        {
          "_id": "test-eventfilter-bulk-update-5",
          "author": "root",
          "description": "break filter",
          "enabled": true,
          "event_pattern": [
            [
              {
                "cond": {
                  "type": "eq",
                  "value": "test-eventfilter-bulk-update-5-pattern"
                },
                "field": "resource"
              }
            ]
          ],
          "priority": 3,
          "type": "break"
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

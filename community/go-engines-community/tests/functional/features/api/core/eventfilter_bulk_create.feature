Feature: Bulk create eventfilters
  I need to be able to bulk create eventfilters
  Only admin should be able to bulk create eventfilters

  Scenario: given bulk create request and no auth should not allow access
    When I do POST /api/v4/bulk/eventfilters
    Then the response code should be 401

  Scenario: given bulk create request and auth by api key without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/bulk/eventfilters
    Then the response code should be 403

  Scenario: given bulk create request should return multi status and should be handled independently
    When I am admin
    When I do POST /api/v4/bulk/eventfilters:
    """json
    [
      {
        "_id": "test-eventfilter-bulk-create-1",
        "description": "test create 1",
        "type": "enrichment",
        "event_pattern": [
          [
            {
              "field": "connector",
              "cond": {
                "type": "eq",
                "value": "test-eventfilter-create-1-pattern"
              }
            }
          ]
        ],
        "priority": 0,
        "enabled": true,
        "config": {
          "actions": [
            {
              "type": "set_field",
              "name": "connector",
              "value": "kafka_connector"
            }
          ],
          "on_success": "pass",
          "on_failure": "pass"
        },
        "external_data": {
          "test": {
            "type": "mongo"
          }
        }
      },
      {
        "_id": "test-eventfilter-bulk-create-1",
        "description": "test create 1",
        "type": "enrichment",
        "event_pattern": [
          [
            {
              "field": "connector",
              "cond": {
                "type": "eq",
                "value": "test-eventfilter-create-1-pattern"
              }
            }
          ]
        ],
        "priority": 0,
        "enabled": true,
        "config": {
          "actions": [
            {
              "type": "set_field",
              "name": "connector",
              "value": "kafka_connector"
            }
          ],
          "on_success": "pass",
          "on_failure": "pass"
        },
        "external_data": {
          "test": {
            "type": "mongo"
          }
        }
      },
      {
        "_id": "test-eventfilter-bulk-create-2",
        "description": "test create 2",
        "type": "enrichment",
        "entity_pattern": [
          [
            {
              "field": "name",
              "cond": {
                "type": "eq",
                "value": "test-eventfilter-create-2-pattern"
              }
            }
          ]
        ],
        "priority": 0,
        "enabled": true,
        "config": {
          "actions": [
            {
              "type": "set_field",
              "name": "connector",
              "value": "kafka_connector"
            }
          ],
          "on_success": "pass",
          "on_failure": "pass"
        },
        "external_data": {
          "test": {
            "type": "mongo"
          }
        }
      },
      {
        "_id": "test-eventfilter-bulk-create-3",
        "description": "test create 3",
        "type": "enrichment",
        "event_pattern": [
          [
            {
              "field": "connector",
              "cond": {
                "type": "eq",
                "value": "test-eventfilter-create-3-pattern"
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
                "value": "test-eventfilter-create-4-pattern"
              }
            }
          ]
        ],
        "priority": 0,
        "enabled": true,
        "config": {
          "actions": [
            {
              "type": "set_field",
              "name": "connector",
              "value": "kafka_connector"
            }
          ],
          "on_success": "pass",
          "on_failure": "pass"
        },
        "external_data": {
          "test": {
            "type": "mongo"
          }
        }
      },
      {
        "_id": "test-eventfilter-bulk-create-4",
        "description": "test create 4",
        "type": "enrichment",
        "event_pattern": [
          [
            {
              "field": "connector",
              "cond": {
                "type": "eq",
                "value": "test-eventfilter-create-4-pattern"
              }
            }
          ]
        ],
        "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
        "priority": 0,
        "enabled": true,
        "config": {
          "actions": [
            {
              "type": "set_field",
              "name": "connector",
              "value": "kafka_connector"
            }
          ],
          "on_success": "pass",
          "on_failure": "pass"
        },
        "external_data": {
          "test": {
            "type": "mongo"
          }
        }
      },
      {
        "_id": "test-eventfilter-bulk-create-5",
        "description": "test create 4",
        "type": "enrichment",
        "event_pattern": [
          [
            {
              "field": "connector",
              "cond": {
                "type": "eq",
                "value": "test-eventfilter-create-4-pattern"
              }
            }
          ]
        ],
        "corporate_entity_pattern": "test-pattern-not-exist",
        "priority": 0,
        "enabled": true,
        "config": {
          "actions": [
            {
              "type": "set_field",
              "name": "connector",
              "value": "kafka_connector"
            }
          ],
          "on_success": "pass",
          "on_failure": "pass"
        },
        "external_data": {
          "test": {
            "type": "mongo"
          }
        }
      },
      {
        "type": "unspecified"
      },
      {
        "type": "enrichment",
        "description": "some",
        "config": {
          "actions": []
        }
      },
      {
        "type": "enrichment",
        "description": "some",
        "config": {
          "actions": [
            {
              "type":"set_entity_info_from_template",
              "name":"test",
              "value":"{{ `{{.ExternalData.test}}` }}",
              "description":"test"
            }
          ],
          "on_failure": "continue",
          "on_success": "continue"
        }
      },
      {
        "type":"enrichment",
        "description":"Another entity copy",
        "event_pattern":[[]],
        "entity_pattern": [[]],
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
        "description": "test",
        "type": "change_entity",
        "event_pattern":[[
          {
            "field": "connector",
            "cond": {
              "type": "eq",
              "value": "test_connector"
            }
          }
        ]],
        "enabled": true
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "id": "test-eventfilter-bulk-create-1",
        "status": 200,
        "item": {
          "_id": "test-eventfilter-bulk-create-1",
          "description": "test create 1",
          "type": "enrichment",
          "event_pattern": [
            [
              {
                "field": "connector",
                "cond": {
                  "type": "eq",
                  "value": "test-eventfilter-create-1-pattern"
                }
              }
            ]
          ],
          "priority": 0,
          "enabled": true,
          "config": {
            "actions": [
              {
                "type": "set_field",
                "name": "connector",
                "value": "kafka_connector"
              }
            ],
            "on_success": "pass",
            "on_failure": "pass"
          },
          "external_data": {
            "test": {
              "type": "mongo"
            }
          }
        }
      },
      {
        "status": 400,
        "item": {
          "_id": "test-eventfilter-bulk-create-1",
          "description": "test create 1",
          "type": "enrichment",
          "event_pattern": [
            [
              {
                "field": "connector",
                "cond": {
                  "type": "eq",
                  "value": "test-eventfilter-create-1-pattern"
                }
              }
            ]
          ],
          "priority": 0,
          "enabled": true,
          "config": {
            "actions": [
              {
                "type": "set_field",
                "name": "connector",
                "value": "kafka_connector"
              }
            ],
            "on_success": "pass",
            "on_failure": "pass"
          },
          "external_data": {
            "test": {
              "type": "mongo"
            }
          }
        },
        "errors": {
          "_id": "ID already exists."
        }
      },
      {
        "status": 200,
        "item": {
          "_id": "test-eventfilter-bulk-create-2",
          "description": "test create 2",
          "type": "enrichment",
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-eventfilter-create-2-pattern"
                }
              }
            ]
          ],
          "priority": 0,
          "enabled": true,
          "config": {
            "actions": [
              {
                "type": "set_field",
                "name": "connector",
                "value": "kafka_connector"
              }
            ],
            "on_success": "pass",
            "on_failure": "pass"
          },
          "external_data": {
            "test": {
              "type": "mongo"
            }
          }
        }
      },
      {
        "status": 200,
        "item": {
          "_id": "test-eventfilter-bulk-create-3",
          "description": "test create 3",
          "type": "enrichment",
          "event_pattern": [
            [
              {
                "field": "connector",
                "cond": {
                  "type": "eq",
                  "value": "test-eventfilter-create-3-pattern"
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
                  "value": "test-eventfilter-create-4-pattern"
                }
              }
            ]
          ],
          "priority": 0,
          "enabled": true,
          "config": {
            "actions": [
              {
                "type": "set_field",
                "name": "connector",
                "value": "kafka_connector"
              }
            ],
            "on_success": "pass",
            "on_failure": "pass"
          },
          "external_data": {
            "test": {
              "type": "mongo"
            }
          }
        }
      },
      {
        "status": 200,
        "item": {
          "_id": "test-eventfilter-bulk-create-4",
          "description": "test create 4",
          "type": "enrichment",
          "event_pattern": [
            [
              {
                "field": "connector",
                "cond": {
                  "type": "eq",
                  "value": "test-eventfilter-create-4-pattern"
                }
              }
            ]
          ],
          "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
          "priority": 0,
          "enabled": true,
          "config": {
            "actions": [
              {
                "type": "set_field",
                "name": "connector",
                "value": "kafka_connector"
              }
            ],
            "on_success": "pass",
            "on_failure": "pass"
          },
          "external_data": {
            "test": {
              "type": "mongo"
            }
          }
        }
      },
      {
        "status": 400,
        "item": {
          "_id": "test-eventfilter-bulk-create-5",
          "description": "test create 4",
          "type": "enrichment",
          "event_pattern": [
            [
              {
                "field": "connector",
                "cond": {
                  "type": "eq",
                  "value": "test-eventfilter-create-4-pattern"
                }
              }
            ]
          ],
          "corporate_entity_pattern": "test-pattern-not-exist",
          "priority": 0,
          "enabled": true,
          "config": {
            "actions": [
              {
                "type": "set_field",
                "name": "connector",
                "value": "kafka_connector"
              }
            ],
            "on_success": "pass",
            "on_failure": "pass"
          },
          "external_data": {
            "test": {
              "type": "mongo"
            }
          }
        },
        "errors": {
          "corporate_entity_pattern": "CorporateEntityPattern doesn't exist."
        }
      },
      {
        "status": 400,
        "item": {
          "type": "unspecified"
        },
        "errors": {
          "type": "Type must be one of [break drop enrichment change_entity]."
        }
      },
      {
        "status": 400,
        "item": {
          "type": "enrichment",
          "description": "some",
          "config": {
            "actions": []
          }
        },
        "errors": {
          "actions": "Actions is missing.",
          "on_failure": "OnFailure is required when Type enrichment is defined.",
          "on_success": "OnSuccess is required when Type enrichment is defined."
        }
      },
      {
        "status": 400,
        "item": {
          "type": "enrichment",
          "description": "some",
          "config": {
            "actions": [
              {
                "type":"set_entity_info_from_template",
                "name":"test",
                "value":"{{ `{{.ExternalData.test}}` }}",
                "description":"test"
              }
            ],
            "on_failure": "continue",
            "on_success": "continue"
          }
        },
        "errors": {
          "on_failure": "OnFailure must be one of [pass drop break].",
          "on_success": "OnSuccess must be one of [pass drop break]."
        }
      },
      {
        "status": 400,
        "item": {
          "type":"enrichment",
          "description":"Another entity copy",
          "event_pattern":[[]],
          "entity_pattern": [[]],
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
          "entity_pattern": "EntityPattern is invalid entity pattern.",
          "event_pattern": "EventPattern is invalid event pattern."
        }
      },
      {
        "status": 400,
        "item": {
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
          "description": "test",
          "type": "change_entity",
          "event_pattern":[[
            {
              "field": "connector",
              "cond": {
                "type": "eq",
                "value": "test_connector"
              }
            }
          ]],
          "enabled": true
        },
        "errors": {
          "config": "Config is missing."
        }
      }
    ]
    """
    When I do GET /api/v4/eventfilter/rules?search=test-eventfilter-bulk-create
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-eventfilter-bulk-create-1",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "description": "test create 1",
          "type": "enrichment",
          "event_pattern": [
            [
              {
                "field": "connector",
                "cond": {
                  "type": "eq",
                  "value": "test-eventfilter-create-1-pattern"
                }
              }
            ]
          ],
          "priority": 0,
          "enabled": true,
          "config": {
            "actions": [
              {
                "type": "set_field",
                "name": "connector",
                "value": "kafka_connector"
              }
            ],
            "on_success": "pass",
            "on_failure": "pass"
          },
          "external_data": {
            "test": {
              "type": "mongo"
            }
          }
        },
        {
          "_id": "test-eventfilter-bulk-create-2",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "description": "test create 2",
          "type": "enrichment",
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-eventfilter-create-2-pattern"
                }
              }
            ]
          ],
          "priority": 0,
          "enabled": true,
          "config": {
            "actions": [
              {
                "type": "set_field",
                "name": "connector",
                "value": "kafka_connector"
              }
            ],
            "on_success": "pass",
            "on_failure": "pass"
          },
          "external_data": {
            "test": {
              "type": "mongo"
            }
          }
        },
        {
          "_id": "test-eventfilter-bulk-create-3",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "description": "test create 3",
          "type": "enrichment",
          "event_pattern": [
            [
              {
                "field": "connector",
                "cond": {
                  "type": "eq",
                  "value": "test-eventfilter-create-3-pattern"
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
                  "value": "test-eventfilter-create-4-pattern"
                }
              }
            ]
          ],
          "priority": 0,
          "enabled": true,
          "config": {
            "actions": [
              {
                "type": "set_field",
                "name": "connector",
                "value": "kafka_connector"
              }
            ],
            "on_success": "pass",
            "on_failure": "pass"
          },
          "external_data": {
            "test": {
              "type": "mongo"
            }
          }
        },
        {
          "_id": "test-eventfilter-bulk-create-4",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "description": "test create 4",
          "type": "enrichment",
          "event_pattern": [
            [
              {
                "field": "connector",
                "cond": {
                  "type": "eq",
                  "value": "test-eventfilter-create-4-pattern"
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
          "priority": 0,
          "enabled": true,
          "config": {
            "actions": [
              {
                "type": "set_field",
                "name": "connector",
                "value": "kafka_connector"
              }
            ],
            "on_success": "pass",
            "on_failure": "pass"
          },
          "external_data": {
            "test": {
              "type": "mongo"
            }
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 4
      }
    }
    """

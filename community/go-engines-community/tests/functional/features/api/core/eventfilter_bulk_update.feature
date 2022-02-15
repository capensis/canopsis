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
        "description": "test update 1",
        "type": "enrichment",
        "patterns": [
          {
            "connector": "test-eventfilter-update-1-pattern"
          }
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
        "_id": "test-eventfilter-bulk-update-1",
        "description": "test update 1111",
        "type": "enrichment",
        "patterns": [
          {
            "connector": "test-eventfilter-update-1-pattern"
          }
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
        "_id": "test-eventfilter-bulk-update-2",
        "type": "unspecified"
      },
      {
        "_id": "test-eventfilter-bulk-update-2",
        "type": "enrichment",
        "description": "some",
        "config": {
          "actions": []
        }
      },
      {
        "_id": "test-eventfilter-bulk-update-2",
        "type": "enrichment",
        "description": "some",
        "config": {
          "actions": [
            {
              "type":"set_entity_info_from_template",
              "name":"test",
              "value":"{{.ExternalData.test}}",
              "description":"test"
            }
          ],
          "on_failure": "continue",
          "on_success": "continue"
        }
      },
      {
        "_id": "test-eventfilter-bulk-update-2",
        "description": "test update 2",
        "type": "enrichment",
        "patterns": [
          {
            "connector": "test-eventfilter-update-2-pattern"
          }
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
        "_id": "test-eventfilter-bulk-update-3",
        "description": "test update 3",
        "type": "enrichment",
        "patterns": [
          {
            "connector": "test-eventfilter-update-3-pattern"
          }
        ],
        "priority": 0,
        "enabled": false,
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
        "_id": "test-eventfilter-bulk-update-4",
        "type":"enrichment",
        "description":"some description",
        "patterns":[{}],
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
        "description":"some description",
        "patterns":null,
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
        "description":"some description",
        "patterns":[4],
        "priority":0,
        "enabled":true,
        "config": {
          "actions":[{"value":"ExternalData.entity","name":"Entity","type":"copy"}],
          "on_success":"pass",
          "on_failure":"pass"
        },
        "external_data":{"entity":{"type":"entity"}},
        "author": "root"
      },
      {
        "_id": "test-eventfilter-bulk-update-5",
        "type":"enrichment",
        "description":"Invalid pattern with empty document",
        "patterns":[{},{"connector": "test-eventfilter-update-1-pattern"}],
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
        "description": "test",
        "type": "change_entity",
        "patterns": [
          {
            "connector": "test_connector",
            "customer_tags": {
              "regex_match": "CMDB:(?P<SI_CMDB>.*?)($|,)"
            }
          }
        ],
        "enabled": true
      },
      {
        "_id": "test-eventfilter-bulk-update-5",
        "description": "test",
        "type": "change_entity",
        "patterns": [
          {
            "connector": "test_connector",
            "customer_tags": {
              "regex_match": "CMDB:(?P<SI_CMDB>.*?)($|,)"
            }
          }
        ],
        "config": {
          "component": "",
          "connector": "",
          "resource": "",
          "connector_name": ""
        },
        "enabled": true
      },
      {
        "_id": "test-eventfilter-bulk-update-5",
        "description": "test",
        "type": "change_entity",
        "patterns": [
          {
            "connector": "test_connector",
            "customer_tags": {
              "regex_match": "CMDB:(?P<SI_CMDB>.*?)($|,)"
            }
          }
        ],
        "config": {},
        "enabled": true
      },
      {
        "_id": "test-eventfilter-bulk-update-6",
        "type":"enrichment",
        "description":"some description",
        "patterns":[{}],
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
        "id": "test-eventfilter-bulk-update-1",
        "status": 200,
        "item": {
          "_id": "test-eventfilter-bulk-update-1",
          "description": "test update 1",
          "type": "enrichment",
          "patterns": [
            {
              "connector": "test-eventfilter-update-1-pattern"
            }
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
          "_id": "test-eventfilter-bulk-update-1",
          "description": "test update 1111",
          "type": "enrichment",
          "patterns": [
            {
              "connector": "test-eventfilter-update-1-pattern"
            }
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
                "value":"{{.ExternalData.test}}",
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
        "id": "test-eventfilter-bulk-update-2",
        "status": 200,
        "item": {
          "_id": "test-eventfilter-bulk-update-2",
          "description": "test update 2",
          "type": "enrichment",
          "patterns": [
            {
              "connector": "test-eventfilter-update-2-pattern"
            }
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
        "id": "test-eventfilter-bulk-update-3",
        "status": 200,
        "item": {
          "_id": "test-eventfilter-bulk-update-3",
          "description": "test update 3",
          "type": "enrichment",
          "patterns": [
            {
              "connector": "test-eventfilter-update-3-pattern"
            }
          ],
          "priority": 0,
          "enabled": false,
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
          "_id": "test-eventfilter-bulk-update-4",
          "type":"enrichment",
          "description":"some description",
          "patterns":[{}],
          "priority":0,
          "enabled":true,
          "config": {
            "actions":[{"value":"ExternalData.entity","name":"Entity","type":"copy"}],
            "on_success":"pass",
            "on_failure":"pass"
          },
          "external_data":{"entity":{"type":"entity"}}
        }
      },
      {
        "status": 200,
        "item": {
          "_id": "test-eventfilter-bulk-update-5",
          "type":"enrichment",
          "description":"some description",
          "patterns":null,
          "priority":0,
          "enabled":true,
          "config": {
            "actions":[{"value":"ExternalData.entity","name":"Entity","type":"copy"}],
            "on_success":"pass",
            "on_failure":"pass"
          },
          "external_data":{"entity":{"type":"entity"}}
        }
      },
      {
        "status": 400,
        "item": {
          "_id": "test-eventfilter-bulk-update-5",
          "type":"enrichment",
          "description":"some description",
          "patterns":[4],
          "priority":0,
          "enabled":true,
          "config": {
            "actions":[{"value":"ExternalData.entity","name":"Entity","type":"copy"}],
            "on_success":"pass",
            "on_failure":"pass"
          },
          "external_data":{"entity":{"type":"entity"}},
          "author": "root"
        },
        "error": "error decoding key list: unable to parse event pattern list element"
      },
      {
        "status": 400,
        "item": {
          "_id": "test-eventfilter-bulk-update-5",
          "type":"enrichment",
          "description":"Invalid pattern with empty document",
          "patterns":[{},{"connector": "test-eventfilter-update-1-pattern"}],
          "priority":0,
          "enabled":true,
          "config": {
            "actions":[{"value":"ExternalData.entity","name":"Entity","type":"copy"}],
            "on_success":"pass",
            "on_failure":"pass"
          },
          "external_data":{"entity":{"type":"entity"}}
        },
        "error": "error decoding key list: unable to parse event pattern list element"
      },
      {
        "status": 400,
        "item": {
          "_id": "test-eventfilter-bulk-update-5",
          "description": "test",
          "type": "change_entity",
          "patterns": [
            {
              "connector": "test_connector",
              "customer_tags": {
                "regex_match": "CMDB:(?P<SI_CMDB>.*?)($|,)"
              }
            }
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
          "description": "test",
          "type": "change_entity",
          "patterns": [
            {
              "connector": "test_connector",
              "customer_tags": {
                "regex_match": "CMDB:(?P<SI_CMDB>.*?)($|,)"
              }
            }
          ],
          "config": {
            "component": "",
            "connector": "",
            "resource": "",
            "connector_name": ""
          },
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
          "description": "test",
          "type": "change_entity",
          "patterns": [
            {
              "connector": "test_connector",
              "customer_tags": {
                "regex_match": "CMDB:(?P<SI_CMDB>.*?)($|,)"
              }
            }
          ],
          "config": {},
          "enabled": true
        },
        "errors": {
          "config": "Config is missing."
        }
      },
      {
        "status": 404,
        "item": {
          "_id": "test-eventfilter-bulk-update-6",
          "type":"enrichment",
          "description":"some description",
          "patterns":[{}],
          "priority":0,
          "enabled":true,
          "config": {
            "actions":[{"value":"ExternalData.entity","name":"Entity","type":"copy"}],
            "on_success":"pass",
            "on_failure":"pass"
          },
          "external_data":{"entity":{"type":"entity"}}
        },
        "error": "Not found"
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
          "description": "test update 1111",
          "type": "enrichment",
          "patterns": [
            {
              "connector": "test-eventfilter-update-1-pattern"
            }
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
          "_id": "test-eventfilter-bulk-update-2",
          "author": "root",
          "description": "test update 2",
          "type": "enrichment",
          "patterns": [
            {
              "connector": "test-eventfilter-update-2-pattern"
            }
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
          "_id": "test-eventfilter-bulk-update-3",
          "author": "root",
          "description": "test update 3",
          "type": "enrichment",
          "patterns": [
            {
              "connector": "test-eventfilter-update-3-pattern"
            }
          ],
          "priority": 0,
          "enabled": false,
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
          "_id": "test-eventfilter-bulk-update-4",
          "type":"enrichment",
          "description":"some description",
          "patterns":[{}],
          "priority":0,
          "enabled":true,
          "config": {
            "actions":[{"value":"ExternalData.entity","name":"Entity","type":"copy"}],
            "on_success":"pass",
            "on_failure":"pass"
          },
          "external_data":{"entity":{"type":"entity"}},
          "author": "root"
        },
        {
          "_id": "test-eventfilter-bulk-update-5",
          "type":"enrichment",
          "description":"some description",
          "patterns":null,
          "priority":0,
          "enabled":true,
          "config": {
            "actions":[{"value":"ExternalData.entity","name":"Entity","type":"copy"}],
            "on_success":"pass",
            "on_failure":"pass"
          },
          "external_data":{"entity":{"type":"entity"}},
          "author": "root"
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

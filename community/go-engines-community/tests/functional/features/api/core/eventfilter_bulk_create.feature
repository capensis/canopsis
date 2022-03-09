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
        "description": "test bulk create 1",
        "type": "enrichment",
        "patterns": [
          {
            "connector": "test-eventfilter-bulk-create-1-pattern"
          }
        ],
        "priority": 0,
        "enabled": true,
        "actions": [
          {
            "type": "set_field",
            "name": "connector",
            "value": "kafka_connector"
          }
        ],
        "external_data": {
          "clear": "sky",
          "type": "no",
          "arr": [
            1,
            2,
            3
          ]
        },
        "on_success": "pass",
        "on_failure": "pass"
      },
      {
        "_id": "test-eventfilter-bulk-create-1",
        "description": "test bulk create 1",
        "type": "enrichment",
        "patterns": [
          {
            "connector": "test-eventfilter-bulk-create-1-pattern"
          }
        ],
        "priority": 0,
        "enabled": true,
        "actions": [
          {
            "type": "set_field",
            "name": "connector",
            "value": "kafka_connector"
          }
        ],
        "external_data": {
          "clear": "sky",
          "type": "no",
          "arr": [
            1,
            2,
            3
          ]
        },
        "on_success": "pass",
        "on_failure": "pass"
      },
      {
        "type": "unspecified"
      },
      {
        "type": "enrichment",
        "actions": []
      },
      {
        "type": "enrichment",
        "on_failure": "continue",
        "on_success": "continue"
      },
      {
        "_id": "test-eventfilter-check-id"
      },
      {
        "type": "enrichment",
        "description": "More entity copy",
        "patterns": [
          4
        ],
        "priority": 0,
        "enabled": true,
        "actions": [
          {
            "from": "ExternalData.entity",
            "to": "Entity",
            "type": "copy"
          }
        ],
        "external_data": {
          "entity": {
            "type": "entity"
          }
        },
        "on_success": "pass",
        "on_failure": "pass",
        "author": "root"
      },
      {
        "type": "enrichment",
        "description": "Invalid pattern with empty document",
        "patterns": [
          {},
          {
            "connector": "test-eventfilter-bulk-create-1-pattern"
          }
        ],
        "priority": 0,
        "enabled": true,
        "actions": [
          {
            "from": "ExternalData.entity",
            "to": "Entity",
            "type": "copy"
          }
        ],
        "external_data": {
          "entity": {
            "type": "entity"
          }
        },
        "on_success": "pass",
        "on_failure": "pass",
        "author": "root"
      },
      {
        "_id": "test-eventfilter-bulk-create-2",
        "description": "test bulk create 2",
        "type": "enrichment",
        "patterns": [
          {
            "connector": "test-eventfilter-bulk-create-2-pattern"
          }
        ],
        "priority": 0,
        "enabled": false,
        "actions": [
          {
            "type": "set_field",
            "name": "connector",
            "value": "kafka_connector"
          }
        ],
        "external_data": {
          "clear": "sky",
          "type": "no",
          "arr": [
            1,
            2,
            3
          ]
        },
        "on_success": "pass",
        "on_failure": "pass"
      },
      {
        "_id": "test-eventfilter-bulk-create-3",
        "type": "enrichment",
        "description": "test bulk create 3",
        "patterns": [
          null
        ],
        "priority": 0,
        "enabled": true,
        "actions": [
          {
            "from": "ExternalData.entity",
            "to": "Entity",
            "type": "copy"
          }
        ],
        "external_data": {
          "entity": {
            "type": "entity"
          }
        },
        "on_success": "pass",
        "on_failure": "pass",
        "author": "root"
      },
      {
        "_id": "test-eventfilter-bulk-create-4",
        "type": "enrichment",
        "description": "test bulk create 4",
        "patterns": [
          {}
        ],
        "priority": 0,
        "enabled": true,
        "actions": [
          {
            "from": "ExternalData.entity",
            "to": "Entity",
            "type": "copy"
          }
        ],
        "external_data": {
          "entity": {
            "type": "entity"
          }
        },
        "on_success": "pass",
        "on_failure": "pass",
        "author": "root"
      },
      {
        "_id": "test-eventfilter-bulk-create-5",
        "type": "enrichment",
        "description": "test bulk create 5",
        "patterns": null,
        "priority": 0,
        "enabled": true,
        "actions": [
          {
            "from": "ExternalData.entity",
            "to": "Entity",
            "type": "copy"
          }
        ],
        "external_data": {
          "entity": {
            "type": "entity"
          }
        },
        "on_success": "pass",
        "on_failure": "pass"
      },
      []
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
          "description": "test bulk create 1",
          "type": "enrichment",
          "patterns": [
            {
              "connector": "test-eventfilter-bulk-create-1-pattern"
            }
          ],
          "priority": 0,
          "enabled": true,
          "actions": [
            {
              "type": "set_field",
              "name": "connector",
              "value": "kafka_connector"
            }
          ],
          "external_data": {
            "clear": "sky",
            "type": "no",
            "arr": [
              1,
              2,
              3
            ]
          },
          "on_success": "pass",
          "on_failure": "pass"
        }
      },
      {
        "status": 400,
        "item": {
          "_id": "test-eventfilter-bulk-create-1",
          "description": "test bulk create 1",
          "type": "enrichment",
          "patterns": [
            {
              "connector": "test-eventfilter-bulk-create-1-pattern"
            }
          ],
          "priority": 0,
          "enabled": true,
          "actions": [
            {
              "type": "set_field",
              "name": "connector",
              "value": "kafka_connector"
            }
          ],
          "external_data": {
            "clear": "sky",
            "type": "no",
            "arr": [
              1,
              2,
              3
            ]
          },
          "on_success": "pass",
          "on_failure": "pass"
        },
        "errors": {
          "_id": "ID already exists."
        }
      },
      {
        "status": 400,
        "item": {
          "type": "unspecified"
        },
        "errors": {
          "type": "Type must be one of [break drop enrichment]."
        }
      },
      {
        "status": 400,
        "item": {
          "type": "enrichment",
          "actions": []
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
          "on_failure": "continue",
          "on_success": "continue"
        },
        "errors": {
          "on_failure": "OnFailure must be one of [pass drop break].",
          "on_success": "OnSuccess must be one of [pass drop break]."
        }
      },
      {
        "status": 400,
        "item": {
          "_id": "test-eventfilter-check-id"
        },
        "errors": {
          "_id": "ID already exists."
        }
      },
      {
        "status": 400,
        "item": {
          "type": "enrichment",
          "description": "More entity copy",
          "patterns": [
            4
          ],
          "priority": 0,
          "enabled": true,
          "actions": [
            {
              "from": "ExternalData.entity",
              "to": "Entity",
              "type": "copy"
            }
          ],
          "external_data": {
            "entity": {
              "type": "entity"
            }
          },
          "on_success": "pass",
          "on_failure": "pass",
          "author": "root"
        },
        "error": "error decoding key list: unable to parse event pattern list element"
      },
      {
        "status": 400,
        "item": {
          "type": "enrichment",
          "description": "Invalid pattern with empty document",
          "patterns": [
            {},
            {
              "connector": "test-eventfilter-bulk-create-1-pattern"
            }
          ],
          "priority": 0,
          "enabled": true,
          "actions": [
            {
              "from": "ExternalData.entity",
              "to": "Entity",
              "type": "copy"
            }
          ],
          "external_data": {
            "entity": {
              "type": "entity"
            }
          },
          "on_success": "pass",
          "on_failure": "pass",
          "author": "root"
        },
        "error": "error decoding key list: unable to parse event pattern list element"
      },
      {
        "id": "test-eventfilter-bulk-create-2",
        "status": 200,
        "item": {
          "_id": "test-eventfilter-bulk-create-2",
          "description": "test bulk create 2",
          "type": "enrichment",
          "patterns": [
            {
              "connector": "test-eventfilter-bulk-create-2-pattern"
            }
          ],
          "priority": 0,
          "enabled": false,
          "actions": [
            {
              "type": "set_field",
              "name": "connector",
              "value": "kafka_connector"
            }
          ],
          "external_data": {
            "clear": "sky",
            "type": "no",
            "arr": [
              1,
              2,
              3
            ]
          },
          "on_success": "pass",
          "on_failure": "pass"
        }
      },
      {
        "id": "test-eventfilter-bulk-create-3",
        "status": 200,
        "item": {
          "_id": "test-eventfilter-bulk-create-3",
          "type": "enrichment",
          "description": "test bulk create 3",
          "patterns": [
            null
          ],
          "priority": 0,
          "enabled": true,
          "actions": [
            {
              "from": "ExternalData.entity",
              "to": "Entity",
              "type": "copy"
            }
          ],
          "external_data": {
            "entity": {
              "type": "entity"
            }
          },
          "on_success": "pass",
          "on_failure": "pass",
          "author": "root"
        }
      },
      {
        "id": "test-eventfilter-bulk-create-4",
        "status": 200,
        "item": {
          "_id": "test-eventfilter-bulk-create-4",
          "type": "enrichment",
          "description": "test bulk create 4",
          "patterns": [
            {}
          ],
          "priority": 0,
          "enabled": true,
          "actions": [
            {
              "from": "ExternalData.entity",
              "to": "Entity",
              "type": "copy"
            }
          ],
          "external_data": {
            "entity": {
              "type": "entity"
            }
          },
          "on_success": "pass",
          "on_failure": "pass",
          "author": "root"
        }
      },
      {
        "id": "test-eventfilter-bulk-create-5",
        "status": 200,
        "item": {
          "_id": "test-eventfilter-bulk-create-5",
          "type": "enrichment",
          "description": "test bulk create 5",
          "patterns": null,
          "priority": 0,
          "enabled": true,
          "actions": [
            {
              "from": "ExternalData.entity",
              "to": "Entity",
              "type": "copy"
            }
          ],
          "external_data": {
            "entity": {
              "type": "entity"
            }
          },
          "on_success": "pass",
          "on_failure": "pass"
        }
      },
      {
        "status": 400,
        "item": [],
        "error": "value doesn't contain object; it contains array"
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
          "author": "root",
          "description": "test bulk create 1",
          "type": "enrichment",
          "patterns": [
            {
              "connector": "test-eventfilter-bulk-create-1-pattern"
            }
          ],
          "priority": 0,
          "enabled": true,
          "actions": [
            {
              "type": "set_field",
              "name": "connector",
              "value": "kafka_connector"
            }
          ],
          "external_data": {
            "clear": "sky",
            "type": "no",
            "arr": [
              1,
              2,
              3
            ]
          },
          "on_success": "pass",
          "on_failure": "pass"
        },
        {
          "_id": "test-eventfilter-bulk-create-2",
          "description": "test bulk create 2",
          "type": "enrichment",
          "patterns": [
            {
              "connector": "test-eventfilter-bulk-create-2-pattern"
            }
          ],
          "priority": 0,
          "enabled": false,
          "actions": [
            {
              "type": "set_field",
              "name": "connector",
              "value": "kafka_connector"
            }
          ],
          "external_data": {
            "clear": "sky",
            "type": "no",
            "arr": [
              1,
              2,
              3
            ]
          },
          "on_success": "pass",
          "on_failure": "pass"
        },
        {
          "_id": "test-eventfilter-bulk-create-3",
          "type": "enrichment",
          "description": "test bulk create 3",
          "patterns": [
            {}
          ],
          "priority": 0,
          "enabled": true,
          "actions": [
            {
              "from": "ExternalData.entity",
              "to": "Entity",
              "type": "copy"
            }
          ],
          "external_data": {
            "entity": {
              "type": "entity"
            }
          },
          "on_success": "pass",
          "on_failure": "pass",
          "author": "root"
        },
        {
          "_id": "test-eventfilter-bulk-create-4",
          "type": "enrichment",
          "description": "test bulk create 4",
          "patterns": [
            {}
          ],
          "priority": 0,
          "enabled": true,
          "actions": [
            {
              "from": "ExternalData.entity",
              "to": "Entity",
              "type": "copy"
            }
          ],
          "external_data": {
            "entity": {
              "type": "entity"
            }
          },
          "on_success": "pass",
          "on_failure": "pass",
          "author": "root"
        },
        {
          "_id": "test-eventfilter-bulk-create-5",
          "type": "enrichment",
          "description": "test bulk create 5",
          "patterns": null,
          "priority": 0,
          "enabled": true,
          "actions": [
            {
              "from": "ExternalData.entity",
              "to": "Entity",
              "type": "copy"
            }
          ],
          "external_data": {
            "entity": {
              "type": "entity"
            }
          },
          "on_success": "pass",
          "on_failure": "pass",
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

  Scenario: given invalid create request should return errors
    When I am admin
    When I do POST /api/v4/bulk/eventfilters:
    """json
    [
      {
        "type": "unspecified"
      },
      {
        "type": "enrichment",
        "actions": []
      },
      {
        "type": "enrichment",
        "on_failure": "continue",
        "on_success": "continue"
      },
      {
        "_id": "test-eventfilter-check-id"
      },
      {
        "type": "enrichment",
        "description": "More entity copy",
        "patterns": [
          4
        ],
        "priority": 0,
        "enabled": true,
        "actions": [
          {
            "from": "ExternalData.entity",
            "to": "Entity",
            "type": "copy"
          }
        ],
        "external_data": {
          "entity": {
            "type": "entity"
          }
        },
        "on_success": "pass",
        "on_failure": "pass",
        "author": "root"
      },
      {
        "type": "enrichment",
        "description": "Invalid pattern with empty document",
        "patterns": [
          {},
          {
            "connector": "test-eventfilter-bulk-create-1-pattern"
          }
        ],
        "priority": 0,
        "enabled": true,
        "actions": [
          {
            "from": "ExternalData.entity",
            "to": "Entity",
            "type": "copy"
          }
        ],
        "external_data": {
          "entity": {
            "type": "entity"
          }
        },
        "on_success": "pass",
        "on_failure": "pass",
        "author": "root"
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 400,
        "item": {
          "type": "unspecified"
        },
        "errors": {
          "type": "Type must be one of [break drop enrichment]."
        }
      },
      {
        "status": 400,
        "item": {
          "type": "enrichment",
          "actions": []
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
          "on_failure": "continue",
          "on_success": "continue"
        },
        "errors": {
          "on_failure": "OnFailure must be one of [pass drop break].",
          "on_success": "OnSuccess must be one of [pass drop break]."
        }
      },
      {
        "status": 400,
        "item": {
          "_id": "test-eventfilter-check-id"
        },
        "errors": {
          "_id": "ID already exists."
        }
      },
      {
        "status": 400,
        "item": {
          "type": "enrichment",
          "description": "More entity copy",
          "patterns": [
            4
          ],
          "priority": 0,
          "enabled": true,
          "actions": [
            {
              "from": "ExternalData.entity",
              "to": "Entity",
              "type": "copy"
            }
          ],
          "external_data": {
            "entity": {
              "type": "entity"
            }
          },
          "on_success": "pass",
          "on_failure": "pass",
          "author": "root"
        },
        "error": "error decoding key list: unable to parse event pattern list element"
      },
      {
        "status": 400,
        "item": {
          "type": "enrichment",
          "description": "Invalid pattern with empty document",
          "patterns": [
            {},
            {
              "connector": "test-eventfilter-bulk-create-1-pattern"
            }
          ],
          "priority": 0,
          "enabled": true,
          "actions": [
            {
              "from": "ExternalData.entity",
              "to": "Entity",
              "type": "copy"
            }
          ],
          "external_data": {
            "entity": {
              "type": "entity"
            }
          },
          "on_success": "pass",
          "on_failure": "pass",
          "author": "root"
        },
        "error": "error decoding key list: unable to parse event pattern list element"
      }
    ]
    """

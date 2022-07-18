Feature: Create an eventfilter
  I need to be able to create an eventfilter
  Only admin should be able to create an eventfilter

  Scenario: given create request with event pattern should return ok
    When I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
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
    """
    Then the response code should be 201
    When I do GET /api/v4/eventfilter/rules/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "author": "root",
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
    """

  Scenario: given create request with entity pattern should return ok
    When I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
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
    """
    Then the response code should be 201
    When I do GET /api/v4/eventfilter/rules/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "author": "root",
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
    """

  Scenario: given create request with both patterns should return ok
    When I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
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
    """
    Then the response code should be 201
    When I do GET /api/v4/eventfilter/rules/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "author": "root",
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
    """

  Scenario: given create request with corporate pattern should return ok
    When I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
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
    """
    Then the response code should be 201
    When I do GET /api/v4/eventfilter/rules/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "author": "root",
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
    """

  Scenario: given create request with absent corporate pattern should return error
    When I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
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
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "corporate_entity_pattern": "CorporateEntityPattern doesn't exist."
      }
    }
    """

  Scenario: given create request with wrong type should return bad request
    When I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
       "type": "unspecified"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "type": "Type must be one of [break drop enrichment change_entity]."
      }
    }
    """

  Scenario: given create request with missing fields should return bad request
    When I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "type": "enrichment",
      "description": "some",
      "config": {
        "actions": []
      }
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "actions": "Actions is missing.",
        "on_failure": "OnFailure is required when Type enrichment is defined.",
        "on_success": "OnSuccess is required when Type enrichment is defined."
      }
    }
    """

  Scenario: given create request with missing fields should return bad request
    When I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
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
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "on_failure": "OnFailure must be one of [pass drop break].",
        "on_success": "OnSuccess must be one of [pass drop break]."
      }
    }
    """

  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/eventfilter/rules
    Then the response code should be 401

  Scenario: given create request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/eventfilter/rules
    Then the response code should be 403

  Scenario: given create request with already exists id should return error
    When I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "_id": "test-eventfilter-check-id"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "_id": "ID already exists."
      }
    }
    """

  Scenario: create request with empty patterns should return error
    When I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
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
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "entity_pattern": "EntityPattern is invalid entity pattern.",
        "event_pattern": "EventPattern is invalid event pattern."
      }
    }
    """

  Scenario: create request with invalid event patterns should return error
    When I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
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
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "event_pattern": "EventPattern is invalid event pattern."
      }
    }
    """
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "type":"enrichment",
      "description":"Another entity copy",
      "event_pattern":[[
        {
          "field": "connector",
          "cond": {
            "type": "gt",
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
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "event_pattern": "EventPattern is invalid event pattern."
      }
    }
    """
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "type":"enrichment",
      "description":"Another entity copy",
      "event_pattern":[[
        {
          "field": "extra.test",
          "field_type": "string",
          "cond": {
            "type": "gt",
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
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "event_pattern": "EventPattern is invalid event pattern."
      }
    }
    """

  Scenario: create request with empty patterns should return error
    When I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "type":"enrichment",
      "description":"Another entity copy",
      "priority":0,
      "enabled":true,
      "config": {
        "actions":[{"value":"ExternalData.entity","name":"Entity","type":"copy"}],
        "on_success":"pass",
        "on_failure":"pass"
      },
      "external_data":{"entity":{"type":"entity"}}
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "event_pattern": "EventPattern or EntityPattern is required.",
        "entity_pattern": "EntityPattern or EventPattern is required."
      }
    }
    """

  Scenario: given POST change_entity rule requests should return error, because of empty config
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
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
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "config": "Config is missing."
      }
    }
    """
    When I do POST /api/v4/eventfilter/rules:
    """
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
      "config": {
        "component": "",
        "connector": "",
        "resource": "",
        "connector_name": ""
      },
      "enabled": true
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "config": "Config is missing."
      }
    }
    """
    When I do POST /api/v4/eventfilter/rules:
    """
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
      "config": {},
      "enabled": true
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "config": "Config is missing."
      }
    }
    """

  Scenario: given create request with set_entity_info where info value is string value should return success
    When I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "description": "test create 5",
      "type": "enrichment",
      "event_pattern": [
        [
          {
            "field": "connector",
            "cond": {
              "type": "eq",
              "value": "test-eventfilter-create-5-pattern"
            }
          }
        ]
      ],
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
      }
    }
    """
    Then the response code should be 201

  Scenario: given create request with set_entity_info where info value is int value should return success
    When I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "description": "test create 6",
      "type": "enrichment",
      "event_pattern": [
        [
          {
            "field": "connector",
            "cond": {
              "type": "eq",
              "value": "test-eventfilter-create-6-pattern"
            }
          }
        ]
      ],
      "enabled": true,
      "config": {
        "actions": [
          {
            "type": "set_field",
            "name": "connector",
            "value": 123
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      }
    }
    """
    Then the response code should be 201

  Scenario: given create request with set_entity_info where info value is bool value should return success
    When I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "description": "test create 7",
      "type": "enrichment",
      "event_pattern": [
        [
          {
            "field": "connector",
            "cond": {
              "type": "eq",
              "value": "test-eventfilter-create-7-pattern"
            }
          }
        ]
      ],
      "enabled": true,
      "config": {
        "actions": [
          {
            "type": "set_field",
            "name": "connector",
            "value": true
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      }
    }
    """
    Then the response code should be 201

  Scenario: given create request with set_entity_info where info value is array of strings value should return success
    When I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "description": "test create 8",
      "type": "enrichment",
      "event_pattern": [
        [
          {
            "field": "connector",
            "cond": {
              "type": "eq",
              "value": "test-eventfilter-create-8-pattern"
            }
          }
        ]
      ],
      "enabled": true,
      "config": {
        "actions": [
          {
            "type": "set_field",
            "name": "connector",
            "value": ["kafka_connector_1", "kafka_connector_2"]
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      }
    }
    """
    Then the response code should be 201

  Scenario: given create request with set_entity_info where info value is float should return error
    When I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "description": "test create 9",
      "type": "enrichment",
      "event_pattern": [
        [
          {
            "field": "connector",
            "cond": {
              "type": "eq",
              "value": "test-eventfilter-create-9-pattern"
            }
          }
        ]
      ],
      "enabled": true,
      "config": {
        "actions": [
          {
            "type": "set_field",
            "name": "connector",
            "value": 1.2
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      }
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "config.actions.0.value": "info value should be int, string, bool or array of strings"
      }
    }
    """

  Scenario: given create request with set_entity_info where info value is array of various types should return error
    When I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "description": "test create 10",
      "type": "enrichment",
      "event_pattern": [
        [
          {
            "field": "connector",
            "cond": {
              "type": "eq",
              "value": "test-eventfilter-create-10-pattern"
            }
          }
        ]
      ],
      "enabled": true,
      "config": {
        "actions": [
          {
            "type": "set_field",
            "name": "connector",
            "value": ["test", 1, "test 2", false]
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      }
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "config.actions.0.value": "info value should be int, string, bool or array of strings"
      }
    }
    """

  Scenario: given create request with set_entity_info where info value is object should return error
    When I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "description": "test create 11",
      "type": "enrichment",
      "event_pattern": [
        [
          {
            "field": "connector",
            "cond": {
              "type": "eq",
              "value": "test-eventfilter-create-11-pattern"
            }
          }
        ]
      ],
      "enabled": true,
      "config": {
        "actions": [
          {
            "type": "set_field",
            "name": "connector",
            "value": {
              "test": "abc",
              "test2": 1
            }
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      }
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "config.actions.0.value": "info value should be int, string, bool or array of strings"
      }
    }
    """

Feature: Update an eventfilter
  I need to be able to update an eventfilter
  Only admin should be able to update an eventfilter

  Scenario: given update request with event_pattern should update eventfilter
    When I am admin
    Then I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-update-1:
    """
    {
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
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "_id": "test-eventfilter-to-update-1",
      "author": {
        "_id": "root",
        "name": "root"
      },
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
    """

  Scenario: given update request with entity_pattern should update eventfilter
    When I am admin
    Then I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-update-2:
    """
    {
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
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "_id": "test-eventfilter-to-update-2",
      "author": {
        "_id": "root",
        "name": "root"
      },
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
    """

  Scenario: given update request with both patterns should update eventfilter
    When I am admin
    Then I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-update-3:
    """
    {
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
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "_id": "test-eventfilter-to-update-3",
      "author": {
        "_id": "root",
        "name": "root"
      },
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
    """

  Scenario: given update request with corporate pattern should update eventfilter
    When I am admin
    Then I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-update-4:
    """
    {
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
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "_id": "test-eventfilter-to-update-4",
      "author": {
        "_id": "root",
        "name": "root"
      },
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
      "enabled": true
    }
    """

  Scenario: given get request and no auth user should not allow access
    When I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-update-1
    Then the response code should be 401

  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-update-1
    Then the response code should be 403

  Scenario: given update request with not exist id should return not found error
    When I am admin
    When I do PUT /api/v4/eventfilter/rules/test-eventfilter-not-found:
    """
    {
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
    """
    Then the response code should be 404
    Then the response body should be:
    """
    {
      "error": "Not found"
    }
    """
    
  Scenario: given PUT change_entity rule requests should return error, because of empty config
    Given I am admin
    When I do PUT /api/v4/eventfilter/rules/test-update-change-entity:
    """
    {
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
    When I do PUT /api/v4/eventfilter/rules/test-update-change-entity:
    """
    {
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
    When I do PUT /api/v4/eventfilter/rules/test-update-change-entity:
    """
    {
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

  Scenario: given update request with bad pattern should return error
    When I am admin
    Then I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-update-1:
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
    Then I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-update-1:
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
    Then I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-update-1:
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

  Scenario: given update request with absent corporate pattern should return error
    When I am admin
    Then I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-update-1:
    """
    {
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

  Scenario: update request with empty patterns should return error
    When I am admin
    When I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-update-1:
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

  Scenario: given update request should update eventfilter without changes in old patterns,
            but should unset old patterns if new patterns are present
    When I am admin
    Then I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-backward-compatibility-to-update:
    """
    {
      "description": "changed description",
      "type": "drop",
      "priority": 0,
      "enabled": true
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "description": "changed description",
      "type": "drop",
      "priority": 0,
      "enabled": true,
      "old_patterns": [
        {
          "resource": {
            "regex_match": "test-eventfilter-to-backward-compatibility-to-update"
          }
        }
      ]
    }
    """
    Then I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-backward-compatibility-to-update:
    """
    {
      "description": "changed description",
      "type": "drop",
      "priority": 0,
      "enabled": true,
      "event_pattern": [
        [
          {
            "field": "resource",
            "cond": {
              "type": "eq",
              "value": "test-eventfilter-to-backward-compatibility-to-update"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "description": "changed description",
      "type": "drop",
      "priority": 0,
      "enabled": true,
      "old_patterns": null,
      "event_pattern": [
        [
          {
            "field": "resource",
            "cond": {
              "type": "eq",
              "value": "test-eventfilter-to-backward-compatibility-to-update"
            }
          }
        ]
      ]
    }
    """

  Scenario: given create request with set_entity_info where info value is string value should return success
    When I am admin
    When I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-update-5:
    """
    {
      "description": "test update 5",
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
    Then the response code should be 200

  Scenario: given create request with set_entity_info where info value is int value should return success
    When I am admin
    When I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-update-6:
    """
    {
      "description": "test update 6",
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
    Then the response code should be 200

  Scenario: given create request with set_entity_info where info value is bool value should return success
    When I am admin
    When I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-update-7:
    """
    {
      "description": "test update 6",
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
            "value": true
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      }
    }
    """
    Then the response code should be 200

  Scenario: given create request with set_entity_info where info value is array of strings value should return success
    When I am admin
    When I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-update-8:
    """
    {
      "description": "test update 8",
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
            "value": ["kafka_connector_1", "kafka_connector_2"]
          }
        ],
        "on_success": "pass",
        "on_failure": "pass"
      }
    }
    """
    Then the response code should be 200

  Scenario: given create request with set_entity_info where info value is float should return error
    When I am admin
    When I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-update-9:
    """
    {
      "description": "test update 9",
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
    When I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-update-9:
    """
    {
      "description": "test update 9",
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
    When I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-update-9:
    """
    {
      "description": "test update 9",
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

  Scenario: given update request with start and stop should return success
    When I am admin
    When I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-update-10:
    """
    {
      "description": "test update 10",
      "type": "enrichment",
      "start": 1663316803,
      "stop": 1663326803,
      "event_pattern": [
        [
          {
            "field": "connector",
            "cond": {
              "type": "eq",
              "value": "test-eventfilter-update-10-pattern"
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
    Then the response code should be 200
    When I do GET /api/v4/eventfilter/rules/test-eventfilter-to-update-10
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "description": "test update 10",
      "type": "enrichment",
      "start": 1663316803,
      "stop": 1663326803,
      "event_pattern": [
        [
          {
            "field": "connector",
            "cond": {
              "type": "eq",
              "value": "test-eventfilter-update-10-pattern"
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

  Scenario: given update request with start but without stop should return error
    When I am admin
    When I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-update-11:
    """
    {
      "description": "test update 11",
      "type": "enrichment",
      "start": 1663316803,
      "event_pattern": [
        [
          {
            "field": "connector",
            "cond": {
              "type": "eq",
              "value": "test-eventfilter-update-11-pattern"
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
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "stop": "Stop is required when Start is present."
      }
    }
    """

  Scenario: given update request with stop but without start should return error
    When I am admin
    When I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-update-12:
    """
    {
      "description": "test update 12",
      "type": "enrichment",
      "stop": 1663316803,
      "event_pattern": [
        [
          {
            "field": "connector",
            "cond": {
              "type": "eq",
              "value": "test-eventfilter-update-12-pattern"
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
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "start": "Start is required when Stop is present."
      }
    }
    """

  Scenario: given update request with stop but without start should return error
    When I am admin
    When I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-update-13:
    """
    {
      "description": "test update 13",
      "type": "enrichment",
      "start": 1663317803,
      "stop": 1663316803,
      "event_pattern": [
        [
          {
            "field": "connector",
            "cond": {
              "type": "eq",
              "value": "test-eventfilter-update-13-pattern"
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
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "stop": "Stop should be greater than Start."
      }
    }
    """

  Scenario: given update request with rrule should return success
    When I am admin
    When I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-update-14:
    """
    {
      "description": "test update 14",
      "type": "enrichment",
      "start": 1463314803,
      "stop": 1463326803,
      "rrule": "FREQ=DAILY",
      "event_pattern": [
        [
          {
            "field": "connector",
            "cond": {
              "type": "eq",
              "value": "test-eventfilter-update-14-pattern"
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
    Then the response code should be 200
    When I do GET /api/v4/eventfilter/rules/test-eventfilter-to-update-14
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "description": "test update 14",
      "type": "enrichment",
      "start": 1463314803,
      "stop": 1463326803,
      "rrule": "FREQ=DAILY",
      "event_pattern": [
        [
          {
            "field": "connector",
            "cond": {
              "type": "eq",
              "value": "test-eventfilter-update-14-pattern"
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

  Scenario: given update request with rrule but without interval should return error
    When I am admin
    When I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-update-15:
    """
    {
      "description": "test update 15",
      "type": "enrichment",
      "rrule": "FREQ=DAILY",
      "event_pattern": [
        [
          {
            "field": "connector",
            "cond": {
              "type": "eq",
              "value": "test-eventfilter-update-15-pattern"
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
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "start": "Start is required when RRule is present.",
        "stop": "Stop is required when RRule is present."
      }
    }
    """

  Scenario: given update request with invalid rrule should return error
    When I am admin
    When I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-update-16:
    """
    {
      "description": "test update 16",
      "type": "enrichment",
      "rrule": "FREQ=DAILYYY",
      "start": 1663316803,
      "stop": 1663326803,
      "event_pattern": [
        [
          {
            "field": "connector",
            "cond": {
              "type": "eq",
              "value": "test-eventfilter-update-16-pattern"
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
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "rrule": "RRule is invalid recurrence rule."
      }
    }
    """

  Scenario: given update request with exdates should return success
    When I am admin
    When I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-update-17:
    """
    {
      "description": "test update 17",
      "type": "enrichment",
      "start": 1663316803,
      "stop": 1663326803,
      "event_pattern": [
        [
          {
            "field": "connector",
            "cond": {
              "type": "eq",
              "value": "test-eventfilter-update-17-pattern"
            }
          }
        ]
      ],
      "exdates": [
        {
          "begin": 1591164001,
          "end": 1591167601
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
    """
    Then the response code should be 200
    When I do GET /api/v4/eventfilter/rules/test-eventfilter-to-update-17
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "description": "test update 17",
      "type": "enrichment",
      "start": 1663316803,
      "stop": 1663326803,
      "event_pattern": [
        [
          {
            "field": "connector",
            "cond": {
              "type": "eq",
              "value": "test-eventfilter-update-17-pattern"
            }
          }
        ]
      ],
      "exdates": [
        {
          "begin": 1591164001,
          "end": 1591167601
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
    """

  Scenario: given update request with invalid exdates should return error
    When I am admin
    When I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-update-18:
    """
    {
      "description": "test update 18",
      "type": "enrichment",
      "start": 1663316803,
      "stop": 1663326803,
      "event_pattern": [
        [
          {
            "field": "connector",
            "cond": {
              "type": "eq",
              "value": "test-eventfilter-update-18-pattern"
            }
          }
        ]
      ],
      "exdates": [
        {
          "begin": 1691164001,
          "end": 1591167601
        },
        {},
        {
          "begin": 1591167601
        },
        {
          "end": 1591167601
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
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "exdates.0.end": "End should be greater than Begin.",
        "exdates.1.begin": "Begin is missing.",
        "exdates.1.end": "End is missing.",
        "exdates.2.end": "End should be greater than Begin.",
        "exdates.3.begin": "Begin is missing."
      }
    }
    """

  Scenario: given update request with exceptions should return success
    When I am admin
    When I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-update-19:
    """
    {
      "description": "test update 19",
      "type": "enrichment",
      "start": 1663316803,
      "stop": 1663326803,
      "event_pattern": [
        [
          {
            "field": "connector",
            "cond": {
              "type": "eq",
              "value": "test-eventfilter-update-19-pattern"
            }
          }
        ]
      ],
      "exceptions":  ["test-exception-to-pbh-edit"],
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
    Then the response code should be 200
    When I do GET /api/v4/eventfilter/rules/test-eventfilter-to-update-19
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "description": "test update 19",
      "type": "enrichment",
      "start": 1663316803,
      "stop": 1663326803,
      "event_pattern": [
        [
          {
            "field": "connector",
            "cond": {
              "type": "eq",
              "value": "test-eventfilter-update-19-pattern"
            }
          }
        ]
      ],
      "exceptions": [
        {
          "_id": "test-exception-to-pbh-edit"
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
    """

  Scenario: given update request with exceptions should return success
    When I am admin
    When I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-update-20:
    """
    {
      "description": "test update 20",
      "type": "enrichment",
      "start": 1663316803,
      "stop": 1663326803,
      "event_pattern": [
        [
          {
            "field": "connector",
            "cond": {
              "type": "eq",
              "value": "test-eventfilter-update-20-pattern"
            }
          }
        ]
      ],
      "exceptions":  ["test-exception-not-exist"],
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
        "exceptions": "Exceptions doesn't exist."
      }
    }
    """

  Scenario: given update request with start and stop should return success
    When I am admin
    When I do PUT /api/v4/eventfilter/rules/test-eventfilter-to-update-21:
    """
    {
      "description": "test update 21",
      "type": "enrichment",
      "event_pattern": [
        [
          {
            "field": "connector",
            "cond": {
              "type": "eq",
              "value": "test-eventfilter-update-21-pattern"
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
    Then the response code should be 200
    When I do GET /api/v4/eventfilter/rules/test-eventfilter-to-update-21
    Then the response code should be 200
    Then the response key "start" should not exist
    Then the response key "stop" should not exist

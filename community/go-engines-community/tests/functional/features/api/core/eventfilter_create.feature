Feature: Create an eventfilter
  I need to be able to create an eventfilter
  Only admin should be able to create an eventfilter

  Scenario: given create request should return ok
    When I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "description": "test create 1",
      "type": "enrichment",
      "patterns": [
        {
          "connector": "test-eventfilter-create-1-pattern"
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
    Then the response code should be 201
    Then the response body should contain:
    """
    {
      "author": "root",
      "description": "test create 1",
      "type": "enrichment",
      "patterns": [
        {
          "connector": "test-eventfilter-create-1-pattern"
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
    When I do GET /api/v4/eventfilter/rules/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "author": "root",
      "description": "test create 1",
      "type": "enrichment",
      "patterns": [
        {
          "connector": "test-eventfilter-create-1-pattern"
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

  Scenario: given create request with enabled is false should return ok
    When I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "description": "test create 2",
      "type": "enrichment",
      "patterns": [
        {
          "connector": "test-eventfilter-create-2-pattern"
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
    """
    Then the response code should be 201
    Then the response body should contain:
    """
    {
      "author": "root",
      "description": "test create 2",
      "type": "enrichment",
      "patterns": [
        {
          "connector": "test-eventfilter-create-2-pattern"
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
            "value":"{{.ExternalData.test}}",
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

  Scenario: create request with empty patterns should return ok
    When I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "type":"enrichment",
      "description":"Another entity copy",
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
    """
    Then the response code should be 201
    Then the response body should contain:
    """
    {
      "type":"enrichment",
      "description":"Another entity copy",
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
    }
    """
    When I do GET /api/v4/eventfilter/rules/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "type":"enrichment",
      "description":"Another entity copy",
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
    }
    """
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "type":"enrichment",
      "description":"More entity copy",
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
    """
    Then the response code should be 201
    Then the response body should contain:
    """
    {
      "type":"enrichment",
      "description":"More entity copy",
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
    """
    When I do GET /api/v4/eventfilter/rules/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "type":"enrichment",
      "description":"More entity copy",
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
    """
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "type":"enrichment",
      "description":"More entity copy",
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
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "error":"request has invalid structure"
    }
    """
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "type":"enrichment",
      "description":"Invalid pattern with empty document",
      "patterns":[{},{"connector": "test-eventfilter-create-1-pattern"}],
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
      "error":"request has invalid structure"
    }
    """
  Scenario: given POST change_entity rule requests should return error, because of empty config
    Given I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
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

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
        "type": "Type must be one of [break drop enrichment]."
      }
    }
    """

  Scenario: given create request with missing fields should return bad request
    When I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
       "type": "enrichment",
       "actions": []
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
      "on_failure": "continue",
      "on_success": "continue"
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

  Scenario: create request with patterns: [null] should return ok
    When I am admin
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "type":"enrichment","description":"Recopie entite","patterns":[null],
      "priority":0,"enabled":true,"actions":[{"from":"ExternalData.entity","to":"Entity","type":"copy"}],
      "external_data":{"entity":{"type":"entity"}},"on_success":"pass","on_failure":"pass","author":"root"
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """
    {
      "type":"enrichment","description":"Recopie entite","patterns":[ {} ],
      "priority":0,"enabled":true,"actions":[{"from":"ExternalData.entity","to":"Entity","type":"copy"}],
      "external_data":{"entity":{"type":"entity"}},"on_success":"pass","on_failure":"pass","author":"root"
    }
    """
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "type":"enrichment","description":"Another entity copy","patterns":[ { } ],
      "priority":0,"enabled":true,"actions":[{"from":"ExternalData.entity","to":"Entity","type":"copy"}],
      "external_data":{"entity":{"type":"entity"}},"on_success":"pass","on_failure":"pass","author":"root"
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """
    {
      "type":"enrichment","description":"Another entity copy","patterns":[ { } ],
      "priority":0,"enabled":true,"actions":[{"from":"ExternalData.entity","to":"Entity","type":"copy"}],
      "external_data":{"entity":{"type":"entity"}},"on_success":"pass","on_failure":"pass","author":"root"
    }
    """
    When I do GET /api/v4/eventfilter/rules?search=Another%20entity%20copy
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "type":"enrichment","description":"Another entity copy","patterns":[ { } ],
          "priority":0,"enabled":true,"actions":[{"from":"ExternalData.entity","to":"Entity","type":"copy"}],
          "external_data":{"entity":{"type":"entity"}},"on_success":"pass","on_failure":"pass","author":"root"
        }
      ]
    }
    """
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "type":"enrichment","description":"More entity copy","patterns":null,
      "priority":0,"enabled":true,"actions":[{"from":"ExternalData.entity","to":"Entity","type":"copy"}],
      "external_data":{"entity":{"type":"entity"}},"on_success":"pass","on_failure":"pass"
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """
    {
      "type":"enrichment","description":"More entity copy","patterns":null,
      "priority":0,"enabled":true,"actions":[{"from":"ExternalData.entity","to":"Entity","type":"copy"}],
      "external_data":{"entity":{"type":"entity"}},"on_success":"pass","on_failure":"pass","author":"root"
    }
    """
    When I do GET /api/v4/eventfilter/rules?search=More%20entity%20copy
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "type":"enrichment","description":"More entity copy","patterns":null,
          "priority":0,"enabled":true,"actions":[{"from":"ExternalData.entity","to":"Entity","type":"copy"}],
          "external_data":{"entity":{"type":"entity"}},"on_success":"pass","on_failure":"pass","author":"root"
        }
      ]
    }
    """
    When I do POST /api/v4/eventfilter/rules:
    """
    {
      "type":"enrichment","description":"More entity copy","patterns":[4],
      "priority":0,"enabled":true,"actions":[{"from":"ExternalData.entity","to":"Entity","type":"copy"}],
      "external_data":{"entity":{"type":"entity"}},"on_success":"pass","on_failure":"pass","author":"root"
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
      "type":"enrichment","description":"Invalid pattern with empty document",
      "patterns":[
        {},{"connector": "test-eventfilter-create-1-pattern"}],
      "priority":0,"enabled":true,"actions":[{"from":"ExternalData.entity","to":"Entity","type":"copy"}],
      "external_data":{"entity":{"type":"entity"}},"on_success":"pass","on_failure":"pass","author":"root"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "error":"request has invalid structure"
    }
    """

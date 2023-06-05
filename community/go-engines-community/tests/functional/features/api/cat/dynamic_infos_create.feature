Feature: Create an dynamic-infos
  I need to be able to create a dynamic-infos
  Only admin should be able to create a dynamic-infos

  Scenario: given create request should return ok
    When I am admin
    When I do POST /api/v4/cat/dynamic-infos:
    """json
    {
      "entity_pattern": [
        [
          {
            "field": "infos.alert_name",
            "field_type": "string",
            "cond": {
              "type": "eq",
              "value": "test-dynamic-infos-to-create-1-pattern"
            }
          }
        ]
      ],
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-dynamic-infos-to-create-1-pattern"
            }
          }
        ]
      ],
      "name": "test-dynamic-infos-to-create-1-name",
      "description": "test-dynamic-infos-to-create-1-description",
      "enabled": true,
      "infos": [
        {
          "name": "test-dynamic-infos-to-create-1-info-1-name",
          "value": "test-dynamic-infos-to-create-1-info-1-value"
        },
        {
          "name": "test-dynamic-infos-to-create-1-info-2-name",
          "value": "test-dynamic-infos-to-create-1-info-2-value"
        }
      ]
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "entity_pattern": [
        [
          {
            "field": "infos.alert_name",
            "field_type": "string",
            "cond": {
              "type": "eq",
              "value": "test-dynamic-infos-to-create-1-pattern"
            }
          }
        ]
      ],
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-dynamic-infos-to-create-1-pattern"
            }
          }
        ]
      ],
      "name": "test-dynamic-infos-to-create-1-name",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "description": "test-dynamic-infos-to-create-1-description",
      "enabled": true,
      "infos": [
        {
          "name": "test-dynamic-infos-to-create-1-info-1-name",
          "value": "test-dynamic-infos-to-create-1-info-1-value"
        },
        {
          "name": "test-dynamic-infos-to-create-1-info-2-name",
          "value": "test-dynamic-infos-to-create-1-info-2-value"
        }
      ]
    }
    """
    When I do GET /api/v4/cat/dynamic-infos/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "entity_pattern": [
        [
          {
            "field": "infos.alert_name",
            "field_type": "string",
            "cond": {
              "type": "eq",
              "value": "test-dynamic-infos-to-create-1-pattern"
            }
          }
        ]
      ],
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-dynamic-infos-to-create-1-pattern"
            }
          }
        ]
      ],
      "name": "test-dynamic-infos-to-create-1-name",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "description": "test-dynamic-infos-to-create-1-description",
      "enabled": true,
      "infos": [
        {
          "name": "test-dynamic-infos-to-create-1-info-1-name",
          "value": "test-dynamic-infos-to-create-1-info-1-value"
        },
        {
          "name": "test-dynamic-infos-to-create-1-info-2-name",
          "value": "test-dynamic-infos-to-create-1-info-2-value"
        }
      ]
    }
    """

  Scenario: given new rule should return dynamic infos by pattern search request
    When I am admin
    When I do POST /api/v4/cat/dynamic-infos:
    """json
    {
      "name": "test-dynamic-infos-to-create-2-name",
      "entity_pattern": [
        [
          {
            "field": "infos.alert_name",
            "field_type": "string",
            "cond": {
              "type": "regexp",
              "value": "test-dynamic-infos-to-create-2-entity-pattern"
            }
          }
        ]
      ],
      "alarm_pattern": [
        [
          {
            "field": "v.connector",
            "cond": {
              "type": "regexp",
              "value": "test-dynamic-infos-to-create-2-alarm-pattern"
            }
          }
        ]
      ],
      "description": "test-dynamic-infos-to-create-2-description",
      "enabled": true,
      "infos": [
        {
          "name": "test-dynamic-infos-to-create-2-info-1-name",
          "value": "test-dynamic-infos-to-create-2-info-1-value"
        },
        {
          "name": "test-dynamic-infos-to-create-2-info-2-name",
          "value": "test-dynamic-infos-to-create-2-info-2-value"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/cat/dynamic-infos?search=pattern%20LIKE%20"infos-to-create-2-entity-pattern"
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-dynamic-infos-to-create-2-name"
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
    When I do GET /api/v4/cat/dynamic-infos?search=pattern%20LIKE%20"infos-to-create-2-alarm-pattern"
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "name": "test-dynamic-infos-to-create-2-name"
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """

  Scenario: given invalid create request should return bad request
    When I am admin
    When I do POST /api/v4/cat/dynamic-infos:
    """json
    {}
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "alarm_pattern": "AlarmPattern or EntityPattern is required.",
        "entity_pattern": "EntityPattern or AlarmPattern is required.",
        "description": "Description is missing.",
        "name": "Name is missing.",
        "infos": "Infos is missing.",
        "enabled": "Enabled is missing."
      }
    }
    """
    When I do POST /api/v4/cat/dynamic-infos:
    """json
    {
      "infos": []
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "infos": "Infos should not be blank."
      }
    }
    """
    When I do POST /api/v4/cat/dynamic-infos:
    """json
    {
      "infos": [
        {},
        {"name": "test-name"}
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "infos.0.name": "Name is missing.",
        "infos.1.value": "Value is missing."
      }
    }
    """

  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/cat/dynamic-infos
    Then the response code should be 401

  Scenario: given create request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/cat/dynamic-infos
    Then the response code should be 403

  Scenario: given create request with already exists id should return error
    When I am admin
    When I do POST /api/v4/cat/dynamic-infos:
    """json
    {
      "_id": "test-dynamic-infos-to-check-unique"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "_id": "ID already exists."
      }
    }
    """

  Scenario: given create request where info value is int value should return success
    When I am admin
    When I do POST /api/v4/cat/dynamic-infos:
    """json
    {
      "entity_pattern": [
        [
          {
            "field": "infos.alert_name",
            "field_type": "string",
            "cond": {
              "type": "eq",
              "value": "test-dynamic-infos-to-create-3-pattern"
            }
          }
        ]
      ],
      "name": "test-dynamic-infos-to-create-3-name",
      "description": "test-dynamic-infos-to-create-3-description",
      "enabled": true,
      "infos": [
        {
          "name": "test-dynamic-infos-to-create-3-info-1-name",
          "value": 123
        }
      ]
    }
    """
    Then the response code should be 201

  Scenario: given create request where info value is string value should return success
    When I am admin
    When I do POST /api/v4/cat/dynamic-infos:
    """json
    {
      "entity_pattern": [
        [
          {
            "field": "infos.alert_name",
            "field_type": "string",
            "cond": {
              "type": "eq",
              "value": "test-dynamic-infos-to-create-4-pattern"
            }
          }
        ]
      ],
      "name": "test-dynamic-infos-to-create-4-name",
      "description": "test-dynamic-infos-to-create-4-description",
      "enabled": true,
      "infos": [
        {
          "name": "test-dynamic-infos-to-create-4-info-1-name",
          "value": "test"
        }
      ]
    }
    """
    Then the response code should be 201

  Scenario: given create request where info value is bool value should return success
    When I am admin
    When I do POST /api/v4/cat/dynamic-infos:
    """json
    {
      "entity_pattern": [
        [
          {
            "field": "infos.alert_name",
            "field_type": "string",
            "cond": {
              "type": "eq",
              "value": "test-dynamic-infos-to-create-5-pattern"
            }
          }
        ]
      ],
      "name": "test-dynamic-infos-to-create-5-name",
      "description": "test-dynamic-infos-to-create-5-description",
      "enabled": true,
      "infos": [
        {
          "name": "test-dynamic-infos-to-create-5-info-1-name",
          "value": false
        }
      ]
    }
    """
    Then the response code should be 201

  Scenario: given create request where info value is array of strings value should return success
    When I am admin
    When I do POST /api/v4/cat/dynamic-infos:
    """json
    {
      "entity_pattern": [
        [
          {
            "field": "infos.alert_name",
            "field_type": "string",
            "cond": {
              "type": "eq",
              "value": "test-dynamic-infos-to-create-6-pattern"
            }
          }
        ]
      ],
      "name": "test-dynamic-infos-to-create-6-name",
      "description": "test-dynamic-infos-to-create-6-description",
      "enabled": true,
      "infos": [
        {
          "name": "test-dynamic-infos-to-create-6-info-1-name",
          "value": ["test", "test2"]
        }
      ]
    }
    """
    Then the response code should be 201

  Scenario: given create request where info value is float should return error
    When I am admin
    When I do POST /api/v4/cat/dynamic-infos:
    """json
    {
      "entity_pattern": [
        [
          {
            "field": "infos.alert_name",
            "field_type": "string",
            "cond": {
              "type": "eq",
              "value": "test-dynamic-infos-to-create-7-pattern"
            }
          }
        ]
      ],
      "name": "test-dynamic-infos-to-create-7-name",
      "description": "test-dynamic-infos-to-create-7-description",
      "enabled": true,
      "infos": [
        {
          "name": "test-dynamic-infos-to-create-7-info-1-name",
          "value": 1.2
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "infos.0.value": "info value should be int, string, bool or array of strings"
      }
    }
    """

  Scenario: given create request where info value is array of various types should return error
    When I am admin
    When I do POST /api/v4/cat/dynamic-infos:
    """json
    {
      "entity_pattern": [
        [
          {
            "field": "infos.alert_name",
            "field_type": "string",
            "cond": {
              "type": "eq",
              "value": "test-dynamic-infos-to-create-8-pattern"
            }
          }
        ]
      ],
      "name": "test-dynamic-infos-to-create-8-name",
      "description": "test-dynamic-infos-to-create-8-description",
      "enabled": true,
      "infos": [
        {
          "name": "test-dynamic-infos-to-create-8-info-1-name",
          "value": ["test 1", 1, "test 2"]
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "infos.0.value": "info value should be int, string, bool or array of strings"
      }
    }
    """

  Scenario: given create request where info value is object should return error
    When I am admin
    When I do POST /api/v4/cat/dynamic-infos:
    """json
    {
      "entity_pattern": [
        [
          {
            "field": "infos.alert_name",
            "field_type": "string",
            "cond": {
              "type": "eq",
              "value": "test-dynamic-infos-to-create-9-pattern"
            }
          }
        ]
      ],
      "name": "test-dynamic-infos-to-create-9-name",
      "description": "test-dynamic-infos-to-create-9-description",
      "enabled": true,
      "infos": [
        {
          "name": "test-dynamic-infos-to-create-9-info-1-name",
          "value": {
            "test": "abc",
            "test2": 1
          }
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "infos.0.value": "info value should be int, string, bool or array of strings"
      }
    }
    """

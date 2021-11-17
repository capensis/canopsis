Feature: Create an dynamic-infos
  I need to be able to create a dynamic-infos
  Only admin should be able to create a dynamic-infos

  Scenario: given create request should return ok
    When I am admin
    When I do POST /api/v4/cat/dynamic-infos:
    """json
    {
      "entity_patterns": [
        {
          "infos": {
            "alert_name": {
              "value": {
                "regex_match": "test-dynamic-infos-to-create-1-pattern"
              }
            }
          }
        }
      ],
      "name": "test-dynamic-infos-to-create-1-name",
      "alarm_patterns": null,
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
      "entity_patterns": [
        {
          "infos": {
            "alert_name": {
              "value": {
                "regex_match": "test-dynamic-infos-to-create-1-pattern"
              }
            }
          }
        }
      ],
      "name": "test-dynamic-infos-to-create-1-name",
      "author": "root",
      "alarm_patterns": null,
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
      "entity_patterns": [
        {
          "infos": {
            "alert_name": {
              "value": {
                "regex_match": "test-dynamic-infos-to-create-1-pattern"
              }
            }
          }
        }
      ],
      "name": "test-dynamic-infos-to-create-1-name",
      "author": "root",
      "alarm_patterns": null,
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

  Scenario: given search DSL request should return dynamic infos
    When I am admin
    When I do POST /api/v4/cat/dynamic-infos:
    """json
    {
      "entity_patterns": [
        {
          "infos": {
            "alert_name": {
              "value": {
                "regex_match": "test-dynamic-infos-to-create-2-entity-pattern"
              }
            }
          }
        }
      ],
      "name": "test-dynamic-infos-to-create-2-name",
      "alarm_patterns": [
        {
          "v": {
            "connector": {
              "regex_match": "test-dynamic-infos-to-create-2-alarm-pattern"
            }
          }
        }
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
        "description": "Description is missing.",
        "name": "Name is missing.",
        "enabled": "Enabled is missing."
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
    When I do POST /api/v4/cat/dynamic-infos:
    """json
    {
      "entity_patterns": [
        {"not-exist-field": "test-value"}
      ],
      "alarm_patterns": [
        {"not-exist-field": "test-value"}
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "alarm_patterns": "Invalid alarm patterns.",
        "entity_patterns": "Invalid entity patterns."
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

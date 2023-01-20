Feature: Update a dynamic infos
  I need to be able to update a dynamic infos
  Only admin should be able to update a dynamic infos

  Scenario: given update request should update dynamic infos
    When I am admin
    Then I do PUT /api/v4/cat/dynamic-infos/test-dynamic-infos-to-update-1:
    """json
    {
      "alarm_pattern": [
        [
          {
            "field": "v.connector",
            "cond": {
              "type": "eq",
              "value": "test-dynamic-infos-to-update-1-alarm-pattern-updated"
            }
          }
        ]
      ],
      "entity_pattern": [
        [
          {
            "field": "_id",
            "cond": {
              "type": "eq",
              "value": "test-dynamic-infos-to-update-1-alarm-pattern-updated"
            }
          }
        ]
      ],
      "name": "test-dynamic-infos-to-update-1-name-updated",
      "description": "test-dynamic-infos-to-update-1-description-updated",
      "enabled": true,
      "infos": [
        {
          "name": "test-dynamic-infos-to-update-1-info-3-name",
          "value": "test-dynamic-infos-to-update-1-info-3-value"
        },
        {
          "name": "test-dynamic-infos-to-update-1-info-2-name",
          "value": "test-dynamic-infos-to-update-1-info-2-value-updated"
        }
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-dynamic-infos-to-update-1",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "created": 1581423405,
      "description": "test-dynamic-infos-to-update-1-description-updated",
      "disable_during_periods": null,
      "enabled": true,
      "alarm_pattern": [
        [
          {
            "field": "v.connector",
            "cond": {
              "type": "eq",
              "value": "test-dynamic-infos-to-update-1-alarm-pattern-updated"
            }
          }
        ]
      ],
      "entity_pattern": [
        [
          {
            "field": "_id",
            "cond": {
              "type": "eq",
              "value": "test-dynamic-infos-to-update-1-alarm-pattern-updated"
            }
          }
        ]
      ],
      "infos": [
        {
          "name": "test-dynamic-infos-to-update-1-info-3-name",
          "value": "test-dynamic-infos-to-update-1-info-3-value"
        },
        {
          "name": "test-dynamic-infos-to-update-1-info-2-name",
          "value": "test-dynamic-infos-to-update-1-info-2-value-updated"
        }
      ],
      "name": "test-dynamic-infos-to-update-1-name-updated"
    }
    """

  Scenario: given updated rule should return dynamic infos by pattern search request
    When I am admin
    Then I do PUT /api/v4/cat/dynamic-infos/test-dynamic-infos-to-update-2:
    """json
    {
      "name": "test-dynamic-infos-to-update-2-name-updated",
      "alarm_patterns": null,
      "description": "test-dynamic-infos-to-update-2-description-updated",
      "enabled": true,
      "infos": [
        {
          "name": "test-dynamic-infos-to-update-2-info-3-name",
          "value": "test-dynamic-infos-to-update-2-info-3-value"
        },
        {
          "name": "test-dynamic-infos-to-update-2-info-2-name",
          "value": "test-dynamic-infos-to-update-2-info-2-value-updated"
        }
      ]
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/cat/dynamic-infos?search=pattern%20LIKE%20"test-dynamic-infos-to-update-2-entity-pattern"
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-dynamic-infos-to-update-2",
          "old_alarm_patterns": [
            {
              "v": {
                "connector": "test-dynamic-infos-to-update-2-alarm-pattern"
              }
            }
          ],
          "old_entity_patterns": [
            {
              "_id": "test-dynamic-infos-to-update-2-entity-pattern"
            }
          ]
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
    Then I do PUT /api/v4/cat/dynamic-infos/test-dynamic-infos-to-update-2:
    """json
    {
      "name": "test-dynamic-infos-to-update-2-name-updated",
      "alarm_patterns": null,
      "description": "test-dynamic-infos-to-update-2-description-updated",
      "enabled": true,
      "alarm_pattern": [
        [
          {
            "field": "v.connector",
            "cond": {
              "type": "eq",
              "value": "test-dynamic-infos-to-update-2-alarm-pattern-updated"
            }
          }
        ]
      ],
      "entity_pattern": [
        [
          {
            "field": "_id",
            "cond": {
              "type": "eq",
              "value": "test-dynamic-infos-to-update-2-entity-pattern-updated"
            }
          }
        ]
      ],
      "infos": [
        {
          "name": "test-dynamic-infos-to-update-2-info-3-name",
          "value": "test-dynamic-infos-to-update-2-info-3-value"
        },
        {
          "name": "test-dynamic-infos-to-update-2-info-2-name",
          "value": "test-dynamic-infos-to-update-2-info-2-value-updated"
        }
      ]
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/cat/dynamic-infos?search=pattern%20LIKE%20"test-dynamic-infos-to-update-2-entity-pattern-updated"
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-dynamic-infos-to-update-2",
          "alarm_pattern": [
            [
              {
                "field": "v.connector",
                "cond": {
                  "type": "eq",
                  "value": "test-dynamic-infos-to-update-2-alarm-pattern-updated"
                }
              }
            ]
          ],
          "entity_pattern": [
            [
              {
                "field": "_id",
                "cond": {
                  "type": "eq",
                  "value": "test-dynamic-infos-to-update-2-entity-pattern-updated"
                }
              }
            ]
          ],
          "old_alarm_patterns": null,
          "old_entity_patterns": null
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

  Scenario: given invalid update request should return bad request
    When I am admin
    When I do PUT /api/v4/cat/dynamic-infos/test-dynamic-infos-not-found:
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
    When I do PUT /api/v4/cat/dynamic-infos/test-dynamic-infos-not-found:
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

  Scenario: given get request and no auth user should not allow access
    When I do PUT /api/v4/cat/dynamic-infos/test-dynamic-infos-to-update
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/cat/dynamic-infos/test-dynamic-infos-to-update
    Then the response code should be 403

  Scenario: given update request with not exist id should return not found error
    When I am admin
    When I do PUT /api/v4/cat/dynamic-infos/test-dynamic-infos-not-found:
    """json
    {
      "entity_pattern": [
        [
          {
            "field": "infos.alert_name",
            "field_type": "string",
            "cond": {
              "type": "eq",
              "value": "test-dynamic-infos-not-found-pattern"
            }
          }
        ]
      ],
      "name": "test-dynamic-infos-not-found-name",
      "description": "test-dynamic-infos-not-found-description",
      "enabled": true,
      "infos": [
        {
          "name": "test-dynamic-infos-not-found-info-1-name",
          "value": "test-dynamic-infos-not-found-info-1-value"
        }
      ]
    }
    """
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """

  Scenario: given create request where info value is int value should return success
    When I am admin
    When I do PUT /api/v4/cat/dynamic-infos/test-dynamic-infos-to-update-3:
    """json
    {
      "entity_pattern": [
        [
          {
            "field": "infos.alert_name",
            "field_type": "string",
            "cond": {
              "type": "eq",
              "value": "test-dynamic-infos-to-update-3-pattern"
            }
          }
        ]
      ],
      "name": "test-dynamic-infos-to-update-3-name",
      "description": "test-dynamic-infos-to-update-3-description",
      "enabled": true,
      "infos": [
        {
          "name": "test-dynamic-infos-to-update-3-info-1-name",
          "value": 123
        }
      ]
    }
    """
    Then the response code should be 200

  Scenario: given create request where info value is string value should return success
    When I am admin
    When I do PUT /api/v4/cat/dynamic-infos/test-dynamic-infos-to-update-4:
    """json
    {
      "entity_pattern": [
        [
          {
            "field": "infos.alert_name",
            "field_type": "string",
            "cond": {
              "type": "eq",
              "value": "test-dynamic-infos-to-update-4-pattern"
            }
          }
        ]
      ],
      "name": "test-dynamic-infos-to-update-4-name",
      "description": "test-dynamic-infos-to-update-4-description",
      "enabled": true,
      "infos": [
        {
          "name": "test-dynamic-infos-to-update-4-info-1-name",
          "value": "test"
        }
      ]
    }
    """
    Then the response code should be 200

  Scenario: given create request where info value is bool value should return success
    When I am admin
    When I do PUT /api/v4/cat/dynamic-infos/test-dynamic-infos-to-update-5:
    """json
    {
      "entity_pattern": [
        [
          {
            "field": "infos.alert_name",
            "field_type": "string",
            "cond": {
              "type": "eq",
              "value": "test-dynamic-infos-to-update-5-pattern"
            }
          }
        ]
      ],
      "name": "test-dynamic-infos-to-update-5-name",
      "description": "test-dynamic-infos-to-update-5-description",
      "enabled": true,
      "infos": [
        {
          "name": "test-dynamic-infos-to-update-5-info-1-name",
          "value": false
        }
      ]
    }
    """
    Then the response code should be 200

  Scenario: given create request where info value is array of strings value should return success
    When I am admin
    When I do PUT /api/v4/cat/dynamic-infos/test-dynamic-infos-to-update-6:
    """json
    {
      "entity_pattern": [
        [
          {
            "field": "infos.alert_name",
            "field_type": "string",
            "cond": {
              "type": "eq",
              "value": "test-dynamic-infos-to-update-6-pattern"
            }
          }
        ]
      ],
      "name": "test-dynamic-infos-to-update-6-name",
      "description": "test-dynamic-infos-to-update-6-description",
      "enabled": true,
      "infos": [
        {
          "name": "test-dynamic-infos-to-update-6-info-1-name",
          "value": ["test", "test2"]
        }
      ]
    }
    """
    Then the response code should be 200

  Scenario: given create request where info value is float should return error
    When I am admin
    When I do PUT /api/v4/cat/dynamic-infos/test-dynamic-infos-to-update-7:
    """json
    {
      "entity_pattern": [
        [
          {
            "field": "infos.alert_name",
            "field_type": "string",
            "cond": {
              "type": "eq",
              "value": "test-dynamic-infos-to-update-7-pattern"
            }
          }
        ]
      ],
      "name": "test-dynamic-infos-to-update-7-name",
      "description": "test-dynamic-infos-to-update-7-description",
      "enabled": true,
      "infos": [
        {
          "name": "test-dynamic-infos-to-update-7-info-1-name",
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
    When I do PUT /api/v4/cat/dynamic-infos/test-dynamic-infos-to-update-7:
    """json
    {
      "entity_pattern": [
        [
          {
            "field": "infos.alert_name",
            "field_type": "string",
            "cond": {
              "type": "eq",
              "value": "test-dynamic-infos-to-update-7-pattern"
            }
          }
        ]
      ],
      "name": "test-dynamic-infos-to-update-7-name",
      "description": "test-dynamic-infos-to-update-7-description",
      "enabled": true,
      "infos": [
        {
          "name": "test-dynamic-infos-to-update-7-info-1-name",
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
    When I do PUT /api/v4/cat/dynamic-infos/test-dynamic-infos-to-update-7:
    """json
    {
      "entity_pattern": [
        [
          {
            "field": "infos.alert_name",
            "field_type": "string",
            "cond": {
              "type": "eq",
              "value": "test-dynamic-infos-to-update-7-pattern"
            }
          }
        ]
      ],
      "name": "test-dynamic-infos-to-update-7-name",
      "description": "test-dynamic-infos-to-update-7-description",
      "enabled": true,
      "infos": [
        {
          "name": "test-dynamic-infos-to-update-7-info-1-name",
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

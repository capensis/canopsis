Feature: Update a dynamic infos
  I need to be able to update a dynamic infos
  Only admin should be able to update a dynamic infos

  Scenario: given update request should update dynamic infos
    When I am admin
    Then I do PUT /api/v4/cat/dynamic-infos/test-dynamic-infos-to-update-1:
    """json
    {
      "entity_patterns": [
        {
          "infos": {
            "alert_name": {
              "value": {
                "regex_match": "test-dynamic-infos-to-update-1-pattern-updated"
              }
            }
          }
        }
      ],
      "name": "test-dynamic-infos-to-update-1-name-updated",
      "alarm_patterns": null,
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
      "alarm_patterns": null,
      "author": "root",
      "creation_date": 1581423405,
      "description": "test-dynamic-infos-to-update-1-description-updated",
      "disable_during_periods": null,
      "enabled": true,
      "entity_patterns": [
        {
          "infos": {
            "alert_name": {
              "value": {
                "regex_match": "test-dynamic-infos-to-update-1-pattern-updated"
              }
            }
          }
        }
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

  Scenario: given search DSL request should return dynamic infos
    When I am admin
    Then I do PUT /api/v4/cat/dynamic-infos/test-dynamic-infos-to-update-2:
    """json
    {
      "entity_patterns": [
        {
          "infos": {
            "alert_name": {
              "value": {
                "regex_match": "test-dynamic-infos-to-update-2-pattern-updated"
              }
            }
          }
        }
      ],
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
    When I do GET /api/v4/cat/dynamic-infos?search=pattern%20LIKE%20"test-dynamic-infos-to-update-2-pattern-updated"
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-dynamic-infos-to-update-2"
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
        "description": "Description is missing.",
        "name": "Name is missing.",
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
    When I do PUT /api/v4/cat/dynamic-infos/test-dynamic-infos-not-found:
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
      "entity_patterns": [
        {
          "infos": {
            "alert_name": {
              "value": {
                "regex_match": "test-dynamic-infos-not-found-pattern"
              }
            }
          }
        }
      ],
      "name": "test-dynamic-infos-not-found-name",
      "alarm_patterns": null,
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

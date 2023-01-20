Feature: Update entity service
  I need to be able to update a entity service

  Scenario: given update request should update entity
    When I am admin
    When I do PUT /api/v4/entityservices/test-entityservice-to-update-1:
    """json
    {
      "name": "test-entityservice-to-update-1-name",
      "output_template": "test-entityservice-to-update-1-output-updated",
      "category": "test-category-to-entityservice-edit",
      "impact_level": 2,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-entityservice-to-update-1-pattern-updated"
            }
          }
        ]
      ],
      "sli_avail_state": 1,
      "infos": [
        {
          "description": "test-entityservice-to-update-1-info-1-description",
          "name": "test-entityservice-to-update-1-info-1-name",
          "value": "test-entityservice-to-update-1-info-1-value"
        },
        {
          "description": "test-entityservice-to-update-1-info-2-description",
          "name": "test-entityservice-to-update-1-info-2-name",
          "value": false
        },
        {
          "description": "test-entityservice-to-update-1-info-3-description",
          "name": "test-entityservice-to-update-1-info-3-name",
          "value": 1022
        },
        {
          "description": "test-entityservice-to-update-1-info-4-description",
          "name": "test-entityservice-to-update-1-info-4-name",
          "value": 10.45
        },
        {
          "description": "test-entityservice-to-update-1-info-5-description",
          "name": "test-entityservice-to-update-1-info-5-name",
          "value": null
        },
        {
          "description": "test-entityservice-to-update-1-info-6-description",
          "name": "test-entityservice-to-update-1-info-6-name",
          "value": ["test-entityservice-to-update-1-info-6-value", false, 1022, 10.45, null]
        }
      ],
      "coordinates": {
        "lat": 62.34960927573042,
        "lng": 74.02834455685206
      }
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-entityservice-to-update-1",
      "category": {
        "_id": "test-category-to-entityservice-edit",
        "name": "test-category-to-entityservice-edit-name"
      },
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-entityservice-to-update-1-pattern-updated"
            }
          }
        ]
      ],
      "old_entity_patterns": null,
      "impact_level": 2,
      "infos": {
        "test-entityservice-to-update-1-info-1-name": {
          "description": "test-entityservice-to-update-1-info-1-description",
          "name": "test-entityservice-to-update-1-info-1-name",
          "value": "test-entityservice-to-update-1-info-1-value"
        },
        "test-entityservice-to-update-1-info-2-name": {
          "description": "test-entityservice-to-update-1-info-2-description",
          "name": "test-entityservice-to-update-1-info-2-name",
          "value": false
        },
        "test-entityservice-to-update-1-info-3-name": {
          "description": "test-entityservice-to-update-1-info-3-description",
          "name": "test-entityservice-to-update-1-info-3-name",
          "value": 1022
        },
        "test-entityservice-to-update-1-info-4-name": {
          "description": "test-entityservice-to-update-1-info-4-description",
          "name": "test-entityservice-to-update-1-info-4-name",
          "value": 10.45
        },
        "test-entityservice-to-update-1-info-5-name": {
          "description": "test-entityservice-to-update-1-info-5-description",
          "name": "test-entityservice-to-update-1-info-5-name",
          "value": null
        },
        "test-entityservice-to-update-1-info-6-name": {
          "description": "test-entityservice-to-update-1-info-6-description",
          "name": "test-entityservice-to-update-1-info-6-name",
          "value": ["test-entityservice-to-update-1-info-6-value", false, 1022, 10.45, null]
        }
      },
      "name": "test-entityservice-to-update-1-name",
      "output_template": "test-entityservice-to-update-1-output-updated",
      "sli_avail_state": 1,
      "type": "service",
      "coordinates": {
        "lat": 62.34960927573042,
        "lng": 74.02834455685206
      }
    }
    """

  Scenario: given update request with old pattern should update entity
    When I am admin
    When I do PUT /api/v4/entityservices/test-entityservice-to-update-2:
    """json
    {
      "name": "test-entityservice-to-update-2-name",
      "output_template": "test-entityservice-to-update-2-output-updated",
      "category": "test-category-to-entityservice-edit",
      "impact_level": 2,
      "enabled": true,
      "sli_avail_state": 1,
      "infos": []
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-entityservice-to-update-2",
      "category": {
        "_id": "test-category-to-entityservice-edit",
        "name": "test-category-to-entityservice-edit-name"
      },
      "enabled": true,
      "old_entity_patterns": [{"name": "test-entityservice-to-update-2-pattern"}],
      "impact_level": 2,
      "infos": {},
      "name": "test-entityservice-to-update-2-name",
      "output_template": "test-entityservice-to-update-2-output-updated",
      "sli_avail_state": 1,
      "type": "service"
    }
    """
    When I do PUT /api/v4/entityservices/test-entityservice-to-update-2:
    """json
    {
      "name": "test-entityservice-to-update-2-name",
      "output_template": "test-entityservice-to-update-2-output-updated",
      "category": "test-category-to-entityservice-edit",
      "impact_level": 2,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-entityservice-to-update-2-pattern-updated"
            }
          }
        ]
      ],
      "sli_avail_state": 1,
      "infos": []
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-entityservice-to-update-2",
      "category": {
        "_id": "test-category-to-entityservice-edit",
        "name": "test-category-to-entityservice-edit-name"
      },
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-entityservice-to-update-2-pattern-updated"
            }
          }
        ]
      ],
      "old_entity_patterns": null,
      "impact_level": 2,
      "infos": {},
      "name": "test-entityservice-to-update-2-name",
      "output_template": "test-entityservice-to-update-2-output-updated",
      "sli_avail_state": 1,
      "type": "service"
    }
    """

  Scenario: given invalid update request should return errors
    When I am admin
    When I do PUT /api/v4/entityservices/test-entityservice-not-found:
    """json
    {}
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "enabled": "Enabled is missing.",
        "impact_level": "ImpactLevel is missing.",
        "entity_pattern": "EntityPattern is missing.",
        "name": "Name is missing.",
        "output_template": "OutputTemplate is missing.",
        "sli_avail_state": "SliAvailState is missing."
      }
    }
    """

  Scenario: given invalid update request should return errors
    When I am admin
    When I do PUT /api/v4/entityservices/test-entityservice-to-update:
    """json
    {
      "category": "test-category-not-exist",
      "infos": [
        {}
      ],
      "sli_avail_state": 4
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "category": "Category doesn't exist.",
        "infos.0.name": "Name is missing.",
        "sli_avail_state": "SliAvailState should be 3 or less."
      }
    }
    """

  Scenario: given update request with already exists name should return error
    When I am admin
    When I do PUT /api/v4/entityservices/test-entityservice-to-update:
    """json
    {
      "name": "test-entityservice-to-check-unique-name-name"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
          "name": "Name already exists."
      }
    }
    """

  Scenario: given update request and no auth user should not allow access
    When I do PUT /api/v4/entityservices/test-entityservice-not-found
    Then the response code should be 401

  Scenario: given update request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/entityservices/test-entityservice-not-found
    Then the response code should be 403

  Scenario: given update request with not exist id should return not found error
    When I am admin
    When I do PUT /api/v4/entityservices/test-entityservice-not-found:
    """json
    {
      "name": "test-entityservice-to-update-not-found-name",
      "output_template": "test-entityservice-to-update-not-found-output",
      "category": "test-category-to-entityservice-edit",
      "impact_level": 2,
      "enabled": true,
      "sli_avail_state": 1,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-entityservice-to-update-not-found-pattern"
            }
          }
        ]
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

Feature: Create entity service
  I need to be able to create a entity service

  Scenario: given create request should return ok
    When I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-to-create-1-name",
      "output_template": "test-entityservice-to-create-1-output",
      "category": "test-category-to-entityservice-edit",
      "impact_level": 1,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-entityservice-to-create-1-pattern"
            }
          }
        ]
      ],
      "sli_avail_state": 1,
      "infos": [
        {
          "description": "test-entityservice-to-create-1-info-1-description",
          "name": "test-entityservice-to-create-1-info-1-name",
          "value": "test-entityservice-to-create-1-info-1-value"
        },
        {
          "description": "test-entityservice-to-create-1-info-2-description",
          "name": "test-entityservice-to-create-1-info-2-name",
          "value": false
        },
        {
          "description": "test-entityservice-to-create-1-info-3-description",
          "name": "test-entityservice-to-create-1-info-3-name",
          "value": 1022
        },
        {
          "description": "test-entityservice-to-create-1-info-4-description",
          "name": "test-entityservice-to-create-1-info-4-name",
          "value": 10.45
        },
        {
          "description": "test-entityservice-to-create-1-info-5-description",
          "name": "test-entityservice-to-create-1-info-5-name",
          "value": null
        },
        {
          "description": "test-entityservice-to-create-1-info-6-description",
          "name": "test-entityservice-to-create-1-info-6-name",
          "value": ["test-entityservice-to-create-1-info-6-value", false, 1022, 10.45, null]
        }
      ],
      "coordinates": {
        "lat": 62.34960927573042,
        "lng": 74.02834455685206
      }
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
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
              "value": "test-entityservice-to-create-1-pattern"
            }
          }
        ]
      ],
      "impact_level": 1,
      "infos": {
        "test-entityservice-to-create-1-info-1-name": {
          "description": "test-entityservice-to-create-1-info-1-description",
          "name": "test-entityservice-to-create-1-info-1-name",
          "value": "test-entityservice-to-create-1-info-1-value"
        },
        "test-entityservice-to-create-1-info-2-name": {
          "description": "test-entityservice-to-create-1-info-2-description",
          "name": "test-entityservice-to-create-1-info-2-name",
          "value": false
        },
        "test-entityservice-to-create-1-info-3-name": {
          "description": "test-entityservice-to-create-1-info-3-description",
          "name": "test-entityservice-to-create-1-info-3-name",
          "value": 1022
        },
        "test-entityservice-to-create-1-info-4-name": {
          "description": "test-entityservice-to-create-1-info-4-description",
          "name": "test-entityservice-to-create-1-info-4-name",
          "value": 10.45
        },
        "test-entityservice-to-create-1-info-5-name": {
          "description": "test-entityservice-to-create-1-info-5-description",
          "name": "test-entityservice-to-create-1-info-5-name",
          "value": null
        },
        "test-entityservice-to-create-1-info-6-name": {
          "description": "test-entityservice-to-create-1-info-6-description",
          "name": "test-entityservice-to-create-1-info-6-name",
          "value": ["test-entityservice-to-create-1-info-6-value", false, 1022, 10.45, null]
        }
      },
      "name": "test-entityservice-to-create-1-name",
      "output_template": "test-entityservice-to-create-1-output",
      "sli_avail_state": 1,
      "type": "service",
      "coordinates": {
        "lat": 62.34960927573042,
        "lng": 74.02834455685206
      }
    }
    """

  Scenario: given create request should return ok to get request
    When I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "name": "test-entityservice-to-create-2-name",
      "output_template": "test-entityservice-to-create-2-output",
      "category": "test-category-to-entityservice-edit",
      "impact_level": 2,
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-entityservice-to-create-2-pattern"
            }
          }
        ]
      ],
      "sli_avail_state": 1,
      "infos": [
        {
          "description": "test-entityservice-to-create-2-customer-description",
          "name": "test-entityservice-to-create-2-customer-name",
          "value": "test-entityservice-to-create-2-customer-value"
        }
      ]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/entityservices/{{ .lastResponse._id}}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
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
              "value": "test-entityservice-to-create-2-pattern"
            }
          }
        ]
      ],
      "impact_level": 2,
      "infos": {
        "test-entityservice-to-create-2-customer-name": {
          "description": "test-entityservice-to-create-2-customer-description",
          "name": "test-entityservice-to-create-2-customer-name",
          "value": "test-entityservice-to-create-2-customer-value"
        }
      },
      "name": "test-entityservice-to-create-2-name",
      "output_template": "test-entityservice-to-create-2-output",
      "sli_avail_state": 1,
      "type": "service"
    }
    """

  Scenario: given invalid create request should return errors
    When I am admin
    When I do POST /api/v4/entityservices:
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
    When I do POST /api/v4/entityservices:
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
    When I do POST /api/v4/entityservices:
    """json
    {
      "coordinates": {}
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "coordinates.lat": "Lat is missing.",
        "coordinates.lng": "Lng is missing."
      }
    }
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "coordinates": {
        "lat": 214983904,
        "lng": 214983904
      }
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "coordinates.lat": "Lat must contain valid latitude coordinates.",
        "coordinates.lng": "Lng must contain valid longitude coordinates."
      }
    }
    """

  Scenario: given create request with already exists name should return error
    When I am admin
    When I do POST /api/v4/entityservices:
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

  Scenario: given invalid create request with invalid entity pattern should return error
    When I am admin
    When I do POST /api/v4/entityservices:
    """json
    {
      "entity_pattern": [[]]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "entity_pattern": "EntityPattern is invalid entity pattern."
      }
    }
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-entityservice-to-create-3-pattern"
            }
          },
          {
            "field": "last_event_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 1,
                "unit": "m"
              }
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "entity_pattern": "EntityPattern is invalid entity pattern."
      }
    }
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-entityservice-to-create-3-pattern"
            }
          }
        ],
        [
          {
            "field": "impact",
            "cond": {
              "type": "has_one_of",
              "value": ["test-entityservice-to-create-3-pattern"]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "entity_pattern": "EntityPattern is invalid entity pattern."
      }
    }
    """
    When I do POST /api/v4/entityservices:
    """json
    {
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-entityservice-to-create-3-pattern"
            }
          }
        ],
        [
          {
            "field": "depends",
            "cond": {
              "type": "has_one_of",
              "value": ["test-entityservice-to-create-3-pattern"]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "entity_pattern": "EntityPattern is invalid entity pattern."
      }
    }
    """

  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/entityservices
    Then the response code should be 401

  Scenario: given create request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/entityservices
    Then the response code should be 403

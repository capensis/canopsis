Feature: Create entity service
  I need to be able to create a entity service

  Scenario: given create request should return ok
    When I am admin
    When I do POST /api/v4/entityservices:
    """
    {
      "name": "test-entityservice-to-create-1-name",
      "output_template": "test-entityservice-to-create-1-output",
      "category": "test-category-to-entityservice-edit",
      "impact_level": 1,
      "enabled": true,
      "entity_patterns": [{"name": "test-entityservice-to-create-1-pattern"}],
      "infos": [
        {
          "description": "test-entityservice-to-create-info-1-description",
          "name": "test-entityservice-to-create-info-1-name",
          "value": "test-entityservice-to-create-info-1-value"
        },
        {
          "description": "test-entityservice-to-create-info-2-description",
          "name": "test-entityservice-to-create-info-2-name",
          "value": false
        },
        {
          "description": "test-entityservice-to-create-info-3-description",
          "name": "test-entityservice-to-create-info-3-name",
          "value": 1022
        },
        {
          "description": "test-entityservice-to-create-info-4-description",
          "name": "test-entityservice-to-create-info-4-name",
          "value": 10.45
        },
        {
          "description": "test-entityservice-to-create-info-5-description",
          "name": "test-entityservice-to-create-info-5-name",
          "value": null
        },
        {
          "description": "test-entityservice-to-create-info-6-description",
          "name": "test-entityservice-to-create-info-6-name",
          "value": ["test-entityservice-to-create-info-6-value", false, 1022, 10.45, null]
        }
      ]
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """
    {
      "category": {
        "_id": "test-category-to-entityservice-edit",
        "name": "test-category-to-entityservice-edit-name",
        "author": "test-category-to-entityservice-edit-author",
        "created": 1592215337,
        "updated": 1592215337
      },
      "depends": [],
      "enabled": true,
      "enable_history": [],
      "enabled": true,
      "entity_patterns": [
        {
          "name": "test-entityservice-to-create-1-pattern"
        }
      ],
      "impact": [],
      "impact_level": 1,
      "infos": {
        "test-entityservice-to-create-info-1-name": {
          "description": "test-entityservice-to-create-info-1-description",
          "name": "test-entityservice-to-create-info-1-name",
          "value": "test-entityservice-to-create-info-1-value"
        },
        "test-entityservice-to-create-info-2-name": {
          "description": "test-entityservice-to-create-info-2-description",
          "name": "test-entityservice-to-create-info-2-name",
          "value": false
        },
        "test-entityservice-to-create-info-3-name": {
          "description": "test-entityservice-to-create-info-3-description",
          "name": "test-entityservice-to-create-info-3-name",
          "value": 1022
        },
        "test-entityservice-to-create-info-4-name": {
          "description": "test-entityservice-to-create-info-4-description",
          "name": "test-entityservice-to-create-info-4-name",
          "value": 10.45
        },
        "test-entityservice-to-create-info-5-name": {
          "description": "test-entityservice-to-create-info-5-description",
          "name": "test-entityservice-to-create-info-5-name",
          "value": null
        },
        "test-entityservice-to-create-info-6-name": {
          "description": "test-entityservice-to-create-info-6-description",
          "name": "test-entityservice-to-create-info-6-name",
          "value": ["test-entityservice-to-create-info-6-value", false, 1022, 10.45, null]
        }
      },
      "measurements": null,
      "name": "test-entityservice-to-create-1-name",
      "output_template": "test-entityservice-to-create-1-output",
      "type": "service"
    }
    """

  Scenario: given create request should return ok to get request
    When I am admin
    When I do POST /api/v4/entityservices:
    """
    {
      "name": "test-entityservice-to-create-2-name",
      "output_template": "test-entityservice-to-create-2-output",
      "category": "test-category-to-entityservice-edit",
      "impact_level": 2,
      "enabled": true,
      "entity_patterns": [{"name": "test-entityservice-to-create-2-pattern"}],
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
    """
    {
      "category": {
        "_id": "test-category-to-entityservice-edit",
        "name": "test-category-to-entityservice-edit-name",
        "author": "test-category-to-entityservice-edit-author",
        "created": 1592215337,
        "updated": 1592215337
      },
      "depends": [],
      "enabled": true,
      "enable_history": [],
      "enabled": true,
      "entity_patterns": [
        {
          "name": "test-entityservice-to-create-2-pattern"
        }
      ],
      "impact": [],
      "impact_level": 2,
      "infos": {
        "test-entityservice-to-create-2-customer-name": {
          "description": "test-entityservice-to-create-2-customer-description",
          "name": "test-entityservice-to-create-2-customer-name",
          "value": "test-entityservice-to-create-2-customer-value"
        }
      },
      "measurements": null,
      "name": "test-entityservice-to-create-2-name",
      "output_template": "test-entityservice-to-create-2-output",
      "type": "service"
    }
    """

  Scenario: given invalid create request should return errors
    When I am admin
    When I do POST /api/v4/entityservices:
    """
    {}
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "enabled": "Enabled is missing.",
        "impact_level": "ImpactLevel is missing.",
        "name": "Name is missing.",
        "output_template": "OutputTemplate is missing."
      }
    }
    """

  Scenario: given invalid create request should return errors
    When I am admin
    When I do POST /api/v4/entityservices:
    """
    {
      "name": "test-entityservice-to-create-3-name",
      "output_template": "test-entityservice-to-create-3-output",
      "category": "test-category-not-exist",
      "impact_level": 3,
      "enabled": true,
      "entity_patterns": [{"name": "test-entityservice-to-create-3-pattern"}],
      "infos": [
        {}
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "category": "Category doesn't exist.",
        "infos.0.name": "Name is missing."
      }
    }
    """

  Scenario: given create request with already exists name should return error
    When I am admin
    When I do POST /api/v4/entityservices:
    """
    {
      "name": "test-entityservice-to-check-unique-name-name",
      "output_template": "test-entityservice-to-create-4-output",
      "category": "test-category-to-entityservice-edit",
      "impact_level": 4,
      "enabled": true,
      "entity_patterns": [{"name": "test-entityservice-to-create-4-pattern"}]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
          "name": "Name already exists."
      }
    }
    """

  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/entityservices
    Then the response code should be 401

  Scenario: given create request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/entityservices
    Then the response code should be 403

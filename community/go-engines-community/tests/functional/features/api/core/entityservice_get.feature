Feature: Get entity service
  I need to be able to get a entity service

  Scenario: given get request should return entity
    When I am admin
    When I do GET /api/v4/entityservices/test-entityservice-to-get
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "_id": "test-entityservice-to-get",
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
          "name": "test-entityservice-to-get-pattern"
        }
      ],
      "impact": [],
      "impact_level": 3,
      "infos": {
        "test-entityservice-to-get-info-1": {
          "name": "test-entityservice-to-get-info-1-name",
          "description": "test-entityservice-to-get-info-1-description",
          "value": "test-entityservice-to-get-info-1-value"
        },
        "test-entityservice-to-get-info-2": {
          "name": "test-entityservice-to-get-info-2-name",
          "description": "test-entityservice-to-get-info-2-description",
          "value": false
        },
        "test-entityservice-to-get-info-3": {
          "name": "test-entityservice-to-get-info-3-name",
          "description": "test-entityservice-to-get-info-3-description",
          "value": 1022
        },
        "test-entityservice-to-get-info-4": {
          "name": "test-entityservice-to-get-info-4-name",
          "description": "test-entityservice-to-get-info-4-description",
          "value": 10.45
        },
        "test-entityservice-to-get-info-5": {
          "name": "test-entityservice-to-get-info-5-name",
          "description": "test-entityservice-to-get-info-5-description",
          "value": null
        },
        "test-entityservice-to-get-info-6": {
          "name": "test-entityservice-to-get-info-6-name",
          "description": "test-entityservice-to-get-info-6-description",
          "value": ["test-entityservice-to-get-info-6-value", false, 1022, 10.45, null]
        },
        "test-entityservice-to-get-info-7": {
          "name": "test-entityservice-to-get-info-7",
          "description": "test-entityservice-to-get-info-7-description",
          "value": "test-entityservice-to-get-info-7-value"
        }
      },
      "measurements": null,
      "name": "test-entityservice-to-get-name",
      "output_template": "test-entityservice-to-get-output",
      "type": "service"
    }
    """

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/entityservices/test-entityservice-not-found
    Then the response code should be 401

  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/entityservices/test-entityservice-not-found
    Then the response code should be 403

  Scenario: given get request with not exist id should return not found error
    When I am admin
    When I do GET /api/v4/entityservices/test-entityservice-not-found
    Then the response code should be 404
    Then the response body should be:
    """
    {
      "error": "Not found"
    }
    """

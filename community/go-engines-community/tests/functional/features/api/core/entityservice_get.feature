Feature: Get entity service
  I need to be able to get a entity service

  Scenario: given get request should return entity
    When I am admin
    When I do GET /api/v4/entityservices/test-entityservice-to-get-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-entityservice-to-get-1",
      "category": {
        "_id": "test-category-to-entityservice-edit",
        "name": "test-category-to-entityservice-edit-name",
        "author": "test-category-to-entityservice-edit-author",
        "created": 1592215337,
        "updated": 1592215337
      },
      "depends": [],
      "enabled": true,
      "enable_history": null,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-entityservice-to-get-1-pattern"
            }
          }
        ]
      ],
      "old_entity_patterns": null,
      "impact": [],
      "impact_level": 3,
      "infos": {
        "test-entityservice-to-get-1-info-1": {
          "name": "test-entityservice-to-get-1-info-1-name",
          "description": "test-entityservice-to-get-1-info-1-description",
          "value": "test-entityservice-to-get-1-info-1-value"
        },
        "test-entityservice-to-get-1-info-2": {
          "name": "test-entityservice-to-get-1-info-2-name",
          "description": "test-entityservice-to-get-1-info-2-description",
          "value": false
        },
        "test-entityservice-to-get-1-info-3": {
          "name": "test-entityservice-to-get-1-info-3-name",
          "description": "test-entityservice-to-get-1-info-3-description",
          "value": 1022
        },
        "test-entityservice-to-get-1-info-4": {
          "name": "test-entityservice-to-get-1-info-4-name",
          "description": "test-entityservice-to-get-1-info-4-description",
          "value": 10.45
        },
        "test-entityservice-to-get-1-info-5": {
          "name": "test-entityservice-to-get-1-info-5-name",
          "description": "test-entityservice-to-get-1-info-5-description",
          "value": null
        },
        "test-entityservice-to-get-1-info-6": {
          "name": "test-entityservice-to-get-1-info-6-name",
          "description": "test-entityservice-to-get-1-info-6-description",
          "value": ["test-entityservice-to-get-1-info-6-value", false, 1022, 10.45, null]
        },
        "test-entityservice-to-get-1-info-7": {
          "name": "test-entityservice-to-get-1-info-7",
          "description": "test-entityservice-to-get-1-info-7-description",
          "value": "test-entityservice-to-get-1-info-7-value"
        }
      },
      "measurements": null,
      "name": "test-entityservice-to-get-1-name",
      "output_template": "test-entityservice-to-get-1-output",
      "sli_avail_state": 0,
      "type": "service"
    }
    """

  Scenario: given get request with old pattern should return entity
    When I am admin
    When I do GET /api/v4/entityservices/test-entityservice-to-get-2
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-entityservice-to-get-2",
      "category": null,
      "depends": [],
      "enabled": true,
      "enable_history": null,
      "old_entity_patterns": [{"name": "test-entityservice-to-get-2-pattern"}],
      "impact": [],
      "impact_level": 1,
      "infos": {},
      "measurements": null,
      "name": "test-entityservice-to-get-2-name",
      "output_template": "test-entityservice-to-get-2-output",
      "sli_avail_state": 0,
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
    """json
    {
      "error": "Not found"
    }
    """

Feature: Get entity basic
  I need to be able to get a entity basic

  Scenario: given get request should return entity
    When I am admin
    When I do GET /api/v4/entitybasics?_id=test-entitybasic-to-get-resource/test-entitybasic-to-get-component
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-entitybasic-to-get-resource/test-entitybasic-to-get-component",
      "category": {
        "_id": "test-category-to-entitybasic-edit",
        "name": "test-category-to-entitybasic-edit-name"
      },
      "component": "test-entitybasic-to-get-component",
      "connector": "test-entitybasic-to-get-connector/test-entitybasic-to-get-connector-name",
      "description": "test-entitybasic-to-get-resource-description",
      "enabled": true,
      "old_entity_patterns": null,
      "impact_level": 3,
      "infos": {
        "test-entitybasic-to-get-info-1": {
          "name": "test-entitybasic-to-get-info-1-name",
          "description": "test-entitybasic-to-get-info-1-description",
          "value": "test-entitybasic-to-get-info-1-value"
        },
        "test-entitybasic-to-get-info-2": {
          "name": "test-entitybasic-to-get-info-2-name",
          "description": "test-entitybasic-to-get-info-2-description",
          "value": false
        },
        "test-entitybasic-to-get-info-3": {
          "name": "test-entitybasic-to-get-info-3-name",
          "description": "test-entitybasic-to-get-info-3-description",
          "value": 1022
        },
        "test-entitybasic-to-get-info-4": {
          "name": "test-entitybasic-to-get-info-4-name",
          "description": "test-entitybasic-to-get-info-4-description",
          "value": 10.45
        },
        "test-entitybasic-to-get-info-5": {
          "name": "test-entitybasic-to-get-info-5-name",
          "description": "test-entitybasic-to-get-info-5-description",
          "value": null
        },
        "test-entitybasic-to-get-info-6": {
          "name": "test-entitybasic-to-get-info-6-name",
          "description": "test-entitybasic-to-get-info-6-description",
          "value": ["test-entitybasic-to-get-info-6-value", false, 1022, 10.45, null]
        },
        "test-entitybasic-to-get-info-7": {
          "name": "test-entitybasic-to-get-info-7",
          "description": "test-entitybasic-to-get-info-7-description",
          "value": "test-entitybasic-to-get-info-7-value"
        }
      },
      "name": "test-entitybasic-to-get-resource",
      "sli_avail_state": 0,
      "type": "resource",
      "coordinates": {
        "lat": 64.52269494598361,
        "lng": 54.037685420804365
      }
    }
    """

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/entitybasics?_id=test-entitybasic-not-found
    Then the response code should be 401

  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/entitybasics?_id=test-entitybasic-not-found
    Then the response code should be 403

  Scenario: given get request with not exist id should return not found error
    When I am admin
    When I do GET /api/v4/entitybasics?_id=test-entitybasic-not-found
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """

Feature: Update entity basic
  I need to be able to update a entity basic

  Scenario: given update request should update entity
    When I am admin
    When I do PUT /api/v4/entitybasics?_id=test-entitybasic-to-update-connector/test-entitybasic-to-update-connector-name:
    """json
    {
      "description": "test-entitybasic-to-update-connector-description-updated",
      "enabled": true,
      "category": "test-category-to-entitybasic-edit",
      "impact_level": 3,
      "sli_avail_state": 1,
      "infos": [
        {
          "description": "test-entitybasic-to-update-info-1-description",
          "name": "test-entitybasic-to-update-info-1-name",
          "value": "test-entitybasic-to-update-info-1-value"
        },
        {
          "description": "test-entitybasic-to-update-info-2-description",
          "name": "test-entitybasic-to-update-info-2-name",
          "value": false
        },
        {
          "description": "test-entitybasic-to-update-info-3-description",
          "name": "test-entitybasic-to-update-info-3-name",
          "value": 1022
        },
        {
          "description": "test-entitybasic-to-update-info-4-description",
          "name": "test-entitybasic-to-update-info-4-name",
          "value": 10.45
        },
        {
          "description": "test-entitybasic-to-update-info-5-description",
          "name": "test-entitybasic-to-update-info-5-name",
          "value": null
        },
        {
          "description": "test-entitybasic-to-update-info-6-description",
          "name": "test-entitybasic-to-update-info-6-name",
          "value": ["test-entitybasic-to-update-info-6-value", false, 1022, 10.45, null]
        }
      ],
      "coordinates": {
        "lat": 62.34960927573042,
        "lng": 74.02834455685206
      }
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-entitybasic-to-update-connector/test-entitybasic-to-update-connector-name",
      "category": {
        "_id": "test-category-to-entitybasic-edit",
        "name": "test-category-to-entitybasic-edit-name"
      },
      "description": "test-entitybasic-to-update-connector-description-updated",
      "enabled": true,
      "impact_level": 3,
      "infos": {
        "test-entitybasic-to-update-info-1-name": {
          "description": "test-entitybasic-to-update-info-1-description",
          "name": "test-entitybasic-to-update-info-1-name",
          "value": "test-entitybasic-to-update-info-1-value"
        },
        "test-entitybasic-to-update-info-2-name": {
          "description": "test-entitybasic-to-update-info-2-description",
          "name": "test-entitybasic-to-update-info-2-name",
          "value": false
        },
        "test-entitybasic-to-update-info-3-name": {
          "description": "test-entitybasic-to-update-info-3-description",
          "name": "test-entitybasic-to-update-info-3-name",
          "value": 1022
        },
        "test-entitybasic-to-update-info-4-name": {
          "description": "test-entitybasic-to-update-info-4-description",
          "name": "test-entitybasic-to-update-info-4-name",
          "value": 10.45
        },
        "test-entitybasic-to-update-info-5-name": {
          "description": "test-entitybasic-to-update-info-5-description",
          "name": "test-entitybasic-to-update-info-5-name",
          "value": null
        },
        "test-entitybasic-to-update-info-6-name": {
          "description": "test-entitybasic-to-update-info-6-description",
          "name": "test-entitybasic-to-update-info-6-name",
          "value": ["test-entitybasic-to-update-info-6-value", false, 1022, 10.45, null]
        }
      },
      "name": "test-entitybasic-to-update-connector-name",
      "sli_avail_state": 1,
      "type": "connector",
      "coordinates": {
        "lat": 62.34960927573042,
        "lng": 74.02834455685206
      }
    }
    """
    When I do GET /api/v4/entities/context-graph?_id=test-entitybasic-to-update-connector/test-entitybasic-to-update-connector-name
    Then the response code should be 200
    Then the response array key "depends" should contain:
    """json
    [
      "test-entitybasic-to-update-component-1",
      "test-entitybasic-to-update-component-2"
    ]
    """
    Then the response array key "impact" should contain:
    """json
    [
      "test-entitybasic-to-update-resource-1/test-entitybasic-to-update-component-1",
      "test-entitybasic-to-update-resource-2/test-entitybasic-to-update-component-2"
    ]
    """

  Scenario: given invalid update request should return errors
    When I am admin
    When I do PUT /api/v4/entitybasics?_id=test-entitybasic-not-found:
    """json
    { }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "enabled": "Enabled is missing.",
        "impact_level": "ImpactLevel is missing.",
        "sli_avail_state": "SliAvailState is missing."
      }
    }
    """

  Scenario: given invalid update request should return errors
    When I am admin
    When I do PUT /api/v4/entitybasics?_id=test-entitybasic-to-update:
    """json
    {
      "description": "test-entitybasic-to-update-description",
      "enabled": true,
      "category": "test-category-not-exist",
      "impact_level": 11,
      "sli_avail_state": 4,
      "infos": [
        {}
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "impact_level": "ImpactLevel should be 10 or less.",
        "sli_avail_state": "SliAvailState should be 3 or less.",
        "category": "Category doesn't exist.",
        "infos.0.name": "Name is missing."
      }
    }
    """

  Scenario: given invalid update request should return errors
    When I am admin
    When I do PUT /api/v4/entitybasics?_id=test-entitybasic-to-update:
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
    When I do PUT /api/v4/entitybasics?_id=test-entitybasic-to-update:
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

  Scenario: given update request and no auth user should not allow access
    When I do PUT /api/v4/entitybasics?_id=test-entitybasic-not-found
    Then the response code should be 401

  Scenario: given update request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/entitybasics?_id=test-entitybasic-not-found
    Then the response code should be 403

  Scenario: given update request with not exist id should return not found error
    When I am admin
    When I do PUT /api/v4/entitybasics?_id=test-entitybasic-not-found:
    """json
    {
      "description": "test-entitybasic-not-found-description",
      "enabled": true,
      "category": "test-category-to-entitybasic-edit",
      "impact_level": 1,
      "sli_avail_state": 0,
      "infos": []
    }
    """
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """

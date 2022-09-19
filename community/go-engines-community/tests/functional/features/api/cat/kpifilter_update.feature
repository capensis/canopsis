Feature: Update a kpi filter
  I need to be able to update a kpi filter
  Only admin should be able to update a kpi filter

  Scenario: given update request should update filter
    When I am admin
    Then I do PUT /api/v4/cat/kpi-filters/test-kpi-filter-to-update-1:
    """json
    {
      "name": "test-kpi-filter-to-update-1-name-updated",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-kpi-filter-to-update-1-resource-updated"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "created": 1619083733,
      "name": "test-kpi-filter-to-update-1-name-updated",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-kpi-filter-to-update-1-resource-updated"
            }
          }
        ]
      ]
    }
    """

  Scenario: given get request and no auth user should not allow access
    When I do PUT /api/v4/cat/kpi-filters/test-kpi-filter-to-update
    Then the response code should be 401

  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/cat/kpi-filters/test-kpi-filter-to-update
    Then the response code should be 403

  Scenario: given update request with not exist id should return not found error
    When I am admin
    When I do PUT /api/v4/cat/kpi-filters/test-kpi-filter-not-found:
    """json
    {
      "name": "test-kpi-filter-to-update-not-found",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-kpi-filter-to-update-not-found"
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

  Scenario: given update request with missing fields should return bad request
    When I am admin
    Then I do PUT /api/v4/cat/kpi-filters/test-kpi-filter-to-update:
    """
    {}
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "entity_pattern": "EntityPattern is missing.",
        "name": "Name is missing."
      }
    }
    """

  Scenario: given update request with already exists name should return error
    When I am admin
    When I do PUT /api/v4/cat/kpi-filters/test-kpi-filter-to-update-1:
    """json
    {
      "name": "test-kpi-filter-to-check-unique-name"
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

  Scenario: given update request for filter with old pattern should not update old pattern
    When I am admin
    Then I do PUT /api/v4/cat/kpi-filters/test-kpi-filter-to-update-2:
    """json
    {
      "name": "test-kpi-filter-to-update-2-name-updated"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "created": 1619083733,
      "name": "test-kpi-filter-to-update-2-name-updated",
      "old_entity_patterns": [{"name": "test-kpi-filter-to-update-2-resource"}]
    }
    """

Feature: get check state settings for an entity
  I need to be able to check state settings for an entity
  Only admin should be able to check state settings for an entity

  @concurrent
  Scenario: given check request and no auth user should not allow access
    When I do POST /api/v4/entities/check-state-setting
    Then the response code should be 401

  @concurrent
  Scenario: given check request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/entities/check-state-setting
    Then the response code should be 403

  @concurrent
  Scenario: given check request should be ok
    When I am admin
    When I do POST /api/v4/entities/check-state-setting:
    """json
    {
      "_id": "name",
      "name": "test-service-state-settings-to-check-1",
      "type": "service",
      "impact_level": 1
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "title": "state-settings-to-check-1-title"
    }
    """

  @concurrent
  Scenario: given check request when entity matches 2 settings by name should get the first by priority
    When I am admin
    When I do POST /api/v4/entities/check-state-setting:
    """json
    {
      "_id": "name",
      "name": "test-service-state-settings-to-check-2",
      "type": "service",
      "impact_level": 1
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "title": "state-settings-to-check-2-1-title"
    }
    """

  @concurrent
  Scenario: given check request without entity matches should return empty string
    When I am admin
    When I do POST /api/v4/entities/check-state-setting:
    """json
    {
      "_id": "name",
      "name": "test-service-state-settings-to-check-not-matched",
      "type": "service",
      "impact_level": 1
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "title": ""
    }
    """

  @concurrent
  Scenario: given check request should return different state settings depending on infos
    When I am admin
    When I do POST /api/v4/entities/check-state-setting:
    """json
    {
      "_id": "name",
      "name": "test-service-state-settings-to-check-3",
      "type": "service",
      "impact_level": 1,
      "infos": [
        {
          "description": "test-description",
          "name": "infos_1",
          "value": "value_1"
        }
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "title": "state-settings-to-check-3-1-title"
    }
    """
    When I do POST /api/v4/entities/check-state-setting:
    """json
    {
      "_id": "name",
      "name": "test-service-state-settings-to-check-3",
      "type": "service",
      "impact_level": 1,
      "infos": [
        {
          "description": "test-description",
          "name": "infos_1",
          "value": "value_2"
        }
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "title": "state-settings-to-check-3-2-title"
    }
    """
    When I do POST /api/v4/entities/check-state-setting:
    """json
    {
      "_id": "name",
      "name": "test-service-state-settings-to-check-3",
      "type": "service",
      "impact_level": 1,
      "infos": [
        {
          "description": "test-description",
          "name": "infos_1",
          "value": "value_3"
        }
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "title": ""
    }
    """

  @concurrent
  Scenario: given invalid check request should return error
    When I am admin
    When I do POST /api/v4/entities/check-state-setting:
    """json
    {}
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "impact_level": "ImpactLevel is missing.",
        "name": "Name is missing.",
        "type": "Type is missing."
      }
    }
    """

  @concurrent
  Scenario: given invalid check request with wrong type should return error
    When I am admin
    When I do POST /api/v4/entities/check-state-setting:
    """json
    {
      "type": "resource"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "type": "Type must be one of [component service]."
      }
    }
    """

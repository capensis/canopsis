Feature: Create a pbehavior type
  I need to be able to create a pbehavior type
  Only admin should be able to create a pbehavior type

  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/pbehavior-types
    Then the response code should be 401

  Scenario: given create request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/pbehavior-types
    Then the response code should be 403

  Scenario: given invalid create request should return errors
    When I am admin
    When I do POST /api/v4/pbehavior-types:
    """json
    {}
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "description": "Description is missing.",
        "icon_name": "IconName is missing.",
        "color": "Color is missing.",
        "name": "Name is missing.",
        "priority": "Priority is missing.",
        "type": "Type is missing."
      }
    }
    """
    When I do POST /api/v4/pbehavior-types:
    """json
    {
      "type": "notexist"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "type": "Type must be one of [active inactive maintenance pause]."
      }
    }
    """

  Scenario: given create request should return ok
    When I am admin
    When I do POST /api/v4/pbehavior-types:
    """json
    {
      "name": "Active State",
      "description": "Active state type",
      "type": "active",
      "priority": 177,
      "icon_name": "exclamation-mark.png",
      "color": "#FFFFFF"
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "name": "Active State",
      "description": "Active state type",
      "type": "active",
      "priority": 177,
      "icon_name": "exclamation-mark.png",
      "color": "#FFFFFF"
    }
    """

  Scenario: given create request with custom id should return ok
    When I am admin
    When I do POST /api/v4/pbehavior-types:
    """json
    {
      "_id": "custom-id",
      "name": "Active State custom-id",
      "description": "Active state type",
      "type": "active",
      "priority": 277,
      "icon_name": "exclamation-mark.png",
      "color": "#FFFFFF"
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/pbehavior-types/custom-id
    Then the response code should be 200

  Scenario: given create request with already exists custom id should return error
    When I am admin
    When I do POST /api/v4/pbehavior-types:
    """json
    {
      "_id": "test-type-to-update"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "_id": "ID already exists."
      }
    }
    """

  Scenario: given create request with already exists priority should return error
    When I am admin
    When I do POST /api/v4/pbehavior-types:
    """json
    {
      "priority": 10
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "priority": "Priority already exists."
      }
    }
    """

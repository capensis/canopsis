Feature: Update a pbehavior type
  I need to be able to update a pbehavior type
  Only admin should be able to update a pbehavior type

  Scenario: given update request and no auth user should not allow access
    When I do PUT /api/v4/pbehavior-types/test-to-update
    Then the response code should be 401

  Scenario: given update request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/pbehavior-types/test-to-update
    Then the response code should be 403

  Scenario: given update request with not exist id should return not found error
    When I am admin
    When I do PUT /api/v4/pbehavior-types/test-to-update-unavailable:
    """json
    {
      "name": "Maintenance State",
      "description": "Maintenance state type",
      "type": "maintenance",
      "priority": 399,
      "icon_name": "exclamation-mark.png"
    }
    """
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """

  Scenario: given invalid update request should return errors
    When I am admin
    When I do PUT /api/v4/pbehavior-types/test-to-update:
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
        "name": "Name is missing.",
        "priority": "Priority is missing.",
        "type": "Type is missing."
      }
    }
    """
    When I do PUT /api/v4/pbehavior-types/test-to-update:
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
        "type": "type must be one of [active inactive maintenance pause]."
      }
    }
    """

  Scenario: given update request should update type
    When I am admin
    When I do PUT /api/v4/pbehavior-types/test-type-to-update:
    """json
    {
      "name": "Maintenance State",
      "description": "Maintenance state type",
      "type": "maintenance",
      "priority": 399,
      "icon_name": "exclamation-mark.png"
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-type-to-update",
      "name": "Maintenance State",
      "description": "Maintenance state type",
      "type": "maintenance",
      "priority": 399,
      "icon_name": "exclamation-mark.png"
    }
    """

  Scenario: given update request with already exists priority and name should return error
    When I am admin
    When I do PUT /api/v4/pbehavior-types/test-type-to-update:
    """json
    {
      "name": "Some State",
      "priority": 4
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "name": "Name already exists.",
        "priority": "Priority already exists."
      }
    }
    """

  Scenario: given update request for default type should return error
    When I am admin
    When I do PUT /api/v4/pbehavior-types/test-default-pause-type:
    """json
    {
      "name": "Default Type Pause",
      "description": "Maintenance state type",
      "type": "maintenance",
      "priority": 3,
      "icon_name": "exclamation-mark.png"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "error": "type is default"
    }
    """

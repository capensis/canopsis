Feature: PBehavior Type update

  Scenario: PUT as unauthorized
    When I do PUT /api/v4/pbehavior-types/test-to-update:
            """
            {
            Name: "Maintenance State",
            Description: "Maintenance state type",
            Type: "maintenance",
            Priority: 11,
            icon_name: "exclamation-mark.png"
            }
            """
    Then the response code should be 401

  Scenario: PUT without permissions
    When I am noperms
    When I do PUT /api/v4/pbehavior-types/test-to-update:
            """
            {
            Name: "Maintenance State",
            Description: "Maintenance state type",
            Type: "maintenance",
            Priority: 11,
            icon_name: "exclamation-mark.png"
            }
            """
    Then the response code should be 403

  Scenario: PUT a valid PBehavior type that doesn't exist
    When I am admin
    When I do PUT /api/v4/pbehavior-types/test-to-update-unavailable:
            """
            {
                "Name": "Maintenance State",
                "Description": "Maintenance state type",
                "Type": "maintenance",
                "Priority": 399,
                "icon_name": "exclamation-mark.png"
            }
            """
    Then the response code should be 404

  Scenario: PUT a valid PBehavior Type with invalid priority
    When I am admin
    When I do PUT /api/v4/pbehavior-types/test-to-update:
            """
            {
                "name": "Active State",
                "description": "Active state type",
                "type": "active",
                "priority": "invalid",
                "icon_name": "exclamation-mark.png"
            }
            """
    Then the response code should be 400

  Scenario: PUT a valid PBehavior Type
    When I am admin
    When I do PUT /api/v4/pbehavior-types/test-type-to-update:
    """
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
    """
    {
        "_id": "test-type-to-update",
        "name": "Maintenance State",
        "description": "Maintenance state type",
        "type": "maintenance",
        "priority": 399,
        "icon_name": "exclamation-mark.png"
    }
    """

  Scenario: PUT a valid PBehavior Type without any changes
    When I am admin
    When I do PUT /api/v4/pbehavior-types/test-type-to-update:
    """
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
    """
    {
        "_id": "test-type-to-update",
        "name": "Maintenance State",
        "description": "Maintenance state type",
        "type": "maintenance",
        "priority": 399,
        "icon_name": "exclamation-mark.png"
    }
    """

  Scenario: PUT a valid PBehavior Type with already existed priority and name
    When I am admin
    When I do PUT /api/v4/pbehavior-types/test-type-to-update:
    """
    {
        "name": "Some State",
        "description": "Maintenance state type",
        "type": "maintenance",
        "priority": 4,
        "icon_name": "exclamation-mark.png"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "name": "Name already exists.",
        "priority": "Priority already exists."
      }
    }
    """

  Scenario: Given default type Should return error
    When I am admin
    When I do PUT /api/v4/pbehavior-types/test-default-pause-type
    """
    {
        "name": "Some State",
        "description": "Maintenance state type",
        "type": "maintenance",
        "priority": 4,
        "icon_name": "exclamation-mark.png"
    }
    """
    Then the response code should be 400

Feature: PBehavior Type create
  Create PBehavior Type item

  Scenario: POST as unautohorized
    When I do POST /api/v4/pbehavior-types:
    """
    {
      "name": "Active State",
      "description": "Active state type",
      "type": "active",
      "priority": 12,
      "icon_name": "exclamation-mark.png"
    }
    """
    Then the response code should be 401

  Scenario: POST without permissions
    When I am noperms
    When I do POST /api/v4/pbehavior-types:
    """
    {
      "name": "Active State",
      "description": "some_description"
    }
    """
    Then the response code should be 403

  Scenario: POST a valid PBehavior Type with invalid payload
    When I am admin
    When I do POST /api/v4/pbehavior-types:
    """
    {
      "name": "Active State",
      "description": "Active state type",
      "type": "active",
      "priority": 15,
      "icon_name": "exclamation-mark.png"
    }
    """
    Then the response code should be 400

  Scenario: POST a valid PBehavior Type instance
    When I am admin
    When I do POST /api/v4/pbehavior-types:
    """
    {
      "name": "Active State",
      "description": "Active state type",
      "type": "active",
      "priority": 177,
      "icon_name": "exclamation-mark.png"
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """
    {
      "name": "Active State",
      "description": "Active state type",
      "type": "active",
      "priority": 177,
      "icon_name": "exclamation-mark.png"
    }
    """

  Scenario: POST a valid PBehavior Type with custom id
    When I am admin
    When I do POST /api/v4/pbehavior-types:
    """
    {
      "_id": "custom-id",
      "name": "Active State custom-id",
      "description": "Active state type",
      "type": "active",
      "priority": 277,
      "icon_name": "exclamation-mark.png"
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/pbehavior-types/custom-id
    Then the response code should be 200

  Scenario: POST a valid PBehavior Type with custom id that already exist should cause dup error
    When I am admin
    When I do POST /api/v4/pbehavior-types:
    """
    {
      "_id": "test-type-to-update",
      "name": "Active State custom-id 2",
      "description": "Active state type",
      "type": "active",
      "priority": 278,
      "icon_name": "exclamation-mark.png"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "_id": "ID already exists"
      }
    }
    """

  Scenario: POST a duplicate PBehavior Type priority
    When I am admin
    When I do POST /api/v4/pbehavior-types:
    """
    {
      "name": "Duplicate 10-Active State",
      "description": "Active state type",
      "type": "active",
      "priority": 10,
      "icon_name": "exclamation-mark.png"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "priority": "Priority already exists"
      }
    }
    """

  Scenario: POST a valid PBehavior Type instance
    When I am admin
    When I do POST /api/v4/pbehavior-types:
    """
    {
      "name": "Active State 2",
      "description": "Active state type",
      "type": "active",
      "priority": 188,
      "icon_name": "exclamation-mark.png"
    }
    """
    When I do GET /api/v4/pbehavior-types/{{ .lastResponse._id}}
    Then the response code should be 200

  Scenario: POST a valid PBehavior Type with color
    When I am admin
    When I do POST /api/v4/pbehavior-types:
    """
    {
      "_id": "test-type-active-green",
      "name": "Active State 2 Green",
      "description": "Active state type",
      "type": "active",
      "priority": 191,
      "icon_name": "exclamation-mark.png",
      "color": "green"
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/pbehavior-types/test-type-active-green
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "_id": "test-type-active-green",
      "name": "Active State 2 Green",
      "description": "Active state type",
      "type": "active",
      "priority": 191,
      "icon_name": "exclamation-mark.png",
      "color": "green"
    }
    """

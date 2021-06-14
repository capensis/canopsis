Feature: get a PBehavior Type

  Scenario: GET unauthorized
    When I do GET /api/v4/pbehavior-types/test-type-to-get
    Then the response code should be 401

  Scenario: GET without permissions
    When I am noperms
    When I do GET /api/v4/pbehavior-types/test-type-to-get
    Then the response code should be 403

  Scenario: Get one with success
    When I am admin
    When I do GET /api/v4/pbehavior-types/test-type-to-get
    Then the response code should be 200
    Then the response body should be:
    """
    {
        "_id": "test-type-to-get",
        "description": "Some State type",
        "icon_name": "test-to-get-icon",
        "name": "Some State",
        "priority": 8,
        "type": "active"
    }
    """

  Scenario: Get a PBehavior Type with not found response
    When I am admin
    When I do GET /api/v4/pbehavior-types/test-type-to-get-not-found
    Then the response code should be 404
    Then the response body should be:
    """
    {
        "error": "Not found"
    }
    """

  Scenario: Get "search=Find State" with success
    When I am admin
    When I do GET /api/v4/pbehavior-types?search=Find%20State
    Then the response code should be 200
    Then the response body should be:
    """
    {
        "meta": {
            "page": 1,
            "page_count": 1,
            "per_page": 10,
            "total_count": 1
        },
        "data": [
            {
                "_id": "test-type-to-find",
                "description": "Find State type",
                "icon_name": "test-to-find-icon",
                "name": "Find State",
                "priority": 9,
                "type": "active"
            }
        ]
    }
    """

  Scenario: Find default types
    When I am admin
    When I do GET /api/v4/pbehavior-types?default=true&with_flags=true
    Then the response code should be 200
    Then the response body should be:
    """
    {
        "meta": {
            "page": 1,
            "page_count": 1,
            "per_page": 10,
            "total_count": 4
        },
        "data": [
            {
              "_id": "test-default-active-type",
              "description": "Default Type Active",
              "icon_name": "test-active-icon",
              "name": "Default Type Active",
              "priority": 2,
              "type": "active",
              "editable": false,
              "deletable": false
            },
            {
              "_id": "test-default-inactive-type",
              "description": "Default Type Inactive",
              "icon_name": "test-inactive-icon",
              "name": "Default Type Inactive",
              "priority": 1,
              "type": "inactive",
              "editable": false,
              "deletable": false
            },
            {
              "_id": "test-default-maintenance-type",
              "description": "Default Type Maintenance",
              "icon_name": "test-maintenance-icon",
              "name": "Default Type Maintenance",
              "priority": 4,
              "type": "maintenance",
              "editable": false,
              "deletable": false
            },
            {
              "_id": "test-default-pause-type",
              "description": "Default Type Pause",
              "icon_name": "test-pause-icon",
              "name": "Default Type Pause",
              "priority": 3,
              "type": "pause",
              "editable": false,
              "deletable": false
            }
        ]
    }
    """

Feature: Get pbehavior exception
  I need to be able to read a pbehavior exception

  Scenario: given get all request should return exceptions
    When I am admin
    When I do GET /api/v4/pbehavior-exceptions?search=test-exception-to-get
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "test-exception-to-get-1",
          "name": "test-exception-to-get-1-name",
          "description": "test-exception-to-get-1-description",
          "exdates": [
            {
              "begin": 1591164001,
              "end": 1591167601,
              "type": {
                "_id": "test-type-to-exception-edit-2",
                "description": "Exception edit 2 State type",
                "icon_name": "test-to-exception-edit-2-icon",
                "color": "#2FAB63",
                "name": "Exception edit 2 State",
                "priority": 14,
                "type": "active"
              }
            }
          ],
          "created": 1592215037
        },
        {
          "_id": "test-exception-to-get-2",
          "name": "test-exception-to-get-2-name",
          "description": "test-exception-to-get-2-description",
          "exdates": [
            {
              "begin": 1591164001,
              "end": 1591167601,
              "type": {
                "_id": "test-type-to-exception-edit-2",
                "description": "Exception edit 2 State type",
                "icon_name": "test-to-exception-edit-2-icon",
                "color": "#2FAB63",
                "name": "Exception edit 2 State",
                "priority": 14,
                "type": "active"
              }
            }
          ],
          "created": 1592215037
        },
        {
          "_id": "test-exception-to-get-3",
          "name": "test-exception-to-get-3-name",
          "description": "test-exception-to-get-3-description",
          "exdates": [
            {
              "begin": 1591164001,
              "end": 1591167601,
              "type": {
                "_id": "test-type-to-exception-edit-3",
                "description": "Exception edit 3 State type",
                "icon_name": "test-to-exception-edit-3-icon",
                "color": "#2FAB63",
                "name": "Exception edit 3 State",
                "priority": 15,
                "type": "active"
              }
            }
          ],
          "created": 1592215037
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 3
      }
    }
    """

  Scenario: given get all request should return exceptions with flags
    When I am admin
    When I do GET /api/v4/pbehavior-exceptions?search=test-exception-to-get&with_flags=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-exception-to-get-1",
          "deletable": true
        },
        {
          "_id": "test-exception-to-get-2",
          "deletable": false
        },
        {
          "_id": "test-exception-to-get-3",
          "deletable": false
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 3
      }
    }
    """

  Scenario: given get all request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/pbehavior-exceptions
    Then the response code should be 403

  Scenario: given get request should return exception
    When I am admin
    When I do GET /api/v4/pbehavior-exceptions/test-exception-to-get-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-exception-to-get-1",
      "name": "test-exception-to-get-1-name",
      "description": "test-exception-to-get-1-description",
      "exdates": [
        {
          "begin": 1591164001,
          "end": 1591167601,
          "type": {
            "_id": "test-type-to-exception-edit-2",
            "description": "Exception edit 2 State type",
            "icon_name": "test-to-exception-edit-2-icon",
            "color": "#2FAB63",
            "name": "Exception edit 2 State",
            "priority": 14,
            "type": "active"
          }
        }
      ],
      "created": 1592215037
    }
    """

  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/pbehavior-exceptions/test-exception-to-get
    Then the response code should be 403

  Scenario: given invalid get request should return error
    When I am admin
    When I do GET /api/v4/pbehavior-exceptions/notexist
    Then the response code should be 404

Feature: Update pbehavior exception
  I need to be able to update a pbehavior exception

  Scenario: given update request should update exception
    When I am admin
    When I do PUT /api/v4/pbehavior-exceptions/test-exception-to-update:
    """json
    {
      "name": "Christmas Update",
      "description": "Public holidays",
      "exdates":[
        {
          "begin": 1591164001,
          "end": 1591167601,
          "type": "test-type-to-exception-edit-2"
        }
      ]
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-exception-to-update",
      "name": "Christmas Update",
      "description": "Public holidays",
      "exdates":[
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

  Scenario: PUT a valid exception without any changes
    When I am admin
    When I do PUT /api/v4/pbehavior-exceptions/test-exception-to-update:
    """json
    {
      "name": "Christmas Update",
      "description": "Public holidays",
      "exdates":[
        {
          "begin": 1591164001,
          "end": 1591167601,
          "type": "test-type-to-exception-edit-2"
        }
      ]
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-exception-to-update",
      "name": "Christmas Update",
      "description": "Public holidays",
      "exdates":[
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

  Scenario: PUT a valid exception with already existed name
    When I am admin
    When I do PUT /api/v4/pbehavior-exceptions/test-exception-to-update:
    """json
    {
      "name": "test-exception-to-get-1-name",
      "description": "Public holidays",
      "exdates":[
        {
          "begin": 1591164001,
          "end": 1591167601,
          "type": "test-type-to-exception-edit-2"
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "name": "Name already exists."
      }
    }
    """

  Scenario: given update request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/pbehavior-exceptions/test-exception-to-update:
    """json
    {
      "name": "Christmas Update noperms",
      "description": "Public holidays",
      "exdates":[
        {
          "begin": 1591164001,
          "end": 1591167601,
          "type": "test-type-to-exception-edit-2"
        }
      ]
    }
    """
    Then the response code should be 403

  Scenario: given no exist exception id should return error
    When I am admin
    When I do PUT /api/v4/pbehavior-exceptions/notexist:
    """json
    {
      "name": "Christmas Update not exist",
      "description": "Public holidays",
      "exdates":[
        {
          "begin": 1591164001,
          "end": 1591167601,
          "type": "test-type-to-exception-edit-1"
        }
      ]
    }
    """
    Then the response code should be 404

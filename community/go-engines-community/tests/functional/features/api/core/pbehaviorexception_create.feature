Feature: Create pbehavior exception
  I need to be able to create a pbehavior exception

  Scenario: given create request should return ok
    When I am admin
    When I do POST /api/v4/pbehavior-exceptions:
    """
    {
      "name": "Christmas Create",
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
    Then the response code should be 201
    Then the response body should contain:
    """
    {
      "name": "Christmas Create",
      "description": "Public holidays",
      "exdates":[
        {
          "begin": 1591164001,
          "end": 1591167601,
          "type": {
            "_id": "test-type-to-exception-edit-1"
          }
        }
      ]
    }
    """

  Scenario: given create request with custom id should return ok
    When I am admin
    When I do POST /api/v4/pbehavior-exceptions:
    """
    {
      "_id": "custom-id",
      "name": "Christmas Create custom-id",
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
    Then the response code should be 201
    When I do GET /api/v4/pbehavior-exceptions/custom-id
    Then the response code should be 200

  Scenario: given create request with custom id that exists should cause dup error
    When I am admin
    When I do POST /api/v4/pbehavior-exceptions:
    """
    {
      "_id": "test-exception-to-update",
      "name": "Christmas Create custom-id 2",
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
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "_id": "ID already exists."
      }
    }
    """

  Scenario: given create request should return exception
    When I am admin
    When I do POST /api/v4/pbehavior-exceptions:
    """
    {
      "name": "Christmas Create check response",
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
    When I do GET /api/v4/pbehavior-exceptions/{{ .lastResponse._id}}
    Then the response code should be 200

  Scenario: given create request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/pbehavior-exceptions:
    """
    {
      "name": "Christmas Create no perms",
      "description": "Public holidays",
      "exdates":[
        {
          "begin": 1591164001,
          "end": 1591167601,
          "type": "test-type-1"
        }
      ]
    }
    """
    Then the response code should be 403

  Scenario: given invalid type id should return error
    When I am admin
    When I do POST /api/v4/pbehavior-exceptions:
    """
    {
      "name": "Christmas Create invalid type",
      "description": "Public holidays",
      "exdates":[
        {
          "begin": 1591164001,
          "end": 1591167601,
          "type": "notexist"
        }
      ]
    }
    """
    Then the response code should be 400

  Scenario: given invalid create request should return errors
    When I am admin
    When I do POST /api/v4/pbehavior-exceptions
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "name": "Name is missing.",
        "description": "Description is missing.",
        "exdates": "Exdates is missing."
      }
    }
    """

  Scenario: given invalid create request should return errors
    When I am admin
    When I do POST /api/v4/pbehavior-exceptions:
    """
    {
      "name": "",
      "description": "",
      "exdates": []
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "name": "Name is missing.",
        "description": "Description is missing.",
        "exdates": "Exdates should not be blank."
      }
    }
    """

  Scenario: given invalid create request should return errors
    When I am admin
    When I do POST /api/v4/pbehavior-exceptions:
    """
    {
      "name": "test",
      "description": "test",
      "exdates": [
        {
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "exdates.0.begin":"Begin is missing.",
        "exdates.0.end":"End is missing.",
        "exdates.0.type":"Type is missing."
      }
    }
    """

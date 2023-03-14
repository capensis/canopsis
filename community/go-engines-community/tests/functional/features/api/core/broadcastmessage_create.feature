Feature: Create an broadcast message
  I need to be able to create a broadcast message
  Only admin should be able to create a broadcast message

  @standalone
  Scenario: given create request should return ok
    When I am admin
    When I connect to websocket
    When I subscribe to websocket room "broadcast-messages"
    When I save response start={{ nowAdd "-2h" }}
    When I save response end={{ nowAdd "24h" }}
    When I do POST /api/v4/broadcast-message:
    """json
    {
      "color": "rgb(159, 13, 0)",
      "message": "test-broadcast-create",
      "start": {{ .start }},
      "end": {{ .end }}
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "color": "rgb(159, 13, 0)",
      "message": "test-broadcast-create",
      "start": {{ .start }},
      "end": {{ .end }}
    }
    """
    When I do GET /api/v4/broadcast-message/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "color": "rgb(159, 13, 0)",
      "message": "test-broadcast-create",
      "start": {{ .start }},
      "end": {{ .end }}
    }
    """
    Then I wait message from websocket room "broadcast-messages" which contains:
    """json
    [
      {
        "_id": "test-broadcast-to-get-4",
        "color": "#e75e40",
        "message": "broadcast-test-to-get-4"
      },
      {
        "color": "rgb(159, 13, 0)",
        "message": "test-broadcast-create",
        "start": {{ .start }},
        "end": {{ .end }}
      }
    ]
    """

  Scenario: given create request with wrong type should return bad request
    When I am admin
    When I do POST /api/v4/broadcast-message:
    """json
    {
      "color": "#23"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "color": "Color is not valid."
      }
    }
    """

  Scenario: given create request with missing fields should return bad request
    When I am admin
    When I do POST /api/v4/broadcast-message:
    """json
    {
      "color": "rgb(159, 13, 0, 23)"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "color": "Color is not valid."
      }
    }
    """

  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/broadcast-message
    Then the response code should be 401

  Scenario: given create request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/broadcast-message
    Then the response code should be 403

  Scenario: given create request with already exists name should return error
    When I am admin
    When I do POST /api/v4/broadcast-message:
    """json
    {
      "_id": "test-broadcast-to-check-unique"
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

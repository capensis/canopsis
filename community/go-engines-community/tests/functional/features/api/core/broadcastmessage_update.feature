Feature: Update an broadcast message
  I need to be able to update an broadcast message
  Only admin should be able to update an broadcast message

  @standalone
  Scenario: given update request should update broadcast message
    When I am admin
    When I connect to websocket
    When I subscribe to websocket room "broadcast-messages"
    When I save response start={{ nowAdd "-2h" }}
    When I save response end={{ nowAdd "24h" }}
    Then I do PUT /api/v4/broadcast-message/test-broadcast-to-update:
    """json
    {
      "color": "#32a852",
      "message": "test-broadcast-to-update-updated",
      "start": {{ .start }},
      "end": {{ .end }},
      "created": 1612139798
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-broadcast-to-update",
      "color": "#32a852",
      "message": "test-broadcast-to-update-updated",
      "start": {{ .start }},
      "end": {{ .end }},
      "created": 1612139798
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
        "color": "#32a852",
        "message": "test-broadcast-to-update-updated",
        "start": {{ .start }},
        "end": {{ .end }}
      }
    ]
    """

  Scenario: given get request and no auth user should not allow access
    When I do PUT /api/v4/broadcast-message/test-broadcast-to-update
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/broadcast-message/test-broadcast-to-update
    Then the response code should be 403

  Scenario: given update request with not exist id should return not found error
    When I am admin
    When I do PUT /api/v4/broadcast-message/test-broadcast-not-found:
    """json
    {
      "color": "#32a852",
      "message": "test-broadcast-to-update-updated",
      "start": 1612139798,
      "end": 2612296900
    }
    """
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """

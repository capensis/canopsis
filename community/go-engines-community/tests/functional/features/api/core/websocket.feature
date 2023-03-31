Feature: Connect to websocket
  I need to be able to connect to websocket
  Only admin should be able to connect to websocket

  Scenario: given invalid message should return error
    When I connect to websocket
    When I send message to websocket:
    """json
    {
      "type": -1
    }
    """
    Then I wait next message from websocket:
    """json
    {
      "type": 2,
      "error": 400,
      "msg": "unknown message type"
    }
    """

  Scenario: given empty auth token should return error
    When I connect to websocket
    When I send message to websocket:
    """json
    {
      "type": 3
    }
    """
    Then I wait next message from websocket:
    """json
    {
      "type": 2,
      "error": 400,
      "msg": "token is missing"
    }
    """

  Scenario: given not exist auth token should return error
    When I connect to websocket
    When I send message to websocket:
    """json
    {
      "type": 3,
      "token": "not-exist"
    }
    """
    Then I wait next message from websocket:
    """json
    {
      "type": 2,
      "error": 401
    }
    """

  Scenario: given logout auth token should return error
    When I connect to websocket
    When I do POST /api/v4/login:
    """json
    {
      "username": "root",
      "password": "test"
    }
    """
    Then the response code should be 200
    When I send message to websocket:
    """json
    {
      "type": 3,
      "token": "{{ .lastResponse.access_token }}"
    }
    """
    Then I wait next message from websocket:
    """json
    {
      "type": 4
    }
    """
    When I set header Authorization=Bearer {{ .lastResponse.access_token }}
    When I do POST /api/v4/logout
    Then the response code should be 204
    Then I wait next message from websocket:
    """json
    {
      "type": 2,
      "error": 401
    }
    """

  Scenario: given empty room should return error
    When I connect to websocket
    When I send message to websocket:
    """json
    {
      "type": 1
    }
    """
    Then I wait next message from websocket:
    """json
    {
      "type": 2,
      "error": 400,
      "msg": "room is missing"
    }
    """
    When I send message to websocket:
    """json
    {
      "type": 2
    }
    """
    Then I wait next message from websocket:
    """json
    {
      "type": 2,
      "error": 400,
      "msg": "room is missing"
    }
    """

  Scenario: given unknown room should return error
    When I connect to websocket
    When I send message to websocket:
    """json
    {
      "type": 1,
      "room": "not-exist"
    }
    """
    Then I wait next message from websocket:
    """json
    {
      "type": 2,
      "error": 404,
      "room": "not-exist"
    }
    """

Feature: Connect to websocket
  I need to be able to connect to websocket
  Only admin should be able to connect to websocket

  Scenario: given not exist auth token and protected room should return error
    When I connect to websocket
    When I send message to websocket:
    """json
    {
      "type": 1,
      "room": "healthcheck"
    }
    """
    Then I wait next message from websocket:
    """json
    {
      "type": 2,
      "error": 401,
      "room": "healthcheck"
    }
    """

  Scenario: given auth token without permissions and protected room should return error
    When I connect to websocket
    When I do POST /api/v4/login:
    """json
    {
      "username": "nopermsuser",
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
    When I send message to websocket:
    """json
    {
      "type": 1,
      "room": "healthcheck"
    }
    """
    Then I wait next message from websocket:
    """json
    {
      "type": 2,
      "error": 403,
      "room": "healthcheck"
    }
    """

  Scenario: given removed permission and protected room should return error
    When I connect to websocket
    When I do POST /api/v4/login:
    """json
    {
      "username": "test-user-to-websocket-connect",
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
    When I send message to websocket:
    """json
    {
      "type": 1,
      "room": "healthcheck"
    }
    """
    When I am admin
    When I do PUT /api/v4/roles/test-role-to-websocket-connect:
    """json
    {
    }
    """
    Then the response code should be 200
    Then I wait message from websocket:
    """json
    {
      "type": 2,
      "error": 403,
      "room": "healthcheck"
    }
    """

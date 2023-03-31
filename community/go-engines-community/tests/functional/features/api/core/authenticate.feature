Feature: Authenticate user
  I need to be able to authenticate user.

  Scenario: given user username and password should log in user
    When I do POST /api/v4/login:
    """json
    {
      "username": "root",
      "password": "test"
    }
    """
    Then the response code should be 200

  Scenario: given auth user should allow access
    When I do POST /api/v4/login:
    """json
    {
      "username": "root",
      "password": "test"
    }
    """
    Then the response code should be 200
    When I set header Authorization=Bearer {{ .lastResponse.access_token }}
    When I do GET /api/v4/alarms
    Then the response code should be 200

  Scenario: given auth user should log out user
    When I do POST /api/v4/login:
    """json
    {
      "username": "root",
      "password": "test"
    }
    """
    When I set header Authorization=Bearer {{ .lastResponse.access_token }}
    When I do POST /api/v4/logout
    Then the response code should be 204

  Scenario: given unauth user should not allow access
    When I do POST /api/v4/login:
    """json
    {
      "username": "root",
      "password": "test"
    }
    """
    When I set header Authorization=Bearer {{ .lastResponse.access_token }}
    When I do POST /api/v4/logout
    Then the response code should be 204
    When I do GET /api/v4/alarms
    Then the response code should be 401

  Scenario: given invalid username and password should return error
    When I do POST /api/v4/login:
    """json
    {
      "username": "nouser",
      "password": "nopass"
    }
    """
    Then the response code should be 401

  Scenario: given no username and password should return errors
    When I do POST /api/v4/login:
    """json
    {}
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "password":"Password is missing.",
        "username":"Username is missing."
      }
    }
	"""

  Scenario: given unauth user should not log out
    When I do POST /api/v4/logout
    Then the response code should be 401

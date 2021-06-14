Feature: Session stats
  I need to be able to collect session stats by user.

  Scenario: given auth user should update stats
    When I do POST /auth:
    """
    {
      "username": "root",
      "password": "test"
    }
    """
    When I do GET /api/v2/sessionstart
    Then the response code should be 200

  Scenario: given auth user should update stats
    When I do POST /auth:
    """
    {
      "username": "root",
      "password": "test"
    }
    """
    When I do POST /api/v2/keepalive:
    """
    {
      "path": ["view", "tab"],
      "visible": true
    }
    """
    Then the response code should be 200

  Scenario: given auth user should update stats
    When I do POST /auth:
    """
    {
      "username": "root",
      "password": "test"
    }
    """
    When I do POST /api/v2/keepalive:
    """
    {
      "path": ["view", "tab"],
      "visible": false
    }
    """
    Then the response code should be 200

  Scenario: given auth user should update stats
    When I do POST /auth:
    """
    {
      "username": "root",
      "password": "test"
    }
    """
    When I do POST /api/v2/session_tracepath:
    """
    {
      "path": ["view", "tab"]
    }
    """
    Then the response code should be 200

  Scenario: given auth user should get stats
    When I am admin
    When I do GET /api/v2/sessions
    Then the response code should be 200

  Scenario: given no session should not allow access to
    When I do GET /api/v2/sessionstart
    Then the response code should be 401

  Scenario: given no session should not allow access to
    When I do POST /api/v2/keepalive
    Then the response code should be 401

  Scenario: given no session should not allow access to
    When I do POST /api/v2/session_tracepath
    Then the response code should be 401

  Scenario: given no auth user should not allow access to
    When I do GET /api/v2/sessions
    Then the response code should be 401

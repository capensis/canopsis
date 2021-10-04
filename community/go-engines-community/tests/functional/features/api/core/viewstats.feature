Feature: View stats
  I need to be able to collect view stats by user.

  Scenario: given auth user should update stats
    When I am admin
    When I do POST /api/v4/view-stats
    Then the response code should be 201

  Scenario: given auth user should update stats
    When I am admin
    When I do POST /api/v4/view-stats
    Then the response code should be 201
    When I do PUT /api/v4/view-stats/{{ .lastResponse._id }}:
    """json
    {
      "path": ["view", "tab"],
      "visible": true
    }
    """
    Then the response code should be 200

  Scenario: given auth user should update stats
    When I am admin
    When I do POST /api/v4/view-stats
    Then the response code should be 201
    When I do PUT /api/v4/view-stats/{{ .lastResponse._id }}:
    """json
    {
      "path": ["view", "tab"],
      "visible": false
    }
    """
    Then the response code should be 200

  Scenario: given auth user should get stats
    When I am admin
    When I do GET /api/v4/view-stats
    Then the response code should be 200

  Scenario: given unauth user should not allow access to
    When I do GET /api/v4/view-stats
    Then the response code should be 401

  Scenario: given unauth user should not allow access to
    When I do POST /api/v4/view-stats
    Then the response code should be 401

  Scenario: given unauth should not allow access to
    When I do PUT /api/v4/view-stats/not-exist
    Then the response code should be 401

  Scenario: given not exist stats should return not found error
    When I am admin
    When I do PUT /api/v4/view-stats/not-exist:
    """json
    {
      "path": ["view", "tab"],
      "visible": false
    }
    """
    Then the response code should be 404

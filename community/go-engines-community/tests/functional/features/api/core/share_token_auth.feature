Feature: Authenticate with a share token
  I need to be able to authenticate a share token
  Only admin should be able to authenticate a share token

  Scenario: given new token should authenticate user
    When I am admin
    When I do POST /api/v4/share-tokens:
    """json
    {
    }
    """
    Then the response code should be 201
    When I set header Authorization=Bearer {{ .lastResponse.value }}
    When I do GET /api/v4/alarms
    Then the response code should be 200

  Scenario: given expired token should not authenticate user
    When I am admin
    When I do POST /api/v4/share-tokens:
    """json
    {
      "duration": {
        "value": 2,
        "unit": "s"
      }
    }
    """
    Then the response code should be 201
    When I set header Authorization=Bearer {{ .lastResponse.value }}
    When I do GET /api/v4/alarms
    Then the response code should be 200
    When I do GET /api/v4/alarms until response code is 401

  Scenario: given share token should not logout
    When I am admin
    When I do POST /api/v4/share-tokens:
    """json
    {}
    """
    Then the response code should be 201
    When I set header Authorization=Bearer {{ .lastResponse.value }}
    When I do GET /api/v4/alarms
    Then the response code should be 200
    When I do POST /api/v4/logout
    Then the response code should be 401

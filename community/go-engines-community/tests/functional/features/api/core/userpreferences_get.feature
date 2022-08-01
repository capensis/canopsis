Feature: Get user preferences
  I need to be able to get user preferences.

  Scenario: given user should get user preferences
    When I am test-role-to-user-preferences-edit
    When I do GET /api/v4/user-preferences/test-widget-to-user-preferences-get-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "widget": "test-widget-to-user-preferences-get-1",
      "content": {
        "test-int": 1,
        "test-str": "test",
        "test-array": ["test1", "test2"],
        "test-map": {
          "nested": 1
        }
      }
    }
	"""
    When I do GET /api/v4/user-preferences/test-widget-to-user-preferences-get-2
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "widget": "test-widget-to-user-preferences-get-2",
      "content": {
        "test-int": 1,
        "test-str": "test",
        "test-array": ["test1", "test2"],
        "test-map": {
          "nested": 1
        }
      }
    }
	"""

  Scenario: given get request with not exist id should return not found error
    When I am test-role-to-user-preferences-edit
    When I do GET /api/v4/user-preferences/not-found
    Then the response code should be 403

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/user-preferences/not-found
    Then the response code should be 401

  Scenario: given get request and auth user without view permission should not allow access
    When I am admin
    When I do GET /api/v4/user-preferences/test-widget-to-user-preferences-get-1
    Then the response code should be 403

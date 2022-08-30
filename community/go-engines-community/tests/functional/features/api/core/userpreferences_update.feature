Feature: Update user preferences
  I need to be able to update user preferences.

  Scenario: given user should update user preferences
    When I am test-role-to-user-preferences-edit
    When I do GET /api/v4/user-preferences/test-widget-to-user-preferences-update-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "widget": "test-widget-to-user-preferences-update-1",
      "content": {},
      "filters": []
    }
	"""
    When I do PUT /api/v4/user-preferences:
    """json
    {
      "widget": "test-widget-to-user-preferences-update-1",
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
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "widget": "test-widget-to-user-preferences-update-1",
      "content": {
        "test-int": 1,
        "test-str": "test",
        "test-array": ["test1", "test2"],
        "test-map": {
          "nested": 1
        }
      },
      "filters": []
    }
	"""
    When I do GET /api/v4/user-preferences/test-widget-to-user-preferences-update-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "widget": "test-widget-to-user-preferences-update-1",
      "content": {
        "test-int": 1,
        "test-str": "test",
        "test-array": ["test1", "test2"],
        "test-map": {
          "nested": 1
        }
      },
      "filters": []
    }
	"""
    When I do PUT /api/v4/user-preferences:
    """json
    {
      "widget": "test-widget-to-user-preferences-update-1",
      "content": {
        "test-int": 2,
        "test-str": "test-updated",
        "test-array": ["test1-updated", "test2-updated"],
        "test-map": {
          "nested": 3
        }
      }
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "widget": "test-widget-to-user-preferences-update-1",
      "content": {
        "test-int": 2,
        "test-str": "test-updated",
        "test-array": ["test1-updated", "test2-updated"],
        "test-map": {
          "nested": 3
        }
      },
      "filters": []
    }
	"""
    When I do GET /api/v4/user-preferences/test-widget-to-user-preferences-update-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "widget": "test-widget-to-user-preferences-update-1",
      "content": {
        "test-int": 2,
        "test-str": "test-updated",
        "test-array": ["test1-updated", "test2-updated"],
        "test-map": {
          "nested": 3
        }
      },
      "filters": []
    }
	"""

  Scenario: given update request with not exist id should return not found error
    When I am test-role-to-user-preferences-edit
    When I do PUT /api/v4/user-preferences:
    """json
    {
      "widget": "not-found",
      "content": {
        "test": 1
      }
    }
    """
    Then the response code should be 403

  Scenario: given update request and no auth user should not allow access
    When I do PUT /api/v4/user-preferences
    Then the response code should be 401

  Scenario: given update request and auth user without view permission should not allow access
    When I am test-role-to-user-preferences-edit
    When I do PUT /api/v4/user-preferences:
    """json
    {
      "widget": "test-widget-to-user-preferences-update-2",
      "content": {
        "test-int": 1
      }
    }
    """
    Then the response code should be 403

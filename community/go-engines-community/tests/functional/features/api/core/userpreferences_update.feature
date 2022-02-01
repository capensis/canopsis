Feature: Update user preferences
  I need to be able to update user preferences.

  Scenario: given user should update user preferences
    When I am test-role-to-user-preferences-update
    When I do GET /api/v4/user-preferences/test-view-to-update-user-preferences-tab-1-widget-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "widget": "test-view-to-update-user-preferences-tab-1-widget-1",
      "content": {}
    }
	"""
    When I do PUT /api/v4/user-preferences:
    """json
    {
      "widget": "test-view-to-update-user-preferences-tab-1-widget-1",
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
      "widget": "test-view-to-update-user-preferences-tab-1-widget-1",
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
    When I do GET /api/v4/user-preferences/test-view-to-update-user-preferences-tab-1-widget-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "widget": "test-view-to-update-user-preferences-tab-1-widget-1",
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
    When I do PUT /api/v4/user-preferences:
    """json
    {
      "widget": "test-view-to-update-user-preferences-tab-1-widget-1",
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
      "widget": "test-view-to-update-user-preferences-tab-1-widget-1",
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
    When I do GET /api/v4/user-preferences/test-view-to-update-user-preferences-tab-1-widget-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "widget": "test-view-to-update-user-preferences-tab-1-widget-1",
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

  Scenario: given update request with not exist id should return not found error
    When I am test-role-to-user-preferences-update
    When I do PUT /api/v4/user-preferences:
    """json
    {
      "widget": "not-found",
      "content": {
        "test": 1
      }
    }
    """
    Then the response code should be 404

  Scenario: given get request with not exist id should return not found error
    When I am test-role-to-user-preferences-update
    When I do GET /api/v4/user-preferences/not-found
    Then the response code should be 404

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/user-preferences/not-found
    Then the response code should be 401

  Scenario: given update request and no auth user should not allow access
    When I do PUT /api/v4/user-preferences
    Then the response code should be 401

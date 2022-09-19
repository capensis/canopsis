Feature: Delete user preferences on widget delete.
  I need to be able to delete user preferences.

  Scenario: given deleted view should delete user preferences
    When I am test-role-to-user-preferences-edit
    When I do POST /api/v4/views:
    """json
    {
      "enabled": true,
      "title": "test-view-to-user-preferences-delete-1-title",
      "description": "test-view-to-user-preferences-delete-1-description",
      "group": "test-viewgroup-to-user-preferences-delete"
    }
    """
    Then the response code should be 201
    When I save response viewID={{ .lastResponse._id }}
    When I do POST /api/v4/view-tabs:
    """json
    {
      "title": "test-tab-to-user-preferences-delete-1-title",
      "view": "{{ .viewID }}"
    }
    """
    Then the response code should be 201
    When I save response tabID={{ .lastResponse._id }}
    When I do POST /api/v4/widgets:
    """json
    {
      "title": "test-widget-to-user-preferences-delete-1-title",
      "type": "test-widget-to-user-preferences-delete-1-type",
      "tab": "{{ .tabID }}"
    }
    """
    Then the response code should be 201
    When I save response widgetID={{ .lastResponse._id }}
    When I do PUT /api/v4/user-preferences:
    """json
    {
      "widget": "{{ .widgetID }}",
      "content": {
        "test": 1
      }
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/user-preferences/{{ .widgetID }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "widget": "{{ .widgetID }}",
      "content": {
        "test": 1
      },
      "filters": []
    }
	"""
    When I do DELETE /api/v4/views/{{ .viewID }}
    Then the response code should be 204
    When I do GET /api/v4/user-preferences/{{ .widgetID }}
    Then the response code should be 403

  Scenario: given deleted widget should delete user preferences
    When I am test-role-to-user-preferences-edit
    When I do POST /api/v4/views:
    """json
    {
      "enabled": true,
      "title": "test-view-to-user-preferences-delete-2-title",
      "description": "test-view-to-user-preferences-delete-2-description",
      "group": "test-viewgroup-to-user-preferences-delete"
    }
    """
    Then the response code should be 201
    When I save response viewID={{ .lastResponse._id }}
    When I do POST /api/v4/view-tabs:
    """json
    {
      "title": "test-tab-to-user-preferences-delete-2-title",
      "view": "{{ .viewID }}"
    }
    """
    Then the response code should be 201
    When I save response tabID={{ .lastResponse._id }}
    When I do POST /api/v4/widgets:
    """json
    {
      "title": "test-widget-to-user-preferences-delete-2-1-title",
      "type": "test-widget-to-user-preferences-delete-2-1-type",
      "tab": "{{ .tabID }}"
    }
    """
    Then the response code should be 201
    When I save response widget1ID={{ .lastResponse._id }}
    When I do POST /api/v4/widgets:
    """json
    {
      "title": "test-widget-to-user-preferences-delete-2-2-title",
      "type": "test-widget-to-user-preferences-delete-2-2-type",
      "tab": "{{ .tabID }}"
    }
    """
    Then the response code should be 201
    When I save response widget2ID={{ .lastResponse._id }}
    When I do PUT /api/v4/user-preferences:
    """json
    {
      "widget": "{{ .widget1ID }}",
      "content": {
        "test": 1
      }
    }
    """
    Then the response code should be 200
    When I do PUT /api/v4/user-preferences:
    """json
    {
      "widget": "{{ .widget2ID }}",
      "content": {
        "test": 2
      }
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/user-preferences/{{ .widget1ID }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "widget": "{{ .widget1ID }}",
      "content": {
        "test": 1
      },
      "filters": []
    }
	"""
    When I do GET /api/v4/user-preferences/{{ .widget2ID }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "widget": "{{ .widget2ID }}",
      "content": {
        "test": 2
      },
      "filters": []
    }
	"""
    When I do DELETE /api/v4/widgets/{{ .widget2ID }}
    Then the response code should be 204
    When I do GET /api/v4/user-preferences/{{ .widget1ID }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "widget": "{{ .widget1ID }}",
      "content": {
        "test": 1
      },
      "filters": []
    }
	"""
    When I do GET /api/v4/user-preferences/{{ .widget2ID }}
    Then the response code should be 403

  Scenario: given deleted tab should delete user preferences
    When I am test-role-to-user-preferences-edit
    When I do POST /api/v4/views:
    """json
    {
      "enabled": true,
      "title": "test-view-to-user-preferences-delete-3-title",
      "description": "test-view-to-user-preferences-delete-3-description",
      "group": "test-viewgroup-to-user-preferences-delete"
    }
    """
    Then the response code should be 201
    When I save response viewID={{ .lastResponse._id }}
    When I do POST /api/v4/view-tabs:
    """json
    {
      "title": "test-tab-to-user-preferences-delete-3-title",
      "view": "{{ .viewID }}"
    }
    """
    Then the response code should be 201
    When I save response tabID={{ .lastResponse._id }}
    When I do POST /api/v4/widgets:
    """json
    {
      "title": "test-widget-to-user-preferences-delete-3-title",
      "type": "test-widget-to-user-preferences-delete-3-type",
      "tab": "{{ .tabID }}"
    }
    """
    Then the response code should be 201
    When I save response widgetID={{ .lastResponse._id }}
    When I do PUT /api/v4/user-preferences:
    """json
    {
      "widget": "{{ .widgetID }}",
      "content": {
        "test": 1
      },
      "filters": []
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/user-preferences/{{ .widgetID }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "widget": "{{ .widgetID }}",
      "content": {
        "test": 1
      },
      "filters": []
    }
	"""
    When I do DELETE /api/v4/view-tabs/{{ .tabID }}
    Then the response code should be 204
    When I do GET /api/v4/user-preferences/{{ .widgetID }}
    Then the response code should be 403

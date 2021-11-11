Feature: Delete user preferences on widget delete.
  I need to be able to delete user preferences.

  Scenario: given deleted view should delete user preferences
    When I am test-role-to-user-preferences-delete-1
    When I do POST /api/v4/views:
    """json
    {
      "enabled": true,
      "title": "test-view-to-user-preferences-delete-1-title",
      "description": "test-view-to-user-preferences-delete-1-description",
      "group": "test-viewgroup-to-user-preferences-delete",
      "tabs": [
        {
          "_id": "test-view-to-user-preferences-delete-1-tab-1",
          "title": "test-view-to-user-preferences-delete-1-tab-1-title",
          "widgets": [
            {
              "_id": "test-view-to-user-preferences-delete-1-tab-1-widget-1",
              "title": "test-view-to-user-preferences-delete-1-tab-1-widget-1-title",
              "type": "test-view-to-user-preferences-delete-1-tab-1-widget-1-type"
            }
          ]
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response viewID={{ .lastResponse._id }}
    When I do PUT /api/v4/user-preferences:
    """json
    {
      "widget": "test-view-to-user-preferences-delete-1-tab-1-widget-1",
      "content": {
        "test": 1
      }
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/user-preferences/test-view-to-user-preferences-delete-1-tab-1-widget-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "widget": "test-view-to-user-preferences-delete-1-tab-1-widget-1",
      "content": {
        "test": 1
      }
    }
	"""
    When I do DELETE /api/v4/views/{{ .viewID }}
    Then the response code should be 204
    When I do GET /api/v4/user-preferences/test-view-to-user-preferences-delete-1-tab-1-widget-1
    Then the response code should be 404

  Scenario: given deleted widget should delete user preferences
    When I am test-role-to-user-preferences-delete-2
    When I do POST /api/v4/views:
    """json
    {
      "enabled": true,
      "title": "test-view-to-user-preferences-delete-2-title",
      "description": "test-view-to-user-preferences-delete-2-description",
      "group": "test-viewgroup-to-user-preferences-delete",
      "tabs": [
        {
          "_id": "test-view-to-user-preferences-delete-2-tab-1",
          "title": "test-view-to-user-preferences-delete-2-tab-1-title",
          "widgets": [
            {
              "_id": "test-view-to-user-preferences-delete-2-tab-1-widget-1",
              "title": "test-view-to-user-preferences-delete-2-tab-1-widget-1-title",
              "type": "test-view-to-user-preferences-delete-2-tab-1-widget-1-type"
            },
            {
              "_id": "test-view-to-user-preferences-delete-2-tab-1-widget-2",
              "title": "test-view-to-user-preferences-delete-2-tab-1-widget-2-title",
              "type": "test-view-to-user-preferences-delete-2-tab-1-widget-2-type"
            }
          ]
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response viewID={{ .lastResponse._id }}
    When I do PUT /api/v4/user-preferences:
    """json
    {
      "widget": "test-view-to-user-preferences-delete-2-tab-1-widget-1",
      "content": {
        "test": 1
      }
    }
    """
    Then the response code should be 200
    When I do PUT /api/v4/user-preferences:
    """json
    {
      "widget": "test-view-to-user-preferences-delete-2-tab-1-widget-2",
      "content": {
        "test": 2
      }
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/user-preferences/test-view-to-user-preferences-delete-2-tab-1-widget-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "widget": "test-view-to-user-preferences-delete-2-tab-1-widget-1",
      "content": {
        "test": 1
      }
    }
	"""
    When I do GET /api/v4/user-preferences/test-view-to-user-preferences-delete-2-tab-1-widget-2
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "widget": "test-view-to-user-preferences-delete-2-tab-1-widget-2",
      "content": {
        "test": 2
      }
    }
	"""
    When I do PUT /api/v4/views/{{ .viewID }}:
    """json
    {
      "enabled": true,
      "title": "test-view-to-user-preferences-delete-2-title",
      "description": "test-view-to-user-preferences-delete-2-description",
      "group": "test-viewgroup-to-user-preferences-delete",
      "tabs": [
        {
          "_id": "test-view-to-user-preferences-delete-2-tab-1",
          "title": "test-view-to-user-preferences-delete-2-tab-1-title",
          "widgets": [
            {
              "_id": "test-view-to-user-preferences-delete-2-tab-1-widget-1",
              "title": "test-view-to-user-preferences-delete-2-tab-1-widget-1-title",
              "type": "test-view-to-user-preferences-delete-2-tab-1-widget-1-type"
            }
          ]
        }
      ]
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/user-preferences/test-view-to-user-preferences-delete-2-tab-1-widget-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "widget": "test-view-to-user-preferences-delete-2-tab-1-widget-1",
      "content": {
        "test": 1
      }
    }
	"""
    When I do GET /api/v4/user-preferences/test-view-to-user-preferences-delete-2-tab-1-widget-2
    Then the response code should be 404

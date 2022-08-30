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
      "content": {}
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
      }
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
      }
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
      }
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
      }
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

  Scenario: given updated widget filters should update user preferences with object main filter
    When I am test-role-to-user-preferences-edit
    When I do POST /api/v4/widgets:
    """json
    {
      "title": "test-widget-to-user-preferences-update-3-title",
      "tab": "test-tab-to-user-preferences-update-3",
      "type": "test-widget-to-user-preferences-update-3-type",
      "parameters": {
        "viewFilters": [
          {
            "title": "test-widget-to-user-preferences-update-3-filter-1-title",
            "filter": "test-widget-to-user-preferences-update-3-filter-1-value"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I save response widgetID={{ .lastResponse._id }}
    When I do PUT /api/v4/user-preferences:
    """json
    {
      "widget": "{{ .widgetID }}",
      "content": {
        "mainFilter": {
          "title": "test-widget-to-user-preferences-update-3-filter-1-title",
          "filter": "test-widget-to-user-preferences-update-3-filter-1-value"
        }
      }
    }
    """
    Then the response code should be 200
    Then I do PUT /api/v4/widgets/{{ .widgetID }}:
    """json
    {
      "title": "test-widget-to-user-preferences-update-3-title",
      "tab": "test-tab-to-user-preferences-update-3",
      "type": "test-widget-to-user-preferences-update-3-type",
      "parameters": {
        "viewFilters": [
          {
            "title": "test-widget-to-user-preferences-update-3-filter-1-title",
            "filter": "test-widget-to-user-preferences-update-3-filter-1-value"
          },
          {
            "title": "test-widget-to-user-preferences-update-3-filter-2-title",
            "filter": "test-widget-to-user-preferences-update-3-filter-2-value"
          }
        ]
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
        "mainFilter": {
          "title": "test-widget-to-user-preferences-update-3-filter-1-title",
          "filter": "test-widget-to-user-preferences-update-3-filter-1-value"
        }
      }
    }
    """
    Then I do PUT /api/v4/widgets/{{ .widgetID }}:
    """json
    {
      "title": "test-widget-to-user-preferences-update-3-title",
      "tab": "test-tab-to-user-preferences-update-3",
      "type": "test-widget-to-user-preferences-update-3-type",
      "parameters": {
        "viewFilters": [
          {
            "title": "test-widget-to-user-preferences-update-3-filter-1-title",
            "filter": "test-widget-to-user-preferences-update-3-filter-1-value-updated"
          },
          {
            "title": "test-widget-to-user-preferences-update-3-filter-2-title",
            "filter": "test-widget-to-user-preferences-update-3-filter-2-value"
          }
        ]
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
        "mainFilter": {
          "title": "test-widget-to-user-preferences-update-3-filter-1-title",
          "filter": "test-widget-to-user-preferences-update-3-filter-1-value-updated"
        }
      }
    }
    """
    Then I do PUT /api/v4/widgets/{{ .widgetID }}:
    """json
    {
      "title": "test-widget-to-user-preferences-update-3-title",
      "tab": "test-tab-to-user-preferences-update-3",
      "type": "test-widget-to-user-preferences-update-3-type",
      "parameters": {
        "viewFilters": [
          {
            "title": "test-widget-to-user-preferences-update-3-filter-1-title",
            "filter": "test-widget-to-user-preferences-update-3-filter-1-value-updated"
          }
        ]
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
        "mainFilter": {
          "title": "test-widget-to-user-preferences-update-3-filter-1-title",
          "filter": "test-widget-to-user-preferences-update-3-filter-1-value-updated"
        }
      }
    }
    """
    Then I do PUT /api/v4/widgets/{{ .widgetID }}:
    """json
    {
      "title": "test-widget-to-user-preferences-update-3-title",
      "tab": "test-tab-to-user-preferences-update-3",
      "type": "test-widget-to-user-preferences-update-3-type",
      "parameters": {
        "viewFilters": [
          {
            "title": "test-widget-to-user-preferences-update-3-filter-1-title-updated",
            "filter": "test-widget-to-user-preferences-update-3-filter-1-value-updated"
          }
        ]
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
        "mainFilter": {
          "title": "test-widget-to-user-preferences-update-3-filter-1-title",
          "filter": "test-widget-to-user-preferences-update-3-filter-1-value-updated"
        }
      }
    }
    """
    Then I do PUT /api/v4/widgets/{{ .widgetID }}:
    """json
    {
      "title": "test-widget-to-user-preferences-update-3-title",
      "tab": "test-tab-to-user-preferences-update-3",
      "type": "test-widget-to-user-preferences-update-3-type",
      "parameters": {
        "viewFilters": []
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
        "mainFilter": {
          "title": "test-widget-to-user-preferences-update-3-filter-1-title",
          "filter": "test-widget-to-user-preferences-update-3-filter-1-value-updated"
        }
      }
    }
    """

  Scenario: given updated widget filters should update user preferences with array main filter
    When I am test-role-to-user-preferences-edit
    Then I do POST /api/v4/widgets:
    """json
    {
      "title": "test-widget-to-user-preferences-update-4-title",
      "tab": "test-tab-to-user-preferences-update-3",
      "type": "test-widget-to-user-preferences-update-4-type",
      "parameters": {
        "viewFilters": [
          {
            "title": "test-widget-to-user-preferences-update-4-filter-1-title",
            "filter": "test-widget-to-user-preferences-update-4-filter-1-value"
          },
          {
            "title": "test-widget-to-user-preferences-update-4-filter-2-title",
            "filter": "test-widget-to-user-preferences-update-4-filter-2-value"
          },
          {
            "title": "test-widget-to-user-preferences-update-4-filter-3-title",
            "filter": "test-widget-to-user-preferences-update-4-filter-3-value"
          },
          {
            "title": "test-widget-to-user-preferences-update-4-filter-4-title",
            "filter": "test-widget-to-user-preferences-update-4-filter-4-value"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    When I save response widgetID={{ .lastResponse._id }}
    When I do PUT /api/v4/user-preferences:
    """json
    {
      "widget": "{{ .widgetID }}",
      "content": {
        "mainFilter": [
          {
            "title": "test-widget-to-user-preferences-update-4-filter-1-title",
            "filter": "test-widget-to-user-preferences-update-4-filter-1-value"
          },
          {
            "title": "test-widget-to-user-preferences-update-4-filter-2-title",
            "filter": "test-widget-to-user-preferences-update-4-filter-2-value"
          },
          {
            "title": "test-widget-to-user-preferences-update-4-filter-3-title",
            "filter": "test-widget-to-user-preferences-update-4-filter-3-value"
          }
        ]
      }
    }
    """
    Then the response code should be 200
    Then I do PUT /api/v4/widgets/{{ .widgetID }}:
    """json
    {
      "title": "test-widget-to-user-preferences-update-4-title",
      "tab": "test-tab-to-user-preferences-update-3",
      "type": "test-widget-to-user-preferences-update-4-type",
      "parameters": {
        "viewFilters": [
          {
            "title": "test-widget-to-user-preferences-update-4-filter-1-title",
            "filter": "test-widget-to-user-preferences-update-4-filter-1-value-updated"
          },
          {
            "title": "test-widget-to-user-preferences-update-4-filter-2-title-updated",
            "filter": "test-widget-to-user-preferences-update-4-filter-2-value-updated"
          },
          {
            "title": "test-widget-to-user-preferences-update-4-filter-3-title",
            "filter": "test-widget-to-user-preferences-update-4-filter-3-value"
          },
          {
            "title": "test-widget-to-user-preferences-update-4-filter-4-title",
            "filter": "test-widget-to-user-preferences-update-4-filter-4-value"
          }
        ]
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
        "mainFilter": [
          {
            "title": "test-widget-to-user-preferences-update-4-filter-1-title",
            "filter": "test-widget-to-user-preferences-update-4-filter-1-value-updated"
          },
          {
            "title": "test-widget-to-user-preferences-update-4-filter-2-title",
            "filter": "test-widget-to-user-preferences-update-4-filter-2-value"
          },
          {
            "title": "test-widget-to-user-preferences-update-4-filter-3-title",
            "filter": "test-widget-to-user-preferences-update-4-filter-3-value"
          }
        ]
      }
    }
    """
    Then I do PUT /api/v4/widgets/{{ .widgetID }}:
    """json
    {
      "title": "test-widget-to-user-preferences-update-4-title",
      "tab": "test-tab-to-user-preferences-update-3",
      "type": "test-widget-to-user-preferences-update-4-type",
      "parameters": {
        "viewFilters": [
          {
            "title": "test-widget-to-user-preferences-update-4-filter-1-title",
            "filter": "test-widget-to-user-preferences-update-4-filter-1-value-updated"
          },
          {
            "title": "test-widget-to-user-preferences-update-4-filter-2-title-updated",
            "filter": "test-widget-to-user-preferences-update-4-filter-2-value-updated"
          },
          {
            "title": "test-widget-to-user-preferences-update-4-filter-3-title",
            "filter": "test-widget-to-user-preferences-update-4-filter-3-value"
          }
        ]
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
        "mainFilter": [
          {
            "title": "test-widget-to-user-preferences-update-4-filter-1-title",
            "filter": "test-widget-to-user-preferences-update-4-filter-1-value-updated"
          },
          {
            "title": "test-widget-to-user-preferences-update-4-filter-2-title",
            "filter": "test-widget-to-user-preferences-update-4-filter-2-value"
          },
          {
            "title": "test-widget-to-user-preferences-update-4-filter-3-title",
            "filter": "test-widget-to-user-preferences-update-4-filter-3-value"
          }
        ]
      }
    }
    """
    Then I do PUT /api/v4/widgets/{{ .widgetID }}:
    """json
    {
      "title": "test-widget-to-user-preferences-update-4-title",
      "tab": "test-tab-to-user-preferences-update-3",
      "type": "test-widget-to-user-preferences-update-4-type",
      "parameters": {
        "viewFilters": []
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
        "mainFilter": [
          {
            "title": "test-widget-to-user-preferences-update-4-filter-1-title",
            "filter": "test-widget-to-user-preferences-update-4-filter-1-value-updated"
          },
          {
            "title": "test-widget-to-user-preferences-update-4-filter-2-title",
            "filter": "test-widget-to-user-preferences-update-4-filter-2-value"
          },
          {
            "title": "test-widget-to-user-preferences-update-4-filter-3-title",
            "filter": "test-widget-to-user-preferences-update-4-filter-3-value"
          }
        ]
      }
    }
    """

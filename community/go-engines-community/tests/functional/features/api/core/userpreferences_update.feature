Feature: Update user preferences
  I need to be able to update user preferences.

  Scenario: given user should update user preferences
    When I am test-role-to-user-preferences-update-1
    When I do GET /api/v4/user-preferences/test-view-to-update-user-preferences-1-tab-1-widget-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "widget": "test-view-to-update-user-preferences-1-tab-1-widget-1",
      "content": {}
    }
	"""
    When I do PUT /api/v4/user-preferences:
    """json
    {
      "widget": "test-view-to-update-user-preferences-1-tab-1-widget-1",
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
      "widget": "test-view-to-update-user-preferences-1-tab-1-widget-1",
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
    When I do GET /api/v4/user-preferences/test-view-to-update-user-preferences-1-tab-1-widget-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "widget": "test-view-to-update-user-preferences-1-tab-1-widget-1",
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
      "widget": "test-view-to-update-user-preferences-1-tab-1-widget-1",
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
      "widget": "test-view-to-update-user-preferences-1-tab-1-widget-1",
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
    When I do GET /api/v4/user-preferences/test-view-to-update-user-preferences-1-tab-1-widget-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "widget": "test-view-to-update-user-preferences-1-tab-1-widget-1",
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
    When I am test-role-to-user-preferences-update-1
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
    When I am test-role-to-user-preferences-update-1
    When I do GET /api/v4/user-preferences/not-found
    Then the response code should be 404

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/user-preferences/not-found
    Then the response code should be 401

  Scenario: given update request and no auth user should not allow access
    When I do PUT /api/v4/user-preferences
    Then the response code should be 401

  Scenario: given updated widget filters should update user preferences with object main filter
    When I am test-role-to-user-preferences-update-2
    When I do POST /api/v4/views:
    """json
    {
      "title": "test-view-to-update-user-preferences-2-title",
      "description": "test-view-to-update-user-preferences-2-description",
      "enabled": true,
      "group": "test-viewgroup-to-view-edit",
      "tabs": [
        {
          "_id": "test-view-to-update-user-preferences-2-tab-1",
          "title": "test-view-to-update-user-preferences-2-tab-1-title",
          "widgets": [
            {
              "_id": "test-view-to-update-user-preferences-2-tab-1-widget-1",
              "title": "test-view-to-update-user-preferences-2-tab-1-widget-1-title",
              "type": "test-view-to-update-user-preferences-2-tab-1-widget-1-type",
              "parameters": {
                "viewFilters": [
                  {
                    "title": "test-view-to-update-user-preferences-2-tab-1-widget-1-filter-1-title",
                    "filter": "test-view-to-update-user-preferences-2-tab-1-widget-1-filter-1-value"
                  }
                ]
              }
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
      "widget": "test-view-to-update-user-preferences-2-tab-1-widget-1",
      "content": {
        "mainFilter": {
          "title": "test-view-to-update-user-preferences-2-tab-1-widget-1-filter-1-title",
          "filter": "test-view-to-update-user-preferences-2-tab-1-widget-1-filter-1-value"
        }
      }
    }
    """
    Then the response code should be 200
    When I do PUT /api/v4/views/{{ .viewID }}:
    """json
    {
      "title": "test-view-to-update-user-preferences-2-title",
      "description": "test-view-to-update-user-preferences-2-description",
      "enabled": true,
      "group": "test-viewgroup-to-view-edit",
      "tabs": [
        {
          "_id": "test-view-to-update-user-preferences-2-tab-1",
          "title": "test-view-to-update-user-preferences-2-tab-1-title",
          "widgets": [
            {
              "_id": "test-view-to-update-user-preferences-2-tab-1-widget-1",
              "title": "test-view-to-update-user-preferences-2-tab-1-widget-1-title",
              "type": "test-view-to-update-user-preferences-2-tab-1-widget-1-type",
              "parameters": {
                "viewFilters": [
                  {
                    "title": "test-view-to-update-user-preferences-2-tab-1-widget-1-filter-1-title",
                    "filter": "test-view-to-update-user-preferences-2-tab-1-widget-1-filter-1-value"
                  },
                  {
                    "title": "test-view-to-update-user-preferences-2-tab-1-widget-1-filter-2-title",
                    "filter": "test-view-to-update-user-preferences-2-tab-1-widget-1-filter-2-value"
                  }
                ]
              }
            }
          ]
        }
      ]
    }
    """
    When I do GET /api/v4/user-preferences/test-view-to-update-user-preferences-2-tab-1-widget-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "widget": "test-view-to-update-user-preferences-2-tab-1-widget-1",
      "content": {
        "mainFilter": {
          "title": "test-view-to-update-user-preferences-2-tab-1-widget-1-filter-1-title",
          "filter": "test-view-to-update-user-preferences-2-tab-1-widget-1-filter-1-value"
        }
      }
    }
    """
    When I do PUT /api/v4/views/{{ .viewID }}:
    """json
    {
      "title": "test-view-to-update-user-preferences-2-title",
      "description": "test-view-to-update-user-preferences-2-description",
      "enabled": true,
      "group": "test-viewgroup-to-view-edit",
      "tabs": [
        {
          "_id": "test-view-to-update-user-preferences-2-tab-1",
          "title": "test-view-to-update-user-preferences-2-tab-1-title",
          "widgets": [
            {
              "_id": "test-view-to-update-user-preferences-2-tab-1-widget-1",
              "title": "test-view-to-update-user-preferences-2-tab-1-widget-1-title",
              "type": "test-view-to-update-user-preferences-2-tab-1-widget-1-type",
              "parameters": {
                "viewFilters": [
                  {
                    "title": "test-view-to-update-user-preferences-2-tab-1-widget-1-filter-1-title",
                    "filter": "test-view-to-update-user-preferences-2-tab-1-widget-1-filter-1-value-updated"
                  },
                  {
                    "title": "test-view-to-update-user-preferences-2-tab-1-widget-1-filter-2-title",
                    "filter": "test-view-to-update-user-preferences-2-tab-1-widget-1-filter-2-value"
                  }
                ]
              }
            }
          ]
        }
      ]
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/user-preferences/test-view-to-update-user-preferences-2-tab-1-widget-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "widget": "test-view-to-update-user-preferences-2-tab-1-widget-1",
      "content": {
        "mainFilter": {
          "title": "test-view-to-update-user-preferences-2-tab-1-widget-1-filter-1-title",
          "filter": "test-view-to-update-user-preferences-2-tab-1-widget-1-filter-1-value-updated"
        }
      }
    }
    """
    When I do PUT /api/v4/views/{{ .viewID }}:
    """json
    {
      "title": "test-view-to-update-user-preferences-2-title",
      "description": "test-view-to-update-user-preferences-2-description",
      "enabled": true,
      "group": "test-viewgroup-to-view-edit",
      "tabs": [
        {
          "_id": "test-view-to-update-user-preferences-2-tab-1",
          "title": "test-view-to-update-user-preferences-2-tab-1-title",
          "widgets": [
            {
              "_id": "test-view-to-update-user-preferences-2-tab-1-widget-1",
              "title": "test-view-to-update-user-preferences-2-tab-1-widget-1-title",
              "type": "test-view-to-update-user-preferences-2-tab-1-widget-1-type",
              "parameters": {
                "viewFilters": [
                  {
                    "title": "test-view-to-update-user-preferences-2-tab-1-widget-1-filter-1-title",
                    "filter": "test-view-to-update-user-preferences-2-tab-1-widget-1-filter-1-value-updated"
                  }
                ]
              }
            }
          ]
        }
      ]
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/user-preferences/test-view-to-update-user-preferences-2-tab-1-widget-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "widget": "test-view-to-update-user-preferences-2-tab-1-widget-1",
      "content": {
        "mainFilter": {
          "title": "test-view-to-update-user-preferences-2-tab-1-widget-1-filter-1-title",
          "filter": "test-view-to-update-user-preferences-2-tab-1-widget-1-filter-1-value-updated"
        }
      }
    }
    """
    When I do PUT /api/v4/views/{{ .viewID }}:
    """json
    {
      "title": "test-view-to-update-user-preferences-2-title",
      "description": "test-view-to-update-user-preferences-2-description",
      "enabled": true,
      "group": "test-viewgroup-to-view-edit",
      "tabs": [
        {
          "_id": "test-view-to-update-user-preferences-2-tab-1",
          "title": "test-view-to-update-user-preferences-2-tab-1-title",
          "widgets": [
            {
              "_id": "test-view-to-update-user-preferences-2-tab-1-widget-1",
              "title": "test-view-to-update-user-preferences-2-tab-1-widget-1-title",
              "type": "test-view-to-update-user-preferences-2-tab-1-widget-1-type",
              "parameters": {
                "viewFilters": [
                  {
                    "title": "test-view-to-update-user-preferences-2-tab-1-widget-1-filter-1-title-updated",
                    "filter": "test-view-to-update-user-preferences-2-tab-1-widget-1-filter-1-value-updated"
                  }
                ]
              }
            }
          ]
        }
      ]
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/user-preferences/test-view-to-update-user-preferences-2-tab-1-widget-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "widget": "test-view-to-update-user-preferences-2-tab-1-widget-1",
      "content": {
        "mainFilter": {
          "title": "test-view-to-update-user-preferences-2-tab-1-widget-1-filter-1-title",
          "filter": "test-view-to-update-user-preferences-2-tab-1-widget-1-filter-1-value-updated"
        }
      }
    }
    """
    When I do PUT /api/v4/views/{{ .viewID }}:
    """json
    {
      "title": "test-view-to-update-user-preferences-2-title",
      "description": "test-view-to-update-user-preferences-2-description",
      "enabled": true,
      "group": "test-viewgroup-to-view-edit",
      "tabs": [
        {
          "_id": "test-view-to-update-user-preferences-2-tab-1",
          "title": "test-view-to-update-user-preferences-2-tab-1-title",
          "widgets": [
            {
              "_id": "test-view-to-update-user-preferences-2-tab-1-widget-1",
              "title": "test-view-to-update-user-preferences-2-tab-1-widget-1-title",
              "type": "test-view-to-update-user-preferences-2-tab-1-widget-1-type",
              "parameters": {
                "viewFilters": []
              }
            }
          ]
        }
      ]
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/user-preferences/test-view-to-update-user-preferences-2-tab-1-widget-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "widget": "test-view-to-update-user-preferences-2-tab-1-widget-1",
      "content": {
        "mainFilter": {
          "title": "test-view-to-update-user-preferences-2-tab-1-widget-1-filter-1-title",
          "filter": "test-view-to-update-user-preferences-2-tab-1-widget-1-filter-1-value-updated"
        }
      }
    }
    """

  Scenario: given updated widget filters should update user preferences with array main filter
    When I am test-role-to-user-preferences-update-2
    When I do POST /api/v4/views:
    """json
    {
      "title": "test-view-to-update-user-preferences-3-title",
      "description": "test-view-to-update-user-preferences-3-description",
      "enabled": true,
      "group": "test-viewgroup-to-view-edit",
      "tabs": [
        {
          "_id": "test-view-to-update-user-preferences-3-tab-1",
          "title": "test-view-to-update-user-preferences-3-tab-1-title",
          "widgets": [
            {
              "_id": "test-view-to-update-user-preferences-3-tab-1-widget-1",
              "title": "test-view-to-update-user-preferences-3-tab-1-widget-1-title",
              "type": "test-view-to-update-user-preferences-3-tab-1-widget-1-type",
              "parameters": {
                "viewFilters": [
                  {
                    "title": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-1-title",
                    "filter": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-1-value"
                  },
                  {
                    "title": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-2-title",
                    "filter": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-2-value"
                  },
                  {
                    "title": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-3-title",
                    "filter": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-3-value"
                  },
                  {
                    "title": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-4-title",
                    "filter": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-4-value"
                  }
                ]
              }
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
      "widget": "test-view-to-update-user-preferences-3-tab-1-widget-1",
      "content": {
        "mainFilter": [
          {
            "title": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-1-title",
            "filter": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-1-value"
          },
          {
            "title": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-2-title",
            "filter": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-2-value"
          },
          {
            "title": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-3-title",
            "filter": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-3-value"
          }
        ]
      }
    }
    """
    Then the response code should be 200
    When I do PUT /api/v4/views/{{ .viewID }}:
    """json
    {
      "title": "test-view-to-update-user-preferences-3-title",
      "description": "test-view-to-update-user-preferences-3-description",
      "enabled": true,
      "group": "test-viewgroup-to-view-edit",
      "tabs": [
        {
          "_id": "test-view-to-update-user-preferences-3-tab-1",
          "title": "test-view-to-update-user-preferences-3-tab-1-title",
          "widgets": [
            {
              "_id": "test-view-to-update-user-preferences-3-tab-1-widget-1",
              "title": "test-view-to-update-user-preferences-3-tab-1-widget-1-title",
              "type": "test-view-to-update-user-preferences-3-tab-1-widget-1-type",
              "parameters": {
                "viewFilters": [
                  {
                    "title": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-1-title",
                    "filter": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-1-value-updated"
                  },
                  {
                    "title": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-2-title-updated",
                    "filter": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-2-value-updated"
                  },
                  {
                    "title": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-3-title",
                    "filter": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-3-value"
                  },
                  {
                    "title": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-4-title",
                    "filter": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-4-value"
                  }
                ]
              }
            }
          ]
        }
      ]
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/user-preferences/test-view-to-update-user-preferences-3-tab-1-widget-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "widget": "test-view-to-update-user-preferences-3-tab-1-widget-1",
      "content": {
        "mainFilter": [
          {
            "title": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-1-title",
            "filter": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-1-value-updated"
          },
          {
            "title": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-2-title",
            "filter": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-2-value"
          },
          {
            "title": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-3-title",
            "filter": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-3-value"
          }
        ]
      }
    }
    """
    When I do PUT /api/v4/views/{{ .viewID }}:
    """json
    {
      "title": "test-view-to-update-user-preferences-3-title",
      "description": "test-view-to-update-user-preferences-3-description",
      "enabled": true,
      "group": "test-viewgroup-to-view-edit",
      "tabs": [
        {
          "_id": "test-view-to-update-user-preferences-3-tab-1",
          "title": "test-view-to-update-user-preferences-3-tab-1-title",
          "widgets": [
            {
              "_id": "test-view-to-update-user-preferences-3-tab-1-widget-1",
              "title": "test-view-to-update-user-preferences-3-tab-1-widget-1-title",
              "type": "test-view-to-update-user-preferences-3-tab-1-widget-1-type",
              "parameters": {
                "viewFilters": [
                  {
                    "title": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-1-title",
                    "filter": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-1-value-updated"
                  },
                  {
                    "title": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-2-title-updated",
                    "filter": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-2-value-updated"
                  },
                  {
                    "title": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-3-title",
                    "filter": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-3-value"
                  }
                ]
              }
            }
          ]
        }
      ]
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/user-preferences/test-view-to-update-user-preferences-3-tab-1-widget-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "widget": "test-view-to-update-user-preferences-3-tab-1-widget-1",
      "content": {
        "mainFilter": [
          {
            "title": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-1-title",
            "filter": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-1-value-updated"
          },
          {
            "title": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-2-title",
            "filter": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-2-value"
          },
          {
            "title": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-3-title",
            "filter": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-3-value"
          }
        ]
      }
    }
    """
    When I do PUT /api/v4/views/{{ .viewID }}:
    """json
    {
      "title": "test-view-to-update-user-preferences-3-title",
      "description": "test-view-to-update-user-preferences-3-description",
      "enabled": true,
      "group": "test-viewgroup-to-view-edit",
      "tabs": [
        {
          "_id": "test-view-to-update-user-preferences-3-tab-1",
          "title": "test-view-to-update-user-preferences-3-tab-1-title",
          "widgets": [
            {
              "_id": "test-view-to-update-user-preferences-3-tab-1-widget-1",
              "title": "test-view-to-update-user-preferences-3-tab-1-widget-1-title",
              "type": "test-view-to-update-user-preferences-3-tab-1-widget-1-type",
              "parameters": {
                "viewFilters": []
              }
            }
          ]
        }
      ]
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/user-preferences/test-view-to-update-user-preferences-3-tab-1-widget-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "widget": "test-view-to-update-user-preferences-3-tab-1-widget-1",
      "content": {
        "mainFilter": [
          {
            "title": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-1-title",
            "filter": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-1-value-updated"
          },
          {
            "title": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-2-title",
            "filter": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-2-value"
          },
          {
            "title": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-3-title",
            "filter": "test-view-to-update-user-preferences-3-tab-1-widget-1-filter-3-value"
          }
        ]
      }
    }
    """

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
      },
      "filters": [
        {
          "_id": "test-widgetfilter-to-user-preferences-get-1",
          "title": "test-widgetfilter-to-user-preferences-get-1-title",
          "is_private": true,
          "author": {
            "_id": "test-user-to-user-preferences-edit",
            "name": "test-user-to-user-preferences-edit"
          },
          "created": 1611229670,
          "updated": 1611229670,
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-widgetfilter-to-user-preferences-get-1-pattern"
                }
              }
            ]
          ]
        },
        {
          "_id": "test-widgetfilter-to-user-preferences-get-2",
          "title": "test-widgetfilter-to-user-preferences-get-2-title",
          "is_private": true,
          "author": {
            "_id": "test-user-to-user-preferences-edit",
            "name": "test-user-to-user-preferences-edit"
          },
          "created": 1611229670,
          "updated": 1611229670,
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-widgetfilter-to-user-preferences-get-2-pattern"
                }
              }
            ]
          ]
        }
      ]
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
      },
      "filters": []
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

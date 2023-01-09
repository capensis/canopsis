Feature: Copy a view tab
  I need to be able to copy a view tab
  Only admin should be able to copy a view tab

  Scenario: given copy request should return ok
    When I am admin
    When I do POST /api/v4/view-tab-copy/test-tab-to-copy-1:
    """json
    {
      "title": "test-tab-to-copy-1-title-updated",
      "view": "test-view-to-tab-copy-2"
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "title": "test-tab-to-copy-1-title-updated",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "widgets": [
        {
          "title": "test-widget-to-tab-copy-1-title",
          "type": "test-widget-to-tab-copy-1-type",
          "grid_parameters": {
            "desktop": {"x": 0, "y": 0}
          },
          "parameters": {
            "test-widget-to-tab-copy-1-parameter-1": {
              "test-widget-to-tab-copy-1-parameter-1-subparameter": "test-widget-to-tab-copy-1-parameter-1-subvalue"
            },
            "test-widget-to-tab-copy-1-parameter-2": [
              {
                "test-widget-to-tab-copy-1-parameter-2-subparameter": "test-widget-to-tab-copy-1-parameter-2-subvalue"
              }
            ]
          },
          "filters": [
            {
              "title": "test-widgetfilter-to-tab-copy-1-title",
              "is_private": false,
              "author": {
                "_id": "root",
                "name": "root"
              }
            },
            {
              "title": "test-widgetfilter-to-tab-copy-2-title",
              "is_private": false,
              "author": {
                "_id": "root",
                "name": "root"
              }
            }
          ],
          "author": {
            "_id": "root",
            "name": "root"
          }
        }
      ]
    }
    """
    When I do GET /api/v4/view-tabs/{{ .lastResponse._id}}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "title": "test-tab-to-copy-1-title-updated",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "widgets": [
        {
          "title": "test-widget-to-tab-copy-1-title",
          "type": "test-widget-to-tab-copy-1-type",
          "grid_parameters": {
            "desktop": {"x": 0, "y": 0}
          },
          "parameters": {
            "test-widget-to-tab-copy-1-parameter-1": {
              "test-widget-to-tab-copy-1-parameter-1-subparameter": "test-widget-to-tab-copy-1-parameter-1-subvalue"
            },
            "test-widget-to-tab-copy-1-parameter-2": [
              {
                "test-widget-to-tab-copy-1-parameter-2-subparameter": "test-widget-to-tab-copy-1-parameter-2-subvalue"
              }
            ]
          },
          "filters": [
            {
              "title": "test-widgetfilter-to-tab-copy-1-title",
              "author": {
                "_id": "root",
                "name": "root"
              },
              "is_private": false,
              "alarm_pattern": [
                [
                  {
                    "field": "v.component",
                    "cond": {
                      "type": "eq",
                      "value": "test-widgetfilter-to-tab-copy-1-pattern"
                    }
                  }
                ]
              ]
            },
            {
              "title": "test-widgetfilter-to-tab-copy-2-title",
              "author": {
                "_id": "root",
                "name": "root"
              },
              "is_private": false,
              "alarm_pattern": [
                [
                  {
                    "field": "v.component",
                    "cond": {
                      "type": "eq",
                      "value": "test-widgetfilter-to-tab-copy-2-pattern"
                    }
                  }
                ]
              ]
            }
          ],
          "author": {
            "_id": "root",
            "name": "root"
          }
        }
      ]
    }
    """
    Then the response key "_id" should not be "test-tab-to-copy-1"
    Then the response key "widgets.0._id" should not be "test-widget-to-tab-copy-1"
    Then the response key "widgets.0.filters.0._id" should not be "test-widgetfilter-to-tab-copy-1"
    Then the response key "widgets.0.filters.1._id" should not be "test-widgetfilter-to-tab-copy-2"
    Then the response key "widgets.0.parameters.mainFilter" should not be "test-widgetfilter-to-tab-copy-1"
    When I do GET /api/v4/views/test-view-to-tab-copy-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-view-to-tab-copy-2",
      "tabs": [
        {
          "_id": "test-tab-to-copy-2"
        },
        {
          "title": "test-tab-to-copy-1-title-updated"
        }
      ]
    }
    """

  Scenario: given copy request and no auth user should not allow access
    When I do POST /api/v4/view-tab-copy/test-tab-to-copy-1
    Then the response code should be 401

  Scenario: given copy request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/view-tab-copy/test-tab-to-copy-1
    Then the response code should be 403

  Scenario: given invalid copy request should return errors
    When I am admin
    When I do POST /api/v4/view-tab-copy/test-tab-to-copy-1:
    """json
    {
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "view": "View is missing.",
        "title": "Title is missing."
      }
    }
    """

  Scenario: given copy request and auth user without view permission should not allow access
    When I am admin
    When I do POST /api/v4/view-tab-copy/test-tab-to-copy-1:
    """json
    {
      "view": "test-view-to-tab-check-access",
      "title": "test-tab-to-copy-1-title-updated"
    }
    """
    Then the response code should be 403
    When I do POST /api/v4/view-tab-copy/test-tab-to-copy-3:
    """json
    {
      "view": "test-view-to-tab-copy-2",
      "title": "test-tab-to-copy-1-title-updated"
    }
    """
    Then the response code should be 403

  Scenario: given copy request with not exist tab should return not allow access error
    When I am admin
    When I do POST /api/v4/view-tab-copy/test-tab-to-copy-1:
    """json
    {
      "view": "test-view-not-found",
      "title": "test-tab-to-copy-1-title-updated"
    }
    """
    Then the response code should be 403

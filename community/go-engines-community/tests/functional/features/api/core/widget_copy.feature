Feature: Copy a widget
  I need to be able to copy a widget
  Only admin should be able to copy a widget

  Scenario: given copy request should return ok
    When I am admin
    When I do POST /api/v4/widget-copy/test-widget-to-copy-1:
    """json
    {
      "tab": "test-tab-to-widget-copy-2",
      "title": "test-widget-to-copy-1-title-updated",
      "type": "test-widget-to-copy-1-type-updated",
      "grid_parameters": {
        "desktop": {"x": 0, "y": 1}
      },
      "parameters": {
        "mainFilter": "test-widgetfilter-to-widget-copy-1",
        "test-widget-to-copy-1-parameter-1": {
          "test-widget-to-copy-1-parameter-1-subparameter": "test-widget-to-copy-1-parameter-1-subvalue"
        },
        "test-widget-to-copy-1-parameter-2": [
          {
            "test-widget-to-copy-1-parameter-2-subparameter": "test-widget-to-copy-1-parameter-2-subvalue"
          }
        ]
      }
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "title": "test-widget-to-copy-1-title-updated",
      "type": "test-widget-to-copy-1-type-updated",
      "grid_parameters": {
        "desktop": {"x": 0, "y": 1}
      },
      "parameters": {
        "test-widget-to-copy-1-parameter-1": {
          "test-widget-to-copy-1-parameter-1-subparameter": "test-widget-to-copy-1-parameter-1-subvalue"
        },
        "test-widget-to-copy-1-parameter-2": [
          {
            "test-widget-to-copy-1-parameter-2-subparameter": "test-widget-to-copy-1-parameter-2-subvalue"
          }
        ]
      },
      "filters": [
        {
          "title": "test-widgetfilter-to-widget-copy-1-title",
          "is_private": false,
          "author": {
            "_id": "root",
            "name": "root"
          },
          "corporate_alarm_pattern": "test-pattern-to-filter-edit-1",
          "corporate_alarm_pattern_title": "test-pattern-to-filter-edit-1-title",
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-pattern-to-filter-edit-1-pattern"
                }
              }
            ]
          ]
        },
        {
          "title": "test-widgetfilter-to-widget-copy-2-title",
          "is_private": false,
          "author": {
            "_id": "root",
            "name": "root"
          },
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-widgetfilter-to-widget-copy-2-pattern"
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
    """
    When I do GET /api/v4/widgets/{{ .lastResponse._id}}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "title": "test-widget-to-copy-1-title-updated",
      "type": "test-widget-to-copy-1-type-updated",
      "grid_parameters": {
        "desktop": {"x": 0, "y": 1}
      },
      "parameters": {
        "test-widget-to-copy-1-parameter-1": {
          "test-widget-to-copy-1-parameter-1-subparameter": "test-widget-to-copy-1-parameter-1-subvalue"
        },
        "test-widget-to-copy-1-parameter-2": [
          {
            "test-widget-to-copy-1-parameter-2-subparameter": "test-widget-to-copy-1-parameter-2-subvalue"
          }
        ]
      },
      "filters": [
        {
          "title": "test-widgetfilter-to-widget-copy-1-title",
          "is_private": false,
          "author": {
            "_id": "root",
            "name": "root"
          },
          "corporate_alarm_pattern": "test-pattern-to-filter-edit-1",
          "corporate_alarm_pattern_title": "test-pattern-to-filter-edit-1-title",
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-pattern-to-filter-edit-1-pattern"
                }
              }
            ]
          ]
        },
        {
          "title": "test-widgetfilter-to-widget-copy-2-title",
          "is_private": false,
          "author": {
            "_id": "root",
            "name": "root"
          },
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-widgetfilter-to-widget-copy-2-pattern"
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
    """
    Then the response key "_id" should not be "test-widget-to-copy-1"
    Then the response key "filters.0._id" should not be "test-widgetfilter-to-widget-copy-1"
    Then the response key "filters.1._id" should not be "test-widgetfilter-to-widget-copy-2"
    Then the response key "parameters.mainFilter" should not be "test-widgetfilter-to-widget-copy-1"
    When I do GET /api/v4/views/test-view-to-widget-copy-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-view-to-widget-copy-2",
      "tabs": [
        {
          "_id": "test-tab-to-widget-copy-2",
          "widgets": [
            {
              "_id": "test-widget-to-copy-2"
            },
            {
              "title": "test-widget-to-copy-1-title-updated"
            }
          ]
        }
      ]
    }
    """

  Scenario: given copy request and no auth user should not allow access
    When I do POST /api/v4/widget-copy/test-widget-to-copy-1
    Then the response code should be 401

  Scenario: given copy request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/widget-copy/test-widget-to-copy-1
    Then the response code should be 403

  Scenario: given invalid copy request should return errors
    When I am admin
    When I do POST /api/v4/widget-copy/test-widget-to-copy-1:
    """json
    {
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "tab": "Tab is missing.",
        "type": "Type is missing."
      }
    }
    """

  Scenario: given copy request and auth user without view permission should not allow access
    When I am admin
    When I do POST /api/v4/widget-copy/test-widget-to-copy-1:
    """json
    {
      "tab": "test-tab-to-widget-check-access",
      "title": "test-widget-to-copy-1-title-updated",
      "type": "test-widget-to-copy-1-type-updated"
    }
    """
    Then the response code should be 403
    When I do POST /api/v4/widget-copy/test-widget-to-copy-3:
    """json
    {
      "tab": "test-tab-to-widget-copy-2",
      "title": "test-widget-to-copy-1-title-updated",
      "type": "test-widget-to-copy-1-type-updated"
    }
    """
    Then the response code should be 403

  Scenario: given copy request with not exist tab should return not allow access error
    When I am admin
    When I do POST /api/v4/widget-copy/test-widget-to-copy-1:
    """json
    {
      "tab": "test-tab-not-found",
      "title": "test-widget-to-copy-1-title-updated",
      "type": "test-widget-to-copy-1-type-updated"
    }
    """
    Then the response code should be 403

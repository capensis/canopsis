Feature: Copy a view
  I need to be able to copy a view
  Only admin should be able to copy a view

  Scenario: given copy request should return ok
    When I am admin
    When I do POST /api/v4/view-copy/test-view-to-copy-1:
    """json
    {
      "description": "test-view-to-copy-1-description-updated",
      "enabled": true,
      "group": "test-viewgroup-to-view-copy",
      "title": "test-view-to-copy-1-title-updated",
      "periodic_refresh": {
        "enabled": true,
        "value": 600,
        "unit": "m"
      },
      "tags": [
        "test-view-to-copy-1-tag-updated"
      ]
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "description": "test-view-to-copy-1-description-updated",
      "enabled": true,
      "group": {
        "_id": "test-viewgroup-to-view-copy",
        "author": {
          "_id": "root",
          "name": "root"
        },
        "created": 1611229670,
        "title": "test-viewgroup-to-view-copy-title",
        "updated": 1611229670
      },
      "title": "test-view-to-copy-1-title-updated",
      "periodic_refresh": {
        "enabled": true,
        "value": 600,
        "unit": "m"
      },
      "tags": [
        "test-view-to-copy-1-tag-updated"
      ]
    }
    """
    When I do GET /api/v4/views/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "description": "test-view-to-copy-1-description-updated",
      "enabled": true,
      "group": {
        "_id": "test-viewgroup-to-view-copy",
        "author": {
          "_id": "root",
          "name": "root"
        },
        "created": 1611229670,
        "title": "test-viewgroup-to-view-copy-title",
        "updated": 1611229670
      },
      "title": "test-view-to-copy-1-title-updated",
      "periodic_refresh": {
        "enabled": true,
        "value": 600,
        "unit": "m"
      },
      "tags": [
        "test-view-to-copy-1-tag-updated"
      ],
      "tabs": [
        {
          "title": "test-tab-to-view-copy-1-title",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "widgets": [
            {
              "title": "test-widget-to-view-copy-1-title",
              "type": "test-widget-to-view-copy-1-type",
              "grid_parameters": {
                "desktop": {"x": 0, "y": 0}
              },
              "parameters": {
                "test-widget-to-view-copy-1-parameter-1": {
                  "test-widget-to-view-copy-1-parameter-1-subparameter": "test-widget-to-view-copy-1-parameter-1-subvalue"
                },
                "test-widget-to-view-copy-1-parameter-2": [
                  {
                    "test-widget-to-view-copy-1-parameter-2-subparameter": "test-widget-to-view-copy-1-parameter-2-subvalue"
                  }
                ]
              },
              "filters": [
                {
                  "title": "test-widgetfilter-to-view-copy-1-title",
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
                          "value": "test-widgetfilter-to-view-copy-1-pattern"
                        }
                      }
                    ]
                  ]
                },
                {
                  "title": "test-widgetfilter-to-view-copy-2-title",
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
                          "value": "test-widgetfilter-to-view-copy-2-pattern"
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
      ]
    }
    """
    Then the response key "_id" should not be "test-view-to-copy-1"
    Then the response key "tabs.0._id" should not be "test-tab-to-view-copy-1"
    Then the response key "tabs.0.widgets.0._id" should not be "test-widget-to-view-copy-1"
    Then I save response viewId={{ .lastResponse._id }}
    When I do GET /api/v4/permissions?search={{ .viewId }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "{{ .viewId }}",
          "name": "{{ .viewId }}",
          "description": "Rights on view : test-view-to-copy-1-title-updated",
          "view": {
            "_id": "{{ .viewId }}",
            "title": "test-view-to-copy-1-title-updated"
          },
          "view_group": {
            "_id": "test-viewgroup-to-view-copy",
            "title": "test-viewgroup-to-view-copy-title"
          },
          "type": "RW"
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
    When I do GET /api/v4/view-groups?search=test-viewgroup-to-view-copy&with_views=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "views": [
            {
              "_id": "test-view-to-copy-2"
            },
            {
              "title": "test-view-to-copy-1-title-updated"
            }
          ]
        }
      ]
    }
    """

  Scenario: given invalid copy request should return errors
    When I am admin
    When I do POST /api/v4/view-copy/test-view-to-copy-1:
    """json
    {
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "enabled": "Enabled is missing.",
        "group": "Group is missing.",
        "title": "Title is missing."
      }
    }
    """

  Scenario: given copy request and no auth user should not allow access
    When I do POST /api/v4/view-copy/test-view-to-copy-1
    Then the response code should be 401

  Scenario: given copy request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/view-copy/test-view-to-copy-1
    Then the response code should be 403

  Scenario: given copy request and auth user without view permission should not allow access
    When I am admin
    When I do POST /api/v4/view-copy/test-view-to-copy-3
    Then the response code should be 403

  Scenario: given copy request with not exist tab should return not allow access error
    When I am admin
    When I do POST /api/v4/view-copy/test-view-not-found
    Then the response code should be 403

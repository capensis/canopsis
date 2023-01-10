Feature: Get a view tab
  I need to be able to get a view tab
  Only admin should be able to get a view tab

  Scenario: given get request should return tab
    When I am admin
    When I do GET /api/v4/view-tabs/test-tab-to-get
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-tab-to-get",
      "title": "test-tab-to-get-title",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "widgets": [
        {
          "_id": "test-widget-to-tab-get-1",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "created": 1611229670,
          "grid_parameters": {
            "desktop": {"x": 0,"y": 0}
          },
          "parameters": {
            "test-widget-to-tab-get-1-parameter-1": {
              "test-widget-to-tab-get-1-parameter-1-subparameter": "test-widget-to-tab-get-1-parameter-1-subvalue"
            },
            "test-widget-to-tab-get-1-parameter-2": [
              {
                "test-widget-to-tab-get-1-parameter-2-subparameter": "test-widget-to-tab-get-1-parameter-2-subvalue"
              }
            ]
          },
          "filters": [
            {
              "_id": "test-widgetfilter-to-tab-get-1",
              "title": "test-widgetfilter-to-tab-get-1-title",
              "is_private": false,
              "author": {
                "_id": "nopermsuser",
                "name": "nopermsuser"
              },
              "created": 1611229670,
              "updated": 1611229670,
              "alarm_pattern": [
                [
                  {
                    "field": "v.component",
                    "cond": {
                      "type": "eq",
                      "value": "test-widgetfilter-to-tab-get-1-pattern"
                    }
                  }
                ]
              ]
            }
          ],
          "title": "test-widget-to-tab-get-1-title",
          "type": "test-widget-to-tab-get-1-type",
          "updated": 1611229670
        },
        {
          "_id": "test-widget-to-tab-get-2",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "created": 1611229670,
          "grid_parameters": {
            "desktop": {"x": 0,"y": 1}
          },
          "parameters": {
            "test-widget-to-tab-get-2-parameter-1": {
              "test-widget-to-tab-get-2-parameter-1-subparameter": "test-widget-to-tab-get-2-parameter-1-subvalue"
            },
            "test-widget-to-tab-get-2-parameter-2": [
              {
                "test-widget-to-tab-get-2-parameter-2-subparameter": "test-widget-to-tab-get-2-parameter-2-subvalue"
              }
            ]
          },
          "filters": [
            {
              "_id": "test-widgetfilter-to-tab-get-2",
              "title": "test-widgetfilter-to-tab-get-2-title",
              "is_private": false,
              "author": {
                "_id": "root",
                "name": "root"
              },
              "created": 1611229670,
              "updated": 1611229670,
              "alarm_pattern": [
                [
                  {
                    "field": "v.component",
                    "cond": {
                      "type": "eq",
                      "value": "test-widgetfilter-to-tab-get-2-pattern"
                    }
                  }
                ]
              ]
            }
          ],
          "title": "test-widget-to-tab-get-2-title",
          "type": "test-widget-to-tab-get-2-type",
          "updated": 1611229670
        }
      ],
      "created": 1611229670,
      "updated": 1611229670
    }
    """

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/view-tabs/test-tab-to-get
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/view-tabs/test-tab-to-get
    Then the response code should be 403

  Scenario: given get request and auth user without view permissions should not allow access
    When I am admin
    When I do GET /api/v4/view-tabs/test-tab-to-check-access
    Then the response code should be 403

  Scenario: given get request with not exist id should return error
    When I am admin
    When I do GET /api/v4/view-tabs/test-tab-not-found
    Then the response code should be 404

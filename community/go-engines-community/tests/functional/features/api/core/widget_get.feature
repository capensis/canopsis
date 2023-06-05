Feature: Get a widget
  I need to be able to get a widget
  Only admin should be able to get a widget

  Scenario: given get request should return widget
    When I am admin
    When I do GET /api/v4/widgets/test-widget-to-get
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-widget-to-get",
      "title": "test-widget-to-get-title",
      "type": "test-widget-to-get-type",
      "grid_parameters": {
        "desktop": {"x": 0, "y": 0}
      },
      "parameters": {
        "test-widget-to-get-parameter-1": {
          "test-widget-to-get-parameter-1-subparameter": "test-widget-to-get-parameter-1-subvalue"
        },
        "test-widget-to-get-parameter-2": [
          {
            "test-widget-to-get-parameter-2-subparameter": "test-widget-to-get-parameter-2-subvalue"
          }
        ]
      },
      "filters": [
        {
          "_id": "test-widgetfilter-to-widget-get-1",
          "title": "test-widgetfilter-to-widget-get-1-title",
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
                  "value": "test-widgetfilter-to-widget-get-1-pattern"
                }
              }
            ]
          ]
        },
        {
          "_id": "test-widgetfilter-to-widget-get-2",
          "title": "test-widgetfilter-to-widget-get-2-title",
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
                  "value": "test-widgetfilter-to-widget-get-2-pattern"
                }
              }
            ]
          ]
        }
      ],
      "author": {
        "_id": "root",
        "name": "root"
      },
      "created": 1611229670,
      "updated": 1611229670
    }
    """

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/widgets/test-widget-to-get
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/widgets/test-widget-to-get
    Then the response code should be 403

  Scenario: given get request and auth user without view permissions should not allow access
    When I am admin
    When I do GET /api/v4/widgets/test-widget-to-check-access
    Then the response code should be 403

  Scenario: given get request with not exist id should return error
    When I am admin
    When I do GET /api/v4/widgets/test-widget-not-found
    Then the response code should be 404

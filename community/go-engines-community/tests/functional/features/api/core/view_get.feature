Feature: Get a view
  I need to be able to get a view
  Only admin should be able to get a view

  Scenario: given search request should return views
    When I am admin
    When I do GET /api/v4/views?search=test-view-to-get
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "test-view-to-get-1",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "created": 1611229670,
          "description": "test-view-to-get-1-description",
          "enabled": true,
          "group": {
            "_id": "test-viewgroup-to-view-edit",
            "author": {
              "_id": "root",
              "name": "root"
            },
            "created": 1611229670,
            "title": "test-viewgroup-to-view-edit-title",
            "updated": 1611229670
          },
          "title": "test-view-to-get-1-title",
          "periodic_refresh": {
            "enabled": true,
            "value": 1,
            "unit": "s"
          },
          "tags": [
            "test-view-to-get-1-tag"
          ],
          "updated": 1611229670
        },
        {
          "_id": "test-view-to-get-2",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "created": 1611229670,
          "description": "test-view-to-get-2-description",
          "enabled": true,
          "group": {
            "_id": "test-viewgroup-to-view-edit",
            "author": {
              "_id": "root",
              "name": "root"
            },
            "created": 1611229670,
            "title": "test-viewgroup-to-view-edit-title",
            "updated": 1611229670
          },
          "title": "test-view-to-get-2-title",
          "periodic_refresh": {
            "enabled": true,
            "value": 1,
            "unit": "s"
          },
          "tags": [
            "test-view-to-get-2-tag"
          ],
          "updated": 1611229670
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 2
      }
    }
    """

  Scenario: given get request should return view
    When I am admin
    When I do GET /api/v4/views/test-view-to-get-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-view-to-get-1",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "created": 1611229670,
      "description": "test-view-to-get-1-description",
      "enabled": true,
      "group": {
        "_id": "test-viewgroup-to-view-edit",
        "author": {
          "_id": "root",
          "name": "root"
        },
        "created": 1611229670,
        "title": "test-viewgroup-to-view-edit-title",
        "updated": 1611229670
      },
      "title": "test-view-to-get-1-title",
      "periodic_refresh": {
        "enabled": true,
        "value": 1,
        "unit": "s"
      },
      "tabs": [
        {
          "_id": "test-tab-to-view-get-1",
          "title": "test-tab-to-view-get-1-title",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "created": 1611229670,
          "updated": 1611229670,
          "widgets": [
            {
              "_id": "test-widget-to-view-get-1",
              "author": {
                "_id": "root",
                "name": "root"
              },
              "created": 1611229670,
              "updated": 1611229670,
              "grid_parameters": {
                "desktop": {"x": 0, "y": 0}
              },
              "parameters": {
                "test-widget-to-view-get-1-parameter-1": {
                  "test-widget-to-view-get-1-parameter-1-subparameter": "test-widget-to-view-get-1-parameter-1-subvalue"
                },
                "test-widget-to-view-get-1-parameter-2": [
                  {
                    "test-widget-to-view-get-1-parameter-2-subparameter": "test-widget-to-view-get-1-parameter-2-subvalue"
                  }
                ]
              },
              "filters": [
                {
                  "_id": "test-widgetfilter-to-view-get-1",
                  "title": "test-widgetfilter-to-view-get-1-title",
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
                          "value": "test-widgetfilter-to-view-get-1-pattern"
                        }
                      }
                    ]
                  ]
                },
                {
                  "_id": "test-widgetfilter-to-view-get-2",
                  "title": "test-widgetfilter-to-view-get-2-title",
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
                          "value": "test-widgetfilter-to-view-get-2-pattern"
                        }
                      }
                    ]
                  ]
                }
              ],
              "title": "test-widget-to-view-get-1-title",
              "type": "test-widget-to-view-get-1-type"
            },
            {
              "_id": "test-widget-to-view-get-2",
              "author": {
                "_id": "root",
                "name": "root"
              },
              "created": 1611229670,
              "updated": 1611229670,
              "grid_parameters": {
                "desktop": {"x": 0, "y": 1}
              },
              "parameters": {
                "test-widget-to-view-get-2-parameter-1": {
                  "test-widget-to-view-get-2-parameter-1-subparameter": "test-widget-to-view-get-2-parameter-1-subvalue"
                },
                "test-widget-to-view-get-2-parameter-2": [
                  {
                    "test-widget-to-view-get-2-parameter-2-subparameter": "test-widget-to-view-get-2-parameter-2-subvalue"
                  }
                ]
              },
              "filters": [],
              "title": "test-widget-to-view-get-2-title",
              "type": "test-widget-to-view-get-2-type"
            }
          ]
        },
        {
          "_id": "test-tab-to-view-get-2",
          "title": "test-tab-to-view-get-2-title",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "widgets": [],
          "created": 1611229670,
          "updated": 1611229670
        }
      ],
      "tags": [
        "test-view-to-get-1-tag"
      ],
      "updated": 1611229670
    }
    """
    When I do GET /api/v4/views/test-view-to-get-2
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-view-to-get-2",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "created": 1611229670,
      "description": "test-view-to-get-2-description",
      "enabled": true,
      "group": {
        "_id": "test-viewgroup-to-view-edit",
        "author": {
          "_id": "root",
          "name": "root"
        },
        "created": 1611229670,
        "title": "test-viewgroup-to-view-edit-title",
        "updated": 1611229670
      },
      "title": "test-view-to-get-2-title",
      "periodic_refresh": {
        "enabled": true,
        "value": 1,
        "unit": "s"
      },
      "tags": [
        "test-view-to-get-2-tag"
      ],
      "tabs": [],
      "updated": 1611229670
    }
    """

  Scenario: given search request should not return views without access
    When I am admin
    When I do GET /api/v4/views?search=test-view-to-check-access
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """

  Scenario: given get all request and no auth user should not allow access
    When I do GET /api/v4/views
    Then the response code should be 401

  Scenario: given get all request and auth user without view permission should not allow access
    When I am noperms
    When I do GET /api/v4/views
    Then the response code should be 403

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/views/test-view-to-get-1
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/views/test-view-to-get-1
    Then the response code should be 403

  Scenario: given get request and auth user without permissions should not allow access
    When I am admin
    When I do GET /api/v4/views/test-view-to-check-access
    Then the response code should be 403

  Scenario: given get request with not exist id should return not allow access error
    When I am admin
    When I do GET /api/v4/views/test-view-not-found
    Then the response code should be 403

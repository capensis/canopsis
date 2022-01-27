Feature: Get a view
  I need to be able to get a view
  Only admin should be able to get a view

  Scenario: given search request should return views
    When I am admin
    When I do GET /api/v4/views?search=test-view-to-get
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "data": [
        {
          "_id": "test-view-to-get-1",
          "author": "test-view-to-get-1-author",
          "created": 1611229670,
          "description": "test-view-to-get-1-description",
          "enabled": true,
          "group": {
            "_id": "test-viewgroup-to-view-edit",
            "author": "test-viewgroup-to-view-edit-author",
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
              "_id": "test-view-to-get-1-tab-1",
              "title": "test-view-to-get-1-tab-1-title",
              "widgets": [
                {
                  "_id": "test-view-to-get-1-tab-1-widget-1",
                  "grid_parameters": {
                    "test-view-to-get-1-tab-1-widget-1-gridparameter": "test-view-to-get-1-tab-1-widget-1-gridparameter-value"
                  },
                  "parameters": {
                    "test-view-to-get-1-tab-1-widget-1-parameter-1": {
                      "test-view-to-get-1-tab-1-widget-1-parameter-1-subparameter": "test-view-to-get-1-tab-1-widget-1-parameter-1-subvalue"
                    },
                    "test-view-to-get-1-tab-1-widget-1-parameter-2": [
                      {
                        "test-view-to-get-1-tab-1-widget-1-parameter-2-subparameter": "test-view-to-get-1-tab-1-widget-1-parameter-2-subvalue"
                      }
                    ]
                  },
                  "title": "test-view-to-get-1-tab-1-widget-1-title",
                  "type": "test-view-to-get-1-tab-1-widget-1-type"
                }
              ]
            }
          ],
          "tags": [
            "test-view-to-get-1-tag"
          ],
          "updated": 1611229670
        },
        {
          "_id": "test-view-to-get-2",
          "author": "test-view-to-get-2-author",
          "created": 1611229670,
          "description": "test-view-to-get-2-description",
          "enabled": true,
          "group": {
            "_id": "test-viewgroup-to-view-edit",
            "author": "test-viewgroup-to-view-edit-author",
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
          "tabs": [
            {
              "_id": "test-view-to-get-2-tab-1",
              "title": "test-view-to-get-2-tab-1-title",
              "widgets": [
                {
                  "_id": "test-view-to-get-2-tab-1-widget-1",
                  "grid_parameters": {
                    "test-view-to-get-2-tab-1-widget-1-gridparameter": "test-view-to-get-2-tab-1-widget-1-gridparameter-value"
                  },
                  "parameters": {
                    "test-view-to-get-2-tab-1-widget-1-parameter": "test-view-to-get-2-tab-1-widget-1-parameter-value"
                  },
                  "title": "test-view-to-get-2-tab-1-widget-1-title",
                  "type": "test-view-to-get-2-tab-1-widget-1-type"
                }
              ]
            }
          ],
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
    """
    {
      "_id": "test-view-to-get-1",
      "author": "test-view-to-get-1-author",
      "created": 1611229670,
      "description": "test-view-to-get-1-description",
      "enabled": true,
      "group": {
        "_id": "test-viewgroup-to-view-edit",
        "author": "test-viewgroup-to-view-edit-author",
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
          "_id": "test-view-to-get-1-tab-1",
          "title": "test-view-to-get-1-tab-1-title",
          "widgets": [
            {
              "_id": "test-view-to-get-1-tab-1-widget-1",
              "grid_parameters": {
                "test-view-to-get-1-tab-1-widget-1-gridparameter": "test-view-to-get-1-tab-1-widget-1-gridparameter-value"
              },
              "parameters": {
                "test-view-to-get-1-tab-1-widget-1-parameter-1": {
                  "test-view-to-get-1-tab-1-widget-1-parameter-1-subparameter": "test-view-to-get-1-tab-1-widget-1-parameter-1-subvalue"
                },
                "test-view-to-get-1-tab-1-widget-1-parameter-2": [
                  {
                    "test-view-to-get-1-tab-1-widget-1-parameter-2-subparameter": "test-view-to-get-1-tab-1-widget-1-parameter-2-subvalue"
                  }
                ]
              },
              "title": "test-view-to-get-1-tab-1-widget-1-title",
              "type": "test-view-to-get-1-tab-1-widget-1-type"
            }
          ]
        }
      ],
      "tags": [
        "test-view-to-get-1-tag"
      ],
      "updated": 1611229670
    }
    """

  Scenario: given search request should not return views without access
    When I am admin
    When I do GET /api/v4/views?search=test-view-to-check-access
    Then the response code should be 200
    Then the response body should be:
    """
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

  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/views/test-view-to-get-1
    Then the response code should be 403

  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am admin
    When I do GET /api/v4/views/test-view-to-check-access
    Then the response code should be 403

  Scenario: given get request with not exist id should return not allow access error
    When I am admin
    When I do GET /api/v4/views/test-view-not-found
    Then the response code should be 403

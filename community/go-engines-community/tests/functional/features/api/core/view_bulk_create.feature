Feature: Bulk create views
  I need to be able to create multiple view
  Only admin should be able to create multiple view

  Scenario: given bulk create request should return ok
    When I am admin
    When I do POST /api/v4/bulk/views:
    """
    [
      {
        "enabled": true,
        "title": "test-view-to-bulk-create-1-title",
        "description": "test-view-to-bulk-create-1-description",
        "group": "test-viewgroup-to-view-edit",
        "tags": ["test-view-to-bulk-create-1-tag"],
        "periodic_refresh": {
          "enabled": true,
          "value": 10,
          "unit": "m"
        },
        "tabs": [
          {
            "_id": "test-view-to-bulk-create-1-tab-1",
            "title": "test-view-to-bulk-create-1-tab-1-title",
            "widgets": [
              {
                "_id": "test-view-to-bulk-create-1-tab-1-widget-1",
                "grid_parameters": {
                  "test-view-to-bulk-create-1-tab-1-widget-1-gridparameter": "test-view-to-bulk-create-1-tab-1-widget-1-gridparameter-value"
                },
                "parameters": {
                  "test-view-to-bulk-create-1-tab-1-widget-1-parameter": "test-view-to-bulk-create-1-tab-1-widget-1-parameter-value"
                },
                "title": "test-view-to-bulk-create-1-tab-1-widget-1-title",
                "type": "test-view-to-bulk-create-1-tab-1-widget-1-type"
              }
            ]
          }
        ]
      },
      {
        "enabled": true,
        "title": "test-view-to-bulk-create-2-title",
        "description": "test-view-to-bulk-create-2-description",
        "group": "test-viewgroup-to-view-edit",
        "tags": ["test-view-to-bulk-create-2-tag"],
        "periodic_refresh": {
          "enabled": true,
          "value": 10,
          "unit": "m"
        },
        "tabs": [
          {
            "_id": "test-view-to-bulk-create-2-tab-1",
            "title": "test-view-to-bulk-create-2-tab-1-title",
            "widgets": [
              {
                "_id": "test-view-to-bulk-create-2-tab-1-widget-1",
                "grid_parameters": {
                  "test-view-to-bulk-create-2-tab-1-widget-1-gridparameter": "test-view-to-bulk-create-2-tab-1-widget-1-gridparameter-value"
                },
                "parameters": {
                  "test-view-to-bulk-create-2-tab-1-widget-1-parameter": "test-view-to-bulk-create-2-tab-1-widget-1-parameter-value"
                },
                "title": "test-view-to-bulk-create-2-tab-1-widget-1-title",
                "type": "test-view-to-bulk-create-2-tab-1-widget-1-type"
              }
            ]
          }
        ]
      }
    ]
    """
    Then the response code should be 201
    Then the response body should contain:
    """
    [
      {
        "enabled": true,
        "title": "test-view-to-bulk-create-1-title",
        "description": "test-view-to-bulk-create-1-description",
        "group": {
          "_id": "test-viewgroup-to-view-edit",
          "author": "test-viewgroup-to-view-edit-author",
          "created": 1611229670,
          "title": "test-viewgroup-to-view-edit-title",
          "updated": 1611229670
        },
        "author": "root",
        "tags": ["test-view-to-bulk-create-1-tag"],
        "periodic_refresh": {
          "enabled": true,
          "value": 10,
          "unit": "m"
        },
        "tabs": [
          {
            "_id": "test-view-to-bulk-create-1-tab-1",
            "title": "test-view-to-bulk-create-1-tab-1-title",
            "widgets": [
              {
                "_id": "test-view-to-bulk-create-1-tab-1-widget-1",
                "grid_parameters": {
                  "test-view-to-bulk-create-1-tab-1-widget-1-gridparameter": "test-view-to-bulk-create-1-tab-1-widget-1-gridparameter-value"
                },
                "parameters": {
                  "test-view-to-bulk-create-1-tab-1-widget-1-parameter": "test-view-to-bulk-create-1-tab-1-widget-1-parameter-value"
                },
                "title": "test-view-to-bulk-create-1-tab-1-widget-1-title",
                "type": "test-view-to-bulk-create-1-tab-1-widget-1-type"
              }
            ]
          }
        ]
      },
      {
        "enabled": true,
        "title": "test-view-to-bulk-create-2-title",
        "description": "test-view-to-bulk-create-2-description",
        "group": {
          "_id": "test-viewgroup-to-view-edit",
          "author": "test-viewgroup-to-view-edit-author",
          "created": 1611229670,
          "title": "test-viewgroup-to-view-edit-title",
          "updated": 1611229670
        },
        "author": "root",
        "tags": ["test-view-to-bulk-create-2-tag"],
        "periodic_refresh": {
          "enabled": true,
          "value": 10,
          "unit": "m"
        },
        "tabs": [
          {
            "_id": "test-view-to-bulk-create-2-tab-1",
            "title": "test-view-to-bulk-create-2-tab-1-title",
            "widgets": [
              {
                "_id": "test-view-to-bulk-create-2-tab-1-widget-1",
                "grid_parameters": {
                  "test-view-to-bulk-create-2-tab-1-widget-1-gridparameter": "test-view-to-bulk-create-2-tab-1-widget-1-gridparameter-value"
                },
                "parameters": {
                  "test-view-to-bulk-create-2-tab-1-widget-1-parameter": "test-view-to-bulk-create-2-tab-1-widget-1-parameter-value"
                },
                "title": "test-view-to-bulk-create-2-tab-1-widget-1-title",
                "type": "test-view-to-bulk-create-2-tab-1-widget-1-type"
              }
            ]
          }
        ]
      }
    ]
    """

  Scenario: given bulk create request should return ok to get request
    When I am admin
    When I do POST /api/v4/bulk/views:
    """
    [
      {
        "enabled": true,
        "title": "test-view-to-bulk-create-3-title",
        "description": "test-view-to-bulk-create-3-description",
        "group": "test-viewgroup-to-view-edit",
        "tags": ["test-view-to-bulk-create-3-tag"],
        "periodic_refresh": {
          "enabled": true,
          "value": 10,
          "unit": "m"
        },
        "tabs": [
          {
            "_id": "test-view-to-bulk-create-3-tab-1",
            "title": "test-view-to-bulk-create-3-tab-1-title",
            "widgets": [
              {
                "_id": "test-view-to-bulk-create-3-tab-1-widget-1",
                "grid_parameters": {
                  "test-view-to-bulk-create-3-tab-1-widget-1-gridparameter": "test-view-to-bulk-create-3-tab-1-widget-1-gridparameter-value"
                },
                "parameters": {
                  "test-view-to-bulk-create-3-tab-1-widget-1-parameter": "test-view-to-bulk-create-3-tab-1-widget-1-parameter-value"
                },
                "title": "test-view-to-bulk-create-3-tab-1-widget-1-title",
                "type": "test-view-to-bulk-create-3-tab-1-widget-1-type"
              }
            ]
          }
        ]
      }
    ]
    """
    When I do GET /api/v4/views/{{ (index .lastResponse 0)._id}}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "enabled": true,
      "title": "test-view-to-bulk-create-3-title",
      "description": "test-view-to-bulk-create-3-description",
      "group": {
        "_id": "test-viewgroup-to-view-edit",
        "author": "test-viewgroup-to-view-edit-author",
        "created": 1611229670,
        "title": "test-viewgroup-to-view-edit-title",
        "updated": 1611229670
      },
      "author": "root",
      "tags": ["test-view-to-bulk-create-3-tag"],
      "periodic_refresh": {
        "enabled": true,
        "value": 10,
        "unit": "m"
      },
      "tabs": [
        {
          "_id": "test-view-to-bulk-create-3-tab-1",
          "title": "test-view-to-bulk-create-3-tab-1-title",
          "widgets": [
            {
              "_id": "test-view-to-bulk-create-3-tab-1-widget-1",
              "grid_parameters": {
                "test-view-to-bulk-create-3-tab-1-widget-1-gridparameter": "test-view-to-bulk-create-3-tab-1-widget-1-gridparameter-value"
              },
              "parameters": {
                "test-view-to-bulk-create-3-tab-1-widget-1-parameter": "test-view-to-bulk-create-3-tab-1-widget-1-parameter-value"
              },
              "title": "test-view-to-bulk-create-3-tab-1-widget-1-title",
              "type": "test-view-to-bulk-create-3-tab-1-widget-1-type"
            }
          ]
        }
      ]
    }
    """

  Scenario: given bulk create request should create new permission
    When I am admin
    When I do POST /api/v4/bulk/views:
    """
    [
      {
        "enabled": true,
        "title": "test-view-to-bulk-create-4-title",
        "group": "test-viewgroup-to-view-edit"
      },
      {
        "enabled": true,
        "title": "test-view-to-bulk-create-5-title",
        "group": "test-viewgroup-to-view-edit"
      }
    ]
    """
    Then the response code should be 201
    When I save response view1={{ (index .lastResponse 0)._id}}
    When I save response view2={{ (index .lastResponse 1)._id}}
    When I do GET /api/v4/permissions?search={{ .view1 }}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "description": "Rights on view : test-view-to-bulk-create-4-title",
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
    When I do GET /api/v4/permissions?search={{ .view2 }}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "description": "Rights on view : test-view-to-bulk-create-5-title",
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

  Scenario: given invalid bulk create request should return errors
    When I am admin
    When I do POST /api/v4/bulk/views:
    """
    [
      {}
    ]
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "0.enabled": "Enabled is missing.",
        "0.group": "Group is missing.",
        "0.title": "Title is missing."
      }
    }
    """

  Scenario: given bulk create request with one invalid and one valid data
    should return errors
    When I am admin
    When I do POST /api/v4/bulk/views:
    """
    [
      {
        "title": "test-view-to-bulk-create-6-title",
        "enabled": true,
        "group": "test-viewgroup-to-view-edit"
      },
      {}
    ]
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "1.enabled": "Enabled is missing.",
        "1.group": "Group is missing.",
        "1.title": "Title is missing."
      }
    }
    """

  Scenario: given bulk create request with multiple items with the same title
    should return error
    When I am admin
    When I do POST /api/v4/bulk/views:
    """
    [
      {
        "title": "test-view-to-bulk-create-6-title",
        "enabled": true,
        "group": "test-viewgroup-to-view-edit"
      },
      {
        "title": "test-view-to-bulk-create-6-title",
        "enabled": true,
        "group": "test-viewgroup-to-view-edit"
      },
      {
        "title": "test-view-to-bulk-create-6-title",
        "enabled": true,
        "group": "test-viewgroup-to-view-edit"
      }
    ]
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "1.title": "Title already exists.",
        "2.title": "Title already exists."
      }
    }
    """

  Scenario: given bulk create request and no auth user should not allow access
    When I do POST /api/v4/bulk/views
    Then the response code should be 401

  Scenario: given bulk create request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/bulk/views
    Then the response code should be 403

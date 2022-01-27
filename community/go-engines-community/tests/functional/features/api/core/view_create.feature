Feature: Create a view
  I need to be able to create a view
  Only admin should be able to create a view

  Scenario: given create request should return ok
    When I am admin
    When I do POST /api/v4/views:
    """
    {
      "enabled": true,
      "title": "test-view-to-create-1-title",
      "description": "test-view-to-create-1-description",
      "group": "test-viewgroup-to-view-edit",
      "tags": ["test-view-to-create-1-tag"],
      "periodic_refresh": {
        "enabled": true,
        "value": 600,
        "unit": "m"
      },
      "tabs": [
        {
          "_id": "test-view-to-create-1-tab-1",
          "title": "test-view-to-create-1-tab-1-title",
          "widgets": [
            {
              "_id": "test-view-to-create-1-tab-1-widget-1",
              "grid_parameters": {
                "test-view-to-create-1-tab-1-widget-1-gridparameter": "test-view-to-create-1-tab-1-widget-1-gridparameter-value"
              },
              "parameters": {
                "test-view-to-create-1-tab-1-widget-1-parameter": "test-view-to-create-1-tab-1-widget-1-parameter-value"
              },
              "title": "test-view-to-create-1-tab-1-widget-1-title",
              "type": "test-view-to-create-1-tab-1-widget-1-type"
            }
          ]
        }
      ]
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """
    {
      "author": "root",
      "description": "test-view-to-create-1-description",
      "enabled": true,
      "group": {
        "_id": "test-viewgroup-to-view-edit",
        "author": "test-viewgroup-to-view-edit-author",
        "created": 1611229670,
        "title": "test-viewgroup-to-view-edit-title",
        "updated": 1611229670
      },
      "title": "test-view-to-create-1-title",
      "periodic_refresh": {
        "enabled": true,
        "value": 600,
        "unit": "m"
      },
      "tabs": [
        {
          "_id": "test-view-to-create-1-tab-1",
          "title": "test-view-to-create-1-tab-1-title",
          "widgets": [
            {
              "_id": "test-view-to-create-1-tab-1-widget-1",
              "grid_parameters": {
                "test-view-to-create-1-tab-1-widget-1-gridparameter": "test-view-to-create-1-tab-1-widget-1-gridparameter-value"
              },
              "parameters": {
                "test-view-to-create-1-tab-1-widget-1-parameter": "test-view-to-create-1-tab-1-widget-1-parameter-value"
              },
              "title": "test-view-to-create-1-tab-1-widget-1-title",
              "type": "test-view-to-create-1-tab-1-widget-1-type"
            }
          ]
        }
      ],
      "tags": [
        "test-view-to-create-1-tag"
      ]
    }
    """

  Scenario: given create request should return ok to get request
    When I am admin
    When I do POST /api/v4/views:
    """
    {
      "enabled": true,
      "title": "test-view-to-create-2-title",
      "description": "test-view-to-create-2-description",
      "group": "test-viewgroup-to-view-edit",
      "tags": ["test-view-to-create-2-tag"],
      "periodic_refresh": {
        "enabled": true,
        "value": 600,
        "unit": "m"
      },
      "tabs": [
        {
          "_id": "test-view-to-create-2-tab-1",
          "title": "test-view-to-create-2-tab-1-title",
          "widgets": [
            {
              "_id": "test-view-to-create-2-tab-1-widget-1",
              "grid_parameters": {
                "test-view-to-create-2-tab-1-widget-1-gridparameter": "test-view-to-create-2-tab-1-widget-1-gridparameter-value"
              },
              "parameters": {
                "test-view-to-create-2-tab-1-widget-1-parameter": "test-view-to-create-2-tab-1-widget-1-parameter-value"
              },
              "title": "test-view-to-create-2-tab-1-widget-1-title",
              "type": "test-view-to-create-2-tab-1-widget-1-type"
            }
          ]
        }
      ]
    }
    """
    When I do GET /api/v4/views/{{ .lastResponse._id}}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "author": "root",
      "description": "test-view-to-create-2-description",
      "enabled": true,
      "group": {
        "_id": "test-viewgroup-to-view-edit",
        "author": "test-viewgroup-to-view-edit-author",
        "created": 1611229670,
        "title": "test-viewgroup-to-view-edit-title",
        "updated": 1611229670
      },
      "title": "test-view-to-create-2-title",
      "periodic_refresh": {
        "enabled": true,
        "value": 600,
        "unit": "m"
      },
      "tabs": [
        {
          "_id": "test-view-to-create-2-tab-1",
          "title": "test-view-to-create-2-tab-1-title",
          "widgets": [
            {
              "_id": "test-view-to-create-2-tab-1-widget-1",
              "grid_parameters": {
                "test-view-to-create-2-tab-1-widget-1-gridparameter": "test-view-to-create-2-tab-1-widget-1-gridparameter-value"
              },
              "parameters": {
                "test-view-to-create-2-tab-1-widget-1-parameter": "test-view-to-create-2-tab-1-widget-1-parameter-value"
              },
              "title": "test-view-to-create-2-tab-1-widget-1-title",
              "type": "test-view-to-create-2-tab-1-widget-1-type"
            }
          ]
        }
      ],
      "tags": [
        "test-view-to-create-2-tag"
      ]
    }
    """

  Scenario: given create request should create new permission
    When I am admin
    When I do POST /api/v4/views:
    """
    {
      "enabled": true,
      "title": "test-view-to-create-3-title",
      "description": "test-view-to-create-3-description",
      "group": "test-viewgroup-to-view-edit",
      "tags": ["test-view-to-create-3-tag"],
      "tabs": []
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/permissions?search={{ .lastResponse._id}}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "description": "Rights on view : test-view-to-create-3-title",
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

  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/views
    Then the response code should be 401

  Scenario: given create request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/views
    Then the response code should be 403

  Scenario: given invalid create request should return errors
    When I am admin
    When I do POST /api/v4/views:
    """
    {
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "enabled": "Enabled is missing.",
        "group": "Group is missing.",
        "title": "Title is missing."
      }
    }
    """

  Scenario: given invalid create Junit widget request should return errors
    When I am admin
    When I do POST /api/v4/views:
    """
    {
      "enabled": true,
      "title": "test-view-to-create-5-title",
      "description": "test-view-to-create-5-description",
      "group": "test-viewgroup-to-view-edit",
      "tags": [],
      "tabs": [
        {
          "_id": "test-view-to-create-5-tab-1",
          "title": "test-view-to-create-5-tab-1-title",
          "widgets": [
            {
              "_id": "test-view-to-create-5-tab-1-widget-1",
              "grid_parameters": {},
              "parameters": {},
              "title": "test-view-to-create-5-tab-1-widget-1-title",
              "type": "Any"
            },
            {
              "_id": "test-view-to-create-5-tab-1-widget-2",
              "grid_parameters": {},
              "parameters": {},
              "title": "test-view-to-create-5-tab-1-widget-1-title",
              "type": "Junit"
            }
          ]
        },
        {
          "_id": "test-view-to-create-5-tab-2",
          "title": "test-view-to-create-5-tab-2-title",
          "widgets": [
            {
              "_id": "test-view-to-create-5-tab-2-widget-1",
              "grid_parameters": {},
              "parameters": {
                "directory": "test-dir",
                "is_api": true
              },
              "title": "test-view-to-create-5-tab-2-widget-1-title",
              "type": "Junit"
            }
          ]
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "tabs.0.widgets.1.parameters.directory": "Directory is missing.",
        "tabs.1.widgets.0.parameters.directory": "Directory is not empty."
      }
    }
    """

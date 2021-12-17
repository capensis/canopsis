Feature: Bulk update a views
  I need to be able to update multiple views
  Only admin should be able to update multiple view

  Scenario: given bulk update request should update view
    When I am admin
    Then I do PUT /api/v4/bulk/views:
    """
    [
      {
        "_id": "test-view-to-bulk-update-1",
        "enabled": true,
        "title": "test-view-to-bulk-update-1-title",
        "description": "test-view-to-bulk-update-1-description",
        "group": "test-viewgroup-to-view-edit",
        "tags": ["test-view-to-bulk-update-1-tag"],
        "periodic_refresh": {
          "enabled": true,
          "value": 10,
          "unit": "m"
        },
        "tabs": [
          {
            "_id": "test-view-to-bulk-update-1-tab-1",
            "title": "test-view-to-bulk-update-1-tab-1-title",
            "widgets": [
              {
                "_id": "test-view-to-bulk-update-1-tab-1-widget-1",
                "grid_parameters": {
                  "test-view-to-bulk-update-1-tab-1-widget-1-gridparameter": "test-view-to-bulk-update-1-tab-1-widget-1-gridparameter-value"
                },
                "parameters": {
                  "test-view-to-bulk-update-1-tab-1-widget-1-parameter": "test-view-to-bulk-update-1-tab-1-widget-1-parameter-value"
                },
                "title": "test-view-to-bulk-update-1-tab-1-widget-1-title",
                "type": "test-view-to-bulk-update-1-tab-1-widget-1-type"
              }
            ]
          }
        ]
      },
      {
        "_id": "test-view-to-bulk-update-2",
        "enabled": true,
        "title": "test-view-to-bulk-update-2-title",
        "description": "test-view-to-bulk-update-2-description",
        "group": "test-viewgroup-to-view-edit",
        "tags": ["test-view-to-bulk-update-2-tag"],
        "periodic_refresh": {
          "enabled": true,
          "value": 10,
          "unit": "m"
        },
        "tabs": [
          {
            "_id": "test-view-to-bulk-update-2-tab-1",
            "title": "test-view-to-bulk-update-2-tab-1-title",
            "widgets": [
              {
                "_id": "test-view-to-bulk-update-2-tab-1-widget-1",
                "grid_parameters": {
                  "test-view-to-bulk-update-2-tab-1-widget-1-gridparameter": "test-view-to-bulk-update-2-tab-1-widget-1-gridparameter-value"
                },
                "parameters": {
                  "test-view-to-bulk-update-2-tab-1-widget-1-parameter": "test-view-to-bulk-update-2-tab-1-widget-1-parameter-value"
                },
                "title": "test-view-to-bulk-update-2-tab-1-widget-1-title",
                "type": "test-view-to-bulk-update-2-tab-1-widget-1-type"
              }
            ]
          }
        ]
      }
    ]
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    [
      {
        "_id": "test-view-to-bulk-update-1",
        "enabled": true,
        "created": 1611229670,
        "title": "test-view-to-bulk-update-1-title",
        "description": "test-view-to-bulk-update-1-description",
        "group": {
          "_id": "test-viewgroup-to-view-edit",
          "author": "test-viewgroup-to-view-edit-author",
          "created": 1611229670,
          "title": "test-viewgroup-to-view-edit-title",
          "updated": 1611229670
        },
        "author": "root",
        "tags": ["test-view-to-bulk-update-1-tag"],
        "periodic_refresh": {
          "enabled": true,
          "value": 10,
          "unit": "m"
        },
        "tabs": [
          {
            "_id": "test-view-to-bulk-update-1-tab-1",
            "title": "test-view-to-bulk-update-1-tab-1-title",
            "widgets": [
              {
                "_id": "test-view-to-bulk-update-1-tab-1-widget-1",
                "grid_parameters": {
                  "test-view-to-bulk-update-1-tab-1-widget-1-gridparameter": "test-view-to-bulk-update-1-tab-1-widget-1-gridparameter-value"
                },
                "parameters": {
                  "test-view-to-bulk-update-1-tab-1-widget-1-parameter": "test-view-to-bulk-update-1-tab-1-widget-1-parameter-value"
                },
                "title": "test-view-to-bulk-update-1-tab-1-widget-1-title",
                "type": "test-view-to-bulk-update-1-tab-1-widget-1-type"
              }
            ]
          }
        ]
      },
      {
        "_id": "test-view-to-bulk-update-2",
        "enabled": true,
        "created": 1611229670,
        "title": "test-view-to-bulk-update-2-title",
        "description": "test-view-to-bulk-update-2-description",
        "group": {
          "_id": "test-viewgroup-to-view-edit",
          "author": "test-viewgroup-to-view-edit-author",
          "created": 1611229670,
          "title": "test-viewgroup-to-view-edit-title",
          "updated": 1611229670
        },
        "author": "root",
        "tags": ["test-view-to-bulk-update-2-tag"],
        "periodic_refresh": {
          "enabled": true,
          "value": 10,
          "unit": "m"
        },
        "tabs": [
          {
            "_id": "test-view-to-bulk-update-2-tab-1",
            "title": "test-view-to-bulk-update-2-tab-1-title",
            "widgets": [
              {
                "_id": "test-view-to-bulk-update-2-tab-1-widget-1",
                "grid_parameters": {
                  "test-view-to-bulk-update-2-tab-1-widget-1-gridparameter": "test-view-to-bulk-update-2-tab-1-widget-1-gridparameter-value"
                },
                "parameters": {
                  "test-view-to-bulk-update-2-tab-1-widget-1-parameter": "test-view-to-bulk-update-2-tab-1-widget-1-parameter-value"
                },
                "title": "test-view-to-bulk-update-2-tab-1-widget-1-title",
                "type": "test-view-to-bulk-update-2-tab-1-widget-1-type"
              }
            ]
          }
        ]
      }
    ]
    """

  Scenario: given bulk update request with not exist ids should return error
    When I am admin
    When I do PUT /api/v4/bulk/views:
    """
    [
      {
        "_id": "test-view-not-found-1",
        "title": "test-view-not-found-1-title",
        "enabled": true,
        "group": "test-viewgroup-to-view-edit"
      },
      {
        "_id": "test-view-to-bulk-update-1",
        "title": "test-view-to-bulk-update-1-title",
        "enabled": true,
        "group": "test-viewgroup-to-view-edit"
      }
    ]
    """
    Then the response code should be 403

  Scenario: given invalid bulk update request should return errors
    When I am admin
    When I do PUT /api/v4/bulk/views:
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
        "0._id": "ID is missing.",
        "0.enabled": "Enabled is missing.",
        "0.group": "Group is missing.",
        "0.title": "Title is missing."
      }
    }
    """

  Scenario: given bulk update request with one valid item and one invalid item
    should return errors
    When I am admin
    When I do PUT /api/v4/bulk/views:
    """
    [
      {
        "_id": "test-view-to-bulk-update-1",
        "title": "test-view-to-bulk-update-1-title",
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
        "1._id": "ID is missing.",
        "1.enabled": "Enabled is missing.",
        "1.group": "Group is missing.",
        "1.title": "Title is missing."
      }
    }
    """

  Scenario: given bulk update request with multiple items with the same title
    should return error
    When I am admin
    Then I do PUT /api/v4/bulk/views:
    """
    [
      {
        "_id": "test-view-to-bulk-update-1",
        "title": "test-view-to-bulk-update-1-title-same-updated"
      },
      {
        "_id": "test-view-to-bulk-update-2",
        "title": "test-view-to-bulk-update-1-title-same-updated"
      },
      {
        "_id": "test-view-to-bulk-update-3",
        "title": "test-view-to-bulk-update-1-title-same-updated"
      }
    ]
    """
    Then the response code should be 400
    """
    {
      "errors": {
          "1.title": "Title already exists.",
          "2.title": "Title already exists."
      }
    }
    """

  Scenario: given bulk update request with multiple items with the same id
    should return error
    When I am admin
    Then I do PUT /api/v4/bulk/views:
    """
    [
      {
        "_id": "test-view-to-bulk-update-1",
        "title": "test-view-to-bulk-update-1-title-same-updated"
      },
      {
        "_id": "test-view-to-bulk-update-1",
        "title": "test-view-to-bulk-update-1-title-same-updated"
      },
      {
        "_id": "test-view-to-bulk-update-1",
        "title": "test-view-to-bulk-update-1-title-same-updated"
      }
    ]
    """
    Then the response code should be 400
    """
    {
      "errors": {
          "1._id": "ID already exists.",
          "2._id": "ID already exists."
      }
    }
    """

  Scenario: given bulk update request and no auth user should not allow access
    When I do PUT /api/v4/bulk/views
    Then the response code should be 401

  Scenario: given bulk update request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/bulk/views
    Then the response code should be 403

  Scenario: given bulk update request and auth user without view permission should not allow access
    When I am admin
    Then I do PUT /api/v4/bulk/views:
    """
    [
      {
        "_id": "test-view-to-bulk-update-1",
        "enabled": true,
        "group": "test-viewgroup-to-view-edit",
        "title": "test-view-to-bulk-update-1-title"
      },
      {
        "_id": "test-view-to-check-access",
        "enabled": true,
        "group": "test-viewgroup-to-view-edit",
        "title": "test-view-to-check-access-title"
      }
    ]
    """
    Then the response code should be 403

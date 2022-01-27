Feature: Update a view
  I need to be able to update a view
  Only admin should be able to update a view

  Scenario: given update request should update view
    When I am admin
    Then I do PUT /api/v4/views/test-view-to-update:
    """
    {
      "enabled": true,
      "title": "test-view-to-update-1-title",
      "title": "test-view-to-update-1-title-updated",
      "description": "test-view-to-update-1-description-updated",
      "group": "test-viewgroup-to-view-edit",
      "tags": ["test-view-to-update-1-tag-updated"],
      "periodic_refresh": {
        "enabled": true,
        "value": 10,
        "unit": "m"
      },
      "tabs": [
        {
          "_id": "test-view-to-update-1-tab-1-updated",
          "title": "test-view-to-update-1-tab-1-title-updated",
          "widgets": [
            {
              "_id": "test-view-to-update-1-tab-1-widget-1-updated",
              "grid_parameters": {
                "test-view-to-update-1-tab-1-widget-1-gridparameter-updated": "test-view-to-update-1-tab-1-widget-1-gridparameter-value-updated"
              },
              "parameters": {
                "test-view-to-update-1-tab-1-widget-1-parameter-updated": "test-view-to-update-1-tab-1-widget-1-parameter-value-updated"
              },
              "title": "test-view-to-update-1-tab-1-widget-1-title-updated",
              "type": "test-view-to-update-1-tab-1-widget-1-type-updated"
            }
          ]
        }
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "_id": "test-view-to-update",
      "author": "root",
      "created": 1611229670,
      "description": "test-view-to-update-1-description-updated",
      "enabled": true,
      "group": {
        "_id": "test-viewgroup-to-view-edit",
        "author": "test-viewgroup-to-view-edit-author",
        "created": 1611229670,
        "title": "test-viewgroup-to-view-edit-title",
        "updated": 1611229670
      },
      "title": "test-view-to-update-1-title",
      "periodic_refresh": {
        "enabled": true,
        "value": 10,
        "unit": "m"
      },
      "tabs": [
        {
          "_id": "test-view-to-update-1-tab-1-updated",
          "title": "test-view-to-update-1-tab-1-title-updated",
          "widgets": [
            {
              "_id": "test-view-to-update-1-tab-1-widget-1-updated",
              "grid_parameters": {
                "test-view-to-update-1-tab-1-widget-1-gridparameter-updated": "test-view-to-update-1-tab-1-widget-1-gridparameter-value-updated"
              },
              "parameters": {
                "test-view-to-update-1-tab-1-widget-1-parameter-updated": "test-view-to-update-1-tab-1-widget-1-parameter-value-updated"
              },
              "title": "test-view-to-update-1-tab-1-widget-1-title-updated",
              "type": "test-view-to-update-1-tab-1-widget-1-type-updated"
            }
          ]
        }
      ],
      "tags": [
        "test-view-to-update-1-tag-updated"
      ],
      "title": "test-view-to-update-1-title-updated"
    }
    """

  Scenario: given update request and no auth user should not allow access
    When I do PUT /api/v4/views/test-view-to-update
    Then the response code should be 401

  Scenario: given update request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/views/test-view-to-update
    Then the response code should be 403

  Scenario: given update request and auth user without view permission should not allow access
    When I am admin
    When I do PUT /api/v4/views/test-view-to-check-access
    Then the response code should be 403

  Scenario: given invalid update request should return errors
    When I am admin
    When I do PUT /api/v4/views/test-view-to-update:
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

  Scenario: given update request with not exist id should return not allow access error
    When I am admin
    When I do PUT /api/v4/views/test-view-not-found:
    """
    {
      "description": "test-view-not-found-description",
      "enabled": true,
      "group": "test-viewgroup-to-view-edit",
      "title": "test-view-not-found-title",
      "tabs": [],
      "tags": [],
      "title": "test-view-not-found-title"
    }
    """
    Then the response code should be 403

Feature: Update a view
  I need to be able to update a view
  Only admin should be able to update a view

  Scenario: given update request should update view
    When I am admin
    Then I do PUT /api/v4/views/test-view-to-update:
    """json
    {
      "enabled": true,
      "title": "test-view-to-update-1-title-updated",
      "description": "test-view-to-update-1-description-updated",
      "group": "test-viewgroup-to-view-edit",
      "tags": ["test-view-to-update-1-tag-updated"],
      "periodic_refresh": {
        "enabled": true,
        "value": 10,
        "unit": "m"
      }
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-view-to-update",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "created": 1611229670,
      "description": "test-view-to-update-1-description-updated",
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
      "title": "test-view-to-update-1-title-updated",
      "periodic_refresh": {
        "enabled": true,
        "value": 10,
        "unit": "m"
      },
      "tags": [
        "test-view-to-update-1-tag-updated"
      ]
    }
    """
    When I do GET /api/v4/permissions?search=test-view-to-update
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-view-to-update",
          "name": "test-view-to-update",
          "description": "Rights on view : test-view-to-update-1-title-updated",
          "view": {
            "_id": "test-view-to-update",
            "title": "test-view-to-update-1-title-updated"
          },
          "view_group": {
            "_id": "test-viewgroup-to-view-edit",
            "title": "test-viewgroup-to-view-edit-title"
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

  Scenario: given update request and no auth user should not allow access
    When I do PUT /api/v4/views/test-view-to-update
    Then the response code should be 401

  Scenario: given update request and auth user without permissions should not allow access
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

  Scenario: given update request with not exist id should return not allow access error
    When I am admin
    When I do PUT /api/v4/views/test-view-not-found:
    """json
    {
      "description": "test-view-not-found-description",
      "enabled": true,
      "group": "test-viewgroup-to-view-edit",
      "title": "test-view-not-found-title",
      "tags": []
    }
    """
    Then the response code should be 403

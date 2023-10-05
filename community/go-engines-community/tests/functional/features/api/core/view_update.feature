Feature: Update a view
  I need to be able to update a view
  Only admin should be able to update a view

  @concurrent
  Scenario: given update request should update view
    When I am admin
    Then I do PUT /api/v4/views/test-view-to-update-1:
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
      "_id": "test-view-to-update-1",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "created": 1611229670,
      "description": "test-view-to-update-1-description-updated",
      "enabled": true,
      "is_private": false,
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
    When I do GET /api/v4/permissions?search=test-view-to-update-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-view-to-update-1",
          "name": "test-view-to-update-1",
          "description": "Rights on view : test-view-to-update-1-title-updated",
          "view": {
            "_id": "test-view-to-update-1",
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

  @concurrent
  Scenario: given update public view request with not found group should return error
    When I am admin
    Then I do PUT /api/v4/views/test-view-to-update-2:
    """json
    {
      "enabled": true,
      "title": "test-view-to-update-2-title-updated",
      "description": "test-view-to-update-2-description-updated",
      "group": "not found",
      "tags": ["test-view-to-update-2-tag-updated"],
      "periodic_refresh": {
        "enabled": true,
        "value": 10,
        "unit": "m"
      }
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "group": "Group doesn't exist."
      }
    }
    """
    
  @concurrent
  Scenario: given update public view request with private group should return error
    When I am admin
    Then I do PUT /api/v4/views/test-view-to-update-3:
    """json
    {
      "enabled": true,
      "title": "test-view-to-update-3-title-updated",
      "description": "test-view-to-update-3-description-updated",
      "group": "test-private-viewgroup-to-update-view-1",
      "tags": ["test-view-to-update-3-tag-updated"],
      "periodic_refresh": {
        "enabled": true,
        "value": 10,
        "unit": "m"
      }
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "group": "Group is private."
      }
    }
    """

  @concurrent
  Scenario: given update public view request with public group should be ok
    When I am admin
    Then I do PUT /api/v4/views/test-view-to-update-4:
    """json
    {
      "enabled": true,
      "title": "test-view-to-update-4-title-updated",
      "description": "test-view-to-update-4-description-updated",
      "group": "test-viewgroup-to-update-view-1",
      "tags": ["test-view-to-update-4-tag-updated"],
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
      "is_private": false,
      "group": {
        "_id": "test-viewgroup-to-update-view-1"
      }
    }
    """
    When I do GET /api/v4/permissions?search=test-view-to-update-4
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "view_group": {
            "_id": "test-viewgroup-to-update-view-1",
            "title": "test-viewgroup-to-update-view-1-title"
          }
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

  @concurrent
  Scenario: given update request and no auth user should not allow access
    When I do PUT /api/v4/views/test-view-to-update-1
    Then the response code should be 401

  @concurrent
  Scenario: given update request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/views/test-view-to-update-1
    Then the response code should be 403

  @concurrent
  Scenario: given update request and auth user without view permission should not allow access
    When I am admin
    When I do PUT /api/v4/views/test-view-to-check-access
    Then the response code should be 403

  @concurrent
  Scenario: given invalid update request should return errors
    When I am admin
    When I do PUT /api/v4/views/test-view-to-update-1:
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

  @concurrent
  Scenario: given update request with not exist id should return not found error
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
    Then the response code should be 404

  @concurrent
  Scenario: given update owned private view request should return ok
    When I am admin
    When I do PUT /api/v4/views/test-private-view-to-update-1:
    """json
    {
      "enabled": true,
      "title": "test-private-view-to-update-1-title-updated",
      "description": "test-private-view-to-update-1-description-updated",
      "group": "test-private-viewgroup-to-update-view-1",
      "tags": ["test-private-view-to-update-1-tag-updated"],
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
      "_id": "test-private-view-to-update-1",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "description": "test-private-view-to-update-1-description-updated",
      "enabled": true,
      "is_private": true,
      "group": {
        "_id": "test-private-viewgroup-to-update-view-1",
        "author": {
          "_id": "root",
          "name": "root"
        },
        "title": "test-private-viewgroup-to-update-view-1-title"
      },
      "title": "test-private-view-to-update-1-title-updated",
      "periodic_refresh": {
        "enabled": true,
        "value": 10,
        "unit": "m"
      },
      "tags": [
        "test-private-view-to-update-1-tag-updated"
      ]
    }
    """
    
  @concurrent
  Scenario: given update not owned private view request should return not allow access
    When I am admin
    When I do PUT /api/v4/views/test-private-view-to-update-2:
    """json
    {
      "enabled": true,
      "title": "test-private-view-to-update-2-title-updated",
      "description": "test-private-view-to-update-2-description-updated",
      "group": "test-private-viewgroup-to-update-view-2",
      "tags": ["test-private-view-to-update-2-tag-updated"],
      "periodic_refresh": {
        "enabled": true,
        "value": 10,
        "unit": "m"
      }
    }
    """
    Then the response code should be 403
    
  @concurrent
  Scenario: given update private view request with public group should return error
    When I am admin
    When I do PUT /api/v4/views/test-private-view-to-update-3:
    """json
    {
      "enabled": true,
      "title": "test-private-view-to-update-3-title-updated",
      "description": "test-private-view-to-update-3-description-updated",
      "group": "test-viewgroup-to-update-view-1",
      "tags": ["test-private-view-to-update-3-tag-updated"],
      "periodic_refresh": {
        "enabled": true,
        "value": 10,
        "unit": "m"
      }
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "group": "Group is public."
      }
    }
    """

  @concurrent
  Scenario: given update private view request with not owned private group should return error
    When I am admin
    When I do PUT /api/v4/views/test-private-view-to-update-4:
    """json
    {
      "enabled": true,
      "title": "test-private-view-to-update-4-title-updated",
      "description": "test-private-view-to-update-4-description-updated",
      "group": "test-private-viewgroup-to-update-view-2",
      "tags": ["test-private-view-to-update-4-tag-updated"],
      "periodic_refresh": {
        "enabled": true,
        "value": 10,
        "unit": "m"
      }
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "group": "Group is private."
      }
    }
    """

  @concurrent
  Scenario: given update private view request with owned private group should return ok
    When I am admin
    When I do PUT /api/v4/views/test-private-view-to-update-5:
    """json
    {
      "enabled": true,
      "title": "test-private-view-to-update-5-title-updated",
      "description": "test-private-view-to-update-5-description-updated",
      "group": "test-private-viewgroup-to-update-view-1",
      "tags": ["test-private-view-to-update-5-tag-updated"],
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
      "_id": "test-private-view-to-update-5",
      "is_private": true,
      "group": {
        "_id": "test-private-viewgroup-to-update-view-1",
        "title": "test-private-viewgroup-to-update-view-1-title"
      }
    }
    """

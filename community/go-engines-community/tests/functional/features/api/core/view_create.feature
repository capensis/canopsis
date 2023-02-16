Feature: Create a view
  I need to be able to create a view
  Only admin should be able to create a view

  Scenario: given create request should return ok
    When I am admin
    When I do POST /api/v4/views:
    """json
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
      }
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "description": "test-view-to-create-1-description",
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
      "title": "test-view-to-create-1-title",
      "periodic_refresh": {
        "enabled": true,
        "value": 600,
        "unit": "m"
      },
      "tags": [
        "test-view-to-create-1-tag"
      ],
      "tabs": [
        {
          "author": {
            "_id": "root",
            "name": "root"
          },
          "title": "Default",
          "widgets": []
        }
      ]
    }
    """

  Scenario: given create request should return ok to get request
    When I am admin
    When I do POST /api/v4/views:
    """json
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
      }
    }
    """
    When I do GET /api/v4/views/{{ .lastResponse._id}}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "description": "test-view-to-create-2-description",
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
      "title": "test-view-to-create-2-title",
      "periodic_refresh": {
        "enabled": true,
        "value": 600,
        "unit": "m"
      },
      "tags": [
        "test-view-to-create-2-tag"
      ],
      "tabs": [
        {
          "author": {
            "_id": "root",
            "name": "root"
          },
          "title": "Default",
          "widgets": []
        }
      ]
    }
    """

  Scenario: given create request should create new permission
    When I am test-role-to-view-edit
    When I do POST /api/v4/views:
    """json
    {
      "enabled": true,
      "title": "test-view-to-create-3-title",
      "description": "test-view-to-create-3-description",
      "group": "test-viewgroup-to-view-edit",
      "tags": ["test-view-to-create-3-tag"]
    }
    """
    Then the response code should be 201
    Then I save response viewId={{ .lastResponse._id }}
    When I do GET /api/v4/views/{{ .viewId }}
    Then the response code should be 200
    Then I am admin
    When I do GET /api/v4/views/{{ .viewId }}
    Then the response code should be 200
    When I do GET /api/v4/permissions?search={{ .viewId }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "{{ .viewId }}",
          "name": "{{ .viewId }}",
          "description": "Rights on view : test-view-to-create-3-title",
          "view": {
            "_id": "{{ .viewId }}",
            "title": "test-view-to-create-3-title"
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
    When I do GET /api/v4/roles/admin
    Then the response code should be 200
    Then the response array key "permissions" should contain:
    """json
    [
      {
        "_id": "{{ .viewId }}",
        "name": "{{ .viewId }}",
        "description": "Rights on view : test-view-to-create-3-title",
        "type": "RW"
      }
    ]
    """
    When I do GET /api/v4/roles/test-role-to-view-edit
    Then the response code should be 200
    Then the response array key "permissions" should contain:
    """json
    [
      {
        "_id": "{{ .viewId }}",
        "name": "{{ .viewId }}",
        "description": "Rights on view : test-view-to-create-3-title",
        "type": "RW"
      }
    ]
    """

  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/views
    Then the response code should be 401

  Scenario: given create request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/views
    Then the response code should be 403

  Scenario: given invalid create request should return errors
    When I am admin
    When I do POST /api/v4/views:
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

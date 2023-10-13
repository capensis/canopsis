Feature: Create a view
  I need to be able to create a view
  Only admin should be able to create a view

  @concurrent
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
          "is_private": false,
          "widgets": []
        }
      ]
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
      "description": "test-view-to-create-1-description",
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
          "is_private": false,
          "title": "Default",
          "widgets": []
        }
      ]
    }
    """

  @concurrent
  Scenario: given create request should create new permission
    When I am test-role-to-view-edit
    When I do POST /api/v4/views:
    """json
    {
      "enabled": true,
      "title": "test-view-to-create-2-title",
      "description": "test-view-to-create-2-description",
      "group": "test-viewgroup-to-view-edit",
      "tags": ["test-view-to-create-2-tag"]
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
          "description": "Rights on view : test-view-to-create-2-title",
          "view": {
            "_id": "{{ .viewId }}",
            "title": "test-view-to-create-2-title"
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
        "description": "Rights on view : test-view-to-create-2-title",
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
        "description": "Rights on view : test-view-to-create-2-title",
        "type": "RW"
      }
    ]
    """

  @concurrent
  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/views
    Then the response code should be 401

  @concurrent
  Scenario: given create request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/views
    Then the response code should be 403

  @concurrent
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

  @concurrent
  Scenario: given create request with not found group should return error
    When I am admin
    When I do POST /api/v4/views:
    """json
    {
      "enabled": true,
      "title": "test-view-to-create-3-title",
      "description": "test-view-to-create-3-description",
      "group": "not found",
      "tags": ["test-view-to-create-3-tag"]
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
  Scenario: given create request for private group should create private view
    When I am admin
    When I do POST /api/v4/views:
    """json
    {
      "enabled": true,
      "title": "test-view-to-create-4-title",
      "description": "test-view-to-create-4-description",
      "group": "test-private-viewgroup-to-create-view-1",
      "tags": ["test-view-to-create-4-tag"],
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
      "description": "test-view-to-create-4-description",
      "enabled": true,
      "is_private": true,
      "group": {
        "_id": "test-private-viewgroup-to-create-view-1",
        "author": {
          "_id": "root",
          "name": "root"
        },
        "title": "test-private-viewgroup-to-create-view-1-title"
      },
      "title": "test-view-to-create-4-title",
      "periodic_refresh": {
        "enabled": true,
        "value": 600,
        "unit": "m"
      },
      "tags": [
        "test-view-to-create-4-tag"
      ],
      "tabs": [
        {
          "author": {
            "_id": "root",
            "name": "root"
          },
          "title": "Default",
          "is_private": true,
          "widgets": []
        }
      ]
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
      "description": "test-view-to-create-4-description",
      "enabled": true,
      "is_private": true,
      "group": {
        "_id": "test-private-viewgroup-to-create-view-1",
        "author": {
          "_id": "root",
          "name": "root"
        },
        "title": "test-private-viewgroup-to-create-view-1-title"
      },
      "title": "test-view-to-create-4-title",
      "periodic_refresh": {
        "enabled": true,
        "value": 600,
        "unit": "m"
      },
      "tags": [
        "test-view-to-create-4-tag"
      ],
      "tabs": [
        {
          "author": {
            "_id": "root",
            "name": "root"
          },
          "is_private": true,
          "title": "Default",
          "widgets": []
        }
      ]
    }
    """
    When I do GET /api/v4/permissions?search={{ .lastResponse._id}}
    Then the response code should be 200
    Then the response body should contain:
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

  @concurrent
  Scenario: given invalid create request for not owned private view should not allow access
    When I am admin
    When I do POST /api/v4/views:
    """json
    {
      "enabled": true,
      "title": "test-view-to-create-5-title",
      "description": "test-view-to-create-5-description",
      "group": "test-private-viewgroup-to-create-view-2",
      "tags": ["test-view-to-create-5-tag"]
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
  Scenario: given create request for private group with api_private_view_groups
    but without api_view permissions should return filters should create a private view
    When I am test-role-to-private-views-without-view-perm
    When I do POST /api/v4/views:
    """json
    {
      "enabled": true,
      "title": "test-view-to-create-6-title",
      "description": "test-view-to-create-6-description",
      "group": "test-private-viewgroup-to-create-view-3",
      "tags": ["test-view-to-create-6-tag"],
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
        "_id": "test-user-to-private-views-without-view-perm",
        "name": "test-user-to-private-views-without-view-perm"
      },
      "description": "test-view-to-create-6-description",
      "enabled": true,
      "is_private": true,
      "group": {
        "_id": "test-private-viewgroup-to-create-view-3",
        "author": {
          "_id": "test-user-to-private-views-without-view-perm",
          "name": "test-user-to-private-views-without-view-perm"
        },
        "title": "test-private-viewgroup-to-create-view-3-title"
      },
      "title": "test-view-to-create-6-title",
      "periodic_refresh": {
        "enabled": true,
        "value": 600,
        "unit": "m"
      },
      "tags": [
        "test-view-to-create-6-tag"
      ],
      "tabs": [
        {
          "author": {
            "_id": "test-user-to-private-views-without-view-perm",
            "name": "test-user-to-private-views-without-view-perm"
          },
          "title": "Default",
          "is_private": true,
          "widgets": []
        }
      ]
    }
    """
    When I do GET /api/v4/views/{{ .lastResponse._id}}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "test-user-to-private-views-without-view-perm",
        "name": "test-user-to-private-views-without-view-perm"
      },
      "description": "test-view-to-create-6-description",
      "enabled": true,
      "is_private": true,
      "group": {
        "_id": "test-private-viewgroup-to-create-view-3",
        "author": {
          "_id": "test-user-to-private-views-without-view-perm",
          "name": "test-user-to-private-views-without-view-perm"
        },
        "title": "test-private-viewgroup-to-create-view-3-title"
      },
      "title": "test-view-to-create-6-title",
      "periodic_refresh": {
        "enabled": true,
        "value": 600,
        "unit": "m"
      },
      "tags": [
        "test-view-to-create-6-tag"
      ],
      "tabs": [
        {
          "author": {
            "_id": "test-user-to-private-views-without-view-perm",
            "name": "test-user-to-private-views-without-view-perm"
          },
          "title": "Default",
          "is_private": true,
          "widgets": []
        }
      ]
    }
    """

  @concurrent
  Scenario: given create request for public group with api_private_view_groups
    but without api_view permissions should return filters should return error
    When I am test-role-to-private-views-without-view-perm
    When I do POST /api/v4/views:
    """json
    {
      "enabled": true,
      "title": "test-view-to-create-7-title",
      "description": "test-view-to-create-7-description",
      "group": "test-viewgroup-to-view-edit",
      "tags": ["test-view-to-create-7-tag"],
      "periodic_refresh": {
        "enabled": true,
        "value": 700,
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

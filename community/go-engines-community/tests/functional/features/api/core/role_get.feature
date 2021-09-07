Feature: Get a role
  I need to be able to get a role
  Only admin should be able to get a role

  Scenario: given search request should return roles
    When I am admin
    When I do GET /api/v4/roles?search=test-role-to-get
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "data": [
        {
          "_id": "test-role-to-get-1",
          "defaultview": null,
          "description": "test-role-to-get-1-description",
          "name": "test-role-to-get-1",
          "permissions": [
            {
              "_id": "test-permission-to-edit-role-1",
              "actions": [],
              "description": "test-permission-to-edit-role-1-description",
              "name": "test-permission-to-edit-role-1",
              "type": ""
            },
            {
              "_id": "test-permission-to-edit-role-2",
              "actions": [
                "create",
                "delete",
                "read",
                "update"
              ],
              "description": "test-permission-to-edit-role-2-description",
              "name": "test-permission-to-edit-role-2",
              "type": "CRUD"
            },
            {
              "_id": "test-permission-to-edit-role-3",
              "actions": [
                "delete",
                "read",
                "update"
              ],
              "description": "test-permission-to-edit-role-3-description",
              "name": "test-permission-to-edit-role-3",
              "type": "RW"
            }
          ]
        },
        {
          "_id": "test-role-to-get-2",
          "defaultview": {
            "_id": "test-view-to-edit-role",
            "title": "test-view-to-edit-role-title"
          },
          "description": "test-role-to-get-2-description",
          "name": "test-role-to-get-2",
          "permissions": [
            {
              "_id": "test-permission-to-edit-role-2",
              "actions": [
                "delete",
                "read",
                "update"
              ],
              "description": "test-permission-to-edit-role-2-description",
              "name": "test-permission-to-edit-role-2",
              "type": "CRUD"
            },
            {
              "_id": "test-permission-to-edit-role-3",
              "actions": [
                "delete",
                "update"
              ],
              "description": "test-permission-to-edit-role-3-description",
              "name": "test-permission-to-edit-role-3",
              "type": "RW"
            }
          ]
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

  Scenario: given search request should return roles with permission
    When I am admin
    When I do GET /api/v4/roles?permission=api_instruction_approve
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "admin"
        },
        {
          "_id": "approver"
        },
        {
          "_id": "approver2"
        },
        {
          "_id": "manager"
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 4
      }
    }
    """

  Scenario: given search request should return roles with permission
    When I am admin
    When I do GET /api/v4/roles?permission=api_instruction_approve&search=ap
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "approver"
        },
        {
          "_id": "approver2"
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

  Scenario: given get request should return role
    When I am admin
    When I do GET /api/v4/roles/test-role-to-get-1
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "_id": "test-role-to-get-1",
      "defaultview": null,
      "description": "test-role-to-get-1-description",
      "name": "test-role-to-get-1",
      "permissions": [
        {
          "_id": "test-permission-to-edit-role-1",
          "actions": [],
          "description": "test-permission-to-edit-role-1-description",
          "name": "test-permission-to-edit-role-1",
          "type": ""
        },
        {
          "_id": "test-permission-to-edit-role-2",
          "actions": [
            "create",
            "delete",
            "read",
            "update"
          ],
          "description": "test-permission-to-edit-role-2-description",
          "name": "test-permission-to-edit-role-2",
          "type": "CRUD"
        },
        {
          "_id": "test-permission-to-edit-role-3",
          "actions": [
            "delete",
            "read",
            "update"
          ],
          "description": "test-permission-to-edit-role-3-description",
          "name": "test-permission-to-edit-role-3",
          "type": "RW"
        }
      ]
    }
    """

  Scenario: given sort request should return sorted roles
    When I am admin
    When I do GET /api/v4/roles?search=test-role-to-get&sort=desc&sort_by=name
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-role-to-get-2"
        },
        {
          "_id": "test-role-to-get-1"
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

  Scenario: given get all request and no auth user should not allow access
    When I do GET /api/v4/roles
    Then the response code should be 401

  Scenario: given get all request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/roles
    Then the response code should be 403

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/roles/test-role-to-get-1
    Then the response code should be 401

  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/roles/test-role-to-get-1
    Then the response code should be 403

  Scenario: given get request with not exist id should return not found error
    When I am admin
    When I do GET /api/v4/roles/test-role-not-found
    Then the response code should be 404
    Then the response body should be:
    """
    {
      "error": "Not found"
    }
    """
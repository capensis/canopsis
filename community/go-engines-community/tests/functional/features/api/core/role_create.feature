Feature: Create a role
  I need to be able to create a role
  Only admin should be able to create a role

  Scenario: given create request should return ok
    When I am admin
    When I do POST /api/v4/roles:
    """
    {
      "name": "test-role-to-create-1-name",
      "description": "test-role-to-create-1-description",
      "defaultview": "test-view-to-edit-role",
      "permissions": {
        "test-permission-to-edit-role-1": [],
        "test-permission-to-edit-role-2": ["create", "read", "update", "delete"],
        "test-permission-to-edit-role-3": ["read", "update", "delete"]
      }
    }
    """
    Then the response code should be 201
    Then the response body should be:
    """
    {
      "_id": "test-role-to-create-1-name",
      "defaultview": {
        "_id": "test-view-to-edit-role",
        "title": "test-view-to-edit-role-title"
      },
      "description": "test-role-to-create-1-description",
      "name": "test-role-to-create-1-name",
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

  Scenario: given create request should return ok to get request
    When I am admin
    When I do POST /api/v4/roles:
    """
    {
      "name": "test-role-to-create-2-name",
      "description": "test-role-to-create-2-description",
      "defaultview": "test-view-to-edit-role",
      "permissions": {
        "test-permission-to-edit-role-1": [],
        "test-permission-to-edit-role-2": ["create", "read", "update", "delete"],
        "test-permission-to-edit-role-3": ["read", "update", "delete"]
      }
    }
    """
    When I do GET /api/v4/roles/{{ .lastResponse._id}}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "_id": "test-role-to-create-2-name",
      "defaultview": {
        "_id": "test-view-to-edit-role",
        "title": "test-view-to-edit-role-title"
      },
      "description": "test-role-to-create-2-description",
      "name": "test-role-to-create-2-name",
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

  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/roles
    Then the response code should be 401

  Scenario: given create request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/roles
    Then the response code should be 403

  Scenario: given invalid create request should return errors
    When I am admin
    When I do POST /api/v4/roles:
    """
    {
      "permissions": {
        "not-exist": ["create"],
        "test-permission-to-edit-role-1": ["create"],
        "test-permission-to-edit-role-2": ["not-exist"],
        "test-permission-to-edit-role-3": []
      }
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "name": "Name is missing.",
        "permissions.not-exist": "Permissions.not-exist doesn't exist.",
        "permissions.test-permission-to-edit-role-1": "Permissions.test-permission-to-edit-role-1 is not empty.",
        "permissions.test-permission-to-edit-role-2": "Permissions.test-permission-to-edit-role-2 must be one of [create read update delete].",
        "permissions.test-permission-to-edit-role-3": "Permissions.test-permission-to-edit-role-3 must be one of [read update delete]."
      }
    }
    """

  Scenario: given create request with already exists name should return error
    When I am admin
    When I do POST /api/v4/roles:
    """
    {
      "name": "test-role-to-check-unique-name",
      "description": "test-role-to-create-3-description",
      "defaultview": "test-view-to-edit-role",
      "permissions": {}
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
          "name": "Name already exists."
      }
    }
    """

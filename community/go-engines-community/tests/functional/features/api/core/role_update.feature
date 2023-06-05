Feature: Update a role
  I need to be able to update a role
  Only admin should be able to update a role

  Scenario: given update request should update role
    When I am admin
    Then I do PUT /api/v4/roles/test-role-to-update:
    """json
    {
      "description": "test-role-to-update-description-updated",
      "defaultview": "test-view-to-edit-role",
      "permissions": {
        "test-permission-to-edit-role-2": ["read"],
        "test-permission-to-edit-role-3": ["read"]
      },
      "auth_config": {
        "expiration_interval": {
          "value": 1,
          "unit": "m"
        },
        "inactivity_interval": {
          "value": 8,
          "unit": "h"
        }
      }
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-role-to-update",
      "defaultview": {
        "_id": "test-view-to-edit-role",
        "title": "test-view-to-edit-role-title"
      },
      "name": "test-role-to-update",
      "description": "test-role-to-update-description-updated",
      "permissions": [
        {
          "_id": "test-permission-to-edit-role-2",
          "actions": [
            "read"
          ],
          "description": "test-permission-to-edit-role-2-description",
          "name": "test-permission-to-edit-role-2",
          "type": "CRUD"
        },
        {
          "_id": "test-permission-to-edit-role-3",
          "actions": [
            "read"
          ],
          "description": "test-permission-to-edit-role-3-description",
          "name": "test-permission-to-edit-role-3",
          "type": "RW"
        }
      ],
      "auth_config": {
        "expiration_interval": {
          "value": 1,
          "unit": "m"
        },
        "inactivity_interval": {
          "value": 8,
          "unit": "h"
        }
      }
    }
    """

  Scenario: given get request and no auth user should not allow access
    When I do PUT /api/v4/roles/test-role-to-update
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/roles/test-role-to-update
    Then the response code should be 403

  Scenario: given invalid update request should return errors
    When I am admin
    When I do PUT /api/v4/roles/test-role-to-update:
    """json
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
    """json
    {
      "errors": {
        "permissions.not-exist": "Permissions.not-exist doesn't exist.",
        "permissions.test-permission-to-edit-role-1": "Permissions.test-permission-to-edit-role-1 is not empty.",
        "permissions.test-permission-to-edit-role-2": "Permissions.test-permission-to-edit-role-2 must be one of [create read update delete]."
      }
    }
    """

  Scenario: given update request with not exist id should return not found error
    When I am admin
    When I do PUT /api/v4/roles/test-role-not-found:
    """json
    {
      "description": "test-role-to-update-description-updated",
      "defaultview": "test-view-to-edit-role",
      "permissions": {
        "test-permission-to-edit-role-2": ["read"],
        "test-permission-to-edit-role-3": ["read"]
      }
    }
    """
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """

  Scenario: given update request for admin should return validation error
    When I am admin
    When I do PUT /api/v4/roles/admin
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "error": "admin cannot be updated"
    }
    """

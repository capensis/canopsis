Feature: Account auth user
  I need to be able to get auth user account.

  Scenario: given user username and password should get account
    When I do POST /api/v4/login:
    """json
    {
      "username": "test-user-to-get-1",
      "password": "test"
    }
    """
    Then the response code should be 200
    When I set header Authorization=Bearer {{ .lastResponse.access_token }}
    When I do GET /api/v4/account/me
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-user-to-get-1",
      "authkey": "3ct2e1ff-5e9e-4b1f-9d80-d968d61g5202",
      "defaultview": {
        "_id": "test-view-to-edit-user",
        "title": "test-view-to-edit-user-title"
      },
      "email": "test-user-to-get-1-email@canopsis.net",
      "enable": true,
      "external_id": "",
      "firstname": "test-user-to-get-1-firstname",
      "lastname": "test-user-to-get-1-lastname",
      "name": "test-user-to-get-1",
      "permissions": [
        {
          "_id": "test-permission-to-edit-user-1",
          "actions": [],
          "description": "test-permission-to-edit-user-1-description",
          "name": "test-permission-to-edit-user-1",
          "type": ""
        },
        {
          "_id": "test-permission-to-edit-user-2",
          "actions": [
            "create",
            "delete",
            "read",
            "update"
          ],
          "description": "test-permission-to-edit-user-2-description",
          "name": "test-permission-to-edit-user-2",
          "type": "CRUD"
        },
        {
          "_id": "test-permission-to-edit-user-3",
          "actions": [
            "delete",
            "read",
            "update"
          ],
          "description": "test-permission-to-edit-user-3-description",
          "name": "test-permission-to-edit-user-3",
          "type": "RW"
        }
      ],
      "role": {
        "_id": "test-role-to-user-get-1",
        "name": "test-role-to-user-get-1",
        "defaultview": {
          "_id": "test-view-to-edit-user",
          "title": "test-view-to-edit-user-title"
        }
      },
      "source": "",
      "ui_groups_navigation_type": "side-bar",
      "ui_language": "en",
      "ui_theme": "canopsis",
      "ui_tours": {
        "test-tour-to-get-user-1": true
      }
    }
	"""

  Scenario: given unauth request should not allow access
    When I do GET /api/v4/account/me
    Then the response code should be 401

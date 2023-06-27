Feature: Account auth user
  I need to be able to get auth user account.

  Scenario: given user username and password should get account
    When I do POST /api/v4/login:
    """json
    {
      "username": "test-user-to-account-get-1",
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
      "_id": "test-user-to-account-get-1",
      "authkey": "5ez4e3jj-7e1e-5c2g-0e91-e179f72a6426",
      "defaultview": {
        "_id": "test-view-to-edit-user",
        "title": "test-view-to-edit-user-title"
      },
      "email": "test-user-to-account-get-1-email@canopsis.net",
      "enable": true,
      "external_id": "",
      "firstname": "test-user-to-account-get-1-firstname",
      "lastname": "test-user-to-account-get-1-lastname",
      "name": "test-user-to-account-get-1",
      "display_name": "test-user-to-account-get-1 test-user-to-account-get-1-firstname test-user-to-account-get-1-lastname test-user-to-account-get-1-email@canopsis.net",
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
      "roles": [
        {
          "_id": "test-role-to-account-get-2",
          "name": "test-role-to-account-get-2",
          "defaultview": null
        },
        {
          "_id": "test-role-to-account-get-1",
          "name": "test-role-to-account-get-1",
          "defaultview": {
            "_id": "test-view-to-edit-user",
            "title": "test-view-to-edit-user-title"
          }
        }
      ],
      "source": "",
      "ui_groups_navigation_type": "side-bar",
      "ui_language": "en",
      "ui_theme": "canopsis",
      "ui_tours": {
        "test-tour-to-account-get-1": true
      }
    }
	"""

  Scenario: given unauth request should not allow access
    When I do GET /api/v4/account/me
    Then the response code should be 401

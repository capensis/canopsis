Feature: Update an account
  I need to be able to update an account
  Only admin should be able to update an account

  Scenario: given update request should update user
    When I am test-role-to-account-update-1
    Then I do PUT /api/v4/account/me:
    """json
    {
      "ui_language": "fr",
      "ui_theme": "canopsis",
      "ui_groups_navigation_type": "top-bar",
      "defaultview": "test-view-to-edit-user",
      "ui_tours": {
        "test-tour-to-update-user-1": true
      }
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-user-to-account-update-1",
      "authkey": "5ez4e3jj-7e1e-5c2g-0e91-e079f72o6425",
      "defaultview": {
        "_id": "test-view-to-edit-user",
        "title": "test-view-to-edit-user-title"
      },
      "email": "test-user-to-account-update-1-email@canopsis.net",
      "enable": true,
      "external_id": "",
      "firstname": "test-user-to-account-update-1-firstname",
      "lastname": "test-user-to-account-update-1-lastname",
      "name": "test-user-to-account-update-1",
      "role": {
        "_id": "test-role-to-account-update-1",
        "name": "test-role-to-account-update-1",
        "defaultview": null
      },
      "permissions": [],
      "source": "",
      "ui_groups_navigation_type": "top-bar",
      "ui_language": "fr",
      "ui_theme": "canopsis",
      "ui_tours": {
        "test-tour-to-update-user-1": true
      }
    }
    """
    When I am authenticated with username "test-user-to-account-update-1" and password "test"
    When I do GET /api/v4/account/me
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-user-to-account-update-1",
      "name": "test-user-to-account-update-1"
    }
    """

  Scenario: given update request with password should auth user by base auth
    When I am test-role-to-account-update-2
    Then I do PUT /api/v4/account/me:
    """json
    {
      "password": "test-password-updated",
      "name": "test-user-to-account-update-2",
      "firstname": "test-user-to-account-update-2-firstname-updated",
      "lastname": "test-user-to-account-update-2-lastname-updated",
      "email": "test-user-to-account-update-2-email-updated@canopsis.net",
      "role": "test-role-to-edit-user",
      "ui_language": "fr",
      "ui_theme": "canopsis",
      "ui_groups_navigation_type": "top-bar",
      "enable": true,
      "defaultview": "test-view-to-edit-user",
      "ui_tours": {
        "test-tour-to-update-user-2": true
      }
    }
    """
    When I am authenticated with username "test-user-to-account-update-2" and password "test-password-updated"
    When I do GET /api/v4/account/me
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-user-to-account-update-2",
      "name": "test-user-to-account-update-2"
    }
    """

  Scenario: given get request and no auth user should not allow access
    Then I do PUT /api/v4/account/me
    Then the response code should be 401

  Scenario: given invalid update request should return errors
    When I am admin
    Then I do PUT /api/v4/account/me:
    """json
    {
      "defaultview": "not-exist"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "defaultview": "DefaultView doesn't exist."
      }
    }
    """

  Scenario: given update request with invalid password should return error
    When I am admin
    Then I do PUT /api/v4/account/me:
    """json
    {
      "password": "1"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
          "password": "Password should be 8 or more."
      }
    }
    """

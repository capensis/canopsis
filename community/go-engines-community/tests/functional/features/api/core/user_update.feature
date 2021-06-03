Feature: Update a user
  I need to be able to update a user
  Only admin should be able to update a user

  Scenario: given update request should update user
    When I am admin
    Then I do PUT /api/v4/users/test-user-to-update-1:
    """
    {
      "name": "test-user-to-update-1-updated",
      "firstname": "test-user-to-update-1-firstname-updated",
      "lastname": "test-user-to-update-1-lastname-updated",
      "email": "test-user-to-update-1-email-updated@canopsis.net",
      "role": "test-role-to-edit-user",
      "ui_language": "fr",
      "ui_groups_navigation_type": "top-bar",
      "enable": true,
      "defaultview": "test-view-to-edit-user",
      "ui_tours": {
        "test-tour-to-update-user-1": true
      }
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "_id": "test-user-to-update-1",
      "authkey": "5ez4e3jj-7e1e-5c2g-0e91-e079f72o6424",
      "defaultview": {
        "_id": "test-view-to-edit-user",
        "title": "test-view-to-edit-user-title"
      },
      "email": "test-user-to-update-1-email-updated@canopsis.net",
      "enable": true,
      "external_id": "",
      "firstname": "test-user-to-update-1-firstname-updated",
      "lastname": "test-user-to-update-1-lastname-updated",
      "name": "test-user-to-update-1-updated",
      "role": {
        "_id": "test-role-to-edit-user",
        "name": "test-role-to-edit-user"
      },
      "source": "",
      "ui_groups_navigation_type": "top-bar",
      "ui_language": "fr",
      "ui_tours": {
        "test-tour-to-update-user-1": true
      }
    }
    """

  Scenario: given update request with password should auth user by base auth
    When I am admin
    Then I do PUT /api/v4/users/test-user-to-update-2:
    """
    {
      "password": "test-password-updated",
      "name": "test-user-to-update-2",
      "firstname": "test-user-to-update-2-firstname-updated",
      "lastname": "test-user-to-update-2-lastname-updated",
      "email": "test-user-to-update-2-email-updated@canopsis.net",
      "role": "test-role-to-edit-user",
      "ui_language": "fr",
      "ui_groups_navigation_type": "top-bar",
      "enable": true,
      "defaultview": "test-view-to-edit-user",
      "ui_tours": {
        "test-tour-to-update-user-2": true
      }
    }
    """
    When I am authenticated with username "test-user-to-update-2" and password "test-password-updated"
    When I do GET /api/v4/account/me
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "_id": "test-user-to-update-2",
      "name": "test-user-to-update-2"
    }
    """

  Scenario: given get request and no auth user should not allow access
    When I do PUT /api/v4/users/test-user-to-update
    Then the response code should be 401

  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/users/test-user-to-update
    Then the response code should be 403

  Scenario: given invalid update request should return errors
    When I am admin
    When I do PUT /api/v4/users/test-user-to-update:
    """
    {
      "role": "not-exist",
      "defaultview": "not-exist"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "defaultview": "DefaultView doesn't exist.",
        "email": "Email is missing.",
        "enable": "IsEnabled is missing.",
        "name": "Name is missing.",
        "role": "Role doesn't exist."
      }
    }
    """

  Scenario: given update request with not exist id should return not found error
    When I am admin
    When I do PUT /api/v4/users/test-user-not-found:
    """
    {
      "name": "test-user-to-update-name",
      "firstname": "test-user-to-update-firstname-updated",
      "lastname": "test-user-to-update-lastname-updated",
      "email": "test-user-to-update-email-updated@canopsis.net",
      "role": "test-role-to-edit-user",
      "ui_language": "fr",
      "groups_navigation_type": "top-bar",
      "enable": true,
      "defaultview": "test-view-to-edit-user"
    }
    """
    Then the response code should be 404
    Then the response body should be:
    """
    {
      "error": "Not found"
    }
    """

  Scenario: given create request with already exists name should return error
    When I am admin
    When I do PUT /api/v4/users/notexit:
    """
    {
      "name": "test-user-to-check-unique-name-name"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
          "name": "Name already exists."
      }
    }
    """

  Scenario: given create request with already exists id should return error
    When I am admin
    When I do PUT /api/v4/users/notexit:
    """
    {
      "name": "test-user-to-check-unique-name"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
          "name": "Name already exists."
      }
    }
    """
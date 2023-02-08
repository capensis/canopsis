Feature: Create a user
  I need to be able to create a user
  Only admin should be able to create a user

  Scenario: given create request should return ok
    When I am admin
    When I do POST /api/v4/users:
    """json
    {
      "name": "test-user-to-create-1-name",
      "firstname": "test-user-to-create-1-firstname",
      "lastname": "test-user-to-create-1-lastname",
      "email": "test-user-to-create-1-email@canopsis.net",
      "role": "test-role-to-edit-user",
      "ui_language": "fr",
      "ui_theme": "canopsis",
      "ui_groups_navigation_type": "top-bar",
      "enable": true,
      "defaultview": "test-view-to-edit-user",
      "password": "test-password"
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "_id": "test-user-to-create-1-name",
      "defaultview": {
        "_id": "test-view-to-edit-user",
        "title": "test-view-to-edit-user-title"
      },
      "email": "test-user-to-create-1-email@canopsis.net",
      "enable": true,
      "external_id": "",
      "firstname": "test-user-to-create-1-firstname",
      "lastname": "test-user-to-create-1-lastname",
      "name": "test-user-to-create-1-name",
      "role": {
        "_id": "test-role-to-edit-user",
        "name": "test-role-to-edit-user",
        "defaultview": {
          "_id": "test-view-to-edit-user",
          "title": "test-view-to-edit-user-title"
        }
      },
      "source": "",
      "ui_groups_navigation_type": "top-bar",
      "ui_language": "fr",
      "ui_theme": "canopsis"
    }
    """

  Scenario: given create request should return ok to get request
    When I am admin
    When I do POST /api/v4/users:
    """json
    {
      "name": "test-user-to-create-2-name",
      "firstname": "test-user-to-create-2-firstname",
      "lastname": "test-user-to-create-2-lastname",
      "email": "test-user-to-create-2-email@canopsis.net",
      "role": "test-role-to-edit-user",
      "ui_language": "fr",
      "ui_theme": "canopsis",
      "ui_groups_navigation_type": "top-bar",
      "enable": true,
      "defaultview": "test-view-to-edit-user",
      "password": "test-password"
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/users/{{ .lastResponse._id}}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-user-to-create-2-name",
      "defaultview": {
        "_id": "test-view-to-edit-user",
        "title": "test-view-to-edit-user-title"
      },
      "email": "test-user-to-create-2-email@canopsis.net",
      "enable": true,
      "external_id": "",
      "firstname": "test-user-to-create-2-firstname",
      "lastname": "test-user-to-create-2-lastname",
      "name": "test-user-to-create-2-name",
      "role": {
        "_id": "test-role-to-edit-user",
        "name": "test-role-to-edit-user",
        "defaultview": {
          "_id": "test-view-to-edit-user",
          "title": "test-view-to-edit-user-title"
        }
      },
      "source": "",
      "ui_groups_navigation_type": "top-bar",
      "ui_language": "fr",
      "ui_theme": "canopsis"
    }
    """

  Scenario: given create request should auth new user by base auth
    When I am admin
    When I do POST /api/v4/users:
    """json
    {
      "name": "test-user-to-create-3-name",
      "firstname": "test-user-to-create-3-firstname",
      "lastname": "test-user-to-create-3-lastname",
      "email": "test-user-to-create-3-email@canopsis.net",
      "role": "test-role-to-edit-user",
      "ui_language": "fr",
      "ui_theme": "canopsis",
      "ui_groups_navigation_type": "top-bar",
      "enable": true,
      "defaultview": "test-view-to-edit-user",
      "password": "test-password"
    }
    """
    Then the response code should be 201
    When I am authenticated with username "test-user-to-create-3-name" and password "test-password"
    When I do GET /api/v4/account/me
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-user-to-create-3-name",
      "name": "test-user-to-create-3-name"
    }
    """

  Scenario: given create request should auth new user by password
    When I am admin
    When I do POST /api/v4/users:
    """json
    {
      "name": "test-user-to-create-4-name",
      "firstname": "test-user-to-create-4-firstname",
      "lastname": "test-user-to-create-4-lastname",
      "email": "test-user-to-create-4-email@canopsis.net",
      "role": "test-role-to-edit-user",
      "ui_language": "fr",
      "ui_theme": "canopsis",
      "ui_groups_navigation_type": "top-bar",
      "enable": true,
      "defaultview": "test-view-to-edit-user",
      "password": "test-password"
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/login:
    """json
    {
      "username": "test-user-to-create-4-name",
      "password": "test-password"
    }
    """
    When I set header Authorization=Bearer {{ .lastResponse.access_token }}
    When I do GET /api/v4/account/me
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-user-to-create-4-name",
      "name": "test-user-to-create-4-name"
    }
    """

  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/users
    Then the response code should be 401

  Scenario: given create request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/users
    Then the response code should be 403

  Scenario: given invalid create request should return errors
    When I am admin
    When I do POST /api/v4/users:
    """json
    {
      "role": "not-exist",
      "defaultview": "not-exist"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "defaultview": "DefaultView doesn't exist.",
        "email": "Email is missing.",
        "enable": "IsEnabled is missing.",
        "name": "Name is missing.",
        "password": "Password is missing.",
        "role": "Role doesn't exist."
      }
    }
    """

  Scenario: given create request with already exists name should return error
    When I am admin
    When I do POST /api/v4/users:
    """json
    {
      "name": "test-user-to-check-unique-name-name"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
          "name": "Name already exists."
      }
    }
    """

  Scenario: given create request with already exists id should return error
    When I am admin
    When I do POST /api/v4/users:
    """json
    {
      "name": "test-user-to-check-unique-name"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
          "name": "Name already exists."
      }
    }
    """

  Scenario: given create request with invalid password should return error
    When I am admin
    When I do POST /api/v4/users:
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

  Scenario: given create request should create user with source and external_id
    When I am admin
    When I do POST /api/v4/users:
    """json
    {
      "name": "test-user-to-create-5-name",
      "firstname": "test-user-to-create-5-firstname",
      "lastname": "test-user-to-create-5-lastname",
      "email": "test-user-to-create-5-email@canopsis.net",
      "role": "test-role-to-edit-user",
      "ui_language": "fr",
      "ui_groups_navigation_type": "top-bar",
      "enable": true,
      "defaultview": "test-view-to-edit-user",
      "source": "saml",
      "external_id": "saml_id"
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/users/{{ .lastResponse._id}}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-user-to-create-5-name",
      "defaultview": {
        "_id": "test-view-to-edit-user",
        "title": "test-view-to-edit-user-title"
      },
      "email": "test-user-to-create-5-email@canopsis.net",
      "enable": true,
      "firstname": "test-user-to-create-5-firstname",
      "lastname": "test-user-to-create-5-lastname",
      "name": "test-user-to-create-5-name",
      "role": {
        "_id": "test-role-to-edit-user",
        "name": "test-role-to-edit-user",
        "defaultview": {
          "_id": "test-view-to-edit-user",
          "title": "test-view-to-edit-user-title"
        }
      },
      "ui_groups_navigation_type": "top-bar",
      "ui_language": "fr",
      "source": "saml",
      "external_id": "saml_id"
    }
    """

  Scenario: given create request when only source exists should return error
    When I am admin
    When I do POST /api/v4/users:
    """json
    {
      "name": "test-user-to-create-6-name",
      "firstname": "test-user-to-create-6-firstname",
      "lastname": "test-user-to-create-6-lastname",
      "email": "test-user-to-create-6-email@canopsis.net",
      "role": "test-role-to-edit-user",
      "ui_language": "fr",
      "ui_groups_navigation_type": "top-bar",
      "enable": true,
      "defaultview": "test-view-to-edit-user",
      "source": "saml"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
          "external_id": "ExternalID is required when Source is present."
      }
    }
    """

  Scenario: given create request when only external_id exists should return error
    When I am admin
    When I do POST /api/v4/users:
    """json
    {
      "name": "test-user-to-create-6-name",
      "firstname": "test-user-to-create-6-firstname",
      "lastname": "test-user-to-create-6-lastname",
      "email": "test-user-to-create-6-email@canopsis.net",
      "role": "test-role-to-edit-user",
      "ui_language": "fr",
      "ui_groups_navigation_type": "top-bar",
      "enable": true,
      "defaultview": "test-view-to-edit-user",
      "external_id": "saml_id"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
          "source": "Source is required when ExternalID is present."
      }
    }
    """

  Scenario: given create request with wrong source should return error
    When I am admin
    When I do POST /api/v4/users:
    """json
    {
      "name": "test-user-to-create-6-name",
      "firstname": "test-user-to-create-6-firstname",
      "lastname": "test-user-to-create-6-lastname",
      "email": "test-user-to-create-6-email@canopsis.net",
      "role": "test-role-to-edit-user",
      "ui_language": "fr",
      "ui_groups_navigation_type": "top-bar",
      "enable": true,
      "defaultview": "test-view-to-edit-user",
      "source": "some",
      "external_id": "saml_id"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
          "source": "Source must be one of [ldap cas saml] or empty."
      }
    }
    """

  Scenario: given create request with source and password should return error
    When I am admin
    When I do POST /api/v4/users:
    """json
    {
      "name": "test-user-to-create-6-name",
      "firstname": "test-user-to-create-6-firstname",
      "lastname": "test-user-to-create-6-lastname",
      "email": "test-user-to-create-6-email@canopsis.net",
      "role": "test-role-to-edit-user",
      "ui_language": "fr",
      "ui_groups_navigation_type": "top-bar",
      "enable": true,
      "defaultview": "test-view-to-edit-user",
      "password": "test-password",
      "source": "some",
      "external_id": "saml_id",
      "password": "qwerty123"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
          "source": "Can't be present both Source and Password."
      }
    }
    """

  Scenario: given create request without source and without password should return error
    When I am admin
    When I do POST /api/v4/users:
    """json
    {
      "name": "test-user-to-create-6-name",
      "firstname": "test-user-to-create-6-firstname",
      "lastname": "test-user-to-create-6-lastname",
      "email": "test-user-to-create-6-email@canopsis.net",
      "role": "test-role-to-edit-user",
      "ui_language": "fr",
      "ui_groups_navigation_type": "top-bar",
      "enable": true,
      "defaultview": "test-view-to-edit-user"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
          "password": "Password is missing."
      }
    }
    """

Feature: Bulk create users
  I need to be able to bulk create users
  Only admin should be able to bulk create users

  Scenario: given bulk create request and no auth user should not allow access
    When I do POST /api/v4/bulk/users
    Then the response code should be 401

  Scenario: given bulk create request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/bulk/users
    Then the response code should be 403

  Scenario: given bulk create request should return multi status and should be handled independently
    When I am admin
    When I do POST /api/v4/bulk/users:
    """json
    [
      {
        "name": "test-user-to-bulk-create-1-name",
        "firstname": "test-user-to-bulk-create-1-firstname",
        "lastname": "test-user-to-bulk-create-1-lastname",
        "email": "test-user-to-bulk-create-1-email@canopsis.net",
        "role": "test-role-to-edit-user",
        "ui_language": "fr",
        "ui_groups_navigation_type": "top-bar",
        "enable": true,
        "defaultview": "test-view-to-edit-user",
        "password": "test-password"
      },
      {
        "role": "not-exist",
        "defaultview": "not-exist"
      },
      {
        "name": "test-user-to-check-unique-name-name"
      },
      {
        "name": "test-user-to-check-unique-name"
      },
      [],
      {
        "name": "test-user-to-bulk-create-2-name",
        "firstname": "test-user-to-bulk-create-2-firstname",
        "lastname": "test-user-to-bulk-create-2-lastname",
        "email": "test-user-to-bulk-create-2-email@canopsis.net",
        "role": "test-role-to-edit-user",
        "ui_language": "fr",
        "ui_groups_navigation_type": "top-bar",
        "enable": true,
        "defaultview": "test-view-to-edit-user",
        "password": "test-password"
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "id": "test-user-to-bulk-create-1-name",
        "status": 200,
        "item": {
          "name": "test-user-to-bulk-create-1-name",
          "firstname": "test-user-to-bulk-create-1-firstname",
          "lastname": "test-user-to-bulk-create-1-lastname",
          "email": "test-user-to-bulk-create-1-email@canopsis.net",
          "role": "test-role-to-edit-user",
          "ui_language": "fr",
          "ui_groups_navigation_type": "top-bar",
          "enable": true,
          "defaultview": "test-view-to-edit-user",
          "password": "test-password"
        }
      },
      {
        "status": 400,
        "item": {
          "role": "not-exist",
          "defaultview": "not-exist"
        },
        "errors": {
          "defaultview": "DefaultView doesn't exist.",
          "email": "Email is missing.",
          "enable": "IsEnabled is missing.",
          "name": "Name is missing.",
          "password": "Password is missing.",
          "role": "Role doesn't exist."
        }
      },
      {
        "status": 400,
        "item": {
          "name": "test-user-to-check-unique-name-name"
        },
        "errors": {
          "name": "Name already exists."
        }
      },
      {
        "status": 400,
        "item": {
          "name": "test-user-to-check-unique-name"
        },
        "errors": {
          "name": "Name already exists."
        }
      },
      {
        "status": 400,
        "item": [],
        "error": "value doesn't contain object; it contains array"
      },
      {
        "id": "test-user-to-bulk-create-2-name",
        "status": 200,
        "item": {
          "name": "test-user-to-bulk-create-2-name",
          "firstname": "test-user-to-bulk-create-2-firstname",
          "lastname": "test-user-to-bulk-create-2-lastname",
          "email": "test-user-to-bulk-create-2-email@canopsis.net",
          "role": "test-role-to-edit-user",
          "ui_language": "fr",
          "ui_groups_navigation_type": "top-bar",
          "enable": true,
          "defaultview": "test-view-to-edit-user",
          "password": "test-password"
        }
      }
    ]
    """
    When I do GET /api/v4/users?search=test-user-to-bulk-create
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-user-to-bulk-create-1-name",
          "defaultview": {
            "_id": "test-view-to-edit-user",
            "title": "test-view-to-edit-user-title"
          },
          "email": "test-user-to-bulk-create-1-email@canopsis.net",
          "enable": true,
          "external_id": "",
          "firstname": "test-user-to-bulk-create-1-firstname",
          "lastname": "test-user-to-bulk-create-1-lastname",
          "name": "test-user-to-bulk-create-1-name",
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
          "ui_tours": null
        },
        {
          "_id": "test-user-to-bulk-create-2-name",
          "defaultview": {
            "_id": "test-view-to-edit-user",
            "title": "test-view-to-edit-user-title"
          },
          "email": "test-user-to-bulk-create-2-email@canopsis.net",
          "enable": true,
          "external_id": "",
          "firstname": "test-user-to-bulk-create-2-firstname",
          "lastname": "test-user-to-bulk-create-2-lastname",
          "name": "test-user-to-bulk-create-2-name",
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
          "ui_tours": null
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
    When I am authenticated with username "test-user-to-bulk-create-1-name" and password "test-password"
    When I do GET /api/v4/account/me
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-user-to-bulk-create-1-name",
      "name": "test-user-to-bulk-create-1-name"
    }
    """
    When I am authenticated with username "test-user-to-bulk-create-2-name" and password "test-password"
    When I do GET /api/v4/account/me
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-user-to-bulk-create-2-name",
      "name": "test-user-to-bulk-create-2-name"
    }
    """
    When I do POST /api/v4/login:
    """json
    {
      "username": "test-user-to-bulk-create-1-name",
      "password": "test-password"
    }
    """
    When I set header Authorization=Bearer {{ .lastResponse.access_token }}
    When I do GET /api/v4/account/me
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-user-to-bulk-create-1-name",
      "name": "test-user-to-bulk-create-1-name"
    }
    """
    When I do POST /api/v4/login:
    """json
    {
      "username": "test-user-to-bulk-create-2-name",
      "password": "test-password"
    }
    """
    When I set header Authorization=Bearer {{ .lastResponse.access_token }}
    When I do GET /api/v4/account/me
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-user-to-bulk-create-2-name",
      "name": "test-user-to-bulk-create-2-name"
    }
    """

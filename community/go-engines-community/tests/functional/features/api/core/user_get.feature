Feature: Get a user
  I need to be able to get a user
  Only admin should be able to get a user

  Scenario: given search request should return users
    When I am admin
    When I do GET /api/v4/users?search=test-user-to-get
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "data": [
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
          "role": {
            "_id": "test-role-to-edit-user",
            "name": "test-role-to-edit-user"
          },
          "source": "",
          "ui_groups_navigation_type": "side-bar",
          "ui_language": "en",
          "ui_tours": {
            "test-tour-to-get-user-1": true
          }
        },
        {
          "_id": "test-user-to-get-2",
          "authkey": "4du3d2gg-6d0d-5c2g-0e91-e079f72o6313",
          "defaultview": {
            "_id": "test-view-to-edit-user",
            "title": "test-view-to-edit-user-title"
          },
          "email": "test-user-to-get-2-email@canopsis.net",
          "enable": true,
          "external_id": "",
          "firstname": "test-user-to-get-2-firstname",
          "lastname": "test-user-to-get-2-lastname",
          "name": "test-user-to-get-2",
          "role": {
            "_id": "test-role-to-edit-user",
            "name": "test-role-to-edit-user"
          },
          "source": "",
          "ui_groups_navigation_type": "side-bar",
          "ui_language": "en",
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

  Scenario: given search request should return users with permission
    When I am admin
    When I do GET /api/v4/users?permission=api_instruction_approve
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "approveruser"
        },
        {
          "_id": "approveruser2"
        },
        {
          "_id": "manageruser"
        },
        {
          "_id": "root"
        },
        {
          "_id": "test-user-to-test-paused-executions"
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 5
      }
    }
    """

  Scenario: given search request should return users with permission and search query
    When I am admin
    When I do GET /api/v4/users?permission=api_instruction_approve&search=approveruser
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "approveruser"
        },
        {
          "_id": "approveruser2"
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

  Scenario: given get request should return user
    When I am admin
    When I do GET /api/v4/users/test-user-to-get-1
    Then the response code should be 200
    Then the response body should be:
    """
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
      "role": {
        "_id": "test-role-to-edit-user",
        "name": "test-role-to-edit-user"
      },
      "source": "",
      "ui_groups_navigation_type": "side-bar",
      "ui_language": "en",
      "ui_tours": {
        "test-tour-to-get-user-1": true
      }
    }
    """

  Scenario: given sort request should return sorted users
    When I am admin
    When I do GET /api/v4/users?search=test-user-to-get&sort=desc&sort_by=name
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "data": [
        {
          "_id": "test-user-to-get-2"
        },
        {
          "_id": "test-user-to-get-1"
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
    When I do GET /api/v4/users
    Then the response code should be 401

  Scenario: given get all request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/users
    Then the response code should be 403

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/users/test-user-to-get-1
    Then the response code should be 401

  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/users/test-user-to-get-1
    Then the response code should be 403

  Scenario: given get request with not exist id should return not found error
    When I am admin
    When I do GET /api/v4/users/test-user-not-found
    Then the response code should be 404
    Then the response body should be:
    """
    {
      "error": "Not found"
    }
    """
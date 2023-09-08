Feature: Bulk create users
  I need to be able to bulk create users
  Only admin should be able to bulk create users

  @concurrent
  Scenario: given bulk create request and no auth user should not allow access
    When I do POST /api/v4/bulk/users
    Then the response code should be 401

  @concurrent
  Scenario: given bulk create request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/bulk/users
    Then the response code should be 403

  @concurrent
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
        "roles": [
          "test-role-to-user-edit-2",
          "test-role-to-user-edit-1"
        ],
        "ui_language": "fr",
        "ui_theme": "canopsis",
        "ui_groups_navigation_type": "top-bar",
        "enable": true,
        "defaultview": "test-view-to-edit-user",
        "password": "test-password"
      },
      {
        "name": "test-user-to-bulk-create-1-name",
        "firstname": "test-user-to-bulk-create-1-firstname",
        "lastname": "test-user-to-bulk-create-1-lastname",
        "email": "test-user-to-bulk-create-1-email@canopsis.net",
        "roles": [
          "test-role-to-user-edit-1"
        ],
        "ui_language": "fr",
        "ui_theme": "canopsis",
        "ui_groups_navigation_type": "top-bar",
        "enable": true,
        "defaultview": "test-view-to-edit-user",
        "password": "test-password"
      },
      {
        "roles": [
          "test-role-to-user-edit-1",
          "not-exist"
        ],
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
        "roles": [
          "test-role-to-user-edit-1"
        ],
        "ui_language": "fr",
        "ui_theme": "canopsis_dark",
        "ui_groups_navigation_type": "top-bar",
        "enable": true,
        "defaultview": "test-view-to-edit-user",
        "password": "test-password"
      },
      {
        "name": "test-user-to-bulk-create-3-name",
        "firstname": "test-user-to-bulk-create-3-firstname",
        "lastname": "test-user-to-bulk-create-3-lastname",
        "email": "test-user-to-bulk-create-3-email@canopsis.net",
        "roles": [
          "test-role-to-user-edit-1"
        ],
        "ui_language": "fr",
        "ui_groups_navigation_type": "top-bar",
        "enable": true,
        "defaultview": "test-view-to-edit-user",
        "source": "saml",
        "external_id": "test-id"
      },
      {
        "name": "test-user-to-bulk-create-4-name",
        "firstname": "test-user-to-bulk-create-4-firstname",
        "lastname": "test-user-to-bulk-create-4-lastname",
        "email": "test-user-to-bulk-create-4-email@canopsis.net",
        "roles": [
          "test-role-to-user-edit-1"
        ],
        "ui_language": "fr",
        "ui_groups_navigation_type": "top-bar",
        "enable": true,
        "defaultview": "test-view-to-edit-user",
        "password": "test-password",
        "source": "saml",
        "external_id": "test-id"
      },
      {
        "name": "test-user-to-bulk-create-5-name",
        "firstname": "test-user-to-bulk-create-5-firstname",
        "lastname": "test-user-to-bulk-create-5-lastname",
        "email": "test-user-to-bulk-create-5-email@canopsis.net",
        "roles": [
          "test-role-to-user-edit-1"
        ],
        "ui_language": "fr",
        "ui_groups_navigation_type": "top-bar",
        "enable": true,
        "defaultview": "test-view-to-edit-user",
        "source": "saml"
      },
      {
        "name": "test-user-to-bulk-create-6-name",
        "firstname": "test-user-to-bulk-create-6-firstname",
        "lastname": "test-user-to-bulk-create-6-lastname",
        "email": "test-user-to-bulk-create-6-email@canopsis.net",
        "roles": [
          "test-role-to-user-edit-1"
        ],
        "ui_language": "fr",
        "ui_groups_navigation_type": "top-bar",
        "enable": true,
        "defaultview": "test-view-to-edit-user",
        "source": "some",
        "external_id": "test-id"
      },
      {
        "name": "test-user-to-bulk-create-7-name",
        "firstname": "test-user-to-bulk-create-7-firstname",
        "lastname": "test-user-to-bulk-create-7-lastname",
        "email": "test-user-to-bulk-create-7-email@canopsis.net",
        "roles": [
          "test-role-to-user-edit-1"
        ],
        "ui_language": "fr",
        "ui_groups_navigation_type": "top-bar",
        "enable": true,
        "defaultview": "test-view-to-edit-user"
      },
      {
        "name": "test-user-to-bulk-create-8-name",
        "firstname": "test-user-to-bulk-create-8-firstname",
        "lastname": "test-user-to-bulk-create-8-lastname",
        "email": "test-user-to-bulk-create-8-email@canopsis.net",
        "roles": [
          "test-role-to-user-edit-1"
        ],
        "ui_language": "fr",
        "ui_groups_navigation_type": "top-bar",
        "enable": true,
        "defaultview": "test-view-to-edit-user",
        "password": "123"
      },
      {
        "name": "test-user-to-bulk-create-9-name",
        "firstname": "test-user-to-bulk-create-9-firstname",
        "lastname": "test-user-to-bulk-create-9-lastname",
        "email": "test-user-to-bulk-create-9-email@canopsis.net",
        "roles": [
          "test-role-to-user-edit-1"
        ],
        "ui_language": "fr",
        "ui_groups_navigation_type": "top-bar",
        "enable": true,
        "defaultview": "test-view-to-edit-user",
        "source": "saml",
        "external_id": "test-id",
        "ui_theme": "not found"
      },
      {
        "name": "test-user-to-bulk-create-10-name",
        "firstname": "test-user-to-bulk-create-10-firstname",
        "lastname": "test-user-to-bulk-create-10-lastname",
        "email": "test-user-to-bulk-create-10-email@canopsis.net",
        "roles": [
          "test-role-to-user-edit-1"
        ],
        "ui_language": "fr",
        "ui_groups_navigation_type": "top-bar",
        "enable": true,
        "defaultview": "test-view-to-edit-user",
        "ui_theme": "test_theme_to_pick_1",
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
          "roles": [
            "test-role-to-user-edit-2",
            "test-role-to-user-edit-1"
          ],
          "ui_language": "fr",
          "ui_theme": "canopsis",
          "ui_groups_navigation_type": "top-bar",
          "enable": true,
          "defaultview": "test-view-to-edit-user",
          "password": "test-password"
        }
      },
      {
        "status": 400,
        "item": {
          "name": "test-user-to-bulk-create-1-name",
          "firstname": "test-user-to-bulk-create-1-firstname",
          "lastname": "test-user-to-bulk-create-1-lastname",
          "email": "test-user-to-bulk-create-1-email@canopsis.net",
          "roles": [
            "test-role-to-user-edit-1"
          ],
          "ui_language": "fr",
          "ui_theme": "canopsis",
          "ui_groups_navigation_type": "top-bar",
          "enable": true,
          "defaultview": "test-view-to-edit-user",
          "password": "test-password"
        },
        "errors": {
          "name": "Name already exists."
        }
      },
      {
        "status": 400,
        "item": {
          "roles": [
            "test-role-to-user-edit-1",
            "not-exist"
          ],
          "defaultview": "not-exist"
        },
        "errors": {
          "defaultview": "DefaultView doesn't exist.",
          "email": "Email is missing.",
          "enable": "IsEnabled is missing.",
          "name": "Name is missing.",
          "password": "Password is missing.",
          "roles": "Roles doesn't exist."
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
          "roles": [
            "test-role-to-user-edit-1"
          ],
          "ui_language": "fr",
          "ui_theme": "canopsis_dark",
          "ui_groups_navigation_type": "top-bar",
          "enable": true,
          "defaultview": "test-view-to-edit-user",
          "password": "test-password"
        }
      },
      {
        "id": "test-user-to-bulk-create-3-name",
        "status": 200,
        "item": {
          "name": "test-user-to-bulk-create-3-name",
          "firstname": "test-user-to-bulk-create-3-firstname",
          "lastname": "test-user-to-bulk-create-3-lastname",
          "email": "test-user-to-bulk-create-3-email@canopsis.net",
          "roles": [
            "test-role-to-user-edit-1"
          ],
          "ui_language": "fr",
          "ui_groups_navigation_type": "top-bar",
          "enable": true,
          "defaultview": "test-view-to-edit-user",
          "source": "saml",
          "external_id": "test-id"
        }
      },
      {
        "status": 400,
        "item": {
          "name": "test-user-to-bulk-create-4-name",
          "firstname": "test-user-to-bulk-create-4-firstname",
          "lastname": "test-user-to-bulk-create-4-lastname",
          "email": "test-user-to-bulk-create-4-email@canopsis.net",
          "roles": [
            "test-role-to-user-edit-1"
          ],
          "ui_language": "fr",
          "ui_groups_navigation_type": "top-bar",
          "enable": true,
          "defaultview": "test-view-to-edit-user",
          "password": "test-password",
          "source": "saml",
          "external_id": "test-id"
        },
        "errors": {
          "source": "Can't be present both Source and Password."
        }
      },
      {
        "status": 400,
        "item": {
          "name": "test-user-to-bulk-create-5-name",
          "firstname": "test-user-to-bulk-create-5-firstname",
          "lastname": "test-user-to-bulk-create-5-lastname",
          "email": "test-user-to-bulk-create-5-email@canopsis.net",
          "roles": [
            "test-role-to-user-edit-1"
          ],
          "ui_language": "fr",
          "ui_groups_navigation_type": "top-bar",
          "enable": true,
          "defaultview": "test-view-to-edit-user",
          "source": "saml"
        },
        "errors": {
          "external_id": "ExternalID is required when Source is present."
        }
      },
      {
        "status": 400,
        "item": {
          "name": "test-user-to-bulk-create-6-name",
          "firstname": "test-user-to-bulk-create-6-firstname",
          "lastname": "test-user-to-bulk-create-6-lastname",
          "email": "test-user-to-bulk-create-6-email@canopsis.net",
          "roles": [
            "test-role-to-user-edit-1"
          ],
          "ui_language": "fr",
          "ui_groups_navigation_type": "top-bar",
          "enable": true,
          "defaultview": "test-view-to-edit-user",
          "source": "some",
          "external_id": "test-id"
        },
        "errors": {
          "source": "Source must be one of [ldap cas saml] or empty."
        }
      },
      {
        "status": 400,
        "item": {
          "name": "test-user-to-bulk-create-7-name",
          "firstname": "test-user-to-bulk-create-7-firstname",
          "lastname": "test-user-to-bulk-create-7-lastname",
          "email": "test-user-to-bulk-create-7-email@canopsis.net",
          "roles": [
            "test-role-to-user-edit-1"
          ],
          "ui_language": "fr",
          "ui_groups_navigation_type": "top-bar",
          "enable": true,
          "defaultview": "test-view-to-edit-user"
        },
        "errors": {
          "password": "Password is missing."
        }
      },
      {
        "status": 400,
        "item": {
          "name": "test-user-to-bulk-create-8-name",
          "firstname": "test-user-to-bulk-create-8-firstname",
          "lastname": "test-user-to-bulk-create-8-lastname",
          "email": "test-user-to-bulk-create-8-email@canopsis.net",
          "roles": [
            "test-role-to-user-edit-1"
          ],
          "ui_language": "fr",
          "ui_groups_navigation_type": "top-bar",
          "enable": true,
          "defaultview": "test-view-to-edit-user",
          "password": "123"
        },
        "errors": {
          "password": "Password should be 8 or more."
        }
      },
      {
        "status": 400,
        "item": {
          "name": "test-user-to-bulk-create-9-name",
          "firstname": "test-user-to-bulk-create-9-firstname",
          "lastname": "test-user-to-bulk-create-9-lastname",
          "email": "test-user-to-bulk-create-9-email@canopsis.net",
          "roles": [
            "test-role-to-user-edit-1"
          ],
          "ui_language": "fr",
          "ui_groups_navigation_type": "top-bar",
          "enable": true,
          "defaultview": "test-view-to-edit-user",
          "source": "saml",
          "external_id": "test-id",
          "ui_theme": "not found"
        },
        "errors": {
          "ui_theme": "UITheme doesn't exist."
        }
      },
      {
        "id": "test-user-to-bulk-create-10-name",
        "status": 200,
        "item": {
          "name": "test-user-to-bulk-create-10-name",
          "firstname": "test-user-to-bulk-create-10-firstname",
          "lastname": "test-user-to-bulk-create-10-lastname",
          "email": "test-user-to-bulk-create-10-email@canopsis.net",
          "roles": [
            "test-role-to-user-edit-1"
          ],
          "ui_language": "fr",
          "ui_groups_navigation_type": "top-bar",
          "enable": true,
          "defaultview": "test-view-to-edit-user",
          "ui_theme": "test_theme_to_pick_1",
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
          "roles": [
            {
              "_id": "test-role-to-user-edit-2",
              "name": "test-role-to-user-edit-2",
              "defaultview": {
                "_id": "test-view-to-edit-user",
                "title": "test-view-to-edit-user-title"
              }
            },
            {
              "_id": "test-role-to-user-edit-1",
              "name": "test-role-to-user-edit-1",
              "defaultview": {
                "_id": "test-view-to-edit-user",
                "title": "test-view-to-edit-user-title"
              }
            }
          ],
          "source": "",
          "ui_groups_navigation_type": "top-bar",
          "ui_language": "fr",
          "ui_theme": {
            "name": "Canopsis",
            "colors": {
              "main": {
                "primary": "#2fab63",
                "secondary": "#2b3e4f",
                "accent": "#82b1ff",
                "error": "#ff5252",
                "info": "#2196f3",
                "success": "#4caf50",
                "warning": "#fb8c00",
                "background": "#ffffff",
                "active_color": "#000",
                "font_size": 2
              },
              "table": {
                "background": "#fff",
                "row_color": "#fff",
                "hover_row_color": "#eee"
              },
              "state": {
                "ok": "#00a65a",
                "minor": "#fcdc00",
                "major": "#ff9900",
                "critical": "#f56954"
              }
            }
          }
        },
        {
          "_id": "test-user-to-bulk-create-10-name",
          "defaultview": {
            "_id": "test-view-to-edit-user",
            "title": "test-view-to-edit-user-title"
          },
          "email": "test-user-to-bulk-create-10-email@canopsis.net",
          "enable": true,
          "external_id": "",
          "firstname": "test-user-to-bulk-create-10-firstname",
          "lastname": "test-user-to-bulk-create-10-lastname",
          "name": "test-user-to-bulk-create-10-name",
          "roles": [
            {
              "_id": "test-role-to-user-edit-1",
              "name": "test-role-to-user-edit-1",
              "defaultview": {
                "_id": "test-view-to-edit-user",
                "title": "test-view-to-edit-user-title"
              }
            }
          ],
          "source": "",
          "ui_groups_navigation_type": "top-bar",
          "ui_language": "fr",
          "ui_theme": {
            "name": "test_theme_to_pick_1",
            "colors": {
              "main": {
                "primary": "#AAAAAA",
                "secondary": "#AAAAAA",
                "accent": "#AAAAAA",
                "error": "#AAAAAA",
                "info": "#AAAAAA",
                "success": "#AAAAAA",
                "warning": "#AAAAAA",
                "background": "#AAAAAA",
                "active_color": "#AAAAAA",
                "font_size": 2
              },
              "table": {
                "background": "#AAAAAA",
                "row_color": "#AAAAAA",
                "shift_row_color": "#AAAAAA",
                "hover_row_color": "#AAAAAA"
              },
              "state": {
                "ok": "#AAAAAA",
                "minor": "#AAAAAA",
                "major": "#AAAAAA",
                "critical": "#AAAAAA"
              }
            }
          }
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
          "roles": [
            {
              "_id": "test-role-to-user-edit-1",
              "name": "test-role-to-user-edit-1",
              "defaultview": {
                "_id": "test-view-to-edit-user",
                "title": "test-view-to-edit-user-title"
              }
            }
          ],
          "source": "",
          "ui_groups_navigation_type": "top-bar",
          "ui_language": "fr",
          "ui_theme": {
            "name": "Canopsis dark",
            "colors": {
              "main": {
                "primary": "#2fab63",
                "secondary": "#2b3e4f",
                "accent": "#82b1ff",
                "error": "#ff8b8b",
                "info": "#2196f3",
                "success": "#4caf50",
                "warning": "#fb8c00",
                "background": "#303030",
                "active_color": "#fff",
                "font_size": 2
              },
              "table": {
                "background": "#424242",
                "row_color": "#424242",
                "hover_row_color": "#616161"
              },
              "state": {
                "ok": "#00a65a",
                "minor": "#fcdc00",
                "major": "#ff9900",
                "critical": "#f56954"
              }
            }
          }
        },
        {
          "_id": "test-user-to-bulk-create-3-name",
          "defaultview": {
            "_id": "test-view-to-edit-user",
            "title": "test-view-to-edit-user-title"
          },
          "email": "test-user-to-bulk-create-3-email@canopsis.net",
          "enable": true,
          "firstname": "test-user-to-bulk-create-3-firstname",
          "lastname": "test-user-to-bulk-create-3-lastname",
          "name": "test-user-to-bulk-create-3-name",
          "roles": [
            {
              "_id": "test-role-to-user-edit-1",
              "name": "test-role-to-user-edit-1",
              "defaultview": {
                "_id": "test-view-to-edit-user",
                "title": "test-view-to-edit-user-title"
              }
            }
          ],
          "ui_groups_navigation_type": "top-bar",
          "ui_language": "fr",
          "source": "saml",
          "external_id": "test-id",
          "ui_theme": {
            "name": "Canopsis",
            "colors": {
              "main": {
                "primary": "#2fab63",
                "secondary": "#2b3e4f",
                "accent": "#82b1ff",
                "error": "#ff5252",
                "info": "#2196f3",
                "success": "#4caf50",
                "warning": "#fb8c00",
                "background": "#ffffff",
                "active_color": "#000",
                "font_size": 2
              },
              "table": {
                "background": "#fff",
                "row_color": "#fff",
                "hover_row_color": "#eee"
              },
              "state": {
                "ok": "#00a65a",
                "minor": "#fcdc00",
                "major": "#ff9900",
                "critical": "#f56954"
              }
            }
          }          
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

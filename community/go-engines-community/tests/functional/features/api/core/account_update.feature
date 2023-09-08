Feature: Update an account
  I need to be able to update an account
  Only admin should be able to update an account

  @concurrent
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
    Then the response body should contain:
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
      "display_name": "test-user-to-account-update-1 test-user-to-account-update-1-firstname test-user-to-account-update-1-lastname test-user-to-account-update-1-email@canopsis.net",
      "roles": [
        {
          "_id": "test-role-to-account-update-1",
          "name": "test-role-to-account-update-1",
          "defaultview": null
        }
      ],
      "permissions": [],
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
      },
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

  @concurrent
  Scenario: given update request with password should auth user by base auth
    When I am test-role-to-account-update-2
    Then I do PUT /api/v4/account/me:
    """json
    {
      "password": "test-password-updated",
      "ui_language": "fr",
      "ui_theme": "canopsis",
      "ui_groups_navigation_type": "top-bar",
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

  @concurrent
  Scenario: given get request and no auth user should not allow access
    Then I do PUT /api/v4/account/me
    Then the response code should be 401

  @concurrent
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

  @concurrent
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

  @concurrent
  Scenario: given account update request with custom color theme should be ok
  When I am authenticated with username "test-user-to-pick-color-theme-1" and password "test"
  When I do PUT /api/v4/account/me:
  """json
  {
    "ui_language": "fr",
    "ui_groups_navigation_type": "top-bar",
    "defaultview": "test-view-to-edit-user",
    "ui_theme": "test_theme_to_pick_1"
  }
  """
  Then the response code should be 200
  Then the response body should contain:
  """json
  {
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
  }
  """

  @concurrent
  Scenario: given account update request with not found color theme should return error
  When I am authenticated with username "test-user-to-pick-color-theme-2" and password "test"
  When I do PUT /api/v4/account/me:
  """json
  {
    "ui_language": "fr",
    "ui_groups_navigation_type": "top-bar",
    "defaultview": "test-view-to-edit-user",
    "ui_theme": "test_theme_to_pick_not_found"
  }
  """
  Then the response code should be 400
  Then the response body should contain:
  """json
  {
    "errors": {
      "ui_theme": "UITheme doesn't exist."
    }
  }
  """

  @concurrent
  Scenario: given account update request with empty ui_theme should return default theme
  When I am authenticated with username "test-user-to-pick-color-theme-3" and password "test"
  When I do PUT /api/v4/account/me:
  """json
  {
    "ui_language": "fr",
    "ui_groups_navigation_type": "top-bar",
    "defaultview": "test-view-to-edit-user",
    "ui_theme": ""
  }
  """
  Then the response code should be 200
  Then the response body should contain:
  """json
  {
    "ui_theme": {
      "name": "Canopsis"
    }
  }
  """

  @concurrent
  Scenario: given get user without ui_theme field should return default theme
  When I am authenticated with username "test-user-to-pick-color-theme-4" and password "test"
  When I do PUT /api/v4/account/me:
  """json
  {
    "ui_language": "fr",
    "ui_groups_navigation_type": "top-bar",
    "defaultview": "test-view-to-edit-user"
  }
  """
  Then the response code should be 200
  Then the response body should contain:
  """json
  {
    "ui_theme": {
      "name": "Canopsis"
    }
  }
  """

  @concurrent
  Scenario: given account update request and delete color theme, picked theme should be replaced to default
  When I am authenticated with username "test-user-to-pick-color-theme-5" and password "test"
  When I do PUT /api/v4/account/me:
  """json
  {
    "ui_language": "fr",
    "ui_groups_navigation_type": "top-bar",
    "defaultview": "test-view-to-edit-user",
    "ui_theme": "test_theme_to_pick_2"
  }
  """
  Then the response code should be 200
  Then the response body should contain:
  """json
  {
    "ui_theme": {
      "name": "test_theme_to_pick_2",
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
  }
  """
  When I am admin
  When I do DELETE /api/v4/color-themes/test_theme_to_pick_2
  Then the response code should be 204
  When I am authenticated with username "test-user-to-pick-color-theme-5" and password "test"
  When I do GET /api/v4/account/me
  Then the response code should be 200
  Then the response body should contain:
  """json
  {
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
  """

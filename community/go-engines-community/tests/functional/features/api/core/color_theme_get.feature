Feature: Get a color theme
  I need to be able to get a color theme
  Only admin should be able to get a color theme

  @concurrent
  Scenario: given get list request and no auth user should not allow access
    When I do GET /api/v4/color-themes
    Then the response code should be 401

  @concurrent
  Scenario: given get list request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/color-themes
    Then the response code should be 403

  @concurrent
  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/color-themes/test_theme_to_get_1
    Then the response code should be 401

  @concurrent
  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/color-themes/test_theme_to_get_1
    Then the response code should be 403

  @concurrent
  Scenario: given get list request should return ok
  When I am admin
  When I do GET /api/v4/color-themes?search=test_theme_to_get
  Then the response code should be 200
  Then the response body should contain:
  """json
  {
    "data": [
      {
        "name": "test_theme_to_get_1",
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
      },
      {
        "name": "test_theme_to_get_2",
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
    ],
    "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 2
    }
  }
  """

  @concurrent
  Scenario: given get request should return ok
  When I am admin
  When I do GET /api/v4/color-themes/test_theme_to_get_1
  Then the response code should be 200
  Then the response body should contain:
  """json
  {
    "name": "test_theme_to_get_1",
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
  """

  @concurrent
  Scenario: given get list with sort request should return ok
  When I am admin
  When I do GET /api/v4/color-themes?search=test_theme_to_get&sort_by=name&sort=asc
  Then the response code should be 200
  Then the response body should contain:
  """json
  {
    "data": [
      {
        "name": "test_theme_to_get_1"
      },
      {
        "name": "test_theme_to_get_2"
      }
    ],
    "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 2
    }
  }
  """
  When I do GET /api/v4/color-themes?search=test_theme_to_get&sort_by=name&sort=desc
  Then the response code should be 200
  Then the response body should contain:
  """json
  {
    "data": [
      {
        "name": "test_theme_to_get_2"
      },
      {
        "name": "test_theme_to_get_1"
      }
    ],
    "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 2
    }
  }
  """

  @concurrent
  Scenario: given get with not found theme should return error
  When I am admin
  When I do GET /api/v4/color-themes/test_theme_not_found
  Then the response code should be 404
  Then the response body should contain:
  """json
  {
    "error": "Not found"
  }
  """

  @concurrent
  Scenario: given get list with invalid sort request should return error
  When I am admin
  When I do GET /api/v4/color-themes?sort_by=unexpected&sort=asc
  Then the response code should be 400
  Then the response body should contain:
  """json
  {
    "errors": {
      "sort_by": "SortBy must be one of [name updated] or empty."
    }
  }
  """

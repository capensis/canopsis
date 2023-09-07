Feature: Update a color theme
  I need to be able to update a color theme
  Only admin should be able to update a color theme

  @concurrent
  Scenario: given update request and no auth user should not allow access
    When I do PUT /api/v4/color-themes/test_theme_to_update_1
    Then the response code should be 401

  @concurrent
  Scenario: given update request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/color-themes/test_theme_to_update_1
    Then the response code should be 403

  @concurrent
  Scenario: given update request should return ok
  When I am admin
  When I do PUT /api/v4/color-themes/test_theme_to_update_1:
  """json
  {
    "name": "test_theme_to_update_1_updated",
    "colors": {
      "main": {
        "primary": "#FFFFFF",
        "secondary": "#FFFFFF",
        "accent": "#FFFFFF",
        "error": "#FFFFFF",
        "info": "#FFFFFF",
        "success": "#FFFFFF",
        "warning": "#FFFFFF",
        "background": "#FFFFFF"
      },
      "table": {
        "background": "#FFFFFF",
        "active_color": "#FFFFFF",
        "row_color": "#FFFFFF",
        "shift_row_color": "#FFFFFF",
        "hover_row_color": "#FFFFFF"
      },
      "state": {
        "ok": "#FFFFFF",
        "minor": "#FFFFFF",
        "major": "#FFFFFF",
        "critical": "#FFFFFF"
      }
    }
  }
  """
  Then the response code should be 200
  Then the response body should contain:
  """json
  {
    "_id": "test_theme_to_update_1",
    "name": "test_theme_to_update_1_updated",
    "colors": {
      "main": {
        "primary": "#FFFFFF",
        "secondary": "#FFFFFF",
        "accent": "#FFFFFF",
        "error": "#FFFFFF",
        "info": "#FFFFFF",
        "success": "#FFFFFF",
        "warning": "#FFFFFF",
        "background": "#FFFFFF"
      },
      "table": {
        "background": "#FFFFFF",
        "active_color": "#FFFFFF",
        "row_color": "#FFFFFF",
        "shift_row_color": "#FFFFFF",
        "hover_row_color": "#FFFFFF"
      },
      "state": {
        "ok": "#FFFFFF",
        "minor": "#FFFFFF",
        "major": "#FFFFFF",
        "critical": "#FFFFFF"
      }
    }
  }
  """

  @concurrent
  Scenario: given update request with the existing name should return error
  When I am admin
  When I do PUT /api/v4/color-themes/test_theme_to_update_1:
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
        "background": "#AAAAAA"
      },
      "table": {
        "background": "#AAAAAA",
        "active_color": "#AAAAAA",
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
  Then the response code should be 400
  Then the response body should contain:
  """json
  {
    "errors": {
      "name": "Name already exists."
    }
  }
  """

  @concurrent
  Scenario: given update request without required fields should return error
  When I am admin
  When I do PUT /api/v4/color-themes/test_theme_to_update_1:
  """json
  {}
  """
  Then the response code should be 400
  Then the response body should contain:
  """json
  {
    "errors": {
      "colors.main.accent": "Accent is missing.",
      "colors.main.background": "Background is missing.",
      "colors.main.error": "Error is missing.",
      "colors.main.info": "Info is missing.",
      "colors.main.primary": "Primary is missing.",
      "colors.main.secondary": "Secondary is missing.",
      "colors.main.success": "Success is missing.",
      "colors.main.warning": "Warning is missing.",
      "colors.state.critical": "Critical is missing.",
      "colors.state.major": "Major is missing.",
      "colors.state.minor": "Minor is missing.",
      "colors.state.ok": "OK is missing.",
      "colors.table.active_color": "ActiveColor is missing.",
      "colors.table.background": "Background is missing.",
      "colors.table.row_color": "RowColor is missing.",
      "name": "Name is missing."
    }
  }
  """

  @concurrent
  Scenario: given update request with invalid color fields should return error
  When I am admin
  When I do PUT /api/v4/color-themes/test_theme_to_update_1:
  """json
  {
    "name": "test_theme_to_update_1",
    "colors": {
      "main": {
        "primary": "bad_color",
        "secondary": "bad_color",
        "accent": "bad_color",
        "error": "bad_color",
        "info": "bad_color",
        "success": "bad_color",
        "warning": "bad_color",
        "background": "bad_color"
      },
      "table": {
        "background": "bad_color",
        "active_color": "bad_color",
        "row_color": "bad_color",
        "shift_row_color": "bad_color",
        "hover_row_color": "bad_color"
      },
      "state": {
        "ok": "bad_color",
        "minor": "bad_color",
        "major": "bad_color",
        "critical": "bad_color"
      }
    }
  }
  """
  Then the response code should be 400
  Then the response body should contain:
  """json
  {
    "errors": {
      "colors.main.accent": "Accent is not valid.",
      "colors.main.background": "Background is not valid.",
      "colors.main.error": "Error is not valid.",
      "colors.main.info": "Info is not valid.",
      "colors.main.primary": "Primary is not valid.",
      "colors.main.secondary": "Secondary is not valid.",
      "colors.main.success": "Success is not valid.",
      "colors.main.warning": "Warning is not valid.",
      "colors.state.critical": "Critical is not valid.",
      "colors.state.major": "Major is not valid.",
      "colors.state.minor": "Minor is not valid.",
      "colors.state.ok": "OK is not valid.",
      "colors.table.active_color": "ActiveColor is not valid.",
      "colors.table.background": "Background is not valid.",
      "colors.table.hover_row_color": "HoverRowColor is not valid.",
      "colors.table.row_color": "RowColor is not valid.",
      "colors.table.shift_row_color": "ShiftRowColor is not valid."
    }
  }
  """

  @concurrent
  Scenario: given put request for not found theme should return error
  When I am admin
  When I do PUT /api/v4/color-themes/test_not_found:
  """json
  {
    "name": "test_not_found",
    "colors": {
      "main": {
        "primary": "#AAAAAA",
        "secondary": "#AAAAAA",
        "accent": "#AAAAAA",
        "error": "#AAAAAA",
        "info": "#AAAAAA",
        "success": "#AAAAAA",
        "warning": "#AAAAAA",
        "background": "#AAAAAA"
      },
      "table": {
        "background": "#AAAAAA",
        "active_color": "#AAAAAA",
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
  Then the response code should be 404
  Then the response body should contain:
  """json
  {
    "error": "Not found"
  }
  """

  @concurrent
  Scenario: given put request for default theme should return error
  When I am admin
  When I do PUT /api/v4/color-themes/canopsis:
  """json
  {
    "name": "test_not_found",
    "colors": {
      "main": {
        "primary": "#AAAAAA",
        "secondary": "#AAAAAA",
        "accent": "#AAAAAA",
        "error": "#AAAAAA",
        "info": "#AAAAAA",
        "success": "#AAAAAA",
        "warning": "#AAAAAA",
        "background": "#AAAAAA"
      },
      "table": {
        "background": "#AAAAAA",
        "active_color": "#AAAAAA",
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
  Then the response code should be 400
  Then the response body should contain:
  """json
  {
    "error": "can't modify or delete the default color theme"
  }
  """

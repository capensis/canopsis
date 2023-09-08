Feature: Create a color theme
  I need to be able to create a color theme
  Only admin should be able to create a color theme

  @concurrent
  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/color-themes
    Then the response code should be 401

  @concurrent
  Scenario: given create request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/color-themes
    Then the response code should be 403

  @concurrent
  Scenario: given create request should return ok
  When I am admin
  When I do POST /api/v4/color-themes:
  """json
  {
    "_id": "test_1",
    "name": "test_1",
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
        "font_size": 1
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
  Then the response code should be 201
  Then the response body should contain:
  """json
  {
    "_id": "test_1",
    "name": "test_1",
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
        "font_size": 1
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
    },
    "deletable": true
  }
  """

  @concurrent
  Scenario: given create request with the existing name or id should return error
  When I am admin
  When I do POST /api/v4/color-themes:
  """json
  {
    "_id": "test_2_1",
    "name": "test_2_1",
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
        "font_size": 1
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
  Then the response code should be 201
  When I do POST /api/v4/color-themes:
  """json
  {
    "_id": "test_2_2",
    "name": "test_2_1",
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
        "font_size": 1
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
  Then the response code should be 400
  Then the response body should contain:
  """json
  {
    "errors": {
      "name": "Name already exists."
    }
  }
  """
  When I do POST /api/v4/color-themes:
  """json
  {
    "_id": "test_2_1",
    "name": "test_2_2",
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
        "font_size": 1
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
  Then the response code should be 400
  Then the response body should contain:
  """json
  {
    "errors": {
      "_id": "ID already exists."
    }
  }
  """

  @concurrent
  Scenario: given create request without required fields should return error
  When I am admin
  When I do POST /api/v4/color-themes:
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
      "colors.main.active_color": "ActiveColor is missing.",
      "colors.main.font_size": "FontSize is missing.",
      "colors.state.critical": "Critical is missing.",
      "colors.state.major": "Major is missing.",
      "colors.state.minor": "Minor is missing.",
      "colors.state.ok": "OK is missing.",
      "colors.table.background": "Background is missing.",
      "colors.table.row_color": "RowColor is missing.",
      "name": "Name is missing."
    }
  }
  """

  @concurrent
  Scenario: given create request with invalid color and font fields should return error
  When I am admin
  When I do POST /api/v4/color-themes:
  """json
  {
    "_id": "test_3",
    "name": "test_3",
    "colors": {
      "main": {
        "primary": "bad_color",
        "secondary": "bad_color",
        "accent": "bad_color",
        "error": "bad_color",
        "info": "bad_color",
        "success": "bad_color",
        "warning": "bad_color",
        "background": "bad_color",
        "active_color": "bad_color",
        "font_size": 4
      },
      "table": {
        "background": "bad_color",
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
      "colors.main.active_color": "ActiveColor is not valid.",
      "colors.main.font_size": "FontSize must be one of [1 2 3].",
      "colors.state.critical": "Critical is not valid.",
      "colors.state.major": "Major is not valid.",
      "colors.state.minor": "Minor is not valid.",
      "colors.state.ok": "OK is not valid.",
      "colors.table.background": "Background is not valid.",
      "colors.table.hover_row_color": "HoverRowColor is not valid.",
      "colors.table.row_color": "RowColor is not valid.",
      "colors.table.shift_row_color": "ShiftRowColor is not valid."
    }
  }
  """

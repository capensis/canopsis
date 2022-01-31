Feature: Create a widget
  I need to be able to create a widget
  Only admin should be able to create a widget

  Scenario: given create request should return ok
    When I am admin
    When I do POST /api/v4/widgets:
    """json
    {
      "title": "test-widget-to-create-1-title",
      "tab": "test-tab-to-widget-edit",
      "type": "test-widget-to-create-1-type",
      "grid_parameters": {
        "desktop": {"x": 0, "y": 0}
      },
      "parameters": {
        "test-widget-to-create-1-param-str": "teststr",
        "test-widget-to-create-1-param-int": 2,
        "test-widget-to-create-1-param-bool": true,
        "test-widget-to-create-1-param-arr": ["teststr1", "teststr2"],
        "test-widget-to-create-1-param-map": {"testkey": "teststr"}
      }
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "title": "test-widget-to-create-1-title",
      "type": "test-widget-to-create-1-type",
      "grid_parameters": {
        "desktop": {"x": 0, "y": 0}
      },
      "parameters": {
        "test-widget-to-create-1-param-str": "teststr",
        "test-widget-to-create-1-param-int": 2,
        "test-widget-to-create-1-param-bool": true,
        "test-widget-to-create-1-param-arr": ["teststr1", "teststr2"],
        "test-widget-to-create-1-param-map": {"testkey": "teststr"}
      },
      "author": "root"
    }
    """
    When I do GET /api/v4/widgets/{{ .lastResponse._id}}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "title": "test-widget-to-create-1-title",
      "type": "test-widget-to-create-1-type",
      "grid_parameters": {
        "desktop": {"x": 0, "y": 0}
      },
      "parameters": {
        "test-widget-to-create-1-param-str": "teststr",
        "test-widget-to-create-1-param-int": 2,
        "test-widget-to-create-1-param-bool": true,
        "test-widget-to-create-1-param-arr": ["teststr1", "teststr2"],
        "test-widget-to-create-1-param-map": {"testkey": "teststr"}
      },
      "author": "root"
    }
    """

  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/widgets
    Then the response code should be 401

  Scenario: given create request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/widgets
    Then the response code should be 403

  Scenario: given invalid create request should return errors
    When I am admin
    When I do POST /api/v4/widgets:
    """json
    {
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "type": "Type is missing.",
        "tab": "Tab is missing."
      }
    }
    """

  Scenario: given Junit invalid create request should return errors
    When I am admin
    When I do POST /api/v4/widgets:
    """json
    {
      "title": "test-widget-to-create-2-title",
      "tab": "test-tab-to-widget-edit",
      "type": "Junit"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "parameters.directory": "Directory is missing."
      }
    }
    """
    When I do POST /api/v4/widgets:
    """json
    {
      "title": "test-widget-to-create-2-title",
      "tab": "test-tab-to-widget-edit",
      "type": "Junit",
      "parameters": {
        "is_api": true,
        "directory": "testdirectory",
        "screenshot_directories": ["testdirectory"],
        "video_directories": ["testdirectory"]
      }
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "parameters.directory": "Directory is not empty.",
        "parameters.screenshot_directories": "ScreenshotDirectories is not empty.",
        "parameters.video_directories": "VideoDirectories is not empty."
      }
    }
    """
    When I do POST /api/v4/widgets:
    """json
    {
      "title": "test-widget-to-create-2-title",
      "tab": "test-tab-to-widget-edit",
      "type": "Junit",
      "parameters": {
        "directory": "testdirectory",
        "video_filemask": "test",
        "screenshot_filemask": "test",
        "report_fileregexp": "(.*)"
      }
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "parameters.video_filemask": "VideoFilemask is not a valid file mask.",
        "parameters.screenshot_filemask": "ScreenshotFilemask is not a valid file mask.",
        "parameters.report_fileregexp": "ReportFileRegexp is invalid regexp."
      }
    }
    """

  Scenario: given create request and auth user without view permission should not allow access
    When I am admin
    When I do POST /api/v4/widgets:
    """json
    {
      "title": "test-widget-to-create-2-title",
      "type": "test-widget-to-create-2-type",
      "tab": "test-tab-to-widget-check-access"
    }
    """
    Then the response code should be 403

  Scenario: given create request with not exist tab should return not allow access error
    When I am admin
    When I do POST /api/v4/widgets:
    """json
    {
      "title": "test-widget-to-create-2-title",
      "type": "test-widget-to-create-2-type",
      "tab": "test-tab-not-found"
    }
    """
    Then the response code should be 403

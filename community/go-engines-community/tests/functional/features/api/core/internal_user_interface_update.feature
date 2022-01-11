Feature: update user interface
  I need to be able to update user interface
  Only admin should be able to update user interface
  Revert scenario changes to max_matched_items as fixtures loaded

  Scenario: PUT a valid user_interface but unauthorized
    When I do PUT /api/v4/internal/user_interface
    Then the response code should be 401

  Scenario: POST a valid user_interface but without permissions
    When I am noperms
    When I do POST /api/v4/internal/user_interface
    Then the response code should be 403

  Scenario: PUT a valid user_interface without timeout config
    When I am admin
    When I do PUT /api/v4/internal/user_interface:
    """
    {
      "language": "en",
      "footer": "Test footer",
      "app_title": "Canopsis Test",
      "login_page_description": "Test login",
      "max_matched_items": 10000,
      "check_count_request_timeout": 30
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "allow_change_severity_to_info": false,
      "app_title": "Canopsis Test",
      "footer": "Test footer",
      "language": "en",
      "login_page_description": "Test login",
      "max_matched_items": 10000,
      "check_count_request_timeout": 30,
      "popup_timeout": {
        "error": {
          "value": 3,
          "unit": "s"
        },
        "info": {
          "value": 3,
          "unit": "s"
        }
      }
    }
    """

  Scenario: PUT a valid user_interface with timeout config
    When I am admin
    When I do PUT /api/v4/internal/user_interface:
    """
    {
      "language": "en",
      "footer": "Test footer",
      "popup_timeout": {
        "error": {
          "value": 30,
          "unit": "s"
        }
      }
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "allow_change_severity_to_info": false,
      "app_title": "Canopsis Test",
      "footer": "Test footer",
      "language": "en",
      "login_page_description": "Test login",
      "popup_timeout": {
        "error": {
          "value": 30,
          "unit": "s"
        },
        "info": {
          "value": 3,
          "unit": "s"
        }
      }
    }
    """

  Scenario: POST an invalid user_interface
    When I am admin
    When I do POST /api/v4/internal/user_interface:
    """
    {
      "language": "vn",
      "footer": "Test footer"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
         "language": "Language must be one of [fr en] or empty."
      }
    }
    """

  Scenario: PUT a valid user_interface without max_matched_items or check_count_request_timeout should
    set those values to default
    When I am admin
    When I do PUT /api/v4/internal/user_interface:
    """
    {
      "language": "en",
      "footer": "Test footer",
      "app_title": "Canopsis Test",
      "login_page_description": "Test login",
      "max_matched_items": 100,
      "check_count_request_timeout": 100
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "allow_change_severity_to_info": false,
      "app_title": "Canopsis Test",
      "footer": "Test footer",
      "language": "en",
      "login_page_description": "Test login",
      "max_matched_items": 100,
      "check_count_request_timeout": 100,
      "popup_timeout": {
        "error": {
          "value": 3,
          "unit": "s"
        },
        "info": {
          "value": 3,
          "unit": "s"
        }
      }
    }
    """
    When I do PUT /api/v4/internal/user_interface:
    """
    {
      "language": "en",
      "footer": "Test footer",
      "app_title": "Canopsis Test",
      "login_page_description": "Test login"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "allow_change_severity_to_info": false,
      "app_title": "Canopsis Test",
      "footer": "Test footer",
      "language": "en",
      "login_page_description": "Test login",
      "max_matched_items": 10000,
      "check_count_request_timeout": 30,
      "popup_timeout": {
        "error": {
          "value": 3,
          "unit": "s"
        },
        "info": {
          "value": 3,
          "unit": "s"
        }
      }
    }
    """

  Scenario: PUT an invalid user_interface, max_matched_items and check_count_request_timeout should be >0
    When I am admin
    When I do PUT /api/v4/internal/user_interface:
    """
    {
      "language": "en",
      "footer": "Test footer",
      "app_title": "Canopsis Test",
      "login_page_description": "Test login",
      "max_matched_items": -1,
      "check_count_request_timeout": -1
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "check_count_request_timeout": "CheckCountRequestTimeout should be greater than 0.",
        "max_matched_items": "MaxMatchedItems should be greater than 0."
      }
    }
    """
    When I do PUT /api/v4/internal/user_interface:
    """
    {
      "language": "en",
      "footer": "Test footer",
      "app_title": "Canopsis Test",
      "login_page_description": "Test login",
      "max_matched_items": 0,
      "check_count_request_timeout": 0
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "check_count_request_timeout": "CheckCountRequestTimeout should be greater than 0.",
        "max_matched_items": "MaxMatchedItems should be greater than 0."
      }
    }
    """
    When I do PUT /api/v4/internal/user_interface:
    """
    {
      "max_matched_items": 4
    }
    """
    Then the response code should be 200
    When I wait 2s

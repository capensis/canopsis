Feature: update user interface
  I need to be able to update user interface
  Only admin should be able to update user interface

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
          "interval": 3,
          "unit": "s"
        },
        "info": {
          "interval": 3,
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
          "interval": 30,
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
          "interval": 30,
          "unit": "s"
        },
        "info": {
          "interval": 3,
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

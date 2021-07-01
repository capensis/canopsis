Feature: Get login information
  I need to be able to get application information
  Everyone should be able to get this information

  Scenario: given guest user should be able to get login info
    When I do GET /api/v4/internal/login_info
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "login_config": {
        "casconfig": {
          "enable": false
        },
        "ldapconfig": {
          "enable": false
        },
         "saml2config": {
          "enable": false
        }
      }
    }
    """

  Scenario: given admin user should be able to get login info
    When I am admin
    When I do GET /api/v4/internal/login_info
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "login_config": {
        "casconfig": {
          "enable": false
        },
        "ldapconfig": {
          "enable": false
        },
         "saml2config": {
          "enable": false
        }
      },
      "edition": "cat",
      "stack": "go",
      "user_interface": {
        "allow_change_severity_to_info": false,
        "app_title": "Canopsis Test",
        "footer": "Test footer",
        "language": "en",
        "login_page_description": "Test login"
      },
      "version": "3.42.0"
    }
    """

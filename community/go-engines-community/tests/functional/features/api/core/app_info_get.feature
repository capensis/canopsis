Feature: Get application information
  I need to be able to get application information
  Only admin should be able to get this information

  Scenario: given get auth request should return application information
    When I am admin
    When I do GET /api/v4/app-info
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "allow_change_severity_to_info": false,
      "app_title": "Canopsis Test",
      "edition": "pro",
      "footer": "Test footer",
      "language": "en",
      "login_page_description": "Test login",
      "remediation": {
        "job_config_types": [
          {
            "auth_type": "bearer-token",
            "name": "awx"
          },
          {
            "auth_type": "basic-auth",
            "name": "jenkins"
          },
          {
            "auth_type": "header-token",
            "name": "rundeck"
          }
        ]
      },
      "stack": "go",
      "timezone": "Europe/Paris",
      "version": "development"
    }
    """

  Scenario: given get unauth request should return application information
    When I do GET /api/v4/app-info
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "allow_change_severity_to_info": false,
      "app_title": "Canopsis Test",
      "edition": "pro",
      "footer": "Test footer",
      "language": "en",
      "login_page_description": "Test login",
      "stack": "go",
      "version": "development"
    }
    """

  Scenario: given get request without permissions should return application information
    When I am noperms
    When I do GET /api/v4/app-info
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "allow_change_severity_to_info": false,
      "app_title": "Canopsis Test",
      "edition": "pro",
      "footer": "Test footer",
      "language": "en",
      "login_page_description": "Test login",
      "stack": "go",
      "version": "development"
    }
    """

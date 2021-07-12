Feature: Get application information
  I need to be able to get application information
  Only admin should be able to get this information

  Scenario: given get request should return application information
    When I am admin
    When I do GET /api/v4/internal/app_info
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "allow_change_severity_to_info": false,
      "app_title": "Canopsis Test",
      "edition": "cat",
      "footer": "Test footer",
      "language": "en",
      "login_page_description": "Test login",
      "stack": "go",
      "version": "3.42.0"
    }
    """

  Scenario: GET application information but unauthorized
    When I do GET /api/v4/internal/app_info
    Then the response code should be 401

  Scenario: GET application information but without permissions
    When I am noperms
    When I do GET /api/v4/internal/app_info
    Then the response code should be 403
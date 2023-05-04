Feature: Get application information
  I need to be able to get application information
  Only admin should be able to get this information

  @standalone
  Scenario: given get request should return application information
    When I do GET /api/v4/app-info
    Then the response code should be 200
    Then the response body should be:
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
            "name": "awx",
            "with_body": true,
            "with_query": false
          },
          {
            "auth_type": "basic-auth",
            "name": "jenkins",
            "with_body": false,
            "with_query": true
          },
          {
            "auth_type": "header-token",
            "name": "rundeck",
            "with_body": true,
            "with_query": false
          },
          {
            "auth_type": "header-token",
            "name": "vtom",
            "with_body": true,
            "with_query": false
          }
        ]
      },
      "stack": "go",
      "timezone": "Europe/Paris",
      "version": "development",
      "check_count_request_timeout": 30,
      "show_header_on_kiosk_mode": false,
      "file_upload_max_size": 314572800,
      "version_updated": null,
      "max_matched_items": 4,
      "login": {
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

Feature: Get a metrics settings
  I need to be able to get a metrics settings
  Only admin should be able to get a metrics settings

  Scenario: given get request should return metrics settings
    When I am admin
    When I do GET /api/v4/cat/metrics-settings
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "enabled_manual_instructions": true,
      "enabled_not_acked_metrics": true
    }
    """

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/cat/metrics-settings
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/cat/metrics-settings
    Then the response code should be 403

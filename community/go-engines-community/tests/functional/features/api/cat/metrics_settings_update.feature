Feature: Update a metrics settings
  I need to be able to update a metrics settings
  Only admin should be able to update a metrics settings

  Scenario: given update request should update metrics settings
    When I am admin
    When I do PUT /api/v4/cat/metrics-settings:
    """json
    {
      "enabled_manual_instructions": true,
      "enabled_not_acked_metrics": true
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "enabled_manual_instructions": true,
      "enabled_not_acked_metrics": true
    }
    """
    When I do GET /api/v4/cat/metrics-settings
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "enabled_manual_instructions": true,
      "enabled_not_acked_metrics": true
    }
    """

  Scenario: given update request and no auth user should not allow access
    When I do PUT /api/v4/cat/metrics-settings:
    """json
    {
      "enabled_manual_instructions": true,
      "enabled_not_acked_metrics": true
    }
    """
    Then the response code should be 401

  Scenario: given update request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/cat/metrics-settings:
    """json
    {
      "enabled_manual_instructions": true,
      "enabled_not_acked_metrics": true
    }
    """
    Then the response code should be 403

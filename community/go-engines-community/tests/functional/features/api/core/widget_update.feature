Feature: Update a widget
  I need to be able to update a widget
  Only admin should be able to update a widget

  Scenario: given update request should update widget
    When I am admin
    Then I do PUT /api/v4/widgets/test-widget-to-update:
    """json
    {
      "title": "test-widget-to-update-title-updated",
      "tab": "test-tab-to-widget-edit",
      "type": "test-widget-to-update-type",
      "grid_parameters": {
        "desktop": {"x": 0, "y": 0}
      },
      "parameters": {
        "test-widget-to-update-param-str": "teststr",
        "test-widget-to-update-param-int": 2,
        "test-widget-to-update-param-bool": true,
        "test-widget-to-update-param-arr": ["teststr1", "teststr2"],
        "test-widget-to-update-param-map": {"testkey": "teststr"}
      }
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-widget-to-update",
      "title": "test-widget-to-update-title-updated",
      "type": "test-widget-to-update-type",
      "grid_parameters": {
        "desktop": {"x": 0, "y": 0}
      },
      "parameters": {
        "test-widget-to-update-param-str": "teststr",
        "test-widget-to-update-param-int": 2,
        "test-widget-to-update-param-bool": true,
        "test-widget-to-update-param-arr": ["teststr1", "teststr2"],
        "test-widget-to-update-param-map": {"testkey": "teststr"}
      },
      "author": "root",
      "created": 1611229670
    }
    """

  Scenario: given update request and no auth user should not allow access
    When I do PUT /api/v4/widgets/test-widget-to-update
    Then the response code should be 401

  Scenario: given update request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/widgets/test-widget-to-update
    Then the response code should be 403

  Scenario: given update request and auth user without view permission should not allow access
    When I am admin
    When I do PUT /api/v4/widgets/test-widget-to-check-access:
    """json
    {
      "title": "test-widget-to-check-access-title",
      "type": "test-widget-to-check-access-type",
      "tab": "test-tab-to-widget-edit"
    }
    """
    Then the response code should be 403

  Scenario: given invalid update request should return errors
    When I am admin
    When I do PUT /api/v4/widgets/test-widget-to-update:
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

  Scenario: given update request with not exist id should return not allow access error
    When I am admin
    When I do PUT /api/v4/widgets/test-widget-not-found:
    """json
    {
      "title": "test-widget-not-found-title",
      "type": "test-widget-not-found-type",
      "tab": "test-tab-to-widget-edit"
    }
    """
    Then the response code should be 403

  Scenario: given update request and auth user without view permission should not allow access
    When I am admin
    When I do PUT /api/v4/widgets/test-widget-to-update:
    """json
    {
      "title": "test-widget-to-update-title",
      "type": "test-widget-to-update-type",
      "tab": "test-tab-to-widget-check-access"
    }
    """
    Then the response code should be 403

  Scenario: given update request with not exist tab should return not allow access error
    When I am admin
    When I do PUT /api/v4/widgets/test-widget-to-update:
    """json
    {
      "title": "test-widget-to-update-title",
      "type": "test-widget-to-update-type",
      "tab": "test-tab-not-found"
    }
    """
    Then the response code should be 403

Feature: Update a view tab
  I need to be able to update a view tab
  Only admin should be able to update a view tab

  Scenario: given update request should update tab
    When I am admin
    Then I do PUT /api/v4/view-tabs/test-tab-to-update:
    """json
    {
      "title": "test-tab-to-update-title-updated",
      "view": "test-view-to-tab-edit"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-tab-to-update",
      "title": "test-tab-to-update-title-updated",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "created": 1611229670
    }
    """

  Scenario: given update request and no auth user should not allow access
    When I do PUT /api/v4/view-tabs/test-tab-to-update
    Then the response code should be 401

  Scenario: given update request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/view-tabs/test-tab-to-update
    Then the response code should be 403

  Scenario: given update request and auth user without view permission should not allow access
    When I am admin
    When I do PUT /api/v4/view-tabs/test-tab-to-check-access:
    """json
    {
      "title": "test-tab-to-check-access-title",
      "view": "test-view-to-tab-edit"
    }
    """
    Then the response code should be 403

  Scenario: given invalid update request should return errors
    When I am admin
    When I do PUT /api/v4/view-tabs/test-tab-to-update:
    """json
    {
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "title": "Title is missing.",
        "view": "View is missing."
      }
    }
    """

  Scenario: given update request with not exist id should return not allow access error
    When I am admin
    When I do PUT /api/v4/view-tabs/test-tab-not-found:
    """json
    {
      "title": "test-tab-not-found-title",
      "view": "test-view-to-tab-edit"
    }
    """
    Then the response code should be 404

  Scenario: given update request and auth user without view permission should not allow access
    When I am admin
    When I do PUT /api/v4/view-tabs/test-tab-to-update:
    """json
    {
      "title": "test-tab-to-update-title",
      "view": "test-view-to-tab-check-access"
    }
    """
    Then the response code should be 403

  Scenario: given update request with not exist view should return not allow access error
    When I am admin
    When I do PUT /api/v4/view-tabs/test-tab-to-update:
    """json
    {
      "title": "test-tab-to-update-title",
      "view": "test-view-not-found"
    }
    """
    Then the response code should be 403

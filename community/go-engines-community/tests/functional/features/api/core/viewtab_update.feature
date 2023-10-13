Feature: Update a view tab
  I need to be able to update a view tab
  Only admin should be able to update a view tab

  @concurrent
  Scenario: given update request should update tab
    When I am admin
    Then I do PUT /api/v4/view-tabs/test-tab-to-update-1:
    """json
    {
      "title": "test-tab-to-update-1-title-updated"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-tab-to-update-1",
      "title": "test-tab-to-update-1-title-updated",
      "author": {
        "_id": "root",
        "name": "root"
      }
    }
    """
    Then I do GET /api/v4/view-tabs/test-tab-to-update-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-tab-to-update-1",
      "title": "test-tab-to-update-1-title-updated",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "is_private": false
    }
    """

  @concurrent
  Scenario: given update request and no auth user should not allow access
    When I do PUT /api/v4/view-tabs/test-tab-to-update-1
    Then the response code should be 401

  @concurrent
  Scenario: given update request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/view-tabs/test-tab-to-update-1
    Then the response code should be 403

  @concurrent
  Scenario: given update request and auth user without view permission should not allow access
    When I am admin
    When I do PUT /api/v4/view-tabs/test-tab-to-check-access:
    """json
    {
      "title": "test-tab-to-check-access-title"
    }
    """
    Then the response code should be 403

  @concurrent
  Scenario: given invalid update request should return errors
    When I am admin
    When I do PUT /api/v4/view-tabs/test-tab-to-update-1:
    """json
    {
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "title": "Title is missing."
      }
    }
    """

  @concurrent
  Scenario: given update request with not exist id should return not allow access error
    When I am admin
    When I do PUT /api/v4/view-tabs/test-tab-not-found:
    """json
    {
      "title": "test-tab-not-found-title"
    }
    """
    Then the response code should be 404

  @concurrent
  Scenario: given update owned private viewtab request should be ok
    When I am admin
    Then I do PUT /api/v4/view-tabs/test-private-tab-to-update-1:
    """json
    {
      "title": "test-private-tab-to-update-1-title-updated"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-private-tab-to-update-1",
      "title": "test-private-tab-to-update-1-title-updated",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "is_private": true
    }
    """

  @concurrent
  Scenario: given update not owned private viewtab request should not allow access
    When I am admin
    Then I do PUT /api/v4/view-tabs/test-private-tab-to-update-2:
    """json
    {
      "title": "test-private-tab-to-update-2-title-updated"
    }
    """
    Then the response code should be 403

  @concurrent
  Scenario: given update owned private viewtab request with api_private_view_groups
    but without api_view permissions should return filters should be ok
    When I am test-role-to-private-views-without-view-perm
    Then I do PUT /api/v4/view-tabs/test-private-tab-to-update-3:
    """json
    {
      "title": "test-private-tab-to-update-3-title-updated"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-private-tab-to-update-3",
      "title": "test-private-tab-to-update-3-title-updated",
      "author": {
        "_id": "test-user-to-private-views-without-view-perm",
        "name": "test-user-to-private-views-without-view-perm"
      },
      "is_private": true
    }
    """

  @concurrent
  Scenario: given update public viewtab request with api_private_view_groups
    but without api_view permissions should return filters should not allow access
    When I am test-role-to-private-views-without-view-perm
    Then I do PUT /api/v4/view-tabs/test-tab-to-update-1:
    """json
    {
      "title": "test-private-tab-to-update-3-title-updated"
    }
    """
    Then the response code should be 403

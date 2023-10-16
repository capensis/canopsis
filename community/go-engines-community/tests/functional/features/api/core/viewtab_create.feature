Feature: Create a view tab
  I need to be able to create a view tab
  Only admin should be able to create a view tab

  @concurrent
  Scenario: given create request should return ok
    When I am admin
    When I do POST /api/v4/view-tabs:
    """json
    {
      "title": "test-tab-to-create-1-title",
      "view": "test-view-to-tab-edit"
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "title": "test-tab-to-create-1-title",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "is_private": false
    }
    """
    When I do GET /api/v4/view-tabs/{{ .lastResponse._id}}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "title": "test-tab-to-create-1-title",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "is_private": false
    }
    """

  @concurrent
  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/view-tabs
    Then the response code should be 401

  @concurrent
  Scenario: given create request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/view-tabs
    Then the response code should be 403

  @concurrent
  Scenario: given invalid create request should return errors
    When I am admin
    When I do POST /api/v4/view-tabs:
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

  @concurrent
  Scenario: given create request and auth user without view permission should not allow access
    When I am admin
    When I do POST /api/v4/view-tabs:
    """json
    {
      "title": "test-tab-to-create-2-title",
      "view": "test-view-to-tab-check-access"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "view": "Can't modify a view."
      }
    }
    """

  @concurrent
  Scenario: given create request with not exist view should return error
    When I am admin
    When I do POST /api/v4/view-tabs:
    """json
    {
      "title": "test-tab-to-create-2-title",
      "view": "test-view-not-found"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "view": "View doesn't exist."
      }
    }
    """

  @concurrent
  Scenario: given create request with owned private view should create private tab
    When I am admin
    When I do POST /api/v4/view-tabs:
    """json
    {
      "title": "test-tab-to-create-3-title",
      "view": "test-private-view-to-tab-create-1"
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "title": "test-tab-to-create-3-title",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "is_private": true
    }
    """
    When I do GET /api/v4/view-tabs/{{ .lastResponse._id}}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "title": "test-tab-to-create-3-title",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "is_private": true
    }
    """
    When I am manager
    When I do GET /api/v4/view-tabs/{{ .lastResponse._id}}
    Then the response code should be 403

  @concurrent
  Scenario: given create request with not owned private view should create private tab
    When I am admin
    When I do POST /api/v4/view-tabs:
    """json
    {
      "title": "test-tab-to-create-3-title",
      "view": "test-private-view-to-tab-create-2"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "view": "View is private."
      }
    }
    """

  @concurrent
  Scenario: given create request with owned private view with api_private_view_groups
    but without api_view permissions should return filters should create private tab
    When I am test-role-to-private-views-without-view-perm
    When I do POST /api/v4/view-tabs:
    """json
    {
      "title": "test-tab-to-create-4-title",
      "view": "test-private-view-to-tab-create-3"
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "title": "test-tab-to-create-4-title",
      "author": {
        "_id": "test-user-to-private-views-without-view-perm",
        "name": "test-user-to-private-views-without-view-perm"
      },
      "is_private": true
    }
    """
    When I do GET /api/v4/view-tabs/{{ .lastResponse._id}}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "title": "test-tab-to-create-4-title",
      "author": {
        "_id": "test-user-to-private-views-without-view-perm",
        "name": "test-user-to-private-views-without-view-perm"
      },
      "is_private": true
    }
    """

  @concurrent
  Scenario: given create request with public view with api_private_view_groups
    but without api_view permissions should return filters should return error
    When I am test-role-to-private-views-without-view-perm
    When I do POST /api/v4/view-tabs:
    """json
    {
      "title": "test-tab-to-create-4-title",
      "view": "test-view-to-tab-edit"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "view": "Can't modify a view."
      }
    }
    """

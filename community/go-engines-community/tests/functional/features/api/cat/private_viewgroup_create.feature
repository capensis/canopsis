Feature: Create a private view group
  I need to be able to create a private view group
  Only user with permission should be able to create a private view group

  @concurrent
  Scenario: given create request should return ok
    When I am admin
    When I do POST /api/v4/cat/private-view-groups:
    """json
    {
      "title": "test-viewgroup-to-create-1-title"
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "title": "test-viewgroup-to-create-1-title",
      "author": {
        "_id": "root",
        "name": "root"
      }
    }
    """
    When I do GET /api/v4/view-groups/{{ .lastResponse._id}}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "title": "test-viewgroup-to-create-1-title",
      "author": {
        "_id": "root",
        "name": "root"
      }
    }
    """

  @concurrent
  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/cat/private-view-groups
    Then the response code should be 401

  @concurrent
  Scenario: given create request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/cat/private-view-groups
    Then the response code should be 403

  @concurrent
  Scenario: given invalid create request should return errors
    When I am admin
    When I do POST /api/v4/cat/private-view-groups:
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
  Scenario: given create request with the same title as in public view group should return ok
    When I am admin
    When I do POST /api/v4/cat/private-view-groups:
    """json
    {
      "title": "test-private-viewgroup-to-check-unique-title-1-title"
    }
    """
    Then the response code should be 201

  @concurrent
  Scenario: given create request with the same title as in private view group created by another user should return ok
    When I am admin
    When I do POST /api/v4/cat/private-view-groups:
    """json
    {
      "title": "test-private-viewgroup-to-check-unique-title-2-title"
    }
    """
    Then the response code should be 201

  @concurrent
  Scenario: given create request with the same title as in private view group created by the same user should return error
    When I am admin
    When I do POST /api/v4/cat/private-view-groups:
    """json
    {
      "title": "test-private-viewgroup-to-check-unique-title-3-title"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
          "title": "Title already exists."
      }
    }
    """

Feature: Update a private view group
  I need to be able to update a private view group
  Only user with permission should be able to update a private view group

  @concurrent
  Scenario: given update request should update view group
    When I am admin
    Then I do PUT /api/v4/cat/private-view-groups/test-private-viewgroup-to-update-1:
    """
    {
      "title": "test-private-viewgroup-to-update-1-title"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "_id": "test-private-viewgroup-to-update-1",
      "title": "test-private-viewgroup-to-update-1-title",
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
      "_id": "test-private-viewgroup-to-update-1",
      "title": "test-private-viewgroup-to-update-1-title",
      "author": {
        "_id": "root",
        "name": "root"
      }
    }
    """

  @concurrent
  Scenario: given get request and no auth user should not allow access
    When I do PUT /api/v4/cat/private-view-groups/test-private-viewgroup-to-update-1
    Then the response code should be 401

  @concurrent
  Scenario: given invalid update request should return errors
    When I am admin
    When I do PUT /api/v4/cat/private-view-groups/test-private-viewgroup-to-update-1:
    """
    {
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "title": "Title is missing."
      }
    }
    """

  @concurrent
  Scenario: given update request with not exist id should return not found error
    When I am admin
    When I do PUT /api/v4/cat/private-view-groups/test-viewgroup-not-found:
    """
    {
      "title": "test-viewgroup-not-found-title"
    }
    """
    Then the response code should be 404
    Then the response body should be:
    """
    {
      "error": "Not found"
    }
    """

  @concurrent
  Scenario: given update request with the same title as in public view group should return ok
    When I am admin
    When I do PUT /api/v4/cat/private-view-groups/test-private-viewgroup-to-update-2:
    """json
    {
      "title": "test-private-viewgroup-to-check-unique-title-7-title"
    }
    """
    Then the response code should be 200

  @concurrent
  Scenario: given update request with the same title as in private view group updated by another user should return ok
    When I am admin
    When I do PUT /api/v4/cat/private-view-groups/test-private-viewgroup-to-update-2:
    """json
    {
      "title": "test-private-viewgroup-to-check-unique-title-8-title"
    }
    """
    Then the response code should be 200

  @concurrent
  Scenario: given update request with the same title as in private view group updated by the same user should return error
    When I am admin
    When I do PUT /api/v4/cat/private-view-groups/test-private-viewgroup-to-update-2:
    """json
    {
      "title": "test-private-viewgroup-to-check-unique-title-4-title"
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

  @concurrent
  Scenario: given update request for not owned private view should not allow access
    When I am admin
    When I do PUT /api/v4/cat/private-view-groups/test-private-viewgroup-to-update-3:
    """json
    {
      "title": "test-private-viewgroup-to-update-3-title"
    }
    """
    Then the response code should be 403

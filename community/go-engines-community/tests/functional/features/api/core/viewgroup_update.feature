Feature: Update a view group
  I need to be able to update a view group
  Only admin should be able to update a view group

  @concurrent
  Scenario: given update request should update view group
    When I am admin
    Then I do PUT /api/v4/view-groups/test-viewgroup-to-update-1:
    """
    {
      "title": "test-viewgroup-to-update-1-title"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "_id": "test-viewgroup-to-update-1",
      "title": "test-viewgroup-to-update-1-title",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "created": 1611229670
    }
    """

  @concurrent
  Scenario: given get request and no auth user should not allow access
    When I do PUT /api/v4/view-groups/test-viewgroup-to-update-1
    Then the response code should be 401

  @concurrent
  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/view-groups/test-viewgroup-to-update-1
    Then the response code should be 403

  @concurrent
  Scenario: given invalid update request should return errors
    When I am admin
    When I do PUT /api/v4/view-groups/test-viewgroup-to-update-1:
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
    When I do PUT /api/v4/view-groups/test-viewgroup-not-found:
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
  Scenario: given update request with already exists title should return error
    When I am admin
    When I do PUT /api/v4/view-groups/test-viewgroup-to-update-1:
    """
    {
      "title": "test-viewgroup-to-check-unique-title-title"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
          "title": "Title already exists."
      }
    }
    """

  @concurrent
  Scenario: given update request with already exists title should return error
    When I am admin
    When I do PUT /api/v4/view-groups/test-viewgroup-to-update-2:
    """
    {
      "title": "test-private-viewgroup-to-check-unique-title-6-title"
    }
    """
    Then the response code should be 200

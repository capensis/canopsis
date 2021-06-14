Feature: Bulk update a view groups
  I need to be able to update multiple view groups
  Only admin should be able to update multiple view group

  Scenario: given bulk update request should update view group
    When I am admin
    Then I do PUT /api/v4/bulk/view-groups:
    """
    [
      {
        "_id": "test-viewgroup-to-bulk-update-1",
        "title": "test-viewgroup-to-bulk-update-1-title-updated"
      },
      {
        "_id": "test-viewgroup-to-bulk-update-2",
        "title": "test-viewgroup-to-bulk-update-2-title"
      }
    ]
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    [
      {
        "_id": "test-viewgroup-to-bulk-update-1",
        "title": "test-viewgroup-to-bulk-update-1-title-updated",
        "author": "root",
        "created": 1611229670
      },
      {
        "_id": "test-viewgroup-to-bulk-update-2",
        "title": "test-viewgroup-to-bulk-update-2-title",
        "author": "root",
        "created": 1611229670
      }
    ]
    """

  Scenario: given bulk update request with not exist ids should return not found error
    When I am admin
    When I do PUT /api/v4/bulk/view-groups:
    """
    [
      {
        "_id": "test-viewgroup-not-found-1",
        "title": "test-viewgroup-not-found-1-title"
      },
      {
        "_id": "test-viewgroup-to-bulk-update-1",
        "title": "test-viewgroup-to-bulk-update-1-title"
      }
    ]
    """
    Then the response code should be 404
    Then the response body should be:
    """
    {
      "error": "Not found"
    }
    """

  Scenario: given invalid bulk update request should return errors
    When I am admin
    When I do PUT /api/v4/bulk/view-groups:
    """
    [
      {}
    ]
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "0._id": "ID is missing.",
        "0.title": "Title is missing."
      }
    }
    """

  Scenario: given bulk update request with one valid item and one invalid item
    should return errors
    When I am admin
    When I do PUT /api/v4/bulk/view-groups:
    """
    [
      {
        "_id": "test-viewgroup-to-bulk-update-1",
        "title": "test-viewgroup-to-bulk-update-1-title"
      },
      {}
    ]
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "1._id": "ID is missing.",
        "1.title": "Title is missing."
      }
    }
    """

  Scenario: given bulk update request with already exists title should return error
    When I am admin
    When I do PUT /api/v4/bulk/view-groups:
    """
    [
      {
        "_id": "test-viewgroup-to-bulk-update-1",
        "title": "test-viewgroup-to-check-unique-title-title"
      }
    ]
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
          "0.title": "Title already exists."
      }
    }
    """

  Scenario: given bulk update request with multiple items with the same title
    should return error
    When I am admin
    Then I do PUT /api/v4/bulk/view-groups:
    """
    [
      {
        "_id": "test-viewgroup-to-bulk-update-1",
        "title": "test-viewgroup-to-bulk-update-1-title-same-updated"
      },
      {
        "_id": "test-viewgroup-to-bulk-update-2",
        "title": "test-viewgroup-to-bulk-update-1-title-same-updated"
      },
      {
        "_id": "test-viewgroup-to-bulk-update-3",
        "title": "test-viewgroup-to-bulk-update-1-title-same-updated"
      }
    ]
    """
    Then the response code should be 400
    """
    {
      "errors": {
          "1.title": "Title already exists.",
          "2.title": "Title already exists."
      }
    }
    """

  Scenario: given bulk update request with multiple items with the same id
    should return error
    When I am admin
    Then I do PUT /api/v4/bulk/view-groups:
    """
    [
      {
        "_id": "test-viewgroup-to-bulk-update-1",
        "title": "test-viewgroup-to-bulk-update-1-title-same-updated"
      },
      {
        "_id": "test-viewgroup-to-bulk-update-1",
        "title": "test-viewgroup-to-bulk-update-1-title-same-updated"
      },
      {
        "_id": "test-viewgroup-to-bulk-update-1",
        "title": "test-viewgroup-to-bulk-update-1-title-same-updated"
      }
    ]
    """
    Then the response code should be 400
    """
    {
      "errors": {
          "1._id": "ID already exists.",
          "2._id": "ID already exists."
      }
    }
    """

  Scenario: given bulk update request and no auth user should not allow access
    When I do PUT /api/v4/bulk/view-groups
    Then the response code should be 401

  Scenario: given bulk update request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/bulk/view-groups
    Then the response code should be 403
Feature: Bulk create view groups
  I need to be able to create multiple view group
  Only admin should be able to create multiple view group

  Scenario: given bulk create request should return ok
    When I am admin
    When I do POST /api/v4/bulk/view-groups:
    """
    [
      {
        "title": "test-viewgroup-to-bulk-create-1-title"
      },
      {
        "title": "test-viewgroup-to-bulk-create-2-title"
      }
    ]
    """
    Then the response code should be 201
    Then the response body should contain:
    """
    [
      {
        "title": "test-viewgroup-to-bulk-create-1-title",
        "author": "root"
      },
      {
        "title": "test-viewgroup-to-bulk-create-2-title",
        "author": "root"
      }
    ]
    """

  Scenario: given bulk create request should return ok to get request
    When I am admin
    When I do POST /api/v4/bulk/view-groups:
    """
    [
      {
        "title": "test-viewgroup-to-bulk-create-3-title"
      }
    ]
    """
    When I do GET /api/v4/view-groups/{{ (index .lastResponse 0)._id}}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "title": "test-viewgroup-to-bulk-create-3-title",
      "author": "root"
    }
    """

  Scenario: given invalid bulk create request should return errors
    When I am admin
    When I do POST /api/v4/bulk/view-groups:
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
        "0.title": "Title is missing."
      }
    }
    """

  Scenario: given bulk create request with one invalid and one valid data
    should return errors
    When I am admin
    When I do POST /api/v4/bulk/view-groups:
    """
    [
      {
        "title": "test-viewgroup-to-bulk-create-4-title"
      },
      {}
    ]
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "1.title": "Title is missing."
      }
    }
    """

  Scenario: given bulk create request with already exists title should return error
    When I am admin
    When I do POST /api/v4/bulk/view-groups:
    """
    [
      {
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

  Scenario: given bulk create request with multiple items with the same title
    should return error
    When I am admin
    When I do POST /api/v4/bulk/view-groups:
    """
    [
      {
        "title": "test-viewgroup-to-bulk-create-6-title"
      },
      {
        "title": "test-viewgroup-to-bulk-create-6-title"
      },
      {
        "title": "test-viewgroup-to-bulk-create-6-title"
      }
    ]
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "1.title": "Title already exists.",
        "2.title": "Title already exists."
      }
    }
    """

  Scenario: given bulk create request and no auth user should not allow access
    When I do POST /api/v4/bulk/view-groups
    Then the response code should be 401

  Scenario: given bulk create request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/bulk/view-groups
    Then the response code should be 403

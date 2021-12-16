Feature: Create a widget filter
  I need to be able to create a widget filter
  Only admin should be able to create a widget filter

  Scenario: given create request should return ok
    When I am admin
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-to-create-1-title",
      "widget": "test-widget-to-widgetfilter-edit",
      "query": "{\"test\":\"test\"}",
      "personal": false
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "title": "test-widgetfilter-to-create-1-title",
      "query": "{\"test\":\"test\"}",
      "author": "root"
    }
    """
    When I do GET /api/v4/widget-filters/{{ .lastResponse._id}}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "title": "test-widgetfilter-to-create-1-title",
      "query": "{\"test\":\"test\"}",
      "author": "root"
    }
    """

  Scenario: given personal create request should return ok
    When I am admin
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-to-create-2-title",
      "widget": "test-widget-to-widgetfilter-edit",
      "query": "{\"test\":\"test\"}",
      "personal": true
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "title": "test-widgetfilter-to-create-2-title",
      "query": "{\"test\":\"test\"}",
      "author": "root"
    }
    """
    When I save response filterID={{ .lastResponse._id }}
    When I do GET /api/v4/widget-filters/{{ .filterID }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "title": "test-widgetfilter-to-create-2-title",
      "query": "{\"test\":\"test\"}",
      "author": "root"
    }
    """
    When I am test-role-to-widgetfiler-edit
    When I do GET /api/v4/widget-filters/{{ .filterID }}
    Then the response code should be 404

  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/widget-filters
    Then the response code should be 401

  Scenario: given create request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/widget-filters
    Then the response code should be 403

  Scenario: given invalid create request should return errors
    When I am admin
    When I do POST /api/v4/widget-filters:
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
        "query": "Query is missing.",
        "widget": "Widget is missing.",
        "personal": "Personal is missing."
      }
    }
    """
    When I do POST /api/v4/widget-filters:
    """json
    {
      "query": "invalid"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "query": "Query is not valid mongo query."
      }
    }
    """

  Scenario: given create request and auth user without view permission should not allow access
    When I am admin
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-to-create-3-title",
      "query": "{\"test\":\"test\"}",
      "widget": "test-widget-to-check-access",
      "personal": false
    }
    """
    Then the response code should be 403

  Scenario: given create request with not exist tab should return not allow access error
    When I am admin
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-to-create-3-title",
      "query": "{\"test\":\"test\"}",
      "widget": "test-widget-not-found",
      "personal": false
    }
    """
    Then the response code should be 403

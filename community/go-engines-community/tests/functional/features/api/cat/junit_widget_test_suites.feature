Feature: get test suites list for widget
  I need to be able to get test suites list for widget
  Only admin should be able to get test suites list for widget

  Scenario: GET unauthorized
    When I do GET /api/v4/cat/junit/test-suites-widget/test-widget-to-junit-test-suites-get-2
    Then the response code should be 401

  Scenario: GET without permissions
    When I am noperms
    When I do GET /api/v4/cat/junit/test-suites-widget/test-widget-to-junit-test-suites-get-2
    Then the response code should be 403

  Scenario: GET success
    When I am admin
    When I do GET /api/v4/cat/junit/test-suites-widget/test-widget-to-junit-test-suites-get-2
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "1026ae1f-a18b-4216-8427-95d1be75b365",
          "test_suite_id": "fdebd370-3178-4c9e-97aa-a26a79dca770",
          "name": "05-widget.alarms-list.10.edit-mode-advanced-settings",
          "entity_id": "05-widget.alarms-list.10.edit-mode-advanced-settings.10.edit-mode-advanced-settings.xml",
          "errors": 1,
          "failures": 0,
          "total": 16,
          "skipped": 15,
          "state": 2,
          "time": 0.5,
          "timestamp": 1614782420,
          "created": 1614782420,
          "mini_chart": [
            0,
            0,
            0,
            0,
            0.5
          ]
        },
        {
          "_id": "5b4461d4-cfea-41dd-97cd-a5b15c34c1e3",
          "test_suite_id": "7e5d2a21-be12-47d1-a286-0b46b6b2b99b",
          "name": "noveo.app.PublishMissionTest",
          "entity_id": "noveo.app.PublishMissionTest.TEST-noveo.app.PublishMissionTest.xml",
          "errors": 0,
          "failures": 1,
          "state": 3,
          "skipped": 0,
          "total": 1,
          "time": 0.99,
          "timestamp": 1614782435,
          "created": 1614782435,
          "mini_chart": [
            0,
            0,
            0.49,
            0.19,
            0.99
          ]
        }
      ],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 2
      }
    }
    """

  Scenario: GET when widget doesn't contain test suites
    When I am admin
    When I do GET /api/v4/cat/junit/test-suites-widget/test-widget-to-junit-test-suites-get-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 0
      }
    }
    """

  Scenario: GET widget not found
    When I am admin
    When I do GET /api/v4/cat/junit/test-suites-widget/no-widget
    Then the response code should be 404

  Scenario: given updated widget should clear test suites
    When I am admin
    When I do GET /api/v4/cat/junit/test-suites-widget/test-widget-to-junit-edit-dir-param
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 1
      }
    }
    """
    When I do PUT /api/v4/widgets/test-widget-to-junit-edit-dir-param:
    """json
    {
      "title": "test-widget-to-junit-edit-dir-param-title",
      "type": "Junit",
      "tab": "test-tab-to-junit-edit-dir-param",
      "grid_parameters": {},
      "parameters": {
        "directory": "test-widget-to-junit-edit-dir-param-param-dir"
      }
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/cat/junit/test-suites-widget/test-widget-to-junit-edit-dir-param
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 1
      }
    }
    """
    When I do PUT /api/v4/widgets/test-widget-to-junit-edit-dir-param:
    """json
    {
      "title": "test-widget-to-junit-edit-dir-param-title",
      "type": "Junit",
      "tab": "test-tab-to-junit-edit-dir-param",
      "grid_parameters": {},
      "parameters": {
        "directory": "test-widget-to-junit-edit-dir-param-param-dir-updated"
      }
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/cat/junit/test-suites-widget/test-widget-to-junit-edit-dir-param
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [],
      "meta": {
        "page": 1,
        "per_page": 10,
        "page_count": 1,
        "total_count": 0
      }
    }
    """

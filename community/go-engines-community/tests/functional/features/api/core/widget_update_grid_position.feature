Feature: Update widget grid positions
  I need to be able to widget grid positions
  Only admin should be able to widget grid positions

  Scenario: given update request should return ok
    When I am test-role-to-widget-update-grid-position
    When I do POST /api/v4/widgets:
    """json
    {
      "title": "test-widget-to-update-grid-position-1-1-title",
      "type": "test-widget-to-update-grid-position-1-1-type",
      "tab": "test-tab-to-widget-update-grid-position-1",
      "grid_parameters": {
        "desktop": {"x": 0, "y": 0}
      }
    }
    """
    Then the response code should be 201
    When I save response widget1={{ .lastResponse._id }}
    When I do POST /api/v4/widgets:
    """json
    {
      "title": "test-widget-to-update-grid-position-1-2-title",
      "type": "test-widget-to-update-grid-position-1-2-type",
      "tab": "test-tab-to-widget-update-grid-position-1",
      "grid_parameters": {
        "desktop": {"x": 1, "y": 0}
      }
    }
    """
    Then the response code should be 201
    When I save response widget2={{ .lastResponse._id }}
    When I do POST /api/v4/widgets:
    """json
    {
      "title": "test-widget-to-update-grid-position-1-3-title",
      "type": "test-widget-to-update-grid-position-1-3-type",
      "tab": "test-tab-to-widget-update-grid-position-1",
      "grid_parameters": {
        "desktop": {"x": 0, "y": 1}
      }
    }
    """
    Then the response code should be 201
    When I save response widget3={{ .lastResponse._id }}
    # Test created positions
    When I do GET /api/v4/views/test-view-to-widget-update-grid-position-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "tabs": [
        {
          "widgets": [
            { "_id": "{{ .widget1 }}" },
            { "_id": "{{ .widget2 }}" },
            { "_id": "{{ .widget3 }}" }
          ]
        }
      ]
    }
    """
    # Test updated positions
    When I do PUT /api/v4/widget-grid-positions:
    """json
    [
      {
        "_id": "{{ .widget3 }}",
        "grid_parameters": {
          "desktop": {"x": 0, "y": 0}
        }
      },
      {
        "_id": "{{ .widget1 }}",
        "grid_parameters": {
          "desktop": {"x": 0, "y": 1}
        }
      },
      {
        "_id": "{{ .widget2 }}",
        "grid_parameters": {
          "desktop": {"x": 1, "y": 1}
        }
      }
    ]
    """
    Then the response code should be 204
    When I do GET /api/v4/views/test-view-to-widget-update-grid-position-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "tabs": [
        {
          "widgets": [
            { "_id": "{{ .widget3 }}" },
            { "_id": "{{ .widget1 }}" },
            { "_id": "{{ .widget2 }}" }
          ]
        }
      ]
    }
    """

  Scenario: given update request and no auth user should not allow access
    When I do PUT /api/v4/widget-grid-positions
    Then the response code should be 401

  Scenario: given update request and auth user without view permission should not allow access
    When I am noperms
    When I do PUT /api/v4/widget-grid-positions
    Then the response code should be 403

  Scenario: given invalid request should return error
    When I am test-role-to-widget-update-grid-position
    When I do PUT /api/v4/widget-grid-positions:
    """json
    []
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "items": "Items should not be blank."
      }
    }
    """
    When I do PUT /api/v4/widget-grid-positions:
    """json
    [
      {
        "_id": "test-widget-not-exist"
      }
    ]
    """
    Then the response code should be 403
    When I do PUT /api/v4/widget-grid-positions:
    """json
    [
      {
        "_id": "test-widget-to-update-grid-position-2"
      },
      {
        "_id": "test-widget-not-exist"
      }
    ]
    """
    Then the response code should be 403
    When I do PUT /api/v4/widget-grid-positions:
    """json
    [
      {
        "_id": "test-widget-to-update-grid-position-2"
      },
      {
        "_id": "test-widget-to-update-grid-position-3-1"
      },
      {
        "_id": "test-widget-to-update-grid-position-3-2"
      }
    ]
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "error": "widgets are related to different view tabs"
    }
    """
    When I do PUT /api/v4/widget-grid-positions:
    """json
    [
      {
        "_id": "test-widget-to-update-grid-position-3-1"
      }
    ]
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "error": "widgets are missing"
    }
    """
    When I do PUT /api/v4/widget-grid-positions:
    """json
    [
      {
        "_id": "test-widget-to-update-grid-position-2"
      },
      {
        "_id": "test-widget-to-update-grid-position-4"
      }
    ]
    """
    Then the response code should be 403

Feature: Update view tab positions
  I need to be able to view tab positions
  Only admin should be able to view tab positions

  Scenario: given update request should return ok
    When I am test-role-to-tab-update-position
    When I do POST /api/v4/view-tabs:
    """json
    {
      "title": "test-tab-to-update-position-1-1-title",
      "view": "test-view-to-tab-update-position-1"
    }
    """
    Then the response code should be 201
    When I save response tab1={{ .lastResponse._id }}
    When I do POST /api/v4/view-tabs:
    """json
    {
      "title": "test-tab-to-update-position-1-2-title",
      "view": "test-view-to-tab-update-position-1"
    }
    """
    Then the response code should be 201
    When I save response tab2={{ .lastResponse._id }}
    When I do POST /api/v4/view-tabs:
    """json
    {
      "title": "test-tab-to-update-position-1-3-title",
      "view": "test-view-to-tab-update-position-1"
    }
    """
    Then the response code should be 201
    When I save response tab3={{ .lastResponse._id }}
    # Test created positions
    When I do GET /api/v4/views/test-view-to-tab-update-position-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "tabs": [
        { "_id": "{{ .tab1 }}" },
        { "_id": "{{ .tab2 }}" },
        { "_id": "{{ .tab3 }}" }
      ]
    }
    """
    # Test updated positions
    When I do PUT /api/v4/view-tab-positions:
    """json
    [
      "{{ .tab3 }}",
      "{{ .tab1 }}",
      "{{ .tab2 }}"
    ]
    """
    Then the response code should be 204
    When I do GET /api/v4/views/test-view-to-tab-update-position-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "tabs": [
        { "_id": "{{ .tab3 }}" },
        { "_id": "{{ .tab1 }}" },
        { "_id": "{{ .tab2 }}" }
      ]
    }
    """

  Scenario: given update request and no auth user should not allow access
    When I do PUT /api/v4/view-tab-positions
    Then the response code should be 401

  Scenario: given update request and auth user without view permission should not allow access
    When I am noperms
    When I do PUT /api/v4/view-tab-positions
    Then the response code should be 403

  Scenario: given invalid request should return error
    When I am test-role-to-tab-update-position
    When I do PUT /api/v4/view-tab-positions:
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
    When I do PUT /api/v4/view-tab-positions:
    """json
    [
      "test-tab-not-exist"
    ]
    """
    Then the response code should be 404
    When I do PUT /api/v4/view-tab-positions:
    """json
    [
      "test-tab-to-update-position-2",
      "test-tab-not-exist"
    ]
    """
    Then the response code should be 404
    When I do PUT /api/v4/view-tab-positions:
    """json
    [
      "test-tab-to-update-position-2",
      "test-tab-to-update-position-2"
    ]
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "items": "Items already exists."
      }
    }
    """
    When I do PUT /api/v4/view-tab-positions:
    """json
    [
      "test-tab-to-update-position-2",
      "test-tab-to-update-position-3-1",
      "test-tab-to-update-position-3-2"
    ]
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "error": "view tabs are related to different views"
    }
    """
    When I do PUT /api/v4/view-tab-positions:
    """json
    [
      "test-tab-to-update-position-3-1"
    ]
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "error": "view tabs are missing"
    }
    """
    When I do PUT /api/v4/view-tab-positions:
    """json
    [
      "test-tab-to-update-position-2",
      "test-tab-to-update-position-4"
    ]
    """
    Then the response code should be 403

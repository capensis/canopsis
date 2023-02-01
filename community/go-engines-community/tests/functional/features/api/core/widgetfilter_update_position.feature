Feature: Update widget filter positions
  I need to be able to widget filter positions
  Only admin should be able to widget filter positions

  Scenario: given update private filters request should return ok
    When I am test-role-to-filter-update-position-1
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-to-update-position-1-1-title",
      "widget": "test-widget-to-filter-update-position-1",
      "is_private": true,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-widgetfilter-to-update-position-1-1-pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response filter1={{ .lastResponse._id }}
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-to-update-position-1-2-title",
      "widget": "test-widget-to-filter-update-position-1",
      "is_private": true,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-widgetfilter-to-update-position-1-2-pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response filter2={{ .lastResponse._id }}
    When I do POST /api/v4/widget-filters:
    """json
    {
      "title": "test-widgetfilter-to-update-position-1-3-title",
      "widget": "test-widget-to-filter-update-position-1",
      "is_private": true,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-widgetfilter-to-update-position-1-3-pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I save response filter3={{ .lastResponse._id }}
    # Test created positions
    When I do GET /api/v4/user-preferences/test-widget-to-filter-update-position-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "filters": [
        { "_id": "{{ .filter1 }}" },
        { "_id": "{{ .filter2 }}" },
        { "_id": "{{ .filter3 }}" }
      ]
    }
    """
    # Test updated positions
    When I do PUT /api/v4/widget-filter-positions:
    """json
    [
      "{{ .filter3 }}",
      "{{ .filter1 }}",
      "{{ .filter2 }}"
    ]
    """
    Then the response code should be 204
    When I do GET /api/v4/user-preferences/test-widget-to-filter-update-position-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "filters": [
        { "_id": "{{ .filter3 }}" },
        { "_id": "{{ .filter1 }}" },
        { "_id": "{{ .filter2 }}" }
      ]
    }
    """

  Scenario: given update public filters request should return ok
    When I am test-role-to-filter-update-position-2
    When I do POST /api/v4/widgets:
    """json
    {
      "title": "test-widget-to-filter-update-position-2-title",
      "tab": "test-tab-to-filter-update-position",
      "type": "test-widget-to-filter-update-position-2-type",
      "filters": [
        {
          "title": "test-widgetfilter-to-update-position-2-1-title",
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-widgetfilter-to-update-position-2-1-pattern"
                }
              }
            ]
          ]
        },
        {
          "title": "test-widgetfilter-to-update-position-2-2-title",
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-widgetfilter-to-update-position-2-2-pattern"
                }
              }
            ]
          ]
        },
        {
          "title": "test-widgetfilter-to-update-position-2-3-title",
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-widgetfilter-to-update-position-2-3-pattern"
                }
              }
            ]
          ]
        }
      ]
    }
    """
    Then the response code should be 201
    When I save response widgetId={{ .lastResponse._id }}
    When I save response filter1={{ (index .lastResponse.filters 0)._id }}
    When I save response filter2={{ (index .lastResponse.filters 1)._id }}
    When I save response filter3={{ (index .lastResponse.filters 2)._id }}
    # Test created positions
    When I do GET /api/v4/widgets/{{.widgetId}}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "filters": [
        { "_id": "{{ .filter1 }}" },
        { "_id": "{{ .filter2 }}" },
        { "_id": "{{ .filter3 }}" }
      ]
    }
    """
    # Test updated positions
    When I do PUT /api/v4/widget-filter-positions:
    """json
    [
      "{{ .filter3 }}",
      "{{ .filter1 }}",
      "{{ .filter2 }}"
    ]
    """
    Then the response code should be 204
    When I do GET /api/v4/widgets/{{.widgetId}}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "filters": [
        { "_id": "{{ .filter3 }}" },
        { "_id": "{{ .filter1 }}" },
        { "_id": "{{ .filter2 }}" }
      ]
    }
    """
    When I do PUT /api/v4/widgets/{{.widgetId}}:
    """json
    {
      "title": "test-widget-to-filter-update-position-2-title",
      "tab": "test-tab-to-filter-update-position",
      "type": "test-widget-to-filter-update-position-2-type",
      "filters": [
        {
          "_id": "{{ .filter2 }}",
          "title": "test-widgetfilter-to-update-position-2-2-title",
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-widgetfilter-to-update-position-2-2-pattern"
                }
              }
            ]
          ]
        },
        {
          "_id": "{{ .filter3 }}",
          "title": "test-widgetfilter-to-update-position-2-3-title",
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-widgetfilter-to-update-position-2-3-pattern"
                }
              }
            ]
          ]
        },
        {
          "_id": "{{ .filter1 }}",
          "title": "test-widgetfilter-to-update-position-2-1-title",
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-widgetfilter-to-update-position-2-1-pattern"
                }
              }
            ]
          ]
        }
      ]
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/widgets/{{.widgetId}}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "filters": [
        { "_id": "{{ .filter2 }}" },
        { "_id": "{{ .filter3 }}" },
        { "_id": "{{ .filter1 }}" }
      ]
    }
    """

  Scenario: given update request and no auth user should not allow access
    When I do PUT /api/v4/widget-filter-positions
    Then the response code should be 401

  Scenario: given update request and auth user without view permission should not allow access
    When I am noperms
    When I do PUT /api/v4/widget-filter-positions
    Then the response code should be 403

  Scenario: given invalid request should return error
    When I am test-role-to-filter-update-position-2
    When I do PUT /api/v4/widget-filter-positions:
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
    When I do PUT /api/v4/widget-filter-positions:
    """json
    [
      "test-widgetfilter-not-exist"
    ]
    """
    Then the response code should be 404
    When I do PUT /api/v4/widget-filter-positions:
    """json
    [
      "test-widgetfilter-to-update-position-3-1",
      "test-widgetfilter-not-exist"
    ]
    """
    Then the response code should be 400
    When I do PUT /api/v4/widget-filter-positions:
    """json
    [
      "test-widgetfilter-to-update-position-3-1",
      "test-widgetfilter-to-update-position-3-1"
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
    When I do PUT /api/v4/widget-filter-positions:
    """json
    [
      "test-widgetfilter-to-update-position-3-1",
      "test-widgetfilter-to-update-position-4-1"
    ]
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "error": "filters are related to different widgets or users"
    }
    """
    When I do PUT /api/v4/widget-filter-positions:
    """json
    [
      "test-widgetfilter-to-update-position-3-1",
      "test-widgetfilter-to-update-position-3-2"
    ]
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "error": "filters are related to different widgets or users"
    }
    """
    When I do PUT /api/v4/widget-filter-positions:
    """json
    [
      "test-widgetfilter-to-update-position-3-2",
      "test-widgetfilter-to-update-position-3-3"
    ]
    """
    Then the response code should be 404
    When I do PUT /api/v4/widget-filter-positions:
    """json
    [
      "test-widgetfilter-to-update-position-3-1"
    ]
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "error": "filters are related to different widgets or users"
    }
    """
    When I do PUT /api/v4/widget-filter-positions:
    """json
    [
      "test-widgetfilter-to-update-position-3-3"
    ]
    """
    Then the response code should be 404
    When I do PUT /api/v4/widget-filter-positions:
    """json
    [
      "test-widgetfilter-to-update-position-3-2"
    ]
    """
    Then the response code should be 404
    When I do PUT /api/v4/widget-filter-positions:
    """json
    [
      "test-widgetfilter-to-update-position-5"
    ]
    """
    Then the response code should be 403

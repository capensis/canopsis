Feature: Update view positions
  I need to be able to view positions
  Only admin should be able to view positions

  Scenario: given update request should return ok
    When I am test-role-to-update-view-position
    # Create views
    When I do POST /api/v4/view-groups:
    """json
    {
      "title": "test-viewgroup-to-update-position-1-title"
    }
    """
    Then the response code should be 201
    When I save response group1={{ .lastResponse._id }}
    When I do POST /api/v4/view-groups:
    """json
    {
      "title": "test-viewgroup-to-update-position-2-title"
    }
    """
    Then the response code should be 201
    When I save response group2={{ .lastResponse._id }}
    When I do POST /api/v4/views:
    """json
    {
      "enabled": true,
      "title": "test-view-to-update-position-1-title",
      "group": "{{ .group1 }}"
    }
    """
    Then the response code should be 201
    When I save response view1={{ .lastResponse._id }}
    When I do POST /api/v4/views:
    """json
    {
      "enabled": true,
      "title": "test-view-to-update-position-2-title",
      "group": "{{ .group1 }}"
    }
    """
    Then the response code should be 201
    When I save response view2={{ .lastResponse._id }}
    When I do POST /api/v4/views:
    """json
    {
      "enabled": true,
      "title": "test-view-to-update-position-3-title",
      "group": "{{ .group1 }}"
    }
    """
    Then the response code should be 201
    When I save response view3={{ .lastResponse._id }}
    When I do POST /api/v4/views:
    """json
    {
      "enabled": true,
      "title": "test-view-to-update-position-4-title",
      "group": "{{ .group2 }}"
    }
    """
    Then the response code should be 201
    When I save response view4={{ .lastResponse._id }}
    When I do POST /api/v4/views:
    """json
    {
      "enabled": true,
      "title": "test-view-to-update-position-5-title",
      "group": "{{ .group2 }}"
    }
    """
    Then the response code should be 201
    When I save response view5={{ .lastResponse._id }}
    When I do POST /api/v4/views:
    """json
    {
      "enabled": true,
      "title": "test-view-to-update-position-6-title",
      "group": "{{ .group2 }}"
    }
    """
    Then the response code should be 201
    When I save response view6={{ .lastResponse._id }}
    # Test created positions
    When I do GET /api/v4/view-groups?search=test-viewgroup-to-update-position&with_views=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "{{ .group1 }}",
          "views": [
            { "_id": "{{ .view1 }}" },
            { "_id": "{{ .view2 }}" },
            { "_id": "{{ .view3 }}" }
          ]
        },
        {
          "_id": "{{ .group2 }}",
          "views": [
            { "_id": "{{ .view4 }}" },
            { "_id": "{{ .view5 }}" },
            { "_id": "{{ .view6 }}" }
          ]
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 2
      }
    }
    """
    # Test updated positions
    When I do PUT /api/v4/view-positions:
    """json
    [
      {
        "_id": "{{ .group2 }}",
        "views": [
          "{{ .view5 }}",
          "{{ .view6 }}",
          "{{ .view4 }}"
        ]
      },
      {
        "_id": "{{ .group1 }}",
        "views": [
          "{{ .view3 }}",
          "{{ .view1 }}"
        ]
      }
    ]
    """
    Then the response code should be 204
    When I do GET /api/v4/view-groups?search=test-viewgroup-to-update-position&with_views=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "{{ .group2 }}",
          "views": [
            { "_id": "{{ .view5 }}" },
            { "_id": "{{ .view6 }}" },
            { "_id": "{{ .view4 }}" }
          ]
        },
        {
          "_id": "{{ .group1 }}",
          "views": [
            { "_id": "{{ .view3 }}" },
            { "_id": "{{ .view1 }}" },
            { "_id": "{{ .view2 }}" }
          ]
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 2
      }
    }
    """
    # Test moved view
    When I do PUT /api/v4/views/{{ .view4 }}:
    """json
    {
      "enabled": true,
      "title": "test-view-to-update-position-4-title",
      "group": "{{ .group1 }}"
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/view-groups?search=test-viewgroup-to-update-position&with_views=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "{{ .group2 }}",
          "views": [
            { "_id": "{{ .view5 }}" },
            { "_id": "{{ .view6 }}" }
          ]
        },
        {
          "_id": "{{ .group1 }}",
          "views": [
            { "_id": "{{ .view3 }}" },
            { "_id": "{{ .view1 }}" },
            { "_id": "{{ .view2 }}" },
            { "_id": "{{ .view4 }}" }
          ]
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 2
      }
    }
    """
    # Test partially updated positions
    When I do PUT /api/v4/view-positions:
    """json
    [
      {
        "_id": "{{ .group1 }}",
        "views": [
          "{{ .view4 }}",
          "{{ .view2 }}"
        ]
      }
    ]
    """
    Then the response code should be 204
    When I do GET /api/v4/view-groups?search=test-viewgroup-to-update-position&with_views=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "{{ .group2 }}",
          "views": [
            { "_id": "{{ .view5 }}" },
            { "_id": "{{ .view6 }}" }
          ]
        },
        {
          "_id": "{{ .group1 }}",
          "views": [
            { "_id": "{{ .view3 }}" },
            { "_id": "{{ .view1 }}" },
            { "_id": "{{ .view4 }}" },
            { "_id": "{{ .view2 }}" }
          ]
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 2
      }
    }
    """
    # Test moved view by updated positions
    When I do PUT /api/v4/view-positions:
    """json
    [
      {
        "_id": "{{ .group2 }}",
        "views": [
          "{{ .view5 }}",
          "{{ .view4 }}",
          "{{ .view6 }}"
        ]
      },
      {
        "_id": "{{ .group1 }}",
        "views": [
          "{{ .view3 }}",
          "{{ .view1 }}",
          "{{ .view2 }}"
        ]
      }
    ]
    """
    Then the response code should be 204
    When I do GET /api/v4/view-groups?search=test-viewgroup-to-update-position&with_views=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "{{ .group2 }}",
          "views": [
            { "_id": "{{ .view5 }}" },
            { "_id": "{{ .view4 }}" },
            { "_id": "{{ .view6 }}" }
          ]
        },
        {
          "_id": "{{ .group1 }}",
          "views": [
            { "_id": "{{ .view3 }}" },
            { "_id": "{{ .view1 }}" },
            { "_id": "{{ .view2 }}" }
          ]
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 2
      }
    }
    """
    # Test moved view partially updated positions
    When I do PUT /api/v4/view-positions:
    """json
    [
      {
        "_id": "{{ .group2 }}",
        "views": [
          "{{ .view1 }}"
        ]
      }
    ]
    """
    Then the response code should be 204
    When I do GET /api/v4/view-groups?search=test-viewgroup-to-update-position&with_views=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "{{ .group2 }}",
          "views": [
            { "_id": "{{ .view1 }}" },
            { "_id": "{{ .view5 }}" },
            { "_id": "{{ .view4 }}" },
            { "_id": "{{ .view6 }}" }
          ]
        },
        {
          "_id": "{{ .group1 }}",
          "views": [
            { "_id": "{{ .view3 }}" },
            { "_id": "{{ .view2 }}" }
          ]
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 2
      }
    }
    """

  Scenario: given update request and no auth user should not allow access
    When I do PUT /api/v4/view-positions
    Then the response code should be 401

  Scenario: given update request and auth user without view permission should not allow access
    When I am noperms
    When I do PUT /api/v4/view-positions
    Then the response code should be 403

  Scenario: given invalid request should return error
    When I am test-role-to-update-view-position
    When I do PUT /api/v4/view-positions:
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
    When I do PUT /api/v4/view-positions:
    """json
    [
      {
        "_id": "notexist"
      }
    ]
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "items.0.views": "Views is missing."
      }
    }
    """
    When I do PUT /api/v4/view-positions:
    """json
    [
      {
        "_id": "notexist",
        "views": []
      },
      {
        "_id": "notexist",
        "views": []
      }
    ]
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "items": "Item contains duplicate values."
      }
    }
    """
    When I do PUT /api/v4/view-positions:
    """json
    [
      {
        "_id": "notexist1",
        "views": ["notexistview"]
      },
      {
        "_id": "notexist2",
        "views": ["notexistview"]
      }
    ]
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "items": "Item contains duplicate values."
      }
    }
    """

  Scenario: given request with not exist group should return not found error
    When I am test-role-to-update-view-position
    When I do PUT /api/v4/view-positions:
    """json
    [
      {
        "_id": "notexist",
        "views": []
      }
    ]
    """
    Then the response code should be 404

  Scenario: given request with not exist view should return forbidden error
    When I am test-role-to-update-view-position
    When I do POST /api/v4/view-groups:
    """json
    {
      "title": "test-viewgroup-to-check-not-found-err"
    }
    """
    Then the response code should be 201
    When I do PUT /api/v4/view-positions:
    """json
    [
      {
        "_id": "{{ .lastResponse._id }}",
        "views": ["notexist"]
      }
    ]
    """
    Then the response code should be 403

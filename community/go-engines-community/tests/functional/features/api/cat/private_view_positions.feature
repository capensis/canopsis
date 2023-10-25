Feature: Update view positions for private viewgroups and views

  @concurrent
  Scenario: given update positions request with private groups and views shouldn't be possible
    When I am admin
    When I do POST /api/v4/cat/private-view-groups:
    """json
    {
      "title": "test-private-viewgroup-to-update-position-1-title"
    }
    """
    Then the response code should be 201
    When I save response group1={{ .lastResponse._id }}
    When I do POST /api/v4/cat/private-view-groups:
    """json
    {
      "title": "test-private-viewgroup-to-update-position-2-title"
    }
    """
    Then the response code should be 201
    When I save response group2={{ .lastResponse._id }}
    When I do POST /api/v4/views:
    """json
    {
      "enabled": true,
      "title": "test-private-view-to-update-position-1-title",
      "group": "{{ .group1 }}"
    }
    """
    Then the response code should be 201
    When I save response view1={{ .lastResponse._id }}
    When I do POST /api/v4/views:
    """json
    {
      "enabled": true,
      "title": "test-private-view-to-update-position-2-title",
      "group": "{{ .group1 }}"
    }
    """
    Then the response code should be 201
    When I save response view2={{ .lastResponse._id }}
    When I do POST /api/v4/views:
    """json
    {
      "enabled": true,
      "title": "test-private-view-to-update-position-3-title",
      "group": "{{ .group2 }}"
    }
    """
    Then the response code should be 201
    When I save response view3={{ .lastResponse._id }}
    When I do GET /api/v4/view-groups?search=test-private-viewgroup-to-update-position&with_views=true&with_private=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "{{ .group1 }}",
          "views": [
            { "_id": "{{ .view1 }}" },
            { "_id": "{{ .view2 }}" }
          ]
        },
        {
          "_id": "{{ .group2 }}",
          "views": [
            { "_id": "{{ .view3 }}" }
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
    When I do PUT /api/v4/view-positions:
    """json
    [
      {
        "_id": "{{ .group2 }}",
        "views": [
          "{{ .view3 }}"
        ]
      },
      {
        "_id": "{{ .group1 }}",
        "views": [
          "{{ .view1 }}",
          "{{ .view2 }}"
        ]
      }
    ]
    """
    Then the response code should be 403

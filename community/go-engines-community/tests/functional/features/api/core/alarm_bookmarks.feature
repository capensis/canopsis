Feature: Alarm bookmarks
  I need to be able to add and remove bookmarks to the alarm

  @concurrent
  Scenario: given add bookmarks request and no auth user should not allow access
    When I do PUT /api/v4/alarms/test-alarm-bookmark-1-1/bookmark
    Then the response code should be 401

  @concurrent
  Scenario: given add bookmarks request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/alarms/test-alarm-bookmark-1-1/bookmark
    Then the response code should be 403

  @concurrent
  Scenario: given remove bookmarks request and no auth user should not allow access
    When I do DELETE /api/v4/alarms/test-alarm-bookmark-1-1/bookmark
    Then the response code should be 401

  @concurrent
  Scenario: given remove bookmarks request and auth user without permissions should not allow access
    When I am noperms
    When I do DELETE /api/v4/alarms/test-alarm-bookmark-1-1/bookmark
    Then the response code should be 403

  @concurrent
  Scenario: given add bookmarks requests should add bookmarks, bookmarks should be different for different users
    When I am admin
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-bookmark-1-1",
          "bookmark": false
        },
        {
          "_id": "test-alarm-bookmark-1-2",
          "bookmark": false
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
    When I do PUT /api/v4/alarms/test-alarm-bookmark-1-1/bookmark
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-bookmark-1-1",
          "bookmark": true
        },
        {
          "_id": "test-alarm-bookmark-1-2",
          "bookmark": false
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
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-1&only_bookmarks=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-bookmark-1-1",
          "bookmark": true
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
    When I am manager
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-bookmark-1-1",
          "bookmark": false
        },
        {
          "_id": "test-alarm-bookmark-1-2",
          "bookmark": false
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
    When I do PUT /api/v4/alarms/test-alarm-bookmark-1-2/bookmark
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-bookmark-1-1",
          "bookmark": false
        },
        {
          "_id": "test-alarm-bookmark-1-2",
          "bookmark": true
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
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-1&only_bookmarks=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-bookmark-1-2",
          "bookmark": true
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """

  @concurrent
  Scenario: given remove bookmarks requests should remove bookmarks, bookmarks should be removed differently for different users
    When I am admin
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-bookmark-2-1",
          "bookmark": true
        },
        {
          "_id": "test-alarm-bookmark-2-2",
          "bookmark": false
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
    When I do DELETE /api/v4/alarms/test-alarm-bookmark-2-1/bookmark
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-bookmark-2-1",
          "bookmark": false
        },
        {
          "_id": "test-alarm-bookmark-2-2",
          "bookmark": false
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
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-2&only_bookmarks=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """
    When I am manager
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-bookmark-2-1",
          "bookmark": false
        },
        {
          "_id": "test-alarm-bookmark-2-2",
          "bookmark": true
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
    When I do DELETE /api/v4/alarms/test-alarm-bookmark-2-2/bookmark
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-bookmark-2-1",
          "bookmark": false
        },
        {
          "_id": "test-alarm-bookmark-2-2",
          "bookmark": false
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
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-2&only_bookmarks=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """

  @concurrent
  Scenario: given add bookmarks requests for resolved alarms should add bookmarks, bookmarks should be different for different users
    When I am admin
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-3&opened=false
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-bookmark-3-1",
          "bookmark": false
        },
        {
          "_id": "test-alarm-bookmark-3-2",
          "bookmark": false
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
    When I do PUT /api/v4/alarms/test-alarm-bookmark-3-1/bookmark
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-3&opened=false
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-bookmark-3-1",
          "bookmark": true
        },
        {
          "_id": "test-alarm-bookmark-3-2",
          "bookmark": false
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
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-3&only_bookmarks=true&opened=false
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-bookmark-3-1",
          "bookmark": true
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
    When I am manager
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-3&opened=false
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-bookmark-3-1",
          "bookmark": false
        },
        {
          "_id": "test-alarm-bookmark-3-2",
          "bookmark": false
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
    When I do PUT /api/v4/alarms/test-alarm-bookmark-3-2/bookmark
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-3&opened=false
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-bookmark-3-1",
          "bookmark": false
        },
        {
          "_id": "test-alarm-bookmark-3-2",
          "bookmark": true
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
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-3&only_bookmarks=true&opened=false
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-bookmark-3-2",
          "bookmark": true
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """

  @concurrent
  Scenario: given remove bookmarks requests for resolved alarms should remove bookmarks, bookmarks should be removed differently for different users
    When I am admin
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-4&opened=false
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-bookmark-4-1",
          "bookmark": true
        },
        {
          "_id": "test-alarm-bookmark-4-2",
          "bookmark": false
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
    When I do DELETE /api/v4/alarms/test-alarm-bookmark-4-1/bookmark
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-4&opened=false
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-bookmark-4-1",
          "bookmark": false
        },
        {
          "_id": "test-alarm-bookmark-4-2",
          "bookmark": false
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
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-4&only_bookmarks=true&opened=false
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """
    When I am manager
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-4&opened=false
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-bookmark-4-1",
          "bookmark": false
        },
        {
          "_id": "test-alarm-bookmark-4-2",
          "bookmark": true
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
    When I do DELETE /api/v4/alarms/test-alarm-bookmark-4-2/bookmark
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-4&opened=false
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-bookmark-4-1",
          "bookmark": false
        },
        {
          "_id": "test-alarm-bookmark-4-2",
          "bookmark": false
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
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-4&only_bookmarks=true&opened=false
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """

  @concurrent
  Scenario: given alarm export request with bookmarks should export alarm only with user's bookmarks
    When I am admin
    When I do POST /api/v4/alarm-export:
    """json
    {
      "search": "test-resource-alarm-bookmark-5",
      "only_bookmarks": true,
      "fields": [
         {"name": "_id", "label": "ID"}
      ]
    }
    """
    Then the response code should be 200
    When I save response exportID={{ .lastResponse._id }}
    When I do GET /api/v4/alarm-export/{{ .exportID }} until response code is 200 and body contains:
    """json
    {
       "status": 1
    }
    """
    When I do GET /api/v4/alarm-export/{{ .exportID }}/download
    Then the response code should be 200
    Then the response raw body should be:
    """csv
    ID
    test-alarm-bookmark-5-1

    """
    When I am manager
    When I do POST /api/v4/alarm-export:
    """json
    {
      "search": "test-resource-alarm-bookmark-5",
      "only_bookmarks": true,
      "fields": [
         {"name": "_id", "label": "ID"}
      ]
    }
    """
    Then the response code should be 200
    When I save response exportID={{ .lastResponse._id }}
    When I do GET /api/v4/alarm-export/{{ .exportID }} until response code is 200 and body contains:
    """json
    {
       "status": 1
    }
    """
    When I do GET /api/v4/alarm-export/{{ .exportID }}/download
    Then the response code should be 200
    Then the response raw body should be:
    """csv
    ID
    test-alarm-bookmark-5-2

    """

  @concurrent
  Scenario: given alarm bookmark requests for not found alarm should return not found error
    When I am admin
    When I do PUT /api/v4/alarms/test-alarm-bookmark-not-found/bookmark
    Then the response code should be 404
    Then the response body should be:
    """
    {
      "error": "Not found"
    }
    """
    When I do DELETE /api/v4/alarms/test-alarm-bookmark-not-found/bookmark
    Then the response code should be 404
    Then the response body should be:
    """
    {
      "error": "Not found"
    }
    """

  @concurrent
  Scenario: given add bookmarks requests for alarms represented in both main and resolved collections
    should add bookmarks for both documents, bookmarks should be different for different users
    When I am admin
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-6
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-bookmark-6-1",
          "bookmark": false
        },
        {
          "_id": "test-alarm-bookmark-6-2",
          "bookmark": false
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
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-6&opened=false
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-bookmark-6-1",
          "bookmark": false
        },
        {
          "_id": "test-alarm-bookmark-6-2",
          "bookmark": false
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
    When I do PUT /api/v4/alarms/test-alarm-bookmark-6-1/bookmark
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-6
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-bookmark-6-1",
          "bookmark": true
        },
        {
          "_id": "test-alarm-bookmark-6-2",
          "bookmark": false
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
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-6&opened=false
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-bookmark-6-1",
          "bookmark": true
        },
        {
          "_id": "test-alarm-bookmark-6-2",
          "bookmark": false
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
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-6&only_bookmarks=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-bookmark-6-1",
          "bookmark": true
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-6&only_bookmarks=true&opened=false
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-bookmark-6-1",
          "bookmark": true
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
    When I am manager
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-6
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-bookmark-6-1",
          "bookmark": false
        },
        {
          "_id": "test-alarm-bookmark-6-2",
          "bookmark": false
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
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-6&opened=false
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-bookmark-6-1",
          "bookmark": false
        },
        {
          "_id": "test-alarm-bookmark-6-2",
          "bookmark": false
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
    When I do PUT /api/v4/alarms/test-alarm-bookmark-6-2/bookmark
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-6
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-bookmark-6-1",
          "bookmark": false
        },
        {
          "_id": "test-alarm-bookmark-6-2",
          "bookmark": true
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
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-6&opened=false
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-bookmark-6-1",
          "bookmark": false
        },
        {
          "_id": "test-alarm-bookmark-6-2",
          "bookmark": true
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
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-6&only_bookmarks=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-bookmark-6-2",
          "bookmark": true
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-6&only_bookmarks=true&opened=false
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-bookmark-6-2",
          "bookmark": true
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """

  @concurrent
  Scenario: given remove bookmarks requests for alarms represented in both main and resolved collections
    should add bookmarks for both documents, bookmarks should be different for different users
    When I am admin
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-7
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-bookmark-7-1",
          "bookmark": true
        },
        {
          "_id": "test-alarm-bookmark-7-2",
          "bookmark": false
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
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-7&opened=false
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-bookmark-7-1",
          "bookmark": true
        },
        {
          "_id": "test-alarm-bookmark-7-2",
          "bookmark": false
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
    When I do DELETE /api/v4/alarms/test-alarm-bookmark-7-1/bookmark
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-7
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-bookmark-7-1",
          "bookmark": false
        },
        {
          "_id": "test-alarm-bookmark-7-2",
          "bookmark": false
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
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-7&opened=false
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-bookmark-7-1",
          "bookmark": false
        },
        {
          "_id": "test-alarm-bookmark-7-2",
          "bookmark": false
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
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-7&only_bookmarks=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-7&only_bookmarks=true&opened=false
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """
    When I am manager
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-7
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-bookmark-7-1",
          "bookmark": false
        },
        {
          "_id": "test-alarm-bookmark-7-2",
          "bookmark": true
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
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-7&opened=false
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-bookmark-7-1",
          "bookmark": false
        },
        {
          "_id": "test-alarm-bookmark-7-2",
          "bookmark": true
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
    When I do DELETE /api/v4/alarms/test-alarm-bookmark-7-2/bookmark
    Then the response code should be 204
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-7
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-bookmark-7-1",
          "bookmark": false
        },
        {
          "_id": "test-alarm-bookmark-7-2",
          "bookmark": false
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
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-7&opened=false
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-alarm-bookmark-7-1",
          "bookmark": false
        },
        {
          "_id": "test-alarm-bookmark-7-2",
          "bookmark": false
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
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-7&only_bookmarks=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-7&only_bookmarks=true&opened=false
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 0
      }
    }
    """

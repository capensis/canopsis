Feature: Alarm bookmarks
  I need to be able to add and remove bookmarks to the alarm

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
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-1&with_bookmarks=true
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
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-1&with_bookmarks=true
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
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-2&with_bookmarks=true
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
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-2&with_bookmarks=true
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
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-3&with_bookmarks=true&opened=false
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
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-3&with_bookmarks=true&opened=false
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
  Scenario: given remove bookmarks requests should remove bookmarks, bookmarks should be removed differently for different users
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
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-4&with_bookmarks=true&opened=false
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
    When I do GET /api/v4/alarms?search=test-resource-alarm-bookmark-4&with_bookmarks=true&opened=false
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

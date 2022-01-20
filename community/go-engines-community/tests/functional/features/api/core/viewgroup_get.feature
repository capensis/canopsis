Feature: Get a view group
  I need to be able to get a view group
  Only admin should be able to get a view group

  Scenario: given search request should return view groups
    When I am admin
    When I do GET /api/v4/view-groups?search=test-viewgroup-to-get
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "test-viewgroup-to-get-1",
          "title": "test-viewgroup-to-get-1-title",
          "author": "test-viewgroup-to-get-1-author",
          "created": 1611229670,
          "updated": 1611229670
        },
        {
          "_id": "test-viewgroup-to-get-2",
          "title": "test-viewgroup-to-get-2-title",
          "author": "test-viewgroup-to-get-2-author",
          "created": 1611229670,
          "updated": 1611229670
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

  Scenario: given get request should return viewgroup
    When I am admin
    When I do GET /api/v4/view-groups/test-viewgroup-to-get-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-viewgroup-to-get-1",
      "title": "test-viewgroup-to-get-1-title",
      "author": "test-viewgroup-to-get-1-author",
      "created": 1611229670,
      "updated": 1611229670
    }
    """

  Scenario: given with views request should return view groups
    When I am admin
    When I do GET /api/v4/view-groups?search=test-viewgroup-to-get&with_views=true
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "test-viewgroup-to-get-1",
          "title": "test-viewgroup-to-get-1-title",
          "author": "test-viewgroup-to-get-1-author",
          "views": [
            {
              "_id": "test-view-to-viewgroup-get-2",
              "author": "test-view-to-viewgroup-get-2-author",
              "created": 1611229670,
              "description": "test-view-to-viewgroup-get-2-description",
              "enabled": true,
              "group": {
                "_id": "test-viewgroup-to-get-1",
                "author": "test-viewgroup-to-get-1-author",
                "created": 1611229670,
                "title": "test-viewgroup-to-get-1-title",
                "updated": 1611229670
              },
              "title": "test-view-to-viewgroup-get-2-title",
              "periodic_refresh": {
                "enabled": true,
                "value": 1,
                "unit": "s"
              },
              "tags": [
                "test-view-to-viewgroup-get-2-tag"
              ],
              "updated": 1611229670
            },
            {
              "_id": "test-view-to-viewgroup-get-1",
              "author": "test-view-to-viewgroup-get-1-author",
              "created": 1611229670,
              "description": "test-view-to-viewgroup-get-1-description",
              "enabled": true,
              "group": {
                "_id": "test-viewgroup-to-get-1",
                "author": "test-viewgroup-to-get-1-author",
                "created": 1611229670,
                "title": "test-viewgroup-to-get-1-title",
                "updated": 1611229670
              },
              "title": "test-view-to-viewgroup-get-1-title",
              "periodic_refresh": {
                "enabled": true,
                "value": 1,
                "unit": "s"
              },
              "tags": [
                "test-view-to-viewgroup-get-1-tag"
              ],
              "updated": 1611229670
            }
          ],
          "created": 1611229670,
          "updated": 1611229670
        },
        {
          "_id": "test-viewgroup-to-get-2",
          "title": "test-viewgroup-to-get-2-title",
          "author": "test-viewgroup-to-get-2-author",
          "views": [],
          "created": 1611229670,
          "updated": 1611229670
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
    When I do GET /api/v4/view-groups?search=test-viewgroup-to-get&with_views=true&with_tabs=true
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "test-viewgroup-to-get-1",
          "title": "test-viewgroup-to-get-1-title",
          "author": "test-viewgroup-to-get-1-author",
          "views": [
            {
              "_id": "test-view-to-viewgroup-get-2",
              "author": "test-view-to-viewgroup-get-2-author",
              "created": 1611229670,
              "description": "test-view-to-viewgroup-get-2-description",
              "enabled": true,
              "group": {
                "_id": "test-viewgroup-to-get-1",
                "author": "test-viewgroup-to-get-1-author",
                "created": 1611229670,
                "title": "test-viewgroup-to-get-1-title",
                "updated": 1611229670
              },
              "title": "test-view-to-viewgroup-get-2-title",
              "periodic_refresh": {
                "enabled": true,
                "value": 1,
                "unit": "s"
              },
              "tags": [
                "test-view-to-viewgroup-get-2-tag"
              ],
              "tabs": [],
              "updated": 1611229670
            },
            {
              "_id": "test-view-to-viewgroup-get-1",
              "author": "test-view-to-viewgroup-get-1-author",
              "created": 1611229670,
              "description": "test-view-to-viewgroup-get-1-description",
              "enabled": true,
              "group": {
                "_id": "test-viewgroup-to-get-1",
                "author": "test-viewgroup-to-get-1-author",
                "created": 1611229670,
                "title": "test-viewgroup-to-get-1-title",
                "updated": 1611229670
              },
              "title": "test-view-to-viewgroup-get-1-title",
              "periodic_refresh": {
                "enabled": true,
                "value": 1,
                "unit": "s"
              },
              "tags": [
                "test-view-to-viewgroup-get-1-tag"
              ],
              "tabs": [
                {
                  "_id": "test-tab-to-viewgroup-get-1",
                  "title": "test-tab-to-viewgroup-get-1-title",
                  "author": "test-user-to-viewgroup-get",
                  "created": 1611229670,
                  "updated": 1611229670
                },
                {
                  "_id": "test-tab-to-viewgroup-get-2",
                  "title": "test-tab-to-viewgroup-get-2-title",
                  "author": "test-user-to-viewgroup-get",
                  "created": 1611229670,
                  "updated": 1611229670
                }
              ],
              "updated": 1611229670
            }
          ],
          "created": 1611229670,
          "updated": 1611229670
        },
        {
          "_id": "test-viewgroup-to-get-2",
          "title": "test-viewgroup-to-get-2-title",
          "author": "test-viewgroup-to-get-2-author",
          "views": [],
          "created": 1611229670,
          "updated": 1611229670
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

  Scenario: given with flags request should return view groups
    When I am admin
    When I do GET /api/v4/view-groups?search=test-viewgroup-to-get&with_flags=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-viewgroup-to-get-1",
          "deletable": false
        },
        {
          "_id": "test-viewgroup-to-get-2",
          "deletable": true
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
    When I do GET /api/v4/view-groups?search=test-viewgroup-to-get&with_views=true&with_flags=true
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "test-viewgroup-to-get-1",
          "title": "test-viewgroup-to-get-1-title",
          "author": "test-viewgroup-to-get-1-author",
          "deletable": false,
          "views": [
            {
              "_id": "test-view-to-viewgroup-get-2",
              "author": "test-view-to-viewgroup-get-2-author",
              "created": 1611229670,
              "description": "test-view-to-viewgroup-get-2-description",
              "enabled": true,
              "group": {
                "_id": "test-viewgroup-to-get-1",
                "author": "test-viewgroup-to-get-1-author",
                "created": 1611229670,
                "title": "test-viewgroup-to-get-1-title",
                "updated": 1611229670
              },
              "title": "test-view-to-viewgroup-get-2-title",
              "periodic_refresh": {
                "enabled": true,
                "value": 1,
                "unit": "s"
              },
              "tags": [
                "test-view-to-viewgroup-get-2-tag"
              ],
              "updated": 1611229670
            },
            {
              "_id": "test-view-to-viewgroup-get-1",
              "author": "test-view-to-viewgroup-get-1-author",
              "created": 1611229670,
              "description": "test-view-to-viewgroup-get-1-description",
              "enabled": true,
              "group": {
                "_id": "test-viewgroup-to-get-1",
                "author": "test-viewgroup-to-get-1-author",
                "created": 1611229670,
                "title": "test-viewgroup-to-get-1-title",
                "updated": 1611229670
              },
              "title": "test-view-to-viewgroup-get-1-title",
              "periodic_refresh": {
                "enabled": true,
                "value": 1,
                "unit": "s"
              },
              "tags": [
                "test-view-to-viewgroup-get-1-tag"
              ],
              "updated": 1611229670
            }
          ],
          "created": 1611229670,
          "updated": 1611229670
        },
        {
          "_id": "test-viewgroup-to-get-2",
          "title": "test-viewgroup-to-get-2-title",
          "author": "test-viewgroup-to-get-2-author",
          "deletable": true,
          "views": [],
          "created": 1611229670,
          "updated": 1611229670
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

  Scenario: given get all request and no auth user should not allow access
    When I do GET /api/v4/view-groups
    Then the response code should be 401

  Scenario: given get all request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/view-groups
    Then the response code should be 403

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/view-groups/test-viewgroup-to-get-1
    Then the response code should be 401

  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/view-groups/test-viewgroup-to-get-1
    Then the response code should be 403

  Scenario: given get request with not exist id should return not found error
    When I am admin
    When I do GET /api/v4/view-groups/test-viewgroup-not-found
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """

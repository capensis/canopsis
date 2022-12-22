Feature: Get a playlist
  I need to be able to get a playlist
  Only admin should be able to get a playlist

  Scenario: given search request should return playlists
    When I am admin
    When I do GET /api/v4/playlists?search=test-playlist-to-get
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "test-playlist-to-get-1",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "created": 1608284568,
          "enabled": true,
          "fullscreen": true,
          "interval": {
            "value": 10,
            "unit": "s"
          },
          "name": "test-playlist-to-get-1-name",
          "tabs_list": [
            "test-tab-to-playlist-edit-1"
          ],
          "updated": 1608285370
        },
        {
          "_id": "test-playlist-to-get-2",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "created": 1608284568,
          "enabled": true,
          "fullscreen": true,
          "interval": {
            "value": 20,
            "unit": "s"
          },
          "name": "test-playlist-to-get-2-name",
          "tabs_list": [
            "test-tab-to-playlist-edit-2"
          ],
          "updated": 1608285370
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

  Scenario: given get request should return playlist
    When I am admin
    When I do GET /api/v4/playlists/test-playlist-to-get-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-playlist-to-get-1",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "created": 1608284568,
      "enabled": true,
      "fullscreen": true,
      "interval": {
        "value": 10,
        "unit": "s"
      },
      "name": "test-playlist-to-get-1-name",
      "tabs_list": [
        "test-tab-to-playlist-edit-1"
      ],
      "updated": 1608285370
    }
    """

  Scenario: given sort request should return sorted playlists
    When I am admin
    When I do GET /api/v4/playlists?search=test-playlist-to-get&sort=desc&sort_by=name
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-playlist-to-get-2"
        },
        {
          "_id": "test-playlist-to-get-1"
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

  Scenario: given search request should not return views without access
    When I am admin
    When I do GET /api/v4/playlists?search=test-playlist-to-check-access
    Then the response code should be 200
    Then the response body should be:
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

  Scenario: given get all request and no auth user should not allow access
    When I do GET /api/v4/playlists
    Then the response code should be 401

  Scenario: given get all request and auth user without playlist permission should not allow access
    When I am noperms
    When I do GET /api/v4/playlists
    Then the response code should be 403

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/playlists/notexist
    Then the response code should be 401

  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/playlists/notexist
    Then the response code should be 403

  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am admin
    When I do GET /api/v4/views/test-view-to-check-access
    Then the response code should be 403

  Scenario: given get request with not exist id should return not allow access error
    When I am admin
    When I do GET /api/v4/playlists/notexist
    Then the response code should be 403

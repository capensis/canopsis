Feature: Update an playlist
  I need to be able to update an playlist
  Only admin should be able to update an playlist

  Scenario: given update request should update playlist
    When I am admin
    Then I do PUT /api/v4/playlists/test-playlist-to-update:
    """json
    {
      "fullscreen": false,
      "name": "test-playlist-to-update-name-updated",
      "tabs_list": [
        "test-tab-to-playlist-edit-1"
      ],
      "interval": {
        "value": 120,
        "unit": "m"
      },
      "enabled": false
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-playlist-to-update",
      "author": {
        "_id": "root",
        "name": "root"
      },
      "created": 1608284568,
      "enabled": false,
      "fullscreen": false,
      "interval": {
        "value": 120,
        "unit": "m"
      },
      "name": "test-playlist-to-update-name-updated",
      "tabs_list": [
        "test-tab-to-playlist-edit-1"
      ]
    }
    """
    When I do GET /api/v4/permissions?search=test-playlist-to-update
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "test-playlist-to-update",
          "name": "test-playlist-to-update",
          "description": "Rights on playlist : test-playlist-to-update-name-updated",
          "playlist": {
            "_id": "test-playlist-to-update",
            "name": "test-playlist-to-update-name-updated"
          },
          "type": "RW"
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

  Scenario: given update request with already exists name should return error
    When I am admin
    Then I do PUT /api/v4/playlists/test-playlist-to-update:
    """json
    {
      "name": "test-playlist-to-check-unique-name-name"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "name": "Name already exists."
      }
    }
    """

  Scenario: given get request and no auth user should not allow access
    When I do PUT /api/v4/playlists/notexist
    Then the response code should be 401

  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/playlists/notexist
    Then the response code should be 403

  Scenario: given get request and auth user without view permission should not allow access
    When I am admin
    When I do PUT /api/v4/playlists/test-view-to-check-access
    Then the response code should be 403

  Scenario: given invalid update request should return errors
    When I am admin
    When I do PUT /api/v4/playlists/test-view-to-update:
    """json
    {
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "enabled": "Enabled is missing.",
        "fullscreen": "Fullscreen is missing.",
        "interval.value": "Value is missing.",
        "interval.unit": "Unit is missing.",
        "name": "Name is missing.",
        "tabs_list": "TabsList is missing."
      }
    }
    """

  Scenario: given invalid update request should return errors
    When I am admin
    When I do PUT /api/v4/playlists/test-view-to-update:
    """json
    {
      "fullscreen": true,
      "name": "test-view-to-update-name",
      "tabs_list": ["notexist", "test-tab-to-playlist-edit-1"],
      "interval": {
        "value": 10,
        "unit": "s"
      },
      "enabled": true
    }
    """
    Then the response code should be 403
    When I do PUT /api/v4/playlists/test-view-to-update:
    """json
    {
      "fullscreen": true,
      "name": "test-view-to-update-name",
      "tabs_list": ["test-tab-to-check-access", "test-tab-to-playlist-edit-1"],
      "interval": {
        "value": 10,
        "unit": "s"
      },
      "enabled": true
    }
    """
    Then the response code should be 403

  Scenario: given update request with not exist id should return not allow access error
    When I am admin
    When I do PUT /api/v4/playlists/notexist
    Then the response code should be 403

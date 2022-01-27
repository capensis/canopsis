Feature: Update an playlist
  I need to be able to update an playlist
  Only admin should be able to update an playlist

  Scenario: given update request should update playlist
    When I am admin
    Then I do PUT /api/v4/playlists/test-playlist-to-update:
    """
    {
      "fullscreen": false,
      "name": "test-playlist-to-update",
      "tabs_list": [
        "test-view-to-edit-playlist-tab-1"
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
    """
    {
      "_id": "test-playlist-to-update",
      "author": "root",
      "created": 1608284568,
      "enabled": false,
      "fullscreen": false,
      "interval": {
        "value": 120,
        "unit": "m"
      },
      "name": "test-playlist-to-update",
      "tabs_list": [
        "test-view-to-edit-playlist-tab-1"
      ]
    }
    """

  Scenario: given update request with already exists name should return error
    When I am admin
    Then I do PUT /api/v4/playlists/test-playlist-to-update:
    """
    {
      "name": "test-playlist-to-check-unique-name-name"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
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
    """
    {
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
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

  Scenario: given update request with not exist id should return not allow access error
    When I am admin
    When I do PUT /api/v4/playlists/notexist
    Then the response code should be 403

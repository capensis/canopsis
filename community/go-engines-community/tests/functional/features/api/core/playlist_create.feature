Feature: Create a playlist
  I need to be able to create a playlist
  Only admin should be able to create a playlist

  Scenario: given create request should return ok
    When I am admin
    When I do POST /api/v4/playlists:
    """json
    {
      "fullscreen": true,
      "name": "test-playlist-to-create-1-name",
      "tabs_list": [
        "test-tab-to-playlist-edit-1"
      ],
      "interval": {
        "value": 10,
        "unit": "s"
      },
      "enabled": true
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "fullscreen": true,
      "name": "test-playlist-to-create-1-name",
      "author": "root",
      "tabs_list": [
        "test-tab-to-playlist-edit-1"
      ],
      "interval": {
        "value": 10,
        "unit": "s"
      },
      "enabled": true
    }
    """

  Scenario: given create request should return ok to get request
    When I am admin
    When I do POST /api/v4/playlists:
    """json
    {
      "fullscreen": true,
      "name": "test-playlist-to-create-2-name",
      "tabs_list": [
        "test-tab-to-playlist-edit-1"
      ],
      "interval": {
        "value": 10,
        "unit": "s"
      },
      "enabled": true
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/playlists/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "fullscreen": true,
      "name": "test-playlist-to-create-2-name",
      "author": "root",
      "tabs_list": [
        "test-tab-to-playlist-edit-1"
      ],
      "interval": {
        "value": 10,
        "unit": "s"
      },
      "enabled": true
    }
    """

  Scenario: given create request should create new permission
    When I am admin
    When I do POST /api/v4/playlists:
    """json
    {
      "fullscreen": true,
      "name": "test-playlist-to-create-3-name",
      "tabs_list": [
        "test-tab-to-playlist-edit-1"
      ],
      "interval": {
        "value": 10,
        "unit": "s"
      },
      "enabled": true
    }
    """
    Then the response code should be 201
    Then I save response playlistId={{ .lastResponse._id }}
    When I do GET /api/v4/permissions?search={{ .playlistId }}
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "{{ .playlistId }}",
          "name": "{{ .playlistId }}",
          "description": "Rights on playlist : test-playlist-to-create-3-name",
          "playlist": {
            "_id": "{{ .playlistId }}",
            "name": "test-playlist-to-create-3-name"
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

  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/playlists
    Then the response code should be 401

  Scenario: given create request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/playlists
    Then the response code should be 403

  Scenario: given invalid create request should return errors
    When I am admin
    When I do POST /api/v4/playlists:
    """json
    {
      "tabs_list": []
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
        "tabs_list": "TabsList should not be blank."
      }
    }
    """

  Scenario: given invalid create request should return errors
    When I am admin
    When I do POST /api/v4/playlists:
    """json
    {
      "fullscreen": true,
      "name": "test-playlist-to-create-4-name",
      "tabs_list": ["notexist", "test-tab-to-playlist-edit-1"],
      "interval": {
        "value": 10,
        "unit": "s"
      },
      "enabled": true
    }
    """
    Then the response code should be 403
    When I do POST /api/v4/playlists:
    """json
    {
      "fullscreen": true,
      "name": "test-playlist-to-create-4-name",
      "tabs_list": ["test-tab-to-check-access", "test-tab-to-playlist-edit-1"],
      "interval": {
        "value": 10,
        "unit": "s"
      },
      "enabled": true
    }
    """
    Then the response code should be 403

  Scenario: given create request with already exists name should return error
    When I am admin
    When I do POST /api/v4/playlists:
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

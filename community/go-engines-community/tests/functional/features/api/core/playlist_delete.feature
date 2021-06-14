Feature: Delete a playlist
  I need to be able to delete a playlist
  Only admin should be able to delete a playlist

  Scenario: given delete request should delete playlist
    When I am admin
    When I do DELETE /api/v4/playlists/test-playlist-to-delete
    Then the response code should be 204
    Then I do GET /api/v4/playlists/test-playlist-to-delete
    Then the response code should be 403
    When I do GET /api/v4/permissions?search=test-playlist-to-delete
    Then the response code should be 200
    Then the response body should be:
    """
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

  Scenario: given delete request and no auth user should not allow access
    When I do DELETE /api/v4/playlists/notexist
    Then the response code should be 401

  Scenario: given delete request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do DELETE /api/v4/playlists/notexist
    Then the response code should be 403

  Scenario: given delete request and auth user without playlist permission should not allow access
    When I am admin
    When I do DELETE /api/v4/playlists/test-playlist-to-check-access
    Then the response code should be 403

  Scenario: given delete request with not exist id should return not allow access error
    When I am admin
    When I do DELETE /api/v4/playlists/notexist
    Then the response code should be 403

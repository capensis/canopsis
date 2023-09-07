Feature: Delete a color theme
  I need to be able to delete a color theme
  Only admin should be able to delete a color theme

  @concurrent
  Scenario: given delete request and no auth user should not allow access
    When I do DELETE /api/v4/color-themes/test_theme_to_delete_1
    Then the response code should be 401

  @concurrent
  Scenario: given delete request and auth user without permissions should not allow access
    When I am noperms
    When I do DELETE /api/v4/color-themes/test_theme_to_delete_1
    Then the response code should be 403

  @concurrent
  Scenario: given delete request should return ok
  When I am admin
  When I do DELETE /api/v4/color-themes/test_theme_to_delete_1
  Then the response code should be 204
  When I do GET /api/v4/color-themes/test_theme_to_delete_1
  Then the response code should be 404

  @concurrent
  Scenario: given delete request for not found theme should return error
  When I am admin
  When I do DELETE /api/v4/color-themes/test_theme_not_found
  Then the response code should be 404
  Then the response body should contain:
  """json
  {
    "error": "Not found"
  }
  """

  @concurrent
  Scenario: given delete request for default theme should return error
  When I am admin
  When I do DELETE /api/v4/color-themes/canopsis
  Then the response code should be 400
  Then the response body should contain:
  """json
  {
    "error": "can't modify or delete the default color theme"
  }
  """

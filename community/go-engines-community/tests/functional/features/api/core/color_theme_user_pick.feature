Feature: Pick a color theme
  I need to be able to pick a color theme
  Any user should be able to pick a color theme

  @concurrent
  Scenario: given user update request with default color theme should be ok
  When I am authenticated with username "test-user-to-pick-color-theme-1" and password "test"
  When I do PUT /api/v4/account/me:
  """json
  {
    "ui_language": "fr",
    "ui_groups_navigation_type": "top-bar",
    "defaultview": "test-view-to-edit-user",
    "ui_theme": "canopsis_dark"
  }
  """
  Then the response code should be 200
  Then the response body should contain:
  """json
  {
    "ui_theme": "canopsis_dark"
  }
  """

  @concurrent
  Scenario: given user update request with custom color theme should be ok
  When I am authenticated with username "test-user-to-pick-color-theme-2" and password "test"
  When I do PUT /api/v4/account/me:
  """json
  {
    "ui_language": "fr",
    "ui_groups_navigation_type": "top-bar",
    "defaultview": "test-view-to-edit-user",
    "ui_theme": "test_theme_to_pick_1"
  }
  """
  Then the response code should be 200
  Then the response body should contain:
  """json
  {
    "ui_theme": "test_theme_to_pick_1"
  }
  """

  @concurrent
  Scenario: given user update request with not found color theme should return error
  When I am authenticated with username "test-user-to-pick-color-theme-3" and password "test"
  When I do PUT /api/v4/account/me:
  """json
  {
    "ui_language": "fr",
    "ui_groups_navigation_type": "top-bar",
    "defaultview": "test-view-to-edit-user",
    "ui_theme": "test_theme_to_pick_not_found"
  }
  """
  Then the response code should be 400
  Then the response body should contain:
  """json
  {
    "errors": {
      "ui_theme": "UITheme doesn't exist."
    }
  }
  """

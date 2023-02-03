Feature: Get template env vars
  I need to be able to get template env vars

  @concurrent
  Scenario: given request and no auth should not allow access
    When I do GET /api/v4/template-vars
    Then the response code should be 401

  @concurrent
  Scenario: given request should return env vars
    When I am admin
    When I do GET /api/v4/template-vars:
    """json
    [
      ".Env.Location",
      ".Env.Name"
    ]
    """

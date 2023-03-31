Feature: Get template env vars
  I need to be able to get template env vars

  @concurrent
  Scenario: given request should return env vars
    When I do GET /api/v4/template-vars:
    """json
    {
        "Name": "Test",
        "Location": "FR"
    }
    """

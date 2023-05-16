Feature: Get template env vars
  I need to be able to get template env vars

  @concurrent
  Scenario: given request should return env vars
    When I do GET /api/v4/template-vars
    Then the response body should contain:
    """json
    {
      "Name": "Test",
      "Location": "FR",
      "System": {
        "CPS_MONGO_URL": "{{ .CPS_MONGO_URL }}"
      }
    }
    """

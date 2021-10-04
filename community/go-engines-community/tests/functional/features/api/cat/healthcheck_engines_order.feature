Feature: get healthcheck engines' order
  I need to be able to get healthcheck engines' order
  Only admin should be able to get healthcheck engines' order

  Scenario: given get engines' order request and no auth user should not allow access
    When I do GET /api/v4/cat/healthcheck/engines-order
    Then the response code should be 401

  Scenario: given get engines' order request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/cat/healthcheck/engines-order
    Then the response code should be 403

  Scenario: given get engines' order request should return ok
    When I am admin
    When I do GET /api/v4/cat/healthcheck/engines-order
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "nodes": [
        "engine-fifo",
        "engine-che",
        "engine-pbehavior",
        "engine-axe",
        "engine-axe",
        "engine-correlation",
        "engine-service",
        "engine-dynamic-infos",
        "engine-action"
      ],
      "edges": [
        {
          "from": "engine-fifo",
          "to": "engine-che"
        },
        {
          "from": "engine-che",
          "to": "engine-pbehavior"
        },
        {
          "from": "engine-pbehavior",
          "to": "engine-axe"
        },
        {
          "from": "engine-axe",
          "to": "engine-correlation"
        },
        {
          "from": "engine-axe",
          "to": "engine-remediation"
        },
        {
          "from": "engine-correlation",
          "to": "engine-service"
        },
        {
          "from": "engine-service",
          "to": "engine-dynamic-infos"
        },
        {
          "from": "engine-dynamic-infos",
          "to": "engine-action"
        },
        {
          "from": "engine-action",
          "to": "engine-webhook"
        }
      ]
    }
    """

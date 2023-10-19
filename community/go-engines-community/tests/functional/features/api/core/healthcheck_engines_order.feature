Feature: get healthcheck engines' order
  I need to be able to get healthcheck engines' order
  Only admin should be able to get healthcheck engines' order

  Scenario: given get engines' order request and no auth user should not allow access
    When I do GET /api/v4/healthcheck/engines-order
    Then the response code should be 401

  Scenario: given get engines' order request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/healthcheck/engines-order
    Then the response code should be 403

  Scenario: given get engines' order request should return ok
    When I am admin
    When I do GET /api/v4/healthcheck/engines-order
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "nodes": [
        "engine-fifo",
        "engine-che",
        "engine-axe",
        "engine-correlation",
        "engine-remediation",
        "engine-pbehavior",
        "engine-dynamic-infos",
        "engine-action",
        "engine-webhook"
      ],
      "edges": [
        {
          "from": "engine-fifo",
          "to": "engine-che"
        },
        {
          "from": "engine-che",
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
          "from": "engine-axe",
          "to": "engine-pbehavior"
        },
        {
          "from": "engine-correlation",
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
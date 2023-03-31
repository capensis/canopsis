Feature: to get and update healthcheck parameters
  I need to be able to get and update healthcheck parameters
  Only admin should be able to get and update healthcheck parameters

  Scenario: given get healthcheck parameters request and no auth user should not allow access
    When I do GET /api/v4/cat/healthcheck/parameters
    Then the response code should be 401

  Scenario: given get healthcheck parameters request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/cat/healthcheck/parameters
    Then the response code should be 403

  Scenario: given update healthcheck parameters request and no auth user should not allow access
    When I do PUT /api/v4/cat/healthcheck/parameters
    Then the response code should be 401

  Scenario: given update healthcheck parameters request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/cat/healthcheck/parameters
    Then the response code should be 403

  Scenario: given update and get healthcheck parameters requests should return ok
    When I am admin
    When I do PUT /api/v4/cat/healthcheck/parameters:
    """
    {}
    """
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "queue": {
        "limit": 0,
        "enabled": false
      },
      "messages": {
        "limit": 0,
        "enabled": false
      },
      "engine-fifo": {
        "enabled": false,
        "minimal": 0,
        "optimal": 0
      },
      "engine-pbehavior": {
        "enabled": false,
        "minimal": 0,
        "optimal": 0
      },
      "engine-service": {
        "enabled": false,
        "minimal": 0,
        "optimal": 0
      },
      "engine-action": {
        "enabled": false,
        "minimal": 0,
        "optimal": 0
      },
      "engine-che": {
        "enabled": false,
        "minimal": 0,
        "optimal": 0
      },
      "engine-remediation": {
        "enabled": false,
        "minimal": 0,
        "optimal": 0
      },
      "engine-webhook": {
        "enabled": false,
        "minimal": 0,
        "optimal": 0
      },
      "engine-axe": {
        "enabled": false,
        "minimal": 0,
        "optimal": 0
      },
      "engine-correlation": {
        "enabled": false,
        "minimal": 0,
        "optimal": 0
      },
      "engine-dynamic-infos": {
        "enabled": false,
        "minimal": 0,
        "optimal": 0
      }
    }
    """
    When I do GET /api/v4/cat/healthcheck/parameters
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "queue": {
        "limit": 0,
        "enabled": false
      },
      "messages": {
        "limit": 0,
        "enabled": false
      },
      "engine-fifo": {
        "enabled": false,
        "minimal": 0,
        "optimal": 0
      },
      "engine-pbehavior": {
        "enabled": false,
        "minimal": 0,
        "optimal": 0
      },
      "engine-service": {
        "enabled": false,
        "minimal": 0,
        "optimal": 0
      },
      "engine-action": {
        "enabled": false,
        "minimal": 0,
        "optimal": 0
      },
      "engine-che": {
        "enabled": false,
        "minimal": 0,
        "optimal": 0
      },
      "engine-remediation": {
        "enabled": false,
        "minimal": 0,
        "optimal": 0
      },
      "engine-webhook": {
        "enabled": false,
        "minimal": 0,
        "optimal": 0
      },
      "engine-axe": {
        "enabled": false,
        "minimal": 0,
        "optimal": 0
      },
      "engine-correlation": {
        "enabled": false,
        "minimal": 0,
        "optimal": 0
      },
      "engine-dynamic-infos": {
        "enabled": false,
        "minimal": 0,
        "optimal": 0
      }
    }
    """
    When I do PUT /api/v4/cat/healthcheck/parameters:
    """
    {
      "queue": {
        "limit": 1000,
        "enabled": true
      },
      "messages": {
        "limit": 5000,
        "enabled": true
      },
      "engine-fifo": {
        "enabled": true,
        "minimal": 1,
        "optimal": 2
      },
      "engine-pbehavior": {
        "enabled": true,
        "minimal": 1,
        "optimal": 2
      },
      "engine-service": {
        "enabled": true,
        "minimal": 1,
        "optimal": 3
      },
      "engine-action": {
        "enabled": false
      },
      "engine-che": {
        "enabled": true,
        "minimal": 1,
        "optimal": 1
      },
      "engine-remediation": {
        "enabled": true,
        "minimal": 1,
        "optimal": 1
      },
      "engine-webhook": {
        "enabled": true,
        "minimal": 1,
        "optimal": 4
      },
      "engine-axe": {
        "enabled": true,
        "minimal": 1,
        "optimal": 2
      },
      "engine-correlation": {
        "enabled": true,
        "minimal": 1,
        "optimal": 1
      },
      "engine-dynamic-infos": {
        "enabled": true,
        "minimal": 1,
        "optimal": 1
      }
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "queue": {
        "limit": 1000,
        "enabled": true
      },
      "messages": {
        "limit": 5000,
        "enabled": true
      },
      "engine-fifo": {
        "enabled": true,
        "minimal": 1,
        "optimal": 2
      },
      "engine-pbehavior": {
        "enabled": true,
        "minimal": 1,
        "optimal": 2
      },
      "engine-service": {
        "enabled": true,
        "minimal": 1,
        "optimal": 3
      },
      "engine-action": {
        "enabled": false,
        "minimal": 0,
        "optimal": 0
      },
      "engine-che": {
        "enabled": true,
        "minimal": 1,
        "optimal": 1
      },
      "engine-remediation": {
        "enabled": true,
        "minimal": 1,
        "optimal": 1
      },
      "engine-webhook": {
        "enabled": true,
        "minimal": 1,
        "optimal": 4
      },
      "engine-axe": {
        "enabled": true,
        "minimal": 1,
        "optimal": 2
      },
      "engine-correlation": {
        "enabled": true,
        "minimal": 1,
        "optimal": 1
      },
      "engine-dynamic-infos": {
        "enabled": true,
        "minimal": 1,
        "optimal": 1
      }
    }
    """
    When I do GET /api/v4/cat/healthcheck/parameters
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "queue": {
        "limit": 1000,
        "enabled": true
      },
      "messages": {
        "limit": 5000,
        "enabled": true
      },
      "engine-fifo": {
        "enabled": true,
        "minimal": 1,
        "optimal": 2
      },
      "engine-pbehavior": {
        "enabled": true,
        "minimal": 1,
        "optimal": 2
      },
      "engine-service": {
        "enabled": true,
        "minimal": 1,
        "optimal": 3
      },
      "engine-action": {
        "enabled": false,
        "minimal": 0,
        "optimal": 0
      },
      "engine-che": {
        "enabled": true,
        "minimal": 1,
        "optimal": 1
      },
      "engine-remediation": {
        "enabled": true,
        "minimal": 1,
        "optimal": 1
      },
      "engine-webhook": {
        "enabled": true,
        "minimal": 1,
        "optimal": 4
      },
      "engine-axe": {
        "enabled": true,
        "minimal": 1,
        "optimal": 2
      },
      "engine-correlation": {
        "enabled": true,
        "minimal": 1,
        "optimal": 1
      },
      "engine-dynamic-infos": {
        "enabled": true,
        "minimal": 1,
        "optimal": 1
      }
    }
    """

  Scenario: given update healthcheck parameters requests should return error
    When I am admin
    When I do PUT /api/v4/cat/healthcheck/parameters:
    """
    {
      "queue": {
        "limit": -1,
        "enabled": true
      },
      "messages": {
        "limit": -1,
        "enabled": true
      },
      "engine-fifo": {
        "enabled": true,
        "minimal": 1,
        "optimal": 2
      },
      "engine-pbehavior": {
        "enabled": true,
        "minimal": 1,
        "optimal": 2
      },
      "engine-service": {
        "enabled": true,
        "minimal": 1,
        "optimal": 3
      },
      "engine-action": {
        "enabled": true,
        "minimal": 1,
        "optimal": 2
      },
      "engine-che": {
        "enabled": true,
        "minimal": 1,
        "optimal": 1
      },
      "engine-remediation": {
        "enabled": true,
        "minimal": 1,
        "optimal": 1
      },
      "engine-webhook": {
        "enabled": true,
        "minimal": 1,
        "optimal": 4
      },
      "engine-axe": {
        "enabled": true,
        "minimal": 1,
        "optimal": 2
      },
      "engine-correlation": {
        "enabled": true,
        "minimal": 1,
        "optimal": 1
      },
      "engine-dynamic-infos": {
        "enabled": true,
        "minimal": 1,
        "optimal": 1
      }
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "queue.limit": "Limit should be greater than 0.",
        "messages.limit": "Limit should be greater than 0."
      }
    }
    """
    When I do PUT /api/v4/cat/healthcheck/parameters:
    """
    {
      "queue": {
        "limit": 1000,
        "enabled": true
      },
      "messages": {
        "limit": 5000,
        "enabled": true
      },
      "engine-fifo": {
        "enabled": true,
        "minimal": 1,
        "optimal": 2
      },
      "engine-pbehavior": {
        "enabled": true,
        "minimal": 1,
        "optimal": 2
      },
      "engine-service": {
        "enabled": true,
        "minimal": 1,
        "optimal": 3
      },
      "engine-action": {
        "enabled": true,
        "minimal": 2,
        "optimal": 1
      },
      "engine-che": {
        "enabled": true,
        "minimal": 1,
        "optimal": 1
      },
      "engine-remediation": {
        "enabled": true,
        "minimal": 1,
        "optimal": 1
      },
      "engine-webhook": {
        "enabled": true,
        "minimal": 1,
        "optimal": 4
      },
      "engine-axe": {
        "enabled": true,
        "minimal": 1,
        "optimal": 2
      },
      "engine-correlation": {
        "enabled": true,
        "minimal": 1,
        "optimal": 1
      },
      "engine-dynamic-infos": {
        "enabled": true,
        "minimal": 1,
        "optimal": 1
      }
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "engine-action.optimal": "Optimal should be greater or equal than Minimal."
      }
    }
    """
    When I do PUT /api/v4/cat/healthcheck/parameters:
    """
    {
      "queue": {
        "limit": 1000,
        "enabled": true
      },
      "engine-fifo": {
        "enabled": true,
        "minimal": 1,
        "optimal": 2
      },
      "engine-pbehavior": {
        "enabled": true,
        "minimal": 1,
        "optimal": 2
      },
      "engine-service": {
        "enabled": true,
        "minimal": 1,
        "optimal": 3
      },
      "engine-action": {
        "enabled": true,
        "minimal": 1,
        "optimal": 2
      },
      "engine-che": {
        "enabled": true,
        "minimal": 1,
        "optimal": 1
      },
      "engine-remediation": {
        "enabled": true,
        "minimal": 1,
        "optimal": 1
      },
      "engine-webhook": {
        "enabled": true,
        "minimal": 1,
        "optimal": 4
      },
      "engine-axe": {
        "enabled": true,
        "minimal": 1,
        "optimal": 2
      },
      "engine-correlation": {
        "enabled": true,
        "minimal": 1,
        "optimal": 1
      },
      "engine-dynamic-infos": {
        "enabled": true,
        "minimal": -1,
        "optimal": 1
      }
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "engine-dynamic-infos.minimal": "Minimal should be greater than 0."
      }
    }
    """

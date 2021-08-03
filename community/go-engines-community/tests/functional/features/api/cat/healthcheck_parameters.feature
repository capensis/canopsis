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
    {
      "queue_limit": 1000,
      "engines": {
        "engine-fifo": {
          "minimal": 1,
          "optimal": 2
        },
        "engine-pbehavior": {
          "minimal": 1,
          "optimal": 2
        },
        "engine-service": {
          "minimal": 1,
          "optimal": 3
        },
        "engine-action": {
          "minimal": 1,
          "optimal": 2
        },
        "engine-che": {
          "minimal": 1,
          "optimal": 1
        },
        "engine-remediation": {
          "minimal": 1,
          "optimal": 1
        },
        "engine-webhook": {
          "minimal": 1,
          "optimal": 4
        },
        "engine-axe": {
          "minimal": 1,
          "optimal": 2
        },
        "engine-correlation": {
          "minimal": 1,
          "optimal": 1
        },
        "engine-dynamic-infos": {
          "minimal": 1,
          "optimal": 1
        }
      }
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "queue_limit": 1000,
      "engines": {
        "engine-fifo": {
          "minimal": 1,
          "optimal": 2
        },
        "engine-pbehavior": {
          "minimal": 1,
          "optimal": 2
        },
        "engine-service": {
          "minimal": 1,
          "optimal": 3
        },
        "engine-action": {
          "minimal": 1,
          "optimal": 2
        },
        "engine-che": {
          "minimal": 1,
          "optimal": 1
        },
        "engine-remediation": {
          "minimal": 1,
          "optimal": 1
        },
        "engine-webhook": {
          "minimal": 1,
          "optimal": 4
        },
        "engine-axe": {
          "minimal": 1,
          "optimal": 2
        },
        "engine-correlation": {
          "minimal": 1,
          "optimal": 1
        },
        "engine-dynamic-infos": {
          "minimal": 1,
          "optimal": 1
        }
      }
    }
    """
    When I do GET /api/v4/cat/healthcheck/parameters
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "queue_limit": 1000,
      "engines": {
        "engine-fifo": {
          "minimal": 1,
          "optimal": 2
        },
        "engine-pbehavior": {
          "minimal": 1,
          "optimal": 2
        },
        "engine-service": {
          "minimal": 1,
          "optimal": 3
        },
        "engine-action": {
          "minimal": 1,
          "optimal": 2
        },
        "engine-che": {
          "minimal": 1,
          "optimal": 1
        },
        "engine-remediation": {
          "minimal": 1,
          "optimal": 1
        },
        "engine-webhook": {
          "minimal": 1,
          "optimal": 4
        },
        "engine-axe": {
          "minimal": 1,
          "optimal": 2
        },
        "engine-correlation": {
          "minimal": 1,
          "optimal": 1
        },
        "engine-dynamic-infos": {
          "minimal": 1,
          "optimal": 1
        }
      }
    }
    """

  Scenario: given update healthcheck parameters requests should return error
    When I am admin
    When I do PUT /api/v4/cat/healthcheck/parameters:
    """
    {
      "queue_limit": -1,
      "engines": {
        "engine-fifo": {
          "minimal": 1,
          "optimal": 2
        },
        "engine-pbehavior": {
          "minimal": 1,
          "optimal": 2
        },
        "engine-service": {
          "minimal": 1,
          "optimal": 3
        },
        "engine-action": {
          "minimal": 1,
          "optimal": 2
        },
        "engine-che": {
          "minimal": 1,
          "optimal": 1
        },
        "engine-remediation": {
          "minimal": 1,
          "optimal": 1
        },
        "engine-webhook": {
          "minimal": 1,
          "optimal": 4
        },
        "engine-axe": {
          "minimal": 1,
          "optimal": 2
        },
        "engine-correlation": {
          "minimal": 1,
          "optimal": 1
        },
        "engine-dynamic-infos": {
          "minimal": 1,
          "optimal": 1
        }
      }
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "queue_limit": "QueueLimit should be greater than 0."
      }
    }
    """
    When I do PUT /api/v4/cat/healthcheck/parameters:
    """
    {
      "queue_limit": 1000,
      "engines": {
        "engine-fifo": {
          "minimal": 1,
          "optimal": 2
        },
        "engine-pbehavior": {
          "minimal": 1,
          "optimal": 2
        },
        "engine-service": {
          "minimal": 1,
          "optimal": 3
        },
        "engine-action": {
          "minimal": 2,
          "optimal": 1
        },
        "engine-che": {
          "minimal": 1,
          "optimal": 1
        },
        "engine-remediation": {
          "minimal": 1,
          "optimal": 1
        },
        "engine-webhook": {
          "minimal": 1,
          "optimal": 4
        },
        "engine-axe": {
          "minimal": 1,
          "optimal": 2
        },
        "engine-correlation": {
          "minimal": 1,
          "optimal": 1
        },
        "engine-dynamic-infos": {
          "minimal": 1,
          "optimal": 1
        }
      }
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "engines.engine-action.optimal": "Optimal should be greater or equal than Minimal."
      }
    }
    """
    When I do PUT /api/v4/cat/healthcheck/parameters:
    """
    {
      "queue_limit": 1000,
      "engines": {
        "engine-fifo": {
          "minimal": 1,
          "optimal": 2
        },
        "engine-pbehavior": {
          "minimal": 1,
          "optimal": 2
        },
        "engine-service": {
          "minimal": 1,
          "optimal": 3
        },
        "engine-action": {
          "minimal": 1,
          "optimal": 1
        },
        "engine-che": {
          "minimal": 1,
          "optimal": 1
        },
        "engine-remediation": {
          "minimal": 1,
          "optimal": 1
        },
        "engine-webhook": {
          "minimal": 1,
          "optimal": 4
        },
        "engine-axe": {
          "minimal": 1,
          "optimal": 2
        },
        "engine-correlation": {
          "minimal": 1,
          "optimal": 1
        },
        "engine-dynamic-infos": {
          "minimal": -1,
          "optimal": 1
        }
      }
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "engines.engine-dynamic-infos.minimal": "Minimal should be greater than 0."
      }
    }
    """
    When I do PUT /api/v4/cat/healthcheck/parameters:
    """
    {
      "queue_limit": 1000,
      "engines": {
        "engine-abc": {
          "minimal": 1,
          "optimal": 1
        },
        "engine-fifo": {
          "minimal": 1,
          "optimal": 2
        },
        "engine-pbehavior": {
          "minimal": 1,
          "optimal": 2
        },
        "engine-service": {
          "minimal": 1,
          "optimal": 3
        },
        "engine-action": {
          "minimal": 1,
          "optimal": 1
        },
        "engine-che": {
          "minimal": 1,
          "optimal": 1
        },
        "engine-remediation": {
          "minimal": 1,
          "optimal": 1
        },
        "engine-webhook": {
          "minimal": 1,
          "optimal": 4
        },
        "engine-axe": {
          "minimal": 1,
          "optimal": 2
        },
        "engine-correlation": {
          "minimal": 1,
          "optimal": 1
        },
        "engine-dynamic-infos": {
          "minimal": 1,
          "optimal": 1
        }
      }
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "engines.engine-abc": "engine's name is not valid."
      }
    }
    """
    When I do PUT /api/v4/cat/healthcheck/parameters:
    """
    {
      "queue_limit": 1000,
      "engines": {
        "engine-fifo": {
          "minimal": 1,
          "optimal": 2
        },
        "engine-pbehavior": {
          "minimal": 1,
          "optimal": 2
        },
        "engine-che": {
          "minimal": 1,
          "optimal": 1
        },
        "engine-remediation": {
          "minimal": 1,
          "optimal": 1
        },
        "engine-webhook": {
          "minimal": 1,
          "optimal": 4
        },
        "engine-axe": {
          "minimal": 1,
          "optimal": 2
        },
        "engine-correlation": {
          "minimal": 1,
          "optimal": 1
        },
        "engine-dynamic-infos": {
          "minimal": 1,
          "optimal": 1
        }
      }
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "engines.engine-action": "engine is not defined",
        "engines.engine-service": "engine is not defined"
      }
    }
    """

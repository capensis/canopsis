Feature: update security config
  I need to be able to update security config
  Only admin should be able to update security config

  Scenario: given update request and no auth user should not allow access
    When I do POST /api/v4/security
    Then the response code should be 401

  Scenario: given update request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/security
    Then the response code should be 403

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/security
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/security
    Then the response code should be 403

  Scenario: given update request should return ok
    When I am admin
    When I do POST /api/v4/security:
    """json
    {
      "basic": {
        "expiration_interval": {
          "value": 1,
          "unit": "m"
        },
        "inactivity_interval": {
          "value": 8,
          "unit": "h"
        }
      },
      "ldap": {
        "expiration_interval": {
          "value": 1,
          "unit": "m"
        },
        "inactivity_interval": {
          "value": 8,
          "unit": "h"
        }
      }
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "basic": {
        "expiration_interval": {
          "value": 1,
          "unit": "m"
        },
        "inactivity_interval": {
          "value": 8,
          "unit": "h"
        }
      }
    }
    """
    When I do GET /api/v4/security
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "basic": {
        "expiration_interval": {
          "value": 1,
          "unit": "m"
        },
        "inactivity_interval": {
          "value": 8,
          "unit": "h"
        }
      }
    }
    """

  Scenario: given empty update request should return ok
    When I am admin
    When I do POST /api/v4/security:
    """json
    {}
    """
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "basic": {
        "expiration_interval": null,
        "inactivity_interval": {
          "value": 24,
          "unit": "h"
        }
      }
    }
    """
    When I do GET /api/v4/security
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "basic": {
        "expiration_interval": null,
        "inactivity_interval": {
          "value": 24,
          "unit": "h"
        }
      }
    }
    """

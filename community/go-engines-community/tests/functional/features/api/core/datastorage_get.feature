Feature: Get and update data storage config
  I need to be able to get and update data storage config

  Scenario: given update request should return updated config
    When I am admin
    When I do PUT /api/v4/data-storage:
    """
    {
      "junit": {
        "delete_after": {
          "enabled": true,
          "seconds": 864000,
          "unit": "d"
        }
      }
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "config": {
        "junit": {
          "delete_after": {
            "enabled": true,
            "seconds": 864000,
            "unit": "d"
          }
        }
      },
      "history": {
        "junit": null
      }
    }
    """
    When I do GET /api/v4/data-storage
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "config": {
        "junit": {
          "delete_after": {
            "enabled": true,
            "seconds": 864000,
            "unit": "d"
          }
        }
      },
      "history": {
        "junit": null
      }
    }
    """
    When I do PUT /api/v4/data-storage:
    """
    {}
    """
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "config": {
        "junit": {
          "delete_after": null
        }
      },
      "history": {
        "junit": null
      }
    }
    """
    When I do GET /api/v4/data-storage
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "config": {
        "junit": {
          "delete_after": null
        }
      },
      "history": {
        "junit": null
      }
    }
    """

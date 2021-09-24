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
      },
      "remediation": {
        "accumulate_after": {
          "enabled": true,
          "seconds": 864000,
          "unit": "d"
        },
        "delete_after": {
          "enabled": true,
          "seconds": 1728000,
          "unit": "d"
        }
      },
      "alarm": {
        "archive_after": {
          "enabled": true,
          "seconds": 864000,
          "unit": "d"
        },
        "delete_after": {
          "enabled": true,
          "seconds": 1728000,
          "unit": "d"
        }
      },
      "pbehavior": {
        "delete_after": {
          "enabled": true,
          "seconds": 1728000,
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
        },
        "remediation": {
          "accumulate_after": {
            "enabled": true,
            "seconds": 864000,
            "unit": "d"
          },
          "delete_after": {
            "enabled": true,
            "seconds": 1728000,
            "unit": "d"
          }
        },
        "alarm": {
          "archive_after": {
            "enabled": true,
            "seconds": 864000,
            "unit": "d"
          },
          "delete_after": {
            "enabled": true,
            "seconds": 1728000,
            "unit": "d"
          }
        },
        "pbehavior": {
          "delete_after": {
            "enabled": true,
            "seconds": 1728000,
            "unit": "d"
          }
        }
      },
      "history": {
        "junit": null,
        "remediation": null,
        "alarm": null,
        "entity": null,
        "pbehavior": null
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
        },
        "remediation": {
          "accumulate_after": {
            "enabled": true,
            "seconds": 864000,
            "unit": "d"
          },
          "delete_after": {
            "enabled": true,
            "seconds": 1728000,
            "unit": "d"
          }
        },
        "alarm": {
          "archive_after": {
            "enabled": true,
            "seconds": 864000,
            "unit": "d"
          },
          "delete_after": {
            "enabled": true,
            "seconds": 1728000,
            "unit": "d"
          }
        },
        "pbehavior": {
          "delete_after": {
            "enabled": true,
            "seconds": 1728000,
            "unit": "d"
          }
        }
      },
      "history": {
        "junit": null,
        "remediation": null,
        "alarm": null,
        "entity": null,
        "pbehavior": null
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
        },
        "remediation": {
          "accumulate_after": null,
          "delete_after": null
        },
        "alarm": {
          "archive_after": null,
          "delete_after": null
        },
        "pbehavior": {
          "delete_after": null
        }
      },
      "history": {
        "junit": null,
        "remediation": null,
        "alarm": null,
        "entity": null,
        "pbehavior": null
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
        },
        "remediation": {
          "accumulate_after": null,
          "delete_after": null
        },
        "alarm": {
          "archive_after": null,
          "delete_after": null
        },
        "pbehavior": {
          "delete_after": null
        }
      },
      "history": {
        "junit": null,
        "remediation": null,
        "alarm": null,
        "entity": null,
        "pbehavior": null
      }
    }
    """

  Scenario: given update request should return validation error
    When I am admin
    When I do PUT /api/v4/data-storage:
    """
    {
      "remediation": {
        "accumulate_after": {
          "seconds": 864000,
          "unit": "d"
        },
        "delete_after": {
          "seconds": 864000,
          "unit": "d"
        }
      }
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "remediation.delete_after": "DeleteAfter should be greater than AccumulateAfter."
      }
    }
    """
    When I do PUT /api/v4/data-storage:
    """
    {
      "alarm": {
        "archive_after": {
          "seconds": 864000,
          "unit": "d"
        },
        "delete_after": {
          "seconds": 864000,
          "unit": "d"
        }
      }
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "alarm.delete_after": "DeleteAfter should be greater than ArchiveAfter."
      }
    }
    """
    When I do PUT /api/v4/data-storage:
    """
    {
      "alarm": {
        "delete_after": {
          "enabled": true,
          "seconds": 864000,
          "unit": "d"
        }
      }
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "alarm.archive_after": "ArchiveAfter is required when DeleteAfter is defined."
      }
    }
    """

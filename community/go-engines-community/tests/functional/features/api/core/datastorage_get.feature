Feature: Get and update data storage config
  I need to be able to get and update data storage config

  Scenario: given update request should return updated config
    When I am admin
    When I do PUT /api/v4/data-storage:
    """json
    {
      "junit": {
        "delete_after": {
          "enabled": true,
          "value": 10,
          "unit": "d"
        }
      },
      "remediation": {
        "accumulate_after": {
          "enabled": true,
          "value": 10,
          "unit": "d"
        },
        "delete_after": {
          "enabled": true,
          "value": 20,
          "unit": "d"
        }
      },
      "alarm": {
        "archive_after": {
          "enabled": true,
          "value": 10,
          "unit": "d"
        },
        "delete_after": {
          "enabled": true,
          "value": 20,
          "unit": "d"
        }
      },
      "pbehavior": {
        "delete_after": {
          "enabled": true,
          "value": 20,
          "unit": "d"
        }
      },
      "health_check": {
        "delete_after": {
          "enabled": true,
          "value": 20,
          "unit": "d"
        }
      },
      "webhook": {
        "log_credentials": true,
        "delete_after": {
          "enabled": true,
          "value": 20,
          "unit": "d"
        }
      }
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "config": {
        "junit": {
          "delete_after": {
            "enabled": true,
            "value": 10,
            "unit": "d"
          }
        },
        "remediation": {
          "accumulate_after": {
            "enabled": true,
            "value": 10,
            "unit": "d"
          },
          "delete_after": {
            "enabled": true,
            "value": 20,
            "unit": "d"
          }
        },
        "alarm": {
          "archive_after": {
            "enabled": true,
            "value": 10,
            "unit": "d"
          },
          "delete_after": {
            "enabled": true,
            "value": 20,
            "unit": "d"
          }
        },
        "pbehavior": {
          "delete_after": {
            "enabled": true,
            "value": 20,
            "unit": "d"
          }
        },
        "health_check": {
          "delete_after": {
            "enabled": true,
            "value": 20,
            "unit": "d"
          }
        },
        "webhook": {
          "log_credentials": true,
          "delete_after": {
            "enabled": true,
            "value": 20,
            "unit": "d"
          }
        }
      },
      "history": {
        "junit": null,
        "remediation": null,
        "alarm": null,
        "entity": null,
        "pbehavior": null,
        "health_check": null,
        "webhook": null
      }
    }
    """
    When I do GET /api/v4/data-storage
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "config": {
        "junit": {
          "delete_after": {
            "enabled": true,
            "value": 10,
            "unit": "d"
          }
        },
        "remediation": {
          "accumulate_after": {
            "enabled": true,
            "value": 10,
            "unit": "d"
          },
          "delete_after": {
            "enabled": true,
            "value": 20,
            "unit": "d"
          }
        },
        "alarm": {
          "archive_after": {
            "enabled": true,
            "value": 10,
            "unit": "d"
          },
          "delete_after": {
            "enabled": true,
            "value": 20,
            "unit": "d"
          }
        },
        "pbehavior": {
          "delete_after": {
            "enabled": true,
            "value": 20,
            "unit": "d"
          }
        },
        "health_check": {
          "delete_after": {
            "enabled": true,
            "value": 20,
            "unit": "d"
          }
        },
        "webhook": {
          "log_credentials": true,
          "delete_after": {
            "enabled": true,
            "value": 20,
            "unit": "d"
          }
        }
      },
      "history": {
        "junit": null,
        "remediation": null,
        "alarm": null,
        "entity": null,
        "pbehavior": null,
        "health_check": null,
        "webhook": null
      }
    }
    """
    When I do PUT /api/v4/data-storage:
    """json
    {}
    """
    Then the response code should be 200
    Then the response body should be:
    """json
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
        },
        "health_check": {
          "delete_after": null
        },
        "webhook": {
          "log_credentials": false,
          "delete_after": null
        }
      },
      "history": {
        "junit": null,
        "remediation": null,
        "alarm": null,
        "entity": null,
        "pbehavior": null,
        "health_check": null,
        "webhook": null
      }
    }
    """
    When I do GET /api/v4/data-storage
    Then the response code should be 200
    Then the response body should be:
    """json
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
        },
        "health_check": {
          "delete_after": null
        },
        "webhook": {
          "log_credentials": false,
          "delete_after": null
        }
      },
      "history": {
        "junit": null,
        "remediation": null,
        "alarm": null,
        "entity": null,
        "pbehavior": null,
        "health_check": null,
        "webhook": null
      }
    }
    """

  Scenario: given update request should return validation error
    When I am admin
    When I do PUT /api/v4/data-storage:
    """json
    {
      "remediation": {
        "accumulate_after": {
          "value": 10,
          "unit": "d",
          "enabled": true
        },
        "delete_after": {
          "value": 10,
          "unit": "d",
          "enabled": true
        }
      }
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "remediation.delete_after": "DeleteAfter should be greater than AccumulateAfter."
      }
    }
    """
    When I do PUT /api/v4/data-storage:
    """json
    {
      "alarm": {
        "archive_after": {
          "value": 10,
          "unit": "d",
          "enabled": true
        },
        "delete_after": {
          "value": 10,
          "unit": "d",
          "enabled": true
        }
      }
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "alarm.delete_after": "DeleteAfter should be greater than ArchiveAfter."
      }
    }
    """
    When I do PUT /api/v4/data-storage:
    """json
    {
      "alarm": {
        "delete_after": {
          "enabled": true,
          "value": 10,
          "unit": "d"
        }
      }
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "alarm.archive_after": "ArchiveAfter is required when DeleteAfter is defined."
      }
    }
    """

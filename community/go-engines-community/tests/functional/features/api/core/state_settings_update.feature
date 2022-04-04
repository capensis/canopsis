Feature: update state settings
  I need to be able to update state settings
  Only admin should be able to update state settings

  Scenario: PUT unauthorized
    When I do PUT /api/v4/state-settings/junit
    Then the response code should be 401

  Scenario: PUT without permissions
    When I am noperms
    When I do PUT /api/v4/state-settings/junit

  Scenario: PUT success
    When I am admin
    When I do PUT /api/v4/state-settings/junit:
    """
    {
      "type": "junit",
      "method": "worst_of_share",
      "junit_thresholds": {
        "skipped": {
          "minor": 15,
          "major": 20,
          "critical": 30,
          "type": 1
        },
        "errors": {
          "minor": 10,
          "major": 20,
          "critical": 30,
          "type": 1
        },
        "failures": {
          "minor": 10,
          "major": 20,
          "critical": 30,
          "type": 0
        }
      }
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "_id": "junit",
      "type": "junit",
      "method": "worst_of_share",
      "junit_thresholds": {
        "skipped": {
          "minor": 15,
          "major": 20,
          "critical": 30,
          "type": 1
        },
        "errors": {
          "minor": 10,
          "major": 20,
          "critical": 30,
          "type": 1
        },
        "failures": {
          "minor": 10,
          "major": 20,
          "critical": 30,
          "type": 0
        }
      }
    }
    """

  Scenario: PUT with minor > major
    When I am admin
    When I do PUT /api/v4/state-settings/junit:
    """
    {
      "type": "junit",
      "method": "worst_of_share",
      "junit_thresholds": {
        "skipped": {
          "minor": 25,
          "major": 20,
          "critical": 30,
          "type": 0
        },
        "errors": {
          "minor": 10,
          "major": 20,
          "critical": 30,
          "type": 1
        },
        "failures": {
          "minor": 10,
          "major": 20,
          "critical": 30,
          "type": 1
        }
      }
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "junit_thresholds.skipped.minor": "Minor should be less or equal than Major."
      }
    }
    """

  Scenario: PUT with major > critical
    When I am admin
    When I do PUT /api/v4/state-settings/junit:
    """
    {
      "type": "junit",
      "method": "worst_of_share",
      "junit_thresholds": {
        "skipped": {
          "minor": 10,
          "major": 35,
          "critical": 30,
          "type": 1
        },
        "errors": {
          "minor": 10,
          "major": 20,
          "critical": 30,
          "type": 1
        },
        "failures": {
          "minor": 10,
          "major": 20,
          "critical": 30,
          "type": 1
        }
      }
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "junit_thresholds.skipped.major": "Major should be less or equal than Critical."
      }
    }
    """

  Scenario: PUT worst_of_share without thresholds
    When I am admin
    When I do PUT /api/v4/state-settings/junit:
    """
    {
      "type": "junit",
      "method": "worst_of_share"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
          "junit_thresholds": "junit_thresholds should not be blank."
      }
    }
    """

  Scenario: PUT worst with thresholds
    When I am admin
    When I do PUT /api/v4/state-settings/junit:
    """
    {
      "type": "junit",
      "method": "worst",
      "junit_thresholds": {
        "skipped": {
          "minor": 10,
          "major": 20,
          "critical": 30,
          "type": 1
        },
        "errors": {
          "minor": 10,
          "major": 20,
          "critical": 30,
          "type": 1
        },
        "failures": {
          "minor": 10,
          "major": 20,
          "critical": 30,
          "type": 1
        }
      }
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "junit_thresholds": "junit_thresholds is not empty."
      }
    }
    """

  Scenario: PUT wrong setting type
    When I am admin
    When I do PUT /api/v4/state-settings/junit:
    """
    {
      "type": "junit123"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "type": "Type must be one of [junit]."
      }
    }
    """

  Scenario: PUT wrong setting method
    When I am admin
    When I do PUT /api/v4/state-settings/junit:
    """
    {
      "type": "junit",
      "method": "worst_of_share123"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "method": "Method must be one of [worst,worst_of_share]."
      }
    }
    """

  Scenario: PUT wrong threshold type
    When I am admin
    When I do PUT /api/v4/state-settings/junit:
    """
    {
      "type": "junit",
      "method": "worst_of_share",
      "junit_thresholds": {
        "skipped": {
          "minor": 10,
          "major": 20,
          "critical": 30,
          "type": 3
        },
        "errors": {
          "minor": 10,
          "major": 20,
          "critical": 30,
          "type": 1
        },
        "failures": {
          "minor": 10,
          "major": 20,
          "critical": 30
        }
      }
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "junit_thresholds.skipped.type": "Type must be one of [0,1].",
        "junit_thresholds.failures.type": "Type is missing."
      }
    }
    """

  Scenario: PUT zero state settings
    When I am admin
    When I do PUT /api/v4/state-settings/junit:
    """
    {
      "type": "junit",
      "method": "worst_of_share",
      "junit_thresholds": {
        "skipped": {
          "minor": 0,
          "major": 20,
          "critical": 30,
          "type": 1
        },
        "errors": {
          "minor": 0,
          "major": 20,
          "critical": 30,
          "type": 1
        },
        "failures": {
          "minor": 10,
          "major": 20,
          "critical": 30,
          "type": 0
        }
      }
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "_id": "junit",
      "type": "junit",
      "method": "worst_of_share",
      "junit_thresholds": {
        "skipped": {
          "minor": 0,
          "major": 20,
          "critical": 30,
          "type": 1
        },
        "errors": {
          "minor": 0,
          "major": 20,
          "critical": 30,
          "type": 1
        },
        "failures": {
          "minor": 10,
          "major": 20,
          "critical": 30,
          "type": 0
        }
      }
    }
    """

  Scenario: PUT settings without state value
    When I am admin
    When I do PUT /api/v4/state-settings/junit:
    """
    {
      "type": "junit",
      "method": "worst_of_share",
      "junit_thresholds": {
        "skipped": {
          "major": 20,
          "critical": 30,
          "type": 1
        },
        "errors": {
          "minor": 0,
          "major": 20,
          "critical": 30,
          "type": 1
        },
        "failures": {
          "major": 20,
          "critical": 30,
          "type": 0
        }
      }
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "junit_thresholds.failures.minor": "Minor is missing.",
        "junit_thresholds.skipped.minor": "Minor is missing."
      }
    }
    """

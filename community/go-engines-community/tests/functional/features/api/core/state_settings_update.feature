Feature: update state settings
  I need to be able to update state settings
  Only admin should be able to update state settings

  @concurrent
  Scenario: given update request and no auth user should not allow access
    When I do PUT /api/v4/state-settings/junit
    Then the response code should be 401

  @concurrent
  Scenario: given update request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/state-settings/junit
    Then the response code should be 403

  @concurrent
  Scenario: given update junit state_settings request should return ok
    When I am admin
    When I do PUT /api/v4/state-settings/junit:
    """json
    {
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
          "minor": 0,
          "major": 20,
          "critical": 30,
          "type": 0
        }
      }
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "junit",
      "title": "Junit",
      "method": "worst_of_share",
      "enabled": true,
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
          "minor": 0,
          "major": 20,
          "critical": 30,
          "type": 0
        }
      },
      "editable": true,
      "deletable": false
    }
    """
    When I do GET /api/v4/state-settings/junit
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "junit",
      "title": "Junit",
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
          "minor": 0,
          "major": 20,
          "critical": 30,
          "type": 0
        }
      },
      "editable": true,
      "deletable": false
    }
    """

  @concurrent
  Scenario: given update junit state_settings request with invalid fields should return error
    When I am admin
    When I do PUT /api/v4/state-settings/junit:
    """json
    {
      "_id": "junit",
      "priority": 1,
      "title": "some",
      "method": "worst_of_share",
      "type": "component",
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-service-state-settings-to-update-1"
            }
          }
        ]
      ],
      "inherited_entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-state-settings-to-update-1"
            }
          }
        ]
      ],
      "state_thresholds": {
        "critical": {
          "method": "share",
          "state": "major",
          "cond": "gt",
          "value": 20
        },
        "major": {
          "method": "number",
          "state": "major",
          "cond": "lt",
          "value": 10
        },
        "minor": {
          "method": "number",
          "state": "ok",
          "cond": "lt",
          "value": 10
        },
        "ok": {
          "method": "share",
          "state": "minor",
          "cond": "gt",
          "value": 10
        }
      }
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "title": "Title is not empty.",
        "inherited_entity_pattern": "InheritedEntityPattern is not empty.",
        "junit_thresholds": "JUnitThresholds should not be blank.",
        "entity_pattern": "EntityPattern is not empty.",
        "state_thresholds": "StateThresholds is not empty.",
        "priority": "Priority is not empty."
      }
    }
    """

  @concurrent
  Scenario: given update junit state_settings request with invalid ratio between thresholds should return error
    When I am admin
    When I do PUT /api/v4/state-settings/junit:
    """json
    {
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
          "major": 35,
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
    Then the response body should contain:
    """json
    {
      "errors": {
        "junit_thresholds.skipped.minor": "Minor should be less or equal than Major.",
        "junit_thresholds.errors.major": "Major should be less or equal than Critical."
      }
    }
    """

  @concurrent
  Scenario: given update junit state_settings request with worst method and thresholds should return error
    When I am admin
    When I do PUT /api/v4/state-settings/junit:
    """json
    {
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
    Then the response body should contain:
    """json
    {
      "errors": {
        "junit_thresholds": "JUnitThresholds is not empty."
      }
    }
    """

  @concurrent
  Scenario: given update junit state_settings request with invalid method should return error
    When I am admin
    When I do PUT /api/v4/state-settings/junit:
    """json
    {
      "method": "worst_of_share123"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "method": "Method must be one of [worst,worst_of_share]."
      }
    }
    """

  @concurrent
  Scenario: given update junit state_settings request with invalid threshold type should return error
    When I am admin
    When I do PUT /api/v4/state-settings/junit:
    """json
    {
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
    Then the response body should contain:
    """json
    {
      "errors": {
        "junit_thresholds.skipped.type": "Type must be one of [0,1].",
        "junit_thresholds.failures.type": "Type is missing."
      }
    }
    """

  @concurrent
  Scenario: given update junit state_settings request with absent threshold values should return error
    When I am admin
    When I do PUT /api/v4/state-settings/junit:
    """json
    {
      "method": "worst_of_share",
      "junit_thresholds": {
        "skipped": {
          "major": 20,
          "critical": 30,
          "type": 1
        },
        "errors": {
          "minor": 0,
          "critical": 30,
          "type": 1
        },
        "failures": {
          "minor": 10,
          "major": 20,
          "type": 0
        }
      }
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "junit_thresholds.skipped.minor": "Minor is missing.",
        "junit_thresholds.errors.major": "Major is missing.",
        "junit_thresholds.failures.critical": "Critical is missing."
      }
    }
    """

  @concurrent
  Scenario: given update junit state_settings request with absent threshold values should return error
    When I am admin
    When I do PUT /api/v4/state-settings/service:
    """json
    {
      "method": "worst_of_share"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "error": "can't modify service state settings"
    }
    """

  @concurrent
  Scenario: given update inherited state_settings request should return ok
    When I am admin
    When I do PUT /api/v4/state-settings/inherited-settings-to-update-1:
    """json
    {
      "title": "inherited-settings-to-update-1-title-updated",
      "method": "inherited",
      "enabled": true,
      "type": "component",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-service-state-settings-to-update-1"
            }
          }
        ]
      ],
      "inherited_entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-state-settings-to-update-1"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "inherited-settings-to-update-1",
      "title": "inherited-settings-to-update-1-title-updated",
      "method": "inherited",
      "enabled": true,
      "type": "component",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-service-state-settings-to-update-1"
            }
          }
        ]
      ],
      "inherited_entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-state-settings-to-update-1"
            }
          }
        ]
      ],
      "editable": true,
      "deletable": true
    }
    """
    When I do GET /api/v4/state-settings/inherited-settings-to-update-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "inherited-settings-to-update-1",
      "title": "inherited-settings-to-update-1-title-updated",
      "method": "inherited",
      "enabled": true,
      "type": "component",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-service-state-settings-to-update-1"
            }
          }
        ]
      ],
      "inherited_entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-state-settings-to-update-1"
            }
          }
        ]
      ],
      "editable": true,
      "deletable": true
    }
    """

  @concurrent
  Scenario: given update inherited state_settings request with invalid fields should return error
    When I am admin
    When I do PUT /api/v4/state-settings/inherited-settings-to-update-1:
    """json
    {
      "priority": -1,
      "method": "inherited",
      "state_thresholds": {
        "critical": {
          "method": "share",
          "state": "major",
          "cond": "gt",
          "value": 20
        },
        "major": {
          "method": "number",
          "state": "major",
          "cond": "lt",
          "value": 10
        },
        "minor": {
          "method": "number",
          "state": "ok",
          "cond": "lt",
          "value": 10
        },
        "ok": {
          "method": "share",
          "state": "minor",
          "cond": "gt",
          "value": 10
        }
      },
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
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "title": "Title is missing.",
        "enabled": "Enabled is missing.",
        "entity_pattern": "EntityPattern is missing.",
        "inherited_entity_pattern": "InheritedEntityPattern is missing.",
        "junit_thresholds": "JUnitThresholds is not empty.",
        "state_thresholds": "StateThresholds is not empty.",
        "priority": "Priority should be 0 or more.",
        "type": "Type is required when Method inherited is defined."
      }
    }
    """

  @concurrent
  Scenario: given update state_settings request with invalid patterns should return error
    When I am admin
    When I do PUT /api/v4/state-settings/inherited-state-settings-1:
    """json
    {
      "method": "inherited",
      "enabled": true,
      "type": "component",
      "entity_pattern": [
        [
          {
            "field": "some",
            "cond": {
              "type": "eq",
              "value": "test-service-state-settings-to-update-1"
            }
          }
        ]
      ],
      "inherited_entity_pattern": [
        [
          {
            "field": "some",
            "cond": {
              "type": "eq",
              "value": "test-resource-state-settings-to-update-1"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "entity_pattern": "EntityPattern is invalid entity pattern.",
        "inherited_entity_pattern": "InheritedEntityPattern is invalid entity pattern."
      }
    }
    """
    When I do PUT /api/v4/state-settings/inherited-state-settings-1:
    """json
    {
      "method": "inherited",
      "enabled": true,
      "type": "component",
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-service-state-settings-to-update-1"
            }
          }
        ]
      ],
      "inherited_entity_pattern": [
        [
          {
            "field": "type",
            "cond": {
              "type": "eq",
              "value": "resource"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "entity_pattern": "EntityPattern is invalid entity pattern.",
        "inherited_entity_pattern": "InheritedEntityPattern is invalid entity pattern."
      }
    }
    """

  @concurrent
  Scenario: given update dependencies state_settings request should return ok
    When I am admin
    When I do PUT /api/v4/state-settings/dependencies-settings-to-update-1:
    """json
    {
      "title": "dependencies-settings-to-update-1-title-updated",
      "method": "dependencies",
      "enabled": true,
      "type": "component",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-service-state-settings-to-update-2"
            }
          }
        ]
      ],
      "state_thresholds": {
        "critical": {
          "method": "share",
          "state": "major",
          "cond": "gt",
          "value": 20
        },
        "major": {
          "method": "number",
          "state": "major",
          "cond": "lt",
          "value": 10
        },
        "minor": {
          "method": "number",
          "state": "ok",
          "cond": "lt",
          "value": 10
        },
        "ok": {
          "method": "share",
          "state": "minor",
          "cond": "gt",
          "value": 10
        }
      }
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "dependencies-settings-to-update-1",
      "title": "dependencies-settings-to-update-1-title-updated",
      "method": "dependencies",
      "enabled": true,
      "type": "component",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-service-state-settings-to-update-2"
            }
          }
        ]
      ],
      "state_thresholds": {
        "critical": {
          "method": "share",
          "state": "major",
          "cond": "gt",
          "value": 20
        },
        "major": {
          "method": "number",
          "state": "major",
          "cond": "lt",
          "value": 10
        },
        "minor": {
          "method": "number",
          "state": "ok",
          "cond": "lt",
          "value": 10
        },
        "ok": {
          "method": "share",
          "state": "minor",
          "cond": "gt",
          "value": 10
        }
      },
      "editable": true,
      "deletable": true
    }
    """
    When I do GET /api/v4/state-settings/dependencies-settings-to-update-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "dependencies-settings-to-update-1",
      "title": "dependencies-settings-to-update-1-title-updated",
      "method": "dependencies",
      "enabled": true,
      "type": "component",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-service-state-settings-to-update-2"
            }
          }
        ]
      ],
      "state_thresholds": {
        "critical": {
          "method": "share",
          "state": "major",
          "cond": "gt",
          "value": 20
        },
        "major": {
          "method": "number",
          "state": "major",
          "cond": "lt",
          "value": 10
        },
        "minor": {
          "method": "number",
          "state": "ok",
          "cond": "lt",
          "value": 10
        },
        "ok": {
          "method": "share",
          "state": "minor",
          "cond": "gt",
          "value": 10
        }
      },
      "editable": true,
      "deletable": true
    }
    """
    
  @concurrent
  Scenario: given update dependencies state_settings request with invalid fields should return error
    When I am admin
    When I do PUT /api/v4/state-settings/dependencies-settings-to-update-1:
    """json
    {
      "priority": -1,
      "method": "dependencies",
      "inherited_entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-state-settings-to-update-2"
            }
          }
        ]
      ],
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
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "title": "Title is missing.",
        "enabled": "Enabled is missing.",
        "entity_pattern": "EntityPattern is missing.",
        "inherited_entity_pattern": "InheritedEntityPattern is not empty.",
        "junit_thresholds": "JUnitThresholds is not empty.",
        "state_thresholds": "StateThresholds is missing.",
        "priority": "Priority should be 0 or more.",
        "type": "Type is required when Method dependencies is defined."
      }
    }
    """    
    
  @concurrent
  Scenario: given update dependencies state_settings request with invalid state_thresholds fields should return error
    When I am admin
    When I do PUT /api/v4/state-settings/dependencies-settings-to-update-1:
    """json
    {
      "method": "dependencies",
      "enabled": true,
      "type": "component",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-service-state-settings-invalid"
            }
          }
        ]
      ],
      "state_thresholds": {
        "critical": {
          "method": "some",
          "state": "some",
          "cond": "some",
          "value": 10
        },
        "major": {
          "method": "number",
          "state": "major",
          "cond": "gt",
          "value": -10
        },
        "minor": {
          "method": "share",
          "state": "minor",
          "cond": "gt",
          "value": 100
        }
      }
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "state_thresholds.critical.cond": "Cond must be one of [gt lt].",
        "state_thresholds.critical.method": "Method must be one of [number share].",
        "state_thresholds.critical.state": "State must be one of [critical major minor ok].",
        "state_thresholds.major.value": "Value should be greater than -1.",
        "state_thresholds.minor.value": "Value should be less than 100."
      }
    }
    """

  @concurrent
  Scenario: given update inherited state_settings change method request should clean previous method fields
    When I am admin
    When I do PUT /api/v4/state-settings/inherited-settings-to-update-2:
    """json
    {
      "title": "inherited-settings-to-update-2-title-updated",
      "method": "dependencies",
      "enabled": true,
      "type": "service",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-service-state-settings-to-update-3"
            }
          }
        ]
      ],
      "state_thresholds": {
        "critical": {
          "method": "share",
          "state": "major",
          "cond": "gt",
          "value": 20
        },
        "major": {
          "method": "number",
          "state": "major",
          "cond": "lt",
          "value": 10
        },
        "minor": {
          "method": "number",
          "state": "ok",
          "cond": "lt",
          "value": 10
        },
        "ok": {
          "method": "share",
          "state": "minor",
          "cond": "gt",
          "value": 10
        }
      }
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "inherited-settings-to-update-2",
      "title": "inherited-settings-to-update-2-title-updated",
      "method": "dependencies",
      "enabled": true,
      "type": "service",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-service-state-settings-to-update-3"
            }
          }
        ]
      ],
      "state_thresholds": {
        "critical": {
          "method": "share",
          "state": "major",
          "cond": "gt",
          "value": 20
        },
        "major": {
          "method": "number",
          "state": "major",
          "cond": "lt",
          "value": 10
        },
        "minor": {
          "method": "number",
          "state": "ok",
          "cond": "lt",
          "value": 10
        },
        "ok": {
          "method": "share",
          "state": "minor",
          "cond": "gt",
          "value": 10
        }
      }
    }
    """
    Then the response key "inherited_entity_pattern" should not exist
    When I do GET /api/v4/state-settings/inherited-settings-to-update-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "inherited-settings-to-update-2",
      "title": "inherited-settings-to-update-2-title-updated",
      "method": "dependencies",
      "enabled": true,
      "type": "service",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-service-state-settings-to-update-3"
            }
          }
        ]
      ],
      "state_thresholds": {
        "critical": {
          "method": "share",
          "state": "major",
          "cond": "gt",
          "value": 20
        },
        "major": {
          "method": "number",
          "state": "major",
          "cond": "lt",
          "value": 10
        },
        "minor": {
          "method": "number",
          "state": "ok",
          "cond": "lt",
          "value": 10
        },
        "ok": {
          "method": "share",
          "state": "minor",
          "cond": "gt",
          "value": 10
        }
      }
    }
    """
    Then the response key "inherited_entity_pattern" should not exist

  @concurrent
  Scenario: given update dependencies state_settings change method request should clean previous method fields
    When I am admin
    When I do PUT /api/v4/state-settings/dependencies-settings-to-update-2:
    """json
    {
      "method": "inherited",
      "title": "dependencies-settings-to-update-2-title-updated",
      "enabled": true,
      "type": "service",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-service-state-settings-to-update-4"
            }
          }
        ]
      ],
      "inherited_entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-state-settings-to-update-4"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "dependencies-settings-to-update-2",
      "title": "dependencies-settings-to-update-2-title-updated",
      "method": "inherited",
      "enabled": true,
      "type": "service",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-service-state-settings-to-update-4"
            }
          }
        ]
      ],
      "inherited_entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-state-settings-to-update-4"
            }
          }
        ]
      ]
    }
    """
    Then the response key "state_thresholds" should not exist
    When I do GET /api/v4/state-settings/dependencies-settings-to-update-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "dependencies-settings-to-update-2",
      "title": "dependencies-settings-to-update-2-title-updated",
      "method": "inherited",
      "enabled": true,
      "type": "service",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-service-state-settings-to-update-4"
            }
          }
        ]
      ],
      "inherited_entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-state-settings-to-update-4"
            }
          }
        ]
      ]
    }
    """
    Then the response key "state_thresholds" should not exist

  @concurrent
  Scenario: given dependencies state_settings request with already exist title should return error
    When I am admin
    When I do PUT /api/v4/state-settings/dependencies-settings-to-update-2:
    """json
    {
      "method": "inherited",
      "title": "inherited-settings-to-unique-title",
      "enabled": true,
      "type": "component",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-service-state-settings-to-update-4"
            }
          }
        ]
      ],
      "inherited_entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-state-settings-to-update-4"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "title": "Title already exists."
      }
    }
    """

  @concurrent
  Scenario: given update not found state_settings request should return not found error
    When I am admin
    When I do PUT /api/v4/state-settings/inherited-settings-to-update-not-found:
    """json
    {
      "title": "inherited-settings-to-update-1-title-not-found",
      "method": "inherited",
      "enabled": true,
      "type": "component",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-service-state-settings-to-update-not-found"
            }
          }
        ]
      ],
      "inherited_entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-state-settings-to-update-not-found"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 404
    Then the response body should contain:
    """json
    {
      "error": "Not found"
    }
    """


  @concurrent
  Scenario: given update state_settings request with invalid type for inherited method should return error
    When I am admin
    When I do PUT /api/v4/state-settings/dependencies-settings-to-update-2:
    """json
    {
      "method": "inherited",
      "type": "test"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "type": "Type must be one of [component service]."
      }
    }
    """

  @concurrent
  Scenario: given update state_settings request with invalid type for dependencies method should return error
    When I am admin
    When I do PUT /api/v4/state-settings/dependencies-settings-to-update-2:
    """json
    {
      "method": "dependencies",
      "type": "test"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "type": "Type must be one of [component service]."
      }
    }
    """

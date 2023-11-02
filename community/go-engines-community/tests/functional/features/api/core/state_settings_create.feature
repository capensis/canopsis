Feature: create state settings
  I need to be able to create state settings
  Only admin should be able to create state settings

  @concurrent
  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/state-settings
    Then the response code should be 401

  @concurrent
  Scenario: given create request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/state-settings
    Then the response code should be 403

  @concurrent
  Scenario: given create inherited state_settings request should return ok
    When I am admin
    When I do POST /api/v4/state-settings:
    """json
    {
      "title": "test-inherited-state-setting-1",
      "method": "inherited",
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-state-settings-to-create-1"
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
              "value": "test-resource-state-settings-to-create-1"
            }
          }
        ]
      ],
      "editable": true,
      "deletable": true
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "title": "test-inherited-state-setting-1",
      "method": "inherited",
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-state-settings-to-create-1"
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
              "value": "test-resource-state-settings-to-create-1"
            }
          }
        ]
      ]
    }
    """
    When I do GET /api/v4/state-settings/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "title": "test-inherited-state-setting-1",
      "method": "inherited",
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-state-settings-to-create-1"
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
              "value": "test-resource-state-settings-to-create-1"
            }
          }
        ]
      ]
    }
    """

  @concurrent
  Scenario: given create dependencies state_settings request should return ok
    When I am admin
    When I do POST /api/v4/state-settings:
    """json
    {
      "title": "test-dependencies-state-setting-1",
      "method": "dependencies",
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-state-settings-to-create-2"
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
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "title": "test-dependencies-state-setting-1",
      "method": "dependencies",
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-state-settings-to-create-2"
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
    When I do GET /api/v4/state-settings/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "title": "test-dependencies-state-setting-1",
      "method": "dependencies",
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-state-settings-to-create-2"
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

  @concurrent
  Scenario: given create inherited state_settings request with invalid fields should return error
    When I am admin
    When I do POST /api/v4/state-settings:
    """json
    {
      "method": "inherited",
      "priority": -1,
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
        "enabled": "Enabled is missing.",
        "entity_pattern": "EntityPattern is missing.",
        "inherited_entity_pattern": "InheritedEntityPattern is missing.",
        "junit_thresholds": "JUnitThresholds is not empty.",
        "state_thresholds": "StateThresholds is not empty.",
        "title": "Title is missing.",
        "priority": "Priority should be 0 or more."
      }
    }
    """

  @concurrent
  Scenario: given create dependencies state_settings request with invalid fields should return error
    When I am admin
    When I do POST /api/v4/state-settings:
    """json
    {
      "method": "dependencies",
      "priority": -1,
      "inherited_entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-state-settings-to-create-1"
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
    """
    {
      "errors": {
        "enabled": "Enabled is missing.",
        "entity_pattern": "EntityPattern is missing.",
        "inherited_entity_pattern": "InheritedEntityPattern is not empty.",
        "junit_thresholds": "JUnitThresholds is not empty.",
        "state_thresholds": "StateThresholds is missing.",
        "title": "Title is missing.",
        "priority": "Priority should be 0 or more."
      }
    }
    """

  @concurrent
  Scenario: given create state_settings request with invalid method should return error
    When I am admin
    When I do POST /api/v4/state-settings:
    """json
    {
      "method": "worst"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "method": "Method must be one of [inherited,dependencies]."
      }
    }
    """

  @concurrent
  Scenario: given create state_settings request with invalid patterns should return error
    When I am admin
    When I do POST /api/v4/state-settings:
    """json
    {
      "method": "inherited",
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "some",
            "cond": {
              "type": "eq",
              "value": "test-resource-state-settings-to-create-1"
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
              "value": "test-resource-state-settings-to-create-1"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "entity_pattern": "EntityPattern is invalid entity pattern.",
        "inherited_entity_pattern": "InheritedEntityPattern is invalid entity pattern."
      }
    }
    """
    When I do POST /api/v4/state-settings:
    """json
    {
      "method": "inherited",
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "component",
            "cond": {
              "type": "eq",
              "value": "test-resource-state-settings-to-create-1"
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
    """
    {
      "errors": {
        "entity_pattern": "EntityPattern is invalid entity pattern.",
        "inherited_entity_pattern": "InheritedEntityPattern is invalid entity pattern."
      }
    }
    """

  @concurrent
  Scenario: given create dependencies state_settings request with invalid state_thresholds fields should return error
    When I am admin
    When I do POST /api/v4/state-settings:
    """json
    {
      "method": "dependencies",
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-state-settings-invalid"
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
    """
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
  Scenario: given create state_settings request with already exist title should return error
    When I am admin
    When I do POST /api/v4/state-settings:
    """json
    {
      "title": "inherited-settings-to-unique-title",
      "method": "dependencies",
      "enabled": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-resource-state-settings-to-create-2"
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
        "title": "Title already exists."
      }
    }
    """

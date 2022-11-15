Feature: Create a saved pattern
  I need to be able to create a saved pattern
  Only admin should be able to create a saved pattern

  Scenario: given create alarm pattern request should return ok
    When I am noperms
    When I do POST /api/v4/patterns:
    """json
    {
      "title": "test-pattern-to-create-1-title",
      "type": "alarm",
      "is_corporate": false,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-create-1-pattern"
            }
          },
          {
            "field": "v.duration",
            "cond": {
              "type": "gt",
              "value": {
                "value": 3,
                "unit": "m"
              }
            }
          },
          {
            "field": "v.ack",
            "cond": {
              "type": "exist",
              "value": false
            }
          }
        ],
        [
          {
            "field": "v.last_event_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 1,
                "unit": "m"
              }
            }
          },
          {
            "field": "v.last_update_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 1,
                "unit": "m"
              }
            }
          },
          {
            "field": "v.resolved",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 1,
                "unit": "m"
              }
            }
          },
          {
            "field": "v.creation_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 1,
                "unit": "m"
              }
            }
          },
          {
            "field": "v.activation_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 1,
                "unit": "m"
              }
            }
          },
          {
            "field": "v.ack.t",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 1,
                "unit": "m"
              }
            }
          },
          {
            "field": "v.creation_date",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263992,
                "to": 1605264992
              }
            }
          },
          {
            "field": "v.activation_date",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263992,
                "to": 1605264992
              }
            }
          },
          {
            "field": "v.ack.t",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263992,
                "to": 1605264992
              }
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "nopermsuser",
        "name": "nopermsuser"
      },
      "title": "test-pattern-to-create-1-title",
      "type": "alarm",
      "is_corporate": false,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-create-1-pattern"
            }
          },
          {
            "field": "v.duration",
            "cond": {
              "type": "gt",
              "value": {
                "value": 3,
                "unit": "m"
              }
            }
          },
          {
            "field": "v.ack",
            "cond": {
              "type": "exist",
              "value": false
            }
          }
        ],
        [
          {
            "field": "v.last_event_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 1,
                "unit": "m"
              }
            }
          },
          {
            "field": "v.last_update_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 1,
                "unit": "m"
              }
            }
          },
          {
            "field": "v.resolved",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 1,
                "unit": "m"
              }
            }
          },
          {
            "field": "v.creation_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 1,
                "unit": "m"
              }
            }
          },
          {
            "field": "v.activation_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 1,
                "unit": "m"
              }
            }
          },
          {
            "field": "v.ack.t",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 1,
                "unit": "m"
              }
            }
          },
          {
            "field": "v.creation_date",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263992,
                "to": 1605264992
              }
            }
          },
          {
            "field": "v.activation_date",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263992,
                "to": 1605264992
              }
            }
          },
          {
            "field": "v.ack.t",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263992,
                "to": 1605264992
              }
            }
          }
        ]
      ]
    }
    """
    When I do GET /api/v4/patterns/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "nopermsuser",
        "name": "nopermsuser"
      },
      "title": "test-pattern-to-create-1-title",
      "type": "alarm",
      "is_corporate": false,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-create-1-pattern"
            }
          },
          {
            "field": "v.duration",
            "cond": {
              "type": "gt",
              "value": {
                "value": 3,
                "unit": "m"
              }
            }
          },
          {
            "field": "v.ack",
            "cond": {
              "type": "exist",
              "value": false
            }
          }
        ],
        [
          {
            "field": "v.last_event_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 1,
                "unit": "m"
              }
            }
          },
          {
            "field": "v.last_update_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 1,
                "unit": "m"
              }
            }
          },
          {
            "field": "v.resolved",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 1,
                "unit": "m"
              }
            }
          },
          {
            "field": "v.creation_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 1,
                "unit": "m"
              }
            }
          },
          {
            "field": "v.activation_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 1,
                "unit": "m"
              }
            }
          },
          {
            "field": "v.ack.t",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 1,
                "unit": "m"
              }
            }
          },
          {
            "field": "v.creation_date",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263992,
                "to": 1605264992
              }
            }
          },
          {
            "field": "v.activation_date",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263992,
                "to": 1605264992
              }
            }
          },
          {
            "field": "v.ack.t",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263992,
                "to": 1605264992
              }
            }
          }
        ]
      ]
    }
    """

  Scenario: given create entity pattern request should return ok
    When I am noperms
    When I do POST /api/v4/patterns:
    """json
    {
      "title": "test-pattern-to-create-2-title",
      "type": "entity",
      "is_corporate": false,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-create-2-pattern"
            }
          },
          {
            "field": "last_event_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 1,
                "unit": "m"
              }
            }
          }
        ],
        [
          {
            "field": "connector",
            "cond": {
              "type": "is_one_of",
              "value": ["test-pattern-to-create-2-pattern"]
            }
          },
          {
            "field": "component",
            "cond": {
              "type": "is_one_of",
              "value": ["test-pattern-to-create-2-pattern"]
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "nopermsuser",
        "name": "nopermsuser"
      },
      "title": "test-pattern-to-create-2-title",
      "type": "entity",
      "is_corporate": false,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-create-2-pattern"
            }
          },
          {
            "field": "last_event_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 1,
                "unit": "m"
              }
            }
          }
        ],
        [
          {
            "field": "connector",
            "cond": {
              "type": "is_one_of",
              "value": ["test-pattern-to-create-2-pattern"]
            }
          },
          {
            "field": "component",
            "cond": {
              "type": "is_one_of",
              "value": ["test-pattern-to-create-2-pattern"]
            }
          }
        ]
      ]
    }
    """
    When I do GET /api/v4/patterns/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "nopermsuser",
        "name": "nopermsuser"
      },
      "title": "test-pattern-to-create-2-title",
      "type": "entity",
      "is_corporate": false,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-create-2-pattern"
            }
          },
          {
            "field": "last_event_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 1,
                "unit": "m"
              }
            }
          }
        ],
        [
          {
            "field": "connector",
            "cond": {
              "type": "is_one_of",
              "value": ["test-pattern-to-create-2-pattern"]
            }
          },
          {
            "field": "component",
            "cond": {
              "type": "is_one_of",
              "value": ["test-pattern-to-create-2-pattern"]
            }
          }
        ]
      ]
    }
    """

  Scenario: given create pbehavior pattern request should return ok
    When I am noperms
    When I do POST /api/v4/patterns:
    """json
    {
      "title": "test-pattern-to-create-3-title",
      "type": "pbehavior",
      "is_corporate": false,
      "pbehavior_pattern": [
        [
          {
            "field": "pbehavior_info.type",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-create-3-pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "nopermsuser",
        "name": "nopermsuser"
      },
      "title": "test-pattern-to-create-3-title",
      "type": "pbehavior",
      "is_corporate": false,
      "pbehavior_pattern": [
        [
          {
            "field": "pbehavior_info.type",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-create-3-pattern"
            }
          }
        ]
      ]
    }
    """
    When I do GET /api/v4/patterns/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "nopermsuser",
        "name": "nopermsuser"
      },
      "title": "test-pattern-to-create-3-title",
      "type": "pbehavior",
      "is_corporate": false,
      "pbehavior_pattern": [
        [
          {
            "field": "pbehavior_info.type",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-create-3-pattern"
            }
          }
        ]
      ]
    }
    """

  Scenario: given create corporate pattern request should return ok
    When I am admin
    When I do POST /api/v4/patterns:
    """json
    {
      "title": "test-pattern-to-create-4-title",
      "type": "alarm",
      "is_corporate": true,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-create-4-pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "title": "test-pattern-to-create-4-title",
      "type": "alarm",
      "is_corporate": true,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-create-4-pattern"
            }
          }
        ]
      ]
    }
    """
    When I do GET /api/v4/patterns/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "title": "test-pattern-to-create-4-title",
      "type": "alarm",
      "is_corporate": true,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-create-4-pattern"
            }
          }
        ]
      ]
    }
    """

  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/patterns
    Then the response code should be 401

  Scenario: given create corporate pattern request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/patterns:
    """json
    {
      "title": "test-pattern-to-create-5-title",
      "type": "alarm",
      "is_corporate": true,
      "alarm_pattern": [
        [
          {
            "field": "v.component",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-create-5-pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 403

  Scenario: given create request with missing fields should return bad request error
    When I am admin
    When I do POST /api/v4/patterns:
    """json
    {
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "title": "Title is missing.",
        "type": "Type is missing.",
        "is_corporate": "IsCorporate is missing."
      }
    }
    """
    When I do POST /api/v4/patterns:
    """json
    {
      "type": "unknown"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "type": "Type must be one of [alarm entity pbehavior]."
      }
    }
    """
    When I do POST /api/v4/patterns:
    """json
    {
      "type": "alarm"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "alarm_pattern": "AlarmPattern is missing.",
        "title": "Title is missing.",
        "is_corporate": "IsCorporate is missing."
      }
    }
    """
    When I do POST /api/v4/patterns:
    """json
    {
      "type": "alarm",
      "alarm_pattern": []
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "alarm_pattern": "AlarmPattern is missing.",
        "title": "Title is missing.",
        "is_corporate": "IsCorporate is missing."
      }
    }
    """
    When I do POST /api/v4/patterns:
    """json
    {
      "type": "entity"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "entity_pattern": "EntityPattern is missing.",
        "title": "Title is missing.",
        "is_corporate": "IsCorporate is missing."
      }
    }
    """
    When I do POST /api/v4/patterns:
    """json
    {
      "type": "entity",
      "entity_pattern": []
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "entity_pattern": "EntityPattern is missing.",
        "title": "Title is missing.",
        "is_corporate": "IsCorporate is missing."
      }
    }
    """
    When I do POST /api/v4/patterns:
    """json
    {
      "type": "pbehavior"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "pbehavior_pattern": "PbehaviorPattern is missing.",
        "title": "Title is missing.",
        "is_corporate": "IsCorporate is missing."
      }
    }
    """
    When I do POST /api/v4/patterns:
    """json
    {
      "type": "pbehavior",
      "pbehavior_pattern": []
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "pbehavior_pattern": "PbehaviorPattern is missing.",
        "title": "Title is missing.",
        "is_corporate": "IsCorporate is missing."
      }
    }
    """

  Scenario: given create request with invalid patterns format should return bad request error
    When I am admin
    When I do POST /api/v4/patterns:
    """json
    {
      "alarm_pattern": [
        []
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "alarm_pattern": "AlarmPattern is invalid alarm pattern."
      }
    }
    """
    When I do POST /api/v4/patterns:
    """json
    {
      "alarm_pattern": [
        [
          {}
        ]
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "alarm_pattern": "AlarmPattern is invalid alarm pattern."
      }
    }
    """
    When I do POST /api/v4/patterns:
    """json
    {
      "alarm_pattern": [
        [
          {
            "field": "unknown",
            "cond": {
              "type": "eq",
              "value": "test"
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
        "alarm_pattern": "AlarmPattern is invalid alarm pattern."
      }
    }
    """
    When I do POST /api/v4/patterns:
    """json
    {
      "alarm_pattern": [
        [
          {
            "field": "v.connector",
            "cond": {
              "type": "eq",
              "value": 2
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
        "alarm_pattern": "AlarmPattern is invalid alarm pattern."
      }
    }
    """
    When I do POST /api/v4/patterns:
    """json
    {
      "alarm_pattern": [
        [
          {
            "field": "v.connector",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-create-6-pattern"
            }
          },
          {
            "field": "v.connector",
            "cond": {
              "type": "eq",
              "value": 2
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
        "alarm_pattern": "AlarmPattern is invalid alarm pattern."
      }
    }
    """
    When I do POST /api/v4/patterns:
    """json
    {
      "alarm_pattern": [
        [
          {
            "field": "v.connector",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-create-6-pattern"
            }
          }
        ],
        [
          {
            "field": "v.connector",
            "cond": {
              "type": "eq",
              "value": 2
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
        "alarm_pattern": "AlarmPattern is invalid alarm pattern."
      }
    }
    """

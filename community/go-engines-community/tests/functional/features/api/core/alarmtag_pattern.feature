Feature: Update alarm tag
  I need to be able to update a alarm tag

  @concurrent
  Scenario: given alarm tag and corporate pattern update and delete actions should update patterns in alarm tag
    When I am admin
    When I do POST /api/v4/alarm-tags:
    """json
    {
      "value": "test-alarm-tag-to-edit-patterns",
      "color": "#AABBCC",
      "corporate_entity_pattern": "test-pattern-to-alarm-tag-pattern-1",
      "corporate_alarm_pattern": "test-pattern-to-alarm-tag-pattern-2"
    }
    """
    Then the response code should be 201
    When I save response tagId={{ .lastResponse._id }}
    When I do GET /api/v4/alarm-tags/{{ .tagId }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "corporate_entity_pattern": "test-pattern-to-alarm-tag-pattern-1",
      "corporate_entity_pattern_title": "test-pattern-to-alarm-tag-pattern-1-title",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-alarm-tag-pattern-1-pattern"
            }
          }
        ]
      ],
      "corporate_alarm_pattern": "test-pattern-to-alarm-tag-pattern-2",
      "corporate_alarm_pattern_title": "test-pattern-to-alarm-tag-pattern-2-title",
      "alarm_pattern": [
        [
          {
            "field": "v.state.val",
            "cond": {
              "type": "eq",
              "value": 3
            }
          }
        ],
        [
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
    When I do PUT /api/v4/patterns/test-pattern-to-alarm-tag-pattern-1:
    """json
    {
      "title": "test-pattern-to-alarm-tag-pattern-1-title-updated",
      "type": "entity",
      "is_corporate": true,
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-alarm-tag-pattern-1-pattern-updated"
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
        ]
      ]
    }
    """
    Then the response code should be 200
    When I do PUT /api/v4/patterns/test-pattern-to-alarm-tag-pattern-2:
    """json
    {
      "title": "test-pattern-to-alarm-tag-pattern-2-title-updated",
      "type": "alarm",
      "is_corporate": true,
      "alarm_pattern": [
        [
          {
            "field": "v.state.val",
            "cond": {
              "type": "eq",
              "value": 2
            }
          },
          {
            "field": "v.last_event_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 2,
                "unit": "m"
              }
            }
          },
          {
            "field": "v.last_update_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 2,
                "unit": "m"
              }
            }
          },
          {
            "field": "v.resolved",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 2,
                "unit": "m"
              }
            }
          }
        ],
        [
          {
            "field": "v.creation_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 2,
                "unit": "m"
              }
            }
          },
          {
            "field": "v.activation_date",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 2,
                "unit": "m"
              }
            }
          },
          {
            "field": "v.ack.t",
            "cond": {
              "type": "relative_time",
              "value": {
                "value": 2,
                "unit": "m"
              }
            }
          }
        ],
        [
          {
            "field": "v.creation_date",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263993,
                "to": 1605264993
              }
            }
          },
          {
            "field": "v.activation_date",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263993,
                "to": 1605264993
              }
            }
          },
          {
            "field": "v.ack.t",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263993,
                "to": 1605264993
              }
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 200
    When I do GET /api/v4/alarm-tags/{{ .tagId }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "corporate_entity_pattern": "test-pattern-to-alarm-tag-pattern-1",
      "corporate_entity_pattern_title": "test-pattern-to-alarm-tag-pattern-1-title-updated",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-alarm-tag-pattern-1-pattern-updated"
            }
          }
        ]
      ],
      "corporate_alarm_pattern": "test-pattern-to-alarm-tag-pattern-2",
      "corporate_alarm_pattern_title": "test-pattern-to-alarm-tag-pattern-2-title-updated",
      "alarm_pattern": [
        [
          {
            "field": "v.state.val",
            "cond": {
              "type": "eq",
              "value": 2
            }
          }
        ],
        [
          {
            "field": "v.creation_date",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263993,
                "to": 1605264993
              }
            }
          },
          {
            "field": "v.activation_date",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263993,
                "to": 1605264993
              }
            }
          },
          {
            "field": "v.ack.t",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263993,
                "to": 1605264993
              }
            }
          }
        ]
      ]
    }
    """
    When I do DELETE /api/v4/patterns/test-pattern-to-alarm-tag-pattern-1
    Then the response code should be 204
    When I do DELETE /api/v4/patterns/test-pattern-to-alarm-tag-pattern-2
    Then the response code should be 204
    When I do GET /api/v4/alarm-tags/{{ .tagId }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pattern-to-alarm-tag-pattern-1-pattern-updated"
            }
          }
        ]
      ],
      "alarm_pattern": [
        [
          {
            "field": "v.state.val",
            "cond": {
              "type": "eq",
              "value": 2
            }
          }
        ],
        [
          {
            "field": "v.creation_date",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263993,
                "to": 1605264993
              }
            }
          },
          {
            "field": "v.activation_date",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263993,
                "to": 1605264993
              }
            }
          },
          {
            "field": "v.ack.t",
            "cond": {
              "type": "absolute_time",
              "value": {
                "from": 1605263993,
                "to": 1605264993
              }
            }
          }
        ]
      ]
    }
    """
    Then the response key "corporate_entity_pattern" should not exist
    Then the response key "corporate_entity_pattern_title" should not exist
    Then the response key "corporate_alarm_pattern" should not exist
    Then the response key "corporate_alarm_pattern_title" should not exist

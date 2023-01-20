Feature: Bulk update idlerules
  I need to be able to bulk update idlerules
  Only admin should be able to bulk update idlerules

  Scenario: given bulk update request and no auth should not allow access
    When I do PUT /api/v4/bulk/idle-rules
    Then the response code should be 401

  Scenario: given bulk update request and auth by api key without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/bulk/idle-rules
    Then the response code should be 403

  Scenario: given update request should return multistatus and should be handled independently
    When I am admin
    When I do PUT /api/v4/bulk/idle-rules:
    """json
    [
      {
        "_id": "test-idle-rule-to-bulk-update-1",
        "name": "test-idle-rule-to-bulk-update-1-name",
        "description": "test-idle-rule-to-bulk-update-1-description",
        "type": "alarm",
        "alarm_condition": "last_event",
        "enabled": true,
        "priority": 30,
        "duration": {
          "value": 5,
          "unit": "s"
        },
        "alarm_pattern": [
          [
            {
              "field": "v.component",
              "cond": {
                "type": "eq",
                "value": "test-idle-rule-to-bulk-update-1-alarm-updated"
              }
            }
          ]
        ],
        "entity_pattern": [
          [
            {
              "field": "name",
              "cond": {
                "type": "eq",
                "value": "test-idle-rule-to-bulk-update-1-resource-updated"
              }
            }
          ]
        ],
        "operation": {
          "type": "snooze",
          "parameters": {
            "output": "test-idle-rule-to-bulk-update-1-operation-output-updated",
            "duration": {
              "value": 5,
              "unit": "s"
            }
          }
        },
        "disable_during_periods": ["maintenance"]
      },
      {
        "_id": "test-idle-rule-to-bulk-update-1",
        "name": "test-idle-rule-to-bulk-update-1-name-twice",
        "description": "test-idle-rule-to-bulk-update-1-description",
        "type": "alarm",
        "alarm_condition": "last_event",
        "enabled": true,
        "priority": 30,
        "duration": {
          "value": 5,
          "unit": "s"
        },
        "alarm_pattern": [
          [
            {
              "field": "v.component",
              "cond": {
                "type": "eq",
                "value": "test-idle-rule-to-bulk-update-1-alarm-updated"
              }
            }
          ]
        ],
        "entity_pattern": [
          [
            {
              "field": "name",
              "cond": {
                "type": "eq",
                "value": "test-idle-rule-to-bulk-update-1-resource-updated"
              }
            }
          ]
        ],
        "operation": {
          "type": "snooze",
          "parameters": {
            "output": "test-idle-rule-to-bulk-update-1-operation-output-updated",
            "duration": {
              "value": 5,
              "unit": "s"
            }
          }
        },
        "disable_during_periods": ["maintenance"]
      },
      {},
      {
        "type": "notexists"
      },
      {
        "type": "alarm",
        "operation": {
          "type": "notexists"
        }
      },
      [],
      {
        "_id": "test-idle-rule-to-bulk-update-2",
        "name": "test-idle-rule-to-bulk-update-2-name",
        "description": "test-idle-rule-to-bulk-update-2-description",
        "type": "alarm",
        "alarm_condition": "last_event",
        "enabled": true,
        "priority": 31,
        "duration": {
          "value": 5,
          "unit": "s"
        },
        "alarm_pattern": [
          [
            {
              "field": "v.component",
              "cond": {
                "type": "eq",
                "value": "test-idle-rule-to-bulk-update-2-alarm-updated"
              }
            }
          ]
        ],
        "entity_pattern": [
          [
            {
              "field": "name",
              "cond": {
                "type": "eq",
                "value": "test-idle-rule-to-bulk-update-2-resource-updated"
              }
            }
          ]
        ],
        "operation": {
          "type": "snooze",
          "parameters": {
            "output": "test-idle-rule-to-bulk-update-2-operation-output-updated",
            "duration": {
              "value": 5,
              "unit": "s"
            }
          }
        },
        "disable_during_periods": ["maintenance"]
      },
      {
        "_id": "test-idle-rule-to-bulk-update-3",
        "name": "test-idle-rule-to-bulk-update-3-name",
        "description": "test-idle-rule-to-bulk-update-3-description",
        "type": "alarm",
        "alarm_condition": "last_event",
        "enabled": true,
        "priority": 20,
        "duration": {
          "value": 3,
          "unit": "s"
        },
        "corporate_alarm_pattern": "test-pattern-to-rule-edit-1",
        "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
        "operation": {
          "type": "snooze",
          "parameters": {
            "output": "test-idle-rule-to-bulk-update-3-operation-output",
            "duration": {
              "value": 3,
              "unit": "s"
            }
          }
        },
        "disable_during_periods": ["pause"]
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "id": "test-idle-rule-to-bulk-update-1",
        "status": 200,
        "item": {
          "_id": "test-idle-rule-to-bulk-update-1",
          "name": "test-idle-rule-to-bulk-update-1-name",
          "description": "test-idle-rule-to-bulk-update-1-description",
          "type": "alarm",
          "alarm_condition": "last_event",
          "enabled": true,
          "priority": 30,
          "duration": {
            "value": 5,
            "unit": "s"
          },
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-idle-rule-to-bulk-update-1-alarm-updated"
                }
              }
            ]
          ],
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-idle-rule-to-bulk-update-1-resource-updated"
                }
              }
            ]
          ],
          "operation": {
            "type": "snooze",
            "parameters": {
              "output": "test-idle-rule-to-bulk-update-1-operation-output-updated",
              "duration": {
                "value": 5,
                "unit": "s"
              }
            }
          },
          "disable_during_periods": ["maintenance"]
        }
      },
      {
        "id": "test-idle-rule-to-bulk-update-1",
        "status": 200,
        "item": {
          "_id": "test-idle-rule-to-bulk-update-1",
          "name": "test-idle-rule-to-bulk-update-1-name-twice",
          "description": "test-idle-rule-to-bulk-update-1-description",
          "type": "alarm",
          "alarm_condition": "last_event",
          "enabled": true,
          "priority": 30,
          "duration": {
            "value": 5,
            "unit": "s"
          },
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-idle-rule-to-bulk-update-1-alarm-updated"
                }
              }
            ]
          ],
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-idle-rule-to-bulk-update-1-resource-updated"
                }
              }
            ]
          ],
          "operation": {
            "type": "snooze",
            "parameters": {
              "output": "test-idle-rule-to-bulk-update-1-operation-output-updated",
              "duration": {
                "value": 5,
                "unit": "s"
              }
            }
          },
          "disable_during_periods": ["maintenance"]
        }
      },
      {
        "status": 400,
        "item": {},
        "errors": {
          "_id": "ID is missing.",
          "duration.value": "Value is missing.",
          "duration.unit": "Unit is missing.",
          "enabled": "Enabled is missing.",
          "name": "Name is missing.",
          "priority": "Priority is missing.",
          "type": "Type is missing."
        }
      },
      {
        "status": 400,
        "item": {
          "type": "notexists"
        },
        "errors": {
          "_id": "ID is missing.",
          "type": "Type must be one of [alarm entity]."
        }
      },
      {
        "status": 400,
        "item": {
          "type": "alarm",
          "operation": {
            "type": "notexists"
          }
        },
        "errors": {
          "_id": "ID is missing.",
          "alarm_condition": "AlarmCondition is missing.",
          "alarm_pattern": "AlarmPattern or EntityPattern is required.",
          "entity_pattern": "EntityPattern or AlarmPattern is required.",
          "operation.type": "Type must be one of [ack ackremove cancel assocticket changestate snooze pbehavior]."
        }
      },
      {
        "status": 400,
        "item": [],
        "error": "value doesn't contain object; it contains array"
      },
      {
        "id": "test-idle-rule-to-bulk-update-2",
        "status": 200,
        "item": {
          "_id": "test-idle-rule-to-bulk-update-2",
          "name": "test-idle-rule-to-bulk-update-2-name",
          "description": "test-idle-rule-to-bulk-update-2-description",
          "type": "alarm",
          "alarm_condition": "last_event",
          "enabled": true,
          "priority": 31,
          "duration": {
            "value": 5,
            "unit": "s"
          },
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-idle-rule-to-bulk-update-2-alarm-updated"
                }
              }
            ]
          ],
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-idle-rule-to-bulk-update-2-resource-updated"
                }
              }
            ]
          ],
          "operation": {
            "type": "snooze",
            "parameters": {
              "output": "test-idle-rule-to-bulk-update-2-operation-output-updated",
              "duration": {
                "value": 5,
                "unit": "s"
              }
            }
          },
          "disable_during_periods": ["maintenance"]
        }
      },
      {
        "id": "test-idle-rule-to-bulk-update-3",
        "status": 200,
        "item": {
          "name": "test-idle-rule-to-bulk-update-3-name",
          "description": "test-idle-rule-to-bulk-update-3-description",
          "type": "alarm",
          "alarm_condition": "last_event",
          "enabled": true,
          "priority": 20,
          "duration": {
            "value": 3,
            "unit": "s"
          },
          "corporate_alarm_pattern": "test-pattern-to-rule-edit-1",
          "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
          "operation": {
            "type": "snooze",
            "parameters": {
              "output": "test-idle-rule-to-bulk-update-3-operation-output",
              "duration": {
                "value": 3,
                "unit": "s"
              }
            }
          },
          "disable_during_periods": ["pause"]
        }
      }
    ]
    """
    When I do GET /api/v4/idle-rules?search=test-idle-rule-to-bulk-update
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-idle-rule-to-bulk-update-1",
          "name": "test-idle-rule-to-bulk-update-1-name-twice",
          "description": "test-idle-rule-to-bulk-update-1-description",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "type": "alarm",
          "alarm_condition": "last_event",
          "enabled": true,
          "duration": {
            "value": 5,
            "unit": "s"
          },
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-idle-rule-to-bulk-update-1-alarm-updated"
                }
              }
            ]
          ],
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-idle-rule-to-bulk-update-1-resource-updated"
                }
              }
            ]
          ],
          "operation": {
            "type": "snooze",
            "parameters": {
              "output": "test-idle-rule-to-bulk-update-1-operation-output-updated",
              "duration": {
                "value": 5,
                "unit": "s"
              }
            }
          },
          "disable_during_periods": ["maintenance"]
        },
        {
          "_id": "test-idle-rule-to-bulk-update-2",
          "name": "test-idle-rule-to-bulk-update-2-name",
          "description": "test-idle-rule-to-bulk-update-2-description",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "type": "alarm",
          "alarm_condition": "last_event",
          "enabled": true,
          "duration": {
            "value": 5,
            "unit": "s"
          },
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-idle-rule-to-bulk-update-2-alarm-updated"
                }
              }
            ]
          ],
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-idle-rule-to-bulk-update-2-resource-updated"
                }
              }
            ]
          ],
          "operation": {
            "type": "snooze",
            "parameters": {
              "output": "test-idle-rule-to-bulk-update-2-operation-output-updated",
              "duration": {
                "value": 5,
                "unit": "s"
              }
            }
          },
          "disable_during_periods": ["maintenance"]
        },
        {
          "_id": "test-idle-rule-to-bulk-update-3",
          "name": "test-idle-rule-to-bulk-update-3-name",
          "description": "test-idle-rule-to-bulk-update-3-description",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "type": "alarm",
          "alarm_condition": "last_event",
          "enabled": true,
          "duration": {
            "value": 3,
            "unit": "s"
          },
          "corporate_alarm_pattern": "test-pattern-to-rule-edit-1",
          "corporate_alarm_pattern_title": "test-pattern-to-rule-edit-1-title",
          "alarm_pattern": [
            [
              {
                "field": "v.component",
                "cond": {
                  "type": "eq",
                  "value": "test-pattern-to-rule-edit-1-pattern"
                }
              }
            ]
          ],
          "corporate_entity_pattern": "test-pattern-to-rule-edit-2",
          "corporate_entity_pattern_title": "test-pattern-to-rule-edit-2-title",
          "entity_pattern": [
            [
              {
                "field": "name",
                "cond": {
                  "type": "eq",
                  "value": "test-pattern-to-rule-edit-2-pattern"
                }
              }
            ]
          ],
          "operation": {
            "type": "snooze",
            "parameters": {
              "output": "test-idle-rule-to-bulk-update-3-operation-output",
              "duration": {
                "value": 3,
                "unit": "s"
              }
            }
          }
        }        
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 3
      }
    }
    """

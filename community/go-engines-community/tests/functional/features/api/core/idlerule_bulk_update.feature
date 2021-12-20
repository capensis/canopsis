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
        "alarm_patterns": [
          {
            "_id": "test-idle-rule-to-bulk-update-1-alarm-updated"
          }
        ],
        "entity_patterns": [
          {
            "name": "test-idle-rule-to-bulk-update-1-resource-updated"
          }
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
        "alarm_patterns": [],
        "entity_patterns": [],
        "operation": {
          "type": "notexists"
        }
      },
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
        "alarm_patterns": [
          {
            "_id": "test-idle-rule-to-bulk-update-2-alarm-updated"
          }
        ],
        "entity_patterns": [
          {
            "name": "test-idle-rule-to-bulk-update-2-resource-updated"
          }
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
          "alarm_patterns": [
            {
              "_id": "test-idle-rule-to-bulk-update-1-alarm-updated"
            }
          ],
          "entity_patterns": [
            {
              "name": "test-idle-rule-to-bulk-update-1-resource-updated"
            }
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
          "alarm_patterns": [],
          "entity_patterns": [],
          "operation": {
            "type": "notexists"
          }
        },
        "errors": {
          "_id": "ID is missing.",
          "alarm_condition": "AlarmCondition is missing.",
          "alarm_patterns": "AlarmPatterns or EntityPatterns is required.",
          "entity_patterns": "EntityPatterns or AlarmPatterns is required.",
          "operation.type": "Type must be one of [ack ackremove cancel assocticket changestate snooze pbehavior]."
        }
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
          "alarm_patterns": [
            {
              "_id": "test-idle-rule-to-bulk-update-2-alarm-updated"
            }
          ],
          "entity_patterns": [
            {
              "name": "test-idle-rule-to-bulk-update-2-resource-updated"
            }
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
          "name": "test-idle-rule-to-bulk-update-1-name",
          "description": "test-idle-rule-to-bulk-update-1-description",
          "author": "root",
          "type": "alarm",
          "alarm_condition": "last_event",
          "enabled": true,
          "priority": 30,
          "created": 1616567033,
          "duration": {
            "value": 5,
            "unit": "s"
          },
          "alarm_patterns": [
            {
              "_id": "test-idle-rule-to-bulk-update-1-alarm-updated"
            }
          ],
          "entity_patterns": [
            {
              "name": "test-idle-rule-to-bulk-update-1-resource-updated"
            }
          ],
          "operation": {
            "type": "snooze",
            "parameters": {
              "author": "root",
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
          "author": "root",
          "type": "alarm",
          "alarm_condition": "last_event",
          "enabled": true,
          "priority": 31,
          "created": 1616567033,
          "duration": {
            "value": 5,
            "unit": "s"
          },
          "alarm_patterns": [
            {
              "_id": "test-idle-rule-to-bulk-update-2-alarm-updated"
            }
          ],
          "entity_patterns": [
            {
              "name": "test-idle-rule-to-bulk-update-2-resource-updated"
            }
          ],
          "operation": {
            "type": "snooze",
            "parameters": {
              "author": "root",
              "output": "test-idle-rule-to-bulk-update-2-operation-output-updated",
              "duration": {
                "value": 5,
                "unit": "s"
              }
            }
          },
          "disable_during_periods": ["maintenance"]
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 2
      }
    }
    """

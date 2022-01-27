Feature: Create an metaalarmrule
  I need to be able to create a metaalarmrule
  Only admin should be able to create a metaalarmrule

  Scenario: given create request should return ok
    When I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "_id": "complex-1",
      "auto_resolve": true,
      "name": "complex-test-1",
      "type": "complex",
      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
      "config": {
        "time_interval": {
          "value": 1,
          "unit": "m"
        },
        "alarm_patterns": [
          {
            "v": {"resource": "123"}
          }
        ],
        "threshold_rate": 1
      }
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """
    {
      "_id": "complex-1",
      "auto_resolve": true,
      "name": "complex-test-1",
      "author": "root",
      "type": "complex",
      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
      "config": {
        "time_interval": {
          "value": 1,
          "unit": "m"
        },
        "alarm_patterns": [
          {
            "v": {"resource": "123"}
          }
        ],
        "threshold_rate": 1
      }
    }
    """
    When I do GET /api/v4/cat/metaalarmrules/{{ .lastResponse._id }}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "_id": "complex-1",
      "auto_resolve": true,
      "name": "complex-test-1",
      "author": "root",
      "type": "complex",
      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
      "config": {
        "time_interval": {
          "value": 1,
          "unit": "m"
        },
        "alarm_patterns": [
          {
            "v": {"resource": "123"}
          }
        ],
        "threshold_rate": 1
      }
    }
    """

  Scenario: given create request with wrong type should return bad request
    When I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "_id": "wrong-type-1",
      "auto_resolve": false,
      "name": "wrong-type-1",
      "config": {
      },
      "patterns": null,
      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
      "type": "attribute_path"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "type": "Type must be one of [relation timebased attribute complex valuegroup corel]."
      }
    }
    """

  Scenario: given create request with wrong config type should return bad request
    When I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "_id": "complex-wrong-config-1",
      "auto_resolve": true,
      "name": "complex-test-1",
      "type": "complex",
      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
      "config": {
        "time_interval": {
          "value": 1,
          "unit": "m"
        },
        "alarm_patterns": [
          {
            "v": {"resource": "123"}
          }
        ],
        "threshold_rate": 1,
        "value_paths": ["resource.path"]
      }
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "config": "value_paths config can not be in type complex."
      }
    }
    """

  Scenario: given create request and no auth user should not allow access
    When I do POST /api/v4/cat/metaalarmrules
    Then the response code should be 401

  Scenario: given create request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/cat/metaalarmrules
    Then the response code should be 403

  Scenario: given create request with already exists name should return error
    When I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "_id": "test-metaalarm-to-get-1",
      "auto_resolve": false,
      "name": "test-attribute-type-1",
      "config": {
      },
      "patterns": null,
      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
      "type": "attribute"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "_id": "ID already exists."
      }
    }
    """

  Scenario: given create request with wrong alarm patterns should return bad request
    When I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "_id": "wrong-alarm-pattern-1",
      "auto_resolve": false,
      "name": "wrong-alarm-pattern-1",
      "config": {
        "alarm_patterns": [
         {
           "resource": "123"
         }
       ]
      },
      "patterns": null,
      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
      "type": "complex"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "config.alarm_patterns": "Invalid alarm patterns."
      }
    }
    """

  Scenario: given create request with wrong entity patterns should return bad request
    When I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "auto_resolve": true,
      "name": "complex-test-1",
      "type": "complex",
      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
      "patterns": {
        "entity_patterns": [
          {
            "resource": "123"
          }
        ]
      },
      "config": {
        "time_interval": {
          "value": 1,
          "unit": "m"
        },
        "threshold_rate": 1,
        "entity_patterns": [
          {
            "address": "current"
          }
        ]
      }
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "config.entity_patterns": "Invalid entity patterns.",
        "patterns": "Invalid entity patterns."
      }
    }
    """

  Scenario: given create request with wrong event patterns should return bad request
    When I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "auto_resolve": true,
      "name": "complex-test-1",
      "type": "complex",
      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
      "config": {
        "time_interval": {
          "value": 1,
          "unit": "m"
        },
        "threshold_rate": 1,
        "event_patterns": [
          {
            "address": {
              "regex_match1": "matched"
            }
          }
        ]
      }
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "config.event_patterns":"Invalid event patterns."
      }
    }
    """

  Scenario: given create request with wrong config patterns should return bad request
    When I am admin
    When I do POST /api/v4/cat/metaalarmrules:
    """
    {
      "auto_resolve": true,
      "name": "complex-test-1",
      "type": "complex",
      "output_template": "{{ `{{ .Children.Alarm.Value.State.Message }}` }}",
      "config": {
        "n": [
          {
            "time_interval": {
              "value": 1,
              "unit": "m"
            },
            "threshold_rate": 1,
            "event_patterns": [
              {
                "address": {
                  "regex_match1": "matched"
                }
              }
            ]
          }
        ]
      }
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "error": "request has invalid structure"
    }
    """

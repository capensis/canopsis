Feature: crop alarm steps
  I need to be able to crop alarm steps on event

  @concurrent
  Scenario: given many check events should crop alarm steps
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-axe-crop-steps",
      "connector_name" : "test-connector-name-axe-crop-steps",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-axe-crop-steps",
      "resource" : "test-resource-axe-crop-steps",
      "state" : 1,
      "output" : "test-output-axe-crop-steps-1"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-axe-crop-steps",
      "connector_name" : "test-connector-name-axe-crop-steps",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-axe-crop-steps",
      "resource" : "test-resource-axe-crop-steps",
      "state" : 2,
      "output" : "test-output-axe-crop-steps-2"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-axe-crop-steps",
      "connector_name" : "test-connector-name-axe-crop-steps",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-axe-crop-steps",
      "resource" : "test-resource-axe-crop-steps",
      "state" : 3,
      "output" : "test-output-axe-crop-steps-3"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-axe-crop-steps",
      "connector_name" : "test-connector-name-axe-crop-steps",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-axe-crop-steps",
      "resource" : "test-resource-axe-crop-steps",
      "state" : 2,
      "output" : "test-output-axe-crop-steps-4"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-axe-crop-steps",
      "connector_name" : "test-connector-name-axe-crop-steps",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-axe-crop-steps",
      "resource" : "test-resource-axe-crop-steps",
      "state" : 1,
      "output" : "test-output-axe-crop-steps-5"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-axe-crop-steps",
      "connector_name" : "test-connector-name-axe-crop-steps",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-axe-crop-steps",
      "resource" : "test-resource-axe-crop-steps",
      "state" : 2,
      "output" : "test-output-axe-crop-steps-6"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-axe-crop-steps",
      "connector_name" : "test-connector-name-axe-crop-steps",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-axe-crop-steps",
      "resource" : "test-resource-axe-crop-steps",
      "state" : 3,
      "output" : "test-output-axe-crop-steps-7"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-axe-crop-steps",
      "connector_name" : "test-connector-name-axe-crop-steps",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-axe-crop-steps",
      "resource" : "test-resource-axe-crop-steps",
      "state" : 2,
      "output" : "test-output-axe-crop-steps-8"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-axe-crop-steps",
      "connector_name" : "test-connector-name-axe-crop-steps",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-axe-crop-steps",
      "resource" : "test-resource-axe-crop-steps",
      "state" : 1,
      "output" : "test-output-axe-crop-steps-9"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-axe-crop-steps",
      "connector_name" : "test-connector-name-axe-crop-steps",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-axe-crop-steps",
      "resource" : "test-resource-axe-crop-steps",
      "state" : 2,
      "output" : "test-output-axe-crop-steps-10"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-axe-crop-steps",
      "connector_name" : "test-connector-name-axe-crop-steps",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-axe-crop-steps",
      "resource" : "test-resource-axe-crop-steps",
      "state" : 3,
      "output" : "test-output-axe-crop-steps-11"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-axe-crop-steps",
      "connector_name" : "test-connector-name-axe-crop-steps",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-axe-crop-steps",
      "resource" : "test-resource-axe-crop-steps",
      "state" : 2,
      "output" : "test-output-axe-crop-steps-12"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-axe-crop-steps",
      "connector_name" : "test-connector-name-axe-crop-steps",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-axe-crop-steps",
      "resource" : "test-resource-axe-crop-steps",
      "state" : 1,
      "output" : "test-output-axe-crop-steps-13"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-axe-crop-steps",
      "connector_name" : "test-connector-name-axe-crop-steps",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-axe-crop-steps",
      "resource" : "test-resource-axe-crop-steps",
      "state" : 2,
      "output" : "test-output-axe-crop-steps-14"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-axe-crop-steps",
      "connector_name" : "test-connector-name-axe-crop-steps",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-axe-crop-steps",
      "resource" : "test-resource-axe-crop-steps",
      "state" : 3,
      "output" : "test-output-axe-crop-steps-15"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-axe-crop-steps",
      "connector_name" : "test-connector-name-axe-crop-steps",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-axe-crop-steps",
      "resource" : "test-resource-axe-crop-steps",
      "state" : 2,
      "output" : "test-output-axe-crop-steps-16"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-axe-crop-steps",
      "connector_name" : "test-connector-name-axe-crop-steps",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-axe-crop-steps",
      "resource" : "test-resource-axe-crop-steps",
      "state" : 1,
      "output" : "test-output-axe-crop-steps-17"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-axe-crop-steps",
      "connector_name" : "test-connector-name-axe-crop-steps",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-axe-crop-steps",
      "resource" : "test-resource-axe-crop-steps",
      "state" : 2,
      "output" : "test-output-axe-crop-steps-18"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-axe-crop-steps",
      "connector_name" : "test-connector-name-axe-crop-steps",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-axe-crop-steps",
      "resource" : "test-resource-axe-crop-steps",
      "state" : 3,
      "output" : "test-output-axe-crop-steps-19"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-axe-crop-steps",
      "connector_name" : "test-connector-name-axe-crop-steps",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-axe-crop-steps",
      "resource" : "test-resource-axe-crop-steps",
      "state" : 2,
      "output" : "test-output-axe-crop-steps-20"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-axe-crop-steps",
      "connector_name" : "test-connector-name-axe-crop-steps",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-axe-crop-steps",
      "resource" : "test-resource-axe-crop-steps",
      "state" : 1,
      "output" : "test-output-axe-crop-steps-21"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-axe-crop-steps",
      "connector_name" : "test-connector-name-axe-crop-steps",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-axe-crop-steps",
      "resource" : "test-resource-axe-crop-steps",
      "state" : 2,
      "output" : "test-output-axe-crop-steps-22"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-axe-crop-steps",
      "connector_name" : "test-connector-name-axe-crop-steps",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-axe-crop-steps",
      "resource" : "test-resource-axe-crop-steps",
      "state" : 3,
      "output" : "test-output-axe-crop-steps-23"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-axe-crop-steps",
      "connector_name" : "test-connector-name-axe-crop-steps",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-axe-crop-steps",
      "resource" : "test-resource-axe-crop-steps",
      "state" : 2,
      "output" : "test-output-axe-crop-steps-24"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-axe-crop-steps",
      "connector_name" : "test-connector-name-axe-crop-steps",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-axe-crop-steps",
      "resource" : "test-resource-axe-crop-steps",
      "state" : 1,
      "output" : "test-output-axe-crop-steps-25"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-axe-crop-steps",
      "connector_name" : "test-connector-name-axe-crop-steps",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-axe-crop-steps",
      "resource" : "test-resource-axe-crop-steps",
      "state" : 2,
      "output" : "test-output-axe-crop-steps-26"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-axe-crop-steps",
      "connector_name" : "test-connector-name-axe-crop-steps",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-axe-crop-steps",
      "resource" : "test-resource-axe-crop-steps",
      "state" : 3,
      "output" : "test-output-axe-crop-steps-27"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-axe-crop-steps",
      "connector_name" : "test-connector-name-axe-crop-steps",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-axe-crop-steps",
      "resource" : "test-resource-axe-crop-steps",
      "state" : 2,
      "output" : "test-output-axe-crop-steps-28"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-axe-crop-steps",
      "connector_name" : "test-connector-name-axe-crop-steps",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-axe-crop-steps",
      "resource" : "test-resource-axe-crop-steps",
      "state" : 1,
      "output" : "test-output-axe-crop-steps-29"
    }
    """
    When I send an event and wait the end of event processing:
    """json
    {
      "connector" : "test-connector-axe-crop-steps",
      "connector_name" : "test-connector-name-axe-crop-steps",
      "source_type" : "resource",
      "event_type" : "check",
      "component" :  "test-component-axe-crop-steps",
      "resource" : "test-resource-axe-crop-steps",
      "state" : 2,
      "output" : "test-output-axe-crop-steps-30"
    }
    """
    When I do GET /api/v4/alarms?search=test-resource-axe-crop-steps
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
            "component": "test-component-axe-crop-steps",
            "connector": "test-connector-axe-crop-steps",
            "connector_name": "test-connector-name-axe-crop-steps",
            "resource": "test-resource-axe-crop-steps",
            "total_state_changes": 30
          }
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """
    When I do POST /api/v4/alarm-details:
    """json
    [
      {
        "_id": "{{ (index .lastResponse.data 0)._id }}",
        "steps": {
          "page": 1,
          "limit": 30
        }
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "data": {
          "steps": {
            "data": [
              {
                "_t": "stateinc",
                "val": 1,
                "m": "test-output-axe-crop-steps-1"
              },
              {
                "_t": "statusinc",
                "val": 1,
                "m": "test-output-axe-crop-steps-1"
              },
              {
                "_t": "statecounter",
                "m": "test-output-axe-crop-steps-1",
                "statecounter": {
                 "state:1": 2,
                 "state:2": 5,
                 "state:3": 2,
                 "statechanges": 9,
                 "statedec": 4,
                 "stateinc": 5
                }
              },
              {
                "_t": "stateinc",
                "val": 3,
                "m": "test-output-axe-crop-steps-11"
              },
              {
                "_t": "statedec",
                "val": 2,
                "m": "test-output-axe-crop-steps-12"
              },
              {
                "_t": "statedec",
                "val": 1,
                "m": "test-output-axe-crop-steps-13"
              },
              {
                "_t": "stateinc",
                "val": 2,
                "m": "test-output-axe-crop-steps-14"
              },
              {
                "_t": "stateinc",
                "val": 3,
                "m": "test-output-axe-crop-steps-15"
              },
              {
                "_t": "statedec",
                "val": 2,
                "m": "test-output-axe-crop-steps-16"
              },
              {
                "_t": "statedec",
                "val": 1,
                "m": "test-output-axe-crop-steps-17"
              },
              {
                "_t": "stateinc",
                "val": 2,
                "m": "test-output-axe-crop-steps-18"
              },
              {
                "_t": "stateinc",
                "val": 3,
                "m": "test-output-axe-crop-steps-19"
              },
              {
                "_t": "statedec",
                "val": 2,
                "m": "test-output-axe-crop-steps-20"
              },
              {
                "_t": "statedec",
                "val": 1,
                "m": "test-output-axe-crop-steps-21"
              },
              {
                "_t": "stateinc",
                "val": 2,
                "m": "test-output-axe-crop-steps-22"
              },
              {
                "_t": "stateinc",
                "val": 3,
                "m": "test-output-axe-crop-steps-23"
              },
              {
                "_t": "statedec",
                "val": 2,
                "m": "test-output-axe-crop-steps-24"
              },
              {
                "_t": "statedec",
                "val": 1,
                "m": "test-output-axe-crop-steps-25"
              },
              {
                "_t": "stateinc",
                "val": 2,
                "m": "test-output-axe-crop-steps-26"
              },
              {
                "_t": "stateinc",
                "val": 3,
                "m": "test-output-axe-crop-steps-27"
              },
              {
                "_t": "statedec",
                "val": 2,
                "m": "test-output-axe-crop-steps-28"
              },
              {
                "_t": "statedec",
                "val": 1,
                "m": "test-output-axe-crop-steps-29"
              },
              {
                "_t": "stateinc",
                "val": 2,
                "m": "test-output-axe-crop-steps-30"
              }
            ],
            "meta": {
              "page": 1,
              "page_count": 1,
              "per_page": 30,
              "total_count": 23
            }
          }
        }
      }
    ]
    """

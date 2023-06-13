Feature: create entities on event
  I need to be able to process event on event

  @concurrent
  Scenario: given resource check event and truncate output and long_output field
    Given I am admin
    When I send an event and wait the end of event processing:
    """json
    {
      "connector": "test-connector-che-process-1",
      "connector_name": "test-connector-name-che-process-1",
      "source_type": "resource",
      "event_type": "check",
      "component": "test-component-che-process-1",
      "resource": "test-resource-che-process-1",
      "state": 2,
      "output": "This text is longer than 255 chars. This text is longer than 255 chars. This text is longer than 255 chars. This text is longer than 255 chars. This text is longer than 255 chars. This text is longer than 255 chars....This text should be truncated at hereFAILED",
      "long_output": "This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer......This text should be truncated at hereFAILED"
    }
    """
    When I do GET /api/v4/alarms?search=process-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "v": {
              "long_output": "This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer than 1024 chars. This text is longer......This text should be truncated at here",
              "output": "This text is longer than 255 chars. This text is longer than 255 chars. This text is longer than 255 chars. This text is longer than 255 chars. This text is longer than 255 chars. This text is longer than 255 chars....This text should be truncated at here"
          }
        }
      ]
    }
    """

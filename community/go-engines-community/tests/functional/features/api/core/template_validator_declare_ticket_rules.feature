Feature: Validate templates for declare ticket rules
  I need to be able to validate templates for declare ticket rules

  @concurrent
  Scenario: given validate template request and no auth should not allow access
    When I do POST /api/v4/template-validator/declare-ticket-rules
    Then the response code should be 401

  @concurrent
  Scenario: given validate template request should return success
    When I am admin
    When I do POST /api/v4/template-validator/declare-ticket-rules:
    """
    {
      "text": "{{ `{{ range .Alarms}} {{ end }}` }}"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "is_valid": true
    }
    """
    Then the response key "warnings" should not exist

  @concurrent
  Scenario: given validate template request with unexpected variable should return error
    When I am admin
    When I do POST /api/v4/template-validator/declare-ticket-rules:
    """
    {
      "text": "http://localhost/{{ `{{.Alarmmm}}` }}"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "is_valid": false,
      "err": {
        "line": 1,
        "position": 19,
        "type": 1,
        "message": "No such variable \".Alarmmm\"",
        "var": ".Alarmmm"
      }
    }
    """

  @concurrent
  Scenario: given validate template request with unexpected variable and new lines should return error and valid line value
    When I am admin
    When I do POST /api/v4/template-validator/declare-ticket-rules:
    """
    {
      "text": "test\ntest\ntest\n{{ `{{.Alarmmm}}` }}"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "is_valid": false,
      "err": {
        "line": 4,
        "position": 2,
        "type": 1,
        "message": "No such variable \".Alarmmm\"",
        "var": ".Alarmmm"
      }
    }
    """

  @concurrent
  Scenario: given validate template request with unexpected secondary variable should return error
    When I am admin
    When I do POST /api/v4/template-validator/declare-ticket-rules:
    """
    {
      "text": "{{ `{{ range .Alarms }} {{ .Alarm.Value.Some }} {{ end }}` }}"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "is_valid": false,
      "err": {
        "line": 1,
        "position": 29,
        "type": 2,
        "message": "Invalid variable \"Some\"",
        "var": "Some"
      }
    }
    """

  @concurrent
  Scenario: given validate template request with unexpected block should return error
    When I am admin
    When I do POST /api/v4/template-validator/declare-ticket-rules:
    """
    {
      "text": "test {{ `{{ end }}` }}"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "is_valid": false,
      "err": {
        "line": 1,
        "type": 3,
        "message": "Function or block is missing"
      }
    }
    """

  @concurrent
  Scenario: given validate template request with unexpected symbol should return error
    When I am admin
    When I do POST /api/v4/template-validator/declare-ticket-rules:
    """
    {
      "text": "test {{ `{{ range .Alarms } {{ end }}` }}"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "is_valid": false,
      "err": {
        "line": 1,
        "type": 4,
        "message": "Unexpected \"}\""
      }
    }
    """

  @concurrent
  Scenario: given validate template request with unexpected function should return error
    When I am admin
    When I do POST /api/v4/template-validator/declare-ticket-rules:
    """
    {
      "text": "test {{ `{{ rangee .Alarms }} {{ end }}` }}"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "is_valid": false,
      "err": {
        "line": 1,
        "type": 5,
        "message": "Invalid function \"rangee\""
      }
    }
    """

  @concurrent
  Scenario: given validate template request with unexpected EOF should return error
    When I am admin
    When I do POST /api/v4/template-validator/declare-ticket-rules:
    """
    {
      "text": "test {{ `{{ range .Alarms }}` }}"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "is_valid": false,
      "err": {
        "line": 1,
        "type": 6,
        "message": "Parsing error: invalid template"
      }
    }
    """

  @concurrent
  Scenario: given validate template request with undefined error should return error
    When I am admin
    When I do POST /api/v4/template-validator/declare-ticket-rules:
    """
    {
      "text": "test {{ `{{ break }}` }}"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "is_valid": false,
      "err": {
        "line": 1,
        "type": 0,
        "message": "{{ `{{break}} outside {{range}}`  }}"
      }
    }
    """

  @concurrent
  Scenario: given validate template request with unsafe mapp access should return warning
    When I am admin
    When I do POST /api/v4/template-validator/declare-ticket-rules:
    """
    {
      "text": "test {{ `{ index .Response \"test\" }}` }}"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "is_valid": true,
      "warnings": [
        {
          "type": 0,
          "message": "Variable are out of a template block",
          "var": ".Response"
        }
      ]
    }
    """

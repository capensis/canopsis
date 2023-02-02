Feature: Validate templates for declare ticket rules
  I need to be able to validate templates for declare ticket rules

  @concurrent
  Scenario: given validate template request and no auth should not allow access
    When I do POST /api/v4/template-validator/scenarios
    Then the response code should be 401

  @concurrent
  Scenario: given validate template request should return success
    When I am admin
    When I do POST /api/v4/template-validator/scenarios:
    """
    {
      "text": "{{ `{{ .Alarm.Value.Output}}` }}"
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
  Scenario: given validate template request with unexpected block should return error
    When I am admin
    When I do POST /api/v4/template-validator/scenarios:
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
        "type": 1,
        "message": "Function or block is missing"
      }
    }
    """

  @concurrent
  Scenario: given validate template request with unexpected symbol should return error
    When I am admin
    When I do POST /api/v4/template-validator/scenarios:
    """
    {
      "text": "test {{ `{{ range .Children } {{ end }}` }}"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "is_valid": false,
      "err": {
        "line": 1,
        "type": 2,
        "message": "Unexpected \"}\""
      }
    }
    """

  @concurrent
  Scenario: given validate template request with unexpected function should return error
    When I am admin
    When I do POST /api/v4/template-validator/scenarios:
    """
    {
      "text": "test {{ `{{ rangee .Children }} {{ end }}` }}"
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
        "message": "Invalid function \"rangee\""
      }
    }
    """

  @concurrent
  Scenario: given validate template request with unexpected EOF should return error
    When I am admin
    When I do POST /api/v4/template-validator/scenarios:
    """
    {
      "text": "test {{ `{{ range .Children }}` }}"
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
        "message": "Parsing error: invalid template"
      }
    }
    """

  @concurrent
  Scenario: given validate template request with undefined error should return error
    When I am admin
    When I do POST /api/v4/template-validator/scenarios:
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
  Scenario: given validate template request with var out of block
    When I am admin
    When I do POST /api/v4/template-validator/scenarios:
    """
    {
      "text": "{{ `{{ .AdditionalData.Author }} test { index .Response \"test\" }} {{ .Alarm.Value.Output }}` }}"
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
          "message": "Variable is out of a template block",
          "var": ".Response"
        }
      ]
    }
    """

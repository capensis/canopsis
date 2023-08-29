Feature: Validate templates for event filter rules
  I need to be able to validate templates for event filter rules

  @concurrent
  Scenario: given validate template request and no auth should not allow access
    When I do POST /api/v4/template-validator/event-filter-rules
    Then the response code should be 401

  @concurrent
  Scenario: given validate template request should return success
    When I am admin
    When I do POST /api/v4/template-validator/event-filter-rules:
    """json
    [
      {
        "text": "{{ `{{ .Event.Component }} {{ .ExternalData.component.title }}` }}"
      },
      {
        "text": "{{ `{{ .Event.Component }} {{ .ExternalData.component.title }}` }}"
      }
    ]
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    [
      {
        "is_valid": true
      },
      {
        "is_valid": true
      }
    ]
    """
    Then the response key "0.warnings" should not exist
    Then the response key "1.warnings" should not exist

  @concurrent
  Scenario: given validate template request with unexpected block should return error
    When I am admin
    When I do POST /api/v4/template-validator/event-filter-rules:
    """json
    [
      {
        "text": "test {{ `{{ end }}` }}"
      }
    ]
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    [
      {
        "is_valid": false,
        "err": {
          "line": 1,
          "type": 1,
          "message": "Function or block is missing"
        }
      }
    ]
    """

  @concurrent
  Scenario: given validate template request with unexpected symbol should return error
    When I am admin
    When I do POST /api/v4/template-validator/event-filter-rules:
    """json
    [
      {
        "text": "test {{ `{{ range .Event.ExtraInfos } {{ end }}` }}"
      }
    ]
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    [
      {
        "is_valid": false,
        "err": {
          "line": 1,
          "type": 2,
          "message": "Unexpected \"}\""
        }
      }
    ]
    """

  @concurrent
  Scenario: given validate template request with unexpected function should return error
    When I am admin
    When I do POST /api/v4/template-validator/event-filter-rules:
    """json
    [
      {
        "text": "test {{ `{{ rangee .Event.ExtraInfos }} {{ end }}` }}"
      }
    ]
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    [
      {
        "is_valid": false,
        "err": {
          "line": 1,
          "type": 3,
          "message": "Invalid function \"rangee\""
        }
      }
    ]
    """

  @concurrent
  Scenario: given validate template request with unexpected EOF should return error
    When I am admin
    When I do POST /api/v4/template-validator/event-filter-rules:
    """json
    [
      {
        "text": "test {{ `{{ range .Event.ExtraInfos }}` }}"
      }
    ]
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    [
      {
        "is_valid": false,
        "err": {
          "line": 1,
          "type": 4,
          "message": "Parsing error: invalid template"
        }
      }
    ]
    """

  @concurrent
  Scenario: given validate template request with undefined error should return error
    When I am admin
    When I do POST /api/v4/template-validator/event-filter-rules:
    """json
    [
      {
        "text": "test {{ `{{ break }}` }}"
      }
    ]
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    [
      {
        "is_valid": false,
        "err": {
          "line": 1,
          "type": 0,
          "message": "{{ `{{break}} outside {{range}}`  }}"
        }
      }
    ]
    """

  @concurrent
  Scenario: given validate template request with var out of block
    When I am admin
    When I do POST /api/v4/template-validator/event-filter-rules:
    """json
    [
      {
        "text": "test {{ `test { .Event.Output }}` }}"
      }
    ]
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    [
      {
        "is_valid": true,
        "warnings": [
          {
            "type": 0,
            "message": "Variable is out of a template block",
            "var": ".Event"
          }
        ]
      }
    ]
    """

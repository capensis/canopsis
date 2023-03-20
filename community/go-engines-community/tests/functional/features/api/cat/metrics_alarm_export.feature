Feature: Export alarm metrics
  I need to be able to export alarm metrics
  Only admin should be able to export alarm metrics

  Scenario: given export request should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics-export/alarm?parameters[]=created_alarms&sampling=day&from={{ parseTime "20-11-2021 00:00" }}&to={{ parseTime "24-11-2021 00:00" }}&filter=test-kpi-filter-to-alarm-metrics-get
    Then the response code should be 200
    When I save response exportID={{ .lastResponse._id }}
    When I do GET /api/v4/cat/metrics-export/{{ .exportID }} until response code is 200 and body contains:
    """json
    {
       "status": 1
    }
    """
    When I do GET /api/v4/cat/metrics-export/{{ .exportID }}/download
    Then the response code should be 200
    Then the response raw body should be:
    """csv
    metric,timestamp,value
    created_alarms,{{ parseTime "20-11-2021 00:00" }},0
    created_alarms,{{ parseTime "21-11-2021 00:00" }},1
    created_alarms,{{ parseTime "22-11-2021 00:00" }},0
    created_alarms,{{ parseTime "23-11-2021 00:00" }},3
    created_alarms,{{ parseTime "24-11-2021 00:00" }},0

    """

  Scenario: given export request with empty interval should return metrics with zeros
    When I am admin
    When I do POST /api/v4/cat/metrics-export/alarm?parameters[]=created_alarms&sampling=day&from={{ parseTime "06-09-2020 00:00" }}&to={{ parseTime "08-09-2020 00:00" }}&filter=test-kpi-filter-to-alarm-metrics-get
    Then the response code should be 200
    When I save response exportID={{ .lastResponse._id }}
    When I do GET /api/v4/cat/metrics-export/{{ .exportID }} until response code is 200 and body contains:
    """json
    {
       "status": 1
    }
    """
    When I do GET /api/v4/cat/metrics-export/{{ .exportID }}/download
    Then the response code should be 200
    Then the response raw body should be:
    """csv
    metric,timestamp,value
    created_alarms,{{ parseTime "06-09-2020 00:00" }},0
    created_alarms,{{ parseTime "07-09-2020 00:00" }},0
    created_alarms,{{ parseTime "08-09-2020 00:00" }},0

    """

  Scenario: given export request with filter by entity infos should return metrics
    When I am admin
    When I do POST /api/v4/cat/metrics-export/alarm?parameters[]=created_alarms&sampling=day&from={{ parseTime "20-11-2021 00:00" }}&to={{ parseTime "24-11-2021 00:00" }}&filter=test-kpi-filter-to-alarm-metrics-get-by-entity-infos
    Then the response code should be 200
    When I save response exportID={{ .lastResponse._id }}
    When I do GET /api/v4/cat/metrics-export/{{ .exportID }} until response code is 200 and body contains:
    """json
    {
       "status": 1
    }
    """
    When I do GET /api/v4/cat/metrics-export/{{ .exportID }}/download
    Then the response code should be 200
    Then the response raw body should be:
    """csv
    metric,timestamp,value
    created_alarms,{{ parseTime "20-11-2021 00:00" }},0
    created_alarms,{{ parseTime "21-11-2021 00:00" }},0
    created_alarms,{{ parseTime "22-11-2021 00:00" }},0
    created_alarms,{{ parseTime "23-11-2021 00:00" }},2
    created_alarms,{{ parseTime "24-11-2021 00:00" }},0

    """

  Scenario: given export request with invalid query params should return bad request
    When I am admin
    When I do POST /api/v4/cat/metrics-export/alarm
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "from": "From is missing.",
        "parameters[]": "Parameters is missing.",
        "sampling": "Sampling is missing.",
        "to": "To is missing."
      }
    }
    """
    When I do POST /api/v4/cat/metrics-export/alarm?filter=not-exist&from={{ now }}&to={{ now }}&sampling=day&parameters[]=created_alarms
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "filter": "Filter \"not-exist\" not found."
      }
    }
    """
    When I do POST /api/v4/cat/metrics-export/alarm?parameters[]=not-exist&from={{ now }}&to={{ now }}&sampling=day
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "parameter.0": "Parameter \"not-exist\" is not supported."
      }
    }
    """
    When I do POST /api/v4/cat/metrics-export/alarm?parameters[]=total_user_activity&from={{ now }}&to={{ now }}&sampling=day
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "parameter.0": "Parameter \"total_user_activity\" is not supported."
      }
    }
    """
    When I do POST /api/v4/cat/metrics-export/alarm?sampling=not-exist&from={{ now }}&to={{ now }}&parameters[]=created_alarms
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "sampling": "Sampling must be one of [hour day week month]."
      }
    }
    """

  Scenario: given export request and no auth user should not allow access
    When I do POST /api/v4/cat/metrics-export/alarm
    Then the response code should be 401

  Scenario: given export request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/cat/metrics-export/alarm
    Then the response code should be 403

  Scenario: given export request with all parameters should return all metrics
    When I am admin
    When I do POST /api/v4/cat/metrics-export/alarm?parameters[]=created_alarms&parameters[]=active_alarms&parameters[]=non_displayed_alarms&parameters[]=instruction_alarms&parameters[]=pbehavior_alarms&parameters[]=correlation_alarms&parameters[]=ack_alarms&parameters[]=cancel_ack_alarms&parameters[]=ack_active_alarms&parameters[]=ticket_active_alarms&parameters[]=without_ticket_active_alarms&parameters[]=ratio_correlation&parameters[]=ratio_instructions&parameters[]=ratio_tickets&parameters[]=ratio_non_displayed&parameters[]=average_ack&parameters[]=average_resolve&sampling=day&from={{ parseTime "22-11-2021 00:00" }}&to={{ parseTime "24-11-2021 00:00" }}&filter=test-kpi-filter-to-all-alarm-metrics-get
    Then the response code should be 200
    When I save response exportID={{ .lastResponse._id }}
    When I do GET /api/v4/cat/metrics-export/{{ .exportID }} until response code is 200 and body contains:
    """json
    {
       "status": 1
    }
    """
    When I do GET /api/v4/cat/metrics-export/{{ .exportID }}/download
    Then the response code should be 200
    Then the response raw body should be:
    """csv
    metric,timestamp,value
    created_alarms,{{ parseTime "22-11-2021 00:00" }},0
    created_alarms,{{ parseTime "23-11-2021 00:00" }},3
    created_alarms,{{ parseTime "24-11-2021 00:00" }},0
    active_alarms,{{ parseTime "22-11-2021 00:00" }},0
    active_alarms,{{ parseTime "23-11-2021 00:00" }},3
    active_alarms,{{ parseTime "24-11-2021 00:00" }},3
    non_displayed_alarms,{{ parseTime "22-11-2021 00:00" }},0
    non_displayed_alarms,{{ parseTime "23-11-2021 00:00" }},1
    non_displayed_alarms,{{ parseTime "24-11-2021 00:00" }},0
    instruction_alarms,{{ parseTime "22-11-2021 00:00" }},0
    instruction_alarms,{{ parseTime "23-11-2021 00:00" }},1
    instruction_alarms,{{ parseTime "24-11-2021 00:00" }},0
    pbehavior_alarms,{{ parseTime "22-11-2021 00:00" }},0
    pbehavior_alarms,{{ parseTime "23-11-2021 00:00" }},1
    pbehavior_alarms,{{ parseTime "24-11-2021 00:00" }},0
    correlation_alarms,{{ parseTime "22-11-2021 00:00" }},0
    correlation_alarms,{{ parseTime "23-11-2021 00:00" }},1
    correlation_alarms,{{ parseTime "24-11-2021 00:00" }},0
    ack_alarms,{{ parseTime "22-11-2021 00:00" }},0
    ack_alarms,{{ parseTime "23-11-2021 00:00" }},2
    ack_alarms,{{ parseTime "24-11-2021 00:00" }},0
    cancel_ack_alarms,{{ parseTime "22-11-2021 00:00" }},0
    cancel_ack_alarms,{{ parseTime "23-11-2021 00:00" }},1
    cancel_ack_alarms,{{ parseTime "24-11-2021 00:00" }},0
    ack_active_alarms,{{ parseTime "22-11-2021 00:00" }},0
    ack_active_alarms,{{ parseTime "23-11-2021 00:00" }},1
    ack_active_alarms,{{ parseTime "24-11-2021 00:00" }},1
    ticket_active_alarms,{{ parseTime "22-11-2021 00:00" }},0
    ticket_active_alarms,{{ parseTime "23-11-2021 00:00" }},1
    ticket_active_alarms,{{ parseTime "24-11-2021 00:00" }},1
    without_ticket_active_alarms,{{ parseTime "22-11-2021 00:00" }},0
    without_ticket_active_alarms,{{ parseTime "23-11-2021 00:00" }},2
    without_ticket_active_alarms,{{ parseTime "24-11-2021 00:00" }},2
    ratio_correlation,{{ parseTime "22-11-2021 00:00" }},0
    ratio_correlation,{{ parseTime "23-11-2021 00:00" }},33.33
    ratio_correlation,{{ parseTime "24-11-2021 00:00" }},33.33
    ratio_instructions,{{ parseTime "22-11-2021 00:00" }},0
    ratio_instructions,{{ parseTime "23-11-2021 00:00" }},33.33
    ratio_instructions,{{ parseTime "24-11-2021 00:00" }},33.33
    ratio_tickets,{{ parseTime "22-11-2021 00:00" }},0
    ratio_tickets,{{ parseTime "23-11-2021 00:00" }},33.33
    ratio_tickets,{{ parseTime "24-11-2021 00:00" }},33.33
    ratio_non_displayed,{{ parseTime "22-11-2021 00:00" }},0
    ratio_non_displayed,{{ parseTime "23-11-2021 00:00" }},33.33
    ratio_non_displayed,{{ parseTime "24-11-2021 00:00" }},33.33
    average_ack,{{ parseTime "22-11-2021 00:00" }},0
    average_ack,{{ parseTime "23-11-2021 00:00" }},500
    average_ack,{{ parseTime "24-11-2021 00:00" }},0
    average_resolve,{{ parseTime "22-11-2021 00:00" }},0
    average_resolve,{{ parseTime "23-11-2021 00:00" }},1000
    average_resolve,{{ parseTime "24-11-2021 00:00" }},0

    """

import { describe, expect, it } from 'vitest';

import {
  EVENT_FILTER_TYPES,
  EXTERNAL_DATA_TYPES,
  REQUEST_AUTH_TOKEN_TYPES,
  REQUEST_AUTH_TYPES,
} from '@/constants';

import { eventFilterToForm, formToEventFilter } from '@/helpers/entities/event-filter/rule/form';

describe('formToEventFilter', () => {
  it('keeps the recurrence rule in the event-filter request', () => {
    const eventFilter = formToEventFilter({
      type: EVENT_FILTER_TYPES.drop,
      rrule: 'FREQ=DAILY',
      config: {},
      patterns: {},
      exdates: [],
      exceptions: [],
      external_data: {},
    });

    expect(eventFilter.rrule).toBe('FREQ=DAILY');
  });

  it('converts an external API auth token to the event-filter form', () => {
    const form = eventFilterToForm({
      external_data: [
        {
          type: EXTERNAL_DATA_TYPES.api,
          request: {},
          auth_token: {
            rule: 'token-rule-id',
            header: 'Authorization: Bearer',
          },
        },
      ],
    });

    expect(form.external_data[0].request.auth.type).toBe(REQUEST_AUTH_TYPES.token);
    expect(form.external_data[0].auth_token).toEqual({
      type: REQUEST_AUTH_TOKEN_TYPES.headerAuthorizationBearer,
      parameter: '',
      rule: 'token-rule-id',
    });
  });

  it('adds an external API auth token to the event-filter request', () => {
    const eventFilter = formToEventFilter({
      type: EVENT_FILTER_TYPES.enrichment,
      config: {
        actions: [],
      },
      patterns: {},
      exdates: [],
      exceptions: [],
      external_data: [
        {
          reference: 'external-api',
          type: EXTERNAL_DATA_TYPES.api,
          request: {
            auth: {
              type: REQUEST_AUTH_TYPES.token,
            },
            headers: [],
            retry_delay: {},
            timeout: {},
          },
          auth_token: {
            type: REQUEST_AUTH_TOKEN_TYPES.headerAuthorizationBearer,
            parameter: '',
            rule: 'token-rule-id',
          },
        },
      ],
    });

    expect(eventFilter.external_data[0].auth_token).toEqual({
      rule: 'token-rule-id',
      header: 'Authorization: Bearer',
    });
  });

  it('does not add auth token form data to table external data', () => {
    const eventFilter = formToEventFilter({
      type: EVENT_FILTER_TYPES.enrichment,
      config: {
        actions: [],
      },
      patterns: {},
      exdates: [],
      exceptions: [],
      external_data: [
        {
          reference: 'external-table',
          type: EXTERNAL_DATA_TYPES.table,
          table: 'table-id',
          conditions: [],
          request: {},
          auth_token: {
            type: '',
            parameter: '',
            rule: '',
          },
        },
      ],
    });

    expect(eventFilter.external_data[0]).not.toHaveProperty('auth_token');
  });
});

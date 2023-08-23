import flushPromises from 'flush-promises';

import { generateRenderer } from '@unit/utils/vue';
import { mockDateNow } from '@unit/utils/mock-hooks';
import { EVENT_ENTITY_TYPES } from '@/constants';

import ExtraDetailsTicket from '@/components/widgets/alarm/columns-formatting/extra-details/extra-details-ticket.vue';

describe('extra-details-ticket', () => {
  const nowTimestamp = 1386435500000;
  const prevDateTimestamp = 1386392400000;
  const prevMonthDateTimestamp = 1375894800000;

  mockDateNow(nowTimestamp);

  const tickets = [
    {
      a: 'ticket-author-1',
      t: prevDateTimestamp,
      ticket_rule_name: 'ticket-rule-name-1',
      ticket_comment: 'ticket-comment-1',
      ticket: 'ticket-1',
      _t: EVENT_ENTITY_TYPES.declareTicket,
    },
    {
      a: 'ticket-author-2',
      t: prevMonthDateTimestamp,
      ticket_rule_name: 'ticket-rule-name-2',
    },
  ];

  const snapshotFactory = generateRenderer(ExtraDetailsTicket, {

    attachTo: document.body,
  });

  it('Renders `extra-details-ticket` with full ticket', async () => {
    snapshotFactory({
      propsData: {
        tickets,
      },
    });

    await flushPromises();

    expect(document.body.innerHTML).toMatchSnapshot();
  });

  it('Renders `extra-details-ticket` without value', async () => {
    snapshotFactory({
      propsData: {
        tickets,
      },
    });

    await flushPromises();

    expect(document.body.innerHTML).toMatchSnapshot();
  });

  it('Renders `extra-details-ticket` with date in previous month', async () => {
    snapshotFactory({
      propsData: {
        tickets,
      },
    });

    await flushPromises();

    expect(document.body.innerHTML).toMatchSnapshot();
  });
});

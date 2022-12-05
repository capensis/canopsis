import { omit } from 'lodash';
import flushPromises from 'flush-promises';

import { createVueInstance, generateRenderer } from '@unit/utils/vue';
import { mockDateNow } from '@unit/utils/mock-hooks';

import ExtraDetailsTicket from '@/components/widgets/alarm/columns-formatting/extra-details/extra-details-ticket.vue';

const localVue = createVueInstance();

describe('extra-details-ticket', () => {
  const nowTimestamp = 1386435500000;
  const prevDateTimestamp = 1386392400000;
  const prevMonthDateTimestamp = 1375894800000;

  mockDateNow(nowTimestamp);

  const ticket = {
    a: 'ticket-author',
    t: prevDateTimestamp,
    val: 'ticket-message',
  };

  const snapshotFactory = generateRenderer(ExtraDetailsTicket, {
    localVue,
    attachTo: document.body,
  });

  it('Renders `extra-details-ticket` with full ticket', async () => {
    snapshotFactory({
      propsData: {
        ticket,
      },
    });

    await flushPromises();

    expect(document.body.innerHTML).toMatchSnapshot();
  });

  it('Renders `extra-details-ticket` without value', async () => {
    snapshotFactory({
      propsData: {
        ticket: omit(ticket, ['val']),
      },
    });

    await flushPromises();

    expect(document.body.innerHTML).toMatchSnapshot();
  });

  it('Renders `extra-details-ticket` with date in previous month', async () => {
    snapshotFactory({
      propsData: {
        ticket: {
          ...ticket,
          t: prevMonthDateTimestamp,
        },
      },
    });

    await flushPromises();

    expect(document.body.innerHTML).toMatchSnapshot();
  });
});

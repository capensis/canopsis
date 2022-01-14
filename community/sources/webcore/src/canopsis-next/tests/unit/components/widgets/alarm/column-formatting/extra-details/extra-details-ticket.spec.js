import { omit } from 'lodash';

import { mount, createVueInstance } from '@unit/utils/vue';
import { mockDateNow } from '@unit/utils/mock-hooks';

import ExtraDetailsTicket from '@/components/widgets/alarm/columns-formatting/extra-details/extra-details-ticket.vue';

const localVue = createVueInstance();

const snapshotFactory = (options = {}) => mount(ExtraDetailsTicket, {
  localVue,

  ...options,
});

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

  it('Renders `extra-details-ticket` with full ticket', () => {
    const wrapper = snapshotFactory({
      propsData: {
        ticket,
      },
    });

    const tooltipContent = wrapper.findTooltip();

    expect(wrapper.element).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });

  it('Renders `extra-details-ticket` without value', () => {
    const wrapper = snapshotFactory({
      propsData: {
        ticket: omit(ticket, ['val']),
      },
    });

    const tooltipContent = wrapper.findTooltip();

    expect(wrapper.element).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });

  it('Renders `extra-details-ticket` with date in previous month', () => {
    const wrapper = snapshotFactory({
      propsData: {
        ticket: {
          ...ticket,
          t: prevMonthDateTimestamp,
        },
      },
    });

    const tooltipContent = wrapper.findTooltip();

    expect(wrapper.element).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });
});

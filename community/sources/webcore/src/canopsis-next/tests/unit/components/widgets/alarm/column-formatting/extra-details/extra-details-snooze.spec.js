import { omit } from 'lodash';

import { mount, createVueInstance } from '@unit/utils/vue';
import { mockDateNow } from '@unit/utils/mock-hooks';

import ExtraDetailsSnooze from '@/components/widgets/alarm/columns-formatting/extra-details/extra-details-snooze.vue';

const localVue = createVueInstance();

const snapshotFactory = (options = {}) => mount(ExtraDetailsSnooze, {
  localVue,

  ...options,
});

describe('extra-details-snooze', () => {
  const nowTimestamp = 1386435500000;
  const prevDateTimestamp = 1386392400000;
  const prevMonthDateTimestamp = 1375894800000;

  mockDateNow(nowTimestamp);

  const snooze = {
    a: 'snooze-author',
    t: prevDateTimestamp,
    initiator: 'snooze-initiator',
    m: 'snooze-message',
    val: prevDateTimestamp,
  };

  it('Renders `extra-details-snooze` with full snooze', () => {
    const wrapper = snapshotFactory({
      propsData: {
        snooze,
      },
    });

    const tooltipContent = wrapper.findTooltip();

    expect(wrapper.element).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });

  it('Renders `extra-details-snooze` without initiator', () => {
    const wrapper = snapshotFactory({
      propsData: {
        snooze: omit(snooze, ['initiator']),
      },
    });

    const tooltipContent = wrapper.findTooltip();

    expect(wrapper.element).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });

  it('Renders `extra-details-snooze` without message', () => {
    const wrapper = snapshotFactory({
      propsData: {
        snooze: omit(snooze, ['m']),
      },
    });

    const tooltipContent = wrapper.findTooltip();

    expect(wrapper.element).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });

  it('Renders `extra-details-snooze` with date in previous month', () => {
    const wrapper = snapshotFactory({
      propsData: {
        snooze: {
          ...snooze,
          t: prevMonthDateTimestamp,
          val: prevMonthDateTimestamp,
        },
      },
    });

    const tooltipContent = wrapper.findTooltip();

    expect(wrapper.element).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });
});

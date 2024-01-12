import { omit } from 'lodash';

import { generateRenderer } from '@unit/utils/vue';
import { mockDateNow } from '@unit/utils/mock-hooks';

import ExtraDetailsSnooze from '@/components/widgets/alarm/columns-formatting/extra-details/extra-details-snooze.vue';

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

  const snapshotFactory = generateRenderer(ExtraDetailsSnooze, {
    attachTo: document.body,
  });

  it('Renders `extra-details-snooze` with full snooze', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        snooze,
      },
    });

    expect(wrapper).toMatchSnapshot();
    await wrapper.activateAllTooltips();
    expect(wrapper).toMatchTooltipSnapshot();
  });

  it('Renders `extra-details-snooze` without initiator', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        snooze: omit(snooze, ['initiator']),
      },
    });

    expect(wrapper).toMatchSnapshot();
    await wrapper.activateAllTooltips();
    expect(wrapper).toMatchTooltipSnapshot();
  });

  it('Renders `extra-details-snooze` without message', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        snooze: omit(snooze, ['m']),
      },
    });

    expect(wrapper).toMatchSnapshot();
    await wrapper.activateAllTooltips();
    expect(wrapper).toMatchTooltipSnapshot();
  });

  it('Renders `extra-details-snooze` with date in previous month', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        snooze: {
          ...snooze,
          t: prevMonthDateTimestamp,
          val: prevMonthDateTimestamp,
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
    await wrapper.activateAllTooltips();
    expect(wrapper).toMatchTooltipSnapshot();
  });
});

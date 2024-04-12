import { omit } from 'lodash';

import { generateRenderer } from '@unit/utils/vue';
import { mockDateNow } from '@unit/utils/mock-hooks';

import ExtraDetailsCanceled from '@/components/widgets/alarm/columns-formatting/extra-details/extra-details-canceled.vue';

const stubs = {
  'c-alarm-extra-details-chip': true,
};

describe('extra-details-canceled', () => {
  const nowTimestamp = 1386435500000;
  const prevDateTimestamp = 1386392400000;
  const prevMonthDateTimestamp = 1375894800000;

  mockDateNow(nowTimestamp);

  const canceled = {
    a: 'cancelled-author',
    t: prevDateTimestamp,
    m: 'cancelled-message',
  };

  const snapshotFactory = generateRenderer(ExtraDetailsCanceled, {
    stubs,
    attachTo: document.body,
  });

  it('Renders `extra-details-canceled` with full canceled', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        canceled,
      },
    });

    expect(wrapper).toMatchSnapshot();
    await wrapper.activateAllTooltips();
    expect(wrapper).toMatchTooltipSnapshot();
  });

  it('Renders `extra-details-canceled` without message', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        canceled: omit(canceled, ['m']),
      },
    });

    expect(wrapper).toMatchSnapshot();
    await wrapper.activateAllTooltips();
    expect(wrapper).toMatchTooltipSnapshot();
  });

  it('Renders `extra-details-canceled` with date in previous month', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        canceled: {
          ...canceled,
          t: prevMonthDateTimestamp,
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
    await wrapper.activateAllTooltips();
    expect(wrapper).toMatchTooltipSnapshot();
  });
});

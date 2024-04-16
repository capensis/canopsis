import { omit } from 'lodash';

import { generateRenderer } from '@unit/utils/vue';
import { mockDateNow } from '@unit/utils/mock-hooks';

import ExtraDetailsAck from '@/components/widgets/alarm/columns-formatting/extra-details/extra-details-ack.vue';

const stubs = {
  'c-alarm-extra-details-chip': true,
};

describe('extra-details-ack', () => {
  const nowTimestamp = 1386435500000;
  const prevDateTimestamp = 1386392400000;
  const prevMonthDateTimestamp = 1375894800000;

  mockDateNow(nowTimestamp);

  const snapshotFactory = generateRenderer(ExtraDetailsAck, {
    stubs,
    attachTo: document.body,
  });

  const ack = {
    a: 'ack-author',
    t: prevDateTimestamp,
    initiator: 'ack-initiator',
    m: 'ack-message',
  };

  it('Renders `extra-details-ack` with full ack', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        ack,
      },
    });

    expect(wrapper).toMatchSnapshot();
    await wrapper.activateAllTooltips();
    expect(wrapper).toMatchTooltipSnapshot();
  });

  it('Renders `extra-details-ack` without initiator', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        ack: omit(ack, ['initiator']),
      },
    });

    expect(wrapper).toMatchSnapshot();
    await wrapper.activateAllTooltips();
    expect(wrapper).toMatchTooltipSnapshot();
  });

  it('Renders `extra-details-ack` without message', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        ack: omit(ack, ['m']),
      },
    });

    expect(wrapper).toMatchSnapshot();
    await wrapper.activateAllTooltips();
    expect(wrapper).toMatchTooltipSnapshot();
  });

  it('Renders `extra-details-ack` with date in previous month', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        ack: {
          ...ack,
          t: prevMonthDateTimestamp,
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
    await wrapper.activateAllTooltips();
    expect(wrapper).toMatchTooltipSnapshot();
  });
});

import { omit } from 'lodash';

import { mount, createVueInstance } from '@unit/utils/vue';
import { mockDateNow } from '@unit/utils/mock-hooks';

import ExtraDetailsAck from '@/components/widgets/alarm/columns-formatting/extra-details/extra-details-ack.vue';

const localVue = createVueInstance();

const snapshotFactory = (options = {}) => mount(ExtraDetailsAck, {
  localVue,

  ...options,
});

describe('extra-details-ack', () => {
  const nowTimestamp = 1386435500000;
  const prevDateTimestamp = 1386392400000;
  const prevMonthDateTimestamp = 1375894800000;

  mockDateNow(nowTimestamp);

  const ack = {
    a: 'ack-author',
    t: prevDateTimestamp,
    initiator: 'ack-initiator',
    m: 'ack-message',
  };

  it('Renders `extra-details-ack` with full ack', () => {
    const wrapper = snapshotFactory({
      propsData: {
        ack,
      },
    });

    const tooltipContent = wrapper.findTooltip();

    expect(wrapper.element).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });

  it('Renders `extra-details-ack` without initiator', () => {
    const wrapper = snapshotFactory({
      propsData: {
        ack: omit(ack, ['initiator']),
      },
    });

    const tooltipContent = wrapper.findTooltip();

    expect(wrapper.element).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });

  it('Renders `extra-details-ack` without message', () => {
    const wrapper = snapshotFactory({
      propsData: {
        ack: omit(ack, ['m']),
      },
    });

    const tooltipContent = wrapper.findTooltip();

    expect(wrapper.element).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });

  it('Renders `extra-details-ack` with date in previous month', () => {
    const wrapper = snapshotFactory({
      propsData: {
        ack: {
          ...ack,
          t: prevMonthDateTimestamp,
        },
      },
    });

    const tooltipContent = wrapper.findTooltip();

    expect(wrapper.element).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });
});

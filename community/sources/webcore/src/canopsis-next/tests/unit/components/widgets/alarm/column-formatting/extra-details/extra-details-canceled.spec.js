import { omit } from 'lodash';

import { mount, createVueInstance } from '@unit/utils/vue';
import { mockDateNow } from '@unit/utils/mock-hooks';

import ExtraDetailsCanceled from '@/components/widgets/alarm/columns-formatting/extra-details/extra-details-canceled.vue';

const localVue = createVueInstance();

const snapshotFactory = (options = {}) => mount(ExtraDetailsCanceled, {
  localVue,

  ...options,
});

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

  it('Renders `extra-details-canceled` with full canceled', () => {
    const wrapper = snapshotFactory({
      propsData: {
        canceled,
      },
    });

    const tooltipContent = wrapper.findTooltip();

    expect(wrapper.element).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });

  it('Renders `extra-details-canceled` without message', () => {
    const wrapper = snapshotFactory({
      propsData: {
        canceled: omit(canceled, ['m']),
      },
    });

    const tooltipContent = wrapper.findTooltip();

    expect(wrapper.element).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });

  it('Renders `extra-details-canceled` with date in previous month', () => {
    const wrapper = snapshotFactory({
      propsData: {
        canceled: {
          ...canceled,
          t: prevMonthDateTimestamp,
        },
      },
    });

    const tooltipContent = wrapper.findTooltip();

    expect(wrapper.element).toMatchSnapshot();
    expect(tooltipContent.element).toMatchSnapshot();
  });
});

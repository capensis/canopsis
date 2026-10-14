import { generateRenderer } from '@unit/utils/vue';

import ExtraDetailsCloseDelayValue
  from '@/components/widgets/alarm/columns-formatting/extra-details/extra-details-close-delay-value.vue';

const stubs = {
  'c-alarm-extra-details-chip': true,
  'c-simple-tooltip': true,
};

describe('extra-details-close-delay-value', () => {
  const closeDelayValue = 6;

  const snapshotFactory = generateRenderer(ExtraDetailsCloseDelayValue, {
    stubs,
    attachTo: document.body,
  });

  it('Renders `extra-details-close-delay-value` with full children and rule', () => {
    jest.useFakeTimers();

    const wrapper = snapshotFactory({
      propsData: {
        closeDelayValue,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});

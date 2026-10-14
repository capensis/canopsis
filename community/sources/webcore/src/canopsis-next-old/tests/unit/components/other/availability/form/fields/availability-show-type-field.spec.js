import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createCheckboxInputStub } from '@unit/stubs/input';

import { AVAILABILITY_SHOW_TYPE } from '@/constants';

import AvailabilityShowTypeField from '@/components/other/availability/form/fields/availability-show-type-field.vue';

const stubs = {
  'v-select': createCheckboxInputStub('v-select'),
};

const selectRadioGroup = wrapper => wrapper.find('.v-select');

describe('availability-show-type-field', () => {
  const factory = generateShallowRenderer(AvailabilityShowTypeField, { stubs });
  const snapshotFactory = generateRenderer(AvailabilityShowTypeField);

  test('Show type changed after trigger radio field', () => {
    const wrapper = factory({
      propsData: {
        value: AVAILABILITY_SHOW_TYPE.percent,
      },
    });

    selectRadioGroup(wrapper).triggerCustomEvent('input', AVAILABILITY_SHOW_TYPE.duration);

    expect(wrapper).toEmitInput(AVAILABILITY_SHOW_TYPE.duration);
  });

  test('Renders `availability-show-type-field` with required props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `availability-show-type-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: AVAILABILITY_SHOW_TYPE.duration,
        label: 'Custom label',
        name: 'customName',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});

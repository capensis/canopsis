import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createCheckboxInputStub } from '@unit/stubs/input';

import { AVAILABILITY_SHOW_TYPE } from '@/constants';

import CAvailabilityShowTypeField from '@/components/forms/fields/availability/c-availability-show-type-field.vue';

const stubs = {
  'v-radio-group': createCheckboxInputStub('v-radio-group'),
};

const selectRadioGroup = wrapper => wrapper.find('.v-radio-group');

describe('c-availability-show-type-field', () => {
  const factory = generateShallowRenderer(CAvailabilityShowTypeField, { stubs });
  const snapshotFactory = generateRenderer(CAvailabilityShowTypeField);

  test('Show type changed after trigger radio field', () => {
    const wrapper = factory({
      propsData: {
        value: AVAILABILITY_SHOW_TYPE.percent,
      },
    });

    selectRadioGroup(wrapper).triggerCustomEvent('input', AVAILABILITY_SHOW_TYPE.duration);

    expect(wrapper).toEmitInput(AVAILABILITY_SHOW_TYPE.duration);
  });

  test('Renders `c-availability-show-type-field` with required props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `c-availability-show-type-field` with custom props', () => {
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

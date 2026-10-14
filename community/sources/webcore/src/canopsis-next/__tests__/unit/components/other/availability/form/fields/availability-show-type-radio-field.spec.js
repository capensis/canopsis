import { createCheckboxInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';

import { AVAILABILITY_SHOW_TYPE } from '@/constants';

import AvailabilityShowTypeRadioField from '@/components/other/availability/form/fields/availability-show-type-radio-field.vue';

const stubs = {
  'v-radio-group': createCheckboxInputStub('v-radio-group'),
};

const selectRadioGroup = wrapper => wrapper.find('.v-radio-group');

describe('availability-show-type-radio-field', () => {
  const factory = generateShallowRenderer(AvailabilityShowTypeRadioField, { stubs });

  test('Show type changed after trigger radio field', () => {
    const wrapper = factory({
      propsData: {
        value: AVAILABILITY_SHOW_TYPE.percent,
      },
    });

    selectRadioGroup(wrapper).triggerCustomEvent('input', AVAILABILITY_SHOW_TYPE.duration);

    expect(wrapper).toEmitInput(AVAILABILITY_SHOW_TYPE.duration);
  });
});

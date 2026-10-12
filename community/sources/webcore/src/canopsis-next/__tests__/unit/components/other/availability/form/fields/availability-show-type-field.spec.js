import { createCheckboxInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';

import { AVAILABILITY_SHOW_TYPE } from '@/constants';

import AvailabilityShowTypeField from '@/components/other/availability/form/fields/availability-show-type-field.vue';

const stubs = {
  'v-select': createCheckboxInputStub('v-select'),
};

const selectRadioGroup = wrapper => wrapper.find('.v-select');

describe('availability-show-type-field', () => {
  const factory = generateShallowRenderer(AvailabilityShowTypeField, { stubs });

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

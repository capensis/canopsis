import { createInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';

import { AVAILABILITY_DISPLAY_PARAMETERS } from '@/constants';

import AvailabilityDisplayParameterRadioField from '@/components/other/availability/form/fields/availability-display-parameter-radio-field.vue';

const stubs = {
  'v-radio-group': createInputStub('v-radio-group'),
};

const selectRadioGroup = wrapper => wrapper.find('.v-radio-group');

describe('availability-display-parameter-radio-field', () => {
  const factory = generateShallowRenderer(AvailabilityDisplayParameterRadioField, { stubs });

  test('Show type changed after trigger radio field', () => {
    const wrapper = factory({
      propsData: {
        value: AVAILABILITY_DISPLAY_PARAMETERS.uptime,
      },
    });

    selectRadioGroup(wrapper).triggerCustomEvent('input', AVAILABILITY_DISPLAY_PARAMETERS.downtime);

    expect(wrapper).toEmitInput(AVAILABILITY_DISPLAY_PARAMETERS.downtime);
  });
});

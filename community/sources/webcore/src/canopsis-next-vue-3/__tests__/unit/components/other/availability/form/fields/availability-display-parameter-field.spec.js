import { createInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';

import { AVAILABILITY_DISPLAY_PARAMETERS } from '@/constants';

import AvailabilityDisplayParameterField from '@/components/other/availability/form/fields/availability-display-parameter-field.vue';

const stubs = {
  'v-select': createInputStub('v-select'),
};

const selectSelectField = wrapper => wrapper.find('.v-select');

describe('availability-display-parameter-field', () => {
  const factory = generateShallowRenderer(AvailabilityDisplayParameterField, { stubs });

  test('Show type changed after trigger radio field', () => {
    const wrapper = factory({
      propsData: {
        value: AVAILABILITY_DISPLAY_PARAMETERS.uptime,
      },
    });

    selectSelectField(wrapper).triggerCustomEvent(
      'input',
      AVAILABILITY_DISPLAY_PARAMETERS.downtime,
    );

    expect(wrapper).toEmitInput(AVAILABILITY_DISPLAY_PARAMETERS.downtime);
  });
});

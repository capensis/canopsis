import { createSelectInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';

import { SERVICE_WEATHER_STATE_COUNTERS } from '@/constants';

import CServiceWeatherStateCounterField from '@/components/forms/fields/service-weather/c-service-weather-state-counter-field.vue';

const stubs = {
  'v-select': createSelectInputStub('v-select'),
};

const selectSelectField = wrapper => wrapper.find('c-select-stub');

describe('c-service-weather-state-counter-field', () => {
  const factory = generateShallowRenderer(CServiceWeatherStateCounterField, { stubs });

  it('Value changed after trigger the input', () => {
    const wrapper = factory();

    const selectField = selectSelectField(wrapper);

    selectField.triggerCustomEvent('input', SERVICE_WEATHER_STATE_COUNTERS.all);

    expect(wrapper).toEmitInput(SERVICE_WEATHER_STATE_COUNTERS.all);
  });

  it('Passes required state to the select', () => {
    const wrapper = factory({
      propsData: {
        required: true,
      },
    });

    expect(selectSelectField(wrapper).attributes('required')).toBeDefined();
  });
});

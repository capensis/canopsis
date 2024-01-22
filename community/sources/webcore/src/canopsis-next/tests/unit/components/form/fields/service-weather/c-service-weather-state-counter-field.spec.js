import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createSelectInputStub } from '@unit/stubs/input';

import { SERVICE_WEATHER_STATE_COUNTERS } from '@/constants';

import CServiceWeatherStateCounterField from '@/components/forms/fields/service-weather/c-service-weather-state-counter-field.vue';

const stubs = {
  'v-select': createSelectInputStub('v-select'),
};

const selectSelectField = wrapper => wrapper.find('.v-select');

describe('c-service-weather-state-counter-field', () => {
  const factory = generateShallowRenderer(CServiceWeatherStateCounterField, { stubs });
  const snapshotFactory = generateRenderer(CServiceWeatherStateCounterField);

  it('Value changed after trigger the input', () => {
    const wrapper = factory();

    const selectField = selectSelectField(wrapper);

    selectField.triggerCustomEvent('input', SERVICE_WEATHER_STATE_COUNTERS.all);

    expect(wrapper).toEmit('input', SERVICE_WEATHER_STATE_COUNTERS.all);
  });

  it('Renders `c-service-weather-state-counter-field` with default props', async () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
    await wrapper.activateAllMenus();
    expect(wrapper).toMatchMenuSnapshot();
  });

  it('Renders `c-service-weather-state-counter-field` with custom props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: [SERVICE_WEATHER_STATE_COUNTERS.acked, SERVICE_WEATHER_STATE_COUNTERS.all],
        name: 'customName',
        disabled: true,
        required: true,
        max: 2,
      },
    });

    expect(wrapper).toMatchSnapshot();
    await wrapper.activateAllMenus();
    expect(wrapper).toMatchMenuSnapshot();
  });
});

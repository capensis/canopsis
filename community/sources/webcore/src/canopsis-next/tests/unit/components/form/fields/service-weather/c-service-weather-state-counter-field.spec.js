import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';
import { createSelectInputStub } from '@unit/stubs/input';
import { SERVICE_WEATHER_STATE_COUNTERS } from '@/constants';

import CServiceWeatherStateCounterField from '@/components/forms/fields/service-weather/c-service-weather-state-counter-field.vue';

const localVue = createVueInstance();

const stubs = {
  'v-select': createSelectInputStub('v-select'),
};

const factory = (options = {}) => shallowMount(CServiceWeatherStateCounterField, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(CServiceWeatherStateCounterField, {
  localVue,

  ...options,
});

const selectSelectField = wrapper => wrapper.find('.v-select');

describe('c-service-weather-state-counter-field', () => {
  it('Value changed after trigger the input', () => {
    const wrapper = factory();

    const selectField = selectSelectField(wrapper);

    selectField.vm.$emit('input', SERVICE_WEATHER_STATE_COUNTERS.all);

    expect(wrapper).toEmit('input', SERVICE_WEATHER_STATE_COUNTERS.all);
  });

  it('Renders `c-service-weather-state-counter-field` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  it('Renders `c-service-weather-state-counter-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: [SERVICE_WEATHER_STATE_COUNTERS.acked, SERVICE_WEATHER_STATE_COUNTERS.all],
        name: 'customName',
        disabled: true,
        required: true,
        max: 2,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });
});

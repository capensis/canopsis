import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createMockedStoreModules, createPbehaviorTypesModule } from '@unit/utils/store';

import FieldCountersSelector from '@/components/sidebars/service-weather/form/fields/counters-selector.vue';

const stubs = {
  'widget-settings-item': true,
  'c-enabled-field': true,
  'c-pbehavior-type-field': true,
  'c-service-weather-state-counter-field': true,
};

const selectEnabledFieldByIndex = (wrapper, index) => wrapper
  .findAll('c-enabled-field-stub')
  .at(index);
const selectPbehaviorEnabledField = wrapper => selectEnabledFieldByIndex(wrapper, 0);
const selectPbehaviorTypeField = wrapper => wrapper.find('c-pbehavior-type-field-stub');
const selectStateEnabledField = wrapper => selectEnabledFieldByIndex(wrapper, 1);
const selectServiceWeatherStateCounterField = wrapper => wrapper.find('c-service-weather-state-counter-field-stub');

describe('field-counters-selector', () => {
  const { pbehaviorTypesModule } = createPbehaviorTypesModule();
  const store = createMockedStoreModules([
    pbehaviorTypesModule,
  ]);

  const value = {
    pbehavior_enabled: false,
    pbehavior_types: [],
    state_enabled: false,
    state_types: [],
  };
  const factory = generateShallowRenderer(FieldCountersSelector, {
    stubs,
    store,
  });
  const snapshotFactory = generateRenderer(FieldCountersSelector, {
    stubs,
    store,
  });

  test('Pbehavior counters enabled after trigger enabled field', () => {
    const wrapper = factory({
      propsData: {
        value,
      },
    });

    selectPbehaviorEnabledField(wrapper).triggerCustomEvent('input', true);

    expect(wrapper).toEmitInput({
      ...value,
      pbehavior_enabled: true,
    });
  });

  test('Pbehavior counters changed after trigger pbehavior type field', () => {
    const wrapper = factory({
      propsData: {
        value,
      },
    });

    const newPbehaviorTypes = [Faker.datatype.string()];

    selectPbehaviorTypeField(wrapper).triggerCustomEvent('input', newPbehaviorTypes);

    expect(wrapper).toEmitInput({
      ...value,
      pbehavior_types: newPbehaviorTypes,
    });
  });

  test('State counters enabled after trigger enabled field', () => {
    const wrapper = factory({
      propsData: {
        value,
      },
    });

    selectStateEnabledField(wrapper).triggerCustomEvent('input', true);

    expect(wrapper).toEmitInput({
      ...value,
      state_enabled: true,
    });
  });

  test('State counters changed after trigger state type field', () => {
    const wrapper = factory({
      propsData: {
        value,
      },
    });

    const newStateTypes = [Faker.datatype.string()];

    selectServiceWeatherStateCounterField(wrapper).triggerCustomEvent('input', newStateTypes);

    expect(wrapper).toEmitInput({
      ...value,
      state_types: newStateTypes,
    });
  });

  test('Renders `field-counters-selector` with disabled props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: {
          pbehavior_enabled: false,
          pbehavior_types: [],
          state_enabled: false,
          state_types: [],
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `field-counters-selector` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: {
          pbehavior_enabled: true,
          pbehavior_types: ['pbh-type'],
          state_enabled: true,
          state_types: ['state-type'],
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});

import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import { PATTERN_CUSTOM_ITEM_VALUE, PATTERN_TYPES } from '@/constants';

import PatternForm from '@/components/forms/pattern-form.vue';

const stubs = {
  'c-name-field': true,
  'c-alarm-patterns-field': true,
  'c-entity-patterns-field': true,
  'c-pbehavior-patterns-field': true,
  'c-service-weather-patterns-field': true,
  'pattern-progress': true,
};

const selectAlarmPatternsField = wrapper => wrapper.find('c-alarm-patterns-field-stub');
const selectEntityPatternsField = wrapper => wrapper.find('c-entity-patterns-field-stub');
const selectPbehaviorPatternsField = wrapper => wrapper.find('c-pbehavior-patterns-field-stub');
const selectServiceWeatherPatternsField = wrapper =>
  wrapper.find('c-service-weather-patterns-field-stub');

describe('pattern-form', () => {
  const factory = generateShallowRenderer(PatternForm, { stubs });

  test('Title changed after trigger text field', () => {
    const form = {
      title: '',
      groups: {},
      type: PATTERN_TYPES.alarm,
    };
    factory({
      propsData: {
        form,
      },
    });

    const title = Faker.datatype.string();

    globalThis.__TEST_FORM__.model('title').value = title;

    expect(globalThis.__TEST_FORM__.values.value).toEqual({ ...form, title });
  });

  test('Alarm pattern changed after trigger alarm patterns field', () => {
    const form = {
      title: '',
      id: PATTERN_CUSTOM_ITEM_VALUE,
      groups: [],
      type: PATTERN_TYPES.alarm,
    };
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const id = Faker.datatype.string();
    const groups = [{}];
    const nextForm = { ...form, id, groups };

    selectAlarmPatternsField(wrapper).triggerCustomEvent('update:patterns', nextForm);

    expect(globalThis.__TEST_FORM__.values.value).toEqual(nextForm);
  });

  test('Entity pattern changed after trigger entity patterns field', () => {
    const form = {
      title: '',
      id: PATTERN_CUSTOM_ITEM_VALUE,
      groups: [],
      type: PATTERN_TYPES.entity,
    };
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    expect(selectEntityPatternsField(wrapper).attributes('patterns')).toBe('[object Object]');

    const id = Faker.datatype.string();
    const groups = [{}];
    const nextForm = { ...form, id, groups };

    selectEntityPatternsField(wrapper).triggerCustomEvent('update:patterns', nextForm);

    expect(globalThis.__TEST_FORM__.values.value).toEqual(nextForm);
  });

  test('Pbehavior pattern changed after trigger pbehavior patterns field', () => {
    const form = {
      title: '',
      id: PATTERN_CUSTOM_ITEM_VALUE,
      groups: [],
      type: PATTERN_TYPES.pbehavior,
    };
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    expect(selectPbehaviorPatternsField(wrapper).attributes('patterns')).toBe('[object Object]');

    const id = Faker.datatype.string();
    const groups = [{}];
    const nextForm = { ...form, id, groups };

    selectPbehaviorPatternsField(wrapper).triggerCustomEvent('update:patterns', nextForm);

    expect(globalThis.__TEST_FORM__.values.value).toEqual(nextForm);
  });

  test('Service weather pattern changed after trigger service weather patterns field', () => {
    const form = {
      title: '',
      id: PATTERN_CUSTOM_ITEM_VALUE,
      groups: [],
      type: PATTERN_TYPES.serviceWeather,
    };
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const id = Faker.datatype.string();
    const groups = [{}];
    const nextForm = { ...form, id, groups };
    const serviceWeatherPatternsField = selectServiceWeatherPatternsField(wrapper);

    expect(serviceWeatherPatternsField.attributes('patterns')).toBe('[object Object]');

    serviceWeatherPatternsField.triggerCustomEvent('update:patterns', nextForm);

    expect(globalThis.__TEST_FORM__.values.value).toEqual(nextForm);
  });
});

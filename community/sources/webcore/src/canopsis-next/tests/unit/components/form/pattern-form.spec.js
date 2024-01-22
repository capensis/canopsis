import Faker from 'faker';

import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { createInputStub } from '@unit/stubs/input';
import { PATTERN_CUSTOM_ITEM_VALUE, PATTERN_TYPES } from '@/constants';

import PatternForm from '@/components/forms/pattern-form.vue';

const stubs = {
  'v-text-field': createInputStub('v-text-field'),
  'c-alarm-patterns-field': true,
  'c-entity-patterns-field': true,
  'c-pbehavior-patterns-field': true,
};

const selectTextField = wrapper => wrapper.find('.v-text-field');
const selectAlarmPatternsField = wrapper => wrapper.find('c-alarm-patterns-field-stub');
const selectEntityPatternsField = wrapper => wrapper.find('c-entity-patterns-field-stub');
const selectPbehaviorPatternsField = wrapper => wrapper.find('c-pbehavior-patterns-field-stub');

describe('pattern-form', () => {
  const factory = generateShallowRenderer(PatternForm, { stubs });
  const snapshotFactory = generateRenderer(PatternForm, { stubs });

  test('Title changed after trigger text field', () => {
    const form = {
      title: '',
      groups: {},
      type: PATTERN_TYPES.alarm,
    };
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const title = Faker.datatype.string();

    const textField = selectTextField(wrapper);

    textField.triggerCustomEvent('input', title);

    expect(wrapper).toEmit('input', { ...form, title });
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

    const alarmPatternsField = selectAlarmPatternsField(wrapper);

    alarmPatternsField.triggerCustomEvent('input', { ...form, id, groups });

    expect(wrapper).toEmit('input', { ...form, id, groups });
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

    const id = Faker.datatype.string();
    const groups = [{}];

    const entityPatternsField = selectEntityPatternsField(wrapper);

    entityPatternsField.triggerCustomEvent('input', { ...form, id, groups });

    expect(wrapper).toEmit('input', { ...form, id, groups });
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

    const id = Faker.datatype.string();
    const groups = [{}];

    const pbehaviorPatternsField = selectPbehaviorPatternsField(wrapper);

    pbehaviorPatternsField.triggerCustomEvent('input', { ...form, id, groups });

    expect(wrapper).toEmit('input', { ...form, id, groups });
  });

  test('Renders `pattern-form` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {},
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `pattern-form` with alarm type', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          title: 'Title',
          type: PATTERN_TYPES.alarm,
          alarm_pattern: {},
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `pattern-form` with entity type', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          title: 'Title',
          type: PATTERN_TYPES.entity,
          entity_pattern: {},
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `pattern-form` with entity type', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          title: 'Title',
          type: PATTERN_TYPES.pbehavior,
          pbehavior_pattern: {},
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `pattern-form` with errors', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          title: '',
        },
      },
    });

    const validator = wrapper.getValidator();

    await validator.validateAll();

    expect(wrapper).toMatchSnapshot();
  });
});

import Faker from 'faker';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import {
  ALARM_PATTERN_FIELDS,
  PATTERN_CONDITIONS,
  PATTERN_CUSTOM_ITEM_VALUE,
  PATTERN_OPERATORS,
  PATTERN_TYPES,
  QUICK_RANGES,
  TIME_UNITS,
} from '@/constants';

import CPatternsEditorField from '@/components/forms/fields/pattern/c-patterns-editor-field.vue';

const localVue = createVueInstance();

const stubs = {
  'c-pattern-field': true,
  'c-pattern-groups-field': true,
  'c-patterns-advanced-editor-field': true,
};

const factory = (options = {}) => shallowMount(CPatternsEditorField, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(CPatternsEditorField, {
  localVue,
  stubs,

  ...options,
});

const selectTabItems = wrapper => wrapper.findAll('a.v-tabs__item');
const selectAdvancedTabItems = wrapper => selectTabItems(wrapper).at(1);
const selectPatternField = wrapper => wrapper.find('c-pattern-field-stub');
const selectEditButton = wrapper => wrapper.find('v-btn-stub');

describe('c-patterns-editor-field', () => {
  test('Pattern id changed to custom after trigger input event on the pattern field', () => {
    const wrapper = factory({
      propsData: {
        patterns: {
          id: '',
          groups: [],
        },
        attributes: [],
        withType: true,
      },
    });

    const patternField = selectPatternField(wrapper);

    patternField.vm.$emit('input', { _id: PATTERN_CUSTOM_ITEM_VALUE });

    expect(wrapper).toEmit('input', {
      id: PATTERN_CUSTOM_ITEM_VALUE,
      groups: [],
    });
  });

  test('Pattern changed after trigger input event on the pattern field', () => {
    const wrapper = factory({
      propsData: {
        patterns: {
          id: '',
          groups: [],
        },
        attributes: [],
        withType: true,
      },
    });

    const patternField = selectPatternField(wrapper);

    const id = Faker.datatype.string();
    const alarmPattern = {
      field: ALARM_PATTERN_FIELDS.component,
      cond: {
        type: PATTERN_CONDITIONS.notEqual,
        value: 'component',
      },
    };
    const pattern = {
      _id: id,
      alarm_pattern: [
        [alarmPattern],
      ],
    };

    patternField.vm.$emit('input', pattern);

    expect(wrapper).toEmit('input', {
      id,
      groups: [{
        key: expect.any(String),
        rules: [
          {
            attribute: alarmPattern.field,
            duration: {
              unit: TIME_UNITS.second,
              value: 1,
            },
            dictionary: '',
            field: '',
            key: expect.any(String),
            operator: PATTERN_OPERATORS.notEqual,
            range: {
              from: 0,
              to: 0,
              type: QUICK_RANGES.last1Hour.value,
            },
            value: alarmPattern.cond.value,
          },
        ],
      }],
    });
  });

  test('Pattern changed to custom after click on the edit button', async () => {
    const rules = [
      {
        attribute: ALARM_PATTERN_FIELDS.component,
        duration: {
          unit: TIME_UNITS.second,
          value: 1,
        },
        dictionary: '',
        field: '',
        key: expect.any(String),
        operator: PATTERN_OPERATORS.notEqual,
        range: {
          from: 0,
          to: 0,
          type: QUICK_RANGES.last1Hour.value,
        },
        value: '',
      },
    ];
    const patterns = {
      id: Faker.datatype.string(),
      groups: [{
        key: 'key',
        rules,
      }],
    };
    const wrapper = factory({
      propsData: {
        patterns,
        attributes: [],
        withType: true,
      },
    });

    const editButton = selectEditButton(wrapper);

    await editButton.vm.$emit('click');

    expect(wrapper).toEmit('input', {
      ...patterns,
      id: PATTERN_CUSTOM_ITEM_VALUE,
    });
  });

  test('Renders `c-patterns-editor-field` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        patterns: {
          id: 'pattern-id',
          groups: [],
        },
        attributes: [],
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `c-patterns-editor-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        patterns: {
          id: 'pattern-id',
          groups: [{}],
        },
        attributes: [
          { value: 'attribute-1', text: 'Attribute 1' },
        ],
        disabled: true,
        name: 'name',
        required: true,
        type: PATTERN_TYPES.alarm,
        withType: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `c-patterns-editor-field` with advanced tab', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        patterns: {
          id: 'pattern-id',
          groups: [],
        },
        attributes: [],
      },
    });

    const advancedTab = selectAdvancedTabItems(wrapper);

    await advancedTab.trigger('click');

    expect(wrapper.element).toMatchSnapshot();
  });
});

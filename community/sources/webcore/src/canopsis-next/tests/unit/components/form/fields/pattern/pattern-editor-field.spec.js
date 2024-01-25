import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import {
  ALARM_PATTERN_FIELDS,
  PATTERN_CONDITIONS,
  PATTERN_CUSTOM_ITEM_VALUE,
  PATTERN_FIELD_TYPES,
  PATTERN_OPERATORS,
  PATTERN_TYPES,
  QUICK_RANGES,
  TIME_UNITS,
} from '@/constants';

import PatternsEditorField from '@/components/forms/fields/pattern/pattern-editor-field.vue';

const stubs = {
  'c-pattern-field': true,
  'pattern-groups-field': true,
  'pattern-advanced-editor-field': true,
};

const selectTabItems = wrapper => wrapper.findAll('.v-tab');
const selectAdvancedTab = wrapper => selectTabItems(wrapper).at(1);
const selectPatternField = wrapper => wrapper.find('c-pattern-field-stub');
const selectEditButton = wrapper => wrapper.find('v-btn-stub');
const selectPatternAdvancedEditorField = wrapper => wrapper.find('pattern-advanced-editor-field-stub');

describe('pattern-editor-field', () => {
  const factory = generateShallowRenderer(PatternsEditorField, {
    stubs,
    parentComponent: {
      $_veeValidate: {
        validator: 'new',
      },
    },
  });
  const snapshotFactory = generateRenderer(PatternsEditorField, {
    stubs,
    parentComponent: {
      $_veeValidate: {
        validator: 'new',
      },
    },
  });

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
            fieldType: PATTERN_FIELD_TYPES.string,
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

  test('Pattern changed to custom after click on the edit button', async () => {
    const patterns = {
      id: Faker.datatype.string(),
      groups: [],
    };
    const wrapper = factory({
      propsData: {
        patterns,
        attributes: [],
        withType: true,
      },
    });

    const advancedEditor = selectPatternAdvancedEditorField(wrapper);

    const patternRule = {
      field: ALARM_PATTERN_FIELDS.displayName,
      fieldType: PATTERN_FIELD_TYPES.string,
      cond: {
        type: PATTERN_CONDITIONS.equal,
        value: Faker.datatype.string(),
      },
    };

    advancedEditor.vm.$emit('input', [[
      patternRule,
    ]]);

    expect(wrapper).toEmit('input', {
      ...patterns,
      groups: [{
        key: expect.any(String),
        rules: [
          {
            key: expect.any(String),
            attribute: patternRule.field,
            duration: {
              unit: TIME_UNITS.second,
              value: 1,
            },
            dictionary: '',
            field: '',
            fieldType: PATTERN_FIELD_TYPES.string,
            operator: PATTERN_OPERATORS.equal,
            range: {
              from: 0,
              to: 0,
              type: QUICK_RANGES.last1Hour.value,
            },
            value: patternRule.cond.value,
          },
        ],
      }],
    });
  });

  test('Renders `pattern-editor-field` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        patterns: {
          id: 'pattern-id',
          groups: [],
        },
        attributes: [],
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `pattern-editor-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        patterns: {
          id: 'pattern-id',
          groups: [{
            rules: [],
          }],
        },
        attributes: [
          { value: 'attribute-1', text: 'Attribute 1' },
        ],
        disabled: true,
        name: 'name',
        required: true,
        type: PATTERN_TYPES.alarm,
        withType: true,
        readonly: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `pattern-editor-field` with advanced tab', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        patterns: {
          id: 'pattern-id',
          groups: [],
        },
        attributes: [],
      },
    });

    await selectAdvancedTab(wrapper).trigger('click');

    expect(wrapper).toMatchSnapshot();
  });
});

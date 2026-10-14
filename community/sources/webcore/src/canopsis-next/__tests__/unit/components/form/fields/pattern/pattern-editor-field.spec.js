import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import {
  ALARM_PATTERN_FIELDS,
  PATTERN_CONDITIONS,
  PATTERN_CUSTOM_ITEM_VALUE,
  PATTERN_FIELD_TYPES,
  PATTERN_OPERATORS,
  QUICK_RANGES,
  TIME_UNITS,
} from '@/constants';

import PatternsEditorField from '@/components/forms/fields/pattern/pattern-editor-field.vue';

const stubs = {
  'c-pattern-field': true,
  'pattern-groups-field': true,
  'pattern-advanced-editor-field': true,
  'pattern-count-message': true,
  'v-tabs': {
    template: '<div class="v-tabs"><slot /></div>',
  },
  'v-tab': {
    template: '<div class="v-tab"><slot /></div>',
  },
  'v-window': {
    template: '<div class="v-window"><slot /></div>',
  },
  'v-window-item': {
    template: '<div class="v-window-item"><slot /></div>',
  },
};

const selectPatternField = wrapper => wrapper.find('c-pattern-field-stub');
const selectEditButton = wrapper => wrapper.find('v-btn-stub');
const selectPatternAdvancedEditorField = wrapper =>
  wrapper.find('pattern-advanced-editor-field-stub');

describe('pattern-editor-field', () => {
  const factory = generateShallowRenderer(PatternsEditorField, { stubs });

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

    selectPatternField(wrapper).triggerCustomEvent('update:value', {
      _id: PATTERN_CUSTOM_ITEM_VALUE,
    });

    expect(wrapper).toEmitInput({
      id: PATTERN_CUSTOM_ITEM_VALUE,
      groups: [],
      is_corporate: undefined,
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
      alarm_pattern: [[alarmPattern]],
    };

    selectPatternField(wrapper).triggerCustomEvent('update:value', pattern);

    expect(wrapper).toEmitInput({
      id,
      is_corporate: undefined,
      groups: [
        {
          key: expect.any(String),
          rules: [
            {
              attribute: alarmPattern.field,
              alias: false,
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
                from: '',
                to: '',
                type: QUICK_RANGES.last1Hour.value,
              },
              value: alarmPattern.cond.value,
            },
          ],
        },
      ],
    });
  });

  test('Pattern changed to custom after click on the edit button', async () => {
    const patterns = {
      id: Faker.datatype.string(),
      groups: [
        {
          key: 'key',
          rules: [],
        },
      ],
    };
    const wrapper = factory({
      propsData: {
        patterns,
        attributes: [],
        withType: true,
      },
    });

    await selectEditButton(wrapper).triggerCustomEvent('click');

    expect(wrapper).toEmitInput({
      ...patterns,
      id: PATTERN_CUSTOM_ITEM_VALUE,
    });
  });

  test('Pattern groups changed after trigger advanced editor', () => {
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

    const patternRule = {
      field: ALARM_PATTERN_FIELDS.displayName,
      fieldType: PATTERN_FIELD_TYPES.string,
      cond: {
        type: PATTERN_CONDITIONS.equal,
        value: Faker.datatype.string(),
      },
    };

    selectPatternAdvancedEditorField(wrapper).triggerCustomEvent('update:value', [[patternRule]]);

    expect(wrapper).toEmitInput({
      ...patterns,
      groups: [
        {
          key: expect.any(String),
          rules: [
            {
              key: expect.any(String),
              attribute: patternRule.field,
              alias: false,
              duration: {
                unit: TIME_UNITS.second,
                value: 1,
              },
              dictionary: '',
              field: '',
              fieldType: PATTERN_FIELD_TYPES.string,
              operator: PATTERN_OPERATORS.equal,
              range: {
                from: '',
                to: '',
                type: QUICK_RANGES.last1Hour.value,
              },
              value: patternRule.cond.value,
            },
          ],
        },
      ],
    });
  });
});

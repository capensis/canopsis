import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import {
  ALARM_PATTERN_FIELDS,
  EVENT_FILTER_PATTERN_FIELDS,
  PATTERN_CONDITIONS,
  PATTERN_FIELD_TYPES,
  PATTERN_RULE_TYPES,
  TIME_UNITS,
} from '@/constants';

import PatternAdvancedEditorField from '@/components/forms/fields/pattern/pattern-advanced-editor-field.vue';

const stubs = {
  'c-json-field': true,
};

const selectJsonField = wrapper => wrapper.findComponent({ name: 'c-json-field' });
const invalidPatternsMessage = 'Patterns are invalid or there is a disabled pattern field';

describe('pattern-advanced-editor-field', () => {
  const factory = generateShallowRenderer(PatternAdvancedEditorField, { stubs });

  test('Patterns invalid with wrong structure', () => {
    const wrapper = factory({
      propsData: {
        value: [[]],
        attributes: [],
        name: 'patterns',
      },
    });

    const jsonField = selectJsonField(wrapper);

    jsonField.triggerCustomEvent('update:value', {});
    jsonField.triggerCustomEvent('update:value', [{}]);
    jsonField.triggerCustomEvent('update:value', [[{}]]);
    jsonField.triggerCustomEvent('update:value', [
      [
        {
          field: 'field',
          cond: {
            type: 'cond-type',
          },
        },
      ],
    ]);

    expect(wrapper).not.toHaveBeenEmit('update:value');
    expect(wrapper.vm.errorMessage).toBe(invalidPatternsMessage);
  });

  test('Patterns invalid with disabled attribute', () => {
    const attribute = Faker.datatype.string();

    const wrapper = factory({
      propsData: {
        value: [[]],
        attributes: [
          {
            value: attribute,
            options: { disabled: true },
          },
        ],
        name: 'patterns',
      },
    });

    selectJsonField(wrapper).triggerCustomEvent('update:value', [
      [
        {
          field: attribute,
          cond: {
            type: PATTERN_CONDITIONS.equal,
            value: ' ',
          },
        },
      ],
    ]);

    expect(wrapper).not.toHaveBeenEmit('update:value');
    expect(wrapper.vm.errorMessage).toBe(invalidPatternsMessage);
  });

  test('Patterns invalid with invalid infos attribute', () => {
    const wrapper = factory({
      propsData: {
        value: [[]],
        attributes: [
          {
            value: 'infos',
            options: { type: PATTERN_RULE_TYPES.infos },
          },
          {
            value: 'extra-infos',
            options: { type: PATTERN_RULE_TYPES.extraInfos },
          },
        ],
        name: 'patterns',
      },
    });

    selectJsonField(wrapper).triggerCustomEvent('update:value', [
      [
        {
          field: 'not-infos',
          cond: {
            type: PATTERN_CONDITIONS.equal,
            value: ' ',
          },
        },
      ],
    ]);

    expect(wrapper).not.toHaveBeenEmit('update:value');
    expect(wrapper.vm.errorMessage).toBe(invalidPatternsMessage);
  });

  test('Patterns invalid with not exist attribute', () => {
    const wrapper = factory({
      propsData: {
        value: [[]],
        attributes: [
          {
            value: 'infos',
          },
        ],
        name: 'patterns',
      },
    });

    selectJsonField(wrapper).triggerCustomEvent('update:value', [
      [
        {
          field: 'not-exist',
          cond: {
            type: PATTERN_CONDITIONS.equal,
            value: ' ',
          },
        },
      ],
    ]);

    expect(wrapper).not.toHaveBeenEmit('update:value');
    expect(wrapper.vm.errorMessage).toBe(invalidPatternsMessage);
  });

  test('Patterns valid with valid rules', () => {
    const wrapper = factory({
      propsData: {
        value: [[]],
        attributes: [
          {
            value: ALARM_PATTERN_FIELDS.infos,
            options: { type: PATTERN_RULE_TYPES.infos },
          },
          {
            value: EVENT_FILTER_PATTERN_FIELDS.extraInfos,
            options: { type: PATTERN_RULE_TYPES.extraInfos },
          },
          { value: ALARM_PATTERN_FIELDS.displayName },
          { value: ALARM_PATTERN_FIELDS.status },
          { value: ALARM_PATTERN_FIELDS.creationDate },
          { value: ALARM_PATTERN_FIELDS.duration },
        ],
        name: 'patterns',
      },
    });

    const patterns = [
      [
        {
          field: `${ALARM_PATTERN_FIELDS.infos}.string`,
          cond: {
            type: PATTERN_CONDITIONS.equal,
            value: 'string',
          },
          field_type: PATTERN_FIELD_TYPES.string,
        },
        {
          field: `${EVENT_FILTER_PATTERN_FIELDS.extraInfos}.stringArray`,
          cond: {
            type: PATTERN_CONDITIONS.hasEvery,
            value: ['string'],
          },
          field_type: PATTERN_FIELD_TYPES.stringArray,
        },
        {
          field: `${ALARM_PATTERN_FIELDS.infos}.number`,
          cond: {
            type: PATTERN_CONDITIONS.equal,
            value: 12,
          },
          field_type: PATTERN_FIELD_TYPES.number,
        },
        {
          field: `${ALARM_PATTERN_FIELDS.infos}.boolean`,
          cond: {
            type: PATTERN_CONDITIONS.equal,
            value: true,
          },
          field_type: PATTERN_FIELD_TYPES.boolean,
        },
        {
          field: ALARM_PATTERN_FIELDS.displayName,
          cond: {
            type: PATTERN_CONDITIONS.isEmpty,
            value: true,
          },
        },
        {
          field: ALARM_PATTERN_FIELDS.status,
          cond: {
            type: PATTERN_CONDITIONS.isEmpty,
            value: 2,
          },
        },
        {
          field: ALARM_PATTERN_FIELDS.creationDate,
          cond: {
            type: PATTERN_CONDITIONS.relativeTime,
            value: {
              value: 200,
              unit: TIME_UNITS.second,
            },
          },
        },
        {
          field: ALARM_PATTERN_FIELDS.creationDate,
          cond: {
            type: PATTERN_CONDITIONS.absoluteTime,
            value: {
              from: 0,
              to: 1,
            },
          },
        },
        {
          field: ALARM_PATTERN_FIELDS.duration,
          cond: {
            type: PATTERN_CONDITIONS.greater,
            value: {
              unit: TIME_UNITS.hour,
              value: 1,
            },
          },
        },
      ],
    ];

    selectJsonField(wrapper).triggerCustomEvent('update:value', patterns);

    expect(wrapper).toEmitInput(patterns);
    expect(wrapper.vm.errorMessage).toBeFalsy();
  });
});

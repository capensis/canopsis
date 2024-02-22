import Faker from 'faker';
import { Validator } from 'vee-validate';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

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

const selectJsonFieldNode = wrapper => wrapper.vm.$children[0];

describe('pattern-advanced-editor-field', () => {
  const factory = generateShallowRenderer(PatternAdvancedEditorField, { stubs });
  const snapshotFactory = generateRenderer(PatternAdvancedEditorField, { stubs });

  test('Patterns invalid with wrong structure', () => {
    jest.useFakeTimers('legacy');

    const validator = new Validator();

    const wrapper = factory({
      propsData: {
        value: [[]],
        attributes: [],
        name: 'patterns',
      },
      provide: {
        $validator: validator,
      },
    });

    const jsonFieldNode = selectJsonFieldNode(wrapper);

    jsonFieldNode.$emit('input', {});
    jsonFieldNode.$emit('input', [{}]);
    jsonFieldNode.$emit('input', [[{}]]);
    jsonFieldNode.$emit('input', [[{
      field: 'field',
      cond: {
        type: 'cond-type',
      },
    }]]);

    expect(wrapper).not.toHaveBeenEmit('input');

    jest.runAllTimers();

    const errors = validator.errors.collect('patterns');

    expect(errors).toHaveLength(4);
    expect(errors[0]).toEqual(
      'Patterns are invalid or there is a disabled pattern field',
    );

    jest.useRealTimers();
  });

  test('Patterns invalid with disabled attribute', () => {
    jest.useFakeTimers('legacy');

    const validator = new Validator();
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
      provide: {
        $validator: validator,
      },
    });

    const jsonFieldNode = selectJsonFieldNode(wrapper);

    jsonFieldNode.$emit('input', [[{
      field: attribute,
      cond: {
        type: PATTERN_CONDITIONS.equal,
        value: ' ',
      },
    }]]);

    expect(wrapper).not.toHaveBeenEmit('input');

    jest.runAllTimers();

    const errors = validator.errors.collect('patterns');

    expect(errors).toHaveLength(1);

    jest.useRealTimers();
  });

  test('Patterns invalid with invalid infos attribute', () => {
    jest.useFakeTimers('legacy');

    const validator = new Validator();

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
      provide: {
        $validator: validator,
      },
    });

    const jsonFieldNode = selectJsonFieldNode(wrapper);

    jsonFieldNode.$emit('input', [[{
      field: 'not-infos',
      cond: {
        type: PATTERN_CONDITIONS.equal,
        value: ' ',
      },
    }]]);

    expect(wrapper).not.toHaveBeenEmit('input');

    jest.runAllTimers();

    const errors = validator.errors.collect('patterns');

    expect(errors).toHaveLength(1);

    jest.useRealTimers();
  });

  test('Patterns invalid with not exist attribute', () => {
    jest.useFakeTimers('legacy');

    const validator = new Validator();

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
      provide: {
        $validator: validator,
      },
    });

    const jsonFieldNode = selectJsonFieldNode(wrapper);

    jsonFieldNode.$emit('input', [[{
      field: 'not-exist',
      cond: {
        type: PATTERN_CONDITIONS.equal,
        value: ' ',
      },
    }]]);

    expect(wrapper).not.toHaveBeenEmit('input');

    jest.runAllTimers();

    const errors = validator.errors.collect('patterns');

    expect(errors).toHaveLength(1);

    jest.useRealTimers();
  });

  test('Patterns valid with valid rules', () => {
    jest.useFakeTimers('legacy');

    const validator = new Validator();

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
      provide: {
        $validator: validator,
      },
    });

    const jsonFieldNode = selectJsonFieldNode(wrapper);

    const patterns = [[
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
    ]];

    jsonFieldNode.$emit('input', patterns);

    expect(wrapper).toEmitInput(patterns);
  });

  test('Renders `pattern-advanced-editor-field` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: [[]],
        attributes: [],
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `pattern-advanced-editor-field` with custom props', () => {
    const name = Faker.datatype.string();
    const wrapper = snapshotFactory({
      propsData: {
        value: [
          [{
            field: name,
            cond: { type: PATTERN_CONDITIONS.equal, value: 'value' },
          }],
        ],
        attributes: [{
          value: name,
        }],
        disabled: true,
        name: 'customName',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});

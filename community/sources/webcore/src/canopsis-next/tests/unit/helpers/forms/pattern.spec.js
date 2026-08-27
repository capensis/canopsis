import Faker from 'faker';

import {
  ALARM_PATTERN_FIELDS,
  ENTITY_PATTERN_FIELDS,
  EVENT_FILTER_PATTERN_FIELDS,
  PATTERN_CONDITIONS,
  PATTERN_FIELD_TYPES,
  PATTERN_OPERATORS,
  PATTERN_RULE_INFOS_FIELDS,
  PATTERN_RULE_TYPES,
  PBEHAVIOR_FIELDS,
  PBEHAVIOR_PATTERN_FIELDS,
  QUICK_RANGES,
  TIME_UNITS,
} from '@/constants';

import { formRuleToPatternRule, getOperatorsByRule, patternRuleToForm } from '@/helpers/entities/pattern/form';
import { durationToForm } from '@/helpers/date/duration';

describe('pattern form converters', () => {
  const defaultForm = {
    key: expect.any(String),
    attribute: '',
    operator: '',
    field: '',
    fieldType: PATTERN_FIELD_TYPES.string,
    dictionary: '',
    value: '',
    alias: false,
    range: {
      type: QUICK_RANGES.last1Hour.value,
      from: 0,
      to: 0,
    },
    duration: durationToForm(),
  };

  const expectOperatorsToIncludeOneOf = operators => expect(operators).toEqual(expect.arrayContaining([
    PATTERN_OPERATORS.isOneOf,
    PATTERN_OPERATORS.isNotOneOf,
  ]));

  const expectOperatorsNotToIncludeOneOf = (operators) => {
    expect(operators).not.toContain(PATTERN_OPERATORS.isOneOf);
    expect(operators).not.toContain(PATTERN_OPERATORS.isNotOneOf);
  };

  const stringPatternFieldsWithOneOfOperators = [
    ALARM_PATTERN_FIELDS.displayName,
    ALARM_PATTERN_FIELDS.output,
    ALARM_PATTERN_FIELDS.longOutput,
    ALARM_PATTERN_FIELDS.initialOutput,
    ALARM_PATTERN_FIELDS.initialLongOutput,
    ALARM_PATTERN_FIELDS.component,
    ALARM_PATTERN_FIELDS.connector,
    ALARM_PATTERN_FIELDS.connectorName,
    ALARM_PATTERN_FIELDS.resource,
    ALARM_PATTERN_FIELDS.tags,
    ALARM_PATTERN_FIELDS.lastComment,
    ALARM_PATTERN_FIELDS.lastCommentInitiator,
    ALARM_PATTERN_FIELDS.lastCommentAuthor,
    ALARM_PATTERN_FIELDS.ticketMessage,
    ALARM_PATTERN_FIELDS.ticketValue,
    ALARM_PATTERN_FIELDS.ticketInitiator,
    ALARM_PATTERN_FIELDS.snoozeAuthor,
    ALARM_PATTERN_FIELDS.snoozeInitiator,
    ALARM_PATTERN_FIELDS.ackBy,
    ALARM_PATTERN_FIELDS.ackMessage,
    ALARM_PATTERN_FIELDS.ackInitiator,
    ALARM_PATTERN_FIELDS.canceledInitiator,
    ALARM_PATTERN_FIELDS.stateInitiator,
    ENTITY_PATTERN_FIELDS.id,
    ENTITY_PATTERN_FIELDS.name,
    ENTITY_PATTERN_FIELDS.category,
    ENTITY_PATTERN_FIELDS.type,
    ENTITY_PATTERN_FIELDS.connector,
    ENTITY_PATTERN_FIELDS.component,
    EVENT_FILTER_PATTERN_FIELDS.component,
    EVENT_FILTER_PATTERN_FIELDS.connector,
    EVENT_FILTER_PATTERN_FIELDS.connectorName,
    EVENT_FILTER_PATTERN_FIELDS.resource,
    EVENT_FILTER_PATTERN_FIELDS.output,
    EVENT_FILTER_PATTERN_FIELDS.longOutput,
    EVENT_FILTER_PATTERN_FIELDS.eventType,
    EVENT_FILTER_PATTERN_FIELDS.sourceType,
    EVENT_FILTER_PATTERN_FIELDS.initiator,
    EVENT_FILTER_PATTERN_FIELDS.author,
    PBEHAVIOR_PATTERN_FIELDS.name,
    PBEHAVIOR_FIELDS.name,
    PBEHAVIOR_FIELDS.author,
    PBEHAVIOR_FIELDS.rrule,
    PBEHAVIOR_FIELDS.reason,
    PBEHAVIOR_FIELDS.type,
    PBEHAVIOR_FIELDS.canonicalType,
  ];

  const oneOfOperatorCases = [
    [PATTERN_CONDITIONS.isOneOf, PATTERN_OPERATORS.isOneOf],
    [PATTERN_CONDITIONS.isNotOneOf, PATTERN_OPERATORS.isNotOneOf],
  ];

  const stringPatternFieldOneOfCases = stringPatternFieldsWithOneOfOperators.reduce(
    (acc, field) => [
      ...acc,
      ...oneOfOperatorCases.map(([condition, operator]) => [field, condition, operator]),
    ],
    [],
  );

  const expectFormValueToContainPrimitiveValues = (formValue, value) => {
    const primitiveFormValue = formValue.map(item => item?.value ?? item);

    expect(primitiveFormValue).toEqual(value);
  };

  it('should be converted to form and back to pattern with `equal` operator', () => {
    const value = Faker.lorem.word();

    const patternRule = {
      field: ALARM_PATTERN_FIELDS.displayName,
      cond: { type: PATTERN_CONDITIONS.equal, value },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.displayName,
      value,
      operator: PATTERN_OPERATORS.equal,
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with `not equal` operator', () => {
    const value = Faker.lorem.word();

    const patternRule = {
      field: ALARM_PATTERN_FIELDS.connectorName,
      cond: { type: PATTERN_CONDITIONS.notEqual, value },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.connectorName,
      value,
      operator: PATTERN_OPERATORS.notEqual,
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with `contains` operator', () => {
    const value = Faker.lorem.word();
    const patternRule = {
      field: ALARM_PATTERN_FIELDS.component,
      cond: { type: PATTERN_CONDITIONS.contains, value },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.component,
      operator: PATTERN_OPERATORS.contains,
      value,
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with `not contains` operator', () => {
    const value = Faker.lorem.word();
    const patternRule = {
      field: ALARM_PATTERN_FIELDS.component,
      cond: { type: PATTERN_CONDITIONS.notContains, value },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.component,
      operator: PATTERN_OPERATORS.notContains,
      value,
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with `regexp` operator', () => {
    const value = `^((?!${Faker.lorem.word()}).)*$`;
    const patternRule = {
      field: ALARM_PATTERN_FIELDS.output,
      cond: { type: PATTERN_CONDITIONS.regexp, value },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.output,
      operator: PATTERN_OPERATORS.regexp,
      value,
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with `acked` operator', () => {
    const patternRule = {
      field: ALARM_PATTERN_FIELDS.ack,
      cond: { type: PATTERN_CONDITIONS.exist, value: true },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.ack,
      operator: PATTERN_OPERATORS.acked,
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with `not acked` operator', () => {
    const patternRule = {
      field: ALARM_PATTERN_FIELDS.ack,
      cond: { type: PATTERN_CONDITIONS.exist, value: false },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.ack,
      operator: PATTERN_OPERATORS.notAcked,
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with `snoozed` operator', () => {
    const patternRule = {
      field: ALARM_PATTERN_FIELDS.snooze,
      cond: { type: PATTERN_CONDITIONS.exist, value: true },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.snooze,
      operator: PATTERN_OPERATORS.snoozed,
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with `not snoozed` operator', () => {
    const patternRule = {
      field: ALARM_PATTERN_FIELDS.snooze,
      cond: { type: PATTERN_CONDITIONS.exist, value: false },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.snooze,
      operator: PATTERN_OPERATORS.notSnoozed,
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with `canceled` operator', () => {
    const patternRule = {
      field: ALARM_PATTERN_FIELDS.canceled,
      cond: { type: PATTERN_CONDITIONS.exist, value: true },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.canceled,
      operator: PATTERN_OPERATORS.canceled,
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with `not canceled` operator', () => {
    const patternRule = {
      field: ALARM_PATTERN_FIELDS.canceled,
      cond: { type: PATTERN_CONDITIONS.exist, value: false },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.canceled,
      operator: PATTERN_OPERATORS.notCanceled,
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with `ticket associated` operator', () => {
    const patternRule = {
      field: ALARM_PATTERN_FIELDS.ticket,
      cond: { type: PATTERN_CONDITIONS.exist, value: true },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.ticket,
      operator: PATTERN_OPERATORS.ticketAssociated,
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with `ticket not associated` operator', () => {
    const patternRule = {
      field: ALARM_PATTERN_FIELDS.ticket,
      cond: { type: PATTERN_CONDITIONS.exist, value: false },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.ticket,
      operator: PATTERN_OPERATORS.ticketNotAssociated,
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with `begins with` operator', () => {
    const value = Faker.lorem.word();
    const patternRule = {
      field: ALARM_PATTERN_FIELDS.connector,
      cond: { type: PATTERN_CONDITIONS.beginWith, value },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.connector,
      operator: PATTERN_OPERATORS.beginWith,
      value,
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with `not begin with` operator', () => {
    const value = Faker.lorem.word();
    const patternRule = {
      field: ALARM_PATTERN_FIELDS.resource,
      cond: { type: PATTERN_CONDITIONS.notBeginWith, value },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.resource,
      operator: PATTERN_OPERATORS.notBeginWith,
      value,
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with `ends with` operator', () => {
    const value = Faker.lorem.word();
    const patternRule = {
      field: ALARM_PATTERN_FIELDS.connectorName,
      cond: { type: PATTERN_CONDITIONS.endsWith, value },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.connectorName,
      operator: PATTERN_OPERATORS.endsWith,
      value,
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with `not end with` operator', () => {
    const value = Faker.lorem.word();
    const patternRule = {
      field: ALARM_PATTERN_FIELDS.connectorName,
      cond: { type: PATTERN_CONDITIONS.notEndWith, value },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.connectorName,
      operator: PATTERN_OPERATORS.notEndWith,
      value,
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with alarm infos field and name `exist` operator', () => {
    const dictionary = Faker.lorem.word();

    const patternRule = {
      field: `${ALARM_PATTERN_FIELDS.infos}.${dictionary}`,
      cond: { type: PATTERN_CONDITIONS.exist, value: true },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.infos,
      field: PATTERN_RULE_INFOS_FIELDS.name,
      dictionary,
      operator: PATTERN_OPERATORS.exist,
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with alarm infos field and name `not exist` operator', () => {
    const dictionary = Faker.lorem.word();

    const patternRule = {
      field: `${ALARM_PATTERN_FIELDS.infos}.${dictionary}`,
      cond: { type: PATTERN_CONDITIONS.exist, value: false },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.infos,
      field: PATTERN_RULE_INFOS_FIELDS.name,
      dictionary,
      operator: PATTERN_OPERATORS.notExist,
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with entity infos field and name `exist` operator', () => {
    const dictionary = Faker.lorem.word();

    const patternRule = {
      field: `${ENTITY_PATTERN_FIELDS.infos}.${dictionary}`,
      cond: { type: PATTERN_CONDITIONS.exist, value: true },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ENTITY_PATTERN_FIELDS.infos,
      field: PATTERN_RULE_INFOS_FIELDS.name,
      dictionary,
      operator: PATTERN_OPERATORS.exist,
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with entity infos field and name `not exist` operator', () => {
    const dictionary = Faker.lorem.word();

    const patternRule = {
      field: `${ENTITY_PATTERN_FIELDS.infos}.${dictionary}`,
      cond: { type: PATTERN_CONDITIONS.exist, value: false },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ENTITY_PATTERN_FIELDS.infos,
      field: PATTERN_RULE_INFOS_FIELDS.name,
      dictionary,
      operator: PATTERN_OPERATORS.notExist,
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with `has every` operator', () => {
    const value = [Faker.lorem.word()];
    const patternRule = {
      field: ALARM_PATTERN_FIELDS.connectorName,
      cond: { type: PATTERN_CONDITIONS.hasEvery, value },
    };

    const form = patternRuleToForm(patternRule);

    const { value: defaultValue, ...defaultFormWithoutValue } = defaultForm;
    expect(form).toMatchObject({
      ...defaultFormWithoutValue,
      attribute: ALARM_PATTERN_FIELDS.connectorName,
      operator: PATTERN_OPERATORS.hasEvery,
    });
    expect(form.value).toEqual(value);
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with `has one of` operator', () => {
    const value = [Faker.lorem.word()];
    const patternRule = {
      field: ALARM_PATTERN_FIELDS.connectorName,
      cond: { type: PATTERN_CONDITIONS.hasOneOf, value },
    };

    const form = patternRuleToForm(patternRule);

    const { value: defaultValue, ...defaultFormWithoutValue } = defaultForm;
    expect(form).toMatchObject({
      ...defaultFormWithoutValue,
      attribute: ALARM_PATTERN_FIELDS.connectorName,
      operator: PATTERN_OPERATORS.hasOneOf,
    });
    expect(form.value).toEqual(value);
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it.each(stringPatternFieldOneOfCases)(
    'should keep correct value after converting `%s` field with `%s` condition to form and back',
    (field, condition, operator) => {
      const value = [Faker.lorem.word(), Faker.lorem.word()];
      const patternRule = {
        field,
        cond: { type: condition, value },
      };

      const form = patternRuleToForm(patternRule);

      const { value: defaultValue, ...defaultFormWithoutValue } = defaultForm;
      expect(form).toMatchObject({
        ...defaultFormWithoutValue,
        attribute: field,
        operator,
      });
      expect(form.value).not.toEqual([]);
      expectFormValueToContainPrimitiveValues(form.value, value);
      expect(formRuleToPatternRule(form)).toEqual(patternRule);
    },
  );

  it.each(oneOfOperatorCases)('should keep correct value after converting `%s` condition for an infos value field to form and back', (condition, operator) => {
    const value = [Faker.lorem.word(), Faker.lorem.word()];
    const dictionary = Faker.lorem.word();
    const patternRule = {
      field: `${ENTITY_PATTERN_FIELDS.infos}.${dictionary}`,
      field_type: PATTERN_FIELD_TYPES.string,
      cond: { type: condition, value },
    };

    const form = patternRuleToForm(patternRule);

    const { value: defaultValue, ...defaultFormWithoutValue } = defaultForm;
    expect(form).toMatchObject({
      ...defaultFormWithoutValue,
      attribute: ENTITY_PATTERN_FIELDS.infos,
      field: PATTERN_RULE_INFOS_FIELDS.value,
      dictionary,
      operator,
    });
    expect(form.value).toEqual(value.map(item => ({ key: expect.any(String), value: item })));
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it.each(oneOfOperatorCases)('should keep correct value after converting `%s` condition for an extra infos value field to form and back', (condition, operator) => {
    const value = [Faker.lorem.word(), Faker.lorem.word()];
    const dictionary = Faker.lorem.word();
    const patternRule = {
      field: `${EVENT_FILTER_PATTERN_FIELDS.extraInfos}.${dictionary}`,
      field_type: PATTERN_FIELD_TYPES.string,
      cond: { type: condition, value },
    };

    const form = patternRuleToForm(patternRule);

    const { value: defaultValue, ...defaultFormWithoutValue } = defaultForm;
    expect(form).toMatchObject({
      ...defaultFormWithoutValue,
      attribute: EVENT_FILTER_PATTERN_FIELDS.extraInfos,
      field: PATTERN_RULE_INFOS_FIELDS.value,
      dictionary,
      operator,
    });
    expect(form.value).toEqual(value.map(item => ({ key: expect.any(String), value: item })));
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it.each(oneOfOperatorCases)('should keep correct value after converting `%s` condition for an object field to form and back', (condition, operator) => {
    const value = [Faker.lorem.word(), Faker.lorem.word()];
    const dictionary = Faker.lorem.word();
    const patternRule = {
      field: `${ALARM_PATTERN_FIELDS.ticketData}.${dictionary}`,
      cond: { type: condition, value },
    };

    const form = patternRuleToForm(patternRule);

    const { value: defaultValue, ...defaultFormWithoutValue } = defaultForm;
    expect(form).toMatchObject({
      ...defaultFormWithoutValue,
      attribute: ALARM_PATTERN_FIELDS.ticketData,
      dictionary,
      operator,
    });
    expect(form.value).toEqual(value.map(item => ({ key: expect.any(String), value: item })));
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it.each(stringPatternFieldsWithOneOfOperators)(
    'should include one-of operators for string pattern field `%s`',
    (field) => {
      const operators = getOperatorsByRule({
        attribute: field,
        fieldType: PATTERN_FIELD_TYPES.string,
      }, PATTERN_RULE_TYPES.string);

      expectOperatorsToIncludeOneOf(operators);
    },
  );

  it.each([
    PATTERN_RULE_TYPES.infos,
    PATTERN_RULE_TYPES.extraInfos,
    PATTERN_RULE_TYPES.object,
  ])('should include one-of operators for %s value fields', (ruleType) => {
    const operators = getOperatorsByRule({
      field: PATTERN_RULE_INFOS_FIELDS.value,
      fieldType: PATTERN_FIELD_TYPES.string,
    }, ruleType);

    expectOperatorsToIncludeOneOf(operators);
  });

  it.each([
    [
      'infos name field',
      { field: PATTERN_RULE_INFOS_FIELDS.name, fieldType: PATTERN_FIELD_TYPES.string },
      PATTERN_RULE_TYPES.infos,
    ],
    [
      'number field',
      { attribute: ENTITY_PATTERN_FIELDS.impactLevel, fieldType: PATTERN_FIELD_TYPES.number },
      PATTERN_RULE_TYPES.number,
    ],
    [
      'date field',
      { attribute: ALARM_PATTERN_FIELDS.creationDate, fieldType: PATTERN_FIELD_TYPES.string },
      PATTERN_RULE_TYPES.date,
    ],
    [
      'duration field',
      { attribute: ALARM_PATTERN_FIELDS.duration, fieldType: PATTERN_FIELD_TYPES.string },
      PATTERN_RULE_TYPES.duration,
    ],
    [
      'string array field type',
      { attribute: ALARM_PATTERN_FIELDS.tags, fieldType: PATTERN_FIELD_TYPES.stringArray },
      PATTERN_RULE_TYPES.string,
    ],
  ])('should not include one-of operators for %s', (caseName, rule, ruleType) => {
    const operators = getOperatorsByRule({
      ...rule,
    }, ruleType);

    expectOperatorsNotToIncludeOneOf(operators);
  });

  it('should be converted to form and back to pattern with `has not` operator', () => {
    const value = [Faker.lorem.word()];
    const patternRule = {
      field: ALARM_PATTERN_FIELDS.connectorName,
      cond: { type: PATTERN_CONDITIONS.hasNot, value },
    };

    const form = patternRuleToForm(patternRule);

    const { value: defaultValue, ...defaultFormWithoutValue } = defaultForm;
    expect(form).toMatchObject({
      ...defaultFormWithoutValue,
      attribute: ALARM_PATTERN_FIELDS.connectorName,
      operator: PATTERN_OPERATORS.hasNot,
    });
    expect(form.value).toEqual(value);
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with `is empty` operator', () => {
    const patternRule = {
      field: ALARM_PATTERN_FIELDS.connectorName,
      cond: { type: PATTERN_CONDITIONS.isEmpty, value: true },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.connectorName,
      operator: PATTERN_OPERATORS.isEmpty,
      value: [],
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with `is not empty` operator', () => {
    const patternRule = {
      field: ALARM_PATTERN_FIELDS.connectorName,
      cond: { type: PATTERN_CONDITIONS.isEmpty, value: false },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.connectorName,
      operator: PATTERN_OPERATORS.isNotEmpty,
      value: [],
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with `higher than` operator', () => {
    const value = Faker.datatype.number();
    const patternRule = {
      field: ENTITY_PATTERN_FIELDS.impactLevel,
      cond: { type: PATTERN_CONDITIONS.greater, value },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ENTITY_PATTERN_FIELDS.impactLevel,
      operator: PATTERN_OPERATORS.higher,
      value,
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with `lower than` operator', () => {
    const value = Faker.datatype.number();
    const patternRule = {
      field: ENTITY_PATTERN_FIELDS.impactLevel,
      cond: { type: PATTERN_CONDITIONS.less, value },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ENTITY_PATTERN_FIELDS.impactLevel,
      operator: PATTERN_OPERATORS.lower,
      value,
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with `longer` operator', () => {
    const value = {
      value: Faker.datatype.number(),
      unit: TIME_UNITS.second,
    };
    const patternRule = {
      field: ALARM_PATTERN_FIELDS.duration,
      cond: { type: PATTERN_CONDITIONS.greater, value },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.duration,
      operator: PATTERN_OPERATORS.longer,
      duration: value,
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with `shorter` operator', () => {
    const value = {
      value: Faker.datatype.number(),
      unit: TIME_UNITS.second,
    };
    const patternRule = {
      field: ALARM_PATTERN_FIELDS.duration,
      cond: { type: PATTERN_CONDITIONS.less, value },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.duration,
      operator: PATTERN_OPERATORS.shorter,
      duration: value,
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with `relative time` condition with `to` value', () => {
    const patternRule = {
      field: ALARM_PATTERN_FIELDS.creationDate,
      cond: {
        type: PATTERN_CONDITIONS.relativeTime,
        value: {
          to: {
            value: 1,
            unit: TIME_UNITS.hour,
          },
        },
      },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.creationDate,
      operator: PATTERN_OPERATORS.olderThan,
      range: {
        type: QUICK_RANGES.last1Hour.value,
        from: 0,
        to: 0,
      },
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with `absolute time` condition', () => {
    const value = {
      from: Faker.datatype.number(),
      to: Faker.datatype.number(),
    };
    const patternRule = {
      field: ALARM_PATTERN_FIELDS.creationDate,
      cond: { type: PATTERN_CONDITIONS.absoluteTime, value },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.creationDate,
      operator: PATTERN_OPERATORS.inRangeDates,
      range: {
        type: QUICK_RANGES.last1Hour.value,
        from: new Date(value.from * 1000),
        to: new Date(value.to * 1000),
      },
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with infos and number value', () => {
    const value = Faker.datatype.number();
    const dictionary = Faker.lorem.word();
    const patternRule = {
      field: `${ALARM_PATTERN_FIELDS.infos}.${dictionary}`,
      field_type: PATTERN_FIELD_TYPES.number,
      cond: { type: PATTERN_CONDITIONS.equal, value },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.infos,
      field: PATTERN_RULE_INFOS_FIELDS.value,
      fieldType: PATTERN_FIELD_TYPES.number,
      operator: PATTERN_OPERATORS.equal,
      dictionary,
      value,
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with infos and string value', () => {
    const value = Faker.lorem.word();
    const dictionary = Faker.lorem.word();
    const patternRule = {
      field: `${ALARM_PATTERN_FIELDS.infos}.${dictionary}`,
      field_type: PATTERN_FIELD_TYPES.string,
      cond: { type: PATTERN_CONDITIONS.equal, value },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.infos,
      field: PATTERN_RULE_INFOS_FIELDS.value,
      operator: PATTERN_OPERATORS.equal,
      dictionary,
      value,
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with infos and boolean value', () => {
    const value = Faker.datatype.boolean();
    const dictionary = Faker.lorem.word();
    const patternRule = {
      field: `${ALARM_PATTERN_FIELDS.infos}.${dictionary}`,
      field_type: PATTERN_FIELD_TYPES.boolean,
      cond: { type: PATTERN_CONDITIONS.equal, value },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.infos,
      field: PATTERN_RULE_INFOS_FIELDS.value,
      fieldType: PATTERN_FIELD_TYPES.boolean,
      operator: PATTERN_OPERATORS.equal,
      dictionary,
      value,
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with infos and string array value', () => {
    const value = Faker.datatype.array(2);
    const dictionary = Faker.lorem.word();
    const patternRule = {
      field: `${ALARM_PATTERN_FIELDS.infos}.${dictionary}`,
      field_type: PATTERN_FIELD_TYPES.stringArray,
      cond: { type: PATTERN_CONDITIONS.hasNot, value },
    };

    const form = patternRuleToForm(patternRule);

    const { value: defaultValue, ...defaultFormWithoutValue } = defaultForm;
    expect(form).toMatchObject({
      ...defaultFormWithoutValue,
      attribute: ALARM_PATTERN_FIELDS.infos,
      field: PATTERN_RULE_INFOS_FIELDS.value,
      fieldType: PATTERN_FIELD_TYPES.stringArray,
      operator: PATTERN_OPERATORS.hasNot,
      dictionary,
    });
    expect(form.value).toEqual(value.map(item => ({ key: expect.any(String), value: item })));
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with activated', () => {
    const patternRule = {
      field: ALARM_PATTERN_FIELDS.activationDate,
      cond: { type: PATTERN_CONDITIONS.exist, value: true },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.activated,
      operator: PATTERN_OPERATORS.activated,
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with inactive', () => {
    const patternRule = {
      field: ALARM_PATTERN_FIELDS.activationDate,
      cond: { type: PATTERN_CONDITIONS.exist, value: false },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.activated,
      operator: PATTERN_OPERATORS.inactive,
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with is meta alarm', () => {
    const patternRule = {
      field: ALARM_PATTERN_FIELDS.meta,
      cond: { type: PATTERN_CONDITIONS.exist, value: true },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.meta,
      operator: PATTERN_OPERATORS.isMetaAlarm,
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with inactive', () => {
    const patternRule = {
      field: ALARM_PATTERN_FIELDS.meta,
      cond: { type: PATTERN_CONDITIONS.exist, value: false },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.meta,
      operator: PATTERN_OPERATORS.isNotMetaAlarm,
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with `relative time` condition for activation date', () => {
    const patternRule = {
      field: ALARM_PATTERN_FIELDS.activationDate,
      cond: {
        type: PATTERN_CONDITIONS.relativeTime,
        value: {
          value: 15,
          unit: TIME_UNITS.minute,
        },
      },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.activationDate,
      operator: PATTERN_OPERATORS.within,
      range: {
        type: QUICK_RANGES.last15Minutes.value,
        from: 0,
        to: 0,
      },
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with `absolute time` condition for activation date', () => {
    const value = {
      from: Faker.datatype.number(),
      to: Faker.datatype.number(),
    };
    const patternRule = {
      field: ALARM_PATTERN_FIELDS.activationDate,
      cond: { type: PATTERN_CONDITIONS.absoluteTime, value },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.activationDate,
      operator: PATTERN_OPERATORS.inRangeDates,
      range: {
        type: QUICK_RANGES.last1Hour.value,
        from: new Date(value.from * 1000),
        to: new Date(value.to * 1000),
      },
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with `relative time` condition and custom duration for `within` operator', () => {
    const customDuration = {
      value: 17,
      unit: TIME_UNITS.hour,
    };
    const patternRule = {
      field: ALARM_PATTERN_FIELDS.creationDate,
      cond: {
        type: PATTERN_CONDITIONS.relativeTime,
        value: customDuration,
      },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.creationDate,
      operator: PATTERN_OPERATORS.within,
      range: {
        type: QUICK_RANGES.custom.value,
        typeCustom: customDuration,
        from: 0,
        to: 0,
      },
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with `relative time` condition and custom duration for `olderThan` operator', () => {
    const customDuration = {
      value: 55,
      unit: TIME_UNITS.day,
    };
    const patternRule = {
      field: ALARM_PATTERN_FIELDS.creationDate,
      cond: {
        type: PATTERN_CONDITIONS.relativeTime,
        value: {
          to: customDuration,
        },
      },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.creationDate,
      operator: PATTERN_OPERATORS.olderThan,
      range: {
        type: QUICK_RANGES.custom.value,
        typeCustom: customDuration,
        from: 0,
        to: 0,
      },
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with `relative time` condition and custom durations for `inRangePeriod` operator with custom `from` range', () => {
    const customFromDuration = {
      value: 67,
      unit: TIME_UNITS.week,
    };
    const toDuration = {
      value: 1,
      unit: TIME_UNITS.hour,
    };
    const patternRule = {
      field: ALARM_PATTERN_FIELDS.activationDate,
      cond: {
        type: PATTERN_CONDITIONS.relativeTime,
        value: {
          from: customFromDuration,
          to: toDuration,
        },
      },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.activationDate,
      operator: PATTERN_OPERATORS.inRangePeriod,
      range: {
        type: QUICK_RANGES.last1Hour.value,
        from: QUICK_RANGES.custom.value,
        fromCustom: customFromDuration,
        to: QUICK_RANGES.last1Hour.value,
      },
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with `relative time` condition and custom durations for `inRangePeriod` operator with custom `to` range', () => {
    const fromDuration = {
      value: 3,
      unit: TIME_UNITS.hour,
    };
    const customToDuration = {
      value: 23, // Use a value that doesn't match predefined ranges
      unit: TIME_UNITS.day,
    };
    const patternRule = {
      field: ALARM_PATTERN_FIELDS.activationDate,
      cond: {
        type: PATTERN_CONDITIONS.relativeTime,
        value: {
          from: fromDuration,
          to: customToDuration,
        },
      },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.activationDate,
      operator: PATTERN_OPERATORS.inRangePeriod,
      range: {
        type: QUICK_RANGES.last1Hour.value,
        from: QUICK_RANGES.last3Hour.value,
        to: QUICK_RANGES.custom.value,
        toCustom: customToDuration,
      },
    });

    const convertedBack = formRuleToPatternRule(form);

    expect(convertedBack).toEqual({
      field: ALARM_PATTERN_FIELDS.activationDate,
      cond: {
        type: PATTERN_CONDITIONS.relativeTime,
        value: {
          from: fromDuration,
          to: customToDuration,
        },
      },
    });
  });

  it('should be converted to form and back to pattern with `relative time` condition and both custom durations for `inRangePeriod` operator', () => {
    const customFromDuration = {
      value: 17, // Use a specific value that doesn't match predefined ranges
      unit: TIME_UNITS.minute,
    };
    const customToDuration = {
      value: 4, // Use a specific value that doesn't match predefined ranges
      unit: TIME_UNITS.minute,
    };
    const patternRule = {
      field: ALARM_PATTERN_FIELDS.creationDate,
      cond: {
        type: PATTERN_CONDITIONS.relativeTime,
        value: {
          from: customFromDuration,
          to: customToDuration,
        },
      },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.creationDate,
      operator: PATTERN_OPERATORS.inRangePeriod,
      range: {
        type: QUICK_RANGES.last1Hour.value,
        from: QUICK_RANGES.custom.value,
        fromCustom: customFromDuration,
        to: QUICK_RANGES.custom.value,
        toCustom: customToDuration,
      },
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });
});

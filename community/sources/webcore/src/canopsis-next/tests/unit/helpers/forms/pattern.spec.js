import Faker from 'faker';

import {
  ALARM_PATTERN_FIELDS,
  ENTITY_PATTERN_FIELDS,
  PATTERN_CONDITIONS,
  PATTERN_FIELD_TYPES,
  PATTERN_OPERATORS,
  PATTERN_RULE_INFOS_FIELDS,
  QUICK_RANGES,
  TIME_UNITS,
} from '@/constants';

import { formRuleToPatternRule, patternRuleToForm } from '@/helpers/entities/pattern/form';
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
    range: {
      type: QUICK_RANGES.last1Hour.value,
      from: 0,
      to: 0,
    },
    duration: durationToForm(),
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
      cond: { type: PATTERN_CONDITIONS.beginsWith, value },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.connector,
      operator: PATTERN_OPERATORS.beginsWith,
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

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.connectorName,
      operator: PATTERN_OPERATORS.hasEvery,
      value,
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with `has one of` operator', () => {
    const value = [Faker.lorem.word()];
    const patternRule = {
      field: ALARM_PATTERN_FIELDS.connectorName,
      cond: { type: PATTERN_CONDITIONS.hasOneOf, value },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.connectorName,
      operator: PATTERN_OPERATORS.hasOneOf,
      value,
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });

  it('should be converted to form and back to pattern with `has not` operator', () => {
    const value = [Faker.lorem.word()];
    const patternRule = {
      field: ALARM_PATTERN_FIELDS.connectorName,
      cond: { type: PATTERN_CONDITIONS.hasNot, value },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.connectorName,
      operator: PATTERN_OPERATORS.hasNot,
      value,
    });
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

  it('should be converted to form and back to pattern with `relative time` condition', () => {
    const lastHour = 3600;
    const patternRule = {
      field: ALARM_PATTERN_FIELDS.creationDate,
      cond: {
        type: PATTERN_CONDITIONS.relativeTime,
        value: {
          value: lastHour,
          unit: TIME_UNITS.second,
        },
      },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.creationDate,
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
      range: {
        type: QUICK_RANGES.custom.value,
        ...value,
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

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.infos,
      field: PATTERN_RULE_INFOS_FIELDS.value,
      fieldType: PATTERN_FIELD_TYPES.stringArray,
      operator: PATTERN_OPERATORS.hasNot,
      dictionary,
      value,
    });
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

  it('should be converted to form and back to pattern with `relative time` condition for activation date', () => {
    const lastHour = 3600;
    const patternRule = {
      field: ALARM_PATTERN_FIELDS.activationDate,
      cond: {
        type: PATTERN_CONDITIONS.relativeTime,
        value: {
          value: lastHour,
          unit: TIME_UNITS.second,
        },
      },
    };

    const form = patternRuleToForm(patternRule);

    expect(form).toEqual({
      ...defaultForm,
      attribute: ALARM_PATTERN_FIELDS.activationDate,
      range: {
        type: QUICK_RANGES.last1Hour.value,
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
      range: {
        type: QUICK_RANGES.custom.value,
        ...value,
      },
    });
    expect(formRuleToPatternRule(form)).toEqual(patternRule);
  });
});

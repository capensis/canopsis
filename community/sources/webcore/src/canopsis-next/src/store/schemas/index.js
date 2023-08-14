import { schema } from 'normalizr';

import { ENTITIES_TYPES } from '@/constants';

import {
  childProcessStrategy,
  childMergeStrategy,
  parentProcessStrategy,
} from './helpers';

export const pbehaviorSchema = new schema.Entity(ENTITIES_TYPES.pbehavior, {}, {
  idAttribute: '_id',
  processStrategy: childProcessStrategy,
  mergeStrategy: childMergeStrategy,
});

export const alarmSchema = new schema.Entity(ENTITIES_TYPES.alarm, {}, {
  idAttribute: '_id',
  processStrategy: parentProcessStrategy,
});

alarmSchema.disabledCache = true;

export const alarmDetailsSchema = new schema.Entity(ENTITIES_TYPES.alarmDetails, {}, {
  idAttribute: '_id',
});

export const entitySchema = new schema.Entity(ENTITIES_TYPES.entity, {
  pbehaviors: [pbehaviorSchema],
}, {
  idAttribute: '_id',
  processStrategy: parentProcessStrategy,
});

export const serviceSchema = new schema.Entity(ENTITIES_TYPES.service, {}, { idAttribute: '_id' });

export const weatherServiceSchema = new schema.Entity(ENTITIES_TYPES.weatherService, {}, { idAttribute: '_id' });

export const userPreferenceSchema = new schema.Entity(ENTITIES_TYPES.userPreference, {}, {
  idAttribute: 'widget',
});

export const snmpRuleSchema = new schema.Entity(ENTITIES_TYPES.snmpRule, {}, { idAttribute: '_id' });

export const testSuiteSchema = new schema.Entity(ENTITIES_TYPES.testSuite, {}, { idAttribute: '_id' });

export const testSuiteHistorySchema = new schema.Entity(ENTITIES_TYPES.testSuiteHistory, {}, { idAttribute: '_id' });

export default {
  [ENTITIES_TYPES.alarm]: alarmSchema,
  [ENTITIES_TYPES.alarmDetails]: alarmDetailsSchema,
  [ENTITIES_TYPES.entity]: entitySchema,
  [ENTITIES_TYPES.service]: serviceSchema,
  [ENTITIES_TYPES.weatherService]: weatherServiceSchema,
  [ENTITIES_TYPES.pbehavior]: pbehaviorSchema,
  [ENTITIES_TYPES.userPreference]: userPreferenceSchema,
  [ENTITIES_TYPES.snmpRule]: snmpRuleSchema,
  [ENTITIES_TYPES.testSuite]: testSuiteSchema,
  [ENTITIES_TYPES.testSuiteHistory]: testSuiteHistorySchema,
};

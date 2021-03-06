import { schema } from 'normalizr';

import { ENTITIES_TYPES } from '@/constants';

import {
  childProcessStrategy,
  childMergeStrategy,
  parentProcessStrategy,
  viewTabProcessStrategy,
} from './helpers';

export const pbehaviorSchema = new schema.Entity(ENTITIES_TYPES.pbehavior, {}, {
  idAttribute: '_id',
  processStrategy: childProcessStrategy,
  mergeStrategy: childMergeStrategy,
});

export const alarmSchema = new schema.Entity(ENTITIES_TYPES.alarm, {
  pbehavior: pbehaviorSchema,
}, {
  idAttribute: '_id',
  processStrategy: parentProcessStrategy,
});

alarmSchema.define({
  consequences: {
    data: [alarmSchema],
  },
  causes: {
    data: [alarmSchema],
  },
});

alarmSchema.disabledCache = true;

export const entitySchema = new schema.Entity(ENTITIES_TYPES.entity, {
  pbehaviors: [pbehaviorSchema],
}, {
  idAttribute: '_id',
  processStrategy: parentProcessStrategy,
});

export const watcherSchema = new schema.Entity(ENTITIES_TYPES.watcher, {}, { idAttribute: '_id' });

export const watcherEntitySchema = new schema.Entity(ENTITIES_TYPES.watcherEntity, {}, { idAttribute: '_id' });

export const userPreferenceSchema = new schema.Entity(ENTITIES_TYPES.userPreference, {}, {
  idAttribute: '_id',
});

export const widgetSchema = new schema.Entity(ENTITIES_TYPES.widget, {}, {
  idAttribute: '_id',
});

export const viewTabSchema = new schema.Entity(ENTITIES_TYPES.viewTab, {
  widgets: [widgetSchema],
}, {
  idAttribute: '_id',
  processStrategy: viewTabProcessStrategy,
  mergeStrategy: childMergeStrategy,
});

export const viewSchema = new schema.Entity(ENTITIES_TYPES.view, {
  tabs: [viewTabSchema],
}, {
  idAttribute: '_id',
  processStrategy: childProcessStrategy,
  mergeStrategy: childMergeStrategy,
});

export const groupSchema = new schema.Entity(ENTITIES_TYPES.group, {
  views: [viewSchema],
}, {
  idAttribute: '_id',
  processStrategy: parentProcessStrategy,
});

groupSchema.disabledCache = true;

export const userSchema = new schema.Entity(ENTITIES_TYPES.user, {}, { idAttribute: '_id' });

export const roleSchema = new schema.Entity(ENTITIES_TYPES.role, {}, { idAttribute: '_id' });

export const eventFilterRuleSchema = new schema.Entity(ENTITIES_TYPES.eventFilterRule, {}, { idAttribute: '_id' });

export const metaAlarmRuleSchema = new schema.Entity(ENTITIES_TYPES.metaAlarmRule, {}, { idAttribute: '_id' });

export const webhookSchema = new schema.Entity(ENTITIES_TYPES.webhook, {}, { idAttribute: '_id' });

export const snmpRuleSchema = new schema.Entity(ENTITIES_TYPES.snmpRule, {}, { idAttribute: '_id' });

export const actionSchema = new schema.Entity(ENTITIES_TYPES.action, {}, { idAttribute: '_id' });

export const heartbeatSchema = new schema.Entity(ENTITIES_TYPES.heartbeat, {}, { idAttribute: '_id' });

export const dynamicInfoSchema = new schema.Entity(ENTITIES_TYPES.dynamicInfo, {}, { idAttribute: '_id' });

export const broadcastMessageSchema = new schema.Entity(ENTITIES_TYPES.broadcastMessage, {}, { idAttribute: '_id' });

export const playlistSchema = new schema.Entity(ENTITIES_TYPES.playlist, {
  tabs: [viewTabSchema],
}, { idAttribute: '_id' });

export const pbehaviorTypesSchema = new schema.Entity(ENTITIES_TYPES.pbehaviorTypes, {}, { idAttribute: '_id' });

export const pbehaviorReasonsSchema = new schema.Entity(ENTITIES_TYPES.pbehaviorReasons, {}, { idAttribute: '_id' });

export const pbehaviorExceptionsSchema = new schema.Entity(ENTITIES_TYPES.pbehaviorExceptions, {}, { idAttribute: '_id' });

export const remediationInstructionSchema = new schema.Entity(ENTITIES_TYPES.remediationInstruction, {}, { idAttribute: '_id' });

export const remediationJobSchema = new schema.Entity(ENTITIES_TYPES.remediationJob, {}, { idAttribute: '_id' });

export const remediationConfigurationSchema = new schema.Entity(ENTITIES_TYPES.remediationConfiguration, {}, { idAttribute: '_id' });

export default {
  [ENTITIES_TYPES.alarm]: alarmSchema,
  [ENTITIES_TYPES.entity]: entitySchema,
  [ENTITIES_TYPES.watcher]: watcherSchema,
  [ENTITIES_TYPES.watcherEntity]: watcherEntitySchema,
  [ENTITIES_TYPES.pbehavior]: pbehaviorSchema,
  [ENTITIES_TYPES.userPreference]: userPreferenceSchema,
  [ENTITIES_TYPES.group]: groupSchema,
  [ENTITIES_TYPES.view]: viewSchema,
  [ENTITIES_TYPES.viewTab]: viewTabSchema,
  [ENTITIES_TYPES.widget]: widgetSchema,
  [ENTITIES_TYPES.user]: userSchema,
  [ENTITIES_TYPES.role]: roleSchema,
  [ENTITIES_TYPES.eventFilterRule]: eventFilterRuleSchema,
  [ENTITIES_TYPES.webhook]: webhookSchema,
  [ENTITIES_TYPES.snmpRule]: snmpRuleSchema,
  [ENTITIES_TYPES.action]: actionSchema,
  [ENTITIES_TYPES.heartbeat]: heartbeatSchema,
  [ENTITIES_TYPES.dynamicInfo]: dynamicInfoSchema,
  [ENTITIES_TYPES.broadcastMessage]: broadcastMessageSchema,
  [ENTITIES_TYPES.playlist]: playlistSchema,
  [ENTITIES_TYPES.metaAlarmRule]: metaAlarmRuleSchema,
  [ENTITIES_TYPES.pbehaviorTypes]: pbehaviorTypesSchema,
  [ENTITIES_TYPES.pbehaviorReasons]: pbehaviorReasonsSchema,
  [ENTITIES_TYPES.pbehaviorExceptions]: pbehaviorExceptionsSchema,
  [ENTITIES_TYPES.remediationInstruction]: remediationInstructionSchema,
  [ENTITIES_TYPES.remediationJob]: remediationJobSchema,
  [ENTITIES_TYPES.remediationConfiguration]: remediationConfigurationSchema,
};

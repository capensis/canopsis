import { schema } from 'normalizr';

import { ENTITIES_TYPES } from '@/constants';
import { childProcessStrategy, childMergeStrategy, parentProcessStrategy } from '@/helpers/schema';

export const pbehaviorSchema = new schema.Entity(ENTITIES_TYPES.pbehavior, {}, {
  idAttribute: '_id',
  processStrategy: childProcessStrategy,
  mergeStrategy: childMergeStrategy,
});

export const alarmSchema = new schema.Entity(ENTITIES_TYPES.alarm, {
  pbehaviors: [pbehaviorSchema],
}, {
  idAttribute: '_id',
  processStrategy: parentProcessStrategy,
});

export const entitySchema = new schema.Entity(ENTITIES_TYPES.entity, {
  pbehaviors: [pbehaviorSchema],
}, {
  idAttribute: '_id',
  processStrategy: parentProcessStrategy,
});

export const watcherSchema = new schema.Entity(ENTITIES_TYPES.watcher, {}, { idAttribute: 'entity_id' });

export const watcherEntitySchema = new schema.Entity(ENTITIES_TYPES.watcherEntity, {}, {
  idAttribute: 'entity_id',
});

export const userPreferenceSchema = new schema.Entity(ENTITIES_TYPES.userPreference, {}, {
  idAttribute: '_id',
});

export const widgetSchema = new schema.Entity(ENTITIES_TYPES.widget, {}, {
  idAttribute: '_id',
});

export const viewRowSchema = new schema.Entity(ENTITIES_TYPES.viewRow, {
  widgets: [widgetSchema],
}, { idAttribute: '_id' });

export const viewTabSchema = new schema.Entity(ENTITIES_TYPES.viewTab, {
  rows: [viewRowSchema],
}, { idAttribute: '_id' });

export const viewSchema = new schema.Entity(ENTITIES_TYPES.view, {
  tabs: [viewTabSchema],
}, { idAttribute: '_id' });

export const groupSchema = new schema.Entity(ENTITIES_TYPES.group, {
  views: [viewSchema],
}, {
  idAttribute: '_id',
});

export const userSchema = new schema.Entity(ENTITIES_TYPES.user, {}, { idAttribute: '_id' });

export const roleSchema = new schema.Entity(ENTITIES_TYPES.role, {}, { idAttribute: '_id' });

export const eventFilterRuleSchema = new schema.Entity(ENTITIES_TYPES.eventFilterRule, {}, { idAttribute: '_id' });

export const webhookSchema = new schema.Entity(ENTITIES_TYPES.webhook, {}, { idAttribute: '_id' });

export const snmpRuleSchema = new schema.Entity(ENTITIES_TYPES.snmpRule, {}, { idAttribute: '_id' });

export const actionSchema = new schema.Entity(ENTITIES_TYPES.action, {}, { idAttribute: '_id' });

export const heartbeatSchema = new schema.Entity(ENTITIES_TYPES.heartbeat, {}, { idAttribute: '_id' });

export default {
  [ENTITIES_TYPES.alarm]: alarmSchema,
  [ENTITIES_TYPES.entity]: entitySchema,
  [ENTITIES_TYPES.watcher]: watcherSchema,
  [ENTITIES_TYPES.watcherEntity]: watcherEntitySchema,
  [ENTITIES_TYPES.pbehavior]: pbehaviorSchema,
  [ENTITIES_TYPES.userPreference]: userPreferenceSchema,
  [ENTITIES_TYPES.group]: groupSchema,
  [ENTITIES_TYPES.view]: viewSchema,
  [ENTITIES_TYPES.viewRow]: viewRowSchema,
  [ENTITIES_TYPES.viewTab]: viewTabSchema,
  [ENTITIES_TYPES.widget]: widgetSchema,
  [ENTITIES_TYPES.user]: userSchema,
  [ENTITIES_TYPES.role]: roleSchema,
  [ENTITIES_TYPES.eventFilterRule]: eventFilterRuleSchema,
  [ENTITIES_TYPES.webhook]: webhookSchema,
  [ENTITIES_TYPES.snmpRule]: snmpRuleSchema,
  [ENTITIES_TYPES.action]: actionSchema,
  [ENTITIES_TYPES.heartbeat]: heartbeatSchema,
};

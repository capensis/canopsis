import { schema } from 'normalizr';

import { ENTITIES_TYPES } from '@/constants';
import { childProcessStrategy, childMergeStrategy } from '@/helpers/schema';

export const pbehaviorSchema = new schema.Entity(ENTITIES_TYPES.pbehavior, {}, {
  idAttribute: '_id',
  processStrategy: childProcessStrategy,
  mergeStrategy: childMergeStrategy,
});

export const alarmSchema = new schema.Entity(ENTITIES_TYPES.alarm, {
  pbehaviors: [pbehaviorSchema],
}, {
  idAttribute: '_id',
  processStrategy: entity =>
    ({
      ...entity,
      _embedded: {
        type: ENTITIES_TYPES.alarm,
      },
    }),
});

export const entitySchema = new schema.Entity(ENTITIES_TYPES.entity, {}, { idAttribute: '_id' });

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

export const rowSchema = new schema.Entity(ENTITIES_TYPES.row, {
  widgets: [widgetSchema],
}, { idAttribute: '_id' });

export const viewSchema = new schema.Entity(ENTITIES_TYPES.view, {
  rows: [rowSchema],
}, { idAttribute: '_id' });

export const groupSchema = new schema.Entity(ENTITIES_TYPES.group, {
  views: [viewSchema],
}, {
  idAttribute: '_id',
});


export const userSchema = new schema.Entity(ENTITIES_TYPES.user);

export const roleSchema = new schema.Entity(ENTITIES_TYPES.role, {}, { idAttribute: '_id' });


export default {
  [ENTITIES_TYPES.alarm]: alarmSchema,
  [ENTITIES_TYPES.entity]: entitySchema,
  [ENTITIES_TYPES.watcher]: watcherSchema,
  [ENTITIES_TYPES.watcherEntity]: watcherEntitySchema,
  [ENTITIES_TYPES.pbehavior]: pbehaviorSchema,
  [ENTITIES_TYPES.userPreference]: userPreferenceSchema,
  [ENTITIES_TYPES.group]: groupSchema,
  [ENTITIES_TYPES.view]: viewSchema,
  [ENTITIES_TYPES.row]: rowSchema,
  [ENTITIES_TYPES.widget]: widgetSchema,
  [ENTITIES_TYPES.user]: userSchema,
  [ENTITIES_TYPES.role]: roleSchema,
};

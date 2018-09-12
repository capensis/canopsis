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

export const widgetSchema = new schema.Entity(ENTITIES_TYPES.widget);

export const widgetWrapperSchema = new schema.Entity(ENTITIES_TYPES.widgetWrapper, {
  widget: widgetSchema,
});

widgetSchema.define({ items: [widgetWrapperSchema] });

export const viewSchema = new schema.Entity(ENTITIES_TYPES.view, {
  containerwidget: {
    items: [widgetWrapperSchema],
  },
});

export const viewV3Schema = new schema.Entity(ENTITIES_TYPES.viewV3, {}, { idAttribute: '_id' });

export const groupSchema = new schema.Entity(ENTITIES_TYPES.group, { views: [viewV3Schema] }, {
  idAttribute: 'name',
});

export default {
  [ENTITIES_TYPES.alarm]: alarmSchema,
  [ENTITIES_TYPES.entity]: entitySchema,
  [ENTITIES_TYPES.watcher]: watcherSchema,
  [ENTITIES_TYPES.watcherEntity]: watcherEntitySchema,
  [ENTITIES_TYPES.pbehavior]: pbehaviorSchema,
  [ENTITIES_TYPES.userPreference]: userPreferenceSchema,
  [ENTITIES_TYPES.view]: viewSchema,
  [ENTITIES_TYPES.viewV3]: viewV3Schema,
  [ENTITIES_TYPES.group]: groupSchema,
  [ENTITIES_TYPES.widgetWrapper]: widgetWrapperSchema,
  [ENTITIES_TYPES.widget]: widgetSchema,
};

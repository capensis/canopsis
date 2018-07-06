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

export const userPreferenceSchema = new schema.Entity(ENTITIES_TYPES.userPreference, {}, {
  idAttribute: '_id',
});

const widgetWrapper = new schema.Entity(ENTITIES_TYPES.widgetWrapper);

const widgetWrapperList = new schema.Array(widgetWrapper);

export const viewSchema = new schema.Entity(ENTITIES_TYPES.view, {
  containerwidget: {
    items: widgetWrapperList,
  },
});

export default {
  [ENTITIES_TYPES.alarm]: alarmSchema,
  [ENTITIES_TYPES.entity]: entitySchema,
  [ENTITIES_TYPES.pbehavior]: pbehaviorSchema,
  [ENTITIES_TYPES.userPreference]: userPreferenceSchema,
  [ENTITIES_TYPES.view]: viewSchema,
};

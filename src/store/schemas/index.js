import { schema } from 'normalizr';
import { ENTITIES_TYPES } from '@/constants';

export const pbehaviorSchema = new schema.Entity(ENTITIES_TYPES.pbehavior, {}, {
  idAttribute: '_id',
  processStrategy: (entity, parent, key) =>
    ({
      ...entity,
      _embedded: {
        parentType: ENTITIES_TYPES.alarm,
        parentId: parent._id,
        relationType: key,
      },
    }),
});

export const alarmSchema = new schema.Entity(ENTITIES_TYPES.alarm, {
  pbehaviors: [pbehaviorSchema],
}, {
  idAttribute: '_id',
});

export const contextSchema = new schema.Entity(ENTITIES_TYPES.context, {}, { idAttribute: '_id' });

export const userPreferenceSchema = new schema.Entity(ENTITIES_TYPES.userPreference, {}, {
  idAttribute: '_id',
});

export default {
  [ENTITIES_TYPES.alarm]: alarmSchema,
  [ENTITIES_TYPES.context]: contextSchema,
  [ENTITIES_TYPES.pbehavior]: pbehaviorSchema,
  [ENTITIES_TYPES.userPreference]: userPreferenceSchema,
};

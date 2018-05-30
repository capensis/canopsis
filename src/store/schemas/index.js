import { schema } from 'normalizr';

export const types = {
  ALARM: 'alarm',
  CONTEXT: 'context',
  EVENT: 'event',
  PBEHAVIOR: 'pbehavior',
  USER_PREFERENCE: 'userPreference',
};

export const pbehaviorSchema = new schema.Entity(types.PBEHAVIOR, {}, {
  idAttribute: '_id',
  processStrategy: (entity, parent, key) =>
    ({
      ...entity,
      _embedded: {
        parentType: types.ALARM,
        parentId: parent._id,
        relationType: key,
      },
    }),
});

export const alarmSchema = new schema.Entity(types.ALARM, {
  pbehaviors: [pbehaviorSchema],
}, {
  idAttribute: '_id',
});

export const contextSchema = new schema.Entity(types.CONTEXT, {}, { idAttribute: '_id' });

export const eventSchema = new schema.Entity(types.EVENT, {}, { idAttribute: '_id' });

export const userPreferenceSchema = new schema.Entity(types.USER_PREFERENCE, {}, {
  idAttribute: '_id',
});

export default {
  [types.ALARM]: alarmSchema,
  [types.CONTEXT]: contextSchema,
  [types.PBEHAVIOR]: pbehaviorSchema,
  [types.EVENT]: eventSchema,
  [types.USER_PREFERENCE]: userPreferenceSchema,
};

import { schema } from 'normalizr';
import { ENTITIES_TYPES } from '@/constants';

// TODO: move it
const processStrategy = (entity, parent, key) => ({
  ...entity,
  _embedded: {
    parents: [{ type: parent._embedded.type, id: parent._id, key }],
  },
});

// TODO: move it
const mergeStrategy = (entityA, entityB) => ({
  ...entityA,
  ...entityB,
  _embedded: {
    parents: [...entityA._embedded.parents, ...entityB._embedded.parents],
  },
});

export const pbehaviorSchema = new schema.Entity(ENTITIES_TYPES.pbehavior, {}, {
  idAttribute: '_id',
  processStrategy,
  mergeStrategy,
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

export const contextSchema = new schema.Entity(ENTITIES_TYPES.context, {}, { idAttribute: '_id' });

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
  [ENTITIES_TYPES.context]: contextSchema,
  [ENTITIES_TYPES.pbehavior]: pbehaviorSchema,
  [ENTITIES_TYPES.userPreference]: userPreferenceSchema,
  [ENTITIES_TYPES.view]: viewSchema,
};

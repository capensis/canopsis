import { schema } from 'normalizr';
import { get } from 'lodash';

import { SCHEMA_EMBEDDED_KEY } from '@/config';
import { WIDGET_TYPES } from '@/constants';

import { defaultSanitizer, getDefaultSanitizerForArray } from '@/helpers/sanitizer';
import { setSeveralFields } from '@/helpers/immutable';

/**
 * If parent has children we should use this processStrategy
 *
 * @param entity
 */
export function parentProcessStrategy(entity) {
  return {
    ...entity,

    [SCHEMA_EMBEDDED_KEY]: {
      type: this.key,
    },
  };
}

/**
 * If entity has parent we should use this processStrategy
 *
 * @param entity
 * @param parent
 * @param key
 */
export function childProcessStrategy(entity, parent, key) {
  const result = parentProcessStrategy.call(this, entity);

  if (parent && parent[SCHEMA_EMBEDDED_KEY]) {
    result[SCHEMA_EMBEDDED_KEY] = {
      ...result[SCHEMA_EMBEDDED_KEY],

      parents: [{ type: parent[SCHEMA_EMBEDDED_KEY].type, id: parent._id, key }],
    };
  }

  return result;
}

/**
 * If entity has parent we should use this mergeStrategy
 *
 * @param entityA
 * @param entityB
 */
export const childMergeStrategy = (entityA, entityB) => {
  const result = {
    ...entityA,
    ...entityB,
  };

  if (entityA[SCHEMA_EMBEDDED_KEY] || entityB[SCHEMA_EMBEDDED_KEY]) {
    const embeddedParentsKey = `${SCHEMA_EMBEDDED_KEY}.parents`;

    result[SCHEMA_EMBEDDED_KEY] = {
      parents: [
        ...get(entityA, embeddedParentsKey, []),
        ...get(entityB, embeddedParentsKey, []),
      ],
    };
  }

  return result;
};

/**
 * Process strategy for the widget schema
 *
 * @param {Object} widget
 * @returns {{parameters: (Object|Array)}|*}
 */
export const widgetProcessStrategy = (widget) => {
  const templatesMap = {
    [WIDGET_TYPES.alarmList]: {
      moreInfoTemplate: defaultSanitizer,
      infoPopups: getDefaultSanitizerForArray('template'),
    },

    [WIDGET_TYPES.weather]: {
      entityTemplate: defaultSanitizer,
      modalTemplate: defaultSanitizer,
      blockTemplate: defaultSanitizer,
    },

    [WIDGET_TYPES.text]: {
      template: defaultSanitizer,
    },
  };

  const widgetTemplatesMap = templatesMap[widget.type];

  if (widgetTemplatesMap) {
    const parameters = setSeveralFields(widget.parameters, widgetTemplatesMap);

    return {
      ...widget,

      parameters,
    };
  }

  return widget;
};

/* eslint-disable */
/**
 * We reinitialized denormalize method for removing our SCHEMA_EMBEDDED_KEY property if we need
 *
 * @param entity
 * @param unvisit
 * @return {*}
 */
schema.Entity.prototype.denormalize = function denormalize(entity, unvisit) {
  if (!this[SCHEMA_EMBEDDED_KEY] && entity.hasOwnProperty(SCHEMA_EMBEDDED_KEY)) {
    delete entity[SCHEMA_EMBEDDED_KEY];
  }

  Object.keys(this.schema).forEach((key) => {
    if (entity.hasOwnProperty(key)) {
      const schema = this.schema[key];

      entity[key] = unvisit(entity[key], schema);
    }
  });

  return entity;
};
/* eslint-enable */

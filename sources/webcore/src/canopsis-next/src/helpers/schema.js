import { schema } from 'normalizr';
import { get, omit } from 'lodash';

import { SCHEMA_EMBEDDED_KEY } from '@/config';
import { WIDGET_GRID_SIZES_KEYS } from '@/constants';

/**
 * If parent has children we should use this processStrategy
 *
 * @param {Object} entity
 * @return {Object}
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
 * @param {Object} entity
 * @param {Object} parent
 * @param {string} key
 * @return {Object}
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
 * @param {Object} entityA
 * @param {Object} entityB
 * @return {Object}
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
 * Special processStrategy for viewTab entity schema
 *
 * @param {Object} entity
 * @param {Object} parent
 * @param {string} key
 * @return {Object}
 */
export function viewTabProcessStrategy(entity, parent, key) {
  const newEntity = childProcessStrategy.call(this, entity, parent, key);

  if (!newEntity.grid || !newEntity.widgets) {
    newEntity.grid = {};

    newEntity.widgets = newEntity.rows.reduce((acc, { widgets }, rowIndex) => {
      const prevEnd = {
        [WIDGET_GRID_SIZES_KEYS.mobile]: 0,
        [WIDGET_GRID_SIZES_KEYS.tablet]: 0,
        [WIDGET_GRID_SIZES_KEYS.desktop]: 0,
      };

      const GRID_SIZES_MAP = {
        [WIDGET_GRID_SIZES_KEYS.mobile]: 'sm',
        [WIDGET_GRID_SIZES_KEYS.tablet]: 'md',
        [WIDGET_GRID_SIZES_KEYS.desktop]: 'lg',
      };

      widgets.forEach((widget) => {
        const newWidget = omit(widget, 'size');

        newWidget.gridParameters = Object.values(WIDGET_GRID_SIZES_KEYS).reduce((secondAcc, size) => {
          const width = widget.size[GRID_SIZES_MAP[size]];

          // eslint-disable-next-line no-param-reassign
          secondAcc[size] = {
            x: prevEnd[size],
            y: rowIndex,
            w: width,
            h: 0,
            autoHeight: false,
          };

          prevEnd[size] += width;

          return secondAcc;
        }, {});

        acc.push(newWidget);
      });

      return acc;
    }, []);

    delete newEntity.rows;
  }

  return newEntity;
}

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

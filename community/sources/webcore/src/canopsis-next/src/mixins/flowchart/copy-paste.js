import { isObject, isString } from 'lodash';

import uid from '@/helpers/uid';
import { readTextFromClipboard, writeTextToClipboard } from '@/helpers/clipboard';

export const copyPasteShapesMixin = {
  methods: {
    copySelectedShapes() {
      const data = this.selectedIds.reduce((acc, id) => {
        const shape = this.data[id];

        /* Clear connections */
        const connections = shape.connections.filter(connection => this.selectedIds.includes(connection.shapeId));
        const connectedTo = shape.connectedTo.filter(shapeId => this.selectedIds.includes(shapeId));

        acc[id] = {
          ...shape,
          connections,
          connectedTo,
        };

        return acc;
      }, {});

      writeTextToClipboard(JSON.stringify(data));
    },

    async pasteShapes() {
      const data = await readTextFromClipboard();

      if (!isString(data)) {
        return;
      }

      const shapes = JSON.parse(data);

      if (!isObject(shapes)) {
        return;
      }

      Object.entries(shapes)
        .forEach(([id, shape]) => {
          /* Check if id not exist in current shapes */
          if (!this.data[id] && !shapes[id]) {
            return;
          }

          const newId = uid();

          /* Change id in connected shapes */
          shape.connectedTo.forEach((connectedShapeId) => {
            const connectedShape = shapes[connectedShapeId];

            connectedShape.connections = connectedShape.connections.map(connection => ({
              ...connection,
              shapeId: connection.shapeId === id ? newId : connection.shapeId,
            }));
          });

          /* Change id in connected shapes */
          shape.connections.forEach(({ shapeId: connectingShapeId }) => {
            const connectingShape = shapes[connectingShapeId];

            connectingShape.connectedTo = connectingShape.connectedTo
              .map(shapeId => (shapeId === id ? newId : shapeId));
          });

          delete shapes[id];

          shapes[newId] = {
            ...shape,
            _id: newId,
          };
        });

      this.data = {
        ...this.data,
        ...shapes,
      };

      this.setSelected(Object.keys(shapes));
      this.updateShapes(this.data);
    },
  },
};

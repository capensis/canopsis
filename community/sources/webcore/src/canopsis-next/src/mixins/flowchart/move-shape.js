import { SHAPES } from '@/constants';

import { roundByStep } from '@/helpers/flowchart/round';

export const moveShapesMixin = {
  methods: {
    moveShapeById(id, offset) {
      const shape = this.data[id];

      switch (shape.type) {
        case SHAPES.storage:
        case SHAPES.parallelogram:
        case SHAPES.image:
        case SHAPES.circle:
        case SHAPES.rhombus:
        case SHAPES.ellipse:
        case SHAPES.process:
        case SHAPES.document:
        case SHAPES.rect: {
          this.updateShape(shape, { x: shape.x + offset.x, y: shape.y + offset.y });
          break;
        }
        case SHAPES.arrowLine:
        case SHAPES.bidirectionalArrowLine:
        case SHAPES.line: {
          this.updateShape(shape, {
            points: shape.points.map(point => ({
              ...point,
              x: point.x + offset.x,
              y: point.y + offset.y,
            })),
          });
          break;
        }
      }
    },

    moveSelected(offset) {
      this.selectedIds.forEach((id) => {
        this.moveShapeById(id, offset);

        this.clearConnectedTo(id);
      });
    },

    handleShapeMove(event) {
      const { x, y } = this.normalizeCursor({ x: event.clientX, y: event.clientY });

      const newMovingOffsetX = roundByStep(
        x - this.movingStart.x,
        this.gridSize,
      );
      const newMovingOffsetY = roundByStep(
        y - this.movingStart.y,
        this.gridSize,
      );

      this.moveSelected({
        x: newMovingOffsetX - this.movingOffset.x,
        y: newMovingOffsetY - this.movingOffset.y,
      });

      this.movingOffset.x = newMovingOffsetX;
      this.movingOffset.y = newMovingOffsetY;
    },

    moveSelectedDown() {
      this.moveSelected({ x: 0, y: this.gridSize });

      this.updateShapes(this.data);
    },

    moveSelectedUp() {
      this.moveSelected({ x: 0, y: -this.gridSize });

      this.updateShapes(this.data);
    },

    moveSelectedRight() {
      this.moveSelected({ x: this.gridSize, y: 0 });

      this.updateShapes(this.data);
    },

    moveSelectedLeft() {
      this.moveSelected({ x: -this.gridSize, y: 0 });

      this.updateShapes(this.data);
    },
  },
};

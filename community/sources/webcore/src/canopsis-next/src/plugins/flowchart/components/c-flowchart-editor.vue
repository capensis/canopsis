<template lang="pug">
  svg(
    ref="svg",
    width="100%",
    height="100%",
    @click="addShape",
    @mousemove="onContainerMouseMove",
    @mouseup="onContainerMouseUp"
  )
    rect-shape(
      v-for="shape in shapes",
      :key="shape.id",
      :shape="shape",
      :selected="isSelected(shape)",
      @mousedown="onShapeMouseDown(shape, $event)",
      @mouseup="onShapeMouseUp(shape, $event)",
      @start:resize="onShapeStartResize(shape, $event)"
    )
</template>

<script>
import { omit } from 'lodash';

import RectShape from './rect-shape.vue';

const SHAPE_TYPES = {
  rect: 'rect',
};

export default {
  components: { RectShape },
  props: {
    movingStep: {
      type: Number,
      default: 10,
    },
  },
  data() {
    return {
      data: {},
      selected: [],
      moving: false,
      movingStart: {
        x: 0,
        y: 0,
      },
      movingOffset: {
        x: 0,
        y: 0,
      },
    };
  },
  computed: {
    shapes() {
      return this.data;
    },

    hasSelected() {
      return !!this.selected.length;
    },
  },
  mounted() {
    document.addEventListener('keydown', this.onKeyDown);
  },
  beforeDestroy() {
    document.removeEventListener('keydown', this.onKeyDown);
  },
  methods: {
    isSelected(shape) {
      return this.selected.includes(shape.id);
    },

    setSelected(shape) {
      if (!this.isSelected(shape)) {
        this.selected.push(shape.id);
      }
    },

    clearSelected() {
      this.selected = [];
    },

    removeShapeSelected(shape) {
      this.selected = this.selected.filter(id => id !== shape.id);
    },

    convertPositionByStep(coordinate) {
      return Math.round(coordinate / this.movingStep) * this.movingStep;
    },

    addShape() {
      if (this.resizing || this.moving) {
        return;
      }

      if (this.hasSelected) {
        this.clearSelected();
        return;
      }

      const { width, height } = this.$refs.svg.getBoundingClientRect();

      const newShape = {
        id: Date.now(),
        type: SHAPE_TYPES.rect,
        width: 100,
        height: 100,
        fill: 'orange',
        x: this.convertPositionByStep(width / 2),
        y: this.convertPositionByStep(height / 2),
      };

      this.$set(this.data, newShape.id, newShape);

      this.clearSelected();
      this.setSelected(newShape);
    },

    onShapeMouseDown(shape, event) {
      if (!this.hasSelected) {
        this.setSelected(shape);
      }

      if (!this.isSelected(shape) && !event.ctrlKey) {
        this.clearSelected();
        this.setSelected(shape);
      }

      const { offsetX, offsetY } = event;
      this.moving = true;
      this.movingStart = { x: offsetX, y: offsetY };
    },

    onShapeStartResize(shape, direction) {
      this.resizing = {
        id: shape.id,
        direction,
      };
    },

    onShapeMouseUp(shape, event) {
      if (this.resizing || this.movingOffset.x || this.movingOffset.y) {
        return;
      }

      const isShapeSelected = this.isSelected(shape);

      if (isShapeSelected && this.selected.length === 1) {
        return;
      }

      if (!event.ctrlKey) {
        this.clearSelected();
        this.setSelected(shape);

        return;
      }

      if (isShapeSelected) {
        this.removeShapeSelected(shape);
      } else {
        this.setSelected(shape);
      }
    },

    onContainerMouseUp() {
      if (this.resizing) {
        this.resizing = undefined;
        return;
      }

      if (this.moving) {
        this.moving = false;
        this.movingStart = { x: 0, y: 0 };
        this.movingOffset = { x: 0, y: 0 };
      }
    },

    handleShapeResize(event) {
      const { offsetX, offsetY } = event;
      const shape = this.data[this.resizing.id];

      const directionArray = this.resizing.direction.split('');

      /**
       * TODO: Should be used moving step here
       * TODO: Directions should be moved to constants
       */
      directionArray.forEach((direction, index) => {
        switch (direction) {
          case 's': {
            const newHeight = offsetY - shape.y;

            if (newHeight > 0) {
              shape.height = newHeight;
            } else {
              shape.height = Math.abs(newHeight);
              shape.y = offsetY;

              directionArray[index] = 'n';
            }

            break;
          }
          case 'e': {
            const newWidth = offsetX - shape.x;

            if (newWidth > 0) {
              shape.width = newWidth;
            } else {
              shape.width = Math.abs(newWidth);
              shape.x = offsetX;

              directionArray[index] = 'w';
            }

            break;
          }
          case 'w': {
            const newWidth = shape.width + shape.x - offsetX;

            if (newWidth > 0) {
              shape.width = newWidth;
              shape.x = offsetX;
            } else {
              shape.width = Math.abs(newWidth);
              shape.x = offsetX + newWidth;

              directionArray[index] = 'e';
            }

            break;
          }
          case 'n': {
            const newHeight = shape.height + shape.y - offsetY;

            if (newHeight > 0) {
              shape.height = newHeight;
              shape.y = offsetY;
            } else {
              shape.height = Math.abs(newHeight);
              shape.y = offsetY + newHeight;

              directionArray[index] = 's';
            }

            break;
          }
        }
      });

      this.resizing.direction = directionArray.join('');
    },

    handleShapeMove(event) {
      const { offsetX, offsetY } = event;

      const newMovingOffsetX = this.convertPositionByStep(offsetX - this.movingStart.x);
      const newMovingOffsetY = this.convertPositionByStep(offsetY - this.movingStart.y);

      this.selected.forEach((id) => {
        const shape = this.data[id];
        const { x, y } = shape;
        const newX = (x - this.movingOffset.x) + newMovingOffsetX;
        const newY = (y - this.movingOffset.y) + newMovingOffsetY;

        shape.x = newX;
        shape.y = newY;
      });

      this.movingOffset = {
        x: newMovingOffsetX,
        y: newMovingOffsetY,
      };
    },

    onContainerMouseMove(event) {
      if (this.resizing) {
        this.handleShapeResize(event);
        return;
      }

      if (this.moving) {
        this.handleShapeMove(event);
      }
    },

    moveSelectedDown() {
      this.selected.forEach((id) => {
        this.data[id].y += this.movingStep;
      });
    },

    moveSelectedTop() {
      this.selected.forEach((id) => {
        this.data[id].y -= this.movingStep;
      });
    },

    moveSelectedRight() {
      this.selected.forEach((id) => {
        this.data[id].x += this.movingStep;
      });
    },

    moveSelectedLeft() {
      this.selected.forEach((id) => {
        this.data[id].x -= this.movingStep;
      });
    },

    removeSelected() {
      if (this.hasSelected) {
        this.data = omit(this.data, this.selected);
        this.clearSelected();
      }
    },

    onKeyDown(event) {
      const handler = {
        37: this.moveSelectedLeft,
        38: this.moveSelectedTop,
        39: this.moveSelectedRight,
        40: this.moveSelectedDown,
        46: this.removeSelected,
      }[event.keyCode];

      if (handler) {
        event.preventDefault();
        handler();
      }
    },
  },
};
</script>

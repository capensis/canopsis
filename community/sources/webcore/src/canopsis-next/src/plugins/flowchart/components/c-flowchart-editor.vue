<template lang="pug">
  svg(
    ref="svg",
    width="100%",
    height="100%",
    @click="addShape",
    @mousemove="onContainerMouseMove",
    @mouseup="onContainerMouseUp"
  )
    component(
      v-for="shape in shapes",
      v-model="shapes[shape.id]",
      :ref="`shape_${shape.id}`",
      :key="shape.id",
      :is="`${shape.type}-shape`",
      :shape="shape",
      :selected="isSelected(shape)",
      @mousedown="onShapeMouseDown(shape, $event)",
      @mouseup="onShapeMouseUp(shape, $event)",
      @start:resize="onShapeStartResize(shape, $event)"
    )
</template>

<script>
import { omit } from 'lodash';

import Observer from '@/services/observer';

import { roundByStep } from '../utils/round';

import RectShape from './rect-shape.vue';
import LineShape from './line-shape.vue';

const SHAPE_TYPES = {
  rect: 'rect',
  line: 'line',
};

const DIRECTIONS = {
  north: 'n',
  west: 'w',
  south: 's',
  east: 'e',
};

export default {
  provide() {
    return {
      $mouseMove: this.$mouseMove,
      $mouseUp: this.$mouseUp,
    };
  },
  components: { RectShape, LineShape },
  props: {
    movingStep: {
      type: Number,
      default: 5,
    },
  },
  data() {
    return {
      data: {},
      selected: [],
      moving: false,
      cursor: {
        x: 0,
        y: 0,
      },
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
  beforeCreate() {
    this.$mouseMove = new Observer();
    this.$mouseUp = new Observer();
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

    addShape() {
      if (this.resizing || this.moving) {
        return;
      }

      if (this.hasSelected) {
        this.clearSelected();
        return;
      }

      const { width, height } = this.$refs.svg.getBoundingClientRect();
      const types = Object.values(SHAPE_TYPES);
      const index = Math.floor(Math.random() * types.length);
      const type = types[index];
      const centerX = roundByStep(width / 2, this.movingStep);
      const centerY = roundByStep(height / 2, this.movingStep);

      const newShape = type === SHAPE_TYPES.rect
        ? {
          id: Date.now(),
          type,
          width: 100,
          height: 100,
          x: centerX,
          y: centerY,
          style: {
            fill: 'orange',
          },
        }
        : {
          id: Date.now(),
          type,
          x1: centerX - 50,
          y1: centerY,
          x2: centerX + 50,
          y2: centerY,
          style: {
            stroke: 'orange',
            'stroke-width': 1,
          },
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

      this.$mouseUp.notify();
    },

    handleShapeResize(event) {
      const offsetX = roundByStep(event.offsetX, this.movingStep);
      const offsetY = roundByStep(event.offsetY, this.movingStep);
      const shape = this.data[this.resizing.id];

      const directionArray = this.resizing.direction.split('');

      directionArray.forEach((direction, index) => {
        switch (direction) {
          case DIRECTIONS.south: {
            const newHeight = offsetY - shape.y;

            if (newHeight > 0) {
              shape.height = newHeight;
            } else {
              shape.height = Math.abs(newHeight);
              shape.y -= shape.height;

              directionArray[index] = DIRECTIONS.north;
            }

            break;
          }
          case DIRECTIONS.north: {
            const newHeight = shape.height + shape.y - offsetY;

            if (newHeight > 0) {
              shape.height = newHeight;
              shape.y = offsetY;
            } else {
              shape.height = Math.abs(newHeight);
              shape.y = offsetY - shape.height;

              directionArray[index] = DIRECTIONS.south;
            }

            break;
          }
          case DIRECTIONS.east: {
            const newWidth = offsetX - shape.x;

            if (newWidth > 0) {
              shape.width = newWidth;
            } else {
              shape.width = Math.abs(newWidth);
              shape.x = offsetX;

              directionArray[index] = DIRECTIONS.west;
            }

            break;
          }
          case DIRECTIONS.west: {
            const newWidth = shape.width + shape.x - offsetX;

            if (newWidth > 0) {
              shape.width = newWidth;
              shape.x = offsetX;
            } else {
              shape.width = Math.abs(newWidth);
              shape.x = offsetX - shape.width;

              directionArray[index] = DIRECTIONS.east;
            }

            break;
          }
        }
      });

      this.resizing.direction = directionArray.join('');
    },

    handleShapeMove(event) {
      const newMovingOffsetX = roundByStep(event.offsetX - this.movingStart.x, this.movingStep);
      const newMovingOffsetY = roundByStep(event.offsetY - this.movingStart.y, this.movingStep);
      const newMovingOffset = {
        x: newMovingOffsetX,
        y: newMovingOffsetY,
      };

      this.selected.forEach((id) => {
        const shape = this.data[id];

        const [shapeComponent] = this.$refs[`shape_${shape.id}`];

        if (shapeComponent.move) {
          shapeComponent.move(newMovingOffset, this.movingOffset);
        }
      });

      this.movingOffset = newMovingOffset;
    },

    onContainerMouseMove(event) {
      if (this.resizing) {
        this.handleShapeResize(event);
        return;
      }

      if (this.moving) {
        this.handleShapeMove(event);
      }

      const cursor = {
        x: roundByStep(event.offsetX, this.movingStep),
        y: roundByStep(event.offsetY, this.movingStep),
      };

      if (this.cursor.x !== cursor.x || this.cursor.y !== cursor.y) {
        this.cursor = cursor;

        this.$mouseMove.notify(cursor);
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

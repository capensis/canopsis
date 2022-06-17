<template lang="pug">
  svg(
    ref="svg",
    width="100%",
    height="100%",
    @mousemove="onContainerMouseMove",
    @mouseup="onContainerMouseUp",
    @mousedown="onContainerMouseDown"
  )
    component(
      v-for="shape in shapes",
      v-model="shapes[shape.id]",
      :ref="`shape_${shape.id}`",
      :key="shape.id",
      :is="`${shape.type}-shape`",
      :selected="isSelected(shape)",
      :readonly="readonly",
      @mousedown="onShapeMouseDown(shape, $event)",
      @mouseup="onShapeMouseUp(shape, $event)"
    )
</template>

<script>
import { cloneDeep, omit } from 'lodash';

import Observer from '@/services/observer';

import { roundByStep } from '../utils/round';

import RectShape from './rect-shape/rect-shape.vue';
import LineShape from './line-shape/line-shape.vue';
import ArrowLineShape from './arrow-line-shape/arrow-line-shape.vue';
import CircleShape from './circle-shape/circle-shape.vue';
import SquareShape from './square-shape/square-shape.vue';

export default {
  provide() {
    return {
      $mouseMove: this.$mouseMove,
      $mouseUp: this.$mouseUp,
    };
  },
  components: {
    RectShape,
    LineShape,
    ArrowLineShape,
    CircleShape,
    SquareShape,
  },
  model: {
    event: 'input',
    prop: 'shapes',
  },
  props: {
    shapes: {
      type: Object,
      default: () => ({}),
    },
    gridSize: {
      type: Number,
      default: 5,
    },
    readonly: {
      type: Boolean,
      default: false,
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
    hasSelected() {
      return !!this.selected.length;
    },
  },
  watch: {
    shapes: {
      immediate: true,
      handler(value) {
        this.data = cloneDeep(value);
      },
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
    updateShapes(shapes) {
      this.$emit('input', shapes);
    },

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

    getShapeComponentById(shapeId) {
      const [shapeComponent] = this.$refs[`shape_${shapeId}`];

      return shapeComponent;
    },

    onShapeMouseDown(shape, event) {
      if (!this.hasSelected) {
        this.setSelected(shape);
      }

      if (!this.isSelected(shape) && !event.ctrlKey) {
        this.clearSelected();
        this.setSelected(shape);
      }

      if (this.isSelected(shape)) {
        const { offsetX, offsetY } = event;
        this.moving = true;
        this.movingStart = {
          x: roundByStep(offsetX, this.gridSize),
          y: roundByStep(offsetY, this.gridSize),
        };
      }
    },

    onShapeMouseUp(shape, event) {
      if (this.movingOffset.x || this.movingOffset.y) {
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
      if (this.moving) {
        this.moving = false;
        this.movingStart = { x: 0, y: 0 };
        this.movingOffset = { x: 0, y: 0 };
        return;
      }

      this.$mouseUp.notify();
    },

    onContainerMouseDown() {
      this.clearSelected();
    },

    moveSelected(newOffset, oldOffset = { x: 0, y: 0 }) {
      this.selected.forEach((id) => {
        const shapeComponent = this.getShapeComponentById(id);

        if (shapeComponent.move) {
          shapeComponent.move(newOffset, oldOffset);
        }
      });
    },

    handleShapeMove(event) {
      const newMovingOffsetX = roundByStep(
        event.offsetX - this.movingStart.x,
        this.gridSize,
      );
      const newMovingOffsetY = roundByStep(
        event.offsetY - this.movingStart.y,
        this.gridSize,
      );
      const newMovingOffset = {
        x: newMovingOffsetX,
        y: newMovingOffsetY,
      };

      this.moveSelected(newMovingOffset, this.movingOffset);

      this.movingOffset = newMovingOffset;
    },

    onContainerMouseMove(event) {
      if (this.moving) {
        this.handleShapeMove(event);
      }

      const cursor = {
        x: roundByStep(event.offsetX, this.gridSize),
        y: roundByStep(event.offsetY, this.gridSize),
      };

      if (this.cursor.x !== cursor.x || this.cursor.y !== cursor.y) {
        this.cursor = cursor;

        this.$mouseMove.notify(cursor);
      }
    },

    moveSelectedDown() {
      this.moveSelected({ x: 0, y: this.gridSize });
    },

    moveSelectedTop() {
      this.moveSelected({ x: 0, y: -this.gridSize });
    },

    moveSelectedRight() {
      this.moveSelected({ x: this.gridSize, y: 0 });
    },

    moveSelectedLeft() {
      this.moveSelected({ x: -this.gridSize, y: 0 });
    },

    removeSelectedShapes() {
      if (this.hasSelected) {
        this.updateShapes(omit(this.data, this.selected));
        this.clearSelected();
      }
    },

    onKeyDown(event) {
      const handler = {
        37: this.moveSelectedLeft,
        38: this.moveSelectedTop,
        39: this.moveSelectedRight,
        40: this.moveSelectedDown,
        46: this.removeSelectedShapes,
      }[event.keyCode];

      if (handler) {
        event.preventDefault();
        handler();
      }
    },
  },
};
</script>

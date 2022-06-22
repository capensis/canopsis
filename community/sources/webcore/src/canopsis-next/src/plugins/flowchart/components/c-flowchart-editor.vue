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
      v-for="shape in data",
      v-model="data[shape.id]",
      :ref="`shape_${shape.id}`",
      :key="shape.id",
      :is="`${shape.type}-shape`",
      :selected="isSelected(shape)",
      :readonly="readonly",
      :connecting="editing",
      @mousedown="onShapeMouseDown(shape, $event)",
      @mouseup="onShapeMouseUp(shape, $event)",
      @connecting="onConnectMove($event)",
      @connected="onConnectFinish(shape, $event)",
      @unconnect="onUnconnect(shape)",
      @edit:point="startEditPoint(shape, $event)",
      @input="updateConnections"
    )
</template>

<script>
import { cloneDeep, omit } from 'lodash';

import Observer from '@/services/observer';

import { roundByStep } from '../utils/round';

import RectShape from './rect-shape/rect-shape.vue';
import LineShape from './line-shape/line-shape.vue';
import ArrowLineShape from './arrow-line-shape/arrow-line-shape.vue';
import BidirectionalArrowLineShape from './bidirectional-arrow-line-shape/bidirectional-arrow-line-shape.vue';
import CircleShape from './circle-shape/circle-shape.vue';
import EllipseShape from './ellipse-shape/ellipse-shape.vue';
import SquareShape from './square-shape/square-shape.vue';
import ImageShape from './image-shape/image-shape.vue';
import RhombusShape from './rhombus-shape/rhombus-shape.vue';
import ParallelogramShape from './parallelogram-shape/parallelogram-shape.vue';
import StorageShape from './storage-shape/storage-shape.vue';
import { calculateConnectorPointBySide } from '@/plugins/flowchart/utils/connectors';

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
    BidirectionalArrowLineShape,
    CircleShape,
    EllipseShape,
    SquareShape,
    ImageShape,
    RhombusShape,
    ParallelogramShape,
    StorageShape,
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

      cursor: {
        x: 0,
        y: 0,
      },

      editing: false,
      editingShape: false,
      editingPoint: false,

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

    startEditPoint(shape, point) {
      this.editing = true;
      this.editingShape = shape;
      this.editingPoint = point;
    },

    onConnectMove({ x, y }) {
      this.$mouseMove.notify({ x, y });
    },

    onConnectFinish(shape, { side, offset }) {
      const connectingShape = this.data[shape.id];

      connectingShape.connections.push({
        shapeId: this.editingShape.id,
        pointId: this.editingPoint._id,
        offset,
        side,
      });
    },

    onUnconnect(shape) {
      const connectingShape = this.data[shape.id];

      connectingShape.connections = connectingShape.connections.filter(
        connection => connection.shapeId !== this.editingShape.id
        || connection.pointId !== this.editingPoint._id,
      );
    },

    updateConnections(shape) {
      if (shape.connections?.length) {
        shape.connections.forEach(({ shapeId, pointId, offset, side }) => {
          const updatableShape = this.data[shapeId];
          const point = updatableShape.points.find(({ _id: id }) => id === pointId);

          const { x, y } = calculateConnectorPointBySide(shape, side, offset);

          point.x = x;
          point.y = y;
        });
      }
    },

    onContainerMouseUp() {
      if (this.moving) {
        this.moving = false;
        this.movingStart = { x: 0, y: 0 };
        this.movingOffset = { x: 0, y: 0 };
        return;
      }

      if (this.editing) {
        this.editing = false;
        this.editingShape = undefined;
        this.editingPoint = undefined;
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

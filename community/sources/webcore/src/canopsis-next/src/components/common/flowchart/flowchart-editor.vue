<template lang="pug">
  svg.grey.lighten-3(
    ref="svg",
    v-resize="setViewBox",
    :viewBox="viewBoxString",
    width="100%",
    height="100%",
    @mousemove="onContainerMouseMove",
    @mouseup="onContainerMouseUp",
    @mousedown="onContainerMouseDown",
    @contextmenu.stop.prevent=""
  )
    component(
      v-for="shape in data",
      :shape="shape",
      :key="shape._id",
      :is="`${shape.type}-shape`",
      :selected="isSelected(shape._id)",
      :readonly="readonly",
      :connecting="editing",
      @mousedown="onShapeMouseDown(shape, $event)",
      @mouseup="onShapeMouseUp(shape, $event)",
      @connecting="onConnectMove($event)",
      @connected="onConnectFinish(shape, $event)",
      @unconnect="onUnconnect(shape)",
      @edit:point="startEditPoint(shape, $event)",
      @update="updateShape(shape, $event)"
    )
</template>

<script>
import { cloneDeep, isEqual, isObject, isString, omit } from 'lodash';

import Observer from '@/services/observer';

import { SHAPES } from '@/constants';

import { roundByStep } from '@/helpers/flowchart/round';
import { calculateConnectorPointBySide } from '@/helpers/flowchart/connectors';
import uid from '@/helpers/uid';

import { selectedShapesMixin } from '@/mixins/flowchart/selected';

import RectShape from './rect-shape/rect-shape.vue';
import LineShape from './line-shape/line-shape.vue';
import ArrowLineShape from './arrow-line-shape/arrow-line-shape.vue';
import BidirectionalArrowLineShape from './bidirectional-arrow-line-shape/bidirectional-arrow-line-shape.vue';
import CircleShape from './circle-shape/circle-shape.vue';
import EllipseShape from './ellipse-shape/ellipse-shape.vue';
import ImageShape from './image-shape/image-shape.vue';
import RhombusShape from './rhombus-shape/rhombus-shape.vue';
import ParallelogramShape from './parallelogram-shape/parallelogram-shape.vue';
import StorageShape from './storage-shape/storage-shape.vue';
import ProcessShape from './process-shape/process-shape.vue';
import DocumentShape from './document-shape/document-shape.vue';

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
    ImageShape,
    RhombusShape,
    ParallelogramShape,
    StorageShape,
    ProcessShape,
    DocumentShape,
  },
  mixins: [selectedShapesMixin],
  model: {
    event: 'input',
    prop: 'shapes',
  },
  props: {
    shapes: {
      type: Object,
      default: () => ({}),
    },
    viewBox: {
      type: Object,
      required: true,
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
      editorSize: {
        width: 0,
        height: 0,
      },

      data: {},

      cursor: {
        x: 0,
        y: 0,
        shift: false,
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
      viewBoxObject: this.viewBox,

      panning: false,
    };
  },
  computed: {
    viewBoxString() {
      const { x, y, width, height } = this.viewBoxObject;

      return `${x} ${y} ${width} ${height}`;
    },

    widthScale() {
      return this.viewBoxObject.width / this.editorSize.width;
    },

    heightScale() {
      return this.viewBoxObject.height / this.editorSize.height;
    },
  },
  watch: {
    shapes: {
      immediate: true,
      handler(value) {
        if (!isEqual(this.data, value)) {
          this.data = cloneDeep(value);
        }
      },
    },
  },
  beforeCreate() {
    this.$mouseMove = new Observer();
    this.$mouseUp = new Observer();
  },
  mounted() {
    document.addEventListener('keydown', this.onKeyDown);
    this.$refs.svg.addEventListener('wheel', this.onWheel);
  },
  beforeDestroy() {
    document.removeEventListener('keydown', this.onKeyDown);
    this.$refs.svg.removeEventListener('wheel', this.onWheel);
  },
  methods: {
    updateViewBox() {
      this.$emit('update:viewBox', this.viewBoxObject);
    },

    normalizeCursor({ x, y }) {
      const point = this.$refs.svg.createSVGPoint();

      point.x = x;
      point.y = y;

      return point.matrixTransform(this.$refs.svg.getScreenCTM().inverse());
    },

    setViewBox(event) {
      const { width, height } = this.$refs.svg.getBoundingClientRect();

      if (event) {
        const widthDiff = (this.editorSize.width - width) * this.widthScale;
        const heightDiff = (this.editorSize.height - height) * this.heightScale;

        this.viewBoxObject.width -= widthDiff;
        this.viewBoxObject.height -= heightDiff;
      } else {
        this.viewBoxObject.width = width;
        this.viewBoxObject.height = height;
      }

      this.updateViewBox();

      this.editorSize.width = width;
      this.editorSize.height = height;
    },

    onWheel(event) {
      event.preventDefault();

      const delta = event.deltaY;

      if (event.ctrlKey) {
        const percent = delta < 0 ? 0.05 : -0.05;

        const scaleWidth = this.viewBoxObject.width * percent;
        const scaleHeight = this.viewBoxObject.height * percent;

        const deltaWidth = scaleWidth * 2;
        const deltaHeight = scaleHeight * 2;

        const { x, y } = this.normalizeCursor({ x: event.clientX, y: event.clientY });

        const cursorPercentX = (x - this.viewBoxObject.x) / this.viewBoxObject.width;
        const cursorPercentY = (y - this.viewBoxObject.y) / this.viewBoxObject.height;

        const offsetX = deltaWidth * cursorPercentX;
        const offsetY = deltaHeight * cursorPercentY;
        const offsetWidth = scaleWidth + deltaWidth - offsetX;
        const offsetHeight = scaleHeight + deltaHeight - offsetY;

        this.viewBoxObject.x += offsetX;
        this.viewBoxObject.y += offsetY;
        this.viewBoxObject.width -= offsetWidth;
        this.viewBoxObject.height -= offsetHeight;
      }

      if (event.shiftKey) {
        this.viewBoxObject.x += delta;
      }

      if (event.altKey) {
        this.viewBoxObject.y += delta;
      }
    },

    updateShape(shape, data) {
      Object.assign(this.data[shape._id], data);

      this.updateConnections(shape._id);
    },

    updateShapes(shapes) {
      this.$emit('input', shapes);
    },

    onShapeMouseDown(shape, event) {
      if (!this.hasSelected) {
        this.setSelectedShape(shape);
      }

      if (!this.isSelected(shape._id) && !event.ctrlKey) {
        this.clearSelected();
        this.setSelectedShape(shape);
      }

      if (this.isSelected(shape._id)) {
        const { x, y } = this.normalizeCursor({ x: event.clientX, y: event.clientY });

        this.moving = true;
        this.movingStart = {
          x: roundByStep(x, this.gridSize),
          y: roundByStep(y, this.gridSize),
        };
      }
    },

    onShapeMouseUp(shape, event) {
      if (this.movingOffset.x || this.movingOffset.y) {
        return;
      }

      const isShapeSelected = this.isSelected(shape._id);

      if (isShapeSelected && this.selectedIds.length === 1) {
        return;
      }

      if (!event.ctrlKey) {
        this.clearSelected();
        this.setSelectedShape(shape);

        return;
      }

      if (isShapeSelected) {
        this.removeSelectedShape(shape);
      } else {
        this.setSelectedShape(shape);
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

    onConnectFinish(shape, { side }) {
      const connectingShape = this.data[shape._id];

      connectingShape.connections.push({
        shapeId: this.editingShape._id,
        pointId: this.editingPoint._id,
        side,
      });

      const editingShape = this.data[this.editingShape._id];

      editingShape.connectedTo.push(connectingShape._id);
    },

    onUnconnect(shape) {
      const connectingShape = this.data[shape._id];

      connectingShape.connections = connectingShape.connections.filter(
        connection => connection.shapeId !== this.editingShape._id
        || connection.pointId !== this.editingPoint._id,
      );

      const editingShape = this.data[this.editingShape._id];

      editingShape.connectedTo = editingShape.connectedTo
        .filter(shapeId => shapeId !== shape._id);
    },

    updateConnections(id) {
      const shape = this.data[id];

      if (shape.connections?.length) {
        shape.connections.forEach(({ shapeId, pointId, side }) => {
          if (this.isSelected(shapeId)) {
            return;
          }

          const updatingShape = this.data[shapeId];
          const updatingPoint = updatingShape.points.find(point => point._id === pointId);

          const { x, y } = calculateConnectorPointBySide(shape, side);

          updatingPoint.x = x;
          updatingPoint.y = y;
        });
      }
    },

    onContainerMouseUp() {
      if (this.panning) {
        this.updateViewBox();
        this.panning = false;
        return;
      }

      if (this.moving) {
        this.moving = false;
        this.movingStart = { x: 0, y: 0 };
        this.movingOffset = { x: 0, y: 0 };
      }

      if (this.editing) {
        this.editing = false;
        this.editingShape = undefined;
        this.editingPoint = undefined;
      }

      this.$mouseUp.notify();

      this.updateShapes(this.data);
    },

    onContainerMouseDown(event) {
      if (event.ctrlKey || event.shiftKey || event.button === 1) {
        this.panning = true;
        return;
      }

      this.clearSelected();
    },

    onContainerMouseMove(event) {
      if (this.panning) {
        this.viewBoxObject.x -= event.movementX * this.widthScale;
        this.viewBoxObject.y -= event.movementY * this.heightScale;

        return;
      }

      if (this.moving) {
        this.handleShapeMove(event);
      }

      const { x, y } = this.normalizeCursor({ x: event.clientX, y: event.clientY });

      const cursor = {
        x: roundByStep(x, this.gridSize),
        y: roundByStep(y, this.gridSize),
        shift: event.shiftKey,
      };

      if (this.cursor.x !== cursor.x || this.cursor.y !== cursor.y) {
        this.cursor = cursor;

        this.$mouseMove.notify(cursor);
      }
    },

    clearConnectedTo(id) {
      const line = this.data[id];

      if (line.connectedTo) {
        const connectedTo = [];

        line.connectedTo.forEach((shapeId) => {
          const connectedShape = this.data[shapeId];

          if (this.isSelected(connectedShape._id)) {
            connectedTo.push(shapeId);
            return;
          }

          connectedShape.connections = connectedShape.connections.filter(
            connection => connection.shapeId !== id,
          );
        });

        line.connectedTo = connectedTo;
      }
    },

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

    moveSelectedTop() {
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

    removeSelectedShapes() {
      if (this.hasSelected) {
        this.updateShapes(omit(this.data, this.selectedIds));
        this.clearSelected();
      }
    },

    copySelectedShapes(event) {
      if (!event.ctrlKey) {
        return;
      }

      const data = this.selectedIds.reduce((acc, id) => {
        acc[id] = this.data[id];

        return acc;
      }, {});

      navigator.clipboard.writeText(JSON.stringify(data));
    },

    async pasteShapes(event) {
      if (!event.ctrlKey) {
        return;
      }

      const data = await navigator.clipboard.readText();

      if (isString(data)) {
        const shapes = JSON.parse(data);

        if (isObject(shapes)) {
          const preparedShapes = Object.entries(shapes)
            .reduce((acc, [id, shape]) => {
              const resultId = this.data[id] || acc[id] ? `${id}_${uid()}` : id;

              acc[resultId] = {
                ...shape,
                _id: resultId,
              };

              return acc;
            }, {});

          this.data = {
            ...this.data,
            ...preparedShapes,
          };

          this.setSelected(Object.keys(preparedShapes));
          this.updateShapes(this.data);
        }
      }
    },

    onKeyDown(event) {
      const handler = {
        37: this.moveSelectedLeft,
        38: this.moveSelectedTop,
        39: this.moveSelectedRight,
        40: this.moveSelectedDown,

        46: this.removeSelectedShapes,

        67: this.copySelectedShapes,
        86: this.pasteShapes,
      }[event.keyCode];

      if (handler) {
        event.preventDefault();
        handler(event);
      }
    },
  },
};
</script>

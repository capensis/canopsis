<template lang="pug">
  svg(
    ref="svg",
    v-resize.quiet="setViewBox",
    :viewBox="viewBoxString",
    :style="svgStyles",
    width="100%",
    height="100%",
    @mousemove="onContainerMouseMove",
    @mouseup="onContainerMouseUp",
    @mousedown="onContainerMouseDown",
    @contextmenu.stop.prevent="handleContextmenu"
  )
    component(
      v-for="shape in data",
      :key="`${shape._id}-shape`",
      :shape="shape",
      :is="`${shape.type}-shape`",
      :readonly="readonly",
      @contextmenu.stop.prevent="handleShapeContextmenu(shape, $event)",
      @mousedown.left="onShapeMouseDown(shape, $event)",
      @mouseup="onShapeMouseUp(shape, $event)"
    )
    template(v-if="!readonly")
      component(
        v-for="selection in selectionComponents",
        :key="selection.key",
        :is="selection.is",
        :shape="selection.shape",
        :selected="isSelected(selection.shape._id)",
        :connecting="editing",
        @connecting="onConnectMove($event)",
        @connected="onConnectFinish(selection.shape, $event)",
        @unconnect="onUnconnect(selection.shape)",
        @edit:point="startEditLinePoint(selection.shape, $event)",
        @update="updateShape(selection.shape, $event)"
      )
      path(
        v-if="selection",
        :d="selectionPath",
        fill="blue",
        fill-opacity="0.1",
        stroke="blue",
        stroke-width="1"
      )
    slot(name="layers", :data="data")
</template>

<script>
import { cloneDeep, isEqual, omit } from 'lodash';

import { SHAPES } from '@/constants';

import Observer from '@/services/observer';

import { roundByStep } from '@/helpers/flowchart/round';
import { calculateConnectorPointBySide } from '@/helpers/flowchart/connectors';

import { selectedShapesMixin } from '@/mixins/flowchart/selected';
import { copyPasteShapesMixin } from '@/mixins/flowchart/copy-paste';
import { moveShapesMixin } from '@/mixins/flowchart/move-shape';
import { viewBoxMixin } from '@/mixins/flowchart/view-box';
import { contextmenuMixin } from '@/mixins/flowchart/contextmenu';

import RectShape from './rect-shape/rect-shape.vue';
import RhombusShape from './rhombus-shape/rhombus-shape.vue';
import CircleShape from './circle-shape/circle-shape.vue';
import EllipseShape from './ellipse-shape/ellipse-shape.vue';
import ParallelogramShape from './parallelogram-shape/parallelogram-shape.vue';
import ProcessShape from './process-shape/process-shape.vue';
import DocumentShape from './document-shape/document-shape.vue';
import StorageShape from './storage-shape/storage-shape.vue';
import LineShape from './line-shape/line-shape.vue';
import ArrowLineShape from './arrow-line-shape/arrow-line-shape.vue';
import BidirectionalArrowLineShape from './bidirectional-arrow-line-shape/bidirectional-arrow-line-shape.vue';
import ImageShape from './image-shape/image-shape.vue';
import RectShapeSelection from './rect-shape/rect-shape-selection.vue';
import CircleShapeSelection from './circle-shape/circle-shape-selection.vue';
import LineShapeSelection from './line-shape/line-shape-selection.vue';

export default {
  provide() {
    return {
      $mouseMove: this.$mouseMove,
      $mouseUp: this.$mouseUp,
    };
  },
  components: {
    RectShape,
    RhombusShape,
    CircleShape,
    EllipseShape,
    ParallelogramShape,
    ProcessShape,
    DocumentShape,
    StorageShape,
    LineShape,
    ArrowLineShape,
    BidirectionalArrowLineShape,
    ImageShape,
    RectShapeSelection,
    CircleShapeSelection,
    LineShapeSelection,
  },
  mixins: [
    selectedShapesMixin,
    copyPasteShapesMixin,
    moveShapesMixin,
    viewBoxMixin,
    contextmenuMixin,
  ],
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
    backgroundColor: {
      type: String,
      required: false,
    },
    pointSize: {
      type: Number,
      default: 24,
    },
  },
  data() {
    return {
      data: {},

      cursor: {
        x: 0,
        y: 0,
        shift: false,
      },

      editing: false,
      editingShape: false,
      editingLinePoint: false,

      moving: false,
      movingStart: {
        x: 0,
        y: 0,
      },
      movingOffset: {
        x: 0,
        y: 0,
      },

      panning: false,
    };
  },
  computed: {
    selectionComponents() {
      return Object.values(this.data).map((shape) => {
        const component = {
          [SHAPES.rect]: 'rect-shape-selection',
          [SHAPES.rhombus]: 'rect-shape-selection',
          [SHAPES.ellipse]: 'rect-shape-selection',
          [SHAPES.parallelogram]: 'rect-shape-selection',
          [SHAPES.process]: 'rect-shape-selection',
          [SHAPES.document]: 'rect-shape-selection',
          [SHAPES.storage]: 'rect-shape-selection',
          [SHAPES.image]: 'rect-shape-selection',

          [SHAPES.circle]: 'circle-shape-selection',

          [SHAPES.line]: 'line-shape-selection',
          [SHAPES.arrowLine]: 'line-shape-selection',
          [SHAPES.bidirectionalArrowLine]: 'line-shape-selection',
        }[shape.type];

        return {
          shape,
          is: component,
          key: `${shape._id}-shape-selection`,
        };
      });
    },

    svgStyles() {
      return {
        backgroundColor: this.backgroundColor,
      };
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
  },
  beforeDestroy() {
    document.removeEventListener('keydown', this.onKeyDown);
  },
  methods: {
    normalizeCursor({ x, y }) {
      const point = this.$refs.svg.createSVGPoint();

      point.x = x;
      point.y = y;

      return point.matrixTransform(this.$refs.svg.getScreenCTM().inverse());
    },

    updateShape(shape, data) {
      Object.assign(this.data[shape._id], data);

      this.updateConnections(shape._id);
    },

    updateShapes(shapes) {
      this.$emit('input', shapes);
    },

    onShapeMouseDown(shape, event) {
      if (this.readonly) {
        return;
      }

      event.stopPropagation();

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

    startEditLinePoint(shape, point) {
      this.editing = true;
      this.editingShape = shape;
      this.editingLinePoint = point;
    },

    onConnectMove({ x, y }) {
      this.$mouseMove.notify({ x, y });
    },

    onConnectFinish(shape, { side }) {
      const connectingShape = this.data[shape._id];

      connectingShape.connections.push({
        shapeId: this.editingShape._id,
        pointId: this.editingLinePoint._id,
        side,
      });

      const editingShape = this.data[this.editingShape._id];

      editingShape.connectedTo.push(connectingShape._id);
    },

    onUnconnect(shape) {
      const connectingShape = this.data[shape._id];

      connectingShape.connections = connectingShape.connections.filter(
        connection => connection.shapeId !== this.editingShape._id
        || connection.pointId !== this.editingLinePoint._id,
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

    onContainerMouseUp(event) {
      if (this.selection) {
        this.selectShapesByArea(this.selectionStart, this.cursor, event.shiftKey);
        this.selection = false;
        return;
      }

      if (this.panning) {
        this.updateViewBox();
        this.panning = false;
        return;
      }

      if (this.moving) {
        this.moving = false;
        this.movingStart = { x: 0, y: 0 };
        this.movingOffset = { x: 0, y: 0 };

        this.updateShapes(this.data);
      }

      if (this.editing) {
        this.editing = false;
        this.editingShape = undefined;
        this.editingLinePoint = undefined;

        this.updateShapes(this.data);
      }

      this.$mouseUp.notify();
    },

    onContainerMouseDown(event) {
      if (event.ctrlKey || event.button === 1) {
        this.panning = true;
        return;
      }

      if (event.button === 0) {
        this.selection = true;
        this.selectionStart = {
          x: this.cursor.x,
          y: this.cursor.y,
        };
      }
    },

    onContainerMouseMove(event) {
      if (this.panning) {
        this.moveViewBox({ x: event.movementX, y: event.movementY });

        return;
      }

      if (this.readonly) {
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

    removeSelectedShapes() {
      if (this.hasSelected) {
        this.updateShapes(omit(this.data, this.selectedIds));
        this.clearSelected();
      }
    },

    onKeyDown(event) {
      const tag = document.activeElement.tagName.toLowerCase();

      if (['input', 'textarea', 'select'].includes(tag)) {
        return;
      }

      const handlers = {
        37: this.moveSelectedLeft,
        38: this.moveSelectedTop,
        39: this.moveSelectedRight,
        40: this.moveSelectedDown,

        46: this.removeSelectedShapes,
      };
      const ctrlHandler = {
        67: this.copySelectedShapes,
        86: this.pasteShapes,
      };

      const handler = handlers[event.keyCode] || (event.ctrlKey && ctrlHandler[event.keyCode]);

      if (handler) {
        event.preventDefault();
        handler(event);
      }
    },
  },
};
</script>

<template>
  <svg
    ref="svg"
    v-on="svgHandlers"
    :viewBox="viewBoxString"
    :style="svgStyles"
    width="100%"
    height="100%"
    xmlns="http://www.w3.org/2000/svg"
  >
    <component
      v-for="shape in data"
      :key="`${shape._id}-shape`"
      :shape="shape"
      :is="`${shape.type}-shape`"
      :readonly="readonly"
      @click.stop=""
      @contextmenu.stop.prevent="handleShapeContextmenu(shape, $event)"
      @mousedown.left="onShapeMouseDown(shape, $event)"
      @mouseup="onShapeMouseUp(shape, $event)"
      @update="updateShape(shape, $event)"
    />
    <template v-if="!readonly">
      <component
        v-for="selection in selectionComponents"
        :key="selection.key"
        :is="selection.is"
        :shape="selection.shape"
        :selected="isSelected(selection.shape._id)"
        :connecting="editing"
        :color="selectionColor"
        :padding="selectionPadding"
        @connecting="onConnectMove($event)"
        @connected="onConnectFinish(selection.shape, $event)"
        @unconnect="onUnconnect(selection.shape)"
        @edit:point="startEditLinePoint(selection.shape, $event)"
        @update="updateShape(selection.shape, $event)"
      />
      <flowchart-selection
        v-if="selecting"
        :start="selectionStart"
        :end="cursor"
        :color="selectionColor"
      />
    </template>
    <slot
      name="layers"
      :data="data"
    />
  </svg>
</template>

<script>
import { cloneDeep, isEqual, omit } from 'lodash';

import { COLORS } from '@/config';
import { FLOWCHART_KEY_CODES, SHAPES } from '@/constants';

import { roundByStep } from '@/helpers/flowchart/round';
import { calculateConnectorPointBySide } from '@/helpers/flowchart/connectors';

import { selectedShapesMixin } from '@/mixins/flowchart/selected';
import { copyPasteShapesMixin } from '@/mixins/flowchart/copy-paste';
import { moveShapesMixin } from '@/mixins/flowchart/move-shape';
import { viewBoxMixin } from '@/mixins/flowchart/view-box';

import FlowchartSelection from './partials/flowchart-selection.vue';
import RectShape from './shapes/rect-shape/rect-shape.vue';
import RhombusShape from './shapes/rhombus-shape/rhombus-shape.vue';
import CircleShape from './shapes/circle-shape/circle-shape.vue';
import EllipseShape from './shapes/ellipse-shape/ellipse-shape.vue';
import ParallelogramShape from './shapes/parallelogram-shape/parallelogram-shape.vue';
import ProcessShape from './shapes/process-shape/process-shape.vue';
import DocumentShape from './shapes/document-shape/document-shape.vue';
import StorageShape from './shapes/storage-shape/storage-shape.vue';
import LineShape from './shapes/line-shape/line-shape.vue';
import ArrowLineShape from './shapes/arrow-line-shape/arrow-line-shape.vue';
import BidirectionalArrowLineShape from './shapes/bidirectional-arrow-line-shape/bidirectional-arrow-line-shape.vue';
import ImageShape from './shapes/image-shape/image-shape.vue';
import RectShapeSelection from './shapes/rect-shape/rect-shape-selection.vue';
import CircleShapeSelection from './shapes/circle-shape/circle-shape-selection.vue';
import LineShapeSelection from './shapes/line-shape/line-shape-selection.vue';

const DOCUMENT_EVENTS = ['mousemove', 'mouseup', 'keydown'];

export default {
  provide() {
    return {
      $flowchart: this.$flowchart,
    };
  },
  components: {
    FlowchartSelection,
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
    selectionColor: {
      type: String,
      default: COLORS.flowchart.selection,
    },
    selectionPadding: {
      type: Number,
      default: 5,
    },
    cursorStyle: {
      type: String,
      required: false,
    },
  },
  data() {
    return {
      data: {},
      handlers: {},

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
    $flowchart() {
      return {
        on: this.addHandler,
        off: this.removeHandler,
        fire: this.callHandlers,
      };
    },

    svgHandlers() {
      return Object.keys(this.handlers)
        .reduce((acc, event) => {
          if (!this.isDocumentEvent(event)) {
            acc[event] = this.callHandlers;
          }

          return acc;
        }, {});
    },

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

    svgCursor() {
      if (this.panning) {
        return 'move';
      }

      return this.cursorStyle;
    },

    svgStyles() {
      return {
        backgroundColor: this.backgroundColor,
        cursor: this.svgCursor,
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
  mounted() {
    this.$flowchart.on('mousemove', this.onContainerMouseMove);
    this.$flowchart.on('mouseup', this.onContainerMouseUp);
    this.$flowchart.on('mousedown', this.onContainerMouseDown);
    this.$flowchart.on('keydown', this.onKeyDown);
  },
  beforeDestroy() {
    this.$flowchart.off('mousemove', this.onContainerMouseMove);
    this.$flowchart.off('mouseup', this.onContainerMouseUp);
    this.$flowchart.off('mousedown', this.onContainerMouseDown);
    this.$flowchart.off('keydown', this.onKeyDown);
  },
  methods: {
    isDocumentEvent(event) {
      return DOCUMENT_EVENTS.includes(event);
    },

    addHandler(event, func) {
      if (this.handlers[event]) {
        this.handlers[event].push(func);
      } else {
        this.$set(this.handlers, event, [func]);

        if (this.isDocumentEvent(event)) {
          document.addEventListener(event, this.callHandlers);
        }
      }
    },

    removeHandler(event, func) {
      if (this.handlers[event]) {
        const newHandlers = this.handlers[event].filter(handler => handler !== func);

        this.handlers[event] = newHandlers;

        if (!newHandlers.length && this.isDocumentEvent(event)) {
          document.removeEventListener(event, this.callHandlers);
        }
      }
    },

    callHandlers(event, data) {
      this.handlers[event.type]?.forEach(func => func({
        event,
        cursor: this.cursor,
        ...data,
      }));
    },

    handleShapeContextmenu(shape, event) {
      this.$flowchart.fire(event, { event, shape });
    },

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
        this.setSelected([]);
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
        this.setSelected([]);
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
      this.$flowchart.fire({ type: 'mousemove' }, { cursor: { x, y } });
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

    onContainerMouseUp({ event }) {
      if (this.selecting) {
        this.selectShapesByArea(this.selectionStart, this.cursor, event.shiftKey);
        this.selecting = false;
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
      }

      if (this.editing) {
        this.editing = false;
        this.editingShape = undefined;
        this.editingLinePoint = undefined;
      }

      if (!isEqual(this.shapes, this.data)) {
        this.updateShapes(this.data);
      }
    },

    onContainerMouseDown({ event }) {
      if (event.ctrlKey || event.button === 1) {
        this.panning = true;
        return;
      }

      if (event.button === 0) {
        this.selecting = true;
        this.selectionStart = {
          x: this.cursor.x,
          y: this.cursor.y,
        };
      }
    },

    onContainerMouseMove({ event }) {
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

      this.cursor = {
        x: roundByStep(x, this.gridSize),
        y: roundByStep(y, this.gridSize),
        shift: event.shiftKey,
      };
    },

    removeSelectedShapes(event) {
      event.preventDefault();

      if (this.hasSelected) {
        const shapes = cloneDeep(omit(this.data, this.selectedIds));

        Object.values(shapes).forEach((shape) => {
          shapes[shape._id].connections = shape.connections.filter(connection => shapes[connection.shapeId]);
          shapes[shape._id].connectedTo = shape.connectedTo.filter(shapeId => shapes[shapeId]);
        });

        this.updateShapes(shapes);
        this.clearSelected();
      }
    },

    handleCopyShapes(event) {
      if (event.ctrlKey) {
        if (window.getSelection().toString()) {
          return;
        }

        event.preventDefault();

        this.copySelectedShapes();
      }
    },

    handlePasteShapes(event) {
      if (event.ctrlKey) {
        event.preventDefault();

        this.pasteShapes();
      }
    },

    onKeyDown({ event }) {
      const tag = document.activeElement.tagName.toLowerCase();

      if (['input', 'textarea', 'select'].includes(tag)) {
        return;
      }

      const handler = {
        [FLOWCHART_KEY_CODES.arrowUp]: this.moveSelectedUp,
        [FLOWCHART_KEY_CODES.arrowRight]: this.moveSelectedRight,
        [FLOWCHART_KEY_CODES.arrowDown]: this.moveSelectedDown,
        [FLOWCHART_KEY_CODES.arrowLeft]: this.moveSelectedLeft,

        [FLOWCHART_KEY_CODES.delete]: this.removeSelectedShapes,

        [FLOWCHART_KEY_CODES.keyC]: this.handleCopyShapes,
        [FLOWCHART_KEY_CODES.keyV]: this.handlePasteShapes,
      }[event.keyCode];

      if (handler) {
        handler(event);
      }
    },
  },
};
</script>

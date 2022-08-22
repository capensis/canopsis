<template lang="pug">
  div.flowchart-editor
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

    flowchart-contextmenu(
      :value="shownMenu",
      :position-x="pageX",
      :position-y="pageY",
      :items="contextmenuItems",
      @close="closeContextmenu"
    )
    v-menu(
      :value="isDialogOpened",
      :position-x="pageX",
      :position-y="pageY",
      :close-on-content-click="false",
      ignore-click-outside,
      absolute
    )
      point-form-dialog(
        v-if="addingPoint || editingPoint",
        :point="addingPoint || editingPoint",
        :editing="!!editingPoint",
        @cancel="closePointDialog",
        @submit="submitPointDialog",
        @remove=""
      )
</template>

<script>
import { cloneDeep, isEqual, omit } from 'lodash';

import Observer from '@/services/observer';

import { roundByStep } from '@/helpers/flowchart/round';
import { calculateConnectorPointBySide } from '@/helpers/flowchart/connectors';

import { selectedShapesMixin } from '@/mixins/flowchart/selected';
import { copyPasteShapesMixin } from '@/mixins/flowchart/copy-paste';
import { moveShapesMixin } from '@/mixins/flowchart/move-shape';
import { viewBoxMixin } from '@/mixins/flowchart/view-box';
import { pointsMixin } from '@/mixins/flowchart/points';

import PointFormDialog from '@/components/other/map/form/partials/point-form-dialog.vue';

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
import FlowchartContextmenu from './partials/flowchart-contextmenu.vue';

export default {
  provide() {
    return {
      $mouseMove: this.$mouseMove,
      $mouseUp: this.$mouseUp,
    };
  },
  components: {
    FlowchartContextmenu,
    PointFormDialog,
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
  mixins: [
    selectedShapesMixin,
    copyPasteShapesMixin,
    moveShapesMixin,
    viewBoxMixin,
    pointsMixin,
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

      panning: false,
    };
  },
  computed: {
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
        this.moveViewBox({ x: event.movementX, y: event.movementY });

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

<style lang="scss">
.flowchart-editor {
  height: 100%;
  width: 100%;
}
</style>

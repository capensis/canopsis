<template>
  <g class="flowchart-points-editor">
    <component
      is="foreignObject"
      v-for="{ point, x, y } in nonShapesIcons"
      :key="point._id"
      :x="x"
      :y="y"
      :width="iconSize"
      :height="iconSize"
      class="flowchart-points-editor__point"
      @click.stop=""
      @contextmenu.stop.prevent="handleEditContextmenu($event, point)"
      @dblclick.stop="openEditPointByClick($event, point)"
      @mousedown.stop="startMoving(point)"
    >
      <point-icon
        :size="iconSize"
        :entity="point.entity"
      />
    </component>
    <component
      is="foreignObject"
      v-for="{ point, x, y } in shapesIcons"
      :key="point._id"
      :height="iconSize"
      :width="iconSize"
      :x="x"
      :y="y"
      class="flowchart-points-editor__point"
      @mouseup.prevent.stop=""
      @mousedown.prevent.stop=""
      @click.stop=""
      @dblclick.stop="openEditPointByClick($event, point)"
      @contextmenu.stop.prevent="handleEditContextmenu($event, point)"
    >
      <point-icon
        :size="iconSize"
        :entity="point.entity"
      />
    </component>
    <component
      is="foreignObject"
      style="overflow: visible;"
    >
      <point-form-dialog-menu
        :value="isDialogOpened"
        :position-x="clientX"
        :position-y="clientY"
        :point="addingPoint || editingPoint"
        :editing="!!editingPoint"
        :exists-entities="existsEntities"
        @cancel="cancelPointDialog"
        @submit="submitPointDialog"
        @remove="showRemovePointModal"
      />
      <flowchart-point-contextmenu
        :value="shownMenu"
        :position-x="clientX"
        :position-y="clientY"
        :items="contextmenuItems"
        @close="cancelPointDialog"
      />
    </component>
  </g>
</template>

<script>
import { cloneDeep } from 'lodash';

import { FLOWCHART_MAX_POSITION_DIFF, FLOWCHART_MAX_TIMESTAMP_DIFF, MODALS } from '@/constants';

import { flowchartPointToForm } from '@/helpers/entities/map/form';
import { waitVuetifyAnimation } from '@/helpers/vuetify';

import { formMixin } from '@/mixins/form';
import { mapFlowchartPointsMixin } from '@/mixins/map/map-flowchart-points-mixin';

import PointIcon from '@/components/other/map/partials/point-icon.vue';

import PointFormDialogMenu from './point-form-dialog-menu.vue';
import FlowchartPointContextmenu from './flowchart-point-contextmenu.vue';

export default {
  inject: ['$flowchart'],
  components: {
    PointFormDialogMenu,
    FlowchartPointContextmenu,
    PointIcon,
  },
  mixins: [formMixin, mapFlowchartPointsMixin],
  model: {
    prop: 'points',
    event: 'input',
  },
  props: {
    iconSize: {
      type: Number,
      default: 24,
    },
    addOnClick: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      pointsData: [],
      movingPointIndex: undefined,
      shownMenu: false,
      shownPointDialog: false,
      clientX: 0,
      clientY: 0,
      pointX: 0,
      pointY: 0,
      shapeId: undefined,
      addingPoint: undefined,
      editingPoint: undefined,
    };
  },
  computed: {
    existsEntities() {
      return this.points.map(({ entity }) => entity);
    },

    isDialogOpened() {
      return this.shownPointDialog;
    },

    contextmenuItems() {
      if (this.editingPoint) {
        return [
          {
            text: this.$t('map.editPoint'),
            action: this.openEditPointDialog,
          },
          {
            text: this.$t('map.removePoint'),
            action: this.showRemovePointModal,
          },
        ];
      }

      return [
        {
          text: this.$t('map.addPoint'),
          action: this.openAddPointDialog,
        },
      ];
    },
  },
  watch: {
    addOnClick: {
      immediate: true,
      handler(value) {
        if (value) {
          this.$flowchart.on('mousedown', this.handleMouseDown);
          this.$flowchart.on('mouseup', this.handleMouseUp);
        } else {
          this.$flowchart.off('mousedown', this.handleMouseDown);
          this.$flowchart.off('mouseup', this.handleMouseUp);
        }
      },
    },
    points: {
      immediate: true,
      handler(value) {
        this.pointsData = cloneDeep(value);
      },
    },
  },
  mounted() {
    this.$flowchart.on('dblclick', this.openAddPointDialogByClick);
    this.$flowchart.on('contextmenu', this.handleContextmenu);
  },
  beforeDestroy() {
    this.$flowchart.off('contextmenu', this.handleContextmenu);
    this.$flowchart.off('dblclick', this.openAddPointDialogByClick);
    this.$flowchart.off('mousedown', this.handleMouseDown);
    this.$flowchart.off('mouseup', this.handleMouseUp);
  },
  methods: {
    updatePointInModel(data) {
      this.updateModel(
        this.points.map(point => (point._id === data._id ? data : point)),
      );
    },

    addPointToModel(data) {
      this.updateModel([...this.points, data]);
    },

    removePointFromModel(data) {
      this.updateModel(
        this.points.filter(point => data._id !== point._id),
      );
    },

    openContextmenu() {
      this.shownMenu = true;
    },

    closeContextmenu() {
      this.shownMenu = false;
    },

    openPointDialog() {
      this.shownPointDialog = true;
    },

    closePointDialog() {
      this.shownPointDialog = false;
    },

    clearPointDialog() {
      this.addingPoint = undefined;
      this.editingPoint = undefined;
      this.shapeId = undefined;
    },

    async cancelPointDialog() {
      this.closeContextmenu();
      this.closePointDialog();

      await waitVuetifyAnimation();

      this.clearPointDialog();
    },

    setOffsetsByEvent(event) {
      this.clientX = event.clientX;
      this.clientY = event.clientY;
    },

    setOffsetsByPointEvent(event) {
      const { x, y, width, height } = event.target.getBoundingClientRect();

      this.clientX = x + width / 2;
      this.clientY = y + height;
    },

    handleContextmenu({ event, cursor, shape }) {
      event.stopPropagation();
      event.preventDefault();

      if (this.shownMenu || this.isDialogOpened) {
        return;
      }

      if (shape) {
        const editingPoint = this.points.find(point => point.shape === shape._id);

        if (editingPoint) {
          this.editingPoint = editingPoint;
        } else {
          this.shapeId = shape._id;
        }
      } else {
        this.pointX = cursor.x;
        this.pointY = cursor.y;
      }

      this.setOffsetsByEvent(event);
      this.openContextmenu();
    },

    handleEditContextmenu(event, point) {
      if (this.shownMenu || this.isDialogOpened) {
        return;
      }

      this.setOffsetsByPointEvent(event);
      this.openContextmenu();

      this.editingPoint = point;
    },

    openAddPointDialog() {
      this.addingPoint = flowchartPointToForm(
        this.shapeId
          ? { shape: this.shapeId }
          : { x: this.pointX, y: this.pointY },
      );

      this.closeContextmenu();
      this.openPointDialog();
    },

    handleMouseDown({ event, cursor }) {
      if (event.buttons !== 1) {
        return;
      }

      this.mouseDownTimestamp = event.timeStamp;
      this.mouseDownCursorX = cursor.x;
      this.mouseDownCursorY = cursor.y;
    },

    handleMouseUp({ event, cursor }) {
      if (
        Math.abs(cursor.x - this.mouseDownCursorX) < FLOWCHART_MAX_POSITION_DIFF
        && Math.abs(cursor.y - this.mouseDownCursorY) < FLOWCHART_MAX_POSITION_DIFF
        && (event.timeStamp - this.mouseDownTimestamp) < FLOWCHART_MAX_TIMESTAMP_DIFF
      ) {
        this.openAddPointDialogByClick({
          event,
          cursor,
        });
      }
    },

    openAddPointDialogByClick({ event, cursor }) {
      if (this.isDialogOpened) {
        return;
      }

      this.setOffsetsByEvent(event);

      this.pointX = cursor.x;
      this.pointY = cursor.y;

      this.openAddPointDialog();
    },

    openEditPointDialog() {
      this.closeContextmenu();
      this.openPointDialog();
    },

    openEditPointByClick(event, point) {
      if (this.isDialogOpened) {
        return;
      }

      this.setOffsetsByEvent(event);

      this.editingPoint = point;

      this.openEditPointDialog();
    },

    submitPointDialog(data) {
      if (this.editingPoint) {
        this.updatePointInModel(data);
      } else {
        this.addPointToModel(data);
      }

      this.closePointDialog();
      this.clearPointDialog();
    },

    showRemovePointModal() {
      this.closeContextmenu();

      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: () => {
            this.removePointFromModel(this.editingPoint);

            this.closePointDialog();
            this.clearPointDialog();
          },
        },
      });
    },

    movePoint({ cursor }) {
      this.pointsData[this.movingPointIndex].x = cursor.x;
      this.pointsData[this.movingPointIndex].y = cursor.y;
    },

    finishMoving() {
      this.updateModel(this.pointsData);

      this.$flowchart.off('mousemove', this.movePoint);
      this.$flowchart.off('mouseup', this.finishMoving);
    },

    startMoving(point) {
      this.movingPointIndex = this.pointsData.findIndex(({ _id: id }) => point._id === id);

      this.$flowchart.on('mousemove', this.movePoint);
      this.$flowchart.on('mouseup', this.finishMoving);
    },
  },
};
</script>

<style lang="scss">
.flowchart-points-editor {
  &__point {
    user-select: none;
    cursor: pointer;
  }
}
</style>

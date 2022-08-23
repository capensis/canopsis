<template lang="pug">
  g.flowchart-points-editor
    component.flowchart-points-editor__point(
      v-for="(point, index) in pointsData",
      :key="point._id",
      :x="point.x - iconSize / 2",
      :y="point.y - iconSize",
      :width="iconSize",
      :height="iconSize",
      is="foreignObject",
      @contextmenu.stop.prevent="handleEditContextmenu($event, point)",
      @mousedown.stop="startMoving(index)"
    )
      point-icon(:size="iconSize", :entity="point.entity")
    component(is="foreignObject", style="overflow: visible;")
      flowchart-point-dialog-menu(
        v-if="isDialogOpened",
        :position-x="clientX",
        :position-y="clientY",
        :point="addingPoint || editingPoint",
        :editing="!!editingPoint",
        @cancel="closePointDialog",
        @submit="submitPointDialog",
        @remove="showRemovePointModal"
      )
      flowchart-point-contextmenu(
        :value="shownMenu",
        :position-x="clientX",
        :position-y="clientY",
        :items="contextmenuItems",
        @close="closeContextmenu"
      )
</template>

<script>
import { cloneDeep } from 'lodash';

import { MODALS } from '@/constants';

import { flowchartPointToForm } from '@/helpers/forms/map';

import { formMixin } from '@/mixins/form';

import PointIcon from '@/components/other/map/partials/point-icon.vue';
import PointFormDialog from '@/components/other/map/form/partials/point-form-dialog.vue';

import FlowchartPointDialogMenu from './flowchart-point-dialog-menu.vue';
import FlowchartPointContextmenu from './flowchart-point-contextmenu.vue';

export default {
  inject: ['$contextmenu', '$mouseMove', '$mouseUp'],
  components: { FlowchartPointDialogMenu, FlowchartPointContextmenu, PointFormDialog, PointIcon },
  mixins: [formMixin],
  model: {
    prop: 'points',
    event: 'input',
  },
  props: {
    points: {
      type: Array,
      required: true,
    },
    iconSize: {
      type: Number,
      default: 24,
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
      addingPoint: undefined,
      editingPoint: undefined,
    };
  },
  computed: {
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
    points: {
      immediate: true,
      handler(value) {
        this.pointsData = cloneDeep(value);
      },
    },
  },
  created() {
    this.$contextmenu.register(this.handleContextmenu);
  },
  beforeDestroy() {
    this.$contextmenu.unregister(this.handleContextmenu);
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
      this.addingPoint = undefined;
      this.editingPoint = undefined;
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

    handleContextmenu({ event, x, y }) {
      if (this.shownMenu || this.isDialogOpened) {
        return;
      }

      this.pointX = x;
      this.pointY = y;

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
      this.addingPoint = flowchartPointToForm({
        x: this.pointX,
        y: this.pointY,
      });

      this.closeContextmenu();
      this.openPointDialog();
    },

    openEditPointDialog() {
      this.closeContextmenu();
      this.openPointDialog();
    },

    submitPointDialog(data) {
      if (this.editingPoint) {
        this.updatePointInModel(data);
      } else {
        this.addPointToModel(data);
      }

      this.closePointDialog();
    },

    showRemovePointModal() {
      this.closeContextmenu();

      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: () => {
            this.removePointFromModel(this.editingPoint);

            this.closePointDialog();
          },
        },
      });
    },

    movePoint(cursor) {
      this.pointsData[this.movingPointIndex].x = cursor.x;
      this.pointsData[this.movingPointIndex].y = cursor.y;
    },

    finishMoving() {
      this.updateModel(this.pointsData);

      this.$mouseMove.unregister(this.movePoint);
      this.$mouseUp.unregister(this.finishMoving);
    },

    startMoving(index) {
      this.movingPointIndex = index;

      this.$mouseMove.register(this.movePoint);
      this.$mouseUp.register(this.finishMoving);
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

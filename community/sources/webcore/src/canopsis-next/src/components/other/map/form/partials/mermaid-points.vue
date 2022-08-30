<template lang="pug">
  div.mermaid-points(
    ref="container",
    :class="{ 'mermaid-points--addable': addOnClick && !moving && !isFormOpened }",
    @contextmenu.stop.prevent="openContextmenu",
    @click="openAddPointFormByClick",
    @dblclick="openAddPointFormByDoubleClick"
  )
    mermaid-point-marker.mermaid-points__point(
      v-for="(point, index) in pointsData",
      :key="point._id",
      :x="point.x",
      :y="point.y",
      :entity="point.entity",
      :size="markerSize",
      :class="{ 'mermaid-points__point--no-events': moving }",
      @click.stop="",
      @contextmenu.stop.prevent="openEditContextmenu(point)",
      @dblclick.stop="openEditPointFormByDoubleClick(point)",
      @mousedown.stop.prevent.left="startMoving(point, index)"
    )

    v-menu(
      :value="shownMenu",
      :position-x="pageX",
      :position-y="pageY",
      :close-on-content-click="false",
      absolute,
      @input="clearMenuData"
    )
      mermaid-contextmenu(
        :editing="!!editingPoint",
        @add:point="openAddPointForm",
        @edit:point="openEditPointForm",
        @remove:point="showRemovePointModal"
      )
    v-menu(
      :value="isFormOpened",
      :position-x="pageX",
      :position-y="pageY",
      :close-on-content-click="false",
      ignore-click-outside,
      absolute
    )
      point-form-dialog(
        v-if="isFormOpened",
        :point="formPoint",
        :editing="!!editingPoint",
        @cancel="clearMenuData",
        @submit="submitPointForm",
        @remove="showRemovePointModal"
      )
</template>

<script>
import { cloneDeep } from 'lodash';

import { MODALS } from '@/constants';

import { waitVuetifyAnimation } from '@/helpers/vuetify';
import { mermaidPointToForm } from '@/helpers/forms/map';

import { formBaseMixin } from '@/mixins/form';

import MermaidPointMarker from '@/components/other/map/partials/mermaid-point-marker.vue';

import MermaidContextmenu from './mermaid-contextmenu.vue';
import PointFormDialog from './point-form-dialog.vue';

export default {
  components: { MermaidContextmenu, PointFormDialog, MermaidPointMarker },
  mixins: [formBaseMixin],
  model: {
    prop: 'points',
    event: 'input',
  },
  props: {
    points: {
      type: Array,
      required: true,
    },
    addOnClick: {
      type: Boolean,
      default: false,
    },
    markerSize: {
      type: Number,
      default: 24,
    },
  },
  data() {
    return {
      pointsData: [],

      shownMenu: false,

      editing: false,
      adding: false,
      moving: false,

      addingPoint: undefined,
      editingPoint: undefined,
      movingPointIndex: undefined,

      offsetX: 0,
      offsetY: 0,
      pageX: 0,
      pageY: 0,
    };
  },
  computed: {
    isFormOpened() {
      return this.adding || this.editing;
    },

    formPoint() {
      return this.addingPoint || this.editingPoint;
    },
  },
  watch: {
    points: {
      immediate: true,
      handler() {
        this.pointsData = cloneDeep(this.points);
      },
    },
  },
  methods: {
    normalizePosition({ x, y }) {
      return {
        x,
        y: Math.max(this.markerSize, y),
      };
    },

    setOffsetsByEvent(event) {
      this.pageX = event.pageX;
      this.pageY = event.pageY - window.scrollY;
      this.offsetX = event.offsetX;
      this.offsetY = event.offsetY;
    },

    setOffsetsByPoint(point) {
      const { x, y } = this.$refs.container.getBoundingClientRect();

      this.pageX = x + point.x;
      this.pageY = y + point.y;
      this.offsetX = point.x;
      this.offsetY = point.y;
    },

    async clearMenuData() {
      this.shownMenu = false;
      this.editing = false;
      this.adding = false;

      await waitVuetifyAnimation();

      this.editingPoint = undefined;
      this.addingPoint = undefined;
    },

    openContextmenu(event) {
      if (this.shownMenu || this.isFormOpened) {
        return;
      }

      this.setOffsetsByEvent(event);

      this.shownMenu = true;
    },

    openEditContextmenu(point) {
      if (this.shownMenu || this.isFormOpened) {
        return;
      }

      this.setOffsetsByPoint(point);

      this.shownMenu = true;
      this.editingPoint = point;
    },

    openAddPointForm() {
      if (this.isFormOpened) {
        return;
      }

      this.addingPoint = mermaidPointToForm(
        this.normalizePosition({
          x: this.offsetX,
          y: this.offsetY,
        }),
      );

      this.shownMenu = false;
      this.adding = true;
    },

    openAddPointFormByClick(event) {
      if (!this.addOnClick || this.shownMenu || this.isFormOpened || this.moving) {
        return;
      }

      this.setOffsetsByEvent(event);
      this.openAddPointForm();
    },

    openAddPointFormByDoubleClick(event) {
      if (this.isFormOpened || this.moving) {
        return;
      }

      this.setOffsetsByEvent(event);
      this.openAddPointForm();
    },

    openEditPointFormByDoubleClick(point) {
      this.setOffsetsByPoint(point);

      this.editingPoint = point;
      this.editing = true;
    },

    openEditPointForm() {
      this.shownMenu = false;
      this.editing = true;
    },

    submitPointForm(data) {
      if (this.editing) {
        this.updateModel(this.points.map(point => (point._id === data._id ? data : point)));
      } else if (this.adding) {
        this.updateModel([...this.points, data]);
      }

      this.clearMenuData();
    },

    showRemovePointModal() {
      this.shownMenu = false;

      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: () => {
            this.updateModel(this.points.filter(point => this.editingPoint._id !== point._id));

            this.clearMenuData();
          },
        },
      });
    },

    movePointByEvent(event) {
      const { x, y } = this.normalizePosition({
        x: event.offsetX,
        y: event.offsetY,
      });

      this.pointsData[this.movingPointIndex].x = x;
      this.pointsData[this.movingPointIndex].y = y;
    },

    startMoving(point, index) {
      if (this.isFormOpened) {
        return;
      }

      this.movingTimer = setTimeout(() => {
        this.moving = true;
        this.movingPointIndex = index;

        this.$refs.container.addEventListener('mousemove', this.movePointByEvent);
      }, 200);

      document.addEventListener('mouseup', this.finishMoving);
    },

    finishMoving() {
      clearTimeout(this.movingTimer);

      setTimeout(() => {
        this.moving = false;
        this.movingPointIndex = undefined;
      });

      this.$refs.container.removeEventListener('mousemove', this.movePointByEvent);
      document.removeEventListener('mouseup', this.finishMoving);

      this.updateModel(this.pointsData);
    },
  },
};
</script>

<style lang="scss">
.mermaid-points {
  &__point {
    &--no-events {
      pointer-events: none;
    }
  }

  &--addable {
    cursor: crosshair;
  }
}
</style>

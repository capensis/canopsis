<template lang="pug">
  div.mermaid-points(
    ref="container",
    :class="{ 'mermaid-points--addable': addOnClick && !isFormOpened }",
    @contextmenu.stop.prevent="openContextmenu",
    @click="openAddPointFormByClick",
    @dblclick="openAddPointFormByDoubleClick"
  )
    v-icon.mermaid-points__point(
      v-for="point in points",
      :key="point._id",
      :style="getPointStyles(point)",
      @contextmenu.stop.prevent="openEditContextmenu(point)",
      @dblclick.stop="openEditPointFormByDoubleClick(point)",
      @click.stop=""
    ) location_on

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
      :value="editing || adding",
      :position-x="pageX",
      :position-y="pageY",
      :close-on-content-click="false",
      ignore-click-outside,
      absolute
    )
      mermaid-point-form(
        v-if="editingPoint || addingPoint",
        :point="editingPoint || addingPoint",
        :editing="!!editingPoint",
        @cancel="clearMenuData",
        @submit="submitPointForm",
        @remove="showRemovePointModal"
      )
</template>

<script>
import { MODALS } from '@/constants';

import { formBaseMixin } from '@/mixins/form';

import { waitVuetifyAnimation } from '@/helpers/vuetify';

import MermaidContextmenu from './partials/mermaid-contextmenu.vue';
import MermaidPointForm from './partials/mermaid-point-form.vue';

export default {
  components: { MermaidPointForm, MermaidContextmenu },
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
  },
  data() {
    return {
      shownMenu: false,

      editing: false,
      adding: false,

      addingPoint: undefined,
      editingPoint: undefined,

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
  },
  methods: {
    getPointStyles(point) {
      return {
        top: `${point.y}px`,
        left: `${point.x}px`,
      };
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

      this.pageX = event.pageX;
      this.pageY = event.pageY;
      this.offsetX = event.offsetX;
      this.offsetY = event.offsetY;

      this.shownMenu = true;
    },

    setOffsetsByEvent(event) {
      this.pageX = event.pageX;
      this.pageY = event.pageY;
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

      this.addingPoint = {
        x: this.offsetX,
        y: this.offsetY,
        _id: Date.now(),
      };

      this.shownMenu = false;
      this.adding = true;
    },

    openAddPointFormByClick(event) {
      if (!this.addOnClick || this.shownMenu || this.isFormOpened) {
        return;
      }

      this.setOffsetsByEvent(event);
      this.openAddPointForm();
    },

    openAddPointFormByDoubleClick(event) {
      if (this.isFormOpened) {
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
  },
};
</script>

<style lang="scss">
.mermaid-points {
  &__point {
    position: absolute;
    transform: translate(-50%, -100%);
    user-select: none;
    cursor: pointer;
  }

  &--addable {
    cursor: crosshair;
  }
}
</style>

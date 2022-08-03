<template lang="pug">
  div.mermaid-points(
    ref="container",
    @contextmenu.stop.prevent="openContextmenu",
    @dblclick="openAddPointFormByDoubleClick"
  )
    v-icon.mermaid-points__point(
      v-for="point in points",
      :key="point._id",
      :style="getPointStyles(point)",
      @contextmenu.stop.prevent="openEditContextmenu(point)",
      @click.stop.prevent="openEditContextmenu(point)"
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
        :removable="editing",
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
  methods: {
    getPointStyles(point) {
      return {
        top: `${point.y * 100}%`,
        left: `${point.x * 100}%`,
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
      if (this.shownMenu || this.adding || this.editing) {
        return;
      }

      this.pageX = event.pageX;
      this.pageY = event.pageY;
      this.offsetX = event.offsetX;
      this.offsetY = event.offsetY;

      this.shownMenu = true;
    },

    openEditContextmenu(point) {
      if (this.shownMenu || this.adding || this.editing) {
        return;
      }

      const { x, y, width, height } = this.$refs.container.getBoundingClientRect();

      const offsetX = point.x * width;
      const offsetY = point.y * height;

      this.pageX = x + offsetX;
      this.pageY = y + offsetY;
      this.offsetX = offsetX;
      this.offsetY = offsetY;

      this.shownMenu = true;
      this.editingPoint = point;
    },

    openAddPointForm() {
      if (this.adding || this.editing) {
        return;
      }

      const { width, height } = this.$refs.container.getBoundingClientRect();

      this.addingPoint = {
        x: this.offsetX / width,
        y: this.offsetY / height,
        _id: Date.now(),
      };

      this.shownMenu = false;
      this.adding = true;
    },

    openAddPointFormByDoubleClick(event) {
      if (this.adding || this.editing) {
        return;
      }

      this.pageX = event.pageX;
      this.pageY = event.pageY;
      this.offsetX = event.offsetX;
      this.offsetY = event.offsetY;

      this.openAddPointForm();
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
}
</style>

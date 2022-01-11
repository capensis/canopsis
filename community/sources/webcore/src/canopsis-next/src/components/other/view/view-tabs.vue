<template lang="pug">
  v-tabs.view-tabs(
    ref="tabs",
    :key="vTabsKey",
    :value="$route.fullPath",
    :class="{ hidden: this.tabs.length < 2 && !editing, 'tabs-editing': editing }",
    :hide-slider="isTabsChanged",
    color="secondary lighten-2",
    slider-color="primary",
    dark
  )
    draggable.d-flex(
      v-if="tabs.length",
      :value="tabs",
      :options="draggableOptions",
      @end="onDragEnd",
      @input="$emit('update:tabs', $event)"
    )
      v-tab.draggable-item(
        v-for="tab in tabs",
        :key="tab._id",
        :disabled="isTabsChanged",
        :to="getTabHrefById(tab._id)",
        exact,
        ripple
      )
        span {{ tab.title }}
        update-tab-btn(
          v-show="updatable && editing",
          :tab="tab",
          :updateTabMethod="updateTab"
        )
        clone-tab-btn(
          v-show="updatable && editing",
          :tab="tab"
        )
        delete-tab-btn(
          v-show="updatable && editing",
          :tab="tab",
          :view="view",
          :updateViewMethod="updateViewMethod"
        )
    template(v-if="$scopedSlots.default")
      v-tabs-items(touchless)
        v-tab-item(
          v-for="tab in tabs",
          :key="tab._id",
          :value="getTabHrefById(tab._id)",
          lazy
        )
          slot(
            :tab="tab",
            :editing="editing",
            :updateTabMethod="updateTab"
          )
</template>

<script>
import Draggable from 'vuedraggable';

import { VUETIFY_ANIMATION_DELAY } from '@/config';

import vuetifyTabsMixin from '@/mixins/vuetify/tabs';

import UpdateTabBtn from './buttons/update-tab-btn.vue';
import CloneTabBtn from './buttons/clone-tab-btn.vue';
import DeleteTabBtn from './buttons/delete-tab-btn.vue';

export default {
  components: {
    Draggable,
    UpdateTabBtn,
    CloneTabBtn,
    DeleteTabBtn,
  },
  mixins: [
    vuetifyTabsMixin,
  ],
  props: {
    view: {
      type: Object,
      required: true,
    },
    tabs: {
      type: Array,
      required: true,
    },
    updatable: {
      type: Boolean,
      default: false,
    },
    isTabsChanged: {
      type: Boolean,
      default: false,
    },
    editing: {
      type: Boolean,
      default: false,
    },
    updateViewMethod: {
      type: Function,
      default: () => {},
    },
  },
  computed: {
    vTabsKey() {
      return this.view.tabs.map(tab => tab._id).join('-');
    },
    draggableOptions() {
      return {
        animation: VUETIFY_ANIMATION_DELAY,
        disabled: !this.editing,
      };
    },
    getTabHrefById() {
      return (id) => {
        const { href } = this.$router.resolve({ query: { tabId: id } }, this.$route);

        return href.replace('#', '');
      };
    },
  },
  watch: {
    editing() {
      this.$nextTick(this.callTabsOnResizeMethod);
    },
    tabs: {
      immediate: true,
      handler() {
        this.onUpdateTabs();
      },
    },
  },
  methods: {
    updateTab(tab) {
      const view = {
        ...this.view,
        tabs: this.view.tabs.map((viewTab) => {
          if (viewTab._id === tab._id) {
            return tab;
          }

          return viewTab;
        }),
      };

      return this.updateViewMethod(view);
    },

    onUpdateTabs() {
      this.$nextTick(() => {
        this.callTabsOnResizeMethod();
        this.callTabsUpdateTabsMethod();
      });
    },

    onDragEnd() {
      this.onUpdateTabs();
    },
  },
};
</script>

<style lang="scss" scoped>
  .view-tabs.hidden {
    & /deep/ > .v-tabs__bar {
      display: none;
    }
  }

  .draggable-item {
    position: relative;
    transform: translateZ(0);

    .tabs-editing & {
      cursor: move;

      & /deep/ .v-tabs__item {
        cursor: move;
      }
    }

    & /deep/ .v-tabs__item--disabled {
      color: #fff;
      opacity: 1;

      button {
        color: rgba(255, 255, 255, 0.3) !important;
        box-shadow: none !important;
        pointer-events: none;
      }
    }
  }
</style>

<template lang="pug">
  div
    v-tabs.view-tabs(
    ref="tabs",
    :key="vTabsKey",
    :value="value",
    :class="{ hidden: this.tabs.length < 2 }",
    :hide-slider="hideSlider",
    color="secondary lighten-2",
    slider-color="primary",
    dark,
    @change="$emit('input', $event)"
    )
      draggable.d-flex(
      :value="tabs",
      :options="draggableOptions",
      @end="onDragEnd",
      @input="$emit('update:tabs', $event)"
      )
        v-tab.draggable-item(v-if="tabs.length", v-for="tab in tabs", :key="tab._id", ripple)
          span {{ tab.title }}
          v-btn(
          v-show="hasUpdateAccess && isEditingMode",
          :disabled="isTabsChanged",
          small,
          flat,
          icon,
          @click.stop="showUpdateTabModal(tab)"
          )
            v-icon(small) edit
          v-btn(
          v-show="hasUpdateAccess && isEditingMode",
          :disabled="isTabsChanged",
          small,
          flat,
          icon,
          @click.stop="showDeleteTabModal(tab)"
          )
            v-icon(small) delete
      v-tabs-items(v-if="$scopedSlots.default", active-class="active-view-tab")
        v-tab-item(v-for="tab in tabs", :key="tab._id", lazy)
          slot(
          :tab="tab",
          :isEditingMode="isEditingMode",
          :hasUpdateAccess="hasUpdateAccess",
          :updateTabMethod="updateTab"
          )
</template>

<script>
import Draggable from 'vuedraggable';

import { VUETIFY_ANIMATION_DELAY } from '@/config';
import { MODALS } from '@/constants';

import modalMixin from '@/mixins/modal';
import vuetifyTabsMixin from '@/mixins/vuetify/tabs';

export default {
  components: { Draggable },
  mixins: [modalMixin, vuetifyTabsMixin],
  props: {
    view: {
      type: Object,
      required: true,
    },
    tabs: {
      type: Array,
      required: true,
    },
    value: {
      type: Number,
      default: null,
    },
    hasUpdateAccess: {
      type: Boolean,
      default: false,
    },
    isTabsChanged: {
      type: Boolean,
      default: false,
    },
    isEditingMode: {
      type: Boolean,
      default: false,
    },
    updateViewMethod: {
      type: Function,
      required: true,
    },
    hideSlider: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    vTabsKey() {
      return this.view.tabs.map(tab => tab._id).join('-');
    },
    draggableOptions() {
      return {
        animation: VUETIFY_ANIMATION_DELAY,
        disabled: !this.isEditingMode,
      };
    },
  },
  watch: {
    isEditingMode(value) {
      this.$nextTick(this.callTabsOnResizeMethod);

      if (!value) {
        this.updateViewMethod({
          ...this.view,

          tabs: this.tabs,
        });
      }
    },
    tabs: {
      immediate: true,
      handler(value) {
        console.log(value);
        this.onUpdateTabs();
      },
    },
  },
  methods: {
    showUpdateTabModal(tab) {
      this.showModal({
        name: MODALS.textFieldEditor,
        config: {
          title: this.$t('modals.viewTab.edit.title'),
          field: {
            name: 'text',
            label: this.$t('modals.viewTab.fields.title'),
            value: tab.title,
            validationRules: 'required',
          },
          action: (title) => {
            const newTab = { ...tab, title };

            return this.updateTab(newTab);
          },
        },
      });
    },

    showDeleteTabModal(tab) {
      this.showModal({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            const view = {
              ...this.view,
              tabs: this.view.tabs.filter(viewTab => viewTab._id !== tab._id),
            };

            await this.updateViewMethod(view);
          },
        },
      });
    },

    updateTab(newTab) {
      const view = {
        ...this.view,
        tabs: this.view.tabs.map((viewTab) => {
          if (viewTab._id === newTab._id) {
            return newTab;
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
    & /deep/ .v-tabs__bar {
      display: none;
    }
  }

  .draggable-item {
    position: relative;
    transform: translateZ(0);
  }
</style>

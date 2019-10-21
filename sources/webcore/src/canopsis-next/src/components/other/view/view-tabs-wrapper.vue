<template lang="pug">
  div
    view-tabs.tabs-absolute(
      v-if="view && isTabsChanged",
      :view="view",
      :tabs.sync="tabs",
      :isTabsChanged="isTabsChanged",
      :isEditingMode="isEditingMode",
      :hasUpdateAccess="hasUpdateAccess"
    )
    v-fade-transition
      div
        .v-overlay.v-overlay--active(v-show="view && isTabsChanged")
          v-btn(
            data-test="submitMoveTab",
            color="primary",
            @click="submit"
          ) {{ $t('common.submit') }}
          v-btn(@click="cancel") {{ $t('common.cancel') }}
    view-tabs(
      :view="view",
      :tabs.sync="tabs",
      :isTabsChanged="isTabsChanged",
      :isEditingMode="isEditingMode",
      :hasUpdateAccess="hasUpdateAccess",
      :updateViewMethod="data => updateViewMethod(data)"
    )
      template(slot-scope="props")
        grid-layout(:layout.sync="props.tab.layout")
          grid-item(v-for="item in props.tab.layout", :x="item.x", :y="item.y", :w="item.w", :h="item.h", :i="item.i")
            widget-wrapper(
              :widget="findWidgetInTabById(props.tab._id, item.i)",
              :tab="props.tab",
              :isEditingMode="isEditingMode"
            )
</template>

<script>
import { isEqual } from 'lodash';
import { GridLayout, GridItem } from 'vue-grid-layout';

import WidgetWrapper from '@/components/widgets/widget-wrapper.vue';

import ViewTabs from './view-tabs.vue';

export default {
  components: {
    ViewTabs,
    GridLayout,
    GridItem,
    WidgetWrapper,
  },
  props: {
    view: {
      type: Object,
      required: true,
    },
    hasUpdateAccess: {
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
  },
  data() {
    return {
      tabs: [...this.view.tabs],
    };
  },
  computed: {
    isTabsChanged() {
      if (this.view.tabs.length === this.tabs.length) {
        return this.view.tabs.some((tab, index) => this.tabs[index] && tab._id !== this.tabs[index]._id);
      }

      return true;
    },
  },
  watch: {
    'view.tabs': {
      handler(tabs, prevTabs) {
        if (!isEqual(tabs, prevTabs)) {
          this.tabs = [...tabs];
        }
      },
    },
  },
  methods: {
    cancel() {
      this.tabs = [...this.view.tabs];
    },

    async submit() {
      this.updateViewMethod({
        ...this.view,

        tabs: this.tabs,
      });
    },

    findWidgetInTabById(tabId, widgetId) {
      const tabIndex = this.view.tabs.findIndex(tab => tab._id === tabId);

      return this.view.tabs[tabIndex].widgets.find(widget => widget._id === widgetId);
    },
  },
};
</script>

<style lang="scss" scoped>
  .tabs-absolute {
    position: absolute;
    z-index: 6;
    width: 100%;
  }
</style>

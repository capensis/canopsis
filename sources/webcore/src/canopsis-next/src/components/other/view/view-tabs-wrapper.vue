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
      :updateViewMethod="updateViewMethod"
    )
      view-tab-widgets(
        slot-scope="props",
        v-bind="props",
        @update:widgetsFields="$emit('update:widgetsFields', $event)"
      )
</template>

<script>
import { isEqual } from 'lodash';

import ViewTabs from './view-tabs.vue';
import ViewTabWidgets from './view-tab-widgets.vue';

export default {
  components: {
    ViewTabs,
    ViewTabWidgets,
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
  },
};
</script>

<style lang="scss" scoped>
  .tabs-absolute {
    position: absolute;
    z-index: 9;
    width: 100%;
  }
</style>

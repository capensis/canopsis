<template lang="pug">
  div
    view-tabs.absolute(
    v-if="isTabsChanged",
    :value="value",
    :view="view",
    :tabs.sync="tabs",
    :isTabsChanged="isTabsChanged",
    :isEditingMode="isEditingMode",
    :hasUpdateAccess="hasUpdateAccess",
    :updateViewMethod="data => updateViewMethod(data)"
    hide-slider,
    )
    v-fade-transition
      div
        .v-overlay.v-overlay--active(v-show="view && isTabsChanged")
          v-btn(@click="submit") Submit
          v-btn(@click="cancel") Cancel
    view-tabs(
    :value="value",
    :view="view",
    :tabs.sync="tabs",
    :isTabsChanged="isTabsChanged",
    :isEditingMode="isEditingMode",
    :hasUpdateAccess="hasUpdateAccess",
    :updateViewMethod="data => updateViewMethod(data)",
    @input="$emit('input', $event)"
    )
      view-tab-rows(
      slot-scope="props",
      v-bind="props",
      )
</template>

<script>
import isEqual from 'lodash/isEqual';

import ViewTabs from './view-tabs.vue';
import ViewTabRows from './view-tab-rows.vue';

export default {
  components: {
    ViewTabs,
    ViewTabRows,
  },
  props: {
    value: {
      type: Number,
      default: null,
    },
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
    submit() {
      this.updateViewMethod({
        ...this.view,

        tabs: this.tabs,
      });
    },
  },
};
</script>

<style lang="scss" scoped>
  .absolute {
    position: absolute;
    z-index: 6;
    width: 100%;
  }
</style>

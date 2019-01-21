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
          v-btn(color="primary", @click="submit") {{ $t('common.submit') }}
          v-btn(@click="cancel") {{ $t('common.cancel') }}
    view-tabs(
    :view="view",
    :value="value",
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
    async submit() {
      const activeTab = this.view.tabs[this.value];
      const activeTabIndex = this.tabs.findIndex(tab => activeTab._id === tab._id);

      await this.updateViewMethod({
        ...this.view,

        tabs: this.tabs,
      });

      this.$emit('input', activeTabIndex);
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

<template lang="pug">
  div
    view-tabs.tabs-absolute(
      v-if="view && isTabsChanged",
      :tabs.sync="tabs",
      :changed="isTabsChanged",
      :updatable="updatable"
    )
    v-fade-transition
      div.v-overlay.v-overlay--active(v-show="view && isTabsChanged")
        v-btn(color="primary", @click="submit") {{ $t('common.submit') }}
        v-btn(@click="cancel") {{ $t('common.cancel') }}
    view-tabs(
      :view="view",
      :tabs.sync="tabs",
      :changed="isTabsChanged",
      :editing="editing",
      :updatable="updatable"
    )
      view-tab-widgets(
        slot-scope="props",
        v-bind="props"
      )
</template>

<script>
import { isEqual } from 'lodash';

import { mapIds } from '@/helpers/entities';

import { activeViewMixin } from '@/mixins/active-view';
import { entitiesViewTabMixin } from '@/mixins/entities/view/tab';

import ViewTabs from './view-tabs.vue';
import ViewTabWidgets from './view-tab-widgets.vue';

export default {
  components: {
    ViewTabs,
    ViewTabWidgets,
  },
  mixins: [
    activeViewMixin,
    entitiesViewTabMixin,
  ],
  props: {
    updatable: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      tabs: [],
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
      immediate: true,
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
      await this.updateViewTabPositions({
        data: mapIds(this.tabs),
      });

      return this.fetchActiveView();
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

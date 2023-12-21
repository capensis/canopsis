<template>
  <div>
    <view-tabs
      class="view-tabs--absolute"
      v-if="view && isTabsChanged"
      :tabs.sync="tabs"
      :changed="isTabsChanged"
      :updatable="updatable"
    />
    <v-fade-transition>
      <v-overlay
        class="view-tabs__overlay"
        :value="view && isTabsChanged"
        z-index="10"
      >
        <v-btn
          class="ma-2"
          color="primary"
          @click="submit"
        >
          {{ $t('common.submit') }}
        </v-btn>
        <v-btn
          class="ma-2"
          light
          @click="cancel"
        >
          {{ $t('common.cancel') }}
        </v-btn>
      </v-overlay>
    </v-fade-transition>
    <view-tabs
      :tabs.sync="tabs"
      :changed="isTabsChanged"
      :editing="editing"
      :updatable="updatable"
    >
      <template #default="props">
        <view-tab-widgets v-bind="props" />
      </template>
    </view-tabs>
  </div>
</template>

<script>
import { isEqual } from 'lodash';

import { mapIds } from '@/helpers/array';

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
.view-tabs {
  &--absolute {
    position: absolute;
    z-index: 11;
    width: 100%;
  }

  &__overlay {
    align-items: flex-start;
    justify-content: flex-start;
  }
}
</style>

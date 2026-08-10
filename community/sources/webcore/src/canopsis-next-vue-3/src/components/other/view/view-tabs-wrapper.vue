<template>
  <div class="view-tabs-stack">
    <view-tabs
      v-if="view && isTabsChanged"
      :tabs="tabs"
      :changed="isTabsChanged"
      :updatable="updatable"
      class="view-tabs--absolute"
      @update:tabs="tabs = $event"
    />

    <div v-if="view && isTabsChanged" class="view-tabs__scrim" />

    <div v-if="view && isTabsChanged" class="view-tabs__actions">
      <v-btn class="ma-2" color="primary" @click="submit">
        {{ $t('common.submit') }}
      </v-btn>

      <v-btn class="ma-2" variant="elevated" @click="cancel">
        {{ $t('common.cancel') }}
      </v-btn>
    </div>

    <view-tabs
      :tabs="tabs"
      :changed="isTabsChanged"
      :updatable="updatable"
      @update:tabs="tabs = $event"
    >
      <template #default="props">
        <slot v-bind="props">
          <view-tab-widgets v-bind="props" />
        </slot>
      </template>
    </view-tabs>
  </div>
</template>

<script>
import { isEqual } from 'lodash';
import { mapActions } from 'pinia';

import { mapIds } from '@/helpers/array';
import { activeViewMixin } from '@/mixins/active-view';
import { useViewStore } from '@/stores/modules/view';

import ViewTabWidgets from './view-tab-widgets.vue';
import ViewTabs from './view-tabs.vue';

export default {
  components: {
    ViewTabs,
    ViewTabWidgets,
  },
  mixins: [activeViewMixin],
  props: {
    updatable: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      tabs: [],
      initialTabsIds: [],
    };
  },
  computed: {
    isTabsChanged() {
      return !isEqual(mapIds(this.tabs), this.initialTabsIds);
    },
  },
  watch: {
    'view.tabs': {
      immediate: true,
      handler(tabs, prevTabs) {
        if (!isEqual(tabs, prevTabs)) {
          this.tabs = [...tabs];
          this.initialTabsIds = mapIds(tabs);
        }
      },
    },
  },
  methods: {
    ...mapActions(useViewStore, {
      updateViewTabPositions: 'updateViewTabPositions',
    }),

    cancel() {
      this.tabs = [...this.view.tabs];
      this.initialTabsIds = mapIds(this.view.tabs);
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
    z-index: 1006;
    width: 100%;
  }

  &__actions {
    position: absolute;
    z-index: 1007;
    top: -50px;
    left: 0;
  }

  &__scrim {
    position: fixed;
    z-index: 1005;
    inset: 0;
    background: rgba(33, 33, 33, 0.46);
  }
}

.view-tabs-stack {
  position: relative;
}
</style>

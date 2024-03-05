<template>
  <v-layout column>
    <v-layout>
      <v-flex
        class="text-center font-weight-bold"
        xs4
      >
        {{ $tc('common.group') }}
      </v-flex>
      <v-flex
        class="text-center font-weight-bold"
        xs4
      >
        {{ $tc('common.view') }}
      </v-flex>
      <v-flex
        class="text-center font-weight-bold"
        xs4
      >
        {{ $tc('common.tab') }}
      </v-flex>
    </v-layout>
    <v-layout column>
      <c-draggable-list-field
        v-field="tabs"
        :class="{ 'tabs-draggable-panel--empty': isTabsEmpty, 'tabs-draggable-panel--disabled': disabled }"
        :disabled="disabled"
        class="tabs-draggable-panel secondary lighten-1"
      >
        <tab-panel-content
          v-for="{ tab, view, group } in tabsWithDetails"
          :key="tab._id"
          :tab="tab"
          hide-actions="hideActions"
        >
          <template #title="">
            <playlist-tab-item
              :tab="tab"
              :view="view"
              :group="group"
            />
          </template>
        </tab-panel-content>
      </c-draggable-list-field>
    </v-layout>
  </v-layout>
</template>

<script>
import { entitiesViewGroupMixin } from '@/mixins/entities/view/group';

import TabPanelContent from '@/components/other/playlists/partials/tab-panel-content.vue';
import PlaylistTabItem from '@/components/other/playlists/partials/playlist-tab-item.vue';

export default {
  components: { TabPanelContent, PlaylistTabItem },
  mixins: [entitiesViewGroupMixin],
  model: {
    prop: 'tabs',
    event: 'change',
  },
  props: {
    tabs: {
      type: Array,
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    tabsDetailsById() {
      return this.groups.reduce((acc, group) => {
        group.views.forEach((view) => {
          view.tabs.forEach((tab) => {
            acc[tab._id] = {
              tab,
              group,
              view,
            };
          });
        });

        return acc;
      }, {});
    },

    tabsWithDetails() {
      return this.tabs
        .map(tab => this.tabsDetailsById[tab._id])
        .filter(Boolean);
    },

    isTabsEmpty() {
      return this.tabs.length === 0;
    },
  },
};
</script>

<style lang="scss" scoped>
  .tabs-draggable-panel {
    &:not(&--disabled) ::v-deep .tab-panel-item {
      cursor: move;
    }

    &--empty {
      &:after {
        content: '';
        display: block;
        height: 48px;
        border: 4px dashed #4f6479;
        border-radius: 5px;
        position: relative;
      }
    }
  }
</style>

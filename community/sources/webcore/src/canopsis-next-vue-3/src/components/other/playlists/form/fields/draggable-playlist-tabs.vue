<template>
  <c-layout class="draggable-playlist-tabs" column>
    <c-layout class="draggable-playlist-tabs__header">
      <div class="draggable-playlist-tabs__header-cell">
        {{ t('common.group', 1) }}
      </div>
      <div class="draggable-playlist-tabs__header-cell">
        {{ t('common.view', 1) }}
      </div>
      <div class="draggable-playlist-tabs__header-cell">
        {{ t('common.tab', 1) }}
      </div>
    </c-layout>

    <c-draggable-list-field
      :model-value="tabs"
      :class="{
        'tabs-draggable-panel--empty': isTabsEmpty,
        'tabs-draggable-panel--disabled': disabled,
      }"
      :disabled="disabled"
      class="tabs-draggable-panel"
      item-key="_id"
      @update:model-value="$emit('update:tabs', $event)"
    >
      <template #item="{ element }">
        <div class="tabs-draggable-panel__item">
          <playlist-tab-item
            v-if="tabsDetailsById[element._id]"
            :tab="tabsDetailsById[element._id].tab"
            :view="tabsDetailsById[element._id].view"
            :group="tabsDetailsById[element._id].group"
          />
        </div>
      </template>
    </c-draggable-list-field>
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import PlaylistTabItem from '@/components/other/playlists/partials/playlist-tab-item.vue';
import { useI18n } from '@/hooks/i18n';

export default {
  components: { PlaylistTabItem },

  props: {
    tabs: {
      type: Array,
      required: true,
    },
    groups: {
      type: Array,
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:tabs'],

  setup(props) {
    const { t } = useI18n();

    const tabsDetailsById = computed(() =>
      props.groups.reduce((acc, group) => {
        group.views?.forEach((view) => {
          view.tabs?.forEach((tab) => {
            acc[tab._id] = { tab, group, view };
          });
        });

        return acc;
      }, {}),
    );

    const isTabsEmpty = computed(() => props.tabs.length === 0);

    return {
      t,
      tabsDetailsById,
      isTabsEmpty,
    };
  },
};
</script>

<style lang="scss" scoped>
.draggable-playlist-tabs__header {
  min-height: 24px;
}

.draggable-playlist-tabs__header-cell {
  width: 33.333%;
  color: #000;
  font-size: 14px;
  font-weight: 500;
  line-height: 16px;
  text-align: center;
}

.tabs-draggable-panel {
  background: transparent !important;

  &:not(.tabs-draggable-panel--disabled) {
    :deep(.tabs-draggable-panel__item) {
      cursor: move;
    }
  }

  &--empty {
    &:after {
      content: '';
      display: block;
      height: 48px;
      border: 1px dashed #4f6479;
      border-radius: 0;
      position: relative;
    }
  }
}
</style>

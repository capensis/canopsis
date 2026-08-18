<template>
  <div class="manage-playlist-tabs">
    <div class="text-center mb-2">
      {{ t('modals.createPlaylist.groups') }}
    </div>

    <v-expansion-panels
      :model-value="openedPanels"
      variant="accordion"
      multiple
      readonly
    >
      <group-panel
        v-for="group in groups"
        :key="group._id"
        :group="group"
        hide-actions
      >
        <template #title>
          <v-checkbox
            :model-value="isGroupSelected(group)"
            :disabled="isDisabledGroup(group)"
            :label="group.title"
            class="manage-playlist-tabs__checkbox"
            color="primary"
            density="compact"
            hide-details
            @click.stop
            @update:model-value="selectGroupHandler(group, $event)"
          />
        </template>

        <v-expansion-panels
          :model-value="getPanelValueFromArray(group.views)"
          variant="accordion"
          multiple
          readonly
        >
          <v-expansion-panel v-for="view in group.views" :key="view._id" class="tabs-panel">
            <v-expansion-panel-title class="pa-0" hide-actions>
              <group-view-panel :view="view">
                <template #title>
                  <v-checkbox
                    :model-value="isViewSelected(view)"
                    :disabled="isDisabledView(view)"
                    class="manage-playlist-tabs__checkbox"
                    color="primary"
                    density="compact"
                    hide-details
                    @click.stop
                    @update:model-value="selectViewHandler(view, $event)"
                  >
                    <template #label>
                      <span class="text-truncate w-100">
                        {{ view.title }}
                        <span v-show="view.description" class="ml-1">
                          ({{ view.description }})
                        </span>
                      </span>
                    </template>
                  </v-checkbox>
                </template>
              </group-view-panel>
            </v-expansion-panel-title>

            <v-expansion-panel-text>
              <v-expansion-panels
                :model-value="getPanelValueFromArray(view.tabs)"
                variant="accordion"
                multiple
                readonly
              >
                <v-expansion-panel v-for="tab in view.tabs" :key="tab._id">
                  <v-expansion-panel-title class="pa-0" hide-actions>
                    <div class="tab-panel-item secondary lighten-2">
                      <v-checkbox
                        :model-value="selectedTabsIds.includes(tab._id)"
                        :label="tab.title"
                        class="manage-playlist-tabs__checkbox ml-10"
                        color="primary"
                        density="compact"
                        hide-details
                        @click.stop
                        @update:model-value="selectTabHandler(tab, $event)"
                      />
                    </div>
                  </v-expansion-panel-title>
                </v-expansion-panel>
              </v-expansion-panels>
            </v-expansion-panel-text>
          </v-expansion-panel>
        </v-expansion-panels>
      </group-panel>
    </v-expansion-panels>
  </div>
</template>

<script>
import { computed } from 'vue';

import { createRangeArray } from '@/helpers/array';

import { useI18n } from '@/hooks/i18n';

import GroupPanel from '@/components/layout/navigation/partials/groups-side-bar/group-panel.vue';
import GroupViewPanel from '@/components/layout/navigation/partials/groups-side-bar/group-view-panel.vue';

export default {
  components: {
    GroupPanel,
    GroupViewPanel,
  },

  props: {
    groups: {
      type: Array,
      required: true,
    },
    selectedTabs: {
      type: Array,
      default: () => [],
    },
  },

  emits: ['update:selectedTabs'],

  setup(props, { emit }) {
    const { t } = useI18n();

    const openedPanels = computed(() => createRangeArray(props.groups.length));
    const selectedTabsIds = computed(() => props.selectedTabs.map(({ _id }) => _id));

    const getPanelValueFromArray = values => createRangeArray(values?.length ?? 0);
    const isDisabledView = view => !view.tabs?.length;
    const isDisabledGroup = group => !group.views?.length || group.views.every(isDisabledView);
    const isViewSelected = view =>
      !!view.tabs?.length && view.tabs.every(({ _id }) => selectedTabsIds.value.includes(_id));
    const isGroupSelected = group => !!group.views?.length && group.views.every(isViewSelected);

    const updateSelectedTabs = (tabs, checked) => {
      const tabIds = tabs.map(({ _id }) => _id);
      const tabsWithoutSelected = props.selectedTabs.filter(({ _id }) => !tabIds.includes(_id));

      emit(
        'update:selectedTabs',
        checked ? [...tabsWithoutSelected, ...tabs] : tabsWithoutSelected,
      );
    };

    const selectTabHandler = (tab, checked) => updateSelectedTabs([tab], checked);
    const selectViewHandler = (view, checked) => updateSelectedTabs(view.tabs ?? [], checked);
    const selectGroupHandler = (group, checked) => {
      const groupTabs =
        group.views?.reduce((acc, view) => {
          acc.push(...(view.tabs ?? []));

          return acc;
        }, []) ?? [];

      updateSelectedTabs(groupTabs, checked);
    };

    return {
      t,
      openedPanels,
      selectedTabsIds,
      getPanelValueFromArray,
      isDisabledView,
      isDisabledGroup,
      isViewSelected,
      isGroupSelected,
      selectTabHandler,
      selectViewHandler,
      selectGroupHandler,
    };
  },
};
</script>

<style lang="scss" scoped>
.manage-playlist-tabs {
  :deep(.v-expansion-panel-text__wrapper) {
    padding: 0;
  }

  :deep(.group-panel__title) {
    display: flex;
    align-items: center;
    height: 48px;
    padding-top: 0 !important;
    padding-bottom: 0 !important;
  }

  :deep(.group-view-panel),
  :deep(.panel-item-content),
  :deep(.panel-item-content__title) {
    min-height: 48px;
    height: 48px;
    overflow: hidden;
  }

  :deep(.group-view-panel),
  :deep(.panel-item-content),
  :deep(.panel-item-content__title) {
    width: 100%;
    max-width: 100%;
  }

  :deep(.group-view-panel) {
    flex: 1 1 auto;
  }

  :deep(.panel-item-content) {
    padding-top: 0;
    padding-bottom: 0;
  }

  &__checkbox {
    width: 100%;
    min-height: 48px;
    height: 48px;
    margin: 0;
    padding: 0;
    overflow: hidden;

    :deep(.v-input__control),
    :deep(.v-selection-control) {
      min-height: 48px;
      height: 48px;
    }

    :deep(.v-selection-control__wrapper) {
      flex: 0 0 40px;
      height: 48px;
    }

    :deep(.v-label) {
      min-width: 0;
      overflow: hidden;
    }
  }
}

.tabs-panel {
  :deep(.v-expansion-panel-title) {
    padding: 0;
    margin: 0;
  }
}

.tab-panel-item {
  display: flex;
  align-items: center;
  width: 100%;
  height: 48px;
}
</style>

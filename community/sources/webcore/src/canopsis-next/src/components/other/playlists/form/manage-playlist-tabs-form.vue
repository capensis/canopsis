<template>
  <v-layout>
    <v-flex
      class="manage-playlist-tabs mr-2"
      xs12
    >
      <v-flex class="text-center mb-2">
        {{ $t('modals.createPlaylist.groups') }}
      </v-flex>
      <v-expansion-panels
        :value="openedPanels"
        accordion
        readonly
        multiple
        dark
      >
        <group-panel
          v-for="group in groups"
          :key="group._id"
          :group="group"
          hide-actions
        >
          <template #title="">
            <v-layout align-center>
              <v-checkbox
                :input-value="selectedGroupsIds"
                :value="group._id"
                :disabled="isDisabledGroup(group)"
                class="group-checkbox mt-0 pt-0"
                color="primary"
                @change="selectGroupHandler(group, $event)"
              >
                <template #label>
                  {{ group.title }}
                </template>
              </v-checkbox>
            </v-layout>
          </template>
          <v-expansion-panels
            :value="getPanelValueFromArray(group.views)"
            accordion
            readonly
            hide-actions
            multiple
            dark
          >
            <v-expansion-panel
              v-for="view in group.views"
              :key="view._id"
              class="tabs-panel"
            >
              <v-expansion-panel-header
                class="pa-0"
                hide-actions
              >
                <group-view-panel :view="view">
                  <template #title="">
                    <v-checkbox
                      :input-value="selectedViewsIds"
                      :value="view._id"
                      :disabled="isDisabledView(view)"
                      class="group-checkbox mt-0 pt-0"
                      color="primary"
                      @change="selectViewHandler(view, $event)"
                    >
                      <template #label>
                        <span class="text-truncate fill-width">
                          {{ view.title }}
                          <span
                            v-show="view.description"
                            class="ml-1"
                          >
                            ({{ view.description }})
                          </span>
                        </span>
                      </template>
                    </v-checkbox>
                  </template>
                </group-view-panel>
              </v-expansion-panel-header>
              <v-expansion-panel-content hide-actions>
                <v-expansion-panels
                  :value="getPanelValueFromArray(view.tabs)"
                  accordion
                  readonly
                  hide-actions
                  multiple
                  dark
                >
                  <v-expansion-panel
                    v-for="tab in view.tabs"
                    :key="tab._id"
                  >
                    <v-expansion-panel-header
                      class="pa-0"
                      hide-actions
                    >
                      <tab-panel-content :tab="tab">
                        <template #title="">
                          <v-checkbox
                            :input-value="selectedTabsIds"
                            :value="tab._id"
                            :label="tab.title"
                            class="ml-10"
                            color="primary"
                            @change="selectTabHandler(tab, $event)"
                          >
                            <template #label>
                              {{ tab.title }}
                            </template>
                          </v-checkbox>
                        </template>
                      </tab-panel-content>
                    </v-expansion-panel-header>
                  </v-expansion-panel>
                </v-expansion-panels>
              </v-expansion-panel-content>
            </v-expansion-panel>
          </v-expansion-panels>
        </group-panel>
      </v-expansion-panels>
    </v-flex>
  </v-layout>
</template>

<script>
import { createRangeArray } from '@/helpers/array';

import GroupViewPanel from '@/components/layout/navigation/partials/groups-side-bar/group-view-panel.vue';
import GroupPanel from '@/components/layout/navigation/partials/groups-side-bar/group-panel.vue';

import TabPanelContent from '../partials/tab-panel-content.vue';

export default {
  components: {
    TabPanelContent,
    GroupViewPanel,
    GroupPanel,
  },
  model: {
    prop: 'selectedTabs',
    event: 'input',
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
  computed: {
    openedPanels() {
      return this.getPanelValueFromArray(this.groups);
    },

    selectedTabsIds() {
      return this.selectedTabs.map(({ _id }) => _id);
    },

    selectedViewsIds() {
      return this.groups.reduce((acc, { views }) => {
        views.forEach(({ _id: viewId, tabs = [] }) => {
          if (tabs.length && tabs.every(({ _id: tabId }) => this.selectedTabsIds.includes(tabId))) {
            acc.push(viewId);
          }
        });

        return acc;
      }, []);
    },

    selectedGroupsIds() {
      return this.groups.reduce((acc, { _id: groupId, views }) => {
        if (views.length && views.every(({ _id: viewId }) => this.selectedViewsIds.includes(viewId))) {
          acc.push(groupId);
        }

        return acc;
      }, []);
    },
  },
  methods: {
    getPanelValueFromArray(values = []) {
      return createRangeArray(values.length);
    },

    selectTabHandler(tab, checkedTabs) {
      const checked = checkedTabs.includes(tab._id);

      this.updateSelectedTabs([tab], checked);
    },

    selectViewHandler(view, checkedViewsIds) {
      const checked = checkedViewsIds.includes(view._id);

      this.updateSelectedTabs(view.tabs, checked);
    },

    selectGroupHandler(group, selectedGroups) {
      const checked = selectedGroups.includes(group._id);
      const groupTabs = group.views.reduce((acc, view) => {
        acc.push(...view.tabs);
        return acc;
      }, []);

      this.updateSelectedTabs(groupTabs, checked);
    },

    updateSelectedTabs(tabs, checked) {
      const tabIds = tabs.map(({ _id }) => _id);
      const tabsWithoutSelected = this.selectedTabs.filter(({ _id: tabId }) => !tabIds.includes(tabId));

      this.$emit('input', !checked ? tabsWithoutSelected : [...tabsWithoutSelected, ...tabs]);
    },

    isDisabledGroup(group) {
      return !group.views.length || group.views.every(this.isDisabledView);
    },

    isDisabledView(view) {
      return !view.tabs?.length;
    },
  },
};
</script>

<style lang="scss" scoped>
  .manage-playlist-tabs {
    & ::v-deep .panel-header {
      display: flex;
      flex: inherit;
      align-items: center;
    }
    & ::v-deep .v-expansion-panel__body {
      transition: none !important;
    }
  }
  .tabs-panel {
    & ::v-deep .v-expansion-panel__header {
      padding: 0;
      margin: 0;
    }
  }
</style>

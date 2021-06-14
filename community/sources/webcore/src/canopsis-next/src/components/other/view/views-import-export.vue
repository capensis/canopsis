<template lang="pug">
  div
    v-layout
      v-flex.export-views-block.mr-0.ma-4(xs6)
        v-checkbox(
          v-model="isAllSelected",
          :label="$t('importExportViews.selectAll')",
          color="primary"
        )
        v-expansion-panel(readonly, hide-actions, expand, dark, focusable, :value="openedPanels")
          group-panel(
            v-for="group in groups",
            :group="group",
            :key="group._id",
            hideActions
          )
            template(slot="title")
              v-checkbox.group-checkbox.mt-0.pt-0(
                v-model="selectedGroupsIds",
                :value="group._id",
                color="primary",
                @change="changeGroupHandler(group._id, $event)"
              )
              span.group-title {{ group.title }}
            group-view-panel(
              v-for="view in group.views",
              :key="view._id",
              :view="view"
            )
              template(slot="title")
                v-layout(align-center, row, justify-space-between)
                  v-checkbox(
                    v-model="selectedViewIds",
                    :value="view._id",
                    color="primary"
                  )
                  span.ellipsis {{ view.title }}
                    span.ml-1(v-show="view.description") ({{ view.description }})
      v-flex.btn-group(xs2)
        v-layout(column)
          v-btn(:disabled="selectedDataIsEmpty", @click="exportViews") {{ $t('common.export') }}
          file-selector.ma-2.view-import-selector(
            ref="fileSelector",
            multiple,
            hide-details,
            @change="importViews"
          )
            v-btn.import-btn.ma-0(slot="activator") {{ $t('common.import') }}
</template>

<script>
import { EXPORT_VIEWS_AND_GROUPS_PREFIX } from '@/config';
import { MODALS } from '@/constants';

import { saveJsonFile } from '@/helpers/file/files';
import { getFileTextContent } from '@/helpers/file/file-select';
import { getExportedGroupsAndViews } from '@/helpers/forms/view';

import { entitiesViewGroupMixin } from '@/mixins/entities/view/group';
import entitiesViewMixin from '@/mixins/entities/view';

import GroupsSideBarGroup from '@/components/layout/navigation/partial/groups-side-bar/groups-side-bar-group.vue';
import FileSelector from '@/components/forms/fields/file-selector.vue';
import GroupPanel from '@/components/layout/navigation/partial/groups-side-bar/group-panel.vue';
import GroupViewPanel from '@/components/layout/navigation/partial/groups-side-bar/group-view-panel.vue';

export default {
  components: {
    GroupViewPanel,
    FileSelector,
    GroupPanel,
    GroupsSideBarGroup,
  },
  mixins: [
    entitiesViewMixin,
    entitiesViewGroupMixin,
  ],
  data() {
    return {
      selectedGroupsIds: [],
      selectedViewIds: [],
    };
  },
  computed: {
    openedPanels() {
      return new Array(this.groups.length).fill(true);
    },
    selectedDataIsEmpty() {
      return !this.selectedGroupsIds.length && !this.selectedViewIds.length;
    },
    viewIds() {
      return this.groups.reduce((acc, { views }) => {
        acc.push(...views.map(({ _id }) => _id));
        return acc;
      }, []);
    },
    groupIds() {
      return this.groups.map(({ _id }) => _id);
    },
    isAllSelected: {
      get() {
        return this.groupIds.every(id => this.selectedGroupsIds.includes(id))
          && this.viewIds.every(id => this.selectedViewIds.includes(id));
      },
      set(checked) {
        if (checked) {
          this.selectedGroupsIds = [...this.groupIds];
          this.selectedViewIds = [...this.viewIds];
        } else {
          this.resetSelected();
        }
      },
    },
  },
  methods: {
    async importViews([file]) {
      try {
        const content = await getFileTextContent(file);
        const { groups = [], views = [] } = JSON.parse(content);

        this.$modals.show({
          name: MODALS.importExportViews,
          config: {
            importedGroups: groups,
            importedViews: views,
          },
        });
      } catch (e) {
        this.$popups.error({ text: this.$t('errors.default') });
      }

      this.$refs.fileSelector.clear();
    },

    exportViews() {
      const exportData = getExportedGroupsAndViews({
        groups: this.selectedGroupsIds.map(this.getGroupById),
        views: this.selectedViewIds.map(this.getViewById),
      });

      saveJsonFile(exportData, `${EXPORT_VIEWS_AND_GROUPS_PREFIX}${new Date().toLocaleString()}`);

      this.resetSelected();
    },

    resetSelected() {
      this.selectedGroupsIds = [];
      this.selectedViewIds = [];
    },

    changeGroupHandler(groupId, selectedGroupsIds) {
      const checked = selectedGroupsIds.includes(groupId);
      const { views } = this.getGroupById(groupId);

      const selectedViewIds = views.map(({ _id }) => _id);
      const viewsIdsWithoutGroupViews = this.selectedViewIds.filter(viewId => !selectedViewIds.includes(viewId));

      if (checked) {
        this.selectedViewIds = [
          ...viewsIdsWithoutGroupViews,
          ...selectedViewIds,
        ];
      } else {
        this.selectedViewIds = viewsIdsWithoutGroupViews;
      }
    },
  },
};
</script>

<style lang="scss" scoped>
  .btn-group {
    display: flex;
    align-items: center;
  }

  .view-import-selector {
    display: inline-flex;

    & /deep/ .file-selector-button-wrapper {
      width: 100%;
    }

    .import-btn {
      cursor: pointer;
      width: 100%;
    }
  }

  .group-title {
    overflow: auto;
  }
  .export-views-block {
    & /deep/ .panel-header {
      display: flex;
      flex: inherit;
      align-items: center;
    }
  }
</style>

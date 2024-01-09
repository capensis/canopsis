<template>
  <div>
    <v-layout>
      <v-flex
        class="export-views-block ma-6"
        xs6
      >
        <v-checkbox
          v-model="isAllSelected"
          :label="$t('view.selectAll')"
          color="primary"
        />
        <views-export-expansion-panel
          v-model="selected"
          :groups="availableGroups"
        />
      </v-flex>
      <v-flex xs2>
        <v-layout
          column
          justify-center
          fill-height
        >
          <v-btn
            :disabled="selectedEmpty"
            color="primary"
            @click="exportViews"
          >
            <v-icon left>
              file_upload
            </v-icon>
            <span>{{ $t('common.export') }}</span>
          </v-btn>
          <file-selector
            class="my-2 view-import-selector"
            ref="fileSelector"
            multiple
            hide-details
            @change="importViews"
          >
            <template #activator="{ on, ...attrs }">
              <v-btn
                class="import-btn ma-0"
                v-bind="attrs"
                v-on="on"
                color="primary"
              >
                <v-icon left>
                  file_download
                </v-icon>
                <span>{{ $t('common.import') }}</span>
              </v-btn>
            </template>
          </file-selector>
        </v-layout>
      </v-flex>
    </v-layout>
  </div>
</template>

<script>
import { EXPORT_VIEWS_AND_GROUPS_FILENAME_PREFIX } from '@/config';
import { MODALS } from '@/constants';

import { saveJsonFile } from '@/helpers/file/files';
import { getFileTextContent } from '@/helpers/file/file-select';
import { getAllViewsFromGroups, exportedGroupsAndViewsToRequest } from '@/helpers/entities/view/form';

import { entitiesViewGroupMixin } from '@/mixins/entities/view/group';
import { entitiesViewMixin } from '@/mixins/entities/view';

import FileSelector from '@/components/forms/fields/file-selector.vue';
import ViewsExportExpansionPanel from '@/components/other/view/partials/views-export-expansion-panel.vue';

export default {
  components: {
    FileSelector,
    ViewsExportExpansionPanel,
  },
  mixins: [
    entitiesViewMixin,
    entitiesViewGroupMixin,
  ],
  data() {
    return {
      selected: {
        groups: [],
        views: [],
      },
    };
  },
  computed: {
    selectedEmpty() {
      return !this.selected.groups.length && !this.selected.views.length;
    },

    availableGroups() {
      return this.groups.filter(group => !group.is_private);
    },

    groupIds() {
      return this.availableGroups.map(({ _id }) => _id);
    },

    viewIds() {
      return getAllViewsFromGroups(this.availableGroups).map(({ _id }) => _id);
    },

    isAllSelected: {
      get() {
        return this.groupIds.every(id => this.selected.groups.includes(id))
          && this.viewIds.every(id => this.selected.views.includes(id));
      },
      set(checked) {
        if (checked) {
          this.selected = {
            groups: [...this.groupIds],
            views: [...this.viewIds],
          };
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
      } catch (err) {
        console.error(err);

        this.$popups.error({ text: this.$t('errors.default') });
      }

      this.$refs.fileSelector.clear();
    },

    async exportViews() {
      const data = exportedGroupsAndViewsToRequest({
        groups: this.selected.groups.map(this.getGroupById),
        views: this.selected.views.map(this.getViewById),
      });

      const result = await this.exportViewsWithoutStore({ data });

      saveJsonFile(result, `${EXPORT_VIEWS_AND_GROUPS_FILENAME_PREFIX}${new Date().toLocaleString()}`);

      this.resetSelected();
    },

    resetSelected() {
      this.selected = {
        groups: [],
        views: [],
      };
    },
  },
};
</script>

<style lang="scss" scoped>
  .view-import-selector {
    display: inline-flex;

    & ::v-deep .file-selector-button-wrapper {
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
    & ::v-deep .panel-header {
      display: flex;
      flex: inherit;
      align-items: center;
    }
  }
</style>

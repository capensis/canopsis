<template>
  <v-layout class="pa-4 gap-4" column>
    <v-layout>
      <v-checkbox
        :input-value="isAllSelected"
        :label="$t('view.selectAll')"
        color="primary"
        hide-details
        @change="allSelectedChange"
      />
      <v-layout class="gap-2" justify-end align-center>
        <v-btn
          :disabled="selectedEmpty"
          color="primary"
          depressed
          @click="exportViews"
        >
          <v-icon left>
            file_upload
          </v-icon>
          <span>{{ $t('common.export') }}</span>
        </v-btn>
        <file-selector
          ref="fileSelector"
          multiple
          hide-details
          @change="importViews"
        >
          <template #activator="{ on, ...attrs }">
            <v-btn
              v-bind="attrs"
              color="primary"
              outlined
              v-on="on"
            >
              <v-icon left>
                file_download
              </v-icon>
              <span>{{ $t('common.import') }}</span>
            </v-btn>
          </template>
        </file-selector>
      </v-layout>
    </v-layout>
    <views-export-expansion-panel
      v-model="selected"
      :groups="availableGroups"
    />
  </v-layout>
</template>

<script>
import { ref, computed } from 'vue';

import { EXPORT_VIEWS_AND_GROUPS_FILENAME_PREFIX } from '@/config';
import { MODALS } from '@/constants';

import { saveJsonFile } from '@/helpers/file/files';
import { getFileTextContent } from '@/helpers/file/file-select';
import { getAllViewsFromGroups, exportedGroupsAndViewsToRequest } from '@/helpers/entities/view/form';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';
import { useView } from '@/hooks/store/modules/view';

import FileSelector from '@/components/forms/fields/file-selector.vue';
import ViewsExportExpansionPanel from '@/components/other/view/partials/views-export-expansion-panel.vue';

export default {
  components: {
    FileSelector,
    ViewsExportExpansionPanel,
  },
  setup() {
    const { t } = useI18n();
    const modals = useModals();
    const popups = usePopups();

    const {
      groups,
      getGroupById,
      getViewById,
      exportViewsWithoutStore,
    } = useView();

    const selected = ref({
      groups: [],
      views: [],
    });

    const fileSelector = ref(null);

    const selectedEmpty = computed(() => !selected.value.groups.length && !selected.value.views.length);

    const availableGroups = computed(() => groups.value.filter(group => !group.is_private));
    const groupIds = computed(() => availableGroups.value.map(({ _id }) => _id));
    const viewIds = computed(() => getAllViewsFromGroups(availableGroups.value).map(({ _id }) => _id));
    const isAllSelected = computed(() => groupIds.value.every(id => selected.value.groups.includes(id))
      && viewIds.value.every(id => selected.value.views.includes(id)));

    /**
     * Resets the selected groups and views to empty arrays.
     */
    const resetSelected = () => {
      selected.value = {
        groups: [],
        views: [],
      };
    };

    /**
     * Handles the change event when the "select all" checkbox is toggled.
     * If checked, selects all available groups and views. Otherwise, resets the selection.
     *
     * @param {boolean} checked - Whether the checkbox is checked or not.
     */
    const allSelectedChange = (checked) => {
      if (checked) {
        selected.value = {
          groups: [...groupIds.value],
          views: [...viewIds.value],
        };

        return;
      }

      resetSelected();
    };

    /**
     * Imports views and groups from a JSON file.
     * Parses the file content and opens a modal to handle the import process.
     * Shows an error popup if the file parsing fails.
     *
     * @param {File[]} files - Array containing the file to import.
     */
    const importViews = async ([file]) => {
      try {
        const content = await getFileTextContent(file);
        const { groups: importedGroups = [], views: importedViews = [] } = JSON.parse(content);

        modals.show({
          name: MODALS.importExportViews,
          config: {
            importedGroups,
            importedViews,
          },
        });
      } catch (err) {
        console.error(err);

        popups.error({ text: t('errors.default') });
      }

      if (fileSelector.value) {
        fileSelector.value.clear();
      }
    };

    /**
     * Exports the selected views and groups to a JSON file.
     * Converts the selected items to the export format, fetches the data from the server,
     * and saves it as a JSON file with a timestamp in the filename.
     */
    const exportViews = async () => {
      const data = exportedGroupsAndViewsToRequest({
        groups: selected.value.groups.map(getGroupById.value),
        views: selected.value.views.map(getViewById.value),
      });

      const result = await exportViewsWithoutStore({ data });

      saveJsonFile(result, `${EXPORT_VIEWS_AND_GROUPS_FILENAME_PREFIX}${new Date().toLocaleString()}`);

      resetSelected();
    };

    return {
      selected,
      fileSelector,
      selectedEmpty,
      availableGroups,
      isAllSelected,
      allSelectedChange,
      importViews,
      exportViews,
    };
  },
};
</script>

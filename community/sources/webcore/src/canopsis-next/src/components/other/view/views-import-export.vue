<template>
  <c-layout class="pa-4" gap="4" column>
    <c-layout justify="space-between" align="center">
      <c-checkbox
        :model-value="isAllSelected"
        :label="t('view.selectAll')"
        color="primary"
        hide-details
        @update:model-value="allSelectedChange"
      />

      <c-layout gap="2" justify="end" align="center">
        <v-btn :disabled="selectedEmpty" color="primary" @click="exportViews">
          <v-icon start>
            file_upload
          </v-icon>
          <span>{{ t('common.export') }}</span>
        </v-btn>

        <file-selector
          ref="fileSelector"
          multiple
          hide-details
          @change="importViews"
        >
          <template #activator="{ on, attrs }">
            <v-btn
              v-bind="attrs"
              color="primary"
              variant="outlined"
              v-on="on"
            >
              <v-icon start>
                file_download
              </v-icon>
              <span>{{ t('common.import') }}</span>
            </v-btn>
          </template>
        </file-selector>
      </c-layout>
    </c-layout>

    <views-export-expansion-panel v-model:selected="selected" :groups="availableGroups" />
  </c-layout>
</template>

<script>
import { computed, onMounted, ref } from 'vue';

import { EXPORT_VIEWS_AND_GROUPS_FILENAME_PREFIX } from '@/config';
import { MAX_LIMIT, MODALS } from '@/constants';

import {
  exportedGroupsAndViewsToRequest,
  getAllViewsFromGroups,
} from '@/helpers/entities/view/form';
import { getFileTextContent } from '@/helpers/file/file-select';
import { saveJsonFile } from '@/helpers/file/files';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';
import { useView } from '@/hooks/store/modules/view';

import FileSelector from '@/components/forms/fields/file-selector.vue';

import ViewsExportExpansionPanel from './partials/views-export-expansion-panel.vue';

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
      viewGroups,
      getViewGroupById,
      getViewById,
      fetchViewGroupsList,
      exportViewWithoutStore,
    } = useView();

    const selected = ref({
      groups: [],
      views: [],
    });

    const fileSelector = ref(null);

    const selectedEmpty = computed(
      () => !selected.value.groups.length && !selected.value.views.length,
    );

    const availableGroups = computed(() => viewGroups.value.filter(group => !group.is_private));
    const groupIds = computed(() => availableGroups.value.map(({ _id }) => _id));
    const viewIds = computed(() => getAllViewsFromGroups(availableGroups.value).map(({ _id }) => _id),
    );
    const isAllSelected = computed(
      () => groupIds.value.every(id => selected.value.groups.includes(id)) &&
        viewIds.value.every(id => selected.value.views.includes(id)),
    );

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
      if (!file) {
        return;
      }

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

      fileSelector.value?.clear();
    };

    /**
     * Exports the selected views and groups to a JSON file.
     * Converts the selected items to the export format, fetches the data from the server,
     * and saves it as a JSON file with a timestamp in the filename.
     */
    const exportViews = async () => {
      const data = exportedGroupsAndViewsToRequest({
        groups: selected.value.groups.map(getViewGroupById),
        views: selected.value.views.map(getViewById),
      });

      const result = await exportViewWithoutStore({ data });

      saveJsonFile(
        result,
        `${EXPORT_VIEWS_AND_GROUPS_FILENAME_PREFIX}${new Date().toLocaleString()}`,
      );

      resetSelected();
    };

    onMounted(() => fetchViewGroupsList({
        params: {
          limit: MAX_LIMIT,
          page: 1,
          with_views: true,
        },
      }),
    );

    return {
      t,
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

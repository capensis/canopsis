<template>
  <div class="position-relative pa-4">
    <c-progress-overlay :pending="pending" />
    <external-data-table-records-list
      v-model:columns="columns"
      :records="records"
      :pending="recordsPending"
      :has-structure="!!externalDataTable.column_configs?.length"
      :options="options"
      :total-items="meta.total_count"
      :downloading="downloading"
      :columns-settings="columnsSettings"
      :draggable="draggableColumn"
      :resizable="resizableColumn"
      :cells-content-behavior="cellsContentBehavior"
      :dense="dense"
      with-toolbar
      with-actions
      densable
      selectable
      exportable
      expandable
      disabled-types
      @export="exportRecords"
      @import="showImportExternalDataTablesModal"
      @add="showCreateExternalDataTableRecordModal"
      @edit="showEditExternalDataTableRecordModal"
      @duplicate="showDuplicateExternalDataTableRecordModal"
      @remove="showRemoveExternalDataTableRecordModal"
      @remove-selected="showRemoveSelectedExternalDataTableRecordsModal"
      @update:dense="updateDense"
      @update:options="updateOptions"
      @update:columns-settings="updateColumnsSettings"
    />
  </div>
</template>

<script>
import { map } from 'lodash';
import { computed, onMounted, toRef } from 'vue';

import { getExternalDataTableRecordsFileUrl } from '@/helpers/entities/external-data-table/url';
import { prepareWidgetQuery, selectedToQuery, widgetToExportQueryColumns } from '@/helpers/entities/widget/query';
import { openUrlInNewTab } from '@/helpers/url';

import { useObserver } from '@/hooks/observer';
import { useFilePollingWithPending } from '@/hooks/polling';
import { useExternalDataTable } from '@/hooks/store/modules/external-data-table';
import { useWidgetPeriodicRefresh } from '@/hooks/widget/periodic-refresh';

import { useExternalDataTableRecordsList } from '@/components/other/external-data-table/hooks/external-data-table-records';

import ExternalDataTableRecordsList from '@/components/other/external-data-table/partials/external-data-table-records-list.vue';

import {
  useExternalDataTableWidgetColumnsSettings,
  useExternalDataTableWidgetDense,
  useExternalDataTableWidgetTable,
} from './hooks/external-data-table-widget';

export default {
  components: { ExternalDataTableRecordsList },
  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const {
      createExternalDataTableExport,
      fetchExternalDataTableExportStatus,
    } = useExternalDataTable();

    const {
      userPreference,
      updateContentInUserPreference,

      externalDataTable,

      columns,

      pending,
      fetchItem,
    } = useExternalDataTableWidgetTable({ widget: toRef(props, 'widget') });

    const initialQuery = computed(() => prepareWidgetQuery(props.widget, userPreference.value));
    const { observer } = useObserver({ key: '$refresh' });

    const {
      records,
      pending: recordsPending,
      meta,
      query,
      options,

      updateOptions: updateOptionsHook,
      resetQuery,
      fetchList: fetchRecordsList,

      showImportExternalDataTablesModal,
      showCreateExternalDataTableRecordModal,
      showEditExternalDataTableRecordModal,
      showDuplicateExternalDataTableRecordModal,
      showRemoveExternalDataTableRecordModal,
      showRemoveSelectedExternalDataTableRecordsModal,
    } = useExternalDataTableRecordsList({
      externalDataTable,
      initialQuery,
      refreshObserver: observer,
    });

    const updateOptions = (newOptions) => {
      if (newOptions.itemsPerPage !== options.value.itemsPerPage) {
        updateContentInUserPreference({ itemsPerPage: newOptions.itemsPerPage });
      }

      return updateOptionsHook(newOptions);
    };

    const fetchItemWithRecords = async () => {
      await fetchItem();

      resetQuery();

      return fetchRecordsList();
    };

    const {
      draggableColumn,
      resizableColumn,
      cellsContentBehavior,
      columnsSettings,
      updateColumnsSettings,
    } = useExternalDataTableWidgetColumnsSettings({
      widget: toRef(props, 'widget'),
      userPreference,
      updateContentInUserPreference,
    });

    const {
      pending: downloading,
      poll: exportRecords,
    } = useFilePollingWithPending({
      createHandler: selected =>
        createExternalDataTableExport({
          id: externalDataTable.value._id,
          data: {
            search: selectedToQuery(selected, query.value.search),
            separator: props.widget.parameters.exportCsvSeparator,
            fields: widgetToExportQueryColumns(props.widget),
            search_by: map(props.widget.parameters.widgetColumns, 'value'),
          },
        }),
      fetchHandler: fetchExternalDataTableExportStatus,
      endHandler: fileData => openUrlInNewTab(getExternalDataTableRecordsFileUrl(fileData?._id)),
    });

    const { dense, updateDense } = useExternalDataTableWidgetDense({
      widget: toRef(props, 'widget'),
      userPreference,
      updateContentInUserPreference,
    });

    useWidgetPeriodicRefresh({
      options: computed(() => props.widget.parameters.periodic_refresh),
      handler: fetchItemWithRecords,
    });

    observer.register(fetchItemWithRecords);
    onMounted(fetchItemWithRecords);

    return {
      externalDataTable,
      pending,

      columns,

      records,
      recordsPending,
      meta,
      options,
      updateOptions,

      showImportExternalDataTablesModal,
      showCreateExternalDataTableRecordModal,
      showEditExternalDataTableRecordModal,
      showDuplicateExternalDataTableRecordModal,
      showRemoveExternalDataTableRecordModal,
      showRemoveSelectedExternalDataTableRecordsModal,

      columnsSettings,
      draggableColumn,
      resizableColumn,
      cellsContentBehavior,
      updateColumnsSettings,

      downloading,
      exportRecords,

      dense,
      updateDense,
    };
  },
};
</script>

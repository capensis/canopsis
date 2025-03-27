<template>
  <div class="secondary pa-3">
    <v-card>
      <v-card-text>
        <external-data-table-records-list
          v-model="columns"
          :records="records"
          :pending="pending"
          :options="options"
          :total-items="meta.total_count"
          :has-structure="!!externalDataTable.columns?.length"
          with-toolbar
          with-actions
          selectable
          expandable
          @input="updateColumns"
          @import="showImportExternalDataTablesModal"
          @add="showCreateExternalDataTableRecordModal"
          @edit="showEditExternalDataTableRecordModal"
          @duplicate="showDuplicateExternalDataTableRecordModal"
          @remove="showRemoveExternalDataTableRecordModal"
          @remove-selected="showRemoveSelectedExternalDataTableRecordsModal"
          @update:options="updateOptions"
        />
      </v-card-text>
    </v-card>
  </div>
</template>
<script>
import { pick } from 'lodash';
import { ref, toRef } from 'vue';

import { externalDataTableColumnsToForm } from '@/helpers/entities/external-data-table/form';

import { useExternalDataTable } from '@/hooks/store/modules/external-data-table';

import ExternalDataTableRecordsList
  from '@/components/other/external-data-table/partials/external-data-table-records-list.vue';

import { useExternalDataTableRecordsList } from '../hooks/external-data-table-records';

export default {
  components: { ExternalDataTableRecordsList },
  props: {
    externalDataTable: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { updateExternalDataTable } = useExternalDataTable();

    const {
      records,
      pending,
      meta,
      options,
      updateOptions,

      showImportExternalDataTablesModal,
      showCreateExternalDataTableRecordModal,
      showEditExternalDataTableRecordModal,
      showDuplicateExternalDataTableRecordModal,
      showRemoveExternalDataTableRecordModal,
      showRemoveSelectedExternalDataTableRecordsModal,
    } = useExternalDataTableRecordsList({ externalDataTable: toRef(props, 'externalDataTable') });

    const columns = ref(externalDataTableColumnsToForm(
      props.externalDataTable.columns,
      props.externalDataTable.column_types,
    ));

    const updateColumns = newColumns => updateExternalDataTable({
      id: props.externalDataTable._id,
      data: {
        ...pick(props.externalDataTable, ['type', 'name', 'description']),

        column_types: Object.values(newColumns),
      },
    });

    return {
      records,
      meta,
      pending,
      options,
      updateOptions,

      columns,
      updateColumns,

      showImportExternalDataTablesModal,
      showCreateExternalDataTableRecordModal,
      showEditExternalDataTableRecordModal,
      showDuplicateExternalDataTableRecordModal,
      showRemoveExternalDataTableRecordModal,
      showRemoveSelectedExternalDataTableRecordsModal,
    };
  },
};
</script>

<template>
  <v-layout column>
    <field-title v-field="form.title" />
    <field-periodic-refresh v-field="form.parameters" />
    <widget-settings-group :title="$t('externalData.title')">
      <div class="pt-4">
        <external-data-table-table-field
          :value="form.parameters.table"
          required
          @input="updateTable"
        />
      </div>
    </widget-settings-group>
    <widget-settings-group
      :title="$t('settings.advancedSettings')"
      :disabled="!form.parameters.table"
    >
      <c-progress-overlay :pending="pending" />
      <field-columns-without-template
        v-field="form.parameters.widgetColumns"
        :label="$t('settings.columnNames')"
        :items="preparedColumns"
        without-custom-label
        required
      />
      <field-resize-column-behavior v-field="form.parameters.columns" />
      <field-default-sort-column
        v-field="form.parameters.sort"
        :columns="form.parameters.widgetColumns"
        :columns-label="$t('settings.columnName')"
        item-text="column"
        item-value="column"
      />
      <field-default-elements-per-page v-field="form.parameters.itemsPerPage" />
      <field-density v-field="form.parameters.dense" />
      <export-csv-form
        v-field="form.parameters"
        :items="preparedColumns"
        without-template
        without-infos-attributes
        without-custom-label
      />
    </widget-settings-group>
  </v-layout>
</template>

<script>
import { computed, onMounted } from 'vue';

import { setSeveralFields } from '@/helpers/immutable';
import { widgetColumnsToForm } from '@/helpers/entities/widget/column/form';
import { addPriorityColumnToColumnsArray } from '@/helpers/entities/external-data-table/form';

import { useModelField } from '@/hooks/form/model-field';

import { useExternalDataTableColumns } from '@/components/sidebars/external-data-table/form/hooks/external-data-table';

import FieldTitle from '@/components/sidebars/form/fields/title.vue';
import FieldPeriodicRefresh from '@/components/sidebars/form/fields/periodic-refresh.vue';
import FieldColumnsWithoutTemplate from '@/components/sidebars/form/fields/columns-without-template.vue';
import WidgetSettingsGroup from '@/components/sidebars/partials/widget-settings-group.vue';
import ExportCsvForm from '@/components/sidebars/form/export-csv.vue';
import FieldDefaultSortColumn from '@/components/sidebars/form/fields/default-sort-column.vue';
import FieldDensity from '@/components/sidebars/alarm/form/fields/density.vue';
import FieldResizeColumnBehavior from '@/components/sidebars/alarm/form/fields/resize-column-behavior.vue';
import FieldDefaultElementsPerPage from '@/components/sidebars/form/fields/default-elements-per-page.vue';
import CProgressOverlay from '@/components/common/overlay/c-progress-overlay.vue';
import ExternalDataTableTableField
  from '@/components/other/external-data-table/form/fields/external-data-table-table-field.vue';

export default {
  components: {
    ExternalDataTableTableField,
    CProgressOverlay,
    FieldDefaultElementsPerPage,
    FieldResizeColumnBehavior,
    FieldDensity,
    FieldDefaultSortColumn,
    ExportCsvForm,
    WidgetSettingsGroup,
    FieldColumnsWithoutTemplate,
    FieldTitle,
    FieldPeriodicRefresh,
  },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
    widgetId: {
      type: String,
      required: false,
    },
  },
  setup(props, { emit }) {
    const { updateModel } = useModelField(props, emit);
    const { pending, columns, fetchColumns } = useExternalDataTableColumns();

    const preparedColumns = computed(() => addPriorityColumnToColumnsArray(columns.value).map(column => ({
      ...column,

      text: column.name,
      value: column.name,
    })));

    const updateTable = async (tableId) => {
      await fetchColumns(tableId);

      const columnsWithKey = widgetColumnsToForm(preparedColumns.value);

      updateModel(setSeveralFields(props.form, {
        'parameters.table': tableId,
        'parameters.widgetColumns': columnsWithKey,
        'parameters.export_settings.widgetExportColumns': columnsWithKey,
        'parameters.sort.column': '',
      }));
    };

    onMounted(() => {
      if (props.form.parameters.table) {
        fetchColumns(props.form.parameters.table);
      }
    });

    return {
      pending,
      preparedColumns,
      updateTable,
    };
  },
};
</script>

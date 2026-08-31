<template>
  <widget-settings
    :submitting="submitting"
    :dirty="hasChanges"
    divider
    @submit="submit"
  >
    <c-layout column>
      <field-title v-model="form.title" />
      <field-periodic-refresh v-model="form.parameters.periodic_refresh" />
      <widget-settings-group :title="t('externalData.title')">
        <div class="pt-4 px-4">
          <external-data-table-table-field
            :value="form.parameters.table"
            name="parameters.table"
            required
            @update:value="updateTable"
          />
        </div>
      </widget-settings-group>
      <widget-settings-group
        :title="t('settings.advancedSettings')"
        :disabled="!form.parameters.table"
      >
        <c-progress-overlay :pending="pending" />
        <field-columns-without-template
          v-model="form.parameters.widgetColumns"
          :label="t('settings.columnNames')"
          :items="preparedColumns"
          without-custom-label
        />
        <field-resize-column-behavior
          v-model="form.parameters.columns"
          draggable
          resizable
        />
        <field-default-sort-column
          v-model="form.parameters.sort"
          :columns="form.parameters.widgetColumns"
          :columns-label="t('settings.columnName')"
          item-title="column"
          item-value="column"
        />
        <field-default-elements-per-page v-model="form.parameters.itemsPerPage" />
        <field-density v-model="form.parameters.dense" />
        <export-csv-form
          v-model="form.parameters"
          :items="preparedColumns"
          without-template
          without-infos-attributes
          without-custom-label
        />
      </widget-settings-group>
    </c-layout>
  </widget-settings>
</template>

<script setup>
import { isObject } from 'lodash';
import { computed, onMounted, toRef } from 'vue';

import { SIDE_BARS } from '@/constants';

import { addPriorityColumnToColumnsArray } from '@/helpers/entities/external-data-table/form';
import { widgetColumnsToForm } from '@/helpers/entities/widget/column/form';
import { setSeveralFields } from '@/helpers/immutable';

import { useI18n } from '@/hooks/i18n';
import { useWidgetSettings } from '@/hooks/widget/settings';

import ExternalDataTableTableField from '@/components/other/external-data-table/form/fields/external-data-table-table-field.vue';
import FieldResizeColumnBehavior from '@/components/sidebars/alarm/form/fields/resize-column-behavior.vue';
import ExportCsvForm from '@/components/sidebars/form/export-csv.vue';
import FieldColumnsWithoutTemplate from '@/components/sidebars/form/fields/columns-without-template.vue';
import FieldDefaultElementsPerPage from '@/components/sidebars/form/fields/default-elements-per-page.vue';
import FieldDefaultSortColumn from '@/components/sidebars/form/fields/default-sort-column.vue';
import FieldDensity from '@/components/sidebars/form/fields/density.vue';
import FieldPeriodicRefresh from '@/components/sidebars/form/fields/periodic-refresh.vue';
import FieldTitle from '@/components/sidebars/form/fields/title.vue';
import WidgetSettingsGroup from '@/components/sidebars/partials/widget-settings-group.vue';
import WidgetSettings from '@/components/sidebars/partials/widget-settings.vue';

import { useExternalDataTableColumns } from './form/hooks/external-data-table';

defineOptions({
  name: SIDE_BARS.externalDataTableSettings,
});

const props = defineProps({
  sidebar: {
    type: Object,
    required: true,
  },
});

const { t } = useI18n();
const sidebar = toRef(props, 'sidebar');

const { form, hasChanges, submitting, close, submit } = useWidgetSettings({
  sidebar,
});

const { pending, columns, fetchColumns } = useExternalDataTableColumns();

const getTableId = value => (isObject(value) ? value._id : value);

const preparedColumns = computed(() => addPriorityColumnToColumnsArray(columns.value).map(column => ({
    ...column,

    text: column.name,
    value: column.name,
  })),
);

const updateTable = async (table) => {
  const tableId = getTableId(table);

  if (!tableId) {
    return;
  }

  await fetchColumns(tableId);

  const columnsWithKey = widgetColumnsToForm(preparedColumns.value);

  Object.assign(
    form.value,
    setSeveralFields(form.value, {
      'parameters.table': tableId,
      'parameters.widgetColumns': columnsWithKey,
      'parameters.widgetExportColumns': columnsWithKey,
      'parameters.sort.column': '',
    }),
  );
};

onMounted(() => {
  const tableId = getTableId(form.value.parameters.table);

  if (tableId && form.value.parameters.table !== tableId) {
    form.value.parameters.table = tableId;
  }

  if (tableId) {
    fetchColumns(tableId);
  }
});

defineExpose({
  close,
});
</script>

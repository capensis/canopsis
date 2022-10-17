<template lang="pug">
  widget-settings(:submitting="submitting", @submit="submit")
    field-title(v-model="form.title")
    v-divider
    widget-settings-group(:title="$t('settings.advancedSettings')")
      field-default-sort-column(
        v-model="form.parameters.sort",
        :columns="form.parameters.widgetColumns",
        :columns-label="$t('settings.columnName')"
      )
      v-divider
      field-columns(
        v-model="form.parameters.widgetColumns",
        :label="$t('settings.columnNames')"
      )
      v-divider
      field-columns(
        v-model="form.parameters.serviceDependenciesColumns",
        :label="$t('settings.treeOfDependenciesColumnNames')",
        with-color-indicator
      )
      v-divider
      template(v-if="hasAccessToListFilters")
        field-filters(
          v-model="form.parameters.mainFilter",
          :filters.sync="form.filters",
          :widget-id="widget._id",
          :addable="hasAccessToAddFilter",
          :editable="hasAccessToEditFilter",
          with-alarm,
          with-entity,
          with-pbehavior
        )
        v-divider
      field-context-entities-types-filter(v-model="form.parameters.selectedTypes")
      v-divider
      export-csv-form(v-model="form.parameters")
    v-divider
</template>

<script>
import { SIDE_BARS } from '@/constants';

import { authMixin } from '@/mixins/auth';
import { widgetSettingsMixin } from '@/mixins/widget/settings';
import { permissionsWidgetsContextFilters } from '@/mixins/permissions/widgets/context/filters';

import FieldTitle from './fields/common/title.vue';
import FieldDefaultSortColumn from './fields/common/default-sort-column.vue';
import FieldColumns from './fields/common/columns.vue';
import FieldFilters from './fields/common/filters.vue';
import FieldContextEntitiesTypesFilter from './fields/context/context-entities-types-filter.vue';
import ExportCsvForm from './forms/export-csv.vue';
import WidgetSettings from './partials/widget-settings.vue';
import WidgetSettingsGroup from './partials/widget-settings-group.vue';

export default {
  name: SIDE_BARS.contextSettings,
  components: {
    FieldTitle,
    FieldDefaultSortColumn,
    FieldColumns,
    FieldFilters,
    FieldContextEntitiesTypesFilter,
    ExportCsvForm,
    WidgetSettings,
    WidgetSettingsGroup,
  },
  mixins: [
    authMixin,
    widgetSettingsMixin,
    permissionsWidgetsContextFilters,
  ],
};
</script>

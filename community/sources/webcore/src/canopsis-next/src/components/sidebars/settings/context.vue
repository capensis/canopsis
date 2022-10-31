<template lang="pug">
  div
    v-list.pt-0(expand)
      field-title(v-model="form.title", :title="$t('common.title')")
      v-divider
      v-list-group
        template(#activator="")
          v-list-tile {{ $t('settings.advancedSettings') }}
        v-list.grey.lighten-4.px-2.py-0(expand)
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
          field-columns(
            v-model="form.parameters.activeAlarmsColumns",
            :label="$t('settings.activeAlarmsColumns')"
          )
          v-divider
          field-columns(
            v-model="form.parameters.resolvedAlarmsColumns",
            :label="$t('settings.resolvedAlarmsColumns')"
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
    v-btn.primary(
      :loading="submitting",
      :disabled="submitting",
      @click="submit"
    ) {{ $t('common.save') }}
</template>

<script>
import { SIDE_BARS } from '@/constants';

import { authMixin } from '@/mixins/auth';
import { widgetSettingsMixin } from '@/mixins/widget/settings';
import { permissionsWidgetsContextFilters } from '@/mixins/permissions/widgets/context/filters';

import FieldTitle from '@/components/sidebars/settings/fields/common/title.vue';
import FieldDefaultSortColumn from '@/components/sidebars/settings/fields/common/default-sort-column.vue';
import FieldColumns from '@/components/sidebars/settings/fields/common/columns.vue';
import FieldFilters from '@/components/sidebars/settings/fields/common/filters.vue';
import FieldContextEntitiesTypesFilter from '@/components/sidebars/settings/fields/context/context-entities-types-filter.vue';
import ExportCsvForm from '@/components/sidebars/settings/forms/export-csv.vue';

export default {
  name: SIDE_BARS.contextSettings,
  components: {
    FieldTitle,
    FieldDefaultSortColumn,
    FieldColumns,
    FieldFilters,
    FieldContextEntitiesTypesFilter,
    ExportCsvForm,
  },
  mixins: [
    authMixin,
    widgetSettingsMixin,
    permissionsWidgetsContextFilters,
  ],
};
</script>

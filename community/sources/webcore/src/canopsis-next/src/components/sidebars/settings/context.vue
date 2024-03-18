<template lang="pug">
  widget-settings(:submitting="submitting", @submit="submit")
    field-title(v-model="form.title")
    v-divider
    widget-settings-group(:title="$t('settings.advancedSettings')")
      field-default-sort-column(
        v-model="form.parameters.sort",
        :columns="sortablePreparedWidgetColumns",
        :columns-label="$t('settings.columnName')"
      )
      v-divider
      field-columns(
        v-model="form.parameters.widgetColumns",
        :template="form.parameters.widgetColumnsTemplate",
        :templates="entityColumnsWidgetTemplates",
        :templates-pending="widgetTemplatesPending",
        :label="$t('settings.columnNames')",
        :type="$constants.ENTITIES_TYPES.entity",
        @update:template="updateWidgetColumnsTemplate"
      )
      v-divider
      field-columns(
        v-model="form.parameters.serviceDependenciesColumns",
        :template="form.parameters.serviceDependenciesColumnsTemplate",
        :templates="entityColumnsWidgetTemplates",
        :templates-pending="widgetTemplatesPending",
        :label="$t('settings.treeOfDependenciesColumnNames')",
        :type="$constants.ENTITIES_TYPES.entity",
        with-color-indicator,
        @update:template="updateServiceDependenciesColumnsTemplate"
      )
      v-divider
      field-columns(
        v-model="form.parameters.activeAlarmsColumns",
        :template="form.parameters.activeAlarmsColumnsTemplate",
        :templates="alarmColumnsWidgetTemplates",
        :templates-pending="widgetTemplatesPending",
        :label="$t('settings.activeAlarmsColumns')",
        :type="$constants.ENTITIES_TYPES.alarm",
        @update:template="updateActiveAlarmsColumnsTemplate"
      )
      v-divider
      field-columns(
        v-model="form.parameters.resolvedAlarmsColumns",
        :template="form.parameters.resolvedAlarmsColumnsTemplate",
        :templates="alarmColumnsWidgetTemplates",
        :templates-pending="widgetTemplatesPending",
        :label="$t('settings.resolvedAlarmsColumns')",
        :type="$constants.ENTITIES_TYPES.alarm",
        @update:template="updateResolvedAlarmsColumnsTemplate"
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
      field-grid-range-size(
        v-model="form.parameters.expandGridRangeSize",
        :title="$t('settings.expandGridRangeSizeContext')"
      )
      v-divider
      export-csv-form(
        v-model="form.parameters",
        :type="$constants.ENTITIES_TYPES.entity",
        :templates="entityColumnsWidgetTemplates",
        :templates-pending="widgetTemplatesPending"
      )
    v-divider
</template>

<script>
import {
  SIDE_BARS,
  ENTITY_UNSORTABLE_FIELDS,
  ENTITY_FIELDS_TO_LABELS_KEYS,
} from '@/constants';

import { formToWidgetColumns } from '@/helpers/forms/shared/widget-column';
import { getColumnLabel, getSortable } from '@/helpers/widgets';

import { authMixin } from '@/mixins/auth';
import { widgetSettingsMixin } from '@/mixins/widget/settings';
import { entitiesInfosMixin } from '@/mixins/entities/infos';
import { widgetTemplatesMixin } from '@/mixins/widget/templates';
import { permissionsWidgetsContextFilters } from '@/mixins/permissions/widgets/context/filters';

import FieldTitle from './fields/common/title.vue';
import FieldDefaultSortColumn from './fields/common/default-sort-column.vue';
import FieldColumns from './fields/common/columns.vue';
import FieldFilters from './fields/common/filters.vue';
import FieldContextEntitiesTypesFilter from './fields/context/context-entities-types-filter.vue';
import FieldGridRangeSize from './fields/common/grid-range-size.vue';
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
    FieldGridRangeSize,
    ExportCsvForm,
    WidgetSettings,
    WidgetSettingsGroup,
  },
  mixins: [
    authMixin,
    widgetSettingsMixin,
    entitiesInfosMixin,
    widgetTemplatesMixin,
    permissionsWidgetsContextFilters,
  ],
  computed: {
    preparedWidgetColumns() {
      return formToWidgetColumns(this.form.parameters.widgetColumns).map(column => ({
        ...column,

        text: getColumnLabel(column, ENTITY_FIELDS_TO_LABELS_KEYS),
      }));
    },

    sortablePreparedWidgetColumns() {
      return this.preparedWidgetColumns.filter(column => getSortable(column, ENTITY_UNSORTABLE_FIELDS));
    },
  },
  mounted() {
    this.fetchInfos();
  },
  methods: {
    updateWidgetColumnsTemplate(template, columns) {
      this.$set(this.form.parameters, 'widgetColumnsTemplate', template);
      this.$set(this.form.parameters, 'widgetColumns', columns);
    },

    updateServiceDependenciesColumnsTemplate(template, columns) {
      this.$set(this.form.parameters, 'serviceDependenciesColumnsTemplate', template);
      this.$set(this.form.parameters, 'serviceDependenciesColumns', columns);
    },

    updateActiveAlarmsColumnsTemplate(template, columns) {
      this.$set(this.form.parameters, 'activeAlarmsColumnsTemplate', template);
      this.$set(this.form.parameters, 'activeAlarmsColumns', columns);
    },

    updateResolvedAlarmsColumnsTemplate(template, columns) {
      this.$set(this.form.parameters, 'resolvedAlarmsColumnsTemplate', template);
      this.$set(this.form.parameters, 'resolvedAlarmsColumns', columns);
    },
  },
};
</script>

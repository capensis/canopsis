<template>
  <widget-settings
    :submitting="submitting"
    divider
    @submit="submit"
  >
    <field-title v-model="form.title" />
    <widget-settings-group :title="$t('settings.advancedSettings')">
      <field-default-sort-column
        v-model="form.parameters.sort"
        :columns="sortablePreparedWidgetColumns"
        :columns-label="$t('settings.columnName')"
      />
      <field-columns
        v-model="form.parameters.widgetColumns"
        :template="form.parameters.widgetColumnsTemplate"
        :templates="entityColumnsWidgetTemplates"
        :templates-pending="widgetTemplatesPending"
        :label="$t('settings.columnNames')"
        :type="$constants.ENTITIES_TYPES.entity"
        @update:template="updateWidgetColumnsTemplate"
      />
      <field-columns
        v-model="form.parameters.serviceDependenciesColumns"
        :template="form.parameters.serviceDependenciesColumnsTemplate"
        :templates="entityColumnsWidgetTemplates"
        :templates-pending="widgetTemplatesPending"
        :label="$t('settings.treeOfDependenciesColumnNames')"
        :type="$constants.ENTITIES_TYPES.entity"
        with-color-indicator
        @update:template="updateServiceDependenciesColumnsTemplate"
      />
      <field-tree-of-dependencies-settings v-model="form.parameters.treeOfDependenciesShowType" />
      <field-root-cause-settings v-model="form.parameters" />
      <field-columns
        v-model="form.parameters.activeAlarmsColumns"
        :template="form.parameters.activeAlarmsColumnsTemplate"
        :templates="alarmColumnsWidgetTemplates"
        :templates-pending="widgetTemplatesPending"
        :label="$t('settings.activeAlarmsColumns')"
        :type="$constants.ENTITIES_TYPES.alarm"
        @update:template="updateActiveAlarmsColumnsTemplate"
      />
      <field-columns
        v-model="form.parameters.resolvedAlarmsColumns"
        :template="form.parameters.resolvedAlarmsColumnsTemplate"
        :templates="alarmColumnsWidgetTemplates"
        :templates-pending="widgetTemplatesPending"
        :label="$t('settings.resolvedAlarmsColumns')"
        :type="$constants.ENTITIES_TYPES.alarm"
        @update:template="updateResolvedAlarmsColumnsTemplate"
      />
      <field-filters
        v-if="hasAccessToListFilters"
        v-model="form.parameters.mainFilter"
        :filters.sync="form.filters"
        :widget-id="widget._id"
        :addable="hasAccessToAddFilter"
        :editable="hasAccessToEditFilter"
        with-alarm
        with-entity
        with-pbehavior
        entity-counters-type
      />
      <field-context-entities-types-filter v-model="form.parameters.selectedTypes" />
      <export-csv-form
        v-model="form.parameters"
        :type="$constants.ENTITIES_TYPES.entity"
        :templates="entityColumnsWidgetTemplates"
        :templates-pending="widgetTemplatesPending"
        without-infos-attributes
      />
    </widget-settings-group>
    <charts-form v-model="form.parameters.charts" />
  </widget-settings>
</template>

<script>
import { SIDE_BARS, ENTITY_UNSORTABLE_FIELDS, ENTITY_FIELDS_TO_LABELS_KEYS } from '@/constants';

import { formToWidgetColumns } from '@/helpers/entities/widget/column/form';
import { getWidgetColumnLabel, getWidgetColumnSortable } from '@/helpers/entities/widget/list';

import { authMixin } from '@/mixins/auth';
import { widgetSettingsMixin } from '@/mixins/widget/settings';
import { entitiesInfosMixin } from '@/mixins/entities/infos';
import { widgetTemplatesMixin } from '@/mixins/widget/templates';
import { permissionsWidgetsContextFilters } from '@/mixins/permissions/widgets/context/filters';

import FieldTreeOfDependenciesSettings from '@/components/sidebars/form/fields/tree-of-dependencies-settings.vue';

import FieldRootCauseSettings from '../form/fields/root-cause-settings.vue';
import FieldTitle from '../form/fields/title.vue';
import FieldDefaultSortColumn from '../form/fields/default-sort-column.vue';
import FieldColumns from '../form/fields/columns.vue';
import FieldFilters from '../form/fields/filters.vue';
import ExportCsvForm from '../form/export-csv.vue';
import WidgetSettings from '../partials/widget-settings.vue';
import WidgetSettingsGroup from '../partials/widget-settings-group.vue';
import ChartsForm from '../chart/form/charts-form.vue';

import FieldContextEntitiesTypesFilter from './form/fields/context-entities-types-filter.vue';

export default {
  name: SIDE_BARS.contextSettings,
  components: {
    FieldRootCauseSettings,
    FieldTreeOfDependenciesSettings,
    FieldTitle,
    FieldDefaultSortColumn,
    FieldColumns,
    FieldFilters,
    FieldContextEntitiesTypesFilter,
    ExportCsvForm,
    WidgetSettings,
    WidgetSettingsGroup,
    ChartsForm,
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

        text: getWidgetColumnLabel(column, ENTITY_FIELDS_TO_LABELS_KEYS),
      }));
    },

    sortablePreparedWidgetColumns() {
      return this.preparedWidgetColumns.filter(column => getWidgetColumnSortable(column, ENTITY_UNSORTABLE_FIELDS));
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

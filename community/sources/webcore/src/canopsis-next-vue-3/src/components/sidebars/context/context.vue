<template>
  <widget-settings
    :submitting="submitting"
    :dirty="hasChanges"
    divider
    @submit="submit"
  >
    <field-title v-model="form.title" />
    <widget-settings-group :title="t('settings.advancedSettings')">
      <field-default-sort-column
        v-model="form.parameters.sort"
        :columns="sortablePreparedWidgetColumns"
        :columns-label="t('settings.columnName')"
      />
      <field-columns
        v-model="form.parameters.widgetColumns"
        :template="form.parameters.widgetColumnsTemplate"
        :templates="entityColumnsWidgetTemplates"
        :templates-pending="widgetTemplatesPending"
        :label="t('settings.columnNames')"
        :type="ENTITIES_TYPES.entity"
        @update:template="updateWidgetColumnsTemplate"
      />
      <field-columns
        v-model="form.parameters.serviceDependenciesColumns"
        :template="form.parameters.serviceDependenciesColumnsTemplate"
        :templates="entityColumnsWidgetTemplates"
        :templates-pending="widgetTemplatesPending"
        :label="t('settings.treeOfDependenciesColumnNames')"
        :type="ENTITIES_TYPES.entity"
        with-color-indicator
        @update:template="updateServiceDependenciesColumnsTemplate"
      />
      <field-tree-of-dependencies-settings v-model="form.parameters.treeOfDependenciesShowType" />
      <field-availability-graph-settings v-model="form.parameters.availability" />
      <field-root-cause-settings v-model="form.parameters" />
      <field-columns
        v-model="form.parameters.activeAlarmsColumns"
        :template="form.parameters.activeAlarmsColumnsTemplate"
        :templates="alarmColumnsWidgetTemplates"
        :templates-pending="widgetTemplatesPending"
        :label="t('settings.activeAlarmsColumns')"
        :type="ENTITIES_TYPES.alarm"
        @update:template="updateActiveAlarmsColumnsTemplate"
      />
      <field-columns
        v-model="form.parameters.resolvedAlarmsColumns"
        :template="form.parameters.resolvedAlarmsColumnsTemplate"
        :templates="alarmColumnsWidgetTemplates"
        :templates-pending="widgetTemplatesPending"
        :label="t('settings.resolvedAlarmsColumns')"
        :type="ENTITIES_TYPES.alarm"
        @update:template="updateResolvedAlarmsColumnsTemplate"
      />
      <field-filters
        v-if="hasAccessToFilter"
        v-model="form.parameters.mainFilter"
        v-model:filters="form.filters"
        :widget-id="widget._id"
        :widget-type="form.type"
        addable
        editable
        with-alarm
        with-entity
        with-pbehavior
        entity-counters-type
      />
      <field-context-entities-types-filter v-model="form.parameters.selectedTypes" />
      <field-grid-range-size
        v-model="form.parameters.expandGridRangeSize"
        :title="t('settings.expandGridRangeSize')"
      />
      <export-csv-form
        v-model="form.parameters"
        :type="ENTITIES_TYPES.entity"
        :templates="entityColumnsWidgetTemplates"
        :templates-pending="widgetTemplatesPending"
        without-infos-attributes
      />
      <field-switcher
        v-model="form.parameters.defaultInheritedPbehavior"
        :title="t('settings.defaultInheritedPbehavior')"
      />
      <field-switcher
        v-model="form.parameters.keepSelectedAfterAction"
        :title="t('common.massActionsPanel.keepSelectedAfterAction')"
      />
    </widget-settings-group>
    <charts-form v-model="form.parameters.charts" />
  </widget-settings>
</template>

<script setup>
import { computed, onMounted, ref, toRef } from 'vue';

import {
  COLUMNS_WIDGET_TEMPLATES_TYPES,
  ENTITIES_TYPES,
  ENTITY_FIELDS_TO_LABELS_KEYS,
  ENTITY_UNSORTABLE_FIELDS,
  MAX_LIMIT,
  QUICK_ACTIONS_WIDGET_TEMPLATES_TYPES,
  SIDE_BARS,
  USER_PERMISSIONS,
  WIDGET_TEMPLATES_TYPES,
} from '@/constants';

import { formToWidgetColumns, widgetColumnsToForm } from '@/helpers/entities/widget/column/form';
import { getWidgetColumnLabel, getWidgetColumnSortable } from '@/helpers/entities/widget/list';
import { widgetSortColumnsToForm } from '@/helpers/entities/widget/sort-column/form';

import { useCanPermission } from '@/hooks/auth';
import { useI18n } from '@/hooks/i18n';
import { useInfos } from '@/hooks/store/modules/infos';
import { useWidgetTemplate } from '@/hooks/store/modules/widget-template';
import { useWidgetSettings } from '@/hooks/widget/settings';

import FieldAvailabilityGraphSettings from '@/components/sidebars/availability/form/fields/availability-graph-settings.vue';
import ChartsForm from '@/components/sidebars/chart/form/charts-form.vue';
import ExportCsvForm from '@/components/sidebars/form/export-csv.vue';
import FieldColumns from '@/components/sidebars/form/fields/columns.vue';
import FieldDefaultSortColumn from '@/components/sidebars/form/fields/default-sort-column.vue';
import FieldFilters from '@/components/sidebars/form/fields/filters.vue';
import FieldGridRangeSize from '@/components/sidebars/form/fields/grid-range-size.vue';
import FieldRootCauseSettings from '@/components/sidebars/form/fields/root-cause-settings.vue';
import FieldSwitcher from '@/components/sidebars/form/fields/switcher.vue';
import FieldTitle from '@/components/sidebars/form/fields/title.vue';
import FieldTreeOfDependenciesSettings from '@/components/sidebars/form/fields/tree-of-dependencies-settings.vue';
import WidgetSettingsGroup from '@/components/sidebars/partials/widget-settings-group.vue';
import WidgetSettings from '@/components/sidebars/partials/widget-settings.vue';

import FieldContextEntitiesTypesFilter from './form/fields/context-entities-types-filter.vue';

defineOptions({
  name: SIDE_BARS.contextSettings,
});

const props = defineProps({
  sidebar: {
    type: Object,
    required: true,
  },
});

const { t } = useI18n();
const { hasAccess: hasAccessToFilter } = useCanPermission(
  USER_PERMISSIONS.business.context.actions.filter,
);
const sidebar = toRef(props, 'sidebar');
const widgetTemplates = ref([]);
const widgetTemplatesPending = ref(false);
const { fetchInfos } = useInfos();
const { fetchWidgetTemplatesListWithoutStore } = useWidgetTemplate();

const { form, widget, hasChanges, submitting, close, submit } = useWidgetSettings({
  sidebar,
});

const preparedWidgetTemplates = computed(() =>
  widgetTemplates.value.map((template) => {
    if (COLUMNS_WIDGET_TEMPLATES_TYPES.includes(template.type)) {
      return { ...template, columns: widgetColumnsToForm(template.columns) };
    }

    if (QUICK_ACTIONS_WIDGET_TEMPLATES_TYPES.includes(template.type)) {
      return template;
    }

    if (WIDGET_TEMPLATES_TYPES.alarmSortColumns === template.type) {
      return { ...template, sort_columns: widgetSortColumnsToForm(template.sort_columns) };
    }

    return template;
  }),
);

const alarmColumnsWidgetTemplates = computed(() =>
  preparedWidgetTemplates.value.filter(({ type }) => type === WIDGET_TEMPLATES_TYPES.alarmColumns),
);
const entityColumnsWidgetTemplates = computed(() =>
  preparedWidgetTemplates.value.filter(({ type }) => type === WIDGET_TEMPLATES_TYPES.entityColumns),
);

const preparedWidgetColumns = computed(() =>
  formToWidgetColumns(form.value.parameters.widgetColumns).map(column => ({
    ...column,

    text: getWidgetColumnLabel(column, ENTITY_FIELDS_TO_LABELS_KEYS),
  })),
);

const sortablePreparedWidgetColumns = computed(() =>
  preparedWidgetColumns.value.filter(column =>
    getWidgetColumnSortable(column, ENTITY_UNSORTABLE_FIELDS),
  ),
);

const updateWidgetColumnsTemplate = (template, columns) => {
  form.value.parameters.widgetColumnsTemplate = template;
  form.value.parameters.widgetColumns = columns;
};

const updateServiceDependenciesColumnsTemplate = (template, columns) => {
  form.value.parameters.serviceDependenciesColumnsTemplate = template;
  form.value.parameters.serviceDependenciesColumns = columns;
};

const updateActiveAlarmsColumnsTemplate = (template, columns) => {
  form.value.parameters.activeAlarmsColumnsTemplate = template;
  form.value.parameters.activeAlarmsColumns = columns;
};

const updateResolvedAlarmsColumnsTemplate = (template, columns) => {
  form.value.parameters.resolvedAlarmsColumnsTemplate = template;
  form.value.parameters.resolvedAlarmsColumns = columns;
};

const fetchWidgetTemplates = async () => {
  try {
    widgetTemplatesPending.value = true;

    const { data } = await fetchWidgetTemplatesListWithoutStore({ params: { limit: MAX_LIMIT } });

    widgetTemplates.value = data;
  } catch (err) {
    console.warn(err);
  } finally {
    widgetTemplatesPending.value = false;
  }
};

onMounted(() => {
  fetchInfos();
  fetchWidgetTemplates();
});

defineExpose({
  close,
});
</script>

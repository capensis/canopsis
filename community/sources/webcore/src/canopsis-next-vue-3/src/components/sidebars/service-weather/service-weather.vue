<template>
  <widget-settings
    :submitting="submitting"
    :dirty="hasChanges"
    divider
    @submit="submit"
  >
    <field-title v-model="form.title" />
    <field-periodic-refresh v-model="form.parameters.periodic_refresh" />
    <field-filters
      v-if="hasAccessToFilter"
      v-model="form.parameters.mainFilter"
      v-model:filters="form.filters"
      :widget-id="widget._id"
      :widget-type="form.type"
      :entity-types="[ENTITY_TYPES.service]"
      addable
      editable
      with-entity
      with-service-weather
      entity-counters-type
    />
    <alarms-list-modal-form
      v-model:value="form.parameters.alarmsList"
      :templates="preparedWidgetTemplates"
      :templates-pending="widgetTemplatesPending"
    />
    <widget-settings-item :title="t('settings.limit')">
      <c-number-field v-model="form.parameters.limit" />
    </widget-settings-item>
    <widget-settings-item :title="t('settings.colorIndicator.title')">
      <c-color-indicator-field v-model="form.parameters.colorIndicator" hide-details />
    </widget-settings-item>
    <widget-settings-group :title="t('context.treeOfDependencies')">
      <field-columns
        v-model="form.parameters.serviceDependenciesColumns"
        :template="form.parameters.serviceDependenciesColumnsTemplate"
        :templates="entityColumnsWidgetTemplates"
        :templates-pending="widgetTemplatesPending"
        :label="t('settings.columnNames')"
        :type="ENTITIES_TYPES.entity"
        with-color-indicator
        @update:template="updateServiceDependenciesColumnsTemplate"
      />
      <field-tree-of-dependencies-settings v-model="form.parameters.treeOfDependenciesShowType" />
    </widget-settings-group>
    <widget-settings-group :title="t('settings.advancedSettings')">
      <field-default-sort-column
        v-model="form.parameters.sort"
        :columns="sortColumns"
        :columns-label="t('settings.orderBy')"
        item-title="label"
        item-value="value"
      />
      <field-default-elements-per-page
        v-model="form.parameters.modalItemsPerPage"
        :sub-title="t('settings.modal')"
      />
      <field-root-cause-settings v-model="form.parameters" />
      <field-text-editor-with-template
        v-model="form.parameters.blockTemplate"
        :template="form.parameters.blockTemplateTemplate"
        :templates="weatherItemWidgetTemplates"
        :widget-template-type="WIDGET_TEMPLATES_TYPES.weatherItem"
        :variables="blockTemplateVariables"
        :title="t('settings.weatherTemplate')"
        @update:template="updateBlockTemplate"
      />
      <field-text-editor-with-template
        v-model="form.parameters.modalTemplate"
        :template="form.parameters.modalTemplateTemplate"
        :templates="weatherModalWidgetTemplates"
        :widget-template-type="WIDGET_TEMPLATES_TYPES.weatherModal"
        :variables="blockTemplateVariables"
        :title="t('settings.modalTemplate')"
        @update:template="updateModalTemplate"
      />
      <field-text-editor-with-template
        v-model="form.parameters.entityTemplate"
        :template="form.parameters.entityTemplateTemplate"
        :templates="weatherEntityWidgetTemplates"
        :widget-template-type="WIDGET_TEMPLATES_TYPES.weatherEntity"
        :variables="entityTemplateVariables"
        :title="t('settings.entityTemplate')"
        @update:template="updateEntityTemplate"
      />
      <field-grid-size
        v-model="form.parameters.columnMobile"
        :title="t('settings.columnMobile')"
        mobile
      />
      <field-grid-size
        v-model="form.parameters.columnTablet"
        :title="t('settings.columnTablet')"
        tablet
      />
      <field-grid-size
        v-model="form.parameters.columnDesktop"
        :title="t('settings.columnDesktop')"
      />
      <margins-form v-model="form.parameters.margin" />
      <field-slider
        v-model="form.parameters.heightFactor"
        :title="t('settings.height')"
        :min="1"
        :max="20"
      />
      <field-counters-selector v-model="form.parameters.counters" />
      <field-switcher
        v-model="form.parameters.isPriorityEnabled"
        :title="t('settings.isPriorityEnabled')"
      />
      <field-switcher
        v-model="form.parameters.isHideGrayEnabled"
        :title="t('settings.isHideGrayEnabled')"
      />
      <field-switcher
        v-model="form.parameters.isSecondaryIconEnabled"
        :title="t('settings.isSecondaryIconEnabled')"
      />
      <field-switcher
        v-model="form.parameters.defaultInheritedPbehavior"
        :title="t('settings.defaultInheritedPbehavior')"
      />
      <field-modal-type v-model="form.parameters.modalType" />
      <widget-settings-group :title="t('settings.stateSettings')">
        <field-action-required-settings v-model="form.parameters.actionRequiredSettings" />
        <field-no-action-required-settings v-model="form.parameters.noActionRequiredSettings" />
      </widget-settings-group>
      <v-divider />
      <field-switcher
        v-model="form.parameters.entitiesActionsInQueue"
        :title="t('settings.entitiesActionsInQueue')"
      />
    </widget-settings-group>
  </widget-settings>
</template>

<script setup>
import { computed, onMounted, ref, toRef } from 'vue';

import {
  COLUMNS_WIDGET_TEMPLATES_TYPES,
  ENTITIES_TYPES,
  ENTITY_FIELDS,
  ENTITY_FIELDS_TO_LABELS_KEYS,
  ENTITY_TEMPLATE_FIELDS,
  ENTITY_TYPES,
  MAX_LIMIT,
  QUICK_ACTIONS_WIDGET_TEMPLATES_TYPES,
  SIDE_BARS,
  USER_PERMISSIONS,
  VARIABLES_STEP_FIELDS,
  WIDGET_TEMPLATES_TYPES,
} from '@/constants';

import { widgetColumnsToForm } from '@/helpers/entities/widget/column/form';
import { widgetSortColumnsToForm } from '@/helpers/entities/widget/sort-column/form';

import { useCanPermission } from '@/hooks/auth';
import { useI18n } from '@/hooks/i18n';
import { useInfos } from '@/hooks/store/modules/infos';
import { useWidgetTemplate } from '@/hooks/store/modules/widget-template';
import { useWidgetSettings } from '@/hooks/widget/settings';

import AlarmsListModalForm from '@/components/sidebars/alarm/form/alarms-list-modal.vue';
import FieldColumns from '@/components/sidebars/form/fields/columns.vue';
import FieldDefaultElementsPerPage from '@/components/sidebars/form/fields/default-elements-per-page.vue';
import FieldDefaultSortColumn from '@/components/sidebars/form/fields/default-sort-column.vue';
import FieldFilters from '@/components/sidebars/form/fields/filters.vue';
import FieldGridSize from '@/components/sidebars/form/fields/grid-size.vue';
import FieldPeriodicRefresh from '@/components/sidebars/form/fields/periodic-refresh.vue';
import FieldRootCauseSettings from '@/components/sidebars/form/fields/root-cause-settings.vue';
import FieldSlider from '@/components/sidebars/form/fields/slider.vue';
import FieldSwitcher from '@/components/sidebars/form/fields/switcher.vue';
import FieldTextEditorWithTemplate from '@/components/sidebars/form/fields/text-editor-with-template.vue';
import FieldTitle from '@/components/sidebars/form/fields/title.vue';
import FieldTreeOfDependenciesSettings from '@/components/sidebars/form/fields/tree-of-dependencies-settings.vue';
import MarginsForm from '@/components/sidebars/form/margins.vue';
import WidgetSettingsGroup from '@/components/sidebars/partials/widget-settings-group.vue';
import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';
import WidgetSettings from '@/components/sidebars/partials/widget-settings.vue';

import FieldActionRequiredSettings from './form/fields/action-required-settings.vue';
import FieldCountersSelector from './form/fields/counters-selector.vue';
import FieldModalType from './form/fields/modal-type.vue';
import FieldNoActionRequiredSettings from './form/fields/no-action-required-settings.vue';

defineOptions({
  name: SIDE_BARS.serviceWeatherSettings,
});

const props = defineProps({
  sidebar: {
    type: Object,
    required: true,
  },
});

const { t } = useI18n();
const { hasAccess: hasAccessToFilter } = useCanPermission(
  USER_PERMISSIONS.business.serviceWeather.actions.filter,
);
const sidebar = toRef(props, 'sidebar');
const widgetTemplates = ref([]);
const widgetTemplatesPending = ref(false);
const { entityInfos, fetchInfos } = useInfos();
const { fetchWidgetTemplatesListWithoutStore } = useWidgetTemplate();

const { form, widget, hasChanges, submitting, close, submit } = useWidgetSettings({
  sidebar,
});

const sortColumns = computed(() => [
  { label: t('common.name'), value: ENTITY_FIELDS.name },
  { label: t('common.state'), value: ENTITY_FIELDS.state },
]);

const infosSubVariables = computed(() => [
  {
    text: t('common.value'),
    value: VARIABLES_STEP_FIELDS.value,
  },
  {
    text: t('common.description'),
    value: 'description',
  },
]);

const infosVariables = computed(() =>
  entityInfos.value.map(({ value }) => ({
    text: value,
    value,
    variables: infosSubVariables.value,
  })),
);

const entityVariables = computed(() => [
  { value: ENTITY_TEMPLATE_FIELDS.id },
  { value: ENTITY_TEMPLATE_FIELDS.name },
  {
    value: ENTITY_TEMPLATE_FIELDS.infos,
    variables: infosVariables.value.length ? infosVariables.value : undefined,
  },
  { value: ENTITY_TEMPLATE_FIELDS.connector },
  { value: ENTITY_TEMPLATE_FIELDS.connectorName },
  { value: ENTITY_TEMPLATE_FIELDS.component },
  { value: ENTITY_TEMPLATE_FIELDS.resource },
  { value: ENTITY_TEMPLATE_FIELDS.state },
  { value: ENTITY_TEMPLATE_FIELDS.status },
  { value: ENTITY_TEMPLATE_FIELDS.snooze },
  { value: ENTITY_TEMPLATE_FIELDS.ack },
  { value: ENTITY_TEMPLATE_FIELDS.lastUpdateDate },
  { value: ENTITY_TEMPLATE_FIELDS.impactLevel },
  { value: ENTITY_TEMPLATE_FIELDS.impactState },
  { value: ENTITY_TEMPLATE_FIELDS.categoryName },
  { value: ENTITY_TEMPLATE_FIELDS.pbehaviorInfo },
  { value: ENTITY_TEMPLATE_FIELDS.alarmCreationDate },
  { value: ENTITY_TEMPLATE_FIELDS.ticket },
  { value: ENTITY_TEMPLATE_FIELDS.statsOk },
  { value: ENTITY_TEMPLATE_FIELDS.statsKo },
  { value: ENTITY_TEMPLATE_FIELDS.alarmDisplayName },
  { value: ENTITY_TEMPLATE_FIELDS.links },
].map(variable => ({
  ...variable,
  text: t(ENTITY_FIELDS_TO_LABELS_KEYS[variable.value.replace('entity.', '')], 2),
})));

const blockTemplateVariables = computed(() => {
  const excludeFields = [
    ENTITY_TEMPLATE_FIELDS.ticket,
    ENTITY_TEMPLATE_FIELDS.statsKo,
    ENTITY_TEMPLATE_FIELDS.statsOk,
    ENTITY_TEMPLATE_FIELDS.alarmDisplayName,
    ENTITY_TEMPLATE_FIELDS.alarmCreationDate,
  ];

  return entityVariables.value.filter(({ value }) => !excludeFields.includes(value));
});

const entityTemplateVariables = computed(() => [
  ...entityVariables.value,
  {
    value: ENTITY_TEMPLATE_FIELDS.sourceType,
    text: t(ENTITY_FIELDS_TO_LABELS_KEYS[ENTITY_FIELDS.sourceType]),
  },
]);

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

const weatherItemWidgetTemplates = computed(() =>
  preparedWidgetTemplates.value.filter(({ type }) => type === WIDGET_TEMPLATES_TYPES.weatherItem),
);
const weatherModalWidgetTemplates = computed(() =>
  preparedWidgetTemplates.value.filter(({ type }) => type === WIDGET_TEMPLATES_TYPES.weatherModal),
);
const weatherEntityWidgetTemplates = computed(() =>
  preparedWidgetTemplates.value.filter(({ type }) => type === WIDGET_TEMPLATES_TYPES.weatherEntity),
);
const entityColumnsWidgetTemplates = computed(() =>
  preparedWidgetTemplates.value.filter(({ type }) => type === WIDGET_TEMPLATES_TYPES.entityColumns),
);

const updateServiceDependenciesColumnsTemplate = (template, columns) => {
  form.value.parameters.serviceDependenciesColumnsTemplate = template;
  form.value.parameters.serviceDependenciesColumns = columns;
};

const updateBlockTemplate = (template) => {
  if (template && template !== form.value.parameters.blockTemplateTemplate) {
    form.value.parameters.blockTemplateTemplate = template;
  }
};

const updateModalTemplate = (template) => {
  if (template && template !== form.value.parameters.modalTemplateTemplate) {
    form.value.parameters.modalTemplateTemplate = template;
  }
};

const updateEntityTemplate = (template) => {
  if (template && template !== form.value.parameters.entityTemplateTemplate) {
    form.value.parameters.entityTemplateTemplate = template;
  }
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

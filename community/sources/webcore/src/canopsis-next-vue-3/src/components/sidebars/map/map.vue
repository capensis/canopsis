<template>
  <widget-settings
    :submitting="submitting"
    :dirty="hasChanges"
    divider
    @submit="submit"
  >
    <field-title v-model="form.title" />
    <field-periodic-refresh v-model="form.parameters.periodic_refresh" />
    <field-map v-model="form.parameters.map" />

    <widget-settings-group :title="t('settings.entityDisplaySettings')">
      <widget-settings-item :title="colorIndicatorTitle">
        <c-color-indicator-field v-model="form.parameters.color_indicator" hide-details />
      </widget-settings-item>

      <field-switcher
        v-model="form.parameters.entities_under_pbehavior_enabled"
        :title="t('settings.entitiesUnderPbehaviorEnabled')"
      />
    </widget-settings-group>

    <widget-settings-group :title="t('settings.advancedSettings')">
      <field-filters
        v-if="hasAccessToFilter"
        v-model="form.parameters.mainFilter"
        v-model:filters="form.filters"
        :widget-id="widget._id"
        :widget-type="form.type"
        addable
        editable
        with-entity
      />

      <field-text-editor
        v-model="form.parameters.entity_info_template"
        :title="t('settings.entityInfoPopup')"
        :variables="entityInfoTemplateVariables"
      />

      <field-columns
        v-model="form.parameters.alarmsColumns"
        :template="form.parameters.alarmsColumnsTemplate"
        :templates="alarmColumnsWidgetTemplates"
        :templates-pending="widgetTemplatesPending"
        :label="t('settings.alarmsColumns')"
        :type="ENTITIES_TYPES.alarm"
        with-template
        with-html
        @update:template="updateAlarmsColumnsTemplate"
      />

      <field-columns
        v-model="form.parameters.entitiesColumns"
        :template="form.parameters.entitiesColumnsTemplate"
        :templates="entityColumnsWidgetTemplates"
        :templates-pending="widgetTemplatesPending"
        :label="t('settings.entitiesColumns')"
        :type="ENTITIES_TYPES.entity"
        with-color-indicator
        @update:template="updateEntitiesColumnsTemplate"
      />
    </widget-settings-group>
  </widget-settings>
</template>

<script setup>
import { computed, onMounted, ref, toRef } from 'vue';

import {
  ENTITIES_TYPES,
  ENTITY_FIELDS_TO_LABELS_KEYS,
  ENTITY_TEMPLATE_FIELDS,
  MAX_LIMIT,
  SIDE_BARS,
  USER_PERMISSIONS,
  VARIABLES_STEP_FIELDS,
  WIDGET_TEMPLATES_TYPES,
} from '@/constants';

import { prepareWidgetTemplates } from '@/helpers/entities/widget/template/prepare';

import { useCanPermission } from '@/hooks/auth';
import { useI18n } from '@/hooks/i18n';
import { useInfos } from '@/hooks/store/modules/infos';
import { useWidgetTemplate } from '@/hooks/store/modules/widget-template';
import { useWidgetSettings } from '@/hooks/widget/settings';

import FieldColumns from '@/components/sidebars/form/fields/columns.vue';
import FieldFilters from '@/components/sidebars/form/fields/filters.vue';
import FieldPeriodicRefresh from '@/components/sidebars/form/fields/periodic-refresh.vue';
import FieldSwitcher from '@/components/sidebars/form/fields/switcher.vue';
import FieldTextEditor from '@/components/sidebars/form/fields/text-editor.vue';
import FieldTitle from '@/components/sidebars/form/fields/title.vue';
import WidgetSettingsGroup from '@/components/sidebars/partials/widget-settings-group.vue';
import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';
import WidgetSettings from '@/components/sidebars/partials/widget-settings.vue';

import FieldMap from './form/fields/map.vue';

defineOptions({
  name: SIDE_BARS.mapSettings,
});

const props = defineProps({
  sidebar: {
    type: Object,
    required: true,
  },
});

const { t } = useI18n();
const { hasAccess: hasAccessToFilter } = useCanPermission(
  USER_PERMISSIONS.business.map.actions.filter,
);
const sidebar = toRef(props, 'sidebar');
const widgetTemplates = ref([]);
const widgetTemplatesPending = ref(false);
const { entityInfos, fetchInfos } = useInfos();
const { fetchWidgetTemplatesListWithoutStore } = useWidgetTemplate();

const { form, widget, hasChanges, submitting, close, submit } = useWidgetSettings({
  sidebar,
});

const alarmColumnsWidgetTemplates = computed(() =>
  preparedWidgetTemplates.value.filter(({ type }) => type === WIDGET_TEMPLATES_TYPES.alarmColumns),
);
const entityColumnsWidgetTemplates = computed(() =>
  preparedWidgetTemplates.value.filter(({ type }) => type === WIDGET_TEMPLATES_TYPES.entityColumns),
);
const colorIndicatorTitle = computed(() => t('settings.colorIndicator.title'));

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

const entityInfoTemplateVariables = computed(() => {
  const excludedFields = [
    ENTITY_TEMPLATE_FIELDS.alarmLastComment,
    ENTITY_TEMPLATE_FIELDS.tags,
    ENTITY_TEMPLATE_FIELDS.lastComment,
  ];

  return [
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
  ]
    .filter(({ value }) => !excludedFields.includes(value))
    .map(variable => ({
      ...variable,
      text: t(ENTITY_FIELDS_TO_LABELS_KEYS[variable.value.replace('entity.', '')], 2),
    }));
});

const preparedWidgetTemplates = computed(() => prepareWidgetTemplates(widgetTemplates.value));

const updateAlarmsColumnsTemplate = (template, columns) => {
  form.value.parameters.alarmsColumnsTemplate = template;
  form.value.parameters.alarmsColumns = columns;
};

const updateEntitiesColumnsTemplate = (template, columns) => {
  form.value.parameters.entitiesColumnsTemplate = template;
  form.value.parameters.entitiesColumns = columns;
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

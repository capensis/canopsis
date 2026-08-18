<template>
  <widget-settings
    :submitting="submitting"
    :dirty="hasChanges"
    divider
    @submit="submit"
  >
    <field-title v-model="form.title" />

    <field-filters
      v-model:filters="form.filters"
      :widget-id="widget._id"
      :widget-type="form.type"
      addable
      editable
      with-alarm
      with-entity
      with-pbehavior
      hide-selector
    />

    <field-opened-resolved-filter v-model="form.parameters.opened" />

    <alarms-list-modal-form
      v-model:value="form.parameters.alarmsList"
      :templates="preparedWidgetTemplates"
      :templates-pending="widgetTemplatesPending"
    />

    <widget-settings-group :title="t('settings.advancedSettings')">
      <field-template
        v-model="form.parameters.blockTemplate"
        :title="t('settings.blockTemplate')"
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

      <counter-levels-form v-model="form.parameters.levels" />

      <field-switcher
        v-model="form.parameters.isCorrelationEnabled"
        :title="t('settings.isCorrelationEnabled')"
      />
    </widget-settings-group>
  </widget-settings>
</template>

<script setup>
import { computed, onMounted, ref, toRef } from 'vue';

import {
  COLUMNS_WIDGET_TEMPLATES_TYPES,
  MAX_LIMIT,
  QUICK_ACTIONS_WIDGET_TEMPLATES_TYPES,
  SIDE_BARS,
  WIDGET_TEMPLATES_TYPES,
} from '@/constants';

import { widgetColumnsToForm } from '@/helpers/entities/widget/column/form';
import { widgetSortColumnsToForm } from '@/helpers/entities/widget/sort-column/form';

import { useI18n } from '@/hooks/i18n';
import { useInfos } from '@/hooks/store/modules/infos';
import { useWidgetTemplate } from '@/hooks/store/modules/widget-template';
import { useWidgetSettings } from '@/hooks/widget/settings';

import AlarmsListModalForm from '@/components/sidebars/alarm/form/alarms-list-modal.vue';
import FieldOpenedResolvedFilter from '@/components/sidebars/alarm/form/fields/opened-resolved-filter.vue';
import FieldFilters from '@/components/sidebars/form/fields/filters.vue';
import FieldGridSize from '@/components/sidebars/form/fields/grid-size.vue';
import FieldSlider from '@/components/sidebars/form/fields/slider.vue';
import FieldSwitcher from '@/components/sidebars/form/fields/switcher.vue';
import FieldTitle from '@/components/sidebars/form/fields/title.vue';
import MarginsForm from '@/components/sidebars/form/margins.vue';
import WidgetSettingsGroup from '@/components/sidebars/partials/widget-settings-group.vue';
import WidgetSettings from '@/components/sidebars/partials/widget-settings.vue';

import CounterLevelsForm from './form/counter-levels.vue';
import FieldTemplate from './form/fields/template.vue';

defineOptions({
  name: SIDE_BARS.counterSettings,
});

const props = defineProps({
  sidebar: {
    type: Object,
    required: true,
  },
});

const { t } = useI18n();
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

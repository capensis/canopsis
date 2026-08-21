<template>
  <c-layout column>
    <c-name-field
      v-model="title"
      :label="t('common.name')"
      :max-length="255"
      name="title"
      autofocus
      required
    />

    <widget-template-columns-form v-if="isColumnsType" />

    <widget-template-quick-actions-form
      v-else-if="isAlarmQuickActionsType"
      :massive="isMassiveQuickActionsType"
    />

    <widget-template-sort-columns-form
      v-else-if="isAlarmSortColumnsType"
      v-model="sortColumns"
      required
    />

    <widget-template-text-form v-else :entity-infos="entityInfos" />
  </c-layout>
</template>

<script>
import { computed, onMounted } from 'vue';

import {
  ALARM_QUICK_ACTIONS_WIDGET_TEMPLATE_TYPES,
  COLUMNS_WIDGET_TEMPLATES_TYPES,
  WIDGET_TEMPLATES_TYPES,
} from '@/constants';

import { useFormContext, useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useInfos } from '@/hooks/store/modules/infos';

import WidgetTemplateColumnsForm from './widget-template-columns-form.vue';
import WidgetTemplateQuickActionsForm from './widget-template-quick-actions-form.vue';
import WidgetTemplateSortColumnsForm from './widget-template-sort-columns-form.vue';
import WidgetTemplateTextForm from './widget-template-text-form.vue';

export default {
  components: {
    WidgetTemplateColumnsForm,
    WidgetTemplateQuickActionsForm,
    WidgetTemplateSortColumnsForm,
    WidgetTemplateTextForm,
  },

  setup() {
    const { t } = useI18n();
    const form = useFormContext();
    const { entityInfos, fetchInfos } = useInfos();
    const { title } = useFormModels(['title']);
    const sortColumns = form.model('sort_columns');
    const type = form.model('type');

    const isColumnsType = computed(() => COLUMNS_WIDGET_TEMPLATES_TYPES.includes(type.value));
    const isAlarmQuickActionsType = computed(() =>
      ALARM_QUICK_ACTIONS_WIDGET_TEMPLATE_TYPES.includes(type.value),
    );
    const isMassiveQuickActionsType = computed(
      () => type.value === WIDGET_TEMPLATES_TYPES.alarmMassQuickActions,
    );
    const isAlarmSortColumnsType = computed(
      () => type.value === WIDGET_TEMPLATES_TYPES.alarmSortColumns,
    );

    onMounted(() => fetchInfos({ withRules: true }));

    return {
      t,
      title,
      sortColumns,
      entityInfos,
      isColumnsType,
      isAlarmQuickActionsType,
      isMassiveQuickActionsType,
      isAlarmSortColumnsType,
    };
  },
};
</script>

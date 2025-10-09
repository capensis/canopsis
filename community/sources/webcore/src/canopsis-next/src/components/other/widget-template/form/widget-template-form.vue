<template>
  <v-layout column>
    <v-text-field
      v-field="form.title"
      v-validate="'required'"
      :label="$t('common.name')"
      :error-messages="errors.collect('title')"
      name="title"
      autofocus
    />
    <widget-template-columns-form
      v-if="isColumnsType"
      v-field="form"
    />
    <widget-template-quick-actions-form
      v-else-if="isAlarmQuickActionsType"
      v-field="form"
      :massive="isMassiveQuickActionsType"
    />
    <widget-template-text-form
      v-else
      v-field="form"
      :entity-infos="entityInfos"
    />
  </v-layout>
</template>

<script>
import { computed, onMounted } from 'vue';

import {
  COLUMNS_WIDGET_TEMPLATES_TYPES,
  ALARM_QUICK_ACTIONS_WIDGET_TEMPLATE_TYPES,
  WIDGET_TEMPLATES_TYPES,
} from '@/constants';

import { useEntityInfos } from '@/hooks/store/modules/entity-infos';

import WidgetTemplateColumnsForm from './widget-template-columns-form.vue';
import WidgetTemplateTextForm from './widget-template-text-form.vue';
import WidgetTemplateQuickActionsForm from './widget-template-quick-actions-form.vue';

export default {
  inject: ['$validator'],
  components: {
    WidgetTemplateColumnsForm,
    WidgetTemplateTextForm,
    WidgetTemplateQuickActionsForm,
  },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { fetchInfos, entityInfos } = useEntityInfos();

    const isColumnsType = computed(() => COLUMNS_WIDGET_TEMPLATES_TYPES.includes(props.form.type));
    const isAlarmQuickActionsType = computed(() => ALARM_QUICK_ACTIONS_WIDGET_TEMPLATE_TYPES.includes(props.form.type));
    const isMassiveQuickActionsType = computed(() => props.form.type === WIDGET_TEMPLATES_TYPES.alarmMassQuickActions);

    onMounted(() => fetchInfos({ withRules: true }));

    return {
      entityInfos,
      isColumnsType,
      isAlarmQuickActionsType,
      isMassiveQuickActionsType,
    };
  },
};
</script>

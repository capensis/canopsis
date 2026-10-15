<template>
  <v-container>
    <c-layout gap="2" column>
      <span class="text-body-2">
        {{ t('settings.quickActions.title') }}
      </span>
      <span>{{ t('settings.quickActions.description') }}</span>
      <c-widget-template-field
        :model-value="template"
        :templates="templates"
        :pending="templatesPending"
        @update:model-value="updateTemplate"
      />
      <quick-alarm-actions-form
        :actions="modelValue"
        :massive="massive"
        @update:actions="updateValue"
      />
    </c-layout>
  </v-container>
</template>

<script>
import { useI18n } from '@/hooks/i18n';
import { useWidgetTemplateField } from '@/hooks/widget/widget-template';

import QuickAlarmActionsForm from '@/components/common/actions-panel/quick-alarm-actions-form.vue';

export default {
  components: { QuickAlarmActionsForm },

  props: {
    modelValue: {
      type: Array,
      default: () => [],
    },
    massive: {
      type: Boolean,
      default: false,
    },
    template: {
      type: [String, Symbol],
      default: '',
    },
    templates: {
      type: Array,
      default: () => [],
    },
    templatesPending: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:modelValue', 'update:template'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const { updateTemplate, updateValue } = useWidgetTemplateField(props, 'actions', emit);

    return {
      t,
      updateTemplate,
      updateValue,
    };
  },
};
</script>

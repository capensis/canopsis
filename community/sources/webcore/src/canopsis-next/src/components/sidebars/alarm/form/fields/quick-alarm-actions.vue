<template>
  <v-container>
    <v-layout class="gap-2" column>
      <span class="text-body-2">{{ $t('settings.quickActions.title') }}</span>
      <span>{{ $t('settings.quickActions.description') }}</span>
      <c-widget-template-field
        :value="template"
        :templates="templates"
        :pending="templatesPending"
        @input="updateTemplate"
      />
      <quick-alarm-actions-form
        :actions="value"
        :massive="massive"
        @input="updateValue"
      />
    </v-layout>
  </v-container>
</template>

<script>
import { useWidgetTemplateField } from '@/hooks/widget/widget-template';

import QuickAlarmActionsForm from '@/components/common/actions-panel/quick-alarm-actions-form.vue';

export default {
  components: { QuickAlarmActionsForm },
  props: {
    value: {
      type: Array,
      default: () => [],
    },
    massive: {
      type: Boolean,
      default: false,
    },
    title: {
      type: String,
      default: '',
    },
    description: {
      type: String,
      default: '',
    },
    template: {
      type: String,
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
  setup(props, { emit }) {
    const { updateTemplate, updateValue } = useWidgetTemplateField(props, 'actions', emit);

    return {
      updateTemplate,
      updateValue,
    };
  },
};
</script>

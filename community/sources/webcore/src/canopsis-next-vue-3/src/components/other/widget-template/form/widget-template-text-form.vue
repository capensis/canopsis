<template>
  <text-editor-field
    v-model="content"
    :variables="variables"
    class="mt-3"
    name="content"
    required
  />
</template>

<script>
import { computed } from 'vue';

import { WIDGET_TEMPLATES_TYPES } from '@/constants';

import {
  prepareAlarmTemplateVariables,
  prepareEntityTemplateVariables,
} from '@/helpers/entities/widget/template/variables';

import { useFormContext } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import TextEditorField from '@/components/forms/fields/text-editor-field.vue';

export default {
  components: { TextEditorField },

  props: {
    entityInfos: {
      type: Array,
      default: () => [],
    },
  },

  setup(props) {
    const { t } = useI18n();
    const form = useFormContext();
    const content = form.model('content');
    const type = form.model('type');

    const alarmVariables = computed(() => prepareAlarmTemplateVariables(t));
    const entityVariables = computed(() => prepareEntityTemplateVariables(t, props.entityInfos));
    const variables = computed(() =>
      type.value === WIDGET_TEMPLATES_TYPES.alarmMoreInfos
        ? alarmVariables.value
        : entityVariables.value,
    );

    return {
      content,
      variables,
    };
  },
};
</script>

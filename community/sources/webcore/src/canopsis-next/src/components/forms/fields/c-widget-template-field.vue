<template>
  <v-select
    v-field="value"
    :items="templatesWithCustom"
    :label="$tc('common.template', 1)"
    :loading="pending"
    :clearable="clearable"
    return-object
  />
</template>

<script>
import { computed } from 'vue';

import { CUSTOM_WIDGET_TEMPLATE } from '@/constants';

import { widgetTemplateToForm } from '@/helpers/entities/widget/template/form';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    value: {
      type: [String, Symbol],
      required: false,
    },
    templates: {
      type: Array,
      default: () => [],
    },
    pending: {
      type: Boolean,
      default: false,
    },
    valueKey: {
      type: String,
      default: 'value',
    },
    clearable: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const { t } = useI18n();

    const templatesWithCustom = computed(() => [
      {
        ...widgetTemplateToForm(),
        value: CUSTOM_WIDGET_TEMPLATE,
        text: t('common.custom'),
      },

      ...props.templates.map(item => ({
        ...item,

        value: item._id,
        text: item.title,
      })),
    ]);

    return {
      templatesWithCustom,
    };
  },
};
</script>

<template>
  <pbehavior-general-form
    :form="modelValue"
    :name="name"
    :name-label="t('common.namePrefix')"
    :name-tooltip="t('scenario.tooltips.pbehaviorActionsNamePrefix')"
    :depth="depth"
    with-start-on-trigger
    block-child
    no-comments
    @update:form="$emit('update:modelValue', $event)"
  />
</template>

<script>
import { useI18n } from '@/hooks/i18n';

import PbehaviorGeneralForm from '@/components/other/pbehavior/pbehaviors/form/pbehavior-general-form.vue';

export default {
  components: {
    PbehaviorGeneralForm,
  },

  props: {
    modelValue: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      default: 'parameters',
    },
    depth: {
      type: Number,
      default: 0,
    },
  },

  emits: ['update:modelValue'],

  setup() {
    const { t } = useI18n();

    return {
      t,
    };
  },
};
</script>

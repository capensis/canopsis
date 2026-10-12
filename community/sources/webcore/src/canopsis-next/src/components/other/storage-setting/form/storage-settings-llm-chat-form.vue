<template>
  <c-layout column>
    <c-information-block :title="$t('storageSetting.llmChat.title')">
      <c-enabled-duration-field
        :duration="form.delete_after"
        :label="$t('storageSetting.llmChat.deleteHistory')"
        :suffix="$t('storageSetting.olderThan')"
        name="llm_chat.delete_after"
        switcher
        hide-value-on-false
        @update:duration="updateField('delete_after', $event)"
      />
      <c-enabled-duration-field
        :duration="form.delete_token_usage_after"
        :label="$t('storageSetting.llmChat.deleteStats')"
        :suffix="$t('storageSetting.olderThan')"
        name="llm_chat.delete_token_usage_after"
        switcher
        hide-value-on-false
        @update:duration="updateField('delete_token_usage_after', $event)"
      />
    </c-information-block>
  </c-layout>
</template>

<script>
import { useModelField } from '@/hooks/form';

export default {
  props: {
    form: {
      type: Object,
      required: true,
    },
  },

  emits: ['update:form'],

  setup(props, { emit }) {
    const { updateField } = useModelField(props, emit);

    return { updateField };
  },
};
</script>

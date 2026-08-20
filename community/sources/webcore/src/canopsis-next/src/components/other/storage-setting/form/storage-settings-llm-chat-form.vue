<template>
  <v-layout column>
    <c-information-block :title="$t('storageSetting.llmChat.title')">
      <c-enabled-duration-field
        v-field="form.delete_after"
        :label="$t('storageSetting.llmChat.deleteHistory')"
        :suffix="$t('storageSetting.olderThan')"
        name="llm_chat.delete_after"
        switcher
        hide-value-on-false
      />
    </c-information-block>
  </v-layout>
</template>

<script>
export default {
  inject: ['$validator'],
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
};
</script>

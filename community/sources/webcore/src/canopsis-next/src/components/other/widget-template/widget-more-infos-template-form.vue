<template lang="pug">
  v-layout(column)
    v-text-field.mb-2(
      v-field="form.title",
      v-validate="'required'",
      :label="$t('common.name')",
      :error-messages="errors.collect('title')",
      name="title"
    )
    text-editor-field(
      v-model="form.text",
      v-validate="'required'",
      :error-messages="errors.collect('text')",
      :variables="alarmVariables",
      name="text"
    )
</template>

<script>
import { alarmVariablesMixin } from '@/mixins/widget/variables/alarm';

import TextEditorField from '@/components/common/text-editor/text-editor.vue';

export default {
  inject: ['$validator'],
  components: { TextEditorField },
  mixins: [alarmVariablesMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>

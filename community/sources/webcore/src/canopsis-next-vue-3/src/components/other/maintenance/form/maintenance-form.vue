<template>
  <c-form-block>
    <c-form-block-row :label="$t('maintenance.message')">
      <text-editor-field
        v-model="messageField"
        :label="$t('maintenance.message')"
        name="message"
        autofocus
        public
        required
      />
    </c-form-block-row>

    <c-form-block-row :label="$t('common.color')" align-center>
      <c-color-picker-field v-model="colorField" hide-details />
    </c-form-block-row>
  </c-form-block>
</template>

<script>
import { useFormModels } from '@/hooks/form';

import TextEditorField from '@/components/forms/fields/text-editor-field.vue';

export default {
  components: { TextEditorField },

  setup() {
    const { messageField, colorField } = useFormModels({
      messageField: 'message',
      colorField: 'color',
    });

    return {
      messageField,
      colorField,
    };
  },
};
</script>

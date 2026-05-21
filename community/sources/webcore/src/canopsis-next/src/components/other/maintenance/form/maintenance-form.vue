<template>
  <c-form-block>
    <c-form-block-row
      :label="$t('maintenance.message')"
      indented
    >
      <text-editor-field
        v-field="form.message"
        v-validate="'required'"
        :label="$t('maintenance.message')"
        :error-messages="errors.collect('message')"
        name="message"
        autofocus
        public
      />
    </c-form-block-row>

    <c-form-block-row
      :label="$t('common.color')"
      align-center
    >
      <c-color-picker-field v-field="form.color" />
    </c-form-block-row>
  </c-form-block>
</template>

<script>
import TextEditorField from '@/components/forms/fields/text-editor-field.vue';

export default {
  inject: ['$validator'],
  components: { TextEditorField },
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

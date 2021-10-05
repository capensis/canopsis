<template lang="pug">
  div(data-test="addInfoPopupLayout")
    v-select(
      v-field="form.selectedColumn",
      v-validate="'required'",
      :items="columns",
      :error-messages="errors.collect('column')",
      item-text="label",
      name="column",
      return-object
    )
    text-editor-field(
      v-field="form.template",
      v-validate="'required'",
      :error-messages="errors.collect('template')",
      name="template"
    )
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
    columns: {
      type: Array,
      default: () => [],
    },
  },
};
</script>

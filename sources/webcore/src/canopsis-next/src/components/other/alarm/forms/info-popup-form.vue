<template lang="pug">
  div(data-test="addInfoPopupLayout")
    v-select(
      v-field="form.column",
      v-validate="'required'",
      :items="columns",
      :error-messages="errors.collect('column')",
      item-text="label",
      name="column"
    )
    text-editor(
      v-field="form.template",
      v-validate="'required'",
      :error-messages="errors.collect('template')",
      name="template"
    )
</template>

<script>
import TextEditor from '@/components/other/text-editor/text-editor.vue';

export default {
  inject: ['$validator'],
  components: { TextEditor },
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

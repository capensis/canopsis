<template>
  <v-layout column>
    <v-select
      v-field="form.column"
      v-validate="'required'"
      :items="columns"
      :label="$tc('common.column', 1)"
      :error-messages="errors.collect('column')"
      name="column"
    />
    <text-editor-field
      v-field="form.template"
      v-validate="'required'"
      :error-messages="errors.collect('template')"
      name="template"
    />
  </v-layout>
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

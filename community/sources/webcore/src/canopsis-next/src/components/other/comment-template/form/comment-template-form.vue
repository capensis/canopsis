<template>
  <v-layout column>
    <c-name-field
      v-field="form.name"
      :label="$t('common.name')"
      :max-length="255"
      name="name"
      required
    />
    <v-layout
      class="mt-3"
      column
    >
      <comment-template-fields-form v-field="form.fields" />
    </v-layout>
  </v-layout>
</template>

<script>
import CommentTemplateFieldsForm from './comment-template-fields-form.vue';

export default {
  components: { CommentTemplateFieldsForm },
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

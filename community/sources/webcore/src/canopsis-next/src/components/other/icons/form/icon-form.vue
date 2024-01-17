<template>
  <v-layout column>
    <c-name-field
      v-field="form.title"
      :label="$t('common.title')"
      :error-messages="errors.collect('title')"
      name="title"
      required
    />
    <icon-file-selector
      v-field="form.file"
    />
  </v-layout>
</template>

<script>
import IconFileSelector from './fields/icon-file-selector.vue';

export default {
  inject: ['$validator'],
  components: {
    IconFileSelector,
  },
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

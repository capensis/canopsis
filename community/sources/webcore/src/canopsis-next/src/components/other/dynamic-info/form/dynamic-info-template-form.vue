<template>
  <v-layout class="gap-2" column>
    <v-text-field
      v-field="template.title"
      v-validate="'required'"
      :error-messages="errors.collect('title')"
      :label="$t('common.title')"
      name="title"
      autofocus
    />
    <dynamic-info-template-names-field v-field="template.names" />
  </v-layout>
</template>

<script>
import DynamicInfoTemplateNamesField from './fields/dynamic-info-template-names-field.vue';

export default {
  inject: ['$validator'],
  components: { DynamicInfoTemplateNamesField },
  model: {
    prop: 'template',
    event: 'input',
  },
  props: {
    template: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>

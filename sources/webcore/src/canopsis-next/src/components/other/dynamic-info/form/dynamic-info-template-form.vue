<template lang="pug">
  div
    v-text-field(
      v-field="template.title",
      v-validate="'required'",
      :error-messages="errors.collect('title')",
      :label="$t('common.title')",
      name="title"
    )
    h3 {{ $t('modals.createDynamicInfoTemplate.fields.names') }}
    dynamic-info-template-names-field(v-field="template.names")
      v-alert(slot="no-data", :value="true", type="info") {{ $t('modals.createDynamicInfoTemplate.emptyNames') }}
</template>

<script>
import DynamicInfoTemplateNamesField from './partials/dynamic-info-template-names-field.vue';

export default {
  components: { DynamicInfoTemplateNamesField },
  inject: ['$validator'],
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

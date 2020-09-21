<template lang="pug">
  div
    v-layout
      v-text-field(
        v-field="form._id",
        :label="$t('modals.createDynamicInfo.steps.general.fields.id')",
        name="id"
      )
    v-layout
      v-text-field(
        v-validate="'required'",
        v-field="form.name",
        :label="$t('modals.createDynamicInfo.steps.general.fields.name')",
        :error-messages="errors.collect('name')",
        name="name"
      )
    v-layout
      v-textarea(
        v-field="form.description",
        :label="$t('modals.createDynamicInfo.steps.general.fields.description')",
        :error-messages="errors.collect('description')",
        name="description"
      )
</template>

<script>
import formValidationHeaderMixin from '@/mixins/form/validation-header';

export default {
  mixins: [formValidationHeaderMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  inject: ['$validator'],
  props: {
    form: {
      type: Object,
      required: true,
    },
  },
};
</script>

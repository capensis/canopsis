<template lang="pug">
  div
    v-layout
      v-text-field(
        v-field="form._id",
        :label="$t('modals.createDynamicInfo.steps.general.fields.id')",
        :disabled="isDisabledIdField",
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
      disable-during-periods-field(v-field="form.disable_during_periods")
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

import DisableDuringPeriodsField from '@/components/forms/fields/disable-during-periods.vue';

export default {
  components: { DisableDuringPeriodsField },
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
    isDisabledIdField: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

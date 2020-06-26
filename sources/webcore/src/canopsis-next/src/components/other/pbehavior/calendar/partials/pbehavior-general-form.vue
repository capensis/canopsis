<template lang="pug">
  div
    v-layout(data-test="pbehaviorTypeLayout", row)
      v-flex(xs6)
        v-combobox(
          v-field="form.reason",
          v-validate="'required'",
          :label="$t('modals.createPbehavior.steps.general.fields.reason')",
          :loading="pbehaviorReasonsPending",
          :items="reasons",
          :error-messages="errors.collect('reason')",
          name="reason",
          data-test="pbehaviorReason"
        )
      v-flex(xs6)
        v-select.ml-3(
          v-field="form.type_",
          v-validate="'required'",
          :label="$t('modals.createPbehavior.steps.general.fields.type')",
          :items="types",
          :error-messages="errors.collect('type')",
          name="type"
        )
</template>

<script>
import { PBEHAVIOR_TYPES, PAUSE_REASONS } from '@/constants';

import formValidationHeaderMixin from '@/mixins/form/validation-header';
import formMixin from '@/mixins/form';
import pbehaviorReasonsMixin from '@/mixins/entities/pbehavior-reasons';

export default {
  mixins: [formMixin, formValidationHeaderMixin, pbehaviorReasonsMixin],
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
  computed: {
    reasons() {
      return this.pbehaviorReasons.length ? this.pbehaviorReasons : Object.values(PAUSE_REASONS);
    },

    types() {
      return Object.values(PBEHAVIOR_TYPES);
    },
  },
  mounted() {
    this.fetchPbehaviorReasons();
  },
};
</script>

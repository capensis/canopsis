<template lang="pug">
  div
    v-layout(wrap)
      v-flex(xs12)
        v-text-field(
          v-field="form.name",
          v-validate="'required'",
          :label="$t('modals.createPbehavior.steps.general.fields.name')",
          :error-messages="errors.collect('name')",
          name="name"
        )
      v-flex(xs12)
        v-switch(
          v-field="form.enabled",
          :label="$t('modals.createPbehavior.steps.general.fields.enabled')",
          color="primary",
          hide-details
        )
      v-flex.mt-3(xs12)
        v-layout(wrap, justify-space-between)
          v-flex(xs6)
            date-time-picker-field(
              v-validate="tstartRules",
              :value="form.tstart",
              :label="$t('modals.createPbehavior.steps.general.fields.start')",
              name="tstart",
              @input="updateField('tstart', $event)"
            )
          v-flex(xs6)
            date-time-picker-field(
              v-validate="tstopRules",
              :value="form.tstop",
              :label="$t('modals.createPbehavior.steps.general.fields.stop')",
              name="tstop",
              @input="updateField('tstop', $event)"
            )
      v-flex(xs12)
        pbehavior-reasons-field(v-field="form.reason")
      v-flex(xs12)
        v-select(
          v-field="form.type",
          v-validate="'required'",
          :label="$t('modals.createPbehavior.steps.general.fields.type')",
          :items="types",
          :error-messages="errors.collect('type')",
          name="type"
        )
</template>

<script>
import moment from 'moment-timezone';

import { PBEHAVIOR_TYPES, DATETIME_FORMATS } from '@/constants';

import formMixin from '@/mixins/form';
import formValidationHeaderMixin from '@/mixins/form/validation-header';
import pbehaviorReasonsMixin from '@/mixins/entities/pbehavior/reasons';

import DateTimePickerField from '@/components/forms/fields/date-time-picker/date-time-picker-field.vue';
import PbehaviorReasonsField from '@/components/other/pbehavior/reasons/partials/pbehavior-reasons-field.vue';

export default {
  components: { PbehaviorReasonsField, DateTimePickerField },
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
    tstartRules() {
      return {
        required: true,
        date_format: DATETIME_FORMATS.veeValidateDateTimeFormat,
      };
    },

    tstopRules() {
      const rules = { required: true };

      if (this.form.tstart) {
        rules.after = [moment(this.form.tstart).format(DATETIME_FORMATS.dateTimePicker)];
        rules.date_format = DATETIME_FORMATS.veeValidateDateTimeFormat;
      }

      return rules;
    },

    types() {
      return Object.values(PBEHAVIOR_TYPES);
    },
  },
  mounted() {
    this.fetchPbehaviorReasonsList();
  },
};
</script>

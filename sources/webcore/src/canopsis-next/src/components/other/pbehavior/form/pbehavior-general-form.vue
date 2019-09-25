<template lang="pug">
  div
    v-divider
    h3.my-3.grey--text {{ $t('modals.createPbehavior.steps.general.general') }}
    v-divider
    v-layout
      v-switch(
        v-model="form.enabled",
        :label="$t('modals.createPbehavior.steps.general.fields.enabled')",
        color="primary",
        hide-details
      )
    v-layout
      v-text-field(
        v-validate="'required'",
        v-model="form.name",
        :label="$t('modals.createPbehavior.steps.general.fields.name')",
        :error-messages="errors.collect('name')",
        name="name"
      )
    v-layout
      v-combobox(
        v-validate="'required'",
        v-model="form.reason",
        :label="$t('modals.createPbehavior.steps.general.fields.reason')",
        :items="reasons",
        :error-messages="errors.collect('reason')",
        name="reason"
      )
      v-select.ml-3(
        v-validate="'required'",
        v-model="form.type_",
        :label="$t('modals.createPbehavior.steps.general.fields.type')",
        :items="types",
        :error-messages="errors.collect('type')",
        name="type"
      )
    v-divider
    h3.my-3.grey--text {{ $t('modals.createPbehavior.steps.general.dates') }}
    v-divider
    v-layout.mt-3(wrap)
      v-flex(xs12)
        v-layout(wrap, justify-space-between)
          v-flex(xs4)
            date-time-picker-field(
              v-validate="tstartRules",
              v-model="form.tstart",
              :label="$t('modals.createPbehavior.steps.general.fields.start')",
              name="tstart"
            )
          v-flex(xs4)
            date-time-picker-field(
              v-validate="tstopRules",
              v-model="form.tstop",
              :label="$t('modals.createPbehavior.steps.general.fields.stop')",
              name="tstop"
            )
          v-flex(xs3)
            v-select(
              :items="timezones",
              v-model="form.timezone",
              :label="$t('modals.createPbehavior.steps.general.fields.timezone')"
            )
</template>

<script>
import moment from 'moment-timezone';

import { PAUSE_REASONS, PBEHAVIOR_TYPES, DATETIME_FORMATS } from '@/constants';

import formValidationHeaderMixin from '@/mixins/form/validation-header';

import DateTimePickerField from '@/components/forms/fields/date-time-picker/date-time-picker-field.vue';

export default {
  components: {
    DateTimePickerField,
  },
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
  computed: {
    reasons() {
      return Object.values(PAUSE_REASONS);
    },

    types() {
      return Object.values(PBEHAVIOR_TYPES);
    },

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

    timezones() {
      return moment.tz.names();
    },
  },
};
</script>

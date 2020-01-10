<template lang="pug">
  div
    v-divider
    h3.my-3.grey--text {{ $t('modals.createPbehavior.steps.general.general') }}
    v-divider
    v-layout
      v-switch(
        v-field="form.enabled",
        :label="$t('modals.createPbehavior.steps.general.fields.enabled')",
        color="primary",
        hide-details
      )
    v-layout
      v-text-field(
        v-field="form.name",
        v-validate="'required'",
        :label="$t('modals.createPbehavior.steps.general.fields.name')",
        :error-messages="errors.collect('name')",
        name="name",
        data-test="pbehaviorFormName"
      )
    v-layout(data-test="pbehaviorTypeLayout", row)
      v-combobox(
        v-field="form.reason",
        v-validate="'required'",
        :label="$t('modals.createPbehavior.steps.general.fields.reason')",
        :loading="pbehaviorReasonsPending",
        :items="pbehaviorReasons",
        :error-messages="errors.collect('reason')",
        name="reason",
        data-test="pbehaviorReason"
      )
      v-select.ml-3(
        v-field="form.type_",
        v-validate="'required'",
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
          v-flex(data-test="startDateTimePicker", xs4)
            date-time-picker-field(
              v-field="form.tstart",
              v-validate="tstartRules",
              :label="$t('modals.createPbehavior.steps.general.fields.start')",
              name="tstart"
            )
          v-flex(data-test="stopDateTimePicker", xs4)
            date-time-picker-field(
              v-field="form.tstop",
              v-validate="tstopRules",
              :label="$t('modals.createPbehavior.steps.general.fields.stop')",
              name="tstop"
            )
          v-flex(xs3)
            v-autocomplete(
              v-field="form.timezone",
              v-validate="'required'",
              :items="timezones",
              :label="$t('modals.createPbehavior.steps.general.fields.timezone')",
              name="timezone"
            )
</template>

<script>
import moment from 'moment-timezone';

import { PBEHAVIOR_TYPES, DATETIME_FORMATS } from '@/constants';

import formValidationHeaderMixin from '@/mixins/form/validation-header';
import formMixin from '@/mixins/form';
import pbehaviorReasonsMixin from '@/mixins/entities/pbehavior-reasons';

import DateTimePickerField from '@/components/forms/fields/date-time-picker/date-time-picker-field.vue';

export default {
  components: {
    DateTimePickerField,
  },
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
  mounted() {
    this.fetchPbehaviorReasons();
  },
};
</script>

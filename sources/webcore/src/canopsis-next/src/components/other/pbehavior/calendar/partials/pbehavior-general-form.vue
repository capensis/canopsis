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
              :clearable="hasPauseType",
              name="tstop",
              @input="updateField('tstop', $event)"
            )
      v-flex(xs12)
        v-combobox(
          v-field="form.reason",
          v-validate="'required'",
          :label="$t('modals.createPbehavior.steps.general.fields.reason')",
          :loading="pbehaviorReasonsPending",
          :items="reasons",
          :error-messages="errors.collect('reason')",
          name="reason"
        )
      v-flex(xs12)
        pbehavior-type-field(v-field="form.type")
</template>

<script>
import { get } from 'lodash';
import moment from 'moment-timezone';

import { PAUSE_REASONS, DATETIME_FORMATS, PBEHAVIOR_TYPE_TYPES } from '@/constants';

import formMixin from '@/mixins/form';
import formValidationHeaderMixin from '@/mixins/form/validation-header';
import entitiesPbehaviorReasonsMixin from '@/mixins/entities/pbehavior-reasons';

import DateTimePickerField from '@/components/forms/fields/date-time-picker/date-time-picker-field.vue';
import PbehaviorTypeField from '@/components/other/pbehavior/calendar/partials/pbehavior-type-field.vue';

export default {
  components: { PbehaviorTypeField, DateTimePickerField },
  mixins: [
    formMixin,
    formValidationHeaderMixin,
    entitiesPbehaviorReasonsMixin,
  ],
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
    hasPauseType() {
      return get(this.form.type, 'type') === PBEHAVIOR_TYPE_TYPES.pause;
    },

    tstartRules() {
      return {
        required: true,
        date_format: DATETIME_FORMATS.veeValidateDateTimeFormat,
      };
    },

    tstopRules() {
      const rules = { required: !this.hasPauseType };

      if (this.form.tstart) {
        rules.after = [moment(this.form.tstart).format(DATETIME_FORMATS.dateTimePicker)];
        rules.date_format = DATETIME_FORMATS.veeValidateDateTimeFormat;
      }

      return rules;
    },

    reasons() {
      return this.pbehaviorReasons.length ? this.pbehaviorReasons : Object.values(PAUSE_REASONS);
    },
  },
  mounted() {
    this.fetchPbehaviorReasons();
  },
};
</script>

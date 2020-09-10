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
        enabled-field(v-field="form.enabled", hide-details)
      v-flex.mt-3(xs12)
        v-layout(row)
          date-time-range-picker-field(
            v-field="form",
            :startLabel="$t('modals.createPbehavior.steps.general.fields.start')",
            :endLabel="$t('modals.createPbehavior.steps.general.fields.stop')",
            :startRules="tstartRules",
            :endRules="tstopRules",
            :noEnding="noEnding",
            :fullDay="fullDay"
          )
        v-layout(wrap)
          v-checkbox.mt-0(
            v-model="fullDay",
            :label="$t('modals.createPbehavior.steps.general.fields.fullDay')",
            color="primary",
            hide-details
          )
        v-layout(wrap)
          v-checkbox.mt-0.mb-2(
            v-if="hasPauseType",
            v-model="noEnding",
            :label="$t('modals.createPbehavior.steps.general.fields.noEnding')",
            color="primary",
            hide-details
          )
      v-flex(xs12)
        pbehavior-reasons-field(v-field="form.reason")
      v-flex(xs12)
        pbehavior-type-field(v-field="form.type")
</template>

<script>
import { get } from 'lodash';
import moment from 'moment-timezone';

import { DATETIME_FORMATS, PBEHAVIOR_TYPE_TYPES } from '@/constants';

import { isStartOfDay, isEndOfDay } from '@/helpers/date';

import formMixin from '@/mixins/form';
import formValidationHeaderMixin from '@/mixins/form/validation-header';
import entitiesPbehaviorReasonsMixin from '@/mixins/entities/pbehavior/reasons';

import PbehaviorTypeField from '@/components/other/pbehavior/calendar/partials/pbehavior-type-field.vue';
import PbehaviorReasonsField from '@/components/other/pbehavior/reasons/partials/pbehavior-reasons-field.vue';
import DateTimeRangePickerField from '@/components/forms/fields/date-time-range-picker-field.vue';
import EnabledField from '@/components/forms/fields/enabled-field.vue';

export default {
  components: {
    EnabledField,
    DateTimeRangePickerField,
    PbehaviorReasonsField,
    PbehaviorTypeField,
  },
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
  data() {
    const noEnding = !this.form.tstop;

    return {
      noEnding,

      fullDay: isStartOfDay(this.form.tstart) && (noEnding || isEndOfDay(this.form.tstop)),
    };
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
  },
  watch: {
    noEnding(value) {
      if (value) {
        this.updateField('tstop', null);
      } else {
        const unit = this.fullDay ? 'day' : 'hour';
        const tstopMoment = moment(this.form.tstart).add(1, unit);

        if (this.fullDay) {
          tstopMoment.endOf(unit);
        }

        this.updateField('tstop', tstopMoment.toDate());
      }
    },
    fullDay() {
      const tstartMoment = moment(this.form.tstart).startOf('day');

      this.updateField('tstart', tstartMoment.toDate());

      if (!this.noEnding) {
        const tstopMoment = moment(this.form.tstop).endOf('day');

        this.updateField('tstop', tstopMoment.toDate());
      }
    },
    hasPauseType(value) {
      if (!value) {
        this.noEnding = false;
      }
    },
  },
};
</script>

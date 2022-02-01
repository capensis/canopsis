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
      v-flex(v-if="!noEnabled", xs12)
        c-enabled-field(v-field="form.enabled", hide-details)
      v-flex.mt-3(xs12)
        v-layout(v-if="withStartOnTrigger", wrap)
          v-switch.mt-0.mb-1(
            v-model="form.start_on_trigger",
            :label="$t('modals.createPbehavior.steps.general.fields.startOnTrigger')",
            color="primary",
            hide-details,
            @change="changeStartOnTrigger"
          )
        v-layout(v-if="form.start_on_trigger", row)
          c-duration-field(v-field="form.duration", required)
        template(v-else)
          v-layout(row)
            date-time-splitted-range-picker-field(
              :start="form.tstart",
              :end="form.tstop",
              :startLabel="$t('modals.createPbehavior.steps.general.fields.start')",
              :endLabel="$t('modals.createPbehavior.steps.general.fields.stop')",
              :startRules="tstartRules",
              :endRules="tstopRules",
              :noEnding="noEnding",
              :fullDay="fullDay",
              @update:start="updateField('tstart', $event)",
              @update:end="updateField('tstop', $event)"
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
        pbehavior-reasons-field(
          v-field="form.reason",
          v-validate="'required'"
        )
      v-flex(xs12)
        pbehavior-type-field(
          v-field="form.type",
          v-validate="'required'",
          return-object
        )
</template>

<script>
import { get } from 'lodash';

import { DATETIME_FORMATS, PBEHAVIOR_TYPE_TYPES } from '@/constants';

import {
  isStartOfDay,
  isEndOfDay,
  convertDateToStartOfDayDateObject,
  convertDateToEndOfDayDateObject,
  convertDateToString,
  addUnitToDate,
  convertDateToDateObject,
  convertDateToEndOfUnitDateObject,
} from '@/helpers/date/date';

import { formMixin, formValidationHeaderMixin } from '@/mixins/form';
import entitiesPbehaviorReasonsMixin from '@/mixins/entities/pbehavior/reasons';

import DateTimeSplittedRangePickerField from '@/components/forms/fields/date-time-splitted-range-picker-field.vue';
import PbehaviorTypeField from '@/components/other/pbehavior/calendar/partials/pbehavior-type-field.vue';
import PbehaviorReasonsField from '@/components/other/pbehavior/reasons/partials/pbehavior-reasons-field.vue';

export default {
  inject: ['$validator'],
  components: {
    DateTimeSplittedRangePickerField,
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
  props: {
    form: {
      type: Object,
      required: true,
    },
    noEnabled: {
      type: Boolean,
      default: false,
    },
    withStartOnTrigger: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    const noEnding = this.form.tstart && !this.form.tstop;

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
        rules.after = [convertDateToString(this.form.tstart, DATETIME_FORMATS.dateTimePicker)];
        rules.date_format = DATETIME_FORMATS.veeValidateDateTimeFormat;
      }

      return rules;
    },
  },
  watch: {
    noEnding(noEnding) {
      const { tstart } = this.form;

      if (noEnding) {
        this.updateField('tstop', null);
      } else if (tstart) {
        const unit = this.fullDay ? 'day' : 'hour';

        const tstop = addUnitToDate(tstart, 1, unit);
        const tstopDate = this.fullDay
          ? convertDateToEndOfUnitDateObject(tstop, unit)
          : convertDateToDateObject(tstop);

        this.updateField('tstop', tstopDate);
      }
    },
    fullDay() {
      const { tstart, tstop } = this.form;

      if (tstart) {
        this.updateField('tstart', convertDateToStartOfDayDateObject(tstart));

        if (!this.noEnding && tstop) {
          this.updateField('tstop', convertDateToEndOfDayDateObject(tstop));
        }
      }
    },
    hasPauseType(value) {
      if (!value) {
        this.noEnding = false;
      }
    },
  },
  methods: {
    changeStartOnTrigger(value) {
      if (value) {
        this.fullDay = false;
        this.noEnding = false;

        this.updateModel({
          ...this.form,

          tstart: null,
          tstop: null,
        });
      } else {
        this.removeField('duration');
      }
    },
  },
};
</script>

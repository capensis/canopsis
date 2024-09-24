<template lang="pug">
  v-layout(column)
    c-name-field(v-field="form.name", :label="nameLabel", required)
    c-enabled-field(v-if="!noEnabled", v-field="form.enabled", hide-details)
    v-flex.mt-3(xs12)
      c-enabled-field.mt-0.mb-1(
        v-if="withStartOnTrigger",
        :value="form.start_on_trigger",
        :label="$t('modals.createPbehavior.steps.general.fields.startOnTrigger')",
        hide-details,
        @input="updateStartOnTrigger"
      )
      c-duration-field(
        v-if="form.start_on_trigger",
        v-field="form.duration",
        required
      )
      template(v-else)
        date-time-splitted-range-picker-field(
          :start="form.tstart",
          :end="form.tstop",
          :start-label="$t('modals.createPbehavior.steps.general.fields.start')",
          :end-label="$t('modals.createPbehavior.steps.general.fields.stop')",
          :start-rules="tstartRules",
          :end-rules="tstopRules",
          :end-min="tstopMin",
          :end-max="tstopMax",
          :no-ending="noEnding",
          :full-day="fullDay",
          @update:start="updateField('tstart', $event)",
          @update:end="updateTStop"
        )
        v-checkbox.mt-0(
          v-model="fullDay",
          :label="$t('modals.createPbehavior.steps.general.fields.fullDay')",
          color="primary",
          hide-details
        )
        v-checkbox.mt-0.mb-2(
          v-if="hasPauseType",
          v-model="noEnding",
          :label="$t('modals.createPbehavior.steps.general.fields.noEnding')",
          color="primary",
          hide-details
        )
    c-pbehavior-reason-field(v-field="form.reason", required, return-object)
    c-pbehavior-type-field(v-field="form.type", required, return-object)
</template>

<script>
import { get } from 'lodash';

import { MAX_PBEHAVIOR_DATES_DIFF_YEARS } from '@/config';

import { DATETIME_FORMATS, PBEHAVIOR_TYPE_TYPES, TIME_UNITS } from '@/constants';

import {
  isStartOfDay,
  isEndOfDay,
  addUnitToDate,
  getNowTimestamp,
  convertDateToString,
  convertDateToTimestamp,
  convertDateToDateObject,
  convertDateToStartOfDayDateObject,
  convertDateToEndOfUnitDateObject,
  convertDateToEndOfDayDateObject,
} from '@/helpers/date/date';

import { formMixin, formValidationHeaderMixin } from '@/mixins/form';
import { entitiesPbehaviorReasonMixin } from '@/mixins/entities/pbehavior/reasons';
import { entitiesFieldPbehaviorFieldTypeMixin } from '@/mixins/entities/pbehavior/types-field';

import DateTimeSplittedRangePickerField from '@/components/forms/fields/date-time-splitted-range-picker-field.vue';

export default {
  inject: ['$validator'],
  components: {
    DateTimeSplittedRangePickerField,
  },
  mixins: [
    formMixin,
    formValidationHeaderMixin,
    entitiesPbehaviorReasonMixin,
    entitiesFieldPbehaviorFieldTypeMixin,
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
    nameLabel: {
      type: String,
      required: false,
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

    tstopMin() {
      const nowTimestamp = getNowTimestamp();
      const startTimestamp = convertDateToTimestamp(this.form.tstart);

      return convertDateToString(
        Math.min(nowTimestamp, startTimestamp),
        DATETIME_FORMATS.vuetifyDatePicker,
      );
    },

    tstopMax() {
      return convertDateToString(
        addUnitToDate(this.form.tstart, MAX_PBEHAVIOR_DATES_DIFF_YEARS, TIME_UNITS.year),
        DATETIME_FORMATS.vuetifyDatePicker,
      );
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
        this.updateModel({
          ...this.form,

          tstart: convertDateToStartOfDayDateObject(tstart),
          tstop: !this.noEnding && tstop ? convertDateToEndOfDayDateObject(tstop) : tstop,
        });
      }
    },
    hasPauseType(value) {
      if (!value) {
        this.noEnding = false;
      }
    },
  },
  mounted() {
    this.fetchFieldPbehaviorTypesList();
  },
  methods: {
    updateTStop(tstop) {
      this.updateField('tstop', tstop ? convertDateToEndOfUnitDateObject(tstop, TIME_UNITS.minute) : tstop);
    },

    updateStartOnTrigger(value) {
      if (value) {
        this.fullDay = false;
        this.noEnding = false;

        this.updateModel({
          ...this.form,

          start_on_trigger: true,
          tstart: null,
          tstop: null,
        });
      } else {
        const { duration, ...form } = this.form;

        this.updateModel({
          ...form,

          start_on_trigger: false,
        });
      }
    },
  },
};
</script>

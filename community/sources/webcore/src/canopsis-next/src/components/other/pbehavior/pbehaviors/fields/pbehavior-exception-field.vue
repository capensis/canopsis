<template lang="pug">
  v-layout(row, wrap, justify-space-between)
    v-flex(xs6)
      date-time-splitted-range-picker-field(
        :start="value.begin",
        :end="value.end",
        :start-label="$t('common.begin')",
        :end-label="$t('common.end')",
        :start-rules="beginRules",
        :end-rules="endRules",
        :name="datesName",
        :full-day="fullDay",
        :disabled="disabled",
        @update:start="updateField('begin', $event)",
        @update:end="updateField('end', $event)"
      )
    v-flex.pl-2(v-if="withType", :class="disabled ? 'xs6' : 'xs5'")
      c-pbehavior-type-field(
        v-field="value.type",
        :required="!disabled",
        :name="typeName",
        :disabled="disabled",
        return-object
      )
    v-flex(v-if="!disabled", xs1)
      v-btn(color="error", icon, @click="$emit('delete')")
        v-icon delete
    v-flex(xs12)
      v-checkbox.mt-0(
        v-model="fullDay",
        :label="$t('modals.createPbehavior.steps.general.fields.fullDay')",
        :disabled="disabled",
        color="primary",
        hide-details
      )
</template>

<script>
import { DATETIME_FORMATS } from '@/constants';

import {
  convertDateToEndOfDayDateObject,
  convertDateToStartOfDayDateObject,
  convertDateToString,
  isEndOfDay,
  isStartOfDay,
} from '@/helpers/date/date';

import { formMixin } from '@/mixins/form';

import DateTimeSplittedRangePickerField from '@/components/forms/fields/date-time-splitted-range-picker-field.vue';

export default {
  inject: ['$validator'],
  components: { DateTimeSplittedRangePickerField },
  mixins: [formMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    withType: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      fullDay: isStartOfDay(this.value.begin) && isEndOfDay(this.value.end),
    };
  },
  computed: {
    beginRules() {
      return {
        required: true,
        date_format: DATETIME_FORMATS.veeValidateDateTimeFormat,
      };
    },

    endRules() {
      return {
        required: true,
        after: [convertDateToString(this.value.begin, DATETIME_FORMATS.dateTimePicker)],
        date_format: DATETIME_FORMATS.veeValidateDateTimeFormat,
      };
    },

    nameSuffix() {
      return this.value.key ? `-${this.value.key}` : '';
    },

    datesName() {
      return `dates${this.nameSuffix}`;
    },

    typeName() {
      return `type${this.nameSuffix}`;
    },
  },
  watch: {
    fullDay() {
      this.updateModel({
        ...this.value,

        begin: convertDateToStartOfDayDateObject(this.value.begin),
        end: convertDateToEndOfDayDateObject(this.value.end),
      });
    },
  },
};
</script>

<template lang="pug">
  v-layout(row, wrap)
    v-flex(xs6)
      date-time-splitted-range-picker-field(
        :start="value.begin",
        :end="value.end",
        :startLabel="$t('common.begin')",
        :endLabel="$t('common.end')",
        :startRules="beginRules",
        :endRules="endRules",
        :name="datesName",
        :fullDay="fullDay",
        :disabled="disabled",
        @update:start="updateField('begin', $event)",
        @update:end="updateField('end', $event)"
      )
    v-flex.pl-2(:class="disabled ? 'xs6' : 'xs5'")
      pbehavior-type-field(
        v-field="value.type",
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
import moment from 'moment';

import { DATETIME_FORMATS } from '@/constants';
import { isEndOfDay, isStartOfDay } from '@/helpers/date/date';

import formMixin from '@/mixins/form';

import PbehaviorTypeField from '@/components/other/pbehavior/calendar/partials/pbehavior-type-field.vue';
import DateTimeSplittedRangePickerField from '@/components/forms/fields/date-time-splitted-range-picker-field.vue';

export default {
  inject: ['$validator'],
  components: { DateTimeSplittedRangePickerField, PbehaviorTypeField },
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
        after: [moment(this.value.begin).format(DATETIME_FORMATS.dateTimePicker)],
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
      const beginMoment = moment(this.value.begin).startOf('day');
      const endMoment = moment(this.value.end).endOf('day');

      this.updateModel({
        ...this.value,

        begin: beginMoment.toDate(),
        end: endMoment.toDate(),
      });
    },
  },
};
</script>

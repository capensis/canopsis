<template lang="pug">
  v-layout(row, wrap)
    v-flex(xs6)
      date-time-range-picker-field(
        v-model="dateField",
        :startLabel="$t('common.begin')",
        :endLabel="$t('common.end')",
        :startRules="beginRules",
        :endRules="endRules",
        :fullDay="fullDay"
      )
    v-flex.pl-2(xs5)
      pbehavior-type-field(
        v-field="value.type",
        :name="typeName"
      )
    v-flex(xs1)
      v-btn(color="error", icon, @click="$emit('delete')")
        v-icon delete
    v-flex(xs12)
      v-checkbox.mt-0(
        v-model="fullDay",
        :label="$t('modals.createPbehavior.steps.general.fields.fullDay')",
        color="primary",
        hide-details
      )
</template>

<script>
import moment from 'moment';

import { DATETIME_FORMATS } from '@/constants';
import { isEndOfDay, isStartOfDay } from '@/helpers/date';

import formMixin from '@/mixins/form';

import PbehaviorTypeField from '@/components/other/pbehavior/calendar/partials/pbehavior-type-field.vue';
import DateTimeRangePickerField from '@/components/forms/fields/date-time-range-picker-field.vue';

export default {
  inject: ['$validator'],
  components: { DateTimeRangePickerField, PbehaviorTypeField },
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
  },
  data() {
    return {
      fullDay: isStartOfDay(this.value.begin) && isEndOfDay(this.value.end),
    };
  },
  computed: {
    dateField: {
      set({ tstart, tstop }) {
        this.updateModel({
          ...this.value,

          begin: tstart,
          end: tstop,
        });
      },
      get() {
        return {
          tstart: this.value.begin,
          tstop: this.value.end,
        };
      },
    },

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

    beginName() {
      return `begin${this.nameSuffix}`;
    },

    endName() {
      return `end${this.nameSuffix}`;
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

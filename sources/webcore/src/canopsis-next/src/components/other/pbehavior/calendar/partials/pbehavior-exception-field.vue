<template lang="pug">
  v-layout(row)
    v-flex.mr-2(xs3)
      date-time-picker-field(
        v-field="value.begin",
        v-validate="beginRules",
        :label="$t('common.begin')",
        :name="beginName"
      )
    v-flex.mr-2(xs3)
      date-time-picker-field(
        v-field="value.end",
        v-validate="endRules",
        :label="$t('common.end')",
        :name="endName"
      )
    v-flex(xs5)
      pbehavior-type-field(
        v-field="value.type",
        :name="typeName"
      )
    v-flex(xs1)
      v-btn(color="error", icon, @click="$emit('delete')")
        v-icon delete
</template>

<script>
import moment from 'moment';

import { DATETIME_FORMATS } from '@/constants';

import DateTimePickerField from '@/components/forms/fields/date-time-picker/date-time-picker-field.vue';
import PbehaviorTypeField from '@/components/other/pbehavior/calendar/partials/pbehavior-type-field.vue';

export default {
  inject: ['$validator'],
  components: { DateTimePickerField, PbehaviorTypeField },
  props: {
    value: {
      type: Object,
      required: true,
    },
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
};
</script>

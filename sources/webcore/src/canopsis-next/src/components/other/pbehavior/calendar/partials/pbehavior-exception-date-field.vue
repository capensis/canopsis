<template lang="pug">
  v-layout(row)
    v-flex
      date-time-picker-field(
        v-field="value.begin",
        v-validate="'required'",
        :error-messages="errors.collect('begin')",
        label="Begin",
        name="begin"
      )
    v-flex
      date-time-picker-field(
        v-field="value.end",
        v-validate="'required'",
        :error-messages="errors.collect('end')",
        label="End",
        name="end"
      )
    v-flex
      v-select(
        v-field="value.type",
        v-validate="'required'",
        :items="types",
        :error-messages="errors.collect('type')",
        label="Type",
        name="type"
      )
</template>

<script>
import moment from 'moment';

import DateTimePickerField from '@/components/forms/fields/date-time-picker/date-time-picker-field.vue';

export default {
  inject: ['$validator'],
  components: { DateTimePickerField },
  props: {
    value: {
      type: Object,
      default: () => {
        const now = moment().startOf('day');

        return {
          start: now.toDate(),
          end: now.toDate(),
          type: '',
        };
      },
    },
  },
  computed: {
    types() {
      return ['a', 'b', 'c'];
    },
  },
};
</script>

<template lang="pug">
  v-layout(row)
    v-flex
      date-time-picker-field(
        v-field="value.begin",
        v-validate="'required'",
        :error-messages="errors.collect(beginName)",
        label="Begin",
        :name="beginName"
      )
    v-flex
      date-time-picker-field(
        v-field="value.end",
        v-validate="'required'",
        :error-messages="errors.collect(endName)",
        label="End",
        :name="endName"
      )
    v-flex
      v-select(
        v-field="value.type",
        v-validate="'required'",
        :items="types",
        :error-messages="errors.collect(typeName)",
        label="Type",
        :name="typeName"
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
    key: {
      type: String,
      default: undefined,
    },
  },
  computed: {
    types() {
      return ['a', 'b', 'c'];
    },
    beginName() {
      return `${this.key || this.value.key || ''}begin`;
    },
    endName() {
      return `${this.key || this.value.key || ''}end`;
    },
    typeName() {
      return `${this.key || this.value.key || ''}type`;
    },
  },
};
</script>

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
    v-flex
      v-btn(color="error", icon, @click="$emit('delete')")
        v-icon delete
</template>

<script>
import DateTimePickerField from '@/components/forms/fields/date-time-picker/date-time-picker-field.vue';

export default {
  inject: ['$validator'],
  components: { DateTimePickerField },
  props: {
    value: {
      type: Object,
      required: true,
    },
  },
  computed: {
    types() {
      return ['a', 'b', 'c'];
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

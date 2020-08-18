<template lang="pug">
  v-layout(row)
    v-flex(xs3)
      date-time-picker-field(
        v-field="value.begin",
        v-validate="'required'",
        :error-messages="errors.collect(beginName)",
        :label="$t('common.begin')",
        :name="beginName"
      )
    v-flex(xs3)
      date-time-picker-field(
        v-field="value.end",
        v-validate="'required'",
        :error-messages="errors.collect(endName)",
        :label="$t('common.end')",
        :name="endName"
      )
    v-flex(xs5)
      pbehavior-type-field(
        v-field="value.type",
        :name="typeName"
      )
    v-flex
      v-btn(color="error", icon, @click="$emit('delete')")
        v-icon delete
</template>

<script>
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

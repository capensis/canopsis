<template lang="pug">
  v-layout
    v-flex.pr-1(xs5)
      date-time-splitted-picker-field(
        v-validate="startRules",
        :value="start",
        :fullDay="fullDay",
        :disabled="disabled",
        :label="startLabel",
        :name="`${name}_start`",
        @input="$emit('update:start', $event)"
      )
    template(v-if="!noEnding")
      v-flex.pr-1(xs2)
        div.time-dash –
      v-flex(xs5)
        date-time-splitted-picker-field(
          v-validate="endRules",
          :value="end",
          :fullDay="fullDay",
          :disabled="disabled",
          :label="endLabel",
          :name="`${name}_end`",
          :min="endMin",
          :max="endMax",
          reverse,
          @input="$emit('update:end', $event)"
        )
</template>

<script>
import DateTimeSplittedPickerField from '@/components/forms/fields/date-time-picker/date-time-splitted-picker-field.vue';

export default {
  inject: ['$validator'],
  components: {
    DateTimeSplittedPickerField,
  },
  props: {
    start: {
      type: Date,
      default: null,
    },
    end: {
      type: Date,
      default: null,
    },
    endRules: {
      type: Object,
      required: true,
    },
    startRules: {
      type: Object,
      required: true,
    },
    startLabel: {
      type: String,
      required: true,
    },
    endLabel: {
      type: String,
      required: true,
    },
    name: {
      type: String,
      default: 'date-range',
    },
    noEnding: {
      type: Boolean,
      default: false,
    },
    fullDay: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    endMin: {
      type: String,
      required: false,
    },
    endMax: {
      type: String,
      required: false,
    },
  },
};
</script>

<style lang="scss" scoped>
.time-dash {
  line-height: 68px;
  padding: 0 8px;
  text-align: center;
}
</style>

<template>
  <v-layout class="date-time-splitted-range-picker-field">
    <v-flex class="date-time-splitted-range-picker-field__start">
      <date-time-splitted-picker-field
        v-validate="startRules"
        :value="start"
        :full-day="fullDay"
        :disabled="disabled"
        :label="startLabel"
        :name="`${name}_start`"
        @input="$emit('update:start', $event)"
      />
    </v-flex>
    <template v-if="!noEnding">
      <div class="date-time-splitted-range-picker-field__time-dash">
        –
      </div>
      <v-flex class="date-time-splitted-range-picker-field__end">
        <date-time-splitted-picker-field
          v-validate="endRules"
          :value="end"
          :full-day="fullDay"
          :disabled="disabled"
          :label="endLabel"
          :name="`${name}_end`"
          :min="endMin"
          :max="endMax"
          reverse
          @input="$emit('update:end', $event)"
        />
      </v-flex>
    </template>
  </v-layout>
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

<style lang="scss">
.date-time-splitted-range-picker-field {
  --dash-width: 26px;

  &__time-dash {
    width: var(--dash-width);
    line-height: 68px;
    text-align: center;
  }

  &__start, &__end {
    width: calc(50% - calc(var(--dash-width) / 2));
    max-width: calc(50% - calc(var(--dash-width) / 2));
  }
}
</style>

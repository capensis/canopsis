<template>
  <c-layout class="date-time-splitted-range-picker-field">
    <div class="date-time-splitted-range-picker-field__start">
      <date-time-splitted-picker-field
        :value="start"
        :full-day="fullDay"
        :disabled="disabled"
        :control-variant="controlVariant"
        :label="startLabel"
        :name="`${name}_start`"
        :rules="startRules"
        :hide-details="hideDetails"
        @update:value="$emit('update:start', $event)"
      />
    </div>
    <template v-if="!noEnding">
      <div class="date-time-splitted-range-picker-field__time-dash">
        -
      </div>
      <div class="date-time-splitted-range-picker-field__end">
        <date-time-splitted-picker-field
          :value="end"
          :full-day="fullDay"
          :disabled="disabled"
          :control-variant="controlVariant"
          :label="endLabel"
          :name="`${name}_end`"
          :min="endMin"
          :max="endMax"
          :rules="endRules"
          :hide-details="hideDetails"
          reverse
          @update:value="$emit('update:end', $event)"
        />
      </div>
    </template>
  </c-layout>
</template>

<script>
import DateTimeSplittedPickerField from '@/components/forms/fields/date-time-picker/date-time-splitted-picker-field.vue';

export default {
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
    controlVariant: {
      type: String,
      default: undefined,
    },
    endMin: {
      type: String,
      required: false,
    },
    endMax: {
      type: String,
      required: false,
    },
    hideDetails: {
      type: Boolean,
      default: false,
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

  &__start,
  &__end {
    width: calc(50% - calc(var(--dash-width) / 2));
    max-width: calc(50% - calc(var(--dash-width) / 2));
  }
}
</style>

<template>
  <v-input :error-messages="errorMessage ? [errorMessage] : []" :hide-details="hideDetails">
    <v-row :style="rowStyle" class="date-time-splitted-field">
      <v-col class="date-time-splitted-field__date">
        <c-date-picker-field
          :value="dateValue"
          :label="!reverse || fullDay ? label : ''"
          :placeholder="placeholder"
          :error="!!errorMessage"
          :disabled="disabled"
          :control-variant="controlVariant"
          :min="min"
          :max="max"
          header-color="primary"
          hide-details
          @update:value="updateDate"
        />
      </v-col>
      <v-col v-if="!fullDay" class="date-time-splitted-field__time" cols="auto">
        <time-picker-field
          :value="timeValue"
          :label="reverse ? label : ''"
          :error="!!errorMessage"
          :disabled="disabled"
          hide-details
          @update:value="updateTime"
        />
      </v-col>
    </v-row>
  </v-input>
</template>

<script>
import { computed, toRef, watch } from 'vue';

import { DATETIME_FORMATS, TIME_UNITS } from '@/constants';

import { convertDateToStartOfUnitDateObject, convertDateToString } from '@/helpers/date/date';
import { getDateObjectByTime, getDateObjectByDate } from '@/helpers/date/date-time-picker';

import { useField } from '@/hooks/form';

import TimePickerField from '@/components/forms/fields/time-picker/time-picker-field.vue';

import { createDateTimeValidationRules } from './validation-rules';

export default {
  emits: ['update:value'],
  components: { TimePickerField },
  props: {
    value: {
      type: Date,
      default: null,
    },
    fullDay: {
      type: Boolean,
      default: false,
    },
    label: {
      type: String,
      default: '',
    },
    placeholder: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'date',
    },
    reverse: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    min: {
      type: String,
      required: false,
    },
    max: {
      type: String,
      required: false,
    },
    rules: {
      type: Object,
      default: () => ({}),
    },
    controlVariant: {
      type: String,
      default: undefined,
    },
    hideDetails: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const fieldName = toRef(props, 'name');
    const getValidatedValue = (value) => {
      if (!value) {
        return value;
      }

      return convertDateToString(
        convertDateToStartOfUnitDateObject(value, TIME_UNITS.minute),
        DATETIME_FORMATS.dateTimePicker,
      );
    };

    const validateValue = (value) => {
      if (props.rules.required && !value) {
        return 'Field is required';
      }

      const validatedValue = getValidatedValue(value);
      const validationRules = createDateTimeValidationRules(props.rules);
      const failedRule = validationRules.find((validationRule) => {
        const result = validationRule(validatedValue);

        return result !== true;
      });

      return failedRule ? failedRule(validatedValue) : true;
    };

    const { errorMessage, setValue } = useField(fieldName, validateValue, {
      initialValue: props.value,
      validateOnMount: false,
      validateOnValueUpdate: true,
    });

    const dateValue = computed(() => convertDateToString(props.value, DATETIME_FORMATS.vuetifyDatePicker, null),
    );
    const timeValue = computed(() => convertDateToString(props.value, DATETIME_FORMATS.timePicker, null),
    );
    const rowStyle = computed(() => ({
      flexDirection: props.reverse ? 'row-reverse' : undefined,
    }));

    const updateModel = (value) => {
      setValue(value);
      emit('update:value', value);
    };

    const updateTime = time => updateModel(getDateObjectByTime(props.value, time));
    const updateDate = date => updateModel(getDateObjectByDate(props.value, date));

    watch(
      () => props.value,
      value => setValue(value, false),
    );

    return {
      dateValue,
      timeValue,
      errorMessage,
      rowStyle,
      updateTime,
      updateDate,
    };
  },
};
</script>

<style lang="scss">
.date-time-splitted-field {
  width: 100%;

  &__date {
    flex: 1 1 160px;
    min-width: 0;
  }

  &__time {
    align-self: end;
    flex: 0 0 56px;
    width: 56px;
    max-width: 56px;

    .v-field {
      font-size: 16px !important;
      padding-inline: 0;

      .v-field__input {
        padding-inline: 0;
      }
    }
  }
}
</style>

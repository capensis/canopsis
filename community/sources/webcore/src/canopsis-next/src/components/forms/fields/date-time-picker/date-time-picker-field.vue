<template>
  <v-menu
    v-model="opened"
    ref="menu"
    :close-on-content-click="false"
    :disabled="disabled"
    content-class="date-time-picker"
    transition="slide-y-transition"
    max-width="290px"
    right
  >
    <template #activator="{ on }">
      <div :aria-required="isRequired" v-on="on">
        <v-text-field
          :label="label"
          :error-messages="errors.collect(name)"
          :value="dateTextValue"
          :append-icon="clearable ? 'close' : ''"
          :disabled="disabled"
          :hide-details="hideDetails"
          readonly
          @click:append="clear"
        />
      </div>
    </template>
    <date-time-picker
      :value="value"
      :label="label"
      :round-hours="roundHours"
      :allowed-dates="allowedDates"
      @close="close"
      @input="input"
    />
  </v-menu>
</template>

<script>
import { ref, computed, nextTick } from 'vue';

import { DATETIME_FORMATS, TIME_UNITS } from '@/constants';

import { convertDateToStartOfUnitDateObject, convertDateToString } from '@/helpers/date/date';

import { useValidator } from '@/hooks/validator/validator';

import DateTimePicker from './date-time-picker.vue';

/**
 * Date time picker component
 *
 * @warning If you want to use validation on the field you shouldn't use `v-field`
 *
 * @prop {Date} [value=null] - Date value
 * @prop {Boolean} [clearable=false] - if it is true input field will be have cross button with clear event on click
 * @prop {string} [label=''] - Label of the input field
 * @prop {string} [name=null] - Name property in the validation object
 * @prop {Boolean} [roundHours=false] - Deny to change minutes it will be only 0
 *
 * @event value#input
 */
export default {
  $_veeValidate: {
    value() {
      if (!this.value) {
        return this.value;
      }

      return convertDateToStartOfUnitDateObject(this.value, TIME_UNITS.minute);
    },

    name() {
      return this.name;
    },
  },
  inject: ['$validator'],
  components: { DateTimePicker },
  props: {
    clearable: {
      type: Boolean,
      default: false,
    },
    value: {
      type: [Date, Number],
      default: () => new Date(),
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: null,
    },
    roundHours: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    hideDetails: {
      type: Boolean,
      default: false,
    },
    allowedDates: {
      type: Function,
      required: false,
    },
  },
  setup(props, { emit }) {
    const validator = useValidator();

    const opened = ref(false);

    const isRequired = computed(() => validator?.fields?.find?.({ name: props.name })?.isRequired);
    const dateTextValue = computed(() => convertDateToString(props.value, DATETIME_FORMATS.dateTimePicker));

    /**
     * Closes the date time picker menu
     */
    const close = () => opened.value = false;

    /**
     * Clears the date time picker value
     */
    const clear = () => emit('input', null);

    /**
     * Handles input event from date time picker
     *
     * @param {Date|Number} value - Date value
     */
    const input = (value) => {
      emit('input', value);

      if (validator && props.name) {
        nextTick(() => validator.validate(props.name));
      }
    };

    return {
      opened,
      isRequired,
      dateTextValue,
      close,
      clear,
      input,
    };
  },
};
</script>

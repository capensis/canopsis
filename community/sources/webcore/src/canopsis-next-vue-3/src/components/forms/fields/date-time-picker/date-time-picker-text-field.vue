<template>
  <c-layout align="center">
    <c-text-field
      :model-value="value"
      :label="label"
      :name="name"
      :autofocus="autofocus"
      :required="fieldRequired"
      :validation-rules="validationRules"
      @update:model-value="updateModel"
      @focus="focus"
      @blur="blur"
    />
    <date-time-picker-menu
      :value="objectValue"
      :label="label"
      :round-hours="roundHours"
      @update:value="updateObjectField"
    />
  </c-layout>
</template>

<script>
import { isEmpty } from 'lodash';

import { DATETIME_FORMATS } from '@/constants';
import { convertDateToDateObject, convertDateToString } from '@/helpers/date/date';
import { formBaseMixin } from '@/mixins/form';

import DateTimePickerMenu from './date-time-picker-menu.vue';
import { createDateTimeValidationRules, isRequiredRule } from './validation-rules';

export default {
  emits: ['update:value', 'update:objectValue'],
  components: { DateTimePickerMenu },
  mixins: [formBaseMixin],
  props: {
    value: {
      type: String,
      default: '',
    },
    label: {
      type: String,
      default: null,
    },
    name: {
      type: String,
      default: '',
    },
    roundHours: {
      type: Boolean,
      default: false,
    },
    dateObjectPreparer: {
      type: Function,
      default: value => convertDateToDateObject(value, DATETIME_FORMATS.dateTimePicker),
    },
    autofocus: {
      type: Boolean,
      default: false,
    },
    rules: {
      type: Object,
      default: () => ({}),
    },
    required: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      objectValue: new Date(),
      isFocusedTextField: false,
    };
  },
  computed: {
    fieldRequired() {
      return this.required || isRequiredRule(this.rules);
    },

    validationRules() {
      return createDateTimeValidationRules(this.rules, this.dateObjectPreparer);
    },
  },
  watch: {
    value: {
      immediate: true,
      handler(value) {
        if (!this.isFocusedTextField) {
          this.setObjectValue(value);
        }
      },
    },

    isFocusedTextField(focused) {
      if (!focused) {
        this.setObjectValue(this.value);
      }
    },
  },
  methods: {
    updateObjectField(value) {
      this.updateModel(convertDateToString(value, DATETIME_FORMATS.dateTimePicker));
    },

    focus() {
      this.isFocusedTextField = true;
    },

    blur() {
      this.isFocusedTextField = false;
    },

    setObjectValue(value) {
      try {
        if (!isEmpty(value)) {
          this.objectValue = this.dateObjectPreparer(value);
        }
      } catch {
        this.objectValue = null;
      }

      this.$emit('update:objectValue', this.objectValue);
    },
  },
};
</script>

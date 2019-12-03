<template lang="pug">
  div
    v-layout
      v-flex
        v-text-field(
          data-test="dateTimePickerTextField",
          :value="value",
          :label="label",
          :name="name",
          :error-messages="errorMessages",
          @focus="focus",
          @blur="blur",
          @input="updateModel($event)"
        )
      v-flex
        date-time-picker-button(
          :value="objectValue",
          :roundHours="roundHours",
          :useSeconds="useSeconds",
          @input="updateObjectField"
        )
</template>

<script>
import moment from 'moment';
import { isEmpty } from 'lodash';

import { DATETIME_FORMATS } from '@/constants';

import uid from '@/helpers/uid';

import formBaseMixin from '@/mixins/form/base';

import DateTimePickerButton from './date-time-picker-button.vue';

export default {
  $_veeValidate: {
    value() {
      return this.value;
    },

    name() {
      return this.name;
    },
  },
  inject: ['$validator'],
  components: { DateTimePickerButton },
  mixins: [formBaseMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
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
      default: null,
    },
    roundHours: {
      type: Boolean,
      default: false,
    },
    useSeconds: {
      type: Boolean,
      default: false,
    },
    dateObjectPreparer: {
      type: Function,
      default: value => moment(value, DATETIME_FORMATS.dateTimePicker).toDate(),
    },
  },
  data() {
    return {
      objectValue: null,
      isFocusedTextField: false,
    };
  },
  computed: {
    errorMessages() {
      if (this.$validator && this.errors && this.name) {
        return this.errors.collect(this.name);
      }

      return [];
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

    isFocusedTextField(value) {
      if (!value) {
        this.setObjectValue(this.value);
      }
    },
  },
  created() {
    if (this.$validator) {
      const pickerFormatRuleName = `picker_format:${uid()}`;

      this.$validator.extend(pickerFormatRuleName, {
        getMessage: () => this.$t('modals.statsDateInterval.errors.endDateLessOrEqualStartDate'),
        validate: (value) => {
          try {
            if (!isEmpty(value)) {
              return Boolean(this.dateObjectPreparer(value));
            }

            return true;
          } catch (err) {
            return false;
          }
        },
      });

      this.$validator.attach({
        name: this.name,
        rules: pickerFormatRuleName,
        getter: () => this.value,
        context: () => this,
      });
    }
  },
  beforeDestroy() {
    if (this.$validator) {
      this.$validator.detach(this.name);
    }
  },
  methods: {
    updateObjectField(value) {
      this.updateModel(moment(value).format(DATETIME_FORMATS.dateTimePicker));
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
      } catch (err) {
        this.objectValue = null;
      }

      this.$emit('update:objectValue', this.objectValue);
    },
  },
};
</script>

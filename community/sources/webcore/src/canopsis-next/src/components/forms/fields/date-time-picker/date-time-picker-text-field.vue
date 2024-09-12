<template lang="pug">
  div
    v-layout(align-center)
      v-flex
        v-text-field(
          v-field="value",
          :label="label",
          :name="name",
          :error-messages="errorMessages",
          @focus="focus",
          @blur="blur"
        )
      v-flex
        date-time-picker-menu(
          :value="objectValue",
          :label="label",
          :round-hours="roundHours",
          @input="updateObjectField"
        )
</template>

<script>
import { isEmpty } from 'lodash';

import { DATETIME_FORMATS } from '@/constants';

import { convertDateToDateObject, convertDateToString } from '@/helpers/date/date';

import { formBaseMixin } from '@/mixins/form';

import DateTimePickerMenu from './date-time-picker-menu.vue';

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
  components: { DateTimePickerMenu },
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
    dateObjectPreparer: {
      type: Function,
      default: value => convertDateToDateObject(value, DATETIME_FORMATS.dateTimePicker),
    },
  },
  data() {
    return {
      objectValue: new Date(),
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

    ruleName() {
      return `picker_format_${this.name}`;
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
  created() {
    if (this.$validator && this.name) {
      this.$validator.attach({
        name: this.name,
        rules: {
          picker_format: {
            preparer: this.dateObjectPreparer,
          },
        },
        getter: () => this.value,
        vm: this,
      });
    }
  },
  beforeDestroy() {
    if (this.$validator && this.name) {
      this.$validator.detach(this.name);
    }
  },
  methods: {
    updateObjectField(value) {
      this.updateModel(
        convertDateToString(value, DATETIME_FORMATS.dateTimePicker),
      );
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

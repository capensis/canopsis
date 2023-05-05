<template lang="pug">
  v-layout(column)
    v-layout(row, justify-space-between)
      v-flex(xs6)
        date-time-splitted-range-picker-field(
          v-if="editing",
          :start="value.begin",
          :end="value.end",
          :start-label="$t('common.begin')",
          :end-label="$t('common.end')",
          :start-rules="beginRules",
          :end-rules="endRules",
          :name="datesName",
          :full-day="fullDay",
          :disabled="disabled",
          @update:start="updateField('begin', $event)",
          @update:end="updateField('end', $event)"
        )
        date-time-splitted-range-picker-text(
          v-else,
          :start="value.begin",
          :end="value.end",
          :start-label="$t('common.begin')",
          :end-label="$t('common.end')",
          :full-day="fullDay"
        )
      v-flex.pl-2(v-if="withType")
        c-pbehavior-type-field(
          v-if="editing",
          v-field="value.type",
          :required="!disabled",
          :name="typeName",
          :disabled="disabled",
          return-object
        )
        c-pbehavior-type-text(
          v-else,
          :value="value.type"
        )
      v-flex.exception-field-item--actions(v-if="!disabled")
        v-btn.btn--editing(
          :input-value="editing",
          icon,
          fab,
          @click="toggleEditing"
        )
          v-icon edit
          v-icon(color="primary") check
        v-btn.mx-0(
          color="error",
          icon,
          @click="$emit('delete')"
        )
          v-icon delete
    v-layout(row)
      v-checkbox.mt-0(
        v-model="fullDay",
        :label="$t('modals.createPbehavior.steps.general.fields.fullDay')",
        :disabled="disabled || !editing",
        color="primary",
        hide-details
      )
</template>

<script>
import { DATETIME_FORMATS } from '@/constants';

import {
  convertDateToEndOfDayDateObject,
  convertDateToStartOfDayDateObject,
  convertDateToString,
  isEndOfDay,
  isStartOfDay,
} from '@/helpers/date/date';

import { formMixin, validationChildrenMixin } from '@/mixins/form';

import DateTimeSplittedRangePickerField from '@/components/forms/fields/date-time-splitted-range-picker-field.vue';
import DateTimeSplittedRangePickerText from '@/components/forms/fields/date-time-picker/date-time-splitted-range-picker-text.vue';

export default {
  inject: ['$validator'],
  components: {
    DateTimeSplittedRangePickerField,
    DateTimeSplittedRangePickerText,
  },
  mixins: [
    formMixin,
    validationChildrenMixin,
  ],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    withType: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      editing: !this.value.type,
      fullDay: isStartOfDay(this.value.begin) && isEndOfDay(this.value.end),
    };
  },
  computed: {
    beginRules() {
      return {
        required: true,
        date_format: DATETIME_FORMATS.veeValidateDateTimeFormat,
      };
    },

    endRules() {
      return {
        required: true,
        after: [convertDateToString(this.value.begin, DATETIME_FORMATS.dateTimePicker)],
        date_format: DATETIME_FORMATS.veeValidateDateTimeFormat,
      };
    },

    nameSuffix() {
      return this.value.key ? `-${this.value.key}` : '';
    },

    datesName() {
      return `dates${this.nameSuffix}`;
    },

    typeName() {
      return `type${this.nameSuffix}`;
    },
  },
  watch: {
    fullDay() {
      this.updateModel({
        ...this.value,

        begin: convertDateToStartOfDayDateObject(this.value.begin),
        end: convertDateToEndOfDayDateObject(this.value.end),
      });
    },
  },
  methods: {
    async toggleEditing() {
      if (this.editing) {
        await this.validateChildren();
      }

      if (!this.hasChildrenError) {
        this.editing = !this.editing;
      }
    },
  },
};
</script>

<style lang="scss" scoped>
.exception-field-item--actions {
  min-width: 90px;

  .btn--editing {
    height: 36px;
    width: 36px;
  }
}
</style>

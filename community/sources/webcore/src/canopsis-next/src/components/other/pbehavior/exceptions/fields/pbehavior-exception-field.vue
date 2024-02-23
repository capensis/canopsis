<template>
  <v-layout
    class="pbehavior-exception-field"
    column
  >
    <v-layout justify-space-between>
      <v-flex
        class="pbehavior-exception-field__interval"
        xs6
      >
        <date-time-splitted-range-picker-field
          v-if="editing"
          :start="value.begin"
          :end="value.end"
          :start-label="$t('common.begin')"
          :end-label="$t('common.end')"
          :start-rules="beginRules"
          :end-rules="endRules"
          :name="datesName"
          :full-day="fullDay"
          :disabled="disabled"
          @update:start="updateField('begin', $event)"
          @update:end="updateField('end', $event)"
        />
        <date-time-splitted-range-picker-text
          v-else
          :start="value.begin"
          :end="value.end"
          :start-label="$t('common.begin')"
          :end-label="$t('common.end')"
          :full-day="fullDay"
        />
      </v-flex>
      <v-flex
        v-if="withType"
        class="pl-2"
      >
        <c-pbehavior-type-field
          v-if="editing"
          v-field="value.type"
          :required="!disabled"
          :name="typeName"
          :disabled="disabled"
          return-object
        />
        <c-pbehavior-type-text
          v-else
          :value="value.type"
        />
      </v-flex>
      <v-flex
        v-if="!disabled"
        class="pbehavior-exception-field__actions"
      >
        <v-btn
          :input-value="editing"
          class="btn--editing"
          icon
          fab
          @click="toggleEditing"
        >
          <v-icon
            v-if="editing"
            color="primary"
          >
            check
          </v-icon>
          <v-icon v-else>
            edit
          </v-icon>
        </v-btn>
        <v-btn
          class="v-btn-legacy-m--y"
          color="error"
          icon
          @click="$emit('delete')"
        >
          <v-icon>delete</v-icon>
        </v-btn>
      </v-flex>
    </v-layout>
    <v-layout>
      <v-checkbox
        v-model="fullDay"
        :label="$t('modals.createPbehavior.steps.general.fields.fullDay')"
        :disabled="disabled || !editing"
        class="mt-0"
        color="primary"
        hide-details
      />
    </v-layout>
  </v-layout>
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
      fullDay: isStartOfDay(this.value.begin) && (isEndOfDay(this.value.end) || isStartOfDay(this.value.end)),
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
.pbehavior-exception-field {
  &__interval {
    flex-shrink: 0;
  }

  &__actions {
    min-width: 90px;

    .btn--editing {
      height: 36px;
      width: 36px;
    }
  }
}
</style>

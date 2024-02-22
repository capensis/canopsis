<template>
  <v-radio-group
    :value="type"
    :label="$t('common.end')"
    @change="updateType"
  >
    <v-radio
      :value="types.never"
      color="primary"
    >
      <template #label="">
        {{ $t('recurrenceRule.never') }}
      </template>
    </v-radio>
    <v-radio
      :value="types.date"
      class="mb-0"
      color="primary"
    >
      <template #label="">
        <span>{{ $t('recurrenceRule.on') }}</span>
        <date-time-splitted-picker-field
          v-field="value.until"
          v-validate="dateTimeSplittedRules"
          :placeholder="$t('common.date')"
          :disabled="!isDateType"
          class="ml-3"
          name="until"
        />
      </template>
    </v-radio>
    <v-radio
      :value="types.after"
      class="mb-0"
      color="primary"
    >
      <template #label="">
        <span>{{ $t('recurrenceRule.after') }}</span>
        <c-number-field
          v-field="value.count"
          :disabled="!isAfterType"
          :required="isAfterType"
          class="mx-3"
          name="count"
        />
        <span class="text-lowercase">{{ $tc('recurrenceRule.occurrence', value.count || 1) }}</span>
      </template>
    </v-radio>
  </v-radio-group>
</template>

<script>
import { omit, isNumber } from 'lodash';

import { RECURRENCE_RULE_END_TYPES } from '@/constants';

import { formBaseMixin } from '@/mixins/form';

import DateTimeSplittedPickerField from '@/components/forms/fields/date-time-picker/date-time-splitted-picker-field.vue';

export default {
  inject: ['$validator'],
  components: { DateTimeSplittedPickerField },
  mixins: [formBaseMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      required: true,
    },
    types: {
      type: Object,
      default: () => RECURRENCE_RULE_END_TYPES,
    },
  },
  data() {
    let type = RECURRENCE_RULE_END_TYPES.never;

    if (this.value.until) {
      type = RECURRENCE_RULE_END_TYPES.date;
    }

    if (isNumber(this.value.count)) {
      type = RECURRENCE_RULE_END_TYPES.after;
    }

    return {
      type,
    };
  },
  computed: {
    dateTimeSplittedRules() {
      return {
        required: this.isDateType,
      };
    },

    isDateType() {
      return this.type === this.types.date;
    },

    isAfterType() {
      return this.type === this.types.after;
    },
  },
  methods: {
    updateType(type) {
      this.type = type;

      const deletingFields = [];

      if (!this.isDateType) {
        deletingFields.push('until');
      }

      if (!this.isAfterType) {
        deletingFields.push('count');
      }

      deletingFields.forEach(name => this.errors.remove(name));

      this.updateModel(omit(this.value, deletingFields));
    },
  },
};
</script>

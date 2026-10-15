<template>
  <v-row class="c-date-time-interval-field" density="compact">
    <v-col :cols="columnCols">
      <date-time-picker-field
        :value="value.from"
        :label="$t('common.from')"
        :disabled="disabled"
        :name="fromFieldName"
        :hide-details="hideDetails"
        :allowed-dates="isAllowedFromDate"
        :round-hours="roundHours"
        :rules="rules"
        :required="required"
        @update:value="updateField('from', $event)"
      >
        <template v-if="withIcon" #append-inner>
          <v-icon size="18">
            calendar_today
          </v-icon>
        </template>
      </date-time-picker-field>
    </v-col>
    <v-col :cols="columnCols">
      <date-time-picker-field
        :value="value.to"
        :label="$t('common.to')"
        :disabled="disabled"
        :name="toFieldName"
        :hide-details="hideDetails"
        :allowed-dates="isAllowedToDate"
        :round-hours="roundHours"
        :rules="rules"
        :required="required"
        @update:value="updateField('to', $event)"
      >
        <template v-if="withIcon" #append-inner>
          <v-icon size="18">
            calendar_today
          </v-icon>
        </template>
      </date-time-picker-field>
    </v-col>
  </v-row>
</template>

<script>
import { computed } from 'vue';

import { useModelField } from '@/hooks/form/model-field';

import DateTimePickerField from './date-time-picker-field.vue';

export default {
  components: {
    DateTimePickerField,
  },
  emits: ['update:value'],
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      default: 'interval',
    },
    disabled: {
      type: Boolean,
      required: false,
    },
    hideDetails: {
      type: Boolean,
      required: false,
    },
    required: {
      type: Boolean,
      required: false,
    },
    isAllowedFromDate: {
      type: Function,
      required: false,
    },
    isAllowedToDate: {
      type: Function,
      required: false,
    },
    column: {
      type: Boolean,
      default: false,
    },
    roundHours: {
      type: Boolean,
      default: false,
    },
    withIcon: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const { updateField } = useModelField(props, emit);

    const rules = computed(() => ({
      required: props.required,
    }));
    const columnCols = computed(() => (props.column ? 12 : 6));
    const fromFieldName = computed(() => `${props.name}.from`);
    const toFieldName = computed(() => `${props.name}.to`);

    return {
      rules,
      columnCols,
      fromFieldName,
      toFieldName,
      updateField,
    };
  },
};
</script>

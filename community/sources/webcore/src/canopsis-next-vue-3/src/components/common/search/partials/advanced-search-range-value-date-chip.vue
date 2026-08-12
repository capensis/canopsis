<template>
  <v-chip
    :closable="closable"
    :color="color"
    class="c-advanced-search__chip c-advanced-search__array-chip"
    size="small"
    @click.prevent="$emit('click')"
    @click:close="close"
    @focusout="$emit('focusout')"
  >
    <date-time-picker-menu
      v-for="chip in chips"
      :key="chip.key"
      :value="chipValue[chip.key]"
      :allowed-dates="chip.allowedDates"
      @update:value="updateField(chip.key, $event)"
    >
      <template #activator="{ props: activatorProps, value }">
        <v-chip v-bind="activatorProps" size="x-small">
          <span class="font-italic text-medium-emphasis">
            {{ $t(`common.${chip.key}`) }}:
          </span>
          {{ formatDate(value) }}
        </v-chip>
      </template>
    </date-time-picker-menu>
  </v-chip>
</template>

<script>
import { computed, toRef } from 'vue';

import DateTimePickerMenu from '@/components/forms/fields/date-time-picker/date-time-picker-menu.vue';
import { convertDateToString } from '@/helpers/date/date';
import { useDateRangeAllowedDates } from '@/hooks/form/date-range-allowed-dates';
import { useModelField } from '@/hooks/form/model-field';

export default {
  components: {
    DateTimePickerMenu,
  },

  inheritAttrs: false,

  props: {
    modelValue: {
      type: Object,
      default: () => ({ from: 0, to: 0 }),
    },
    closable: {
      type: Boolean,
      default: false,
    },
    color: {
      type: String,
      required: false,
    },
  },

  emits: ['update:modelValue', 'click', 'focusout', 'close'],

  setup(props, { emit }) {
    const { updateField } = useModelField(props, emit);
    const { isAllowedFromDate, isAllowedToDate } = useDateRangeAllowedDates(
      toRef(props, 'modelValue'),
    );

    const chipValue = computed(() => props.modelValue ?? { from: 0, to: 0 });

    const chips = [
      { key: 'from', allowedDates: isAllowedFromDate },
      { key: 'to', allowedDates: isAllowedToDate },
    ];

    const formatDate = value => convertDateToString(value, 'long', '--.--.---- --:--');
    const close = () => emit('close');

    return {
      chipValue,
      chips,

      updateField,
      formatDate,
      close,
    };
  },
};
</script>

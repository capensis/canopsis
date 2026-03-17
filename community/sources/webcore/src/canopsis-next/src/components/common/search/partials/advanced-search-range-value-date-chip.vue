<template>
  <v-chip
    :close="closable"
    :color="color"
    class="c-advanced-search__chip c-advanced-search__array-chip"
    @click.prevent=""
    @click:close="close"
  >
    <date-time-picker-menu
      v-for="chip in chips"
      v-field="value[chip.key]"
      :key="chip.key"
      :allowed-dates="chip.allowedDates"
    >
      <template #activator="{ on, value: chipValue }">
        <v-chip v-on="on">
          <span class="font-italic grey--text">
            {{ $t(`common.${chip.key}`) }}:
          </span> {{ chipValue | date('long', '--.--.---- --:--') }}
        </v-chip>
      </template>
    </date-time-picker-menu>
  </v-chip>
</template>

<script>
import { toRef } from 'vue';

import { useDateRangeAllowedDates } from '@/hooks/form/date-range-allowed-dates';

import DateTimePickerMenu from '@/components/forms/fields/date-time-picker/date-time-picker-menu.vue';

export default {
  components: {
    DateTimePickerMenu,
  },
  props: {
    value: {
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
  setup(props, { emit }) {
    const { isAllowedFromDate, isAllowedToDate } = useDateRangeAllowedDates(toRef(props, 'value'));

    const chips = [
      { key: 'from', allowedDates: isAllowedFromDate },
      { key: 'to', allowedDates: isAllowedToDate },
    ];

    const close = () => emit('close');

    return {
      chips,
      close,
    };
  },
};
</script>

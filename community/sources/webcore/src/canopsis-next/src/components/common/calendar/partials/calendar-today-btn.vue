<template>
  <c-tooltip location="bottom">
    <template #activator="{ props: tooltipProps }">
      <v-btn
        v-bind="tooltipProps"
        class="mx-2 my-1"
        variant="text"
        @click="$emit('click')"
      >
        <span class="calendar-today-btn__label">
          {{ t('calendar.today') }}
        </span>
        <v-icon class="calendar-today-btn__icon">
          today
        </v-icon>
      </v-btn>
    </template>

    <span>{{ todayDate }}</span>
  </c-tooltip>
</template>

<script>
import { computed } from 'vue';

import { convertDateToString, getNowTimestamp } from '@/helpers/date/date';

import { useI18n } from '@/hooks/i18n';

export default {
  emits: ['click'],

  setup() {
    const { t } = useI18n();
    const todayDate = computed(() => convertDateToString(getNowTimestamp(), 'dddd, MMMM D'));

    return {
      t,
      todayDate,
    };
  },
};
</script>

<style lang="scss" scoped>
.calendar-today-btn__icon {
  display: none;
}

@media (max-width: 600px) {
  .calendar-today-btn__label {
    display: none;
  }

  .calendar-today-btn__icon {
    display: inline-flex;
  }
}
</style>

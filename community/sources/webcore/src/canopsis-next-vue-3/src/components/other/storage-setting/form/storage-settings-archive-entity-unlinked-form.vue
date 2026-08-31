<template>
  <c-layout gap="2" column>
    <span class="text-medium-emphasis">
      {{ $t('storageSetting.entityUnlinked.archiveUnlinkedAfter') }}
    </span>
    <c-duration-field
      :duration="modelValue"
      :units-label="$t('common.unit')"
      :units="timeUnits"
      required
      @update:duration="updateModel"
    />
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { AVAILABLE_TIME_UNITS, TIME_UNITS } from '@/constants';

import { useModelField } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    modelValue: {
      type: Object,
      default: () => ({
        value: 60,
        unit: TIME_UNITS.day,
      }),
    },
    name: {
      type: String,
      default: 'duration',
    },
  },
  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const { updateModel } = useModelField(props, emit);
    const timeUnits = computed(() => [
        AVAILABLE_TIME_UNITS.day,
        AVAILABLE_TIME_UNITS.week,
        AVAILABLE_TIME_UNITS.month,
        AVAILABLE_TIME_UNITS.year,
      ].map(({ value, text }) => ({
        value,
        text: t(text, props.modelValue.value || 0),
      })),
    );

    return {
      timeUnits,
      updateModel,
    };
  },
};
</script>

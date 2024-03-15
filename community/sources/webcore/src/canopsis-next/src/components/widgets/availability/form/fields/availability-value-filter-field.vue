<template>
  <v-layout class="availability-value-filter-field gap-3" d-inline-flex align-end>
    <v-select
      v-model="method"
      :label="$t('availability.filterByValue')"
      :items="valueFilterMethods"
      class="availability-value-filter-field__method"
      hide-details
    />
    <template v-if="method">
      <c-percents-field
        v-if="isPercentType"
        v-field="value.value"
        class="availability-value-filter-field__value"
        hide-details
      />
      <c-splitted-duration-field
        v-else
        v-field="value.value"
        :max-value="maxSeconds"
        class="availability-value-filter-field__value"
        hide-details
      />
      <c-action-btn
        :tooltip="$t('common.clear')"
        type="delete"
        icon="clear"
        color="white"
        small
        @click="clearValue"
      />
    </template>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { AVAILABILITY_SHOW_TYPE, AVAILABILITY_VALUE_FILTER_METHODS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModelValue } from '@/hooks/form';

export default {
  props: {
    value: {
      type: Object,
      required: false,
    },
    showType: {
      type: Number,
      default: AVAILABILITY_SHOW_TYPE.percent,
    },
    maxSeconds: {
      type: Number,
      required: false,
    },
  },
  setup(props, { emit }) {
    const { t } = useI18n();
    const modelValue = useModelValue(props, emit);

    const method = computed({
      get() {
        return modelValue.value?.method;
      },
      set(newMethod) {
        modelValue.value = {
          ...modelValue.value,
          value: modelValue.value?.value ?? 0,
          method: newMethod,
        };
      },
    });
    const isPercentType = computed(() => props.showType === AVAILABILITY_SHOW_TYPE.percent);
    const valueFilterMethods = computed(
      () => Object.values(AVAILABILITY_VALUE_FILTER_METHODS).map(valueFilterMethod => ({
        value: valueFilterMethod,
        text: t(`availability.valueFilterMethods.${valueFilterMethod}`),
      })),
    );

    const clearValue = () => {
      modelValue.value = undefined;
    };

    return {
      method,

      isPercentType,
      valueFilterMethods,
      clearValue,
    };
  },
};
</script>

<style lang="scss">
.availability-value-filter-field {
  flex-grow: 0;

  &__method {
    width: 140px;
    flex-grow: 0;
    flex-shrink: 0;
  }

  &__value {
    width: 160px;
  }
}
</style>

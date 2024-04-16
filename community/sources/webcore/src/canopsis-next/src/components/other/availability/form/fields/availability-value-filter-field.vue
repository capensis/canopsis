<template>
  <v-layout class="availability-value-filter-field gap-3" d-inline-flex align-end>
    <availability-value-filter-method-field
      v-model="method"
      class="availability-value-filter-field__method"
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

import { AVAILABILITY_SHOW_TYPE } from '@/constants';

import { useModelField } from '@/hooks/form';

import AvailabilityValueFilterMethodField from './availability-value-filter-method-field.vue';

export default {
  components: { AvailabilityValueFilterMethodField },
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
    const { updateModel } = useModelField(props, emit);
    const isPercentType = computed(() => props.showType === AVAILABILITY_SHOW_TYPE.percent);

    const method = computed({
      get() {
        return props.value?.method;
      },
      set(newMethod) {
        updateModel({
          ...props.value,
          value: props.value?.value ?? 0,
          method: newMethod,
        });
      },
    });

    const clearValue = () => updateModel(undefined);

    return {
      method,

      isPercentType,
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

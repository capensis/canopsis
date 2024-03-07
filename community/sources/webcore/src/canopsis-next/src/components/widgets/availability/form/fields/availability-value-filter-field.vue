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
import { AVAILABILITY_SHOW_TYPE, AVAILABILITY_VALUE_FILTER_METHODS } from '@/constants';

import { formMixin } from '@/mixins/form';

export default {
  mixins: [formMixin],
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
  computed: {
    method: {
      get() {
        return this.value?.method;
      },
      set(method) {
        this.updateModel({ ...this.value, value: this.value?.value ?? 0, method });
      },
    },

    isPercentType() {
      return this.showType === AVAILABILITY_SHOW_TYPE.percent;
    },

    valueFilterMethods() {
      return Object.values(AVAILABILITY_VALUE_FILTER_METHODS).map(method => ({
        value: method,
        text: this.$t(`availability.valueFilterMethods.${method}`),
      }));
    },
  },
  methods: {
    clearValue() {
      this.updateModel(undefined);
    },
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

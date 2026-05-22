<template>
  <v-layout class="gap-2" column>
    <c-label :label="$t('stateSetting.computeMethod')" />
    <v-radio-group
      v-field="value"
      class="mt-0"
      hide-details
      mandatory
    >
      <v-radio
        v-for="method in STATE_SETTING_METHODS"
        :key="method"
        :value="method"
        color="primary"
      >
        <template #label>
          <span class="mr-2">
            {{ $t(`stateSetting.methods.${method}.label`) }}
          </span>
          <c-help-icon
            :text="$t(`stateSetting.methods.${method}.tooltip`)"
            top
          />
        </template>
      </v-radio>
    </v-radio-group>
  </v-layout>
</template>

<script>
import { STATE_SETTING_METHODS } from '@/constants';

export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: STATE_SETTING_METHODS.inherited,
    },
  },
  setup() {
    return {
      STATE_SETTING_METHODS,
    };
  },
};
</script>

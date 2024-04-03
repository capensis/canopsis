<template>
  <v-layout column>
    <span class="subheading font-weight-bold">{{ $t('stateSetting.computeMethod') }}</span>
    <v-radio-group
      v-field="value"
      hide-details
      mandatory
    >
      <v-radio
        v-for="method in methods"
        :key="method"
        :value="method"
        color="primary"
      >
        <template #label>
          <span class="mr-2">{{ $t(`stateSetting.methods.${method}.label`) }}</span>
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
  computed: {
    methods() {
      return STATE_SETTING_METHODS;
    },
  },
};
</script>

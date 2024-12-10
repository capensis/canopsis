<template>
  <v-radio-group v-field="value">
    <template #label="">
      {{ $t('common.type') }}
    </template>
    <v-radio
      v-for="type in types"
      :key="type"
      :value="type"
      :label="type"
      class="text-uppercase"
    />
  </v-radio-group>
</template>

<script>
import { ROLE_TYPES } from '@/constants';

export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: ROLE_TYPES.ui,
    },
  },
  setup() {
    const types = Object.values(ROLE_TYPES);

    return {
      types,
    };
  },
};
</script>

<template>
  <v-list-group>
    <template #activator="">
      <v-list-item>
        <span :class="validationHeaderClass">{{ $tc('common.map') }}</span>
      </v-list-item>
    </template>
    <v-container>
      <c-map-field
        v-field="value"
        name="parameters.map"
      />
    </v-container>
  </v-list-group>
</template>

<script>
import { formValidationHeaderMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [formValidationHeaderMixin],
  props: {
    value: {
      type: String,
      required: false,
    },
  },
};
</script>

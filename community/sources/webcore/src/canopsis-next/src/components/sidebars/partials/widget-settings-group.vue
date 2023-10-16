<template>
  <v-list-group class="widget-settings-group">
    <template #activator="">
      <v-list-item>
        <div :class="validationHeaderClass">
          {{ title }}
        </div>
      </v-list-item>
    </template>
    <v-list
      class="px-2 py-0 background darken-1"
      expand="expand"
    >
      <slot />
    </v-list>
  </v-list-group>
</template>

<script>
import { formValidationHeaderMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [formValidationHeaderMixin],
  props: {
    title: {
      type: String,
      required: true,
    },
  },
};
</script>

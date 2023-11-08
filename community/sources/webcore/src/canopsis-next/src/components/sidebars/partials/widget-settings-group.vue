<template>
  <v-list-group
    class="widget-settings-group"
    color=""
    eager
  >
    <template #activator="">
      <v-list-item-content>
        <v-list-item-title :class="validationHeaderClass">
          {{ title }}
        </v-list-item-title>
      </v-list-item-content>
    </template>
    <v-list
      class="px-2 py-0 background darken-1"
      expand
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

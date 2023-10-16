<template>
  <v-list-group>
    <template #activator="">
      <v-list-item class="widget-settings-item-title">
        <v-layout
          :class="validationHeaderClass"
          align-center
        >
          <slot name="title">
            {{ title }}
          </slot><span
            class="font-italic text-caption ml-1"
            v-if="optional || subTitle"
          ><span v-if="optional">({{ $t('common.optional') }})</span><span v-if="subTitle">{{ subTitle }}</span></span>
        </v-layout>
      </v-list-item>
    </template>
    <v-container>
      <slot />
    </v-container>
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
      required: false,
    },
    subTitle: {
      type: String,
      required: false,
    },
    optional: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

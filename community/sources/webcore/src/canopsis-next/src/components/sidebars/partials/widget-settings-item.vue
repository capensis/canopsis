<template>
  <v-list-group
    v-bind="$attrs"
    :eager="eager"
    color=""
  >
    <template #activator="">
      <v-list-item-content class="widget-settings-item-title">
        <v-list-item-title :class="validationHeaderClass">
          <slot name="title">
            {{ title }}
          </slot>
          <span
            v-if="optional || subTitle"
            class="font-italic text-caption ml-1"
          >
            <span v-if="optional" class="text-lowercase">({{ $t('common.optional') }})</span>
            <span v-if="subTitle">{{ subTitle }}</span>
          </span>
        </v-list-item-title>
      </v-list-item-content>
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
  inheritAttrs: false,
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
    eager: {
      type: Boolean,
      default: true,
    },
  },
};
</script>

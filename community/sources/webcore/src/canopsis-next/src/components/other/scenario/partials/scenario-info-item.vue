<template>
  <v-list-item class="text-subtitle-1">
    <template v-if="!hideIcon && !multiline" #prepend>
      <c-icon>{{ icon }}</c-icon>
    </template>

    <c-layout align="center">
      <div class="text-subtitle-1">
        <span class="font-weight-bold">
          {{ label }}
        </span>
        <span v-if="value"> : {{ value }} </span>
      </div>
    </c-layout>
  </v-list-item>
</template>

<script>
export default {
  props: {
    icon: {
      type: String,
      default: '',
    },
    label: {
      type: String,
      default: '',
    },
    value: {
      type: [String, Number],
      default: '',
    },
    hideIcon: {
      type: Boolean,
      default: false,
    },
    multiline: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

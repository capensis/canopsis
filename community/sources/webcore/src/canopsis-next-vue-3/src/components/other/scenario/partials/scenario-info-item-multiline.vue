<template>
  <v-list-item class="text-subtitle-1">
    <c-layout gap="2" column>
      <div class="text-subtitle-1">
        <span class="font-weight-bold">
          {{ label }}
        </span>
        <span v-if="value">
          :
        </span>
      </div>
      <c-layout
        v-if="value"
        class="text-grey"
        align="center"
        gap="2"
      >
        <c-icon
          v-if="!hideIcon"
          color="grey"
        >
          {{ icon }}
        </c-icon>
        <span>
          {{ value }}
        </span>
      </c-layout>
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
  },
};
</script>

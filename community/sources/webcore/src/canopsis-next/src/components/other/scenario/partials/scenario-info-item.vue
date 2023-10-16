<template>
  <v-list-item class="subheading">
    <v-list-item-avatar v-if="!hideIcon">
      <v-icon>{{ icon }}</v-icon>
    </v-list-item-avatar>
    <v-list-item-content>
      <v-layout align-center>
        <div class="subheading">
          <span class="font-weight-bold">{{ label }}</span>
          <span v-if="value">: {{ value }}</span>
        </div>
      </v-layout>
    </v-list-item-content>
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

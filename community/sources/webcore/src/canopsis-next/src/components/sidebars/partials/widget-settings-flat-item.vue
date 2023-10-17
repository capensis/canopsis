<template>
  <v-list-item>
    <v-list-item-content>
      <v-list-item-title>
        {{ title }}
      </v-list-item-title>
    </v-list-item-content>
    <v-list-item-action class="mr-3">
      <slot name="actions" />
    </v-list-item-action>
  </v-list-item>
</template>

<script>
export default {
  props: {
    title: {
      type: String,
      default: '',
    },
  },
};
</script>

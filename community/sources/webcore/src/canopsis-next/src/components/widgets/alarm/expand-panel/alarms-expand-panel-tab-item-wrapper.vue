<template>
  <v-layout class="pa-3">
    <v-flex :class="cardFlexClass">
      <v-card class="tab-item-card">
        <slot name="prepend" />
        <v-card-text>
          <slot />
        </v-card-text>
      </v-card>
    </v-flex>
  </v-layout>
</template>

<script>
export default {
  props: {
    cardFlexClass: {
      type: [String, Array, Object],
      required: true,
    },
  },
};
</script>

<style lang="scss" scoped>
.tab-item-card {
  margin: auto;
}
</style>

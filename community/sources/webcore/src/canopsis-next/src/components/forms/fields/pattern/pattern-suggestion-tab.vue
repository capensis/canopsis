<template>
  <v-layout
    class="fill-height"
    align-center
  >
    <v-layout
      class="pattern-suggestions__tabs__number"
      align-center
    >
      {{ number }}
    </v-layout>
    <v-expand-x-transition>
      <div
        v-show="active"
        class="pattern-suggestions__tabs__found"
      >
        {{ $t('pattern.foundEntities', { count: entitiesCount }) }}
      </div>
    </v-expand-x-transition>
  </v-layout>
</template>

<script>
export default {
  props: {
    number: {
      type: Number,
      required: true,
    },
    entitiesCount: {
      type: Number,
      required: true,
    },
    active: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

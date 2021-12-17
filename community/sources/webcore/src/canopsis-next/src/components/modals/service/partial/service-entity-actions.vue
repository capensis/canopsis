<template lang="pug">
  div
    v-tooltip(v-for="action in actions", :key="action.type", top)
      template(#activator="{ on }")
        v-btn(
          v-on="on",
          :disabled="action.disabled",
          depressed,
          small,
          light,
          @click.stop="$emit('apply', action)"
        )
          v-icon {{ action.icon }}
      span {{ $t(`serviceWeather.actions.${action.type}`) }}
</template>

<script>
export default {
  props: {
    actions: {
      type: Array,
      default: () => [],
    },
  },
};
</script>

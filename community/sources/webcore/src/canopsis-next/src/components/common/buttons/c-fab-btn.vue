<template lang="pug">
  div.fab
    v-layout(column)
      c-refresh-btn(@click="$listeners.refresh")
      slot(name="additional")
      v-tooltip(v-if="hasAccess && $listeners.create", left)
        template(#activator="{ on }")
          v-btn(
            v-on="on",
            :dark="dark",
            color="primary",
            fab,
            @click.stop="$listeners.create"
          )
            v-icon add
        slot
</template>

<script>
export default {
  props: {
    hasAccess: {
      type: Boolean,
      default: true,
    },
    dark: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      fab: false,
    };
  },
};
</script>

<template lang="pug">
  v-list-group
    v-list-tile(slot="activator")
      v-list-tile-title {{ actionNumber }} - {{ action.type }} - {{ action.name || action.from }}
      v-list-tile-action
        v-layout
          c-action-btn(type="edit", color="success", @click="$emit('edit', action)")
          c-action-btn(type="delete", color="error", @click="$listeners.remove")
    v-layout.px-3.py-2(column)
      div(v-for="(value, key) in action", :key="key") {{ $t(`common.${key}`) }}: {{ `${value}` }}
</template>

<script>
export default {
  props: {
    action: {
      type: Object,
      required: true,
    },
    actionNumber: {
      type: [Number, String],
      required: true,
    },
  },
};
</script>

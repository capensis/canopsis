<template lang="pug">
  v-layout(wrap)
    v-chip(
      v-for="entity in entities",
      :key="entity[itemKey]",
      close,
      @input="$emit('remove', entity)"
    )
      slot(:entity="entity") {{ contentKey ? entity[contentKey] : entity }}
    v-btn(v-if="clearable && entities.length", small, icon, @click="$emit('clear')")
      v-icon(color="error") delete
</template>

<script>
export default {
  props: {
    entities: {
      type: Array,
      required: true,
    },
    itemKey: {
      type: String,
      default: '_id',
    },
    contentKey: {
      type: String,
      required: false,
    },
    clearable: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

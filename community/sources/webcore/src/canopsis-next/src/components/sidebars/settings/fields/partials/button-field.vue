<template lang="pug">
  v-container.pa-3(fluid)
    v-layout(align-center, justify-space-between)
      slot(name="title")
      v-layout(justify-end)
        v-btn.primary(v-if="isEmpty", small, @click="$emit('create', $event)")
          span {{ $t('common.create') }}
        template(v-else)
          v-btn.primary(small, @click="$emit('edit', $event)")
            span {{ $t('common.edit') }}
          v-btn.error(small, @click="$emit('delete', $event)")
            v-icon delete
</template>

<script>
export default {
  props: {
    isEmpty: {
      type: Boolean,
      default: true,
    },
  },
};
</script>

<template lang="pug" functional>
  v-card(
    :class="{ 'fill-min-height': props.fillHeight }",
    :data-test="data.attrs && data.attrs['data-test']"
  )
    v-card-title.primary.white--text(v-if="$slots.title && !$slots.fullTitle")
      v-layout(justify-space-between, align-center)
        span.headline
          slot(name="title")
    slot(name="fullTitle")
    v-card-text(v-if="$slots.text")
      slot(name="text")
    template(v-if="$slots.actions")
      v-divider
      v-card-actions
        v-layout.py-1(justify-end)
          slot(name="actions")
</template>

<script>
export default {
  props: {
    fillHeight: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

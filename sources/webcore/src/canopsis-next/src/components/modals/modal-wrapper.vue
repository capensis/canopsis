<template lang="pug" functional>
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline
          slot(name="title")
    v-card-text
      slot(name="text")
    template(v-if="$slots.actions")
      v-divider
      v-card-actions
        v-layout.py-1(justify-end)
          slot(name="actions")
</template>

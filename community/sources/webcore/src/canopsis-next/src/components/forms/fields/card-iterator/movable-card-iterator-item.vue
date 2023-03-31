<template lang="pug">
  v-card
    v-layout.pt-2(justify-space-between)
      v-flex(xs5)
        span.pl-2
          v-btn(
            :disabled="disabledUp",
            icon,
            @click.prevent="$emit('up')"
          )
            v-icon arrow_upward
          v-btn(
            :disabled="disabledDown",
            icon,
            @click.prevent="$emit('down')"
          )
            v-icon arrow_downward
      v-flex.text-xs-right.pr-2(xs3)
        v-btn(icon, @click.prevent="$emit('remove')")
          v-icon(color="error") close
    v-layout.px-3.pb-3(justify-center, column)
      slot
</template>

<script>
export default {
  props: {
    disabledUp: {
      type: Boolean,
      default: false,
    },
    disabledDown: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

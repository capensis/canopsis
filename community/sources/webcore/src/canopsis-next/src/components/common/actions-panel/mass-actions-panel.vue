<template lang="pug">
  div
    mq-layout(mq="l+")
      c-action-btn(
        v-for="(action, index) in actions",
        :key="index",
        :tooltip="action.title",
        :disabled="action.disabled",
        :loading="action.loading",
        :icon="action.icon",
        :color="action.iconColor",
        :badge-value="action.badgeValue",
        :badge-tooltip="action.badgeTooltip",
        lazy,
        @click="action.method"
      )
    mq-layout(:mq="['m', 't']")
      v-menu(
        bottom,
        left,
        lazy,
        @click.native.stop=""
      )
        template(#activator="{ on }")
          v-btn.ma-0(v-on="on", icon)
            v-icon more_vert
        v-list
          v-list-tile(
            v-for="(action, index) in actions",
            :key="index",
            :disabled="action.disabled || action.loading",
            @click.stop="action.method"
          )
            v-list-tile-title
              v-icon.pr-3(
                :color="action.iconColor",
                :disabled="action.disabled",
                left,
                small
              ) {{ action.icon }}
              span.body-1(:class="action.cssClass") {{ action.title }}
</template>

<script>
export default {
  props: {
    actions: {
      type: Array,
      required: true,
    },
  },
};
</script>

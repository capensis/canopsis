<template>
  <div>
    <mq-layout mq="l+">
      <c-action-btn
        v-for="(action, index) in actions"
        :key="index"
        :tooltip="action.title"
        :disabled="action.disabled"
        :loading="action.loading"
        :icon="action.icon"
        :color="action.iconColor"
        :badge-value="action.badgeValue"
        :badge-tooltip="action.badgeTooltip"
        @click="action.method"
      />
    </mq-layout>
    <mq-layout :mq="['m', 't']">
      <v-menu
        bottom
        left
        @click.native.stop=""
      >
        <template #activator="{ on }">
          <v-btn
            icon
            v-on="on"
          >
            <v-icon>more_vert</v-icon>
          </v-btn>
        </template>
        <v-list>
          <v-list-item
            v-for="(action, index) in actions"
            :key="index"
            :disabled="action.disabled || action.loading"
            @click.stop="action.method"
          >
            <v-list-item-title>
              <v-icon
                :color="action.iconColor"
                :disabled="action.disabled"
                class="pr-3"
                left
                small
              >
                {{ action.icon }}
              </v-icon>
              <span
                :class="action.cssClass"
                class="text-body-1"
              >
                {{ action.title }}
              </span>
            </v-list-item-title>
          </v-list-item>
        </v-list>
      </v-menu>
    </mq-layout>
  </div>
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

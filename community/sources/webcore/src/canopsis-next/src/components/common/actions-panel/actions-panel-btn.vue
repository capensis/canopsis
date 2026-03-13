<template>
  <v-menu v-if="action.items?.length" bottom offset-y>
    <template #activator="{ on }">
      <c-action-btn
        :tooltip="action.title"
        :disabled="action.disabled"
        :loading="action.loading"
        :icon="action.icon"
        :color="action.iconColor"
        :badge-value="action.badgeValue"
        :badge-tooltip="action.badgeTooltip"
        v-on="on"
      />
    </template>
    <v-list>
      <v-list-item
        v-for="subAction in action.items"
        :key="subAction.title"
        @click="subAction.method"
      >
        <v-list-item-title>{{ subAction.title }}</v-list-item-title>
      </v-list-item>
    </v-list>
  </v-menu>
  <c-action-btn
    v-else
    :tooltip="action.title"
    :disabled="action.disabled"
    :loading="action.loading"
    :icon="action.icon"
    :color="action.iconColor"
    :badge-value="action.badgeValue"
    :badge-tooltip="action.badgeTooltip"
    @click="action.method"
  />
</template>

<script>
export default {
  props: {
    action: {
      type: Object,
      required: true,
    },
  },
};
</script>

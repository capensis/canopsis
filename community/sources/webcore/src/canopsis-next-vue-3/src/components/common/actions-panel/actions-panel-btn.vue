<template>
  <v-menu v-if="action.items?.length" location="bottom" offset="4">
    <template #activator="{ props: activatorProps }">
      <c-action-btn
        :tooltip="action.title"
        :disabled="action.disabled"
        :loading="action.loading"
        :icon="action.icon"
        :icon-size="action.iconSize"
        :color="action.iconColor"
        :small="small"
        v-bind="activatorProps"
      />
    </template>

    <v-list density="compact">
      <v-list-item
        v-for="subAction in action.items"
        :key="subAction.title"
        @click="subAction.method"
      >
        <template #title>
          {{ subAction.title }}
        </template>
      </v-list-item>
    </v-list>
  </v-menu>

  <c-action-btn
    v-else
    :tooltip="action.title"
    :disabled="action.disabled"
    :loading="action.loading"
    :icon="action.icon"
    :icon-size="action.iconSize"
    :color="action.iconColor"
    :small="small"
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
    small: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

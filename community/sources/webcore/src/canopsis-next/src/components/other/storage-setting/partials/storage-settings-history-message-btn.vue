<template>
  <v-tooltip :disabled="!history" bottom>
    <template #activator="{ on }">
      <v-btn
        v-bind="$attrs"
        v-on="on"
        @click="$emit('click')"
      >
        <slot />
      </v-btn>
    </template>
    <storage-settings-history-message
      v-if="history"
      :history="history"
      archived-count-message-key="storageSetting.history.entity.archivedCount"
      deleted-count-message-key="storageSetting.history.entity.deletedCount"
      hide-deleted
    />
  </v-tooltip>
</template>

<script>
import StorageSettingsHistoryMessage from './storage-settings-history-message.vue';

export default {
  components: { StorageSettingsHistoryMessage },
  inheritAttrs: false,
  props: {
    history: {
      type: [Object, Number],
      required: false,
    },
  },
};
</script>

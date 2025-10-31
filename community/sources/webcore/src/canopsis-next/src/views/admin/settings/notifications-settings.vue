<template>
  <v-container>
    <c-page>
      <template #header>
        {{ $t('parameters.tabs.notificationsSettings') }}
      </template>
      <v-card-text class="pa-4">
        <notifications-settings />
      </v-card-text>
    </c-page>
  </v-container>
</template>

<script>
import NotificationsSettings from '@/components/other/notification/notifications-settings.vue';

export default {
  components: { NotificationsSettings },
};
</script>

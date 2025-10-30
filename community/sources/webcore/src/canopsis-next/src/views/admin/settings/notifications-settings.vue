<template>
  <c-page>
    <template #header>
      {{ $t('parameters.tabs.notificationsSettings') }}
    </template>
    <notifications-settings />
  </c-page>
</template>

<script>
import NotificationsSettings from '@/components/other/notification/notifications-settings.vue';

export default {
  components: { NotificationsSettings },
};
</script>

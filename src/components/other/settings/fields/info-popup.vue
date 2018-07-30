<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{$t('settings.infoPopup')}}
    v-container
      info-popup-settings-item(:widget="widget")
</template>

<script>
import InfoPopupSettingsItem from '@/components/other/info-popup/settings-item.vue';

export default {
  components: { InfoPopupSettingsItem },
  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
};
</script>

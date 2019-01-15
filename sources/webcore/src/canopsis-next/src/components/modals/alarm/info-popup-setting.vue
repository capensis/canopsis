<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('settings.infoPopup.title') }}
    v-card-text
      v-layout(justify-end)
        v-btn Add
</template>

<script>
import { MODALS } from '@/constants';

export default {
  name: MODALS.infoPopupSetting,
};
</script>

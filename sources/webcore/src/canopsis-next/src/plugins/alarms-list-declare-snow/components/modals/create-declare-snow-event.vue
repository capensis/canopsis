<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.createAckEvent.title') }}
    v-card-text
      v-container
        h2 DECLARE SNOW MODULE
</template>

<script>
import modalInnerItemsMixin from '@/mixins/modal/inner-items';

/**
 * Modal to create an ack event
 */
export default {
  $_veeValidate: {
    validator: 'new',
  },
  mixins: [modalInnerItemsMixin],
};
</script>

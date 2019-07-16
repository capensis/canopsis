<template lang="pug">
  v-container
    h2.text-xs-center.my-3.display-1.font-weight-medium {{ $t('common.parameters') }}
    v-layout(row, wrap)
      v-flex(xs12)
        v-card.ma-2
          v-card-text
            user-interface-form
</template>

<script>
import UserInterfaceForm from '@/components/forms/user-interface.vue';

export default {
  components: { UserInterfaceForm },
};
</script>

<template lang="pug">
  p BONJOUR
</template>

<script>
// import { createNamespacedHelpers } from 'vuex';
// const { mapGetters } = createNamespacedHelpers('modal');

import { MODALS } from '@/constants';


export default {
  name: MODALS.contextInfos,
  props: {
    template: {
      type: String,
    },
  },
};
</script>

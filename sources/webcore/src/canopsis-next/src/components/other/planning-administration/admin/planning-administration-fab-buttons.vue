<template lang="pug">
  fab-expand-buttons(@refresh="$emit('refresh')")
</template>

<script>
import FabExpandButtons from '@/components/other/fab-buttons/fab-expand-buttons.vue';

export default {
  components: {
    FabExpandButtons,
  },
};
</script>

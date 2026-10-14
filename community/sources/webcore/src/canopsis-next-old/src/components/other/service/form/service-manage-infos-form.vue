<template>
  <manage-infos v-field="infos" />
</template>

<script>
import ManageInfos from '@/components/widgets/context/manage-infos.vue';

export default {
  components: { ManageInfos },
  model: {
    prop: 'infos',
    event: 'input',
  },
  props: {
    infos: {
      type: Array,
      default: () => [],
    },
  },
};
</script>

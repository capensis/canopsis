<template>
  <manage-infos
    :infos="infos"
    @update:infos="updateModel"
  />
</template>

<script>
import { useModelField } from '@/hooks/form/model-field';

import ManageInfos from '@/components/widgets/context/manage-infos.vue';

export default {
  components: { ManageInfos },

  props: {
    infos: {
      type: Array,
      default: () => [],
    },
  },

  emits: ['update:infos'],

  setup(props, { emit }) {
    const { updateModel } = useModelField(props, emit);

    return {
      updateModel,
    };
  },
};
</script>

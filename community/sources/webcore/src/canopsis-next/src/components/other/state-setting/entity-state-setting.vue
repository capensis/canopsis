<template>
  <v-layout column>
    <v-text-field
      :value="stateSetting?.title"
      :loading="stateSettingPending"
      disabled
    />
  </v-layout>
</template>

<script>
import { debounce, pick } from 'lodash';
import { ref, watch, onMounted, onBeforeUnmount } from 'vue';

import { useEntity } from '@/hooks/store/modules/entity';
import { usePendingHandler } from '@/hooks/query/pending';

export default {
  props: {
    form: {
      type: Object,
      required: true,
    },
    preparer: {
      type: Function,
      default: () => d => d,
    },
  },
  setup(props) {
    const stateSetting = ref();

    const { checkStateSetting: checkEntityStateSetting } = useEntity();

    const {
      pending: stateSettingPending,
      handler: checkStateSetting,
    } = usePendingHandler(async (data) => {
      const response = await checkEntityStateSetting({
        data: pick(data, ['_id', 'name', 'type', 'connector', 'infos', 'category', 'impact_level']),
      });

      stateSetting.value = response?.title ? response : undefined;
    });

    /**
     * Prepares entity data from the form and fetches the matching state setting when a name is present.
     *
     * @param {Object} form - Entity or service form passed from the parent component
     */
    const checkStateSettingByForm = (form) => {
      const data = props.preparer(form);

      if (!data.name) {
        return;
      }

      checkStateSetting(data);
    };

    const debouncedCheckStateSetting = debounce(checkStateSettingByForm, 500);

    watch(() => props.form, debouncedCheckStateSetting, { deep: true });

    onMounted(() => checkStateSettingByForm(props.form));
    onBeforeUnmount(() => debouncedCheckStateSetting.cancel());

    return {
      stateSetting,
      stateSettingPending,
    };
  },
};
</script>

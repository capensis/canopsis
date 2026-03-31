<template>
  <v-btn :disabled="!modalId" color="primary" @click="fixWithAi">
    <v-icon class="mr-2">
      $vuetify.icons.ai
    </v-icon>
    <span>{{ $t('pattern.fixJsonWithAi') }}</span>
  </v-btn>
</template>

<script>
import { computed, inject } from 'vue';

import { useSidebar } from '@/hooks/sidebar';

export default {
  props: {
    jsonString: {
      type: String,
      default: '',
    },
  },
  setup(props) {
    const modal = inject('$modal');
    const sidebar = useSidebar();

    const modalId = computed(() => modal.id);

    const fixWithAi = () => sidebar.updateConfig({
      id: modalId.value,
      config: { jsonString: props.jsonString },
    });

    return {
      modalId,
      fixWithAi,
    };
  },
};
</script>

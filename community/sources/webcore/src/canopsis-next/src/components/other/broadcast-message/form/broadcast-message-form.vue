<template>
  <div>
    <v-tabs v-model="activeTab" color="primary" align-tabs="center">
      <v-tab value="general">
        {{ t('common.general') }}
      </v-tab>
      <v-tab value="views">
        {{ t('broadcastMessage.viewsAndPages') }}
      </v-tab>
    </v-tabs>

    <v-window v-model="activeTab" class="pt-3">
      <v-window-item value="general">
        <broadcast-message-general-form />
      </v-window-item>
      <v-window-item value="views">
        <broadcast-message-views-form :tree-items="treeItems" />
      </v-window-item>
    </v-window>
  </div>
</template>

<script>
import { ref } from 'vue';

import { useI18n } from '@/hooks/i18n';

import BroadcastMessageGeneralForm from './broadcast-message-general-form.vue';
import BroadcastMessageViewsForm from './broadcast-message-views-form.vue';

export default {
  components: {
    BroadcastMessageGeneralForm,
    BroadcastMessageViewsForm,
  },

  props: {
    treeItems: {
      type: Object,
      default: () => ({}),
    },
  },

  setup() {
    const { t } = useI18n();
    const activeTab = ref('general');

    return {
      t,
      activeTab,
    };
  },
};
</script>

<template>
  <c-layout gap="3" column>
    <c-enabled-field v-model="enabled" with-background />

    <c-name-field v-model="name" autofocus required />

    <c-form-block>
      <c-form-block-row :label="t('common.fullscreen')" align-center>
        <c-enabled-field
          v-model="fullscreen"
          :label="t('common.fullscreen')"
          no-margin
          hide-details
        />
      </c-form-block-row>

      <c-form-block-row :label="t('modals.createPlaylist.manageTabs')" align-center>
        <v-btn
          :aria-label="t('modals.createPlaylist.manageTabs')"
          color="primary"
          variant="outlined"
          @click="showManageTabsModal"
          @keydown.enter="showManageTabsModal"
        >
          {{ t('modals.createPlaylist.manageTabs') }}
        </v-btn>
      </c-form-block-row>
    </c-form-block>

    <c-layout>
      <c-layout v-if="tabsPending" justify-center>
        <v-progress-circular color="primary" indeterminate />
      </c-layout>
      <div v-else class="w-100">
        <draggable-playlist-tabs :tabs="tabsList" :groups="groups" @update:tabs="updateTabsList" />
      </div>
    </c-layout>

    <c-alert v-if="hasTabsError" type="error">
      {{ t('modals.createPlaylist.errors.emptyTabs') }}
    </c-alert>
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import DraggablePlaylistTabs from '@/components/other/playlists/form/fields/draggable-playlist-tabs.vue';
import { MODALS } from '@/constants';
import { useFormContext, useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

export default {
  components: { DraggablePlaylistTabs },

  props: {
    groups: {
      type: Array,
      default: () => [],
    },
    tabsPending: {
      type: Boolean,
      default: false,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const modals = useModals();
    const form = useFormContext();
    const { name, enabled, fullscreen, tabsList } = useFormModels({
      name: 'name',
      enabled: 'enabled',
      fullscreen: 'fullscreen',
      tabsList: 'tabs_list',
    });

    const hasTabsError = computed(() => !!form.errors.value.tabs);

    const updateTabsList = (tabs) => {
      form.setValue('tabs_list', tabs);

      if (tabs.length) {
        form.setFieldError('tabs', undefined);
      }
    };

    const showManageTabsModal = () =>
      modals.show({
        name: MODALS.managePlaylistTabs,
        config: {
          groups: props.groups,
          selectedTabs: tabsList.value,
          action: updateTabsList,
        },
      });

    return {
      t,
      name,
      enabled,
      fullscreen,
      tabsList,
      hasTabsError,
      updateTabsList,
      showManageTabsModal,
    };
  },
};
</script>

<template>
  <c-layout gap="3" column>
    <c-enabled-field v-model="enabled" with-background />

    <c-name-field v-model="name" autofocus required />

    <c-form-block>
      <c-form-block-row :label="t('common.duration')">
        <c-duration-field
          v-model:duration="interval"
          :label="t('common.duration')"
          name="interval"
          required
        />
      </c-form-block-row>

      <c-form-block-row :label="t('common.fullscreen')" align-center>
        <c-enabled-field
          v-model="fullscreen"
          :label="t('common.fullscreen')"
          no-margin
          hide-details
        />
      </c-form-block-row>

      <c-form-block-row :label="t('modals.createPlaylist.manageTabs')" align-center>
        <c-btn-with-error
          :error="tabsErrorMessage"
          :aria-label="t('modals.createPlaylist.manageTabs')"
          variant="outlined"
          @click="showManageTabsModal"
          @keydown.enter="showManageTabsModal"
        >
          {{ t('modals.createPlaylist.manageTabs') }}
        </c-btn-with-error>
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
  </c-layout>
</template>

<script>
import { computed, watch } from 'vue';

import { MODALS } from '@/constants';

import { useFieldError, useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useValidationAttachRequiredForField } from '@/hooks/validator/validation-attach-required';

import DraggablePlaylistTabs from '@/components/other/playlists/form/fields/draggable-playlist-tabs.vue';

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
    const { name, enabled, fullscreen, interval, tabsList } = useFormModels({
      name: 'name',
      enabled: 'enabled',
      fullscreen: 'fullscreen',
      interval: 'interval',
      tabsList: 'tabs_list',
    });

    const tabsError = useFieldError('tabs_list');
    const tabsErrorMessage = computed(() =>
      tabsError.value ? t('modals.createPlaylist.errors.emptyTabs') : '',
    );
    const { asyncValidateRequiredRule } = useValidationAttachRequiredForField(
      'tabs_list',
      () => tabsList.value,
      false,
    );

    const updateTabsList = (tabs) => {
      tabsList.value = tabs;
    };

    watch(tabsList, asyncValidateRequiredRule);

    const showManageTabsModal = () => modals.show({
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
      interval,
      tabsList,
      tabsErrorMessage,
      updateTabsList,
      showManageTabsModal,
    };
  },
};
</script>

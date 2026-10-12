<template>
  <v-layout
    class="gap-3"
    column
  >
    <c-enabled-field
      v-field="form.enabled"
      with-background
    />

    <c-name-field
      v-field="form.name"
      autofocus
      required
    />

    <c-form-block>
      <c-form-block-row :label="$t('common.fullscreen')" align-center>
        <c-enabled-field
          v-field="form.fullscreen"
          :label="$t('common.fullscreen')"
          no-margin
          hide-details
        />
      </c-form-block-row>

      <c-form-block-row
        :label="$t('modals.createPlaylist.manageTabs')"
        align-center
      >
        <v-btn
          :aria-label="$t('modals.createPlaylist.manageTabs')"
          color="primary"
          outlined
          @click="showManageTabsModal"
          @keydown.enter="showManageTabsModal"
        >
          {{ $t('modals.createPlaylist.manageTabs') }}
        </v-btn>
      </c-form-block-row>
    </c-form-block>

    <v-layout>
      <v-layout
        v-if="tabsPending"
        justify-center
      >
        <v-progress-circular
          color="primary"
          indeterminate
        />
      </v-layout>
      <v-flex
        v-else
        xs12
      >
        <draggable-playlist-tabs v-field="form.tabs_list" />
      </v-flex>
    </v-layout>
    <c-alert
      :value="hasTabsError"
      type="error"
    >
      {{ $t('modals.createPlaylist.errors.emptyTabs') }}
    </c-alert>
  </v-layout>
</template>

<script>
import { computed, watch } from 'vue';

import { MODALS } from '@/constants';

import { useModelField } from '@/hooks/form/model-field';
import { useModals } from '@/hooks/modals';
import { useValidator } from '@/hooks/validator/validator';
import { useValidationAttachRequiredForField } from '@/hooks/validator/validation-attach-required';

import DraggablePlaylistTabs from '@/components/other/playlists/form/fields/draggable-playlist-tabs.vue';

const TABS_FIELD_NAME = 'tabs';

export default {
  inject: ['$validator'],
  components: { DraggablePlaylistTabs },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
    groups: {
      type: Array,
      default: () => [],
    },
    tabsPending: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const modals = useModals();
    const { errors } = useValidator();
    const { updateField } = useModelField(props, emit);

    const hasTabs = () => props.form.tabs_list.length > 0;

    const { asyncValidateRequiredRule } = useValidationAttachRequiredForField(
      TABS_FIELD_NAME,
      hasTabs,
      false,
    );

    const hasTabsError = computed(() => errors.has(TABS_FIELD_NAME));

    watch(() => props.form.tabs_list.length, asyncValidateRequiredRule);

    const showManageTabsModal = () => modals.show({
      name: MODALS.managePlaylistTabs,
      config: {
        groups: props.groups,
        selectedTabs: props.form.tabs_list,
        action: (tabs) => {
          updateField('tabs_list', tabs);
          asyncValidateRequiredRule();
        },
      },
    });

    return {
      hasTabsError,
      showManageTabsModal,
    };
  },
};
</script>

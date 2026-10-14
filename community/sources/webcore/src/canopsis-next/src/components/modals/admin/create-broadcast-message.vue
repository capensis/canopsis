<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper text-class="pa-0" close>
      <template #title>
        <span>{{ title }}</span>
      </template>

      <template #text>
        <c-layout column>
          <broadcast-message :message="message" :color="color" />
          <broadcast-message-form :tree-items="treeItems" class="pa-3" />
        </c-layout>
      </template>

      <template #actions>
        <v-btn variant="text" @click="closeModal">
          {{ t('common.cancel') }}
        </v-btn>

        <v-btn
          :disabled="isSubmitting"
          :loading="isSubmitting"
          color="primary"
          type="submit"
        >
          {{ submitLabel }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { computed, onMounted, watch } from 'vue';

import {
  BROADCAST_MESSAGE_VIEWS,
  BROADCAST_MESSAGE_VIEWS_FORM_BLOCKS,
  MAX_LIMIT,
} from '@/constants';

import {
  formToMessage,
  messageToForm,
  prepareMessageViews,
} from '@/helpers/entities/broadcast-message/form';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePlaylist } from '@/hooks/store/modules/playlist';
import { useView } from '@/hooks/store/modules/view';

import BroadcastMessageForm from '@/components/other/broadcast-message/form/broadcast-message-form.vue';
import BroadcastMessage from '@/components/other/broadcast-message/partials/broadcast-message.vue';

export default {
  components: {
    BroadcastMessage,
    BroadcastMessageForm,
  },

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const modals = useModals();
    const config = computed(() => props.modal.config ?? {});

    const isNew = computed(() => !config.value.message?._id);
    const title = computed(
      () => config.value.title || t('modals.createBroadcastMessage.create.title'),
    );
    const submitLabel = computed(() => t(isNew.value ? 'common.create' : 'common.edit'));

    const form = useForm({
      initialValues: messageToForm(config.value.message),
    });
    const { values, handleSubmit, isSubmitting, setFieldError, setFieldTouched } = form;

    const { viewGroups, fetchViewGroupsList } = useView();
    const { playlists, fetchPlaylistsList } = usePlaylist();

    const viewGroupsTree = computed(() => (viewGroups.value || []).map((group = {}) => ({
        value: group._id,
        name: group.title || group.name,
        children: (group.views || []).map(view => ({
          value: view._id,
          name: view.title || view.name,
        })),
      })),
    );

    const playlistsTree = computed(() => (playlists.value || []).map(playlist => ({
        value: playlist.name,
        name: playlist.name,
      })),
    );

    const treeItems = computed(() => ({
      pages: [
        {
          value: BROADCAST_MESSAGE_VIEWS.login,
          name: t('common.login'),
        },
        {
          value: BROADCAST_MESSAGE_VIEWS.exploitation,
          name: t('common.exploitation'),
        },
        {
          value: BROADCAST_MESSAGE_VIEWS.administration,
          name: t('common.administration'),
        },
        {
          value: BROADCAST_MESSAGE_VIEWS.notifications,
          name: t('common.notification', 2),
        },
        {
          value: BROADCAST_MESSAGE_VIEWS.profile,
          name: t('common.profile'),
        },
      ],
      views: [
        {
          value: BROADCAST_MESSAGE_VIEWS.allViews,
          name: t('broadcastMessage.allViews'),
          children: viewGroupsTree.value,
        },
      ],
      playlists: [
        {
          value: BROADCAST_MESSAGE_VIEWS.allPlaylists,
          name: t('broadcastMessage.allPlaylists'),
          children: playlistsTree.value,
        },
      ],
    }));

    const message = computed(
      () => values.value.message || t('modals.createBroadcastMessage.defaultMessage'),
    );
    const color = computed(() => values.value.color);

    const hasSelectedViews = () => Object.values(BROADCAST_MESSAGE_VIEWS_FORM_BLOCKS).some(
        block => values.value.views?.[block]?.length > 0,
      );

    const closeModal = () => modals.hide({ id: props.modal.id });

    const submit = handleSubmit(async () => {
      if (!hasSelectedViews()) {
        setFieldTouched('views', true);
        setFieldError('views', t('broadcastMessage.errors.viewsRequired'));

        return;
      }

      await config.value.action?.(formToMessage(values.value, treeItems.value));

      closeModal();
    });

    const { updateOriginalFormField } = useFormConfirmableCloseModal({
      form: values,
      submit,
      close: closeModal,
    });

    watch(
      treeItems,
      () => {
        const views = prepareMessageViews(values.value.views, treeItems.value);

        form.setValue('views', views);
        updateOriginalFormField('views', views);
      },
      { immediate: true },
    );

    onMounted(() => Promise.all([
        fetchViewGroupsList({
          params: {
            limit: MAX_LIMIT,
            with_views: true,
            with_tabs: true,
            with_widgets: true,
            with_flags: true,
            with_private: true,
          },
        }),
        fetchPlaylistsList({
          params: { limit: MAX_LIMIT },
        }),
      ]),
    );

    return {
      t,
      title,
      submitLabel,
      isSubmitting,
      message,
      color,
      treeItems,
      closeModal,
      submit,
    };
  },
};
</script>

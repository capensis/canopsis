<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper
      text-class="pa-0"
      close
    >
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <v-layout column>
          <broadcast-message
            :message="message"
            :color="form.color"
          />
          <broadcast-message-form
            v-model="form"
            :tree-items="treeItems"
            class="pa-3"
          />
        </v-layout>
      </template>
      <template #actions="">
        <v-btn
          depressed
          text
          @click="close"
        >
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          :disabled="isDisabled"
          class="primary white--text"
          type="submit"
        >
          {{ $t('common.submit') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { computed, ref, watch, onMounted } from 'vue';

import { MODALS, BROADCAST_MESSAGE_VIEWS } from '@/constants';

import { messageToForm, formToMessage, prepareMessageViews } from '@/helpers/entities/broadcast-message/form';

import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useViewGroup } from '@/hooks/store/modules/view';
import { usePlaylist } from '@/hooks/store/modules/playlist';

import BroadcastMessage from '@/components/other/broadcast-message/partials/broadcast-message.vue';
import BroadcastMessageForm from '@/components/other/broadcast-message/form/broadcast-message-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createBroadcastMessage,
  $_veeValidate: {
    validator: 'new',
  },
  components: { BroadcastMessage, BroadcastMessageForm, ModalWrapper },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { t, tc } = useI18n();
    const { config, close } = useInnerModal(props);

    const form = ref(messageToForm(config.value.message));
    const title = computed(() => config.value.title || t('modals.createBroadcastMessage.create.title'));
    const message = computed(() => form.value.message || t('modals.createBroadcastMessage.defaultMessage'));

    const { groups, fetchAllGroupsListWithWidgets } = useViewGroup();
    const { items: playlistItems, fetchList: fetchPlaylistList } = usePlaylist();

    const viewGroupsTree = computed(() => {
      if (!groups.value || !groups.value.length) {
        return [];
      }

      return groups.value.map((group = {}) => ({
        value: group._id,
        name: group.title || group.name,
        children: (group.views || []).map(view => ({
          value: view._id,
          name: view.title || view.name,
        })),
      }));
    });

    const playlistsTree = computed(() => {
      if (!playlistItems.value || !playlistItems.value.length) {
        return [];
      }

      return playlistItems.value.map(playlist => ({
        value: playlist.name,
        name: playlist.name,
      }));
    });

    const treeItems = computed(() => [
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
        name: tc('common.notification', 2),
      },
      {
        value: BROADCAST_MESSAGE_VIEWS.profile,
        name: t('common.profile'),
      },
      {
        value: BROADCAST_MESSAGE_VIEWS.allViews,
        name: t('broadcastMessage.allViews'),
        children: viewGroupsTree.value,
      },
      {
        value: BROADCAST_MESSAGE_VIEWS.allPlaylists,
        name: t('broadcastMessage.allPlaylists'),
        children: playlistsTree.value,
      },
    ]);

    const { submit, isDisabled } = useSubmittableForm({
      form,
      method: async () => {
        await config.value.action?.(formToMessage(form.value, treeItems.value));

        close();
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    watch(treeItems, (newTreeItems) => {
      form.value.views = prepareMessageViews(form.value.views, newTreeItems);
    });

    onMounted(async () => {
      try {
        await Promise.all([
          fetchAllGroupsListWithWidgets(),
          fetchPlaylistList(),
        ]);
      } catch (err) {
        console.error(err);
      }
    });

    return {
      form,
      title,
      message,
      isDisabled,
      treeItems,

      submit,
      close,
    };
  },
};
</script>

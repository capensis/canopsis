<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        <span>{{ title }}</span>
      </template>

      <template #text>
        <playlist-form :groups="playlistGroups" :tabs-pending="tabsPending" />
      </template>

      <template #actions>
        <v-btn variant="text" @click="closeModal">
          {{ t('common.cancel') }}
        </v-btn>

        <v-btn
:disabled="submitting"
:loading="submitting"
color="primary"
type="submit">
          {{ t('common.saveChanges') }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { computed, onMounted, ref } from 'vue';

import { MAX_LIMIT } from '@/constants';

import {
  formToPlaylist,
  getTabsFromGroupsByIds,
  playlistToForm,
} from '@/helpers/entities/playlist/form';

import { useForm } from '@/hooks/form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useView } from '@/hooks/store/modules/view';

import PlaylistForm from '@/components/other/playlists/form/playlist-form.vue';

export default {
  components: { PlaylistForm },

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const modals = useModals();
    const { fetchViewGroupsListWithoutStore } = useView();
    const config = computed(() => props.modal.config ?? {});
    const playlist = computed(() => config.value.playlist);
    const isNew = computed(() => !playlist.value?._id || config.value.duplicating);
    const title = computed(() => {
      const mode = config.value.duplicating ? 'duplicate' : isNew.value ? 'create' : 'edit';

      return config.value.title || t(`modals.createPlaylist.${mode}.title`);
    });
    const playlistGroups = ref([]);
    const tabsPending = ref(false);

    const form = useForm({
      initialValues: playlistToForm(playlist.value),
    });
    const { values, handleSubmit, submitting, reset, setFieldError, setFieldTouched } = form;

    const closeModal = () => modals.hide({ id: props.modal.id });

    const submit = handleSubmit(async () => {
      if (!values.value.tabs_list.length) {
        setFieldTouched('tabs', true);
        setFieldError('tabs', t('modals.createPlaylist.errors.emptyTabs'));

        return;
      }

      const result = await config.value.action?.(formToPlaylist(values.value));

      await config.value.afterSubmit?.(result);

      closeModal();

      return result;
    });

    const { updateOriginalForm } = useFormConfirmableCloseModal({
      form: values,
      submit,
      close: closeModal,
    });

    const fetchGroups = async () => {
      try {
        tabsPending.value = true;

        const { data = [] } = await fetchViewGroupsListWithoutStore({
          params: {
            limit: MAX_LIMIT,
            with_views: true,
            with_tabs: true,
          },
        });

        playlistGroups.value = data;

        if (playlist.value?.tabs_list?.length) {
          reset({
            ...values.value,

            tabs_list: getTabsFromGroupsByIds(data, playlist.value.tabs_list),
          });
          updateOriginalForm();
        }
      } finally {
        tabsPending.value = false;
      }
    };

    onMounted(fetchGroups);

    return {
      t,
      title,
      playlistGroups,
      tabsPending,
      submitting,
      closeModal,
      submit,
    };
  },
};
</script>

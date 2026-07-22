<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        <span>{{ t('modals.createPlaylist.manageTabs') }}</span>
      </template>

      <template #text>
        <manage-playlist-tabs-form v-model:selected-tabs="selectedTabs" :groups="groups" />
      </template>

      <template #actions>
        <v-btn variant="text" @click="closeModal">
          {{ t('common.cancel') }}
        </v-btn>

        <v-btn :disabled="submitting" :loading="submitting" color="primary" type="submit">
          {{ t('common.saveChanges') }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { computed, ref } from 'vue'

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal'
import { useI18n } from '@/hooks/i18n'
import { useModals } from '@/hooks/modals'
import { usePendingHandler } from '@/hooks/query/pending'

import ManagePlaylistTabsForm from '@/components/other/playlists/form/manage-playlist-tabs-form.vue'

export default {
  components: { ManagePlaylistTabsForm },

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n()
    const modals = useModals()
    const config = computed(() => props.modal.config ?? {})
    const selectedTabs = ref(config.value.selectedTabs ?? [])
    const groups = computed(() => config.value.groups ?? [])
    const closeModal = () => modals.hide({ id: props.modal.id })
    const { pending: submitting, handler: submit } = usePendingHandler(async () => {
      await config.value.action?.(selectedTabs.value)

      closeModal()
    })

    useFormConfirmableCloseModal({
      form: selectedTabs,
      submit,
      close: closeModal,
    })

    return {
      t,
      groups,
      selectedTabs,
      submitting,
      closeModal,
      submit,
    }
  },
}
</script>

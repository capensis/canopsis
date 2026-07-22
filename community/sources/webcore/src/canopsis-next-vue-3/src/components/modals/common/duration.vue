<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        <span>{{ title }}</span>
      </template>

      <template #text>
        <c-layout class="gap-3" align="start">
          <c-duration-field
            v-model:duration="form"
            :label="label"
            name="duration.value"
            autofocus
            required
          />
          <c-select
            v-model="form.unit"
            :items="units"
            :label="t('common.unit')"
            class="duration-modal__unit"
            item-title="text"
            item-value="value"
            name="duration.unit"
            required
          />
        </c-layout>
      </template>

      <template #actions>
        <v-btn :disabled="submitting" variant="text" @click="closeModal">
          {{ t('common.cancel') }}
        </v-btn>
        <v-btn :disabled="submitting" :loading="submitting" color="primary" type="submit">
          {{ t('common.submit') }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { computed, ref } from 'vue'

import { MODALS } from '@/constants'

import { durationToForm } from '@/helpers/date/duration'

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal'
import { useI18n } from '@/hooks/i18n'
import { useModals } from '@/hooks/modals'

export default {
  name: MODALS.duration,

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
    const form = ref(durationToForm(config.value.duration))
    const submitting = ref(false)

    const title = computed(() => config.value.title ?? t('common.duration'))
    const label = computed(() => config.value.label ?? t('common.duration'))
    const units = computed(() =>
      config.value.units?.map((unit) => ({ ...unit, text: t(unit.text, form.value.value) })),
    )

    const closeModal = () => modals.hide({ id: props.modal.id })

    const submit = async () => {
      try {
        submitting.value = true

        await config.value.action?.(form.value)
        closeModal()
      } finally {
        submitting.value = false
      }
    }

    useFormConfirmableCloseModal({
      form,
      submit,
      close: closeModal,
    })

    return {
      t,
      form,
      title,
      label,
      units,
      submitting,
      closeModal,
      submit,
    }
  },
}
</script>

<style lang="scss" scoped>
.duration-modal__unit {
  min-width: 160px;
}
</style>

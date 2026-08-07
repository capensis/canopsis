<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        <span>{{ title }}</span>
      </template>

      <template #text>
        <c-text-field
          v-model="titleField"
          :label="t('common.title')"
          name="title"
          autofocus
          required
        />
      </template>

      <template #actions>
        <v-btn
          :disabled="isSubmitting"
          variant="text"
          @click="closeModal"
        >
          {{ t('common.cancel') }}
        </v-btn>

        <v-btn
          :disabled="isSubmitting"
          :loading="isSubmitting"
          color="primary"
          type="submit"
        >
          {{ t('common.submit') }}
        </v-btn>

        <v-tooltip
          v-if="config.deletable"
          :disabled="group.deletable"
          location="top"
        >
          <template #activator="{ props: activatorProps }">
            <div
              class="ml-2"
              v-bind="activatorProps"
            >
              <v-btn
                :disabled="isSubmitting || !group.deletable"
                :loading="removing"
                color="error"
                @click="remove"
              >
                {{ t('common.delete') }}
              </v-btn>
            </div>
          </template>

          <span>{{ t('modals.group.errors.isNotEmpty') }}</span>
        </v-tooltip>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { computed, ref } from 'vue';

import { MODALS } from '@/constants';

import { groupToRequest } from '@/helpers/entities/view/form';

import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';

export default {
  name: MODALS.createGroup,

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const modals = useModals();
    const popups = usePopups();
    const removing = ref(false);
    const config = computed(() => props.modal.config ?? {});
    const group = computed(() => config.value.group ?? {});
    const title = computed(() => config.value.title ?? t('modals.group.create.title'));

    const form = useForm({
      initialValues: {
        title: group.value.title ?? '',
      },
    });
    const { values, model, handleSubmit, isSubmitting } = form;
    const titleField = model('title');

    const closeModal = () => modals.hide({ id: props.modal.id });

    const submit = handleSubmit(async () => {
      await config.value.action?.(groupToRequest({ ...group.value, ...values.value }));

      closeModal();
    });

    const remove = () => {
      modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            try {
              removing.value = true;

              await config.value.remove?.();

              closeModal();
            } catch (err) {
              console.warn(err);

              popups.error({ text: t('modals.group.errors.isNotEmpty') });
            } finally {
              removing.value = false;
            }
          },
        },
      });
    };

    return {
      t,
      config,
      group,
      title,
      titleField,
      isSubmitting,
      removing,
      submit,
      remove,
      closeModal,
    };
  },
};
</script>

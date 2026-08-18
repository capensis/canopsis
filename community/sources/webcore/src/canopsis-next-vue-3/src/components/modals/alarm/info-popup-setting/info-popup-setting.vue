<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        <span>{{ $t('modals.infoPopupSetting.title') }}</span>
      </template>

      <template #text>
        <c-layout justify="end">
          <v-tooltip location="left">
            <template #activator="{ props }">
              <v-btn
                color="primary"
                size="small"
                variant="outlined"
                icon
                v-bind="props"
                @click="addPopup"
              >
                <v-icon>add</v-icon>
              </v-btn>
            </template>
            <span>{{ $t('modals.infoPopupSetting.add') }}</span>
          </v-tooltip>
        </c-layout>

        <c-layout column>
          <v-card
            v-for="(popup, index) in form.popups"
            :key="`${popup.column}-${index}`"
            class="my-1 bg-secondary text-white"
            variant="flat"
          >
            <v-card-title>
              <c-layout justify="space-between" align="center">
                <div>{{ $t('modals.infoPopupSetting.column') }}: {{ popup.column }}</div>
                <c-layout>
                  <v-btn
                    size="small"
                    variant="text"
                    icon
                    @click="deletePopup(index)"
                  >
                    <v-icon color="error">
                      delete
                    </v-icon>
                  </v-btn>
                  <v-btn
                    size="small"
                    variant="text"
                    icon
                    @click="editPopup(index, popup)"
                  >
                    <v-icon color="primary">
                      edit
                    </v-icon>
                  </v-btn>
                </c-layout>
              </c-layout>
            </v-card-title>
            <v-card-text>
              <p>{{ $t('common.template', 1) }}:</p>
              <c-textarea :model-value="popup.template" disabled />
            </v-card-text>
          </v-card>
        </c-layout>
      </template>

      <template #actions>
        <v-btn :disabled="isSubmitting" variant="text" @click="closeModal">
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          :loading="isSubmitting"
          :disabled="isSubmitting"
          color="primary"
          type="submit"
        >
          {{ $t('common.submit') }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { cloneDeep } from 'lodash';
import { computed } from 'vue';

import { MODALS } from '@/constants';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useModals } from '@/hooks/modals';

export default {
  name: MODALS.infoPopupSetting,

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const modals = useModals();
    const config = computed(() => props.modal.config ?? {});
    const formContext = useForm({
      initialValues: {
        form: {
          popups: cloneDeep(config.value.infoPopups ?? []),
        },
      },
    });
    const { handleSubmit, isSubmitting, model, values } = formContext;
    const form = model('form');
    const closeModal = () => modals.hide({ id: props.modal.id });

    const addPopup = () => {
      modals.show({
        name: MODALS.addInfoPopup,
        config: {
          columns: config.value.columns,
          variables: config.value.variables,
          action: popup => form.value.popups.push(popup),
        },
      });
    };

    const deletePopup = (index) => {
      form.value.popups.splice(index, 1);
    };

    const editPopup = (index, popup) => {
      modals.show({
        name: MODALS.addInfoPopup,
        config: {
          columns: config.value.columns,
          variables: config.value.variables,
          popup,
          action: (editedPopup) => {
            form.value.popups.splice(index, 1, editedPopup);
          },
        },
      });
    };

    const submit = handleSubmit(async () => {
      await config.value.action?.(values.value.form.popups);

      closeModal();
    });

    useFormConfirmableCloseModal({
      form,
      submit,
      close: closeModal,
    });

    return {
      form,
      isSubmitting,
      addPopup,
      deletePopup,
      editPopup,
      submit,
      closeModal,
    };
  },
};
</script>

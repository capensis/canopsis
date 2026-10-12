<template>
  <v-form v-click-outside.zIndex="clickOutsideDirective" @submit.prevent.stop="submit">
    <v-card width="400">
      <v-card-title class="bg-primary text-white px-3 py-2">
        <c-layout class="w-100" align="center" justify="space-between">
          <h4 class="text-subtitle-1 ma-0">
            {{ title }}
          </h4>
          <v-btn
            class="ma-0"
            size="x-small"
            variant="text"
            icon
            @click="close"
          >
            <v-icon color="white">
              close
            </v-icon>
          </v-btn>
        </c-layout>
      </v-card-title>
      <v-card-text class="pt-4">
        <point-form
          :form="form"
          :coordinates="coordinates"
          :exists-entities="existsEntities"
          @update:form="updateForm"
          @fly:coordinates="emit('fly:coordinates', $event)"
        />
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn :disabled="isSubmitting" variant="text" @click="close">
          {{ t('common.cancel') }}
        </v-btn>
        <v-btn
          v-if="editing"
          :disabled="isSubmitting"
          color="error"
          variant="text"
          @click.stop="emit('remove')"
        >
          {{ t('common.delete') }}
        </v-btn>
        <v-btn
          :disabled="isDisabled"
          :loading="isSubmitting"
          color="primary"
          variant="elevated"
          type="button"
          @click.stop="submit"
        >
          {{ t('common.submit') }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-form>
</template>

<script>
import { isEqual } from 'lodash';
import { computed, watch } from 'vue';

import { MODALS } from '@/constants';

import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

import PointForm from './point-form.vue';

export default {
  components: { PointForm },
  props: {
    point: {
      type: Object,
      required: true,
    },
    editing: {
      type: Boolean,
      required: false,
    },
    coordinates: {
      type: Boolean,
      required: false,
    },
    existsEntities: {
      type: Array,
      default: () => [],
    },
  },
  emits: ['cancel', 'submit', 'remove', 'fly:coordinates'],
  setup(props, { emit }) {
    const { t } = useI18n();
    const modals = useModals();
    const {
      values: form,
      validate,
      isSubmitting,
      reset,
    } = useForm({
      initialValues: { ...props.point },
    });

    const isDisabled = computed(() => !form.value.entity && !form.value.map);
    const title = computed(() => (props.editing ? t('map.editPoint') : t('map.addPoint')));

    const close = () => emit('cancel');
    const updateForm = value => (form.value = value);

    const submit = async () => {
      const { valid } = await validate();

      if (valid && !isDisabled.value) {
        emit('submit', { ...form.value });
      }
    };

    const showConfirmationModal = () => {
      if (!isEqual(props.point, form.value)) {
        modals.show({
          id: props.point._id,
          name: MODALS.clickOutsideConfirmation,
          dialogProps: {
            persistent: true,
          },
          config: {
            action: (confirmed) => {
              if (confirmed) {
                return submit();
              }

              return close();
            },
          },
        });
      } else {
        close();
      }
    };

    const clickOutsideDirective = computed(() => ({
      handler: showConfirmationModal,
      include: () => [...document.querySelectorAll('.v-overlay__content')],
      closeConditional: () => true,
    }));

    watch(
      () => props.point,
      point => reset({ ...point }),
    );

    return {
      t,
      emit,
      form,
      title,
      isDisabled,
      isSubmitting,
      clickOutsideDirective,
      close,
      submit,
      updateForm,
    };
  },
};
</script>

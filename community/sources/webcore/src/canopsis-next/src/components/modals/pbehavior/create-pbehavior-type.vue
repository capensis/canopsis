<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <pbehavior-type-form
          v-model="form"
          :only-color="onlyColor"
          :pending-priority="pendingPriority"
        />
      </template>
      <template #actions="">
        <v-btn
          :disabled="submitting"
          depressed
          text
          @click="close"
        >
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          :disabled="submitting"
          :loading="submitting"
          class="primary"
          type="submit"
        >
          {{ submitLabel }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { ref, computed, onMounted } from 'vue';

import { MODALS, VALIDATION_DELAY } from '@/constants';

import { pbehaviorTypeToForm, formToPbehaviorType } from '@/helpers/entities/pbehavior/type/form';

import { useI18n } from '@/hooks/i18n';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { usePendingHandler } from '@/hooks/query/pending';
import { usePbehaviorType } from '@/hooks/store/modules/pbehavior-type';

import PbehaviorTypeForm from '@/components/other/pbehavior/types/form/pbehavior-type-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createPbehaviorType,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    ModalWrapper,
    PbehaviorTypeForm,
  },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { t } = useI18n();
    const { close, config } = useInnerModal(props);
    const { fetchNextPbehaviorTypePriority } = usePbehaviorType();

    const form = ref(pbehaviorTypeToForm(config.value.pbehaviorType));

    const pbehaviorType = computed(() => config.value.pbehaviorType);
    const onlyColor = computed(() => pbehaviorType.value?.default);
    const isNew = computed(() => !pbehaviorType.value?._id);

    const title = computed(() => (
      config.value.title || t(isNew.value
        ? 'modals.createPbehaviorType.create.title'
        : 'modals.createPbehaviorType.edit.title')
    ));

    const {
      pending: pendingPriority,
      handler: setMinimalPriority,
    } = usePendingHandler(async () => {
      const { priority } = await fetchNextPbehaviorTypePriority();

      form.value.priority = priority;
    });

    const { submit, submitting, submitLabel } = useSubmittableForm({
      form,
      item: config.value.pbehaviorType,
      method: async () => {
        if (config.value.action) {
          await config.value.action(formToPbehaviorType(form.value));
        }

        close();
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    onMounted(() => {
      if (isNew.value) {
        setMinimalPriority();
      }
    });

    return {
      close,
      form,
      title,
      pendingPriority,
      onlyColor,
      submitting,
      submitLabel,
      submit,
    };
  },
};
</script>

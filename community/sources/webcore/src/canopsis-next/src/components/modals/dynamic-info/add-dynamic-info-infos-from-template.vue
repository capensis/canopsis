<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ $t('modals.addDynamicInfoInfosFromTemplate.title') }}</span>
      </template>
      <template #text="">
        <c-select-field
          v-model="template"
          :disabled="templatesPending"
          :items="templates"
          :label="$t('modals.addDynamicInfoInfosFromTemplate.fields.template')"
          :loading="templatesPending"
          item-text="title"
          item-value="_id"
          name="template"
          return-object
          required
        />
      </template>
      <template #actions="">
        <v-btn
          type="button"
          depressed
          text
          @click="close"
        >
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          :disabled="templatesPending"
          :loading="submitting || templatesPending"
          class="primary"
          type="submit"
        >
          {{ $t('modals.addDynamicInfoInfosFromTemplate.actions.addInfos') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { onMounted, ref } from 'vue';

import { MODALS, VALIDATION_DELAY } from '@/constants';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useInnerModal } from '@/hooks/modals';
import { usePendingHandler } from '@/hooks/query/pending';
import { useDynamicInfoTemplates } from '@/hooks/store/modules/dynamic-info-templates';
import { useSubmittableForm } from '@/hooks/submittable-form';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.addDynamicInfoInfosFromTemplate,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    ModalWrapper,
  },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { config, close } = useInnerModal(props);

    const { fetchDynamicInfoTemplatesList } = useDynamicInfoTemplates();

    const template = ref(null);

    const templates = ref([]);

    const {
      pending: templatesPending,
      handler: fetchTemplatesList,
    } = usePendingHandler(async () => {
      templates.value = await fetchDynamicInfoTemplatesList();
    }, true);

    const { submit, submitting } = useSubmittableForm({
      form: template,
      method: async () => {
        await config.value.action?.(template);

        close();

        return template;
      },
    });

    useFormConfirmableCloseModal({ form: template, submit, close });

    onMounted(fetchTemplatesList);

    return {
      template,
      templates,
      templatesPending,
      submitting,
      submit,
      close,
    };
  },
};
</script>

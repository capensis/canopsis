<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper text-class="position-relative" close>
      <template #title>
        <span>{{ title }}</span>
      </template>

      <template #text>
        <c-layout v-if="noPattern" gap="2" column>
          <pbehavior-general-form
            :form="values"
            :with-inherited="withInherited"
            with-enabled
            @update:form="updateForm"
          />
        </c-layout>

        <pbehavior-form
          v-else
          :form="values"
          :with-inherited="withInherited"
          :pbehavior-id="pbehaviorId"
          pbehavior-counter-type
          @update:form="updateForm"
        />

        <ai-chat-sidebar
          v-if="chatShown"
          v-bind="chatOptions.bind"
          v-on="chatOptions.on"
        />
      </template>

      <template #actions>
        <v-btn variant="text" @click="closeModal">
          {{ t('common.cancel') }}
        </v-btn>

        <v-btn
          :disabled="submitting || chatOptions.bind.pending"
          :loading="submitting"
          color="primary"
          type="submit"
        >
          {{ submitLabel }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script setup>
import { computed, toRef } from 'vue';

import { LLM_SOCKET_CONTEXTS } from '@/constants';

import {
  formToPbehavior,
  pbehaviorToForm,
  pbehaviorToRequest,
} from '@/helpers/entities/pbehavior/form';

import { useAiChatForm } from '@/hooks/ai/ai-chat-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

import AiChatSidebar from '@/components/other/llm/chat/ai-chat-sidebar.vue';
import PbehaviorForm from '@/components/other/pbehavior/pbehaviors/form/pbehavior-form.vue';
import PbehaviorGeneralForm from '@/components/other/pbehavior/pbehaviors/form/pbehavior-general-form.vue';

const props = defineProps({
  modal: {
    type: Object,
    required: true,
  },
});

const { t } = useI18n();
const modals = useModals();

const config = computed(() => props.modal.config ?? {});
const pbehavior = computed(() => config.value.pbehavior);
const isNew = computed(() => !pbehavior.value?._id);
const title = computed(
  () =>
    config.value.title ||
    t(`modals.createPbehavior.${isNew.value ? 'create' : 'edit'}.title`),
);
const noPattern = computed(() => !!config.value.noPattern);
const withInherited = computed(() => !!config.value.withInherited);
const pbehaviorId = computed(() => pbehavior.value?._id);
const submitLabel = computed(() => t(isNew.value ? 'common.create' : 'common.edit'));

const form = useForm({
  initialValues: pbehaviorToForm(pbehavior.value, null, config.value.timezone),
});
const { values, handleSubmit, submitting } = form;

const {
  shown: chatShown,
  options: chatOptions,
} = useAiChatForm({
  form: values,
  modal: toRef(props, 'modal'),
  ruleId: pbehaviorId,
  context: LLM_SOCKET_CONTEXTS.pbehavior,
});

const closeModal = () => modals.hide({ id: props.modal.id });
const updateForm = (value) => {
  values.value = value;
};

const submit = handleSubmit(async () => {
  const result = await config.value.action?.(
    pbehaviorToRequest(formToPbehavior(values.value, config.value.timezone)),
  );

  await config.value.afterSubmit?.(result);

  closeModal();

  return result;
});

useFormConfirmableCloseModal({
  form: values,
  submit,
  close: closeModal,
});
</script>

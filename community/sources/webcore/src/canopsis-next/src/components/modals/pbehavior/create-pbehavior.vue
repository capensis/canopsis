<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper text-class="position-relative" close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <v-layout
          v-if="noPattern"
          class="gap-2"
          column
        >
          <pbehavior-general-form
            v-field="form"
            :with-inherited="withInherited"
            with-enabled
          />
        </v-layout>
        <pbehavior-form
          v-else
          v-model="form"
          :with-inherited="withInherited"
          :pbehavior-id="pbehaviorId"
          pbehavior-counter-type
        />
        <ai-chat-sidebar
          v-if="chatShown"
          v-bind="chatOptions.bind"
          v-on="chatOptions.on"
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
          :disabled="submitting || chatOptions.bind.pending"
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
import { computed, ref, toRef } from 'vue';

import { LLM_SOCKET_CONTEXTS, MODALS, VALIDATION_DELAY } from '@/constants';

import { pbehaviorToForm, formToPbehavior, pbehaviorToRequest } from '@/helpers/entities/pbehavior/form';

import { useAiChatForm } from '@/hooks/ai/ai-chat-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';

import AiChatSidebar from '@/components/other/llm/chat/ai-chat-sidebar.vue';
import PbehaviorForm from '@/components/other/pbehavior/pbehaviors/form/pbehavior-form.vue';
import PbehaviorGeneralForm from '@/components/other/pbehavior/pbehaviors/form/pbehavior-general-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createPbehavior,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: { PbehaviorForm, PbehaviorGeneralForm, ModalWrapper, AiChatSidebar },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { t } = useI18n();
    const { config, close } = useInnerModal(props);

    const { pbehavior, timezone } = props.modal.config;

    const form = ref(pbehaviorToForm(pbehavior, null, timezone));

    const title = computed(() => config.value.title || t('modals.createPbehavior.create.title'));
    const noPattern = computed(() => !!config.value.noPattern);
    const withInherited = computed(() => !!config.value.withInherited);
    const pbehaviorId = computed(() => config.value.pbehavior?._id);

    const {
      shown: chatShown,
      options: chatOptions,
    } = useAiChatForm({
      form,

      modal: toRef(props, 'modal'),
      ruleId: pbehaviorId,
      context: LLM_SOCKET_CONTEXTS.pbehavior,
    });

    const { submit, submitting, submitLabel } = useSubmittableForm({
      form,
      item: config.value.pbehavior,
      method: async () => {
        const result = await config.value.action?.(
          pbehaviorToRequest(formToPbehavior(form.value, config.value.timezone)),
        );

        await config.value.afterSubmit?.(result);

        close();

        return result;
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    return {
      form,
      title,
      noPattern,
      withInherited,
      pbehaviorId,
      submitting,
      chatShown,
      chatOptions,
      submitLabel,
      submit,
      close,
    };
  },
};
</script>

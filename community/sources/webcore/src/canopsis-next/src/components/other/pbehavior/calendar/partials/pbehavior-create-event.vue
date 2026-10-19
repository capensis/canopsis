<template>
  <v-form
    v-click-outside.zIndex="clickOutsideDirective"
    class="pbehavior-form position-relative"
    @submit.prevent="submit"
  >
    <div class="pbehavior-form__body">
      <pbehavior-general-form
        v-if="entityPattern"
        :form="form"
        :with-inherited="withInherited"
        :no-timezone="noTimezone"
        with-enabled
        @update:form="form = $event"
      />

      <pbehavior-form
        v-else
        :form="form"
        :pbehavior-id="pbehavior?._id"
        :with-inherited="withInherited"
        :no-timezone="noTimezone"
        pbehavior-counter-type
        @update:form="form = $event"
      />
    </div>

    <ai-chat-sidebar v-if="chatShown" v-bind="chatOptions.attrs" />

    <c-layout class="pbehavior-form__actions" column>
      <v-divider />

      <c-layout class="pa-3" gap="2" justify="end">
        <v-btn v-show="pbehavior" color="error" @click="remove">
          {{ $t('common.delete') }}
        </v-btn>

        <v-btn variant="text" @click="cancel">
          {{ $t('common.cancel') }}
        </v-btn>

        <v-btn :disabled="chatOptions.attrs.pending" color="primary" type="submit">
          {{ $t('common.submit') }}
        </v-btn>
      </c-layout>
    </c-layout>
  </v-form>
</template>

<script>
import { cloneDeep } from 'lodash';
import { computed, onBeforeUnmount, onMounted, ref } from 'vue';

import { LLM_SOCKET_CONTEXTS, MODALS } from '@/constants';

import { isOmitEqual } from '@/helpers/collection';
import { getLocalTimezone } from '@/helpers/date/date';
import {
  calendarEventToPbehaviorForm,
  formToCalendarEvent,
} from '@/helpers/entities/pbehavior/form';

import { useAiChatForm } from '@/hooks/ai/ai-chat-form';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

import AiChatSidebar from '@/components/other/llm/chat/ai-chat-sidebar.vue';
import PbehaviorForm from '@/components/other/pbehavior/pbehaviors/form/pbehavior-form.vue';
import PbehaviorGeneralForm from '@/components/other/pbehavior/pbehaviors/form/pbehavior-general-form.vue';

export default {
  components: { PbehaviorForm, PbehaviorGeneralForm, AiChatSidebar },

  props: {
    event: {
      type: Object,
      required: false,
    },
    entityPattern: {
      type: Array,
      required: false,
    },
    defaultFields: {
      type: Object,
      required: false,
    },
    withInherited: {
      type: Boolean,
      default: false,
    },
    timezone: {
      type: String,
      default: getLocalTimezone(),
    },
    noTimezone: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['close', 'submit', 'remove'],

  setup(props, { emit }) {
    const initialValues = calendarEventToPbehaviorForm(
      props.event,
      props.entityPattern,
      props.defaultFields,
      props.timezone,
    );

    const { handleSubmit } = useForm({ initialValues });

    const modals = useModals();
    const { t } = useI18n();

    const manualClose = ref(false);
    const form = ref(initialValues);

    const pbehavior = computed(() => props.event?.data?.pbehavior);
    const closeConfirmationModalId = computed(
      () => `pbehavior-close-confirmation-${props.event?.id ?? 'new'}`,
    );

    const ruleId = computed(() => {
      const { _id: id = '' } = props.event?.data?.pbehavior ?? {};

      return id.startsWith('pbehavior_') ? '' : id;
    });

    const { shown: chatShown, options: chatOptions } = useAiChatForm({
      form,
      ruleId,
      context: LLM_SOCKET_CONTEXTS.pbehavior,
    });

    const cacheForm = () => {
      // eslint-disable-next-line no-param-reassign, vue/no-mutating-props
      props.event.data.cachedForm = cloneDeep(form.value);
    };

    onMounted(cacheForm);

    onBeforeUnmount(() => {
      if (manualClose.value) {
        // eslint-disable-next-line no-param-reassign, vue/no-mutating-props
        delete props.event.data.cachedForm;
      } else {
        cacheForm();
      }
    });

    const close = (event, manual = false) => {
      manualClose.value = manual;

      emit('close', event);
    };

    const cancel = (event) => {
      const { cachedForm } = props.event.data;

      if (isOmitEqual(cachedForm, form.value, ['_id'])) {
        return close(event, true);
      }

      if (modals.isModalOpenedById(closeConfirmationModalId.value)) {
        return undefined;
      }

      return modals.show({
        id: closeConfirmationModalId.value,
        name: MODALS.confirmation,
        config: {
          text: t('modals.createPbehavior.cancelConfirmation'),
          action: () => close(event, true),
        },
      });
    };

    const clickOutsideDirective = computed(() => {
      const selectorsForInclude = [
        '.ai-chat-sidebar-wrapper',
        '.ai-chat-sidebar__progress',
        '.c-calendar__popover-wrapper',
        '.v-overlay__content',
      ];

      return {
        handler: cancel,
        include: () => [...document.querySelectorAll(selectorsForInclude.join(','))]
          .filter(element => !element.querySelector('.c-calendar')),
        closeConditional: () => true,
      };
    });

    const remove = (event) => {
      emit('remove', pbehavior.value);
      close(event);
    };

    const submit = handleSubmit(async () => {
      const calendarEvent = formToCalendarEvent(form.value, props.event);

      manualClose.value = true;

      emit('submit', calendarEvent);
    });

    return {
      form,
      pbehavior,
      clickOutsideDirective,
      chatShown,
      chatOptions,
      cancel,
      remove,
      submit,
    };
  },
};
</script>

<style lang="scss" scoped>
.pbehavior-form {
  display: flex;
  flex-direction: column;
  flex: 0 1 auto;
  width: 100%;
  max-height: 85vh;
  min-height: 0;
  overflow: hidden;

  &__body {
    flex: 0 1 auto;
    min-height: 0;
    overflow-y: auto;
    padding: 16px;
  }

  &__actions {
    flex: 0 0 auto;
  }
}
</style>

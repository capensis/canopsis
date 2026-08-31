<template>
  <v-form
    v-click-outside.zIndex="clickOutsideDirective"
    class="pbehavior-form position-relative"
    @submit.prevent="submitHandler"
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

        <v-btn :disabled="errors.any() || chatOptions.attrs.pending" color="primary" type="submit">
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

import { getMenuClassByCalendarEvent } from '@/helpers/calendar/calendar';
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
import { useValidator } from '@/hooks/validator/validator';

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

    useForm({ initialValues });

    const modals = useModals();
    const validator = useValidator();
    const { t } = useI18n();

    const manualClose = ref(false);
    const form = ref(initialValues);

    const pbehavior = computed(() => props.event?.data?.pbehavior);

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

      return modals.show({
        name: MODALS.confirmation,
        config: {
          text: t('modals.createPbehavior.cancelConfirmation'),
          action: () => close(event, true),
        },
      });
    };

    const clickOutsideDirective = computed(() => {
      const selectorsForInclude = [
        '.c-calendar__today-btn',
        '.c-calendar__pagination',
        '.c-calendar__menu-right',
        '.v-event',
        '.v-overlay__content',
        `.${getMenuClassByCalendarEvent(props.event.id)}`,
      ];

      return {
        handler: cancel,
        include: () => [...document.querySelectorAll(selectorsForInclude.join(','))],
        closeConditional: () => true,
      };
    });

    const remove = (event) => {
      emit('remove', pbehavior.value);
      close(event);
    };

    const submitHandler = async (event) => {
      const isValid = await validator.validateAll();

      if (isValid) {
        const calendarEvent = formToCalendarEvent(form.value, props.event);

        manualClose.value = true;

        emit('submit', calendarEvent, event);
      }
    };

    return {
      form,
      pbehavior,
      errors: validator.errors,
      clickOutsideDirective,
      chatShown,
      chatOptions,
      cancel,
      remove,
      submitHandler,
    };
  },
};
</script>

<style lang="scss" scoped>
.pbehavior-form {
  display: flex;
  flex-direction: column;
  flex: 1 1 auto;
  width: 100%;
  max-height: 100%;
  min-height: 0;
  overflow: hidden;

  &__body {
    flex: 1 1 auto;
    min-height: 0;
    overflow-y: auto;
    padding: 16px;
  }

  &__actions {
    flex: 0 0 auto;
  }
}
</style>

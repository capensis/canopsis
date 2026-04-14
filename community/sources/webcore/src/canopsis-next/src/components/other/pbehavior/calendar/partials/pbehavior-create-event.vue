<template>
  <v-form
    v-click-outside.zIndex="clickOutsideDirective"
    class="pbehavior-form position-relative"
    @submit.prevent="submitHandler"
  >
    <pattern-progress
      v-if="chatPending"
      :in-progress-text="chatPendingTexts.inProgress"
      :cancel-button-label="chatPendingTexts.cancel"
      @cancel="chatCancel"
    />
    <pbehavior-form
      v-model="form"
      :pbehavior-id="pbehavior?._id"
      :no-pattern="!!entityPattern"
      :with-inherited="withInherited"
      :no-timezone="noTimezone"
      class="py-3"
      pbehavior-counter-type
    />
    <v-layout
      class="pbehavior-form__actions"
      justify-end
    >
      <v-btn
        v-show="pbehavior"
        :outlined="$system.dark"
        class="error"
        @click="remove"
      >
        {{ $t('common.delete') }}
      </v-btn>
      <v-btn
        depressed
        text
        @click="cancel"
      >
        {{ $t('common.cancel') }}
      </v-btn>
      <v-btn
        :disabled="errors.any()"
        color="primary"
        type="submit"
      >
        {{ $t('common.submit') }}
      </v-btn>
    </v-layout>
  </v-form>
</template>

<script>
import { cloneDeep } from 'lodash';
import { computed, ref, onMounted, onBeforeUnmount } from 'vue';
import dependentMixin from 'vuetify/lib/mixins/dependent';

import { LLM_SOCKET_CONTEXTS, MODALS, VALIDATION_DELAY } from '@/constants';

import { calendarEventToPbehaviorForm, formToCalendarEvent } from '@/helpers/entities/pbehavior/form';
import { isOmitEqual } from '@/helpers/collection';
import { getMenuClassByCalendarEvent } from '@/helpers/calendar/calendar';
import { getLocalTimezone } from '@/helpers/date/date';

import { useAiChatForm } from '@/hooks/ai/ai-chat-form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useValidator } from '@/hooks/validator/validator';
import { useComponentInstance } from '@/hooks/vue';

import PbehaviorForm from '@/components/other/pbehavior/pbehaviors/form/pbehavior-form.vue';
import PatternProgress from '@/components/forms/fields/pattern/pattern-progress.vue';

export default {
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  inject: ['$system'],
  components: { PbehaviorForm, PatternProgress },
  mixins: [dependentMixin],
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
  setup(props, { emit }) {
    const vm = useComponentInstance();
    const modals = useModals();
    const validator = useValidator();
    const { t } = useI18n();

    const manualClose = ref(false);
    const form = ref(
      calendarEventToPbehaviorForm(props.event, props.entityPattern, props.defaultFields, props.timezone),
    );

    const pbehavior = computed(() => props.event?.data?.pbehavior);

    const aiChatSidebarId = `pbehavior-calendar-event-${props.event.id}`;
    const ruleId = computed(() => props.event?.data?.pbehavior?._id);

    const aiChatPatternsForm = computed({
      get: () => form.value.patterns,
      set: (patterns) => {
        form.value = { ...form.value, patterns };
      },
    });

    const {
      pending: chatPending,
      pendingTexts: chatPendingTexts,
      cancel: chatCancel,
    } = useAiChatForm({
      form: aiChatPatternsForm,
      modalId: aiChatSidebarId,
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
        `.${getMenuClassByCalendarEvent(props.event.id)}`,
      ];

      return {
        handler: cancel,
        include: () => [
          ...vm.getOpenDependentElements(),
          ...document.querySelectorAll(selectorsForInclude.join(',')),
        ],
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
      clickOutsideDirective,
      chatPending,
      chatPendingTexts,
      chatCancel,
      cancel,
      remove,
      submitHandler,
    };
  },
};
</script>

<style lang="scss" scoped>
  .pbehavior-form {
    width: 100%;

    &__actions {
      gap: 6px;
    }
  }
</style>

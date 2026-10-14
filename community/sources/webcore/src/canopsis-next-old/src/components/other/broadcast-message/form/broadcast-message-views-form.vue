<template>
  <v-layout class="gap-3 py-2" column>
    <c-form-block>
      <c-form-block-row :label="$t('common.pages')" :error="hasAnyErrors" indented>
        <broadcast-message-views-treeview-field
          v-field="views.pages"
          :items="treeItems.pages"
          @input="asyncValidateRequiredRule"
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('common.views')" :error="hasAnyErrors" indented>
        <broadcast-message-views-treeview-field
          v-field="views.views"
          :items="treeItems.views"
          @input="asyncValidateRequiredRule"
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('common.playlists')" :error="hasAnyErrors" indented>
        <broadcast-message-views-treeview-field
          v-field="views.playlists"
          :items="treeItems.playlists"
          @input="asyncValidateRequiredRule"
        />
      </c-form-block-row>
    </c-form-block>
    <v-messages :value="errorMessages" color="error" />
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { BROADCAST_MESSAGE_VIEWS_FORM_BLOCKS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useValidator } from '@/hooks/validator/validator';
import { useValidationAttachRequiredForField } from '@/hooks/validator/validation-attach-required';

import BroadcastMessageViewsTreeviewField from './fields/broadcast-message-views-treeview-field.vue';

export default {
  inject: ['$validator'],
  components: {
    BroadcastMessageViewsTreeviewField,
  },
  model: {
    prop: 'views',
    event: 'input',
  },
  props: {
    views: {
      type: Object,
      default: () => ({}),
    },
    treeItems: {
      type: Object,
      default: () => ({}),
    },
  },
  setup(props) {
    const { t } = useI18n();
    const { errors } = useValidator();

    const hasSelectedViews = () => Object.values(BROADCAST_MESSAGE_VIEWS_FORM_BLOCKS)
      .some(block => props.views[block]?.length > 0);

    const { asyncValidateRequiredRule } = useValidationAttachRequiredForField('views', hasSelectedViews);

    const hasAnyErrors = computed(() => errors.has('views'));
    const errorMessages = computed(() => (hasAnyErrors.value ? [t('broadcastMessage.errors.viewsRequired')] : []));

    return {
      hasAnyErrors,
      errorMessages,

      asyncValidateRequiredRule,
    };
  },
};
</script>

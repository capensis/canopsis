<template>
  <c-layout class="py-2" gap="3" column>
    <c-form-block>
      <c-form-block-row :label="t('common.pages')" :error="hasAnyErrors" indented>
        <broadcast-message-views-treeview-field v-model="pages" :items="treeItems.pages" />
      </c-form-block-row>

      <c-form-block-row :label="t('common.views')" :error="hasAnyErrors" indented>
        <broadcast-message-views-treeview-field v-model="views" :items="treeItems.views" />
      </c-form-block-row>

      <c-form-block-row :label="t('common.playlists')" :error="hasAnyErrors" indented>
        <broadcast-message-views-treeview-field v-model="playlists" :items="treeItems.playlists" />
      </c-form-block-row>
    </c-form-block>

    <v-messages :messages="errorMessages" color="error" />
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { useFormContext } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import BroadcastMessageViewsTreeviewField from './fields/broadcast-message-views-treeview-field.vue';

export default {
  components: {
    BroadcastMessageViewsTreeviewField,
  },

  props: {
    treeItems: {
      type: Object,
      default: () => ({}),
    },
  },

  setup() {
    const { t } = useI18n();
    const form = useFormContext();

    const { pages, views, playlists } = form.models({
      pages: 'views.pages',
      views: 'views.views',
      playlists: 'views.playlists',
    });

    const hasAnyErrors = computed(() => Boolean(form.errors.value.views));
    const errorMessages = computed(() =>
      hasAnyErrors.value ? [t('broadcastMessage.errors.viewsRequired')] : [],
    );

    return {
      t,
      pages,
      views,
      playlists,
      hasAnyErrors,
      errorMessages,
    };
  },
};
</script>

<template>
  <v-tabs
    slider-color="primary"
    centered
  >
    <v-tab :class="{ 'error--text': hasGeneralError }">
      {{ $t('common.general') }}
    </v-tab>
    <v-tab :class="{ 'error--text': hasPagesError }">
      {{ $t('broadcastMessage.viewsAndPages') }}
    </v-tab>

    <v-tab-item eager>
      <broadcast-message-general-form
        v-field="form"
        ref="generalElement"
      />
    </v-tab-item>
    <v-tab-item eager>
      <broadcast-message-pages-form
        v-field="form.views"
        ref="pagesElement"
      />
    </v-tab-item>
  </v-tabs>
</template>

<script>
import { ref } from 'vue';

import { useValidationElementChildren } from '@/hooks/validator/validation-element-children';

import BroadcastMessageGeneralForm from './broadcast-message-general-form.vue';
import BroadcastMessagePagesForm from './broadcast-message-pages-form.vue';

export default {
  inject: ['$validator'],
  components: {
    BroadcastMessageGeneralForm,
    BroadcastMessagePagesForm,
  },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
  },
  setup() {
    const generalElement = ref(null);
    const pagesElement = ref(null);

    const { hasChildrenError: hasGeneralError } = useValidationElementChildren(generalElement);
    const { hasChildrenError: hasPagesError } = useValidationElementChildren(pagesElement);

    return {
      generalElement,
      pagesElement,

      hasGeneralError,
      hasPagesError,
    };
  },
};
</script>

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
      <broadcast-message-views-form
        v-field="form.views"
        ref="pagesElement"
        :tree-items="treeItems"
      />
    </v-tab-item>
  </v-tabs>
</template>

<script>
import { ref } from 'vue';

import { useValidationElementChildren } from '@/hooks/validator/validation-element-children';

import BroadcastMessageGeneralForm from './broadcast-message-general-form.vue';
import BroadcastMessageViewsForm from './broadcast-message-views-form.vue';

export default {
  inject: ['$validator'],
  components: {
    BroadcastMessageGeneralForm,
    BroadcastMessageViewsForm,
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
    treeItems: {
      type: Array,
      default: () => [],
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

<template>
  <c-form-general-patterns-tabs :patterns-label="$t('broadcastMessage.viewsAndPages')">
    <template #general="{ setRef }">
      <broadcast-message-general-form
        v-field="form"
        :ref="setRef"
      />
    </template>
    <template #patterns="{ setRef }">
      <broadcast-message-views-form
        v-field="form.views"
        :ref="setRef"
        :tree-items="treeItems"
      />
    </template>
  </c-form-general-patterns-tabs>
</template>

<script>
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
      type: Object,
      default: () => ({}),
    },
  },
};
</script>

<template>
  <associate-ticket-event-form
    :form="modelValue"
    :name="name"
    :depth="depth"
    @update:form="$emit('update:modelValue', $event)"
  />
</template>

<script>
import AssociateTicketEventForm from '@/components/other/declare-ticket/form/associate-ticket-event-form.vue';

export default {
  components: { AssociateTicketEventForm },

  props: {
    modelValue: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      default: 'parameters',
    },
    depth: {
      type: Number,
      default: 0,
    },
  },

  emits: ['update:modelValue'],
};
</script>

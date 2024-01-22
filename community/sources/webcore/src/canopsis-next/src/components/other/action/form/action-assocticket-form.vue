<template>
  <associate-ticket-event-form
    v-field="value"
    :name="name"
  />
</template>

<script>
import AssociateTicketEventForm from '@/components/other/declare-ticket/form/associate-ticket-event-form.vue';

export default {
  inject: ['$validator'],
  components: { AssociateTicketEventForm },
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      required: true,
    },
    name: {
      type: String,
      required: 'parameters',
    },
  },
};
</script>

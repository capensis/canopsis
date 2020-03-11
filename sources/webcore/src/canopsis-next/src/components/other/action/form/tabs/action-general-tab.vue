<template lang="pug">
  div
    component(:is="field.component", v-model="form[field.formKey]")
</template>

<script>
import formValidationHeaderMixin from '@/mixins/form/validation-header';

import Snooze from '@/components/other/action/form/fields/snooze.vue';
import Pbehavior from '@/components/other/action/form/fields/pbehavior.vue';
import ChangeState from '@/components/other/action/form/fields/change-state.vue';
import Note from '@/components/other/action/form/fields/note.vue';
import Assocticket from '@/components/other/action/form/fields/assocticket.vue';

export default {
  inject: ['$validator'],
  components: {
    Assocticket,
    Note,
    ChangeState,
    Pbehavior,
    Snooze,
  },
  mixins: [formValidationHeaderMixin],
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
  computed: {
    fieldsMap() {
      return ({
        [this.$constants.ACTION_TYPES.snooze]: {
          component: 'snooze',
          formKey: 'snoozeParameters',
        },
        [this.$constants.ACTION_TYPES.pbehavior]: {
          component: 'pbehavior',
          formKey: 'pbehaviorParameters',
        },
        [this.$constants.ACTION_TYPES.changeState]: {
          component: 'change-state',
          formKey: 'changeStateParameters',
        },
        [this.$constants.ACTION_TYPES.ack]: {
          component: 'note',
          formKey: 'ackParameters',
        },
        [this.$constants.ACTION_TYPES.ackremove]: {
          component: 'note',
          formKey: 'ackremoveParameters',
        },
        [this.$constants.ACTION_TYPES.assocticket]: {
          component: 'assocticket',
          formKey: 'assocticketParameters',
        },
        [this.$constants.ACTION_TYPES.declareticket]: {
          component: 'note',
          formKey: 'declareticketParameters',
        },
        [this.$constants.ACTION_TYPES.cancel]: {
          component: 'note',
          formKey: 'cancelParameters',
        },
      });
    },
    field() {
      return this.fieldsMap[this.form.generalParameters.type];
    },
  },
};
</script>

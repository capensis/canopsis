<template lang="pug">
  div
    component(
      v-model="form[$constants.ACTION_FORM_FIELDS_MAP_BY_TYPE[form.generalParameters.type]]",
      :is="fieldComponent"
    )
</template>

<script>
import formValidationHeaderMixin from '@/mixins/form/validation-header';

import { ACTION_TYPES } from '@/constants';

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
    fieldComponentsMap() {
      return {
        [ACTION_TYPES.snooze]: 'snooze',
        [ACTION_TYPES.pbehavior]: 'pbehavior',
        [ACTION_TYPES.changeState]: 'change-state',
        [ACTION_TYPES.assocticket]: 'assocticket',
        [ACTION_TYPES.ack]: 'note',
        [ACTION_TYPES.ackremove]: 'note',
        [ACTION_TYPES.declareticket]: 'note',
        [ACTION_TYPES.cancel]: 'note',
      };
    },
    fieldComponent() {
      return this.fieldComponentsMap[this.form.generalParameters.type];
    },
  },
};
</script>

<template lang="pug">
  v-form(@submit.prevent="submit")
    v-card
      v-card-title
        span.headline {{ $t('modals.createChangeStateEvent.title') }}
      v-card-text
        v-container
          v-layout(row)
            v-btn-toggle(
            v-model="form.state",
            v-validate="'required'",
            data-vv-name="state"
            )
              v-btn(
              v-for="button in buttons",
              :key="button",
              :value="states[button]",
              :color="colorsMap[states[button]]",
              depressed
              ) {{ $t(`modals.createChangeStateEvent.states.${button}`) }}
          v-layout(row)
            v-text-field(
            :label="$t('modals.createChangeStateEvent.output')",
            :error-messages="errors.collect('output')",
            v-model="form.output",
            v-validate="'required'",
            data-vv-name="output"
            )
      v-card-actions
        v-btn(type="submit", :disabled="errors.any()", color="primary") {{ $t('common.actions.saveChanges') }}
</template>

<script>
import ModalInnerItemMixin from '@/mixins/modal/modal-inner-item';
import EventActionsMixin from '@/mixins/event-actions';
import { EVENT_ENTITY_TYPES, ENTITY_STATES, MODALS } from '@/constants';


export default {
  name: MODALS.createChangeStateEvent,

  $_veeValidate: {
    validator: 'new',
  },
  mixins: [ModalInnerItemMixin, EventActionsMixin],
  data() {
    return {
      form: {
        output: '',
        state: ENTITY_STATES.ok,
      },
    };
  },
  computed: {
    buttons() {
      return Object.keys(ENTITY_STATES);
    },
    states() {
      return ENTITY_STATES;
    },
    colorsMap() {
      return {
        [ENTITY_STATES.ok]: 'info',
        [ENTITY_STATES.minor]: 'yellow',
        [ENTITY_STATES.major]: 'orange',
        [ENTITY_STATES.critical]: 'error',
      };
    },
  },
  mounted() {
    this.form.state = this.item.v.state.val;
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        await this.createEvent(EVENT_ENTITY_TYPES.changeState, this.item, this.form);

        this.hideModal();
      }
    },
  },
};
</script>

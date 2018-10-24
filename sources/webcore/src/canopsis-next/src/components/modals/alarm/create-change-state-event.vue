<template lang="pug">
  v-form(@submit.prevent="submit")
    v-card
      v-card-title.green.darken-3.white--text
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
          v-layout.mt-4(row)
            v-text-field(
            :label="$t('modals.createChangeStateEvent.fields.output')",
            :error-messages="errors.collect('output')",
            v-model="form.output",
            v-validate="'required'",
            data-vv-name="output"
            )
      v-card-actions
        v-btn.green.darken-3.white--text(type="submit", :disabled="errors.any()") {{ $t('common.actions.saveChanges') }}
</template>

<script>
import modalInnerItemsMixin from '@/mixins/modal/modal-inner-items';
import eventActionsMixin from '@/mixins/event-actions';
import { EVENT_ENTITY_TYPES, ENTITIES_STATES, MODALS } from '@/constants';

/**
 * Modal to create a 'change-state' event
 */
export default {
  name: MODALS.createChangeStateEvent,

  $_veeValidate: {
    validator: 'new',
  },
  mixins: [modalInnerItemsMixin, eventActionsMixin],
  data() {
    return {
      form: {
        output: '',
        state: ENTITIES_STATES.ok,
      },
    };
  },
  computed: {
    buttons() {
      return Object.keys(ENTITIES_STATES);
    },
    states() {
      return ENTITIES_STATES;
    },
    colorsMap() {
      return {
        [ENTITIES_STATES.ok]: 'green',
        [ENTITIES_STATES.minor]: 'yellow',
        [ENTITIES_STATES.major]: 'orange',
        [ENTITIES_STATES.critical]: 'error',
      };
    },
  },
  mounted() {
    this.form.state = this.firstItem.v.state.val;
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        await this.createEvent(EVENT_ENTITY_TYPES.changeState, this.items, this.form);

        this.hideModal();
      }
    },
  },
};
</script>

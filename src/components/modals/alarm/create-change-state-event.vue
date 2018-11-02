<template lang="pug">
  v-form(@submit.prevent="submit")
    v-card
      v-card-title.primary.white--text
        v-layout(justify-space-between, align-center)
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
      v-divider
      v-layout.py-1(justify-end)
        v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
        v-btn.primary(type="submit", :disabled="errors.any()") {{ $t('common.actions.saveChanges') }}
</template>

<script>
import modalInnerItemsMixin from '@/mixins/modal/modal-inner-items';
import eventActionsMixin from '@/mixins/event-actions';
import { MODALS } from '@/constants';

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
        state: this.$constants.ENTITIES_STATES.ok,
      },
    };
  },
  computed: {
    buttons() {
      return Object.keys(this.$constants.ENTITIES_STATES);
    },
    states() {
      return this.$constants.ENTITIES_STATES;
    },
    colorsMap() {
      return {
        [this.$constants.ENTITIES_STATES.ok]: 'green',
        [this.$constants.ENTITIES_STATES.minor]: 'yellow',
        [this.$constants.ENTITIES_STATES.major]: 'orange',
        [this.$constants.ENTITIES_STATES.critical]: 'error',
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
        await this.createEvent(this.$constants.EVENT_ENTITY_TYPES.changeState, this.items, this.form);

        this.hideModal();
      }
    },
  },
};
</script>

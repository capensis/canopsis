<template lang="pug">
  v-form(@submit.prevent="submit")
    v-card
      v-card-title
        span.headline {{ $t('modals.createChangeStateEvent.title') }}
      v-card-text
        v-container
          v-layout(row)
            v-btn-toggle(v-model="state")
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
import { createNamespacedHelpers } from 'vuex';

import ModalItemMixin from '@/mixins/modal-item';
import EventEntityMixin from '@/mixins/event-entity';
import { STATES } from '@/config';

const { mapActions } = createNamespacedHelpers('event');

export default {
  $_veeValidate: {
    validator: 'new',
  },
  mixins: [ModalItemMixin, EventEntityMixin],
  data() {
    return {
      form: {
        output: '',
      },
    };
  },
  computed: {
    state: {
      get() {
        return this.item.v.state.val;
      },
      set(value) {
        this.item.v.state.val = value;
      },
    },
    buttons() {
      return Object.keys(STATES);
    },
    states() {
      return STATES;
    },
    colorsMap() {
      return {
        [STATES.info]: 'info',
        [STATES.minor]: 'yellow',
        [STATES.major]: 'orange',
        [STATES.critical]: 'error',
      };
    },
  },
  methods: {
    ...mapActions([
      'changeState',
    ]),
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        this.changeState({
          id: 'a7ed1556-4eda-11e8-841e-0242ac12000a',
          resource: 'res93',
          state: this.form.state,
          customAttributes: {
            output: this.form.output,
          },
        });
        //  todo hide modal
      }
    },
  },
};
</script>

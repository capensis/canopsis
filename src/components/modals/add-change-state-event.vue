<template lang="pug">
  v-form(@submit.prevent="submit")
    v-card
      v-card-title
        span.headline {{ $t('modals.addChangeStateEvent.title') }}
      v-card-text
        v-container
          v-layout(row)
            v-btn-toggle(v-model="form.state")
              v-btn(color="success", value="info", depressed) {{ $t('modals.addChangeStateEvent.states.info') }}
              v-btn(color="warning", value="minor", depressed) {{ $t('modals.addChangeStateEvent.states.minor') }}
              v-btn(color="error", value="critical", depressed) {{ $t('modals.addChangeStateEvent.states.critical') }}
          v-layout(row)
            v-text-field(
            :label="$t('modals.addChangeStateEvent.output')",
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

const { mapActions } = createNamespacedHelpers('event');

export default {
  $_veeValidate: {
    validator: 'new',
  },
  data() {
    return {
      form: {
        state: 'info',
        output: '',
      },
    };
  },
  methods: {
    ...mapActions([
      'changeState',
    ]),
    getStateId() {
      switch (this.form.state) {
        case 'info':
          return 0;
        case 'minor':
          return 1;
        case 'critical':
          return 3;
        default: {
          throw new Error('Unknown alarm state');
        }
      }
    },
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        this.changeState({
          id: 'a7ed1556-4eda-11e8-841e-0242ac12000a',
          resource: 'res93',
          state: this.getStateId(this.form.state),
          customAttributes: {
            output: this.form.output,
          },
        });
      }
    },
  },
};
</script>

<template lang="pug">
  v-dialog(:value="opened", @input="hideModal", max-width="700")
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
          v-btn(type="submit") {{ $t('common.actions.saveChanges') }}
</template>

<script>
import ModalHOC from './modal-hoc';

export default ModalHOC({
  name: 'add-change-state-event',
  data() {
    return {
      form: {
        state: 'info',
        output: '',
      },
    };
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        console.log('SUBMITTED');
      }
    },
  },
});
</script>

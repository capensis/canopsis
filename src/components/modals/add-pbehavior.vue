<template lang="pug">
  v-dialog(:value="opened", @input="hideModal", max-width="700")
    v-form(@submit.prevent="submit")
      v-card
        v-card-title
          span.headline {{ $t('modals.addChangeStateEvent.title') }}
        v-card-text
          v-container
            v-layout(row)
              v-text-field(
              :label="$t('modals.addChangeStateEvent.output')",
              :error-messages="errors.collect('output')",
              v-model="form.output",
              v-validate="'required'",
              data-vv-name="output"
              )
            v-layout(row)
              v-menu(
              ref="menu2",
              :close-on-content-click="false",
              transition="scale-transition",
              offset-y,
              full-width,
              :nudge-right="40",
              min-width="290px",
              :return-value.sync="form.date"
              )
                v-text-field(
                readonly,
                slot="activator",
                :label="$t('modals.addChangeStateEvent.output')",
                :error-messages="errors.collect('date')",
                v-model="form.date",
                v-validate="'required'",
                data-vv-name="date"
                )
                v-date-picker(
                v-model="form.date",
                @input="$refs.menu2.save(form.date)"
                )
        v-card-actions
          v-btn(type="submit") {{ $t('common.actions.saveChanges') }}
</template>

<script>
import ModalHOC from './modal-hoc';

export default ModalHOC({
  name: 'add-pbehavior',
  data() {
    return {
      form: {
        name: '',
        date: null,
      },
    };
  },
});
</script>

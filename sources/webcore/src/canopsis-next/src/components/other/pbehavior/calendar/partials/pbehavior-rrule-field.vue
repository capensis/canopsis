<template lang="pug">
  v-flex
    v-btn.ml-0(
      color="primary",
      @click="showCreateRRuleModal"
    )  {{ hasRRule ? $t('pbehavior.buttons.editRrule') : $t('pbehavior.buttons.addRRule') }}
    v-tooltip(v-show="hasRRule", fixed, top)
      v-btn(slot="activator", icon)
        v-icon(color="grey darken-1") info
      span {{ form.rrule }}
</template>

<script>
import { isEmpty } from 'lodash';

import { MODALS } from '@/constants';

import formMixin from '@/mixins/form/object';

export default {
  mixins: [formMixin],
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
    hasRRule() {
      return !isEmpty(this.form.rrule);
    },
  },
  methods: {
    showCreateRRuleModal() {
      this.$modals.show({
        name: MODALS.createRRule,
        config: {
          rrule: this.form.rrule,
          exdates: this.form.exdates,
          exceptions: this.form.exceptions,
          action: ({ rrule, exdates, exceptions }) => this.updateModel({
            ...this.form,
            rrule,
            exdates,
            exceptions,
          }),
        },
      });
    },
  },
};
</script>

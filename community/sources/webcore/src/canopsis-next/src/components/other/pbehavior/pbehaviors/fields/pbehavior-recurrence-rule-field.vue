<template>
  <v-flex>
    <v-btn
      class="ml-0"
      color="primary"
      @click="showCreateRecurrenceRuleModal"
    >
      {{ hasRecurrenceRule ? $t('pbehavior.buttons.editRrule') : $t('pbehavior.buttons.addRRule') }}
    </v-btn>
    <template v-if="hasRecurrenceRule">
      <v-tooltip
        fixed
        top
      >
        <template #activator="{ on }">
          <v-btn
            v-on="on"
            icon
          >
            <v-icon color="grey darken-1">
              info
            </v-icon>
          </v-btn>
        </template>
        <span>{{ form.rrule }}</span>
      </v-tooltip>
      <c-action-btn
        type="delete"
        @click="showConfirmRemoveRecurrenceRuleModal"
      />
    </template>
  </v-flex>
</template>

<script>
import { isEmpty } from 'lodash';

import { MODALS } from '@/constants';

import { formMixin } from '@/mixins/form';

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
    withExdateType: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    hasRecurrenceRule() {
      return !isEmpty(this.form.rrule);
    },
  },
  methods: {
    showConfirmRemoveRecurrenceRuleModal() {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: () => this.updateField('rrule', ''),
        },
      });
    },

    showCreateRecurrenceRuleModal() {
      this.$modals.show({
        name: MODALS.createRecurrenceRule,
        config: {
          rrule: this.form.rrule,
          exdates: this.form.exdates,
          exceptions: this.form.exceptions,
          start: this.form.tstart,
          withExdateType: this.withExdateType,
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

<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.eventFilterRule.editActions') }}
    v-card-text
      v-card
        v-card-title.primary.white--text {{ $t('modals.eventFilterRule.addAction') }}
        v-card-text
          v-form(ref="form")
            v-select(
            :items="Object.values($constants.EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES)",
            v-model="form.type",
            return-object,
            item-text="value",
            :label="$t('common.type')"
            )
            component(
            v-for="option in form.type.options",
            :is="getComponentByOption(option)",
            :key="option.value",
            v-model="form[option.value]",
            :label="option.text",
            :name="option.value",
            :error-messages="errors.collect(option.value)",
            v-validate="getValidationRulesByOption(option)"
            )
          v-btn.primary(@click="addAction") {{ $t('common.add') }}
    v-container
      h2 {{ $t('modals.eventFilterRule.actions') }}
      v-list(dark)
        draggable(v-model="actions")
          v-list-group(v-for="(action, index) in actions", :key="action.name")
            v-list-tile(slot="activator")
              v-list-tile-title {{index + 1}} - {{ action.type }} - {{ action.name || action.from }}
              v-btn(@click.stop="deleteAction(index)", icon)
                v-icon(color="error") delete
            v-list-tile
              v-layout(column)
                div(v-if="action.name") {{ $t('common.name') }}: {{ action.name }}
                div(v-if="action.value") {{ $t('common.value') }}: {{ action.value }}
                div(v-if="action.description") {{ $t('common.description') }}: {{ action.description }}
                div(v-if="action.from") {{ $t('common.from') }}: {{ action.from }}
                div(v-if="action.to") {{ $t('common.to') }}: {{ action.to }}
    v-divider
    v-layout.pa-2(justify-end)
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary(@click.prevent="submit") {{ $t('common.submit') }}
</template>

<script>
import Draggable from 'vuedraggable';
import { cloneDeep, pick } from 'lodash';

import { MODALS, EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

import MixedField from '@/components/forms/fields/mixed-field.vue';

export default {
  name: MODALS.eventFilterRuleActions,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    Draggable,
    MixedField,
  },
  mixins: [modalInnerMixin],
  data() {
    const enrichmentActionsTypes = cloneDeep(EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES);

    return {
      actions: [],
      form: {
        type: enrichmentActionsTypes.setField,
        name: '',
        value: '',
        description: '',
        from: '',
        to: '',
      },
    };
  },
  computed: {
    getComponentByOption() {
      return (option = {}) => (option.value === 'value' ? 'mixed-field' : 'v-text-field');
    },
    getValidationRulesByOption() {
      return (option = {}) => option.required && 'required';
    },
  },
  mounted() {
    if (this.config.actions) {
      this.actions = [...this.config.actions];
    }
  },
  methods: {
    isRequired(actionType, option) {
      return actionType.options[option.value].required;
    },

    async addAction() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        const action = {
          type: this.form.type.value,
          ...pick(this.form, Object.keys(this.form.type.options)),
        };

        this.actions.push(action);
      }
    },

    deleteAction(index) {
      this.$delete(this.actions, index);
    },

    submit() {
      this.config.action(this.actions);
      this.hideModal();
    },
  },
};
</script>

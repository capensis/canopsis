<template lang="pug">
  modal-wrapper(close)
    template(slot="title")
      span {{ $t('modals.eventFilterRule.editActions') }}
    template(slot="text")
      v-layout(justify-end)
        v-tooltip(top)
          v-btn(slot="activator", icon, @click="showCreateActionModal")
            v-icon.primary--text add
          span {{ $t('common.add') }}
      v-container
        h2 {{ $t('modals.eventFilterRule.actions') }}
        v-list(dark)
          draggable(v-model="form.actions")
            v-list-group(v-for="(action, index) in form.actions", :key="action.name")
              v-list-tile(slot="activator")
                v-list-tile-title {{index + 1}} - {{ action.type }} - {{ action.name || action.from }}
                v-btn(@click.stop="showEditActionModal(index)", icon)
                  v-icon(color="success") edit
                v-btn(@click.stop="deleteAction(index)", icon)
                  v-icon(color="error") delete
              v-list-tile
                v-layout(column)
                  div(v-if="action.name") {{ $t('common.name') }}: {{ action.name }}
                  div(v-if="action.value") {{ $t('common.value') }}: {{ action.value }}
                  div(v-if="action.description") {{ $t('common.description') }}: {{ action.description }}
                  div(v-if="action.from") {{ $t('common.from') }}: {{ action.from }}
                  div(v-if="action.to") {{ $t('common.to') }}: {{ action.to }}
    template(slot="actions")
      v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
      v-btn.primary(@click.prevent="submit") {{ $t('common.submit') }}
</template>

<script>
import { cloneDeep } from 'lodash';
import Draggable from 'vuedraggable';

import { MODALS } from '@/constants';

import submittableMixin from '@/mixins/submittable';
import confirmableModalMixin from '@/mixins/confirmable-modal';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.eventFilterRuleActions,
  $_veeValidate: {
    validator: 'new',
  },
  components: { Draggable, ModalWrapper },
  mixins: [
    submittableMixin(),
    confirmableModalMixin(),
  ],
  data() {
    const { actions = [] } = this.modal.config;

    return {
      form: {
        actions: cloneDeep(actions),
      },
    };
  },
  methods: {
    showCreateActionModal() {
      this.$modals.show({
        name: MODALS.eventFilterRuleCreateAction,
        config: {
          action: ruleAction => this.form.actions.push(ruleAction),
        },
      });
    },

    showEditActionModal(index) {
      this.$modals.show({
        name: MODALS.eventFilterRuleCreateAction,
        config: {
          ruleAction: this.form.actions[index],
          action: ruleAction => this.$set(this.form.actions, index, ruleAction),
        },
      });
    },

    deleteAction(index) {
      this.$delete(this.form.actions, index);
    },

    async submit() {
      await this.config.action(this.form.actions);
      this.$modals.hide();
    },
  },
};
</script>
